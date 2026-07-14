# 月度经营体检结构化报告与 PDF 导出 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将月度体检从纯文本升级为分业务类别的结构化报告，为每项异常提供完整解释，并支持按报告版本导出 PDF。

**Architecture:** 后端确定性 `Builder` 从现有报告输入快照生成 `StructuredReport`，同时渲染兼容纯文本并持久化 JSON 快照；列表页面与 PDF 只读取该版本快照，保证已发布报告不可变。前端将总览、分类检查和异常详情拆为三个组件，历史纯文本报告继续降级展示。

**Tech Stack:** Go、GoFrame、GORM/GoFrame DB、MySQL 8/PostgreSQL、GoPDF、Vue 3、Element Plus、TypeScript、pnpm

---

## 文件结构

### 新增文件

- `server/cmd_tools/migrate/1.5.7_monthly_report_structured.mysql.sql`：MySQL 新增结构化报告字段。
- `server/cmd_tools/migrate/1.5.7_monthly_report_structured.pgsql.sql`：PostgreSQL 新增结构化报告字段。
- `server/internal/logic/compliance/report/types.go`：结构化报告、分类、检查项和异常类型。
- `server/internal/logic/compliance/report/builder.go`：确定性分类、汇总、异常标准化和文本摘要。
- `server/internal/library/reportpdf/generator.go`：从持久化报告快照生成 PDF。
- `server/internal/library/reportpdf/generator_test.go`：PDF 生成、历史摘要和分页测试。
- `web/src/views/frontend/compliance/reports/ReportSummary.vue`：体检结论、资料完整度和风险分布。
- `web/src/views/frontend/compliance/reports/ReportCategories.vue`：六类检查结果。
- `web/src/views/frontend/compliance/reports/ReportAnomalyList.vue`：异常摘要与完整详情。

### 修改文件

- `server/internal/logic/compliance/report/generate.go`：接入结构化结果和接口定义。
- `server/internal/logic/compliance/report/report_test.go`：Builder、排序、兼容和校验测试。
- `server/internal/logic/compliance/report/repository.go`：结构化 JSON 持久化、读取和按会员查询。
- `server/internal/service/compliance_assistant.go`：增加报告 PDF 导出服务方法。
- `server/api/member/member_compliance_assistant.go`：增加 PDF 导出路由定义。
- `server/internal/controller/member/compliance_assistant.go`：设置 PDF 响应头并输出二进制。
- `web/src/api/frontend/compliance/assistant.ts`：结构化类型和 Blob 下载 API。
- `web/src/views/frontend/compliance/reports.vue`：新版页面编排、状态和跳转。
- `docs/06-OPC合规助手-接口与部署.md`：补充结构化报告字段及 PDF 接口。
- `docs/07-OPC合规助手-验收清单.md`：补充三种报告状态和导出验收。

### 不修改

- 现有对账单 `statementpdf` 和 `/compliance/statements/{id}/pdf`。
- 收入、费用、申报业务逻辑。
- 已有报告记录的输入快照和纯文本内容。

## Task 1: 数据库迁移

**Files:**
- Create: `server/cmd_tools/migrate/1.5.7_monthly_report_structured.mysql.sql`
- Create: `server/cmd_tools/migrate/1.5.7_monthly_report_structured.pgsql.sql`

- [ ] **Step 1: 编写 MySQL 幂等迁移**

```sql
SET @column_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'xy_compliance_monthly_report'
    AND COLUMN_NAME = 'structured_json'
);
SET @sql := IF(
  @column_exists = 0,
  'ALTER TABLE `xy_compliance_monthly_report` ADD COLUMN `structured_json` JSON NULL AFTER `risk_snapshot_json`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
```

- [ ] **Step 2: 编写 PostgreSQL 幂等迁移**

```sql
ALTER TABLE xy_compliance_monthly_report
  ADD COLUMN IF NOT EXISTS structured_json jsonb NULL;
```

- [ ] **Step 3: 静态检查迁移不含破坏性语句**

Run:

```bash
rg -n "DROP|DELETE|TRUNCATE|UPDATE" server/cmd_tools/migrate/1.5.7_monthly_report_structured.*.sql
```

Expected: 无输出，退出码 1。

- [ ] **Step 4: 提交迁移**

```bash
git add server/cmd_tools/migrate/1.5.7_monthly_report_structured.*.sql
git commit -m "feat: add structured monthly report storage"
```

## Task 2: 结构化报告模型与 Builder

**Files:**
- Create: `server/internal/logic/compliance/report/types.go`
- Create: `server/internal/logic/compliance/report/builder.go`
- Modify: `server/internal/logic/compliance/report/generate.go`
- Modify: `server/internal/logic/compliance/report/report_test.go`

- [ ] **Step 1: 写失败测试，锁定空数据、风险计数和六类顺序**

在 `report_test.go` 增加：

```go
func TestBuildStructuredReportMarksInsufficientData(t *testing.T) {
	result, err := BuildStructured(Input{PeriodKey: "2026-07"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Summary.Conclusion != ConclusionAttention {
		t.Fatalf("conclusion=%s", result.Summary.Conclusion)
	}
	if len(result.Categories) != 6 {
		t.Fatalf("categories=%d", len(result.Categories))
	}
	want := []string{"business", "documents", "tax", "funds", "employment", "annual"}
	for i, code := range want {
		if result.Categories[i].Code != code {
			t.Fatalf("category[%d]=%s", i, result.Categories[i].Code)
		}
	}
	if result.Summary.DataNotice == "" {
		t.Fatal("missing data notice")
	}
}

func TestBuildStructuredReportCountsAndSortsRisks(t *testing.T) {
	in := Input{PeriodKey: "2026-07", Risks: []map[string]any{
		{"code": "LOW", "title": "低风险", "severity": "low", "categoryCode": "annual", "dueDate": "2026-08-20"},
		{"code": "HIGH_LATE", "title": "高风险后到期", "severity": "high", "categoryCode": "tax", "dueDate": "2026-08-10"},
		{"code": "HIGH_SOON", "title": "高风险先到期", "severity": "high", "categoryCode": "funds", "dueDate": "2026-08-01"},
	}}
	report, err := BuildStructured(in)
	if err != nil {
		t.Fatal(err)
	}
	if report.Summary.Conclusion != ConclusionUrgent || report.Summary.HighCount != 2 || report.Summary.LowCount != 1 {
		t.Fatalf("summary=%+v", report.Summary)
	}
	if report.Anomalies[0].Code != "HIGH_SOON" || report.Anomalies[1].Code != "HIGH_LATE" {
		t.Fatalf("order=%+v", report.Anomalies)
	}
	for _, item := range report.Anomalies {
		if item.Facts == "" || item.Basis == "" || item.Impact == "" || item.Recommendation == "" {
			t.Fatalf("incomplete anomaly=%+v", item)
		}
	}
}
```

- [ ] **Step 2: 运行测试并确认因类型和函数不存在而失败**

Run:

```bash
cd server
go test ./internal/logic/compliance/report -run 'TestBuildStructured' -v
```

Expected: FAIL，包含 `undefined: BuildStructured`。

- [ ] **Step 3: 定义稳定 JSON 类型**

在 `types.go` 定义：

```go
const (
	StructuredSchemaVersion = 1
	ConclusionNormal = "normal"
	ConclusionAttention = "attention"
	ConclusionUrgent = "urgent"
)

type StructuredReport struct {
	SchemaVersion int              `json:"schemaVersion"`
	Summary       ReportSummary    `json:"summary"`
	Categories    []ReportCategory `json:"categories"`
	Anomalies     []ReportAnomaly  `json:"anomalies"`
}

type ReportSummary struct {
	Conclusion       string `json:"conclusion"`
	CompletenessRate int    `json:"completenessRate"`
	HighCount        int    `json:"highCount"`
	MediumCount      int    `json:"mediumCount"`
	LowCount         int    `json:"lowCount"`
	DataNotice       string `json:"dataNotice"`
}

type ReportCategory struct {
	Code    string        `json:"code"`
	Name    string        `json:"name"`
	Status  string        `json:"status"`
	Summary string        `json:"summary"`
	Checks  []ReportCheck `json:"checks"`
}

type ReportCheck struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ReportAnomaly struct {
	Code                 string   `json:"code"`
	CategoryCode         string   `json:"categoryCode"`
	Title                string   `json:"title"`
	Severity             string   `json:"severity"`
	Facts                string   `json:"facts"`
	Basis                string   `json:"basis"`
	Impact               string   `json:"impact"`
	Recommendation       string   `json:"recommendation"`
	RequiredMaterials    []string `json:"requiredMaterials"`
	DueDate              string   `json:"dueDate"`
	RequiresManualReview bool     `json:"requiresManualReview"`
	RuleVersion          string   `json:"ruleVersion"`
}
```

- [ ] **Step 4: 实现确定性 Builder 与兼容摘要**

`BuildStructured` 必须：校验 `YYYY-MM`；固定生成六类；从 `completeness.rate` 读取 0-100 完整度；将风险 map 标准化并为缺失字段填入审慎默认文案；按 high/medium/low 和 dueDate 排序；空输入设置资料不足提示。`RenderContent` 输出标题、结论、完整度、风险分布和免责声明。

`Result` 增加：

```go
Structured StructuredReport `json:"structuredReport"`
```

`Generate` 改为返回 `(Result, error)`，`Service.Build` 先调用确定性 Builder；Narrator 只能替换 `Content`，不能修改 `Structured`。

- [ ] **Step 5: 增加账期校验和 Narrator 不可篡改结构测试**

```go
func TestBuildStructuredRejectsInvalidPeriod(t *testing.T) {
	if _, err := BuildStructured(Input{PeriodKey: "2026-13"}); err == nil {
		t.Fatal("invalid period must fail")
	}
}

func TestNarratorCannotReplaceStructuredFacts(t *testing.T) {
	svc := NewService(nil, fixedNarrator{content: "易读摘要", model: "test"})
	result, err := svc.Build(context.Background(), Input{PeriodKey: "2026-07", Risks: []map[string]any{{"code": "R1", "severity": "high"}}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Content != "易读摘要" || result.Structured.Summary.HighCount != 1 {
		t.Fatalf("result=%+v", result)
	}
}
```

- [ ] **Step 6: 运行报告包测试**

Run:

```bash
cd server
gofmt -w internal/logic/compliance/report/*.go
go test ./internal/logic/compliance/report -v
```

Expected: PASS。

- [ ] **Step 7: 提交 Builder**

```bash
git add server/internal/logic/compliance/report
git commit -m "feat: build structured monthly checkup reports"
```

## Task 3: 持久化与历史报告兼容

**Files:**
- Modify: `server/internal/logic/compliance/report/generate.go`
- Modify: `server/internal/logic/compliance/report/repository.go`
- Modify: `server/internal/logic/compliance/report/report_test.go`

- [ ] **Step 1: 写结构化 JSON 往返与历史降级失败测试**

```go
func TestDecodeStructuredFallsBackForLegacyContent(t *testing.T) {
	item := MonthlyReport{Content: "历史报告内容"}
	item.ApplyStructuredJSON("")
	if item.StructuredReport != nil || !item.Legacy {
		t.Fatalf("item=%+v", item)
	}
}

func TestDecodeStructuredRoundTrip(t *testing.T) {
	want, _ := BuildStructured(Input{PeriodKey: "2026-07"})
	raw, _ := json.Marshal(want)
	item := MonthlyReport{}
	item.ApplyStructuredJSON(string(raw))
	if item.StructuredReport == nil || item.StructuredReport.SchemaVersion != 1 || item.Legacy {
		t.Fatalf("item=%+v", item)
	}
}
```

- [ ] **Step 2: 运行测试确认辅助方法不存在**

Run: `cd server && go test ./internal/logic/compliance/report -run 'TestDecodeStructured' -v`  
Expected: FAIL，包含 `ApplyStructuredJSON undefined`。

- [ ] **Step 3: 扩展数据库模型和仓储接口**

`MonthlyReport` 增加：

```go
StructuredJSON   string            `json:"-" orm:"structured_json"`
StructuredReport *StructuredReport `json:"structuredReport,omitempty" orm:"-"`
Legacy           bool              `json:"legacy" orm:"-"`
CreatedAt        uint64            `json:"createdAt" orm:"create_time"`
```

`ApplyStructuredJSON` 在空值时设置 `Legacy=true`；无效 JSON 返回错误，列表服务记录错误后降级为历史摘要。`SaveDraft` 序列化 `result.Structured` 写入 `structured_json`，返回值同步携带结构化报告。`List` 扫描后逐条解码。

仓储接口增加：

```go
Get(context.Context, uint64, uint64) (*MonthlyReport, error)
```

查询必须同时使用 `id` 与 `member_id`。

- [ ] **Step 4: 保持版本并发安全**

`SaveDraft` 使用事务执行 Max(version)+Insert，保留现有 `uk_report_version` 唯一约束；遇到唯一键冲突返回「报告版本生成冲突，请重试」，不得覆盖已有报告。

- [ ] **Step 5: 运行报告包测试和编译**

Run:

```bash
cd server
gofmt -w internal/logic/compliance/report/*.go
go test ./internal/logic/compliance/report -v
go test ./internal/service ./internal/controller/member ./api/member
```

Expected: PASS。

- [ ] **Step 6: 提交持久化改造**

```bash
git add server/internal/logic/compliance/report server/internal/service/compliance_assistant.go
git commit -m "feat: persist structured monthly report snapshots"
```

## Task 4: 报告 PDF 生成器

**Files:**
- Create: `server/internal/library/reportpdf/generator.go`
- Create: `server/internal/library/reportpdf/generator_test.go`
- Reuse: `server/internal/library/statementpdf/fonts/statement.ttf`

- [ ] **Step 1: 写 PDF 生成失败测试**

```go
func TestGenerateStructuredReportPDF(t *testing.T) {
	data := Data{
		CompanyName: "测试 OPC 有限公司",
		PeriodKey: "2026-07",
		Version: 2,
		Status: "draft",
		GeneratedAt: "2026-07-14 12:00:00",
		Summary: Summary{Conclusion: "urgent", CompletenessRate: 72, HighCount: 1},
		Anomalies: []Anomaly{{Code: "R1", Title: "银行流水缺失", Severity: "high", Facts: "已登记收入但未上传银行流水", Basis: "收入记录与资料清单不一致", Impact: "可能影响收入与回款核对", Recommendation: "上传完整银行流水", RequiredMaterials: []string{"银行流水"}, RequiresManualReview: true}},
	}
	pdf, err := Generate(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(pdf) < 1024 || string(pdf[:4]) != "%PDF" {
		t.Fatalf("invalid pdf: %d", len(pdf))
	}
}

func TestGenerateLegacyReportPDF(t *testing.T) {
	pdf, err := Generate(Data{PeriodKey: "2026-06", Version: 1, LegacyContent: "历史报告摘要"})
	if err != nil || len(pdf) < 1024 {
		t.Fatalf("err=%v size=%d", err, len(pdf))
	}
}
```

- [ ] **Step 2: 运行测试确认包不存在**

Run: `cd server && go test ./internal/library/reportpdf -v`  
Expected: FAIL，提示目录或符号不存在。

- [ ] **Step 3: 实现 A4 PDF**

`Data` 包含企业名称、账期、版本、状态、生成时间、`Summary`、`Category[]`、`Anomaly[]` 和历史摘要。`reportpdf` 不导入 `compliance/report`，由服务层显式映射同名字段，避免 `report -> reportpdf -> report` 导入循环。`Generate` 使用 GoPDF 与现有内嵌中文字体，按标题、企业信息、总览、分类、异常详情、免责声明绘制；统一通过 `ensureSpace(requiredHeight)` 在内容不足时 `AddPage()`，异常字段使用 `MultiCellWithOption` 自动换行。

- [ ] **Step 4: 增加多异常分页测试**

生成 20 个完整异常，断言 PDF 长度大于单异常版本且仍以 `%PDF` 开头；不解析或比较二进制快照。

- [ ] **Step 5: 运行 PDF 测试**

Run:

```bash
cd server
gofmt -w internal/library/reportpdf/*.go
go test ./internal/library/reportpdf -v
```

Expected: PASS。

- [ ] **Step 6: 提交 PDF 生成器**

```bash
git add server/internal/library/reportpdf
git commit -m "feat: generate monthly checkup pdf"
```

## Task 5: PDF 会员接口与越权保护

**Files:**
- Modify: `server/internal/logic/compliance/report/generate.go`
- Modify: `server/internal/logic/compliance/report/repository.go`
- Modify: `server/internal/service/compliance_assistant.go`
- Modify: `server/api/member/member_compliance_assistant.go`
- Modify: `server/internal/controller/member/compliance_assistant.go`
- Modify: `server/internal/logic/compliance/report/report_test.go`

- [ ] **Step 1: 写服务层导出失败测试**

使用内存仓储实现 `Get`，验证请求其他会员的报告返回「报告不存在或无权访问」，合法报告返回 `%PDF` 与文件名 `monthly-checkup-2026-07-v2.pdf`。

- [ ] **Step 2: 运行测试确认 ExportPDF 不存在**

Run: `cd server && go test ./internal/logic/compliance/report -run TestExportPDF -v`  
Expected: FAIL，包含 `ExportPDF undefined`。

- [ ] **Step 3: 增加服务方法**

接口签名：

```go
ExportPDF(ctx context.Context, memberID, id uint64) (filename string, content []byte, err error)
```

实现先调用 `repository.Get(ctx, memberID, id)`，再将持久化快照交给 `reportpdf.Generate`。企业名称通过现有 `shared.LoadOpcByMember` 查询，不使用前端提交值。

- [ ] **Step 4: 增加 API 定义和控制器**

```go
type ComplianceReportPdfReq struct {
	g.Meta `path:"/compliance/reports/{id}/pdf" method:"get" tags:"会员合规助手" summary:"导出月度体检 PDF"`
	Id uint64 `p:"id" in:"path" v:"required|min:1#请指定报告"`
}
type ComplianceReportPdfRes struct{}
```

控制器设置：

```go
r := g.RequestFromCtx(ctx)
r.Response.Header().Set("Content-Type", "application/pdf")
r.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
r.Response.Write(content)
```

- [ ] **Step 5: 运行相关 Go 测试和编译**

Run:

```bash
cd server
gofmt -w api/member/member_compliance_assistant.go internal/controller/member/compliance_assistant.go internal/service/compliance_assistant.go internal/logic/compliance/report/*.go
go test ./internal/logic/compliance/report ./internal/library/reportpdf -v
go test ./internal/controller/member ./api/member ./internal/service
go build ./...
```

Expected: PASS。

- [ ] **Step 6: 提交 PDF 接口**

```bash
git add server/api/member/member_compliance_assistant.go server/internal/controller/member/compliance_assistant.go server/internal/service/compliance_assistant.go server/internal/logic/compliance/report
git commit -m "feat: export monthly checkup reports as pdf"
```

## Task 6: 前端结构化类型与展示组件

**Files:**
- Modify: `web/src/api/frontend/compliance/assistant.ts`
- Create: `web/src/views/frontend/compliance/reports/ReportSummary.vue`
- Create: `web/src/views/frontend/compliance/reports/ReportCategories.vue`
- Create: `web/src/views/frontend/compliance/reports/ReportAnomalyList.vue`

- [ ] **Step 1: 增加 API 类型和 Blob 导出方法**

定义与后端同名的 `ReportSummary`、`ReportCheck`、`ReportCategory`、`ReportAnomaly`、`StructuredReport`，并扩展：

```ts
export interface ComplianceReport {
  id: number
  periodKey: string
  version: number
  status: 'draft' | 'published'
  content: string
  structuredReport?: StructuredReport
  legacy: boolean
  aiModel: string
  publishedAt: number
  createdAt: number
}
```

导出方法：

```ts
export function exportComplianceReportPdf(id: number) {
  return memberRequest.get<Blob>({
    url: `/compliance/reports/${id}/pdf`,
    responseType: 'blob'
  })
}
```

- [ ] **Step 2: 实现 ReportSummary**

Props 为 `summary: ReportSummary`；显示结论中文映射、完整度和高/中/低数量。`dataNotice` 作为明确提示，不把 0% 显示为正常。

- [ ] **Step 3: 实现 ReportCategories**

Props 为 `categories: ReportCategory[]`；使用 Element Plus Tabs 或折叠面板展示固定六类，每个检查项显示状态和 message。无检查项时显示分类 summary。

- [ ] **Step 4: 实现 ReportAnomalyList**

Props 为 `reportId` 和 `anomalies`；每项展示风险等级、标题、事实、依据、可能影响、处理建议、所需材料、期限和人工复核。操作事件：

```ts
defineEmits<{
  addMaterials: [anomaly: ReportAnomaly]
  manualReview: [anomaly: ReportAnomaly]
}>()
```

- [ ] **Step 5: 单文件格式与检查**

Run:

```bash
cd web
corepack pnpm exec prettier --write src/api/frontend/compliance/assistant.ts src/views/frontend/compliance/reports/*.vue
corepack pnpm exec eslint src/api/frontend/compliance/assistant.ts src/views/frontend/compliance/reports/*.vue
corepack pnpm exec stylelint src/views/frontend/compliance/reports/*.vue
```

Expected: 退出码 0。

- [ ] **Step 6: 提交展示组件**

```bash
git add web/src/api/frontend/compliance/assistant.ts web/src/views/frontend/compliance/reports
git commit -m "feat: add structured checkup report components"
```

## Task 7: 月度体检页面编排与下载

**Files:**
- Modify: `web/src/views/frontend/compliance/reports.vue`

- [ ] **Step 1: 替换纯文本卡片布局**

保留账期选择和版本列表；每个版本头部显示账期、版本、中文状态、生成时间、导出 PDF、草稿发布按钮。结构化版本依次使用 `ReportSummary`、异常列表、分类检查；历史版本显示 `content` 和「历史版本仅包含摘要」提示。

- [ ] **Step 2: 增加独立操作状态**

```ts
const loading = ref(false)
const creating = ref(false)
const publishingId = ref<number>()
const exportingId = ref<number>()
```

每个 `try/finally` 只恢复对应状态。账期为空不请求；生成失败不清空现有列表。

- [ ] **Step 3: 实现鉴权 Blob 下载**

```ts
async function exportPdf(item: ComplianceReport) {
  exportingId.value = item.id
  try {
    const blob = await exportComplianceReportPdf(item.id)
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `monthly-checkup-${item.periodKey}-v${item.version}.pdf`
    link.click()
    URL.revokeObjectURL(url)
  } finally {
    exportingId.value = undefined
  }
}
```

- [ ] **Step 4: 实现异常跳转**

「补充材料」跳转 `/user/compliance/documents`；「申请人工复核」跳转 `/user/compliance/tickets`，query 携带 `reportId`、`anomalyCode` 和标题，不在当前任务内自动创建付费工单。

- [ ] **Step 5: 增加响应式样式**

桌面端摘要三列；小于 768px 改为单列。报告头部操作允许换行，异常详情使用 `grid-template-columns: minmax(90px, 130px) 1fr`，移动端改为单列，所有长文本 `overflow-wrap: anywhere`。

- [ ] **Step 6: 页面单文件检查与构建**

Run:

```bash
cd web
corepack pnpm exec prettier --write src/views/frontend/compliance/reports.vue
corepack pnpm exec eslint src/views/frontend/compliance/reports.vue src/views/frontend/compliance/reports/*.vue src/api/frontend/compliance/assistant.ts
corepack pnpm exec stylelint src/views/frontend/compliance/reports.vue src/views/frontend/compliance/reports/*.vue
export PATH="$HOME/.nvm/versions/node/v22.22.3/bin:$PATH"
corepack pnpm build
```

Expected: 单文件检查和生产构建通过；仅允许记录项目已有的 chunk 大小警告。

- [ ] **Step 7: 提交页面改造**

```bash
git add web/src/views/frontend/compliance/reports.vue
git commit -m "feat: redesign monthly checkup report page"
```

## Task 8: 接口文档、集成验证与部署准备

**Files:**
- Modify: `docs/06-OPC合规助手-接口与部署.md`
- Modify: `docs/07-OPC合规助手-验收清单.md`

- [ ] **Step 1: 更新 API 文档**

补充 `GET /compliance/reports/{id}/pdf` 的路径参数、响应头、错误示例；补充报告列表的 `structuredReport` 和 `legacy` 返回字段；明确历史报告降级规则。

- [ ] **Step 2: 更新验收清单**

加入空数据、有异常、无异常、历史报告、草稿导出、已发布导出、越权导出、375px 和桌面宽度验收项。

- [ ] **Step 3: 执行后端全量相关验证**

Run:

```bash
cd server
gofmt -w internal/logic/compliance/report internal/library/reportpdf api/member/member_compliance_assistant.go internal/controller/member/compliance_assistant.go internal/service/compliance_assistant.go
go test ./internal/logic/compliance/report ./internal/library/reportpdf ./internal/controller/member ./api/member ./internal/service -v
go vet ./internal/logic/compliance/report ./internal/library/reportpdf
go build ./...
```

Expected: PASS。

- [ ] **Step 4: 执行前端最终验证**

Run:

```bash
cd web
export PATH="$HOME/.nvm/versions/node/v22.22.3/bin:$PATH"
corepack pnpm exec eslint src/views/frontend/compliance/reports.vue src/views/frontend/compliance/reports/*.vue src/api/frontend/compliance/assistant.ts
corepack pnpm exec stylelint src/views/frontend/compliance/reports.vue src/views/frontend/compliance/reports/*.vue
corepack pnpm build
```

Expected: PASS。

- [ ] **Step 5: 本地迁移与接口验收**

在非生产数据库执行 1.5.7 迁移；用同一会员生成报告、获取列表、发布、导出；确认 PDF 头为 `%PDF`。另用不同会员请求相同报告 ID，确认返回「报告不存在或无权访问」。禁止删除或修改现有报告数据。

- [ ] **Step 6: 浏览器验收**

在 375px 与桌面宽度检查：摘要不重叠、六类可读、异常全部字段可见、长文本换行、按钮可操作、历史报告不空白、PDF 可下载。

- [ ] **Step 7: 最终差异与文档提交**

Run:

```bash
git diff --check
git status --short --branch
```

Expected: 无空白错误，仅包含计划内文件。

```bash
git add docs/06-OPC合规助手-接口与部署.md docs/07-OPC合规助手-验收清单.md
git commit -m "docs: document structured monthly checkup reports"
```

## 部署与回滚

部署顺序：备份数据库结构与线上前端 dist；执行 1.5.7 迁移；发布新二进制和前端资源；重启服务；检查 `/tax_filing/`、报告列表、生成、发布和 PDF 导出。`structured_json` 为可空字段，因此代码回滚后旧版本仍可读取原字段。

代码回滚使用对应提交的 `git revert` 或恢复上一个发布包。数据库默认不回退；若必须删除新字段，需要先确认没有新客户端和新报告依赖，并单独获得生产数据库变更确认，禁止删除报告数据。
