#!/usr/bin/env bash
# 从零初始化本地开发环境
# 用法:
#   ./scripts/init.sh
#   MYSQL_ROOT_PASSWORD=xxx ./scripts/init.sh
#   ./scripts/init.sh --skip-database --skip-migrate
#   ./scripts/init.sh --start

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SERVER_DIR="$ROOT/server"
WEB_DIR="$ROOT/web"
SKIP_DATABASE=false
SKIP_MIGRATE=false
FORCE_IMPORT=false
FRESH_CONFIG=false
START=false
ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --skip-database) SKIP_DATABASE=true ;;
    --skip-migrate) SKIP_MIGRATE=true ;;
    --force-import) FORCE_IMPORT=true ;;
    --fresh-config) FRESH_CONFIG=true ;;
    --start) START=true ;;
    --root-password) ROOT_PASSWORD="$2"; shift ;;
    *) echo "未知参数: $1"; exit 1 ;;
  esac
  shift
done

step() { echo ">> $1"; }
ok() { echo "OK $1"; }
warn() { echo "!! $1"; }

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || { echo "缺少命令: $1"; exit 1; }
}

echo ""
echo "========== XYGo 本地环境初始化 =========="
echo ""

step "检查基础依赖..."
require_cmd go
require_cmd node
require_cmd pnpm
ok "Go / Node / pnpm 已就绪"

step "准备配置文件..."
SERVER_CONFIG="$SERVER_DIR/manifest/config/config.yaml"
SERVER_EXAMPLE="$SERVER_DIR/manifest/config/config.yaml.example"
WEB_ENV_LOCAL="$WEB_DIR/.env.local"
WEB_ENV_DEV="$WEB_DIR/.env.development"

if $FRESH_CONFIG && [[ -f "$SERVER_CONFIG" ]]; then
  rm -f "$SERVER_CONFIG"
fi
if [[ ! -f "$SERVER_CONFIG" ]]; then
  cp "$SERVER_EXAMPLE" "$SERVER_CONFIG"
  ok "已生成 server/manifest/config/config.yaml"
fi
if [[ ! -f "$WEB_ENV_LOCAL" ]]; then
  cp "$WEB_ENV_DEV" "$WEB_ENV_LOCAL"
  ok "已生成 web/.env.local"
fi

HOST="${MYSQL_HOST:-127.0.0.1}"
PORT="${MYSQL_PORT:-3306}"
DATABASE="${MYSQL_DATABASE:-xygo}"
APP_USER="${MYSQL_APP_USER:-yxgo}"
APP_PASSWORD="${MYSQL_APP_PASSWORD:-123456}"

if ! $SKIP_DATABASE; then
  if [[ -z "$ROOT_PASSWORD" ]]; then
    read -rsp "MySQL root 密码: " ROOT_PASSWORD
    echo
  fi
  IMPORT_FLAG=""
  $FORCE_IMPORT && IMPORT_FLAG="--force"
  MYSQL_HOST="$HOST" MYSQL_PORT="$PORT" MYSQL_DATABASE="$DATABASE" \
    MYSQL_APP_USER="$APP_USER" MYSQL_APP_PASSWORD="$APP_PASSWORD" \
    MYSQL_ROOT_PASSWORD="$ROOT_PASSWORD" \
    bash "$ROOT/scripts/import-mysql.sh" --import $IMPORT_FLAG
else
  warn "已跳过数据库初始化 (--skip-database)"
fi

step "写入开发配置..."
LINK="mysql:${APP_USER}:${APP_PASSWORD}@tcp(${HOST}:${PORT})/${DATABASE}?loc=Local&parseTime=true"
if command -v python3 >/dev/null 2>&1; then PY=python3; else PY=python; fi
"$PY" - "$SERVER_CONFIG" "$LINK" <<'PY'
import re, sys
path, link = sys.argv[1], sys.argv[2]
text = open(path, encoding='utf-8').read()
text = re.sub(r'(?m)^(\s*link:\s*")[^"]*(")', rf'\1{link}\2', text, count=1)
text = re.sub(r'(?m)^(\s*adapter:\s*")[^"]*(")', r'\1memory\2', text, count=1)
text = re.sub(r'(?m)^(\s*driver:\s*")[^"]*(")', r'\1disk\2', text, count=1)
open(path, 'w', encoding='utf-8').write(text)
PY
ok "开发配置已更新"

step "安装 GoFrame CLI..."
if ! command -v gf >/dev/null 2>&1; then
  go install github.com/gogf/gf/cmd/gf/v2@latest
fi
ok "gf 已就绪"

step "下载后端依赖..."
( cd "$SERVER_DIR" && go mod download )
ok "后端依赖就绪"

if [[ ! -d "$WEB_DIR/node_modules" ]]; then
  step "安装前端依赖..."
  ( cd "$WEB_DIR" && pnpm install )
  ok "前端依赖安装完成"
else
  warn "node_modules 已存在，跳过 pnpm install"
fi

if ! $SKIP_MIGRATE; then
  step "执行数据库迁移..."
  ( cd "$SERVER_DIR" && go run tools.go migrate up )
  ok "数据库迁移完成"
else
  warn "已跳过数据库迁移 (--skip-migrate)"
fi

echo ""
echo "========== 初始化完成 =========="
echo ""
echo "默认账号:"
echo "  管理端  Super / 123456  -> http://localhost:5173/admin"
echo "  会员端  自行注册        -> http://localhost:5173/user/register"
echo ""
echo "启动开发服务:"
echo "  ./start.sh"
echo ""

if $START; then
  step "启动开发服务..."
  bash "$ROOT/start.sh"
fi
