package statementpdf

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	data := Data{
		CompanyName:      "测试 OPC 有限公司",
		Period:           "2026-05",
		Revenue:          128000.50,
		Cost:             32000.00,
		Profit:           96000.50,
		PrepaidTax:       3840.02,
		CumulativeProfit: 256000.00,
		FilingStatus:     "filed",
		IncomeEntries:    12,
		ExpenseEntries:   5,
		PendingFilings:   0,
		GeneratedAt:      "2026-06-10 10:00:00",
		SentAt:           "2026-06-10 10:05:00",
	}
	pdf, err := Generate(data)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if len(pdf) < 1024 {
		t.Fatalf("PDF too small: %d bytes", len(pdf))
	}
	if pdf[0] != '%' || pdf[1] != 'P' || pdf[2] != 'D' || pdf[3] != 'F' {
		t.Fatalf("invalid PDF header")
	}
}
