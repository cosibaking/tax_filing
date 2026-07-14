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
