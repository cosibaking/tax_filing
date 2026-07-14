package reportpdf

import (
	"strings"
	"testing"
)

func TestGenerateStructuredReport(t *testing.T) {
	pdf, err := Generate(sampleData())
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	assertPDF(t, pdf)
}

func TestLocalizedReportLabelsAndRuleVersion(t *testing.T) {
	if got := reportStatusLabel("draft"); got != "草稿" {
		t.Fatalf("draft report status = %q, want 草稿", got)
	}
	if got := reportStatusLabel("published"); got != "已发布" {
		t.Fatalf("published report status = %q, want 已发布", got)
	}
	for input, want := range map[string]string{"high": "高风险", "medium": "中风险", "low": "低风险"} {
		if got := riskLevelLabel(input); got != want {
			t.Fatalf("risk level %q = %q, want %q", input, got, want)
		}
	}
	for input, want := range map[string]string{"normal": "正常", "attention": "需关注", "urgent": "紧急处理", "insufficient": "数据不足"} {
		if got := checkStatusLabel(input); got != want {
			t.Fatalf("check status %q = %q, want %q", input, got, want)
		}
	}

	fields := anomalyDetailFields(Anomaly{Severity: "high", RuleVersion: "risk-v2026.07"})
	joined := strings.Join(fields, "\n")
	if !strings.Contains(joined, "等级：高风险") || !strings.Contains(joined, "规则版本：risk-v2026.07") {
		t.Fatalf("localized anomaly detail is incomplete:\n%s", joined)
	}
	missingVersion := strings.Join(anomalyDetailFields(Anomaly{}), "\n")
	if !strings.Contains(missingVersion, "规则版本：未记录，请结合原始资料人工复核") {
		t.Fatalf("missing rule version fallback is not cautious:\n%s", missingVersion)
	}
}

func TestGenerateLegacySummary(t *testing.T) {
	pdf, err := Generate(Data{CompanyName: "历史测试企业", Version: 1, Status: "generated", GeneratedAt: "2026-07-14 10:00:00", LegacyContent: "历史体检结论：请补充发票与银行回单。"})
	if err != nil {
		t.Fatalf("Generate legacy summary failed: %v", err)
	}
	assertPDF(t, pdf)
}

func TestGeneratePaginatesLongAnomalies(t *testing.T) {
	onePDF, err := Generate(sampleData())
	if err != nil {
		t.Fatalf("Generate one anomaly failed: %v", err)
	}
	many := sampleData()
	many.Anomalies = make([]Anomaly, 20)
	longText := strings.Repeat("该异常涉及多笔经营流水、合同、发票和银行回单，需要逐项核对原始资料并记录复核结论。", 12)
	for i := range many.Anomalies {
		many.Anomalies[i] = Anomaly{Code: "risk", CategoryCode: "tax", Title: "长期未处理的申报差异", Severity: "high", Facts: longText, Basis: longText, Impact: longText, Recommendation: longText, RequiredMaterials: []string{"合同", "发票", "银行回单", longText}, DueDate: "2026-07-31", RequiresManualReview: true}
	}
	manyPDF, err := Generate(many)
	if err != nil {
		t.Fatalf("Generate many anomalies failed: %v", err)
	}
	assertPDF(t, manyPDF)
	if len(manyPDF) <= len(onePDF)+10*1024 {
		t.Fatalf("paginated PDF is not significantly larger: one=%d many=%d", len(onePDF), len(manyPDF))
	}
}

func TestRenderRepeatsAnomalyHeadingAfterPageBreak(t *testing.T) {
	data := sampleData()
	longText := strings.Repeat("单个异常的事实和处理依据需要持续跨页展示，并在续页明确标识当前异常。", 500)
	data.Anomalies[0].Facts = longText

	pdf, stats, err := generateWithStats(data)
	if err != nil {
		t.Fatalf("Generate long anomaly failed: %v", err)
	}
	assertPDF(t, pdf)
	if stats.pageCount <= 1 {
		t.Fatalf("expected multiple pages, got %d", stats.pageCount)
	}
	if stats.anomalyContinuationCount <= 0 {
		t.Fatal("expected anomaly continuation heading after page break")
	}
}

func TestRenderWrapsLongContinuationHeading(t *testing.T) {
	data := sampleData()
	data.Anomalies[0].Title = strings.Repeat("这是需要在续页中保持可读且不能越过页面宽度的超长异常标题", 20) + "\n第二段标题"
	data.Anomalies[0].Facts = strings.Repeat("异常事实需要跨越多个页面。", 500)

	pdf, stats, err := generateWithStats(data)
	if err != nil {
		t.Fatalf("Generate long continuation heading failed: %v", err)
	}
	assertPDF(t, pdf)
	if stats.pageCount <= 1 || stats.anomalyContinuationCount <= 0 {
		t.Fatalf("expected continuation pages and headings, stats=%+v", stats)
	}

	lines := continuationHeadingLines("异常 1：" + data.Anomalies[0].Title + "（续）")
	if len(lines) == 0 || len(lines) > maxContinuationHeadingLines {
		t.Fatalf("unexpected continuation line count: %d", len(lines))
	}
	for _, line := range lines {
		if len([]rune(line)) > charsPerLine(continuationHeadingFontSize) {
			t.Fatalf("continuation line exceeds content width estimate: %q", line)
		}
	}
}

func sampleData() Data {
	return Data{
		CompanyName: "示例科技有限公司", PeriodKey: "2026-06", Version: 2, Status: "generated", GeneratedAt: "2026-07-14 10:00:00",
		Summary:    Summary{Conclusion: "urgent", CompletenessRate: 88.5, HighCount: 1, DataNotice: "尚缺少部分银行回单。"},
		Categories: []Category{{Code: "tax", Name: "税务与申报", Status: "urgent", Summary: "存在一项高风险。", Checks: []Check{{Code: "tax-overview", Name: "税务检查", Status: "urgent", Message: "请尽快处理。"}}}},
		Anomalies:  []Anomaly{{Code: "tax-1", CategoryCode: "tax", Title: "申报收入与流水不一致", Severity: "high", Facts: "申报收入低于银行流水。", Basis: "按月核对申报表与经营流水。", Impact: "可能影响申报准确性。", Recommendation: "核对差异并更正申报。", RequiredMaterials: []string{"申报表", "银行回单"}, DueDate: "2026-07-31", RequiresManualReview: true, RuleVersion: "v1"}},
	}
}

func assertPDF(t *testing.T, pdf []byte) {
	t.Helper()
	if len(pdf) <= 1024 {
		t.Fatalf("PDF too small: %d bytes", len(pdf))
	}
	if !strings.HasPrefix(string(pdf), "%PDF") {
		t.Fatal("invalid PDF header")
	}
}
