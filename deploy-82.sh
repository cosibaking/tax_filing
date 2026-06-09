#!/usr/bin/env bash
set -euo pipefail

APP=/home/ubuntu/apps/tax_filing
cd "$APP"
git checkout -f mis
git log -1 --oneline

mkdir -p server/manifest/config
cp /home/ubuntu/apps/tax-bridge/manifest/config/config.yaml server/manifest/config/config.yaml

# 前端在本地构建后 rsync 同步（服务器 pnpm build 易失败）
test -f server/resource/public/dist/index.html || { echo "缺少 dist，请先在本地 web 目录执行 vite build"; exit 1; }

cd "$APP/server"
test -x tax-filing-server || { echo "缺少 tax-filing-server 二进制，请本地 GOOS=linux GOARCH=amd64 go build"; exit 1; }

pm2 delete tax-filing 2>/dev/null || true

sudo tee /etc/systemd/system/tax-bridge.service >/dev/null <<'EOF'
[Unit]
Description=tax-filing XYGo Admin Server (mis)
After=network.target mysql.service redis.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/apps/tax_filing/server
ExecStart=/home/ubuntu/apps/tax_filing/server/tax-filing-server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable tax-bridge.service
sudo systemctl restart tax-bridge.service
sleep 4

systemctl is-active tax-bridge.service
curl -sS --max-time 8 http://127.0.0.1:4096/ | head -c 200 || true
echo
