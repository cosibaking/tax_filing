#!/usr/bin/env bash
#
# 线上 Docker 部署脚本（在服务器项目根目录执行）
# 用法:
#   ./scripts/deploy-prod.sh              # 构建并启动（80 被占用时自动回退 8048）
#   ./scripts/deploy-prod.sh --pull       # 先 git pull 再部署
#   ./scripts/deploy-prod.sh --no-build   # 仅 up -d，不重建镜像
#   ./scripts/deploy-prod.sh --restart    # 重启全部服务
#   ./scripts/deploy-prod.sh --port 8048  # 强制指定端口（备份方案）
#
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILES=(-f docker-compose.yml -f docker-compose.prod.yml)
LEGACY_DIR="${LEGACY_STACK_DIR:-$HOME/apps/tax_filling}"
HEALTH_PATH="/tax_filing/api/health"
HEALTH_TIMEOUT_SEC="${HEALTH_TIMEOUT_SEC:-300}"
PULL_BRANCH="${DEPLOY_GIT_BRANCH:-}"
PRIMARY_PORT="${NGINX_PRIMARY_PORT:-80}"
FALLBACK_PORT="${NGINX_FALLBACK_PORT:-8048}"

DO_PULL=false
DO_BUILD=true
DO_RESTART=false
FORCE_PORT=""
USING_FALLBACK=false
NEEDS_APP_REBUILD=false

log() { printf '[deploy] %s\n' "$*"; }
warn() { printf '[deploy][WARN] %s\n' "$*" >&2; }
die() { printf '[deploy][ERROR] %s\n' "$*" >&2; exit 1; }

usage() {
  cat <<'EOF'
线上 Docker 部署脚本（在服务器项目根目录执行）
用法:
  ./scripts/deploy-prod.sh              # 构建并启动（80 被占用时自动回退 8048）
  ./scripts/deploy-prod.sh --pull       # 先 git pull 再部署
  ./scripts/deploy-prod.sh --no-build   # 仅 up -d，不重建镜像
  ./scripts/deploy-prod.sh --restart    # 重启全部服务
  ./scripts/deploy-prod.sh --port 8048  # 强制使用备份端口
EOF
  exit "${1:-0}"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --pull) DO_PULL=true ;;
    --no-build) DO_BUILD=false ;;
    --restart) DO_RESTART=true; DO_BUILD=false ;;
    --port)
      [[ $# -ge 2 ]] || die "--port 需要端口号，例如: --port 8048"
      FORCE_PORT="$2"
      shift
      ;;
    -h|--help) usage 0 ;;
    *) die "未知参数: $1（使用 --help 查看用法）" ;;
  esac
  shift
done

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "缺少命令: $1"
}

detect_compose() {
  if docker compose version >/dev/null 2>&1; then
    COMPOSE_BIN=(docker compose)
  elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE_BIN=(docker-compose)
  else
    die "未找到 docker compose，请先安装 Docker Compose v2"
  fi
}

detect_docker() {
  if docker info >/dev/null 2>&1; then
    COMPOSE=("${COMPOSE_BIN[@]}")
    return
  fi
  if sudo docker info >/dev/null 2>&1; then
    COMPOSE=(sudo "${COMPOSE_BIN[@]}")
    warn "当前用户无 docker 权限，将使用 sudo"
    return
  fi
  die "无法连接 Docker 守护进程，请确认 Docker 已启动且当前用户有权限"
}

load_env() {
  [[ -f .env ]] || die "缺少 .env，请先执行: cp .env.docker.example .env 并编辑"

  if grep -q $'\r' .env 2>/dev/null; then
    warn ".env 含 Windows 换行符(CRLF)，正在自动转换为 LF..."
    sed -i 's/\r$//' .env
  fi

  set -a
  while IFS= read -r line || [[ -n "$line" ]]; do
    line="${line#"${line%%[![:space:]]*}"}"
    line="${line%"${line##*[![:space:]]}"}"
    [[ -z "$line" || "$line" == \#* ]] && continue

    if [[ "$line" =~ ^([A-Za-z_][A-Za-z0-9_]*)=(.*)$ ]]; then
      key="${BASH_REMATCH[1]}"
      val="${BASH_REMATCH[2]}"
      if [[ "$val" =~ ^\"(.*)\"$ ]]; then
        val="${BASH_REMATCH[1]}"
      elif [[ "$val" =~ ^\'(.*)\'$ ]]; then
        val="${BASH_REMATCH[1]}"
      fi
      printf -v "$key" '%s' "$val"
      export "$key"
    fi
  done < .env
  set +a
}

validate_env() {
  local required=(
    MYSQL_ROOT_PASSWORD
    MYSQL_PASSWORD
    MYSQL_USER
    MYSQL_DATABASE
    NEXTAUTH_SECRET
    JWT_SECRET
    ENCRYPTION_KEY
    CRON_SECRET
    NEXT_PUBLIC_APP_URL
  )
  local missing=()
  local key
  for key in "${required[@]}"; do
    [[ -n "${!key:-}" ]] || missing+=("$key")
  done

  if [[ ${#missing[@]} -gt 0 ]]; then
    cat >&2 <<EOF
[deploy][ERROR] .env 缺少 Docker 部署必填项: ${missing[*]}

线上 .env 须包含 MySQL 与 NEXT_PUBLIC_APP_URL 等 Docker 变量（与本地开发的 DATABASE_URL 格式不同）。

修复方式（在 ~/apps/tax_filing）:
  1. 若数据库已在运行：编辑 .env，补全缺失项，勿随意改 MYSQL_PASSWORD
  2. 若是全新环境：cp .env.docker.example .env 后编辑

示例片段:
  MYSQL_ROOT_PASSWORD=your-root-password
  MYSQL_DATABASE=opc_compliance
  MYSQL_USER=tax_filing_app
  MYSQL_PASSWORD=your-app-password
  NEXT_PUBLIC_APP_URL=http://82.156.54.232/tax_filing
  NEXTAUTH_SECRET=...（至少 32 字符）
  JWT_SECRET=...
  ENCRYPTION_KEY=...
  CRON_SECRET=...
EOF
    exit 1
  fi

  if [[ "${MYSQL_PASSWORD}" =~ [,@#/:] ]]; then
    die "MYSQL_PASSWORD 含 , @ # / : 等特殊字符，会导致 DATABASE_URL 解析失败"
  fi

  if [[ "${NEXT_PUBLIC_APP_URL}" != */tax_filing ]]; then
    warn "NEXT_PUBLIC_APP_URL 建议以 /tax_filing 结尾，当前: ${NEXT_PUBLIC_APP_URL}"
  fi
}

is_port_listening() {
  local port="$1"
  if command -v ss >/dev/null 2>&1; then
    ss -tlnH "sport = :${port}" 2>/dev/null | grep -q .
    return
  fi
  if command -v netstat >/dev/null 2>&1; then
    netstat -tln 2>/dev/null | awk -v p=":${port} " '$4 ~ p { found=1 } END { exit !found }'
    return
  fi
  # 无法检测时假定可用，交给 docker compose 报错
  return 1
}

is_our_nginx_on_port() {
  local port="$1"
  local mapping="${port}->80/tcp"
  "${COMPOSE[@]}" "${COMPOSE_FILES[@]}" ps nginx --format '{{.Ports}}' 2>/dev/null \
    | grep -q "$mapping"
}

set_env_file_var() {
  local key="$1"
  local val="$2"
  local escaped="${val//\\/\\\\}"
  escaped="${escaped//|/\\|}"

  if grep -q "^${key}=" .env 2>/dev/null; then
    sed -i "s|^${key}=.*|${key}=${escaped}|" .env
  else
    printf '%s=%s\n' "$key" "$val" >> .env
  fi

  printf -v "$key" '%s' "$val"
  export "$key"
}

build_app_url_for_port() {
  local port="$1"
  local url="${NEXT_PUBLIC_APP_URL:-http://127.0.0.1/tax_filing}"
  local scheme="" host="" path="/tax_filing"

  if [[ "$url" =~ ^(https?)://([^/:]+)(:[0-9]+)?(/.*)$ ]]; then
    scheme="${BASH_REMATCH[1]}"
    host="${BASH_REMATCH[2]}"
    path="${BASH_REMATCH[4]}"
  else
    scheme="http"
    host="127.0.0.1"
  fi

  if [[ "$port" == "80" ]]; then
    printf '%s://%s%s' "$scheme" "$host" "$path"
  else
    printf '%s://%s:%s%s' "$scheme" "$host" "$port" "$path"
  fi
}

resolve_listen_port() {
  local chosen_port=""

  if [[ -n "$FORCE_PORT" ]]; then
    chosen_port="$FORCE_PORT"
    if [[ "$chosen_port" != "$PRIMARY_PORT" ]]; then
      USING_FALLBACK=true
      warn "已强制使用端口 ${chosen_port}"
    fi
  elif [[ -n "${NGINX_HTTP_PORT:-}" && "${NGINX_HTTP_PORT}" != "$PRIMARY_PORT" ]]; then
    chosen_port="${NGINX_HTTP_PORT}"
    USING_FALLBACK=true
    log ".env 已指定 NGINX_HTTP_PORT=${chosen_port}，使用备份端口"
  elif is_port_listening "$PRIMARY_PORT" && ! is_our_nginx_on_port "$PRIMARY_PORT"; then
    chosen_port="$FALLBACK_PORT"
    USING_FALLBACK=true
    warn "端口 ${PRIMARY_PORT} 已被占用，自动切换备份端口 ${FALLBACK_PORT}"
  else
    chosen_port="$PRIMARY_PORT"
  fi

  if is_port_listening "$chosen_port" && ! is_our_nginx_on_port "$chosen_port"; then
    die "端口 ${chosen_port} 也被占用，请释放端口或指定其它端口: ./scripts/deploy-prod.sh --port <port>"
  fi

  local new_app_url
  new_app_url="$(build_app_url_for_port "$chosen_port")"

  if [[ "${NGINX_HTTP_PORT:-}" != "$chosen_port" ]]; then
    set_env_file_var NGINX_HTTP_PORT "$chosen_port"
  fi

  if [[ "${NEXT_PUBLIC_APP_URL}" != "$new_app_url" ]]; then
    warn "同步 NEXT_PUBLIC_APP_URL -> ${new_app_url}"
    set_env_file_var NEXT_PUBLIC_APP_URL "$new_app_url"
    NEEDS_APP_REBUILD=true
  fi

  export NGINX_HTTP_PORT="$chosen_port"
  export NEXT_PUBLIC_APP_URL="$new_app_url"
}

remove_iptables_redirect() {
  local removed=0

  if ! command -v iptables >/dev/null 2>&1; then
    log "跳过 iptables 检查（未安装 iptables）"
    return
  fi

  local sudo_cmd=()
  if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
    sudo_cmd=(sudo)
  fi

  delete_rule() {
    local table="$1"
    shift
    while "${sudo_cmd[@]}" iptables -t nat -C "$table" "$@" 2>/dev/null; do
      "${sudo_cmd[@]}" iptables -t nat -D "$table" "$@"
      removed=$((removed + 1))
      log "已删除 iptables NAT $table: 80 -> 4096 重定向"
    done
  }

  delete_rule PREROUTING -p tcp --dport 80 -j REDIRECT --to-ports 4096
  delete_rule OUTPUT -p tcp --dport 80 -j REDIRECT --to-ports 4096
  delete_rule OUTPUT -p tcp -o lo --dport 80 -j REDIRECT --to-ports 4096

  if [[ "$removed" -eq 0 ]]; then
    log "未发现 iptables 80->4096 重定向规则"
  fi
}

stop_legacy_stack() {
  if [[ ! -d "$LEGACY_DIR" ]] || [[ ! -f "$LEGACY_DIR/docker-compose.yml" ]]; then
    return
  fi

  log "检查旧栈目录: $LEGACY_DIR"

  if (cd "$LEGACY_DIR" && "${COMPOSE[@]}" ps -q 2>/dev/null | grep -q .); then
    warn "检测到旧栈 tax_filling 仍在运行，正在停止..."
    (cd "$LEGACY_DIR" && "${COMPOSE[@]}" down) || warn "停止旧栈失败，请手动检查 $LEGACY_DIR"
  else
    log "旧栈未运行，跳过"
  fi
}

git_pull() {
  require_cmd git
  local branch="${PULL_BRANCH:-$(git rev-parse --abbrev-ref HEAD 2>/dev/null || true)}"
  [[ -n "$branch" ]] || die "无法确定 git 分支，请设置 DEPLOY_GIT_BRANCH 或手动 pull"

  log "git pull origin $branch"
  git pull origin "$branch"
}

compose_up() {
  if $DO_RESTART; then
    log "重启服务..."
    "${COMPOSE[@]}" "${COMPOSE_FILES[@]}" restart
    return
  fi

  local args=(up -d)
  if $DO_BUILD || $NEEDS_APP_REBUILD; then
    args+=(--build)
  fi

  log "启动 Docker Compose（prod，端口 ${NGINX_HTTP_PORT}）: ${args[*]}"
  NGINX_HTTP_PORT="${NGINX_HTTP_PORT}" NEXT_PUBLIC_APP_URL="${NEXT_PUBLIC_APP_URL}" \
    "${COMPOSE[@]}" "${COMPOSE_FILES[@]}" "${args[@]}"
}

wait_for_health() {
  local port="${NGINX_HTTP_PORT:-80}"
  local url="http://127.0.0.1:${port}${HEALTH_PATH}"
  local deadline=$((SECONDS + HEALTH_TIMEOUT_SEC))

  log "等待健康检查: $url （最长 ${HEALTH_TIMEOUT_SEC}s）"

  while (( SECONDS < deadline )); do
    if curl -sf --connect-timeout 3 "$url" >/dev/null 2>&1; then
      log "健康检查通过"
      curl -sf "$url" || true
      echo
      return 0
    fi

    if ! "${COMPOSE[@]}" "${COMPOSE_FILES[@]}" ps --status running 2>/dev/null | grep -q app; then
      warn "app 容器未处于 running，最近日志:"
      "${COMPOSE[@]}" "${COMPOSE_FILES[@]}" logs app --tail 30 || true
    fi

    sleep 5
  done

  die "健康检查超时（${HEALTH_TIMEOUT_SEC}s）。请执行: ${COMPOSE[*]} ${COMPOSE_FILES[*]} logs app"
}

print_summary() {
  local base_url="${NEXT_PUBLIC_APP_URL}"

  log "容器状态:"
  "${COMPOSE[@]}" "${COMPOSE_FILES[@]}" ps

  cat <<EOF

部署完成（监听端口: ${NGINX_HTTP_PORT}）。
  首页:     ${base_url}
  健康检查: ${base_url%/}/api/health
  管理后台: ${base_url%/}/admin/login
EOF

  if $USING_FALLBACK; then
    cat <<EOF

[备份端口模式]
  当前未使用 80 端口，请确认云安全组/防火墙已放行 TCP ${NGINX_HTTP_PORT}。
  恢复 80 端口：释放占用后编辑 .env 设 NGINX_HTTP_PORT=80 并更新 NEXT_PUBLIC_APP_URL，再执行本脚本。
EOF
  fi

  cat <<EOF

常用命令:
  ${COMPOSE[*]} ${COMPOSE_FILES[*]} logs app --tail 100
  ${COMPOSE[*]} ${COMPOSE_FILES[*]} ps
EOF
}

main() {
  require_cmd curl
  detect_compose
  detect_docker

  load_env
  validate_env

  stop_legacy_stack
  resolve_listen_port

  if $DO_PULL; then
    git_pull
  fi

  if ! $USING_FALLBACK && [[ "${NGINX_HTTP_PORT}" == "$PRIMARY_PORT" ]]; then
    remove_iptables_redirect
  elif $USING_FALLBACK; then
    log "备份端口模式，跳过 iptables 80 重定向检查"
  fi

  compose_up
  wait_for_health
  print_summary
}

main "$@"
