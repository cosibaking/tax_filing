#!/usr/bin/env bash
# XYGo Admin 本地开发启动脚本
# 用法: ./start.sh [--init] [--migrate] [--restart] [--stop] [--backend-only] [--frontend-only] [--help]

set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
SERVER_DIR="$ROOT/server"
WEB_DIR="$ROOT/web"
BACKEND_PORT=4096
FRONTEND_PORT=5173

INIT=false
MIGRATE=false
RESTART=false
STOP=false
BACKEND_ONLY=false
FRONTEND_ONLY=false

show_help() {
    echo "XYGo Admin 本地开发启动脚本"
    echo ""
    echo "用法:"
    echo "  ./start.sh                启动后端 + 前端"
    echo "  ./start.sh --init         从零初始化（配置、数据库、依赖、迁移）"
    echo "  ./scripts/init.sh --start 同上，完成后自动启动"
    echo "  ./start.sh --migrate      启动前先执行数据库迁移"
    echo "  ./start.sh --restart      重启服务（先停后启）"
    echo "  ./start.sh --stop         停止服务"
    echo "  ./start.sh --backend-only 仅操作后端 (http://localhost:4096)"
    echo "  ./start.sh --frontend-only 仅操作前端 (http://localhost:5173)"
    echo ""
    echo "示例:"
    echo "  ./start.sh --restart --backend-only   仅重启后端"
    echo "  ./start.sh --stop                       停止全部"
    echo ""
    echo "环境要求: Go 1.24+, gf CLI, Node.js 20.19+, pnpm 8.8+, MySQL 8, Redis 7"
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --init) INIT=true ;;
        --migrate) MIGRATE=true ;;
        --restart) RESTART=true ;;
        --stop) STOP=true ;;
        --backend-only) BACKEND_ONLY=true ;;
        --frontend-only) FRONTEND_ONLY=true ;;
        --help|-h) show_help; exit 0 ;;
        *) echo "未知参数: $1"; show_help; exit 1 ;;
    esac
    shift
done

ensure_config() {
    local server_config="$SERVER_DIR/manifest/config/config.yaml"
    local server_example="$SERVER_DIR/manifest/config/config.yaml.example"
    if [[ ! -f "$server_config" ]]; then
        [[ -f "$server_example" ]] || { echo "缺少后端配置模板: $server_example"; exit 1; }
        cp "$server_example" "$server_config"
        echo "配置: 已生成 server/manifest/config/config.yaml，请按需修改数据库与 Redis 连接"
    fi

    local web_env_local="$WEB_DIR/.env.local"
    local web_env_dev="$WEB_DIR/.env.development"
    if [[ ! -f "$web_env_local" ]]; then
        [[ -f "$web_env_dev" ]] || { echo "缺少前端环境变量模板: $web_env_dev"; exit 1; }
        cp "$web_env_dev" "$web_env_local"
        echo "配置: 已生成 web/.env.local"
    fi
}

run_init() {
    bash "$ROOT/scripts/init.sh" "$@"
}

run_migrate() {
    echo "迁移: 执行数据库迁移..."
    (cd "$SERVER_DIR" && go run tools.go migrate up)
    echo "迁移: 完成"
}

stop_port_listener() {
    local port="$1"
    if command -v lsof >/dev/null 2>&1; then
        local pids
        pids=$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)
        if [[ -n "$pids" ]]; then
            echo "$pids" | xargs kill -9 2>/dev/null || true
            echo "已停止端口 $port 上的进程"
        fi
    elif command -v fuser >/dev/null 2>&1; then
        fuser -k "$port"/tcp 2>/dev/null || true
    fi
}

stop_pid_file() {
    local pid_file="$1"
    if [[ -f "$pid_file" ]]; then
        local pid
        pid=$(cat "$pid_file")
        kill "$pid" 2>/dev/null || true
        rm -f "$pid_file"
    fi
}

stop_services() {
    local backend="$1"
    local frontend="$2"

    if $backend; then
        echo "停止后端 (端口 $BACKEND_PORT)..."
        stop_pid_file "$ROOT/.start-backend.pid"
        stop_port_listener "$BACKEND_PORT"
    fi
    if $frontend; then
        echo "停止前端 (端口 $FRONTEND_PORT)..."
        stop_pid_file "$ROOT/.start-frontend.pid"
        stop_port_listener "$FRONTEND_PORT"
    fi
}

start_backend() {
    command -v gf >/dev/null 2>&1 || { echo "未找到 gf，请安装 GoFrame CLI"; exit 1; }
    echo "后端: 启动中 -> http://localhost:4096"
    (cd "$SERVER_DIR" && gf run main.go) &
    echo $! > "$ROOT/.start-backend.pid"
}

start_frontend() {
    command -v pnpm >/dev/null 2>&1 || { echo "未找到 pnpm"; exit 1; }
    [[ -d "$WEB_DIR/node_modules" ]] || { echo "请先运行: ./start.sh --init"; exit 1; }
    echo "前端: 启动中..."
    (cd "$WEB_DIR" && pnpm dev) &
    echo $! > "$ROOT/.start-frontend.pid"
}

show_startup_summary() {
    local backend="$1"
    local frontend="$2"
    local is_restart="$3"

    echo ""
    echo "========== 访问地址 =========="
    if $backend; then
        echo "  后端 API    http://localhost:${BACKEND_PORT}"
    fi
    if $frontend; then
        echo "  用户界面    http://localhost:${FRONTEND_PORT}/"
        echo "  管理界面    http://localhost:${FRONTEND_PORT}/admin"
    fi
    echo ""

    if $is_restart; then
        echo "已重启服务，按 Ctrl+C 停止。"
    elif $backend && $frontend; then
        echo "后端与前端已启动，按 Ctrl+C 停止。"
    elif $frontend; then
        echo "前端已启动，按 Ctrl+C 停止。"
    else
        echo "后端已启动，按 Ctrl+C 停止。"
    fi
}

cleanup() {
    stop_services true true
}

if $INIT; then
    run_init
    exit 0
fi

START_BACKEND=true
START_FRONTEND=true
$BACKEND_ONLY && START_FRONTEND=false
$FRONTEND_ONLY && START_BACKEND=false

if $STOP; then
    stop_services "$START_BACKEND" "$START_FRONTEND"
    echo "服务已停止。"
    exit 0
fi

ensure_config

if $MIGRATE; then
    run_migrate
fi

if $RESTART; then
    echo "========== 重启服务 =========="
    stop_services "$START_BACKEND" "$START_FRONTEND"
    sleep 2
fi

if $START_BACKEND && $START_FRONTEND && ! $RESTART; then
    trap cleanup EXIT INT TERM
fi

$START_BACKEND && start_backend
$START_FRONTEND && start_frontend

if $START_BACKEND || $START_FRONTEND; then
    show_startup_summary "$START_BACKEND" "$START_FRONTEND" "$RESTART"
fi

if $START_BACKEND && $START_FRONTEND; then
    wait
fi
