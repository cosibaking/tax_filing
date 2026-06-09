package compliancein

// TaxFilingImportInp 报税 Excel 导入入参
type TaxFilingImportInp struct {
	MemberId    uint64
	FileContent []byte
	ExcelFileId uint64
	Ip          string
}

// TaxFilingSummaryPreview 申报汇总预览
type TaxFilingSummaryPreview struct {
	Period          string  `json:"period"`
	Revenue         float64 `json:"revenue"`
	ExpenseTotal    float64 `json:"expenseTotal"`
	Profit          float64 `json:"profit"`
	VatAmount       float64 `json:"vatAmount"`
	CitAmount       float64 `json:"citAmount"`
	SurchargeAmount float64 `json:"surchargeAmount"`
	Remark          string  `json:"remark,omitempty"`
	Valid           bool    `json:"valid"`
	Error           string  `json:"error,omitempty"`
}

// TaxFilingExpensePreviewRow 报销明细预览行
type TaxFilingExpensePreviewRow struct {
	Row         int     `json:"row"`
	OccurredAt  string  `json:"occurredAt"`
	Category    string  `json:"category"`
	CategoryName string `json:"categoryName,omitempty"`
	Amount      float64 `json:"amount"`
	InvoiceType string  `json:"invoiceType"`
	Description string  `json:"description,omitempty"`
	Valid       bool    `json:"valid"`
	Error       string  `json:"error,omitempty"`
}

// TaxFilingImportPreviewModel 报税 Excel 预览出参
type TaxFilingImportPreviewModel struct {
	Summary           TaxFilingSummaryPreview      `json:"summary"`
	ExpenseRows       []TaxFilingExpensePreviewRow `json:"expenseRows"`
	ValidExpenseCount int                          `json:"validExpenseCount"`
	InvalidExpenseCount int                        `json:"invalidExpenseCount"`
}

// TaxFilingImportModel 报税 Excel 导入出参
type TaxFilingImportModel struct {
	SubmissionId    uint64 `json:"submissionId"`
	Period          string `json:"period"`
	ExpenseImported int    `json:"expenseImported"`
	ExpenseFailed   int    `json:"expenseFailed"`
}

// TaxFilingSubmissionListInp 申报提交历史入参
type TaxFilingSubmissionListInp struct {
	MemberId uint64
	Page     int
	PageSize int
}

// TaxFilingSubmissionItem 申报提交历史项
type TaxFilingSubmissionItem struct {
	Id              uint64  `json:"id"`
	Period          string  `json:"period"`
	Revenue         float64 `json:"revenue"`
	ExpenseTotal    float64 `json:"expenseTotal"`
	Profit          float64 `json:"profit"`
	VatAmount       float64 `json:"vatAmount"`
	CitAmount       float64 `json:"citAmount"`
	SurchargeAmount float64 `json:"surchargeAmount"`
	ExpenseImported int     `json:"expenseImported"`
	Remark          string  `json:"remark,omitempty"`
	CreatedAt       string  `json:"createdAt"`
}

// TaxFilingSubmissionListModel 申报提交历史出参
type TaxFilingSubmissionListModel struct {
	List     []TaxFilingSubmissionItem `json:"list"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"pageSize"`
	Total    int                       `json:"total"`
}
