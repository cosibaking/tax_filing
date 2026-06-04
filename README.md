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
| `npm run db:seed` | 初始化数据 |

## MVP 闭环

诊断 → 签约 → OPC 落地 → 记账 → 申报提醒 → 月度对账单

## 环境变量

见 [.env.example](.env.example)
