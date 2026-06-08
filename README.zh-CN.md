<p align="center">
  <img src="https://xygoupload.xingyunwangluo.com/gitee/%E5%8D%95%E7%8B%AClogo.png" width="200" />
</p>
<br />
<h1 align="center">主播 OPC 合规服务</h1>
<p align="center">基于 XYGo Admin（Vue3 + GoFrame）构建的个人主播合规 SaaS，覆盖诊断、签约、OPC 落地、记账申报与对账单全链路。MVP v0.1 已验收通过。</p>
<div align="center">简体中文 | <a href="./README.md">English</a></div>

<br />
<p align="center">
  <a href="https://www.xygoadmin.com">XYGo Admin 官网</a> |
  <a href="https://gitee.com/a751300685a/xygo-admin">Gitee</a> |
  <a href="https://github.com/z312193608/xygo-admin">GitHub</a>
</p>

<div align="center">

[![license](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)
[![Vue](https://img.shields.io/badge/Vue-3.x-42b883.svg)](https://vuejs.org/)
[![GoFrame](https://img.shields.io/badge/GoFrame-v2-00ADD8.svg)](https://goframe.org/)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6.svg)](https://www.typescriptlang.org/)
[![Vite](https://img.shields.io/badge/Vite-6.x-646CFF.svg)](https://vitejs.dev/)
[![Element Plus](https://img.shields.io/badge/Element_Plus-2.x-409EFF.svg)](https://element-plus.org/)
[![MVP](https://img.shields.io/badge/MVP-v0.1-22c55e.svg)](./docs/03-验收记录.md)

</div>
<br />

### 介绍

**主播 OPC 合规服务**帮助个人主播及签约 MCN 主播，以 OPC（一人有限责任公司）为主体完成工商注册、日常记账与依法报税。本项目在 [XYGo Admin](https://www.xygoadmin.com) 通用中后台框架之上，实现了合规业务 MVP（32 项 P1 功能）。

前端基于 [Art Design Pro](https://github.com/Daymychen/art-design-pro)（Vue3 + TypeScript + Element Plus），后端基于 [GoFrame v2](https://goframe.org/)。

### 核心闭环

```
诊断 → 签约 → OPC 落地 → 记账 → 申报提醒 → 月度对账单
```

| 阶段 | 主播端 | 顾问端 |
|------|--------|--------|
| 合规诊断 | `/diagnosis` 匿名问卷 + 三方案税负对比 | — |
| 方案签约 | `/user/compliance/plan` 风险告知 + 电子签约 | — |
| OPC 落地 | `/user/compliance/opc` 资料提交 + 进度时间轴 | `/admin/compliance/opc-tasks` |
| 收支台账 | `/user/compliance/income` · `expense` · `ledger` | — |
| 申报交付 | `/user/compliance/tax` | `/admin/compliance/filing` |
| 对账单 | `/user/compliance/statement` | `/admin/compliance/statements` |

### 主要特性

**合规诊断**：匿名/登录诊断问卷，劳务 / 个体户 / OPC 三方案税负对比与推荐

**签约留痕**：风险告知、方案确认、电子签约，`ComplianceConsent` + 审计日志全程留痕

**OPC 落地**：资料提交、进度时间轴、顾问后台推进工商 / 税务 / 银行节点

**收支台账**：收入 / 费用录入、利润预览、分录与合规规则校验

**申报与对账单**：申报日历、顾问标记已申报、月度对账单生成与通知

**敏感数据保护**：身份证、银行账号 AES 加密存储，接口脱敏展示

**平台能力（XYGo Admin）**：RBAC 权限、代码生成、会员门户、消息队列、系统监控、单体部署

### 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3、TypeScript、Vite、Element Plus、Tailwind CSS、Pinia |
| 后端 | GoFrame v2（Go 1.24+） |
| 数据库 | MySQL 8.0+（主）、PostgreSQL 14+（可选） |
| 缓存 / 队列 | Redis |
| 认证 | JWT（会员 Bearer + 管理端 Session） |

### 快速开始

**环境要求**：Node.js ≥ 20.19、pnpm ≥ 8.8、Go ≥ 1.24、MySQL 8、Redis 7

```bash
# 1. 初始化数据库
mysql -u root -p -e "CREATE DATABASE xygo DEFAULT CHARSET utf8mb4;"
mysql -u root -p xygo < mysql_install.sql

# 2. 后端配置与迁移
cd server
cp manifest/config/config.yaml.example manifest/config/config.yaml
# 编辑 config.yaml：database.default.link、redis、auth.jwt.secret
go run tools.go migrate up    # 执行合规 MVP 等增量迁移
gf run main.go                # 默认 http://localhost:4096

# 3. 前端开发
cd web
pnpm install
cp .env.development .env.local   # 按需调整 VITE_API_PROXY_URL
pnpm dev                         # 默认 http://localhost:5173
```

**冒烟路径**（详见 [docs/03-验收记录.md](./docs/03-验收记录.md)）：

```
/ → /diagnosis → /diagnosis/result
→ /user/register → /user/compliance/plan → /user/compliance/opc
→ (admin) /admin/compliance/opc-tasks
→ /user/compliance/income → expense → ledger → tax → statement
```

### 默认账号

| 角色 | 账号 | 密码 | 入口 |
|------|------|------|------|
| 超级管理员 | Super | 123456 | `/admin` |
| 主播（会员） | 自行注册 | — | `/user/register` |

### 项目结构

```
xygo-admin/
├── docs/                              # 设计与验收文档
│   ├── 01-技术选型与架构设计.md
│   ├── 02-多Agent开发编排.md
│   ├── 03-验收记录.md
│   └── modules/                       # M0–M8 模块设计
├── server/                            # GoFrame 后端
│   ├── api/
│   │   ├── site/site_compliance.go    # 匿名合规 API
│   │   ├── member/member_compliance.go
│   │   └── admin/admin_compliance.go
│   ├── internal/logic/compliance/     # 合规领域逻辑
│   │   ├── diagnosis/  order/  opc/
│   │   ├── ledger/     tax/    statement/
│   │   └── audit/      notice/ dashboard/
│   ├── internal/library/complianceverify/  # OPC 资料校验
│   ├── cmd_tools/migrate/             # 数据库迁移（含 1.4.x 合规表）
│   └── manifest/config/config.yaml.example
├── web/                               # Vue3 前端
│   ├── src/views/frontend/compliance/ # 主播合规页面
│   ├── src/views/backend/compliance/  # 顾问后台页面
│   └── src/config/complianceVerify.ts # 前端 Mock 开关
├── 主播OPC合规服务-产品文档.md
├── 主播OPC合规服务-MVP开发Backlog.md
├── mysql_install.sql
└── version.json
```

### 配置说明

**OPC 资料真实性校验（开发 Mock）**

MVP 阶段身份证 OCR、三要素、手机实名等第三方校验暂未接入，前后端均以 Mock 放行：

| 端 | 配置项 | 开发默认值 |
|----|--------|-----------|
| 后端 | `compliance.verifyProvider`（`config.yaml`） | `mock` |
| 前端 | `VITE_COMPLIANCE_VERIFY_MOCK`（`.env.development`） | `true` |

生产接入第三方后：后端改为 `aliyun` 等并实现 `VerifyMaterials`，前端设 `VITE_COMPLIANCE_VERIFY_MOCK=false`。

### 文档索引

| 文档 | 说明 |
|------|------|
| [主播OPC合规服务-产品文档.md](./主播OPC合规服务-产品文档.md) | PRD、用户旅程、功能清单 |
| [主播OPC合规服务-MVP开发Backlog.md](./主播OPC合规服务-MVP开发Backlog.md) | P1 Story、Sprint 排期 |
| [docs/01-技术选型与架构设计.md](./docs/01-技术选型与架构设计.md) | 架构、API 契约、安全设计 |
| [docs/03-验收记录.md](./docs/03-验收记录.md) | AC 验收清单、冒烟路径 |
| [docs/modules/](./docs/modules/) | M0–M8 模块详细设计 |
| [XYGo Admin 官方文档](https://www.xygoadmin.com/docs) | 框架通用能力 |

### MVP 已知限制

- 无真实电子税务局 API（顾问人工申报 + 系统记录）
- 无腾讯电子签（姓名确认 + PDF 存档）
- 无在线支付（后台手动开通订单）
- 流水导入仅支持 CSV（无 OCR）
- OPC 资料校验为 Mock（见上方配置说明）

### 工具命令

在 `server/` 目录下执行：

| 命令 | 说明 |
|------|------|
| `go run tools.go` | 交互式菜单 |
| `go run tools.go migrate up` | 执行数据库迁移 |
| `go run tools.go migrate status` | 查看迁移状态 |
| `gf gen dao` | 根据数据库表生成 DAO |
| `gf gen service` | 根据 Logic 生成 Service 接口 |

### 特别鸣谢

- [XYGo Admin](https://www.xygoadmin.com) — 基础中后台框架
- [GoFrame](https://goframe.org/) — Go Web 框架
- [Art Design Pro](https://github.com/Daymychen/art-design-pro) — Vue3 后台模板

### 开源协议

[MIT](./LICENSE) — 基于 XYGo Admin 开源协议，可免费商用。
