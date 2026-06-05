# 主播 OPC 合规服务 MVP

帮助个人主播以 OPC（一人有限责任公司）为主体完成工商注册、日常记账与依法报税的全链路合规 SaaS。

## 文档

- [技术选型与架构](docs/01-技术选型与架构设计.md)
- [多 Agent 开发编排](docs/02-多Agent开发编排.md)
- [模块设计 M0~M8](docs/modules/)

## 技术栈

Next.js 15 · TypeScript · Prisma · MySQL · Redis · Tailwind · shadcn/ui

## 快速开始

### 1. 环境准备

```bash
cp .env.example .env
docker compose up -d mysql redis
npm install
```

### 2. 数据库

```bash
npx prisma migrate dev --name init
npm run db:seed
```

### 3. 开发

```bash
npm run dev
```

访问 http://localhost:3000

### 默认账号

| 角色 | 账号 | 密码 |
|------|------|------|
| 管理员 | admin | admin123 |
| 会员 | 手机号注册 | 自定义 |
| 短信验证码（开发） | — | 123456 |

## 核心路径

| 区域 | 路径 |
|------|------|
| 首页 | `/` |
| 免费诊断 | `/diagnosis` |
| 会员中心 | `/user/overview` |
| 管理后台 | `/admin/compliance/customers` |

## 脚本

| 命令 | 说明 |
|------|------|
| `npm run dev` | 开发服务器 |
| `npm run build` | 生产构建 |
| `npm test` | 单元测试 |
| `npm run typecheck` | 类型检查 |
| `npm run db:seed` | 初始化基础数据（管理员、套餐、科目） |
| `npm run db:seed-random` | 生成随机测试会员与 OPC 台账数据 |
| `npm run db:deploy:seed` | 部署迁移 + 基础种子 + 随机测试数据 |

## MVP 闭环

诊断 → 签约 → OPC 落地 → 记账 → 申报提醒 → 月度对账单

## Docker 部署

通过 Nginx 对外暴露，访问路径为 `http://<外网地址>/tax_filing`。

### 1. 准备环境变量

```bash
cp .env.docker.example .env
# 编辑 .env，至少修改数据库密码与 JWT/加密密钥
```

### 2. 开发模式（含随机测试数据）

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
```

启动后自动执行：数据库迁移 → 基础种子（管理员/套餐）→ 随机测试会员数据。

默认账号：

| 角色 | 账号 | 密码 |
|------|------|------|
| 管理员 | admin | admin123 |
| 测试会员 | 见容器日志 | test123 |

### 3. 线上模式（不生成随机测试数据）

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

或使用一键部署脚本（推荐在服务器上执行，含 iptables/旧栈检查与健康等待）：

```bash
chmod +x scripts/deploy-prod.sh
./scripts/deploy-prod.sh          # 构建并启动；80 被占用时自动回退 8048
./scripts/deploy-prod.sh --pull   # 先 git pull 再部署
./scripts/deploy-prod.sh --port 8048   # 强制使用备份端口
```

**80 端口被占用时的备份方案**：脚本会检测 80 是否已被其它进程占用；若占用则自动改用 `8048` 端口，并同步更新 `.env` 中的 `NGINX_HTTP_PORT` 与 `NEXT_PUBLIC_APP_URL`（如 `http://<IP>:8048/tax_filing`）。请在云安全组放行 TCP 8048。

启动后自动执行：数据库迁移 → 基础种子（管理员/套餐），跳过随机测试数据。

### 4. 访问

- 首页：`http://<外网地址>/tax_filing`
- 健康检查：`http://<外网地址>/tax_filing/api/health`

## 环境变量

见 [.env.example](.env.example)（本地开发）与 [.env.docker.example](.env.docker.example)（Docker 部署）
