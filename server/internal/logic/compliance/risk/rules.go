package risk

type Facts struct {
	Revenue, SalesDocumentAmount, InvoiceAmount, ReceiptAmount                 *float64
	Expense, CostDocumentAmount, PersonalTransferCount, PersonalTransferAmount *float64
	HighPriorityTaskDueDays, EmployeeCount, RollingSales, TaxpayerThreshold    *float64
	ZeroFiled, InAnnualReportWindow, AnnualReportConfirmed                     *bool
	AddressChanged, AddressProfileUpdated, EmploymentProfileComplete           *bool
	GapThreshold, ThresholdWarningRatio                                        float64
}

type Finding struct {
	Code, Severity, Summary string
	Evidence                map[string]any
}

func Scan(f Facts) []Finding {
	items := make([]Finding, 0, 10)
	add := func(hit bool, code, severity, summary string, evidence map[string]any) {
		if hit {
			items = append(items, Finding{code, severity, summary, evidence})
		}
	}
	add(f.Revenue != nil && f.SalesDocumentAmount != nil && *f.Revenue > 0 && *f.SalesDocumentAmount == 0, "REV_NO_SALES_DOC", "medium", "有经营收入但缺少销项或未开票收入记录", map[string]any{"revenue": value(f.Revenue)})
	gap := f.GapThreshold
	if gap <= 0 {
		gap = .2
	}
	add(f.InvoiceAmount != nil && f.ReceiptAmount != nil && *f.InvoiceAmount > 0 && abs(*f.InvoiceAmount-*f.ReceiptAmount)/(*f.InvoiceAmount) > gap, "INVOICE_RECEIPT_GAP", "medium", "开票金额与回款金额差异较大", nil)
	add(f.ZeroFiled != nil && f.Revenue != nil && *f.ZeroFiled && *f.Revenue > 0, "ZERO_WITH_CASHFLOW", "high", "零申报期间存在经营流水", nil)
	add(f.PersonalTransferCount != nil && f.PersonalTransferAmount != nil && (*f.PersonalTransferCount >= 3 || *f.PersonalTransferAmount >= 10000), "COMPANY_TO_PERSONAL", "high", "公司向个人账户转账频繁或金额较大", nil)
	add(f.Expense != nil && f.CostDocumentAmount != nil && *f.Expense > 0 && *f.CostDocumentAmount < *f.Expense, "COST_DOC_MISSING", "medium", "成本费用凭证不完整", nil)
	add(f.HighPriorityTaskDueDays != nil && *f.HighPriorityTaskDueDays >= 0 && *f.HighPriorityTaskDueDays <= 3, "TASK_DUE_SOON", "medium", "高风险任务即将截止", nil)
	add(f.InAnnualReportWindow != nil && f.AnnualReportConfirmed != nil && *f.InAnnualReportWindow && !*f.AnnualReportConfirmed, "ANNUAL_REPORT_UNCONFIRMED", "high", "工商年报尚未确认", nil)
	add(f.AddressChanged != nil && f.AddressProfileUpdated != nil && *f.AddressChanged && !*f.AddressProfileUpdated, "ADDRESS_CHANGED", "medium", "经营地址变化但企业档案未更新", nil)
	add(f.EmployeeCount != nil && f.EmploymentProfileComplete != nil && *f.EmployeeCount > 0 && !*f.EmploymentProfileComplete, "EMPLOYEE_PROFILE_GAP", "high", "新增员工但用工资料不完整", nil)
	ratio := f.ThresholdWarningRatio
	if ratio <= 0 {
		ratio = .9
	}
	add(f.RollingSales != nil && f.TaxpayerThreshold != nil && *f.TaxpayerThreshold > 0 && *f.RollingSales >= *f.TaxpayerThreshold*ratio, "SALES_THRESHOLD_NEAR", "high", "滚动销售额接近纳税人类型阈值", nil)
	return items
}

func value(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
