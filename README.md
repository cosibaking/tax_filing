<p align="center">
  <img src="https://xygoupload.xingyunwangluo.com/gitee/%E5%8D%95%E7%8B%AClogo.png" width="200" />
</p>
<br />
<h1 align="center">Streamer OPC Compliance Service</h1>
<p align="center">A compliance SaaS for individual streamers built on XYGo Admin (Vue3 + GoFrame). Covers diagnosis, onboarding, OPC setup, bookkeeping, tax filing, and monthly statements. MVP v0.1 accepted.</p>
<div align="center">English | <a href="./README.zh-CN.md">简体中文</a></div>

<br />
<p align="center">
  <a href="https://www.xygoadmin.com">XYGo Admin</a> |
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

### Overview

**Streamer OPC Compliance Service** helps individual streamers and MCN-signed creators operate through an OPC (One-Person Company) entity — covering business registration, daily bookkeeping, and lawful tax filing. Built on the [XYGo Admin](https://www.xygoadmin.com) full-stack admin framework, it implements the compliance MVP (32 P1 features).

Frontend: [Art Design Pro](https://github.com/Daymychen/art-design-pro) (Vue3 + TypeScript + Element Plus). Backend: [GoFrame v2](https://goframe.org/).

### Core Loop

```
Diagnosis → Sign-up → OPC Setup → Bookkeeping → Filing Reminders → Monthly Statements
```

| Stage | Streamer Portal | Advisor Console |
|-------|-----------------|-----------------|
| Compliance Diagnosis | `/diagnosis` — anonymous questionnaire + 3-way tax comparison | — |
| Plan & Sign-up | `/user/compliance/plan` — risk disclosure + e-sign | — |
| OPC Setup | `/user/compliance/opc` — materials + progress timeline | `/admin/compliance/opc-tasks` |
| Ledger | `/user/compliance/income` · `expense` · `ledger` | — |
| Tax Filing | `/user/compliance/tax` | `/admin/compliance/filing` |
| Statements | `/user/compliance/statement` | `/admin/compliance/statements` |

### Key Features

**Compliance Diagnosis**: Anonymous or logged-in questionnaire with labor / sole proprietorship / OPC tax comparison and recommendations

**Consent & Audit Trail**: Risk disclosure, plan confirmation, e-sign — all recorded in `ComplianceConsent` and audit logs

**OPC Setup**: Material submission, progress timeline, advisor-driven business / tax / bank milestones

**Income & Expense Ledger**: Entry, profit preview, vouchers, and compliance rule validation

**Filing & Statements**: Filing calendar, advisor mark-as-filed, monthly statement generation and notifications

**Sensitive Data Protection**: AES encryption for ID numbers and bank accounts; masked display in APIs

**Platform Capabilities (XYGo Admin)**: RBAC, code generator, member portal, message queue, system monitoring, single-binary deploy

### Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Vue 3, TypeScript, Vite, Element Plus, Tailwind CSS, Pinia |
| Backend | GoFrame v2 (Go 1.24+) |
| Database | MySQL 8.0+ (primary), PostgreSQL 14+ (optional) |
| Cache / Queue | Redis |
| Auth | JWT (member Bearer + admin session) |

### Quick Start

**Requirements**: Node.js ≥ 20.19, pnpm ≥ 8.8, Go ≥ 1.24, MySQL 8, Redis 7

```bash
# 1. Initialize database
mysql -u root -p -e "CREATE DATABASE xygo DEFAULT CHARSET utf8mb4;"
mysql -u root -p xygo < mysql_install.sql

# 2. Backend config & migrations
cd server
cp manifest/config/config.yaml.example manifest/config/config.yaml
# Edit config.yaml: database.default.link, redis, auth.jwt.secret
go run tools.go migrate up    # Apply compliance MVP migrations (1.4.x)
gf run main.go                # Default http://localhost:4096

# 3. Frontend dev server
cd web
pnpm install
cp .env.development .env.local   # Adjust VITE_API_PROXY_URL if needed
pnpm dev                         # Default http://localhost:5173
```

**Smoke test path** (see [docs/03-验收记录.md](./docs/03-验收记录.md)):

```
/ → /diagnosis → /diagnosis/result
→ /user/register → /user/compliance/plan → /user/compliance/opc
→ (admin) /admin/compliance/opc-tasks
→ /user/compliance/income → expense → ledger → tax → statement
```

### Default Accounts

| Role | Username | Password | Entry |
|------|----------|----------|-------|
| Super Admin | Super | 123456 | `/admin` |
| Streamer (member) | Self-register | — | `/user/register` |

### Project Structure

```
xygo-admin/
├── docs/                              # Design & acceptance docs
│   ├── 01-技术选型与架构设计.md
│   ├── 02-多Agent开发编排.md
│   ├── 03-验收记录.md
│   └── modules/                       # M0–M8 module specs
├── server/                            # GoFrame backend
│   ├── api/
│   │   ├── site/site_compliance.go    # Public compliance APIs
│   │   ├── member/member_compliance.go
│   │   └── admin/admin_compliance.go
│   ├── internal/logic/compliance/     # Compliance domain logic
│   │   ├── diagnosis/  order/  opc/
│   │   ├── ledger/     tax/    statement/
│   │   └── audit/      notice/ dashboard/
│   ├── internal/library/complianceverify/  # OPC material verification
│   ├── cmd_tools/migrate/             # DB migrations (incl. 1.4.x compliance)
│   └── manifest/config/config.yaml.example
├── web/                               # Vue3 frontend
│   ├── src/views/frontend/compliance/ # Streamer compliance pages
│   ├── src/views/backend/compliance/  # Advisor console pages
│   └── src/config/complianceVerify.ts # Frontend mock toggle
├── 主播OPC合规服务-产品文档.md
├── 主播OPC合规服务-MVP开发Backlog.md
├── mysql_install.sql
└── version.json
```

### Configuration

**OPC Material Verification (Dev Mock)**

Third-party ID OCR, three-factor verification, and phone real-name checks are not wired in for MVP — both sides use mock mode:

| Side | Setting | Dev default |
|------|---------|-------------|
| Backend | `compliance.verifyProvider` in `config.yaml` | `mock` |
| Frontend | `VITE_COMPLIANCE_VERIFY_MOCK` in `.env.development` | `true` |

For production: set backend provider to `aliyun` (or other) and implement `VerifyMaterials`; set frontend `VITE_COMPLIANCE_VERIFY_MOCK=false`.

**Unified Social Credit Code (执照下发 Mock)**

When an advisor marks a license as issued (`issue_license`), credit code validation also follows the same mock switch:

| Mode | Rule |
|------|------|
| Mock | Non-empty, max 18 chars (no format/checksum) |
| Production | GB 32100 format + checksum (`validateCreditCode` in `opc/validate.go` and `complianceVerify.ts`) |

Optional for production: third-party API to verify credit code matches the approved company name. See `docs/03-验收记录.md` §开发备忘 (2026-06-09).

**OPC Task Advance Forms (执照 / 税务 / 银行 Mock)**

Advisor `opc-task-detail-dialog` advance actions use the same mock switch. In mock mode, required fields are non-empty only; attachments skip DB/MIME checks (`validateAttachmentForMode`). Use `ArtFileUpload` with `value-type="id"` so `licenseFileId` / `bankReceiptFileId` send attachment IDs, not URLs.

### Documentation

| Document | Description |
|----------|-------------|
| [主播OPC合规服务-产品文档.md](./主播OPC合规服务-产品文档.md) | PRD, user journey, feature list |
| [主播OPC合规服务-MVP开发Backlog.md](./主播OPC合规服务-MVP开发Backlog.md) | P1 stories, sprint plan |
| [docs/01-技术选型与架构设计.md](./docs/01-技术选型与架构设计.md) | Architecture, API contract, security |
| [docs/03-验收记录.md](./docs/03-验收记录.md) | Acceptance checklist, smoke path |
| [docs/modules/](./docs/modules/) | M0–M8 module design |
| [XYGo Admin Docs](https://www.xygoadmin.com/docs) | Framework general capabilities |

### MVP Limitations

- No real e-tax bureau API (advisor files manually + system records)
- No Tencent e-sign (name confirmation + PDF archive only)
- No online payment (orders activated manually in admin)
- CSV import only for bank statements (no OCR)
- OPC material verification in mock mode (see Configuration above)

### Tools

Run from `server/`:

| Command | Description |
|---------|-------------|
| `go run tools.go` | Interactive menu |
| `go run tools.go migrate up` | Run database migrations |
| `go run tools.go migrate status` | View migration status |
| `gf gen dao` | Generate DAO from database tables |
| `gf gen service` | Generate service interfaces from logic |

### Acknowledgements

- [XYGo Admin](https://www.xygoadmin.com) — base admin framework
- [GoFrame](https://goframe.org/) — Go web framework
- [Art Design Pro](https://github.com/Daymychen/art-design-pro) — Vue3 admin template

### License

[MIT](./LICENSE) — inherits XYGo Admin license; free for commercial use.
