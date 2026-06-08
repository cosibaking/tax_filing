# 主播 OPC 合规服务 —— MVP 开发 Backlog

> 版本：v1.0  
> 编制时间：2026年6月  
> 依据文档：[主播OPC合规服务-产品文档.md](主播OPC合规服务-产品文档.md)  
> 范围：**仅 P1 功能（32 项）**，目标版本 v0.1 MVP

---

## 一、文档说明

本文档将 PRD 中 32 项 P1 功能拆解为 **52 条 User Story**，组织为 **6 个 Epic**、**4 个 Sprint（各 2 周）**，供产品、设计与开发排期使用。

| 关联文档 | 用途 |
|---------|------|
| [主播OPC合规服务-产品文档.md](主播OPC合规服务-产品文档.md) | 功能定义、业务规则、合规约束 |
| [主播OPC合规服务-信息架构与页面原型.md](主播OPC合规服务-信息架构与页面原型.md) | 页面路径、线框、导航 |
| [README.md](README.md) | 技术栈与现有路由/API 模式 |

**估算说明：** S = 1–2 人日，M = 3–5 人日，L = 6–10 人日

---

## 二、MVP 目标与验收口径

### 2.1 核心闭环

```
诊断 → 签约 → OPC 落地 → 记账 → 申报提醒 → 月度对账单
```

### 2.2 端到端验收标准（MVP 发布门槛）

| # | 验收项 | 通过标准 |
|---|--------|---------|
| AC-01 | 匿名诊断 | 访客完成问卷并获得税负对比 + 方案推荐 |
| AC-02 | 注册签约 | 主播注册登录 → 风险告知 → 方案确认 → 订单创建 |
| AC-03 | OPC 落地 | 资料提交 → 顾问更新工商/税务/银行进度 → 主播可见时间轴 |
| AC-04 | 月度台账 | 主播录入/导入收入与费用 → 自动生成利润预览 |
| AC-05 | 申报任务 | 系统按日历生成增值税/企税任务 → 顾问标记已申报 |
| AC-06 | 对账单 | 次月生成月度对账单 → 推送通知 |
| AC-07 | 合规留痕 | 方案确认、风险告知、改账操作均有审计记录 |
| AC-08 | 敏感数据 | 身份证、银行流水加密存储，接口鉴权 |

### 2.3 Definition of Done（每条 Story）

- [ ] API 返回统一信封 `{ code, message, data }`
- [ ] 对应页面可访问（或 API 有契约测试）
- [ ] 合规相关 Story 有操作留痕
- [ ] PRD 功能编号（F-xx）在 Story 中可追溯
- [ ] 与信息架构文档中的页面路径一致

---

## 三、Epic 总览

| Epic | 名称 | P1 功能 | Story 数 | Sprint |
|------|------|---------|---------|--------|
| E1 | 账号与合规基础 | F-85、F-86、F-89、F-90 | 8 | S1 |
| E2 | 诊断与方案 | F-01、F-02、F-03 | 7 | S1 |
| E3 | 签约与 OPC 落地 | F-09~F-12、F-16~F-19 | 12 | S2 |
| E4 | 收入与成本台账 | F-26~F-30、F-35~F-38 | 11 | S3 |
| E5 | 记账与申报 | F-44~F-46、F-51~F-55、F-54 | 10 | S3–S4 |
| E6 | 报告、触达与合规 | F-61、F-62、F-65、F-68、F-69、F-74 | 8 | S4 |
| **合计** | | **32 项 P1** | **52** | **4 Sprint** |

---

## 四、迭代分期

```mermaid
flowchart LR
  S1[Sprint1_账号诊断] --> S2[Sprint2_签约OPC落地]
  S2 --> S3[Sprint3_台账记账]
  S3 --> S4[Sprint4_申报报告合规]
```

| Sprint | 周期 | Epic | 交付里程碑 |
|--------|------|------|-----------|
| **Sprint 1** | 2 周 | E1 + E2 | 诊断闭环可演示；会员扩展字段就绪 |
| **Sprint 2** | 2 周 | E3 | 签约 + OPC 落地进度可追踪 |
| **Sprint 3** | 2 周 | E4 + E5（前半） | 收支录入 + 利润表 + 税额计算引擎 |
| **Sprint 4** | 2 周 | E5（后半）+ E6 | 申报任务 + 对账单 + 合规引擎 + E2E 联调 |

---

## 五、User Story 明细

### Epic E1：账号与合规基础

| ID | PRD | 角色 | 用户故事 | 验收标准 | 依赖 | 估算 | Sprint |
|----|-----|------|---------|---------|------|------|--------|
| US-E1-01 | F-85 | 主播 | 作为未注册主播，我希望用手机号注册并登录，以便使用合规服务 | 复用 `/member/auth/register`、`/member/auth/login`；登录后跳转诊断或会员中心 | 无 | S | S1 |
| US-E1-02 | F-85 | 主播 | 作为已登录主播，我希望 session 过期后自动刷新 token，以便不中断操作 | 复用 `/member/auth/refresh`；member-http 拦截 401 并刷新 | US-E1-01 | S | S1 |
| US-E1-03 | F-86 | 主播 | 作为签约客户，我希望在资料页绑定 OPC 主体，以便后续台账归属正确 | 展示 OPC 名称、统一社会信用代码、状态；未设立时显示「设立中」 | US-E3-08 | M | S2 |
| US-E1-04 | F-86 | 系统 | 作为系统，一个会员账号默认绑定一个 OPC，以便符合 MVP 单主体模型 | `OpcEntity.memberId` 唯一约束；后台可手动关联 | US-E3-01 | S | S2 |
| US-E1-05 | F-90 | 主播 | 作为新用户，我希望注册前阅读并同意隐私政策与用户协议 | `/legal/privacy`、`/legal/terms` 页面；注册勾选必选 | 无 | S | S1 |
| US-E1-06 | F-90 | 系统 | 作为系统，我需记录用户同意协议的时间与版本号 | `ComplianceConsent` 表写入 type=terms/privacy | US-E1-05 | S | S1 |
| US-E1-07 | F-89 | 系统 | 作为系统，我需对身份证、银行流水等敏感字段加密存储 | 字段级加密或应用层 AES；日志不输出明文 | 无 | M | S1 |
| US-E1-08 | F-89 | 运营 | 作为运营，我仅能查看脱敏后的敏感信息 | 后台 API 返回掩码身份证/卡号 | US-E1-07 | S | S1 |

---

### Epic E2：诊断与方案

| ID | PRD | 角色 | 用户故事 | 验收标准 | 依赖 | 估算 | Sprint |
|----|-----|------|---------|---------|------|------|--------|
| US-E2-01 | F-01 | 访客 | 作为潜在主播，我希望填写合规诊断问卷，以便了解我的合规状态 | 字段：平台多选、月收入区间、年成本估算、现有主体、是否报税、是否被税局联系；支持匿名提交 | 无 | M | S1 |
| US-E2-02 | F-01 | 主播 | 作为已登录主播，我希望诊断结果保存到账户，以便后续签约引用 | 登录后提交关联 `memberId`；`/user/compliance/diagnosis` 可查看历史 | US-E1-01, US-E2-01 | S | S1 |
| US-E2-03 | F-02 | 访客 | 作为访客，我希望输入年收入与可扣除成本，查看三种方案税负对比 | 展示劳务/个体户/OPC 三列：预估税额、实际税负率；含「不报税」高危提示 | US-E2-01 | M | S1 |
| US-E2-04 | F-02 | 系统 | 作为系统，我需按 PRD 6.2–6.4 公式计算 OPC 税负 | 增值税 1%、月 10 万免税；企税小微 5%；公式单元测试覆盖 | US-E2-03 | M | S1 |
| US-E2-05 | F-03 | 访客 | 作为访客，我希望系统根据收入段推荐 OPC/个体户/过渡方案 | 输出一页纸建议：推荐主体、理由、预估税负区间、下一步 CTA | US-E2-01, US-E2-03 | M | S1 |
| US-E2-06 | F-03 | 访客 | 作为访客，我希望在诊断结果页下载/查看方案摘要 PDF | 服务端生成 PDF 或 HTML 打印版；含免责声明 | US-E2-05 | M | S1 |
| US-E2-07 | F-03 | 主播 | 作为主播，我希望从诊断结果一键进入签约流程 | 「确认方案并签约」跳转 `/user/compliance/plan`；未登录先引导登录 | US-E1-01, US-E2-05 | S | S1 |

---

### Epic E3：签约与 OPC 落地

| ID | PRD | 角色 | 用户故事 | 验收标准 | 依赖 | 估算 | Sprint |
|----|-----|------|---------|---------|------|------|--------|
| US-E3-01 | F-09 | 主播 | 作为主播，我希望查看基础/进阶/尊享套餐及月费 | `/pricing` 与签约页展示套餐名称、价格、包含项 | 无 | S | S2 |
| US-E3-02 | F-09 | 主播 | 作为主播，我希望选择套餐并创建服务订单 | 创建 `ServiceOrder`：套餐 ID、金额、状态 pending/active | US-E3-01, US-E1-01 | M | S2 |
| US-E3-03 | F-12 | 主播 | 作为主播，签约前我必须阅读并勾选合规风险告知 | 模态展示 4 条底线声明；全部勾选才可继续；不可跳过 | US-E3-02 | M | S2 |
| US-E3-04 | F-12 | 系统 | 作为系统，我需记录风险告知确认时间与 IP | `ComplianceConsent` type=risk_disclosure | US-E3-03 | S | S2 |
| US-E3-05 | F-11 | 主播 | 作为主播，我希望电子确认推荐方案（含 OPC） | 展示方案摘要；勾选「我已阅读并同意」；写入 consent type=plan_confirm | US-E2-05, US-E3-03 | M | S2 |
| US-E3-06 | F-10 | 主播 | 作为主播，我希望在线签署服务协议 | MVP：协议 PDF 预览 + 勾选确认 + 电子签名板（姓名确认）；存档 PDF | US-E3-05 | L | S2 |
| US-E3-07 | F-10 | 系统 | 作为系统，我需存档签约 PDF 及签署元数据 | 附件表关联 orderId；含签署时间、memberId | US-E3-06 | M | S2 |
| US-E3-08 | F-16 | 主播 | 作为主播，我希望提交 OPC 注册资料 | 表单：身份证正反面、手机号、邮箱、拟经营范围、注册地址；附件上传 | US-E3-06 | M | S2 |
| US-E3-09 | F-17 | 顾问 | 作为顾问，我希望在后台推进工商注册并更新进度 | 状态：资料审核→提交工商→执照下发；记录日期与备注 | US-E3-08 | M | S2 |
| US-E3-10 | F-18 | 顾问 | 作为顾问，执照下发后我需完成税务登记并更新状态 | 状态：税务登记中→已完成；记录电子税务局激活日期 | US-E3-09 | M | S2 |
| US-E3-11 | F-19 | 主播 | 作为主播，我希望查看银行开户指引并上传开户回执 | 展示材料清单与预约说明；上传开户许可证/回执 | US-E3-10 | M | S2 |
| US-E3-12 | F-17 | 主播 | 作为主播，我希望在 OPC 进度页看到全流程时间轴 | 工商→税务→银行→完成；当前步骤高亮；预计 SLA 提示 | US-E3-09~11 | M | S2 |

---

### Epic E4：收入与成本台账

| ID | PRD | 角色 | 用户故事 | 验收标准 | 依赖 | 估算 | Sprint |
|----|-----|------|---------|---------|------|------|--------|
| US-E4-01 | F-26 | 主播 | 作为主播，我希望按平台录入月度收入 | 平台 Tab：抖音/快手/B站/视频号/小红书；列表分页 | US-E3-12 | M | S3 |
| US-E4-02 | F-27 | 主播 | 作为主播，我希望为每笔收入选择类型 | 类型：打赏/带货佣金/广告/坑位费/线下商单/其他 | US-E4-01 | S | S3 |
| US-E4-03 | F-29 | 主播 | 作为主播，我希望记录平台服务费并自动算实收 | 字段：含税收入、平台服务费、实收金额；实收=含税-服务费 | US-E4-01 | M | S3 |
| US-E4-04 | F-28 | 主播 | 作为主播，我希望 CSV 导入平台流水 | 模板下载；校验必填列；导入结果成功/失败条数 | US-E4-01 | M | S3 |
| US-E4-05 | F-30 | 主播 | 作为主播，我希望导入银行流水并与平台收入对账 | CSV 导入；按日期+金额匹配；未匹配项标黄 | US-E4-04 | L | S3 |
| US-E4-06 | F-35 | 主播 | 作为主播，我希望按费用类型录入成本 | 分类：设备/网费/场地/投流/外包/差旅/其他 | US-E3-12 | M | S3 |
| US-E4-07 | F-37 | 主播 | 作为主播，录入费用时我希望看到该类型所需凭证说明 | 费用类型库：名称、所需凭证、合规提示 | US-E4-06 | S | S3 |
| US-E4-08 | F-36 | 主播 | 作为主播，我希望上传发票/凭证附件 | 支持图片/PDF；关联 expenseId；复用附件服务 | US-E4-06 | M | S3 |
| US-E4-09 | F-38 | 系统 | 作为系统，成本占收入比超阈值时我需预警 | 可配置阈值（默认 80%）；提交时展示警告，需二次确认 | US-E4-06 | M | S3 |
| US-E4-10 | F-38 | 系统 | 作为系统，我需拦截标记为「无真实业务」的费用描述 | 关键词黑名单；禁止保存或强制顾问审核 | US-E4-06 | M | S3 |
| US-E4-11 | F-65 | 系统 | 作为系统，我需在每月 1–5 日提醒主播上传上月材料 | 定时任务 + 站内通知；文案含截止日 | US-E4-01, US-E4-06 | M | S4 |

---

### Epic E5：记账与申报

| ID | PRD | 角色 | 用户故事 | 验收标准 | 依赖 | 估算 | Sprint |
|----|-----|------|---------|---------|------|------|--------|
| US-E5-01 | F-44 | 系统 | 作为系统，我需为 OPC 小规模纳税人初始化标准账套 | 预置科目：银行存款、主营业务收入、管理费/推广费等 | US-E3-12 | M | S3 |
| US-E5-02 | F-45 | 系统 | 作为系统，收入/费用确认后我需自动生成会计分录 | 收入借银行贷收入；费用借费用贷银行；写入 `LedgerVoucher` | US-E4-01, US-E4-06, US-E5-01 | L | S3 |
| US-E5-03 | F-46 | 主播 | 作为主播，我希望查看月度/季度/年度利润表 | 展示：收入（不含税）、成本、利润、累计年度利润 | US-E5-02 | M | S3 |
| US-E5-04 | F-51 | 系统 | 作为系统，我需计算小规模增值税（含 10 万免税） | 不含税=含税÷1.01；月销&lt;10万免税；否则 1% | US-E5-03 | M | S4 |
| US-E5-05 | F-55 | 系统 | 作为系统，我需计算附加税费 | 增值税×6%（减半后）；无增值税则为 0 | US-E5-04 | S | S4 |
| US-E5-06 | F-52 | 系统 | 作为系统，我需计算企业所得税季度预缴（小微 5%） | 季度累计利润×5%；利润≤300万适用小微 | US-E5-03 | M | S4 |
| US-E5-07 | F-53 | 系统 | 作为系统，我需生成年度企税汇算任务 | 次年 5/31 截止；任务类型 annual_cit | US-E5-06 | S | S4 |
| US-E5-08 | F-54 | 主播 | 作为主播，我希望在申报日历查看截止日与任务状态 | 月历视图；任务：pending/filed/overdue；倒计时 | US-E5-04~07 | M | S4 |
| US-E5-09 | F-51 | 顾问 | 作为顾问，我希望在后台标记申报完成并上传回执 | 批量列表；状态 filed；上传 PDF 回执 | US-E5-08 | M | S4 |
| US-E5-10 | F-69 | 顾问 | 作为顾问，申报前我需完成 9 项自查清单 | PRD 附录 B 九项；全部勾选才可标记 filed | US-E5-09 | M | S4 |

---

### Epic E6：报告、触达与合规

| ID | PRD | 角色 | 用户故事 | 验收标准 | 依赖 | 估算 | Sprint |
|----|-----|------|---------|---------|------|------|--------|
| US-E6-01 | F-61 | 系统 | 作为系统，我需在次月 16–20 日生成月度对账单 | 含：收入/成本/利润/预缴税额/累计年度利润 | US-E5-03, US-E5-06 | M | S4 |
| US-E6-02 | F-61 | 主播 | 作为主播，我希望查看并下载月度对账单 PDF | `/user/compliance/statement` 列表 + PDF 下载 | US-E6-01 | M | S4 |
| US-E6-03 | F-62 | 主播 | 作为主播，申报完成后我希望收到服务完成通知 | 复用 `MemberNotice`；模板含本月摘要 | US-E6-01, US-E5-09 | S | S4 |
| US-E6-04 | F-68 | 系统 | 作为系统，我需禁止录入「借款/赠与」类收入分类 | 收入类型无此类选项；备注关键词拦截 | US-E4-01 | S | S4 |
| US-E6-05 | F-68 | 系统 | 作为系统，我需禁止故意低报收入的操作 | 申报额低于台账汇总时阻断并提示 | US-E5-09 | M | S4 |
| US-E6-06 | F-69 | 主播 | 作为主播，我可在申报页查看自查清单状态 | 只读展示；未完成项高亮 | US-E5-10 | S | S4 |
| US-E6-07 | F-74 | 系统 | 作为系统，我需记录关键操作审计日志 | 改账、删账、方案确认、申报标记；含操作人、时间、前后值 | 全 Epic | M | S4 |
| US-E6-08 | F-74 | 顾问 | 作为顾问，我希望导出某客户的合规留痕 | 后台按 memberId 导出 consent + audit 列表 | US-E6-07 | M | S4 |

---

## 六、P1 功能覆盖矩阵

| PRD | 功能名称 | 对应 Story |
|-----|---------|-----------|
| F-01 | 合规诊断问卷 | US-E2-01, US-E2-02 |
| F-02 | 税负对比计算器 | US-E2-03, US-E2-04 |
| F-03 | 方案推荐引擎 | US-E2-05, US-E2-06, US-E2-07 |
| F-09 | 服务套餐配置 | US-E3-01, US-E3-02 |
| F-10 | 电子签约 | US-E3-06, US-E3-07 |
| F-11 | 方案书面确认 | US-E3-05 |
| F-12 | 强制风险告知 | US-E3-03, US-E3-04 |
| F-16 | 注册资料收集 | US-E3-08 |
| F-17 | 工商注册代办 | US-E3-09, US-E3-12 |
| F-18 | 税务登记 | US-E3-10 |
| F-19 | 银行开户指引 | US-E3-11 |
| F-26 | 多平台收入台账 | US-E4-01 |
| F-27 | 收入分类 | US-E4-02 |
| F-28 | 平台流水导入 | US-E4-04 |
| F-29 | 平台服务费扣除 | US-E4-03 |
| F-30 | 银行流水对账 | US-E4-05 |
| F-35 | 费用台账 | US-E4-06 |
| F-36 | 发票/凭证上传 | US-E4-08 |
| F-37 | 合规费用类型库 | US-E4-07 |
| F-38 | 虚增成本拦截 | US-E4-09, US-E4-10 |
| F-44 | OPC 标准账套 | US-E5-01 |
| F-45 | 自动分录生成 | US-E5-02 |
| F-46 | 利润表 | US-E5-03 |
| F-51 | 增值税申报 | US-E5-04, US-E5-09 |
| F-52 | 企税季度预缴 | US-E5-06 |
| F-53 | 企税年度汇算 | US-E5-07 |
| F-54 | 申报日历与提醒 | US-E5-08 |
| F-55 | 附加税费计算 | US-E5-05 |
| F-61 | 月度对账单 | US-E6-01, US-E6-02 |
| F-62 | 月度服务完成通知 | US-E6-03 |
| F-65 | 材料催收任务 | US-E4-11 |
| F-68 | 合规底线规则引擎 | US-E6-04, US-E6-05 |
| F-69 | 报税前自查清单 | US-E5-10, US-E6-06 |
| F-74 | 操作留痕与取证 | US-E6-07, US-E6-08 |
| F-85 | 主播端注册登录 | US-E1-01, US-E1-02 |
| F-86 | OPC 主体绑定 | US-E1-03, US-E1-04 |
| F-89 | 数据安全与加密 | US-E1-07, US-E1-08 |
| F-90 | 隐私政策与用户协议 | US-E1-05, US-E1-06 |

**覆盖率：32/32 P1 功能均已映射。**

---

## 七、MVP 明确不做（Out of Scope）

以下能力归属 PRD P2+，不在 v0.1 排期：

| 功能 | PRD | 原因 |
|------|-----|------|
| 收入一致性自动比对 | F-31 | 需复杂勾稽规则，MVP 人工对账 |
| 截图 OCR 导入 | F-28 | MVP 仅 CSV |
| 真实电子税务局 API 对接 | F-51~53 | MVP 顾问人工申报 + 系统记录 |
| 腾讯电子签/法大大 API | F-10 | MVP 勾选确认 + PDF 存档 |
| 在线支付/续费 | F-14 | MVP 线下确认或后台开通 |
| 季度/年度报告 | F-63、F-64 | P2 |
| 管理端 CRM/SLA 预警 | F-79、F-81 | P2 简版仅 OPC 任务队列 |
| 刻章备案 | F-20 | P2 |

---

## 八、附录 A：数据模型草案

> 供 Prisma 建模参考，字段可在实现时微调。

### ComplianceDiagnosis（诊断记录）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BigInt | PK |
| memberId | BigInt? | 匿名时为空 |
| platforms | JSON | 平台列表 |
| monthlyIncomeRange | String | 月收入区间 |
| annualCostEstimate | Decimal | 年成本估算 |
| existingEntity | String | none/individual/opc/other |
| hasFiledTax | Boolean | 是否报税 |
| taxBureauContact | Boolean | 是否被税局联系 |
| recommendedPlan | String | opc/individual/labor/transitional |
| taxComparison | JSON | 三方案税负结果 |
| createdAt | DateTime | |

### OpcEntity（OPC 主体）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BigInt | PK |
| memberId | BigInt | FK，唯一 |
| companyName | String | 公司名 |
| creditCode | String? | 统一社会信用代码 |
| status | Enum | pending/registering/tax/bank/active |
| businessScope | String? | 经营范围 |
| registerAddress | String? | 注册地址 |
| bankAccount | String? | 对公账号（加密） |
| licenseFileId | BigInt? | 执照附件 |

### ServicePlan / ServiceOrder

| 表 | 关键字段 |
|----|---------|
| ServicePlan | name, tier(basic/advanced/premium), monthlyPrice, features(JSON) |
| ServiceOrder | memberId, planId, status, signedAt, contractFileId |

### ComplianceConsent（合规留痕）

| 字段 | 说明 |
|------|------|
| type | terms / privacy / risk_disclosure / plan_confirm / contract_sign |
| memberId, orderId?, ip, userAgent, agreedAt, documentVersion |

### IncomeEntry / ExpenseEntry

| 表 | 关键字段 |
|----|---------|
| IncomeEntry | opcId, platform, category, grossAmount, platformFee, netAmount, occurredAt, source(manual/csv) |
| ExpenseEntry | opcId, category, amount, invoiceType, attachmentId, occurredAt, warningFlag |

### LedgerVoucher / ProfitSummary

| 表 | 关键字段 |
|----|---------|
| LedgerVoucher | opcId, period, debitAccount, creditAccount, amount, refType, refId |
| ProfitSummary | opcId, year, month, revenue, cost, profit, cumulativeProfit |

### TaxFilingTask

| 字段 | 说明 |
|------|------|
| opcId, taxType | vat / cit_quarterly / cit_annual / surcharge |
| period, dueDate, status | pending/filed/overdue |
| calculatedAmount, filedAmount, receiptFileId |

### MonthlyStatement

| 字段 | 说明 |
|------|------|
| opcId, year, month, summary(JSON), pdfFileId, sentAt |

### ComplianceAuditLog

| 字段 | 说明 |
|------|------|
| entityType, entityId, action, operatorId, operatorType(member/admin), before, after, createdAt |

---

## 九、附录 B：API 路由规划

### 公开 API（无需登录）

| 方法 | 路径 | Story | 说明 |
|------|------|-------|------|
| POST | `/site/compliance/diagnosis` | US-E2-01 | 提交诊断问卷 |
| POST | `/site/compliance/calculator` | US-E2-03 | 税负计算 |
| GET | `/site/compliance/plans` | US-E3-01 | 套餐列表 |

### 会员 API（memberRoute）

| 方法 | 路径 | Story | 说明 |
|------|------|-------|------|
| GET | `/member/compliance/diagnosis` | US-E2-02 | 诊断历史 |
| POST | `/member/compliance/order` | US-E3-02 | 创建订单 |
| POST | `/member/compliance/consent` | US-E3-04,05 | 风险告知/方案确认 |
| POST | `/member/compliance/sign` | US-E3-06 | 签署协议 |
| POST | `/member/compliance/opc/materials` | US-E3-08 | 提交注册资料 |
| GET | `/member/compliance/opc/progress` | US-E3-12 | OPC 进度 |
| GET/POST | `/member/compliance/income` | US-E4-01~04 | 收入 CRUD/导入 |
| GET/POST | `/member/compliance/expense` | US-E4-06~08 | 费用 CRUD |
| POST | `/member/compliance/bank/import` | US-E4-05 | 银行流水导入 |
| GET | `/member/compliance/ledger/profit` | US-E5-03 | 利润表 |
| GET | `/member/compliance/tax/calendar` | US-E5-08 | 申报日历 |
| GET | `/member/compliance/tax/checklist` | US-E6-06 | 自查清单 |
| GET | `/member/compliance/statements` | US-E6-02 | 对账单列表 |
| GET | `/member/compliance/opc` | US-E1-03 | OPC 主体信息 |

### 管理 API（adminRoute）

| 方法 | 路径 | Story | 说明 |
|------|------|-------|------|
| GET | `/admin/compliance/customers` | — | 客户列表（P2 完整 CRM） |
| GET/PATCH | `/admin/compliance/opc-tasks` | US-E3-09~11 | OPC 任务队列 |
| GET/PATCH | `/admin/compliance/filing` | US-E5-09 | 申报任务处理 |
| POST | `/admin/compliance/statements/send` | US-E6-01 | 生成/发送对账单 |
| GET | `/admin/compliance/audit/export` | US-E6-08 | 留痕导出 |

---

## 十、附录 C：与现有代码复用点

| 能力 | 现有路径 | 复用方式 |
|------|---------|---------|
| 会员认证 | `app/member/auth/*` | 直接复用 |
| 会员 HTTP | `lib/api-client/member-http.ts` | 扩展 compliance 前缀 |
| 附件上传 | `app/admin/upload/file/route.ts` | 会员端共用或封装 |
| 通知 | `MemberNotice` + member notice API | 催收、申报完成通知 |
| 会员中心路由 | `lib/portal/member-registry.tsx` | 注册 compliance 视图 |
| 后台动态路由 | `lib/admin/registry.tsx` | 注册 compliance 管理视图 |
| 定时任务 | `lib/cron/*` | 催收、对账单生成 |
| CMS 文档 | `/docs/*` | P2 政策科普 |

---

*本文档随 Sprint 推进更新 Story 状态。税率与政策以主管机关最新文件为准。*
