#!/usr/bin/env bash
# 单独执行 MySQL 导入 / 创建应用用户
# 用法:
#   ./scripts/import-mysql.sh
#   MYSQL_ROOT_PASSWORD=xxx ./scripts/import-mysql.sh --import
#   ./scripts/import-mysql.sh --import --database xygo --force

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
HOST="${MYSQL_HOST:-127.0.0.1}"
PORT="${MYSQL_PORT:-3306}"
DATABASE="${MYSQL_DATABASE:-xygo}"
APP_USER="${MYSQL_APP_USER:-yxgo}"
APP_PASSWORD="${MYSQL_APP_PASSWORD:-123456}"
ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-}"
IMPORT=false
FORCE=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --import) IMPORT=true ;;
    --force) FORCE=true ;;
    --host) HOST="$2"; shift ;;
    --port) PORT="$2"; shift ;;
    --database) DATABASE="$2"; shift ;;
    --root-password) ROOT_PASSWORD="$2"; shift ;;
    *) echo "未知参数: $1"; exit 1 ;;
  esac
  shift
done

command -v mysql >/dev/null 2>&1 || { echo "未找到 mysql 命令，请先安装 MySQL 客户端"; exit 1; }

if [[ -z "$ROOT_PASSWORD" ]]; then
  read -rsp "MySQL root 密码: " ROOT_PASSWORD
  echo
fi

mysql_exec() {
  mysql -h "$HOST" -P "$PORT" -u root -p"$ROOT_PASSWORD" -e "$1"
}

echo ">> 测试 MySQL 连接 ${HOST}:${PORT}..."
mysql_exec "SELECT VERSION() AS version;"

table_exists() {
  local count
  count=$(mysql -h "$HOST" -P "$PORT" -u root -p"$ROOT_PASSWORD" -N -e \
    "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='${DATABASE}' AND table_name='xy_admin_user';")
  [[ "$count" == "1" ]]
}

if $IMPORT; then
  if table_exists && ! $FORCE; then
    echo "!! 数据库 ${DATABASE} 已有基础表，跳过导入（使用 --force 强制重导）"
  else
    echo ">> 创建数据库 ${DATABASE} 并导入 mysql_install.sql ..."
    mysql_exec "CREATE DATABASE IF NOT EXISTS \`${DATABASE}\` DEFAULT CHARSET utf8mb4;"
    mysql -h "$HOST" -P "$PORT" -u root -p"$ROOT_PASSWORD" "$DATABASE" < "$ROOT/mysql_install.sql"
    echo "OK SQL 导入完成"
  fi
fi

echo ">> 创建应用用户 ${APP_USER} ..."
mysql_exec "
CREATE USER IF NOT EXISTS '${APP_USER}'@'localhost' IDENTIFIED BY '${APP_PASSWORD}';
CREATE USER IF NOT EXISTS '${APP_USER}'@'127.0.0.1' IDENTIFIED BY '${APP_PASSWORD}';
ALTER USER '${APP_USER}'@'localhost' IDENTIFIED BY '${APP_PASSWORD}';
ALTER USER '${APP_USER}'@'127.0.0.1' IDENTIFIED BY '${APP_PASSWORD}';
GRANT ALL PRIVILEGES ON \`${DATABASE}\`.* TO '${APP_USER}'@'localhost';
GRANT ALL PRIVILEGES ON \`${DATABASE}\`.* TO '${APP_USER}'@'127.0.0.1';
FLUSH PRIVILEGES;
"

echo "OK 完成。请在 server/manifest/config/config.yaml 中确认:"
echo "  mysql:${APP_USER}:${APP_PASSWORD}@tcp(${HOST}:${PORT})/${DATABASE}?loc=Local&parseTime=true"
