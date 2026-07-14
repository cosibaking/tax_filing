package risk

import "testing"

func TestRulesPositiveAndNegative(t *testing.T) {
	tests := []struct {
		code               string
		positive, negative Facts
	}{
		{"REV_NO_SALES_DOC", Facts{Revenue: num(100), SalesDocumentAmount: num(0)}, Facts{Revenue: num(100), SalesDocumentAmount: num(100)}},
		{"INVOICE_RECEIPT_GAP", Facts{InvoiceAmount: num(200), ReceiptAmount: num(100), GapThreshold: 0.2}, Facts{InvoiceAmount: num(100), ReceiptAmount: num(90), GapThreshold: 0.2}},
		{"ZERO_WITH_CASHFLOW", Facts{ZeroFiled: flag(true), Revenue: num(1)}, Facts{ZeroFiled: flag(false), Revenue: num(1)}},
		{"COMPANY_TO_PERSONAL", Facts{PersonalTransferCount: num(5), PersonalTransferAmount: num(50000)}, Facts{PersonalTransferCount: num(1), PersonalTransferAmount: num(100)}},
		{"COST_DOC_MISSING", Facts{Expense: num(100), CostDocumentAmount: num(0)}, Facts{Expense: num(100), CostDocumentAmount: num(100)}},
		{"TASK_DUE_SOON", Facts{HighPriorityTaskDueDays: num(1)}, Facts{HighPriorityTaskDueDays: num(8)}},
		{"ANNUAL_REPORT_UNCONFIRMED", Facts{InAnnualReportWindow: flag(true), AnnualReportConfirmed: flag(false)}, Facts{InAnnualReportWindow: flag(true), AnnualReportConfirmed: flag(true)}},
		{"ADDRESS_CHANGED", Facts{AddressChanged: flag(true), AddressProfileUpdated: flag(false)}, Facts{AddressChanged: flag(false), AddressProfileUpdated: flag(false)}},
		{"EMPLOYEE_PROFILE_GAP", Facts{EmployeeCount: num(1), EmploymentProfileComplete: flag(false)}, Facts{EmployeeCount: num(1), EmploymentProfileComplete: flag(true)}},
		{"SALES_THRESHOLD_NEAR", Facts{RollingSales: num(480), TaxpayerThreshold: num(500), ThresholdWarningRatio: 0.9}, Facts{RollingSales: num(200), TaxpayerThreshold: num(500), ThresholdWarningRatio: 0.9}},
	}
	for _, tc := range tests {
		t.Run(tc.code+"_hit", func(t *testing.T) {
			if !hasCode(Scan(tc.positive), tc.code) {
				t.Fatalf("expected %s", tc.code)
			}
		})
		t.Run(tc.code+"_boundary", func(t *testing.T) {
			if hasCode(Scan(tc.negative), tc.code) {
				t.Fatalf("unexpected %s", tc.code)
			}
		})
		t.Run(tc.code+"_missing", func(t *testing.T) {
			if hasCode(Scan(Facts{}), tc.code) {
				t.Fatalf("missing data hit %s", tc.code)
			}
		})
	}
}

func num(v float64) *float64 { return &v }
func flag(v bool) *bool      { return &v }
func hasCode(items []Finding, code string) bool {
	for _, item := range items {
		if item.Code == code {
			return true
		}
	}
	return false
}
