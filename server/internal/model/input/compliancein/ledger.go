package compliancein

// ProfitQueryInp 利润表查询入参
type ProfitQueryInp struct {
	MemberId uint64
	Period   string // monthly / quarterly / yearly
	Year     int
	Month    int
	Quarter  int
}

// ProfitPeriodItem 利润期间项
type ProfitPeriodItem struct {
	Label            string  `json:"label"`
	Year             int     `json:"year"`
	Month            int     `json:"month,omitempty"`
	Quarter          int     `json:"quarter,omitempty"`
	Revenue          float64 `json:"revenue"`
	Cost             float64 `json:"cost"`
	Profit           float64 `json:"profit"`
	CumulativeProfit float64 `json:"cumulativeProfit,omitempty"`
}

// ProfitSummaryModel 利润表出参
type ProfitSummaryModel struct {
	Period             string             `json:"period"`
	Year               int                `json:"year"`
	Items              []ProfitPeriodItem `json:"items"`
	Revenue            float64            `json:"revenue"`
	Cost               float64            `json:"cost"`
	Profit             float64            `json:"profit"`
	CumulativeProfit   float64            `json:"cumulativeProfit,omitempty"`
}

// IncomeListInp 收入台账列表入参
type IncomeListInp struct {
	MemberId uint64
	Page     int
	PageSize int
	Month    string
	Platform string
	Category string
}

// IncomeItem 收入台账条目
type IncomeItem struct {
	Id          uint64  `json:"id"`
	Platform    string  `json:"platform"`
	Category    string  `json:"category"`
	GrossAmount float64 `json:"grossAmount"`
	PlatformFee float64 `json:"platformFee"`
	NetAmount   float64 `json:"netAmount"`
	OccurredAt  string  `json:"occurredAt"`
	Source      string  `json:"source"`
	Remark      string  `json:"remark,omitempty"`
}

// IncomeListModel 收入台账列表出参
type IncomeListModel struct {
	List     []IncomeItem  `json:"list"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
	Total    int           `json:"total"`
	Summary  IncomeSummary `json:"summary"`
}

// IncomeSummary 收入汇总
type IncomeSummary struct {
	TotalGross float64 `json:"totalGross"`
	TotalFee   float64 `json:"totalFee"`
	TotalNet   float64 `json:"totalNet"`
}

// IncomeCreateInp 创建收入入参
type IncomeCreateInp struct {
	MemberId    uint64
	Platform    string
	Category    string
	GrossAmount float64
	PlatformFee float64
	OccurredAt  string
	Remark      string
	Ip          string
}

// IncomeCreateModel 创建收入出参
type IncomeCreateModel struct {
	Id          uint64  `json:"id"`
	Platform    string  `json:"platform"`
	Category    string  `json:"category"`
	GrossAmount float64 `json:"grossAmount"`
	PlatformFee float64 `json:"platformFee"`
	NetAmount   float64 `json:"netAmount"`
	OccurredAt  string  `json:"occurredAt"`
	Source      string  `json:"source"`
}

// IncomeImportInp CSV 导入收入入参
type IncomeImportInp struct {
	MemberId   uint64
	CsvContent string
	Ip         string
}

// IncomeImportRowError 导入行错误
type IncomeImportRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// IncomeImportModel CSV 导入出参
type IncomeImportModel struct {
	Imported int                    `json:"imported"`
	Failed   int                    `json:"failed"`
	Errors   []IncomeImportRowError `json:"errors,omitempty"`
}

// ExpenseListInp 费用台账列表入参
type ExpenseListInp struct {
	MemberId uint64
	Page     int
	PageSize int
	Month    string
	Category string
}

// ExpenseItem 费用台账条目
type ExpenseItem struct {
	Id           uint64  `json:"id"`
	Category     string  `json:"category"`
	CategoryName string  `json:"categoryName,omitempty"`
	Amount       float64 `json:"amount"`
	InvoiceType  string  `json:"invoiceType"`
	AttachmentId uint64  `json:"attachmentId,omitempty"`
	Description  string  `json:"description,omitempty"`
	WarningFlag  bool    `json:"warningFlag"`
	OccurredAt   string  `json:"occurredAt"`
}

// ExpenseListModel 费用台账列表出参
type ExpenseListModel struct {
	List     []ExpenseItem  `json:"list"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Total    int            `json:"total"`
	Summary  ExpenseSummary `json:"summary"`
}

// ExpenseSummary 费用汇总
type ExpenseSummary struct {
	TotalAmount float64 `json:"totalAmount"`
	CostRatio   float64 `json:"costRatio"`
}

// ExpenseCreateInp 创建费用入参
type ExpenseCreateInp struct {
	MemberId       uint64
	Category       string
	Amount         float64
	InvoiceType    string
	AttachmentId   uint64
	Description    string
	OccurredAt     string
	ConfirmWarning bool
	Ip             string
}

// ExpenseCreateModel 创建费用出参
type ExpenseCreateModel struct {
	Id           uint64  `json:"id"`
	Category     string  `json:"category"`
	Amount       float64 `json:"amount"`
	InvoiceType  string  `json:"invoiceType"`
	AttachmentId uint64  `json:"attachmentId,omitempty"`
	Description  string  `json:"description,omitempty"`
	WarningFlag  bool    `json:"warningFlag"`
	WarningMsg   string  `json:"warningMsg,omitempty"`
	OccurredAt   string  `json:"occurredAt"`
	CostRatio    float64 `json:"costRatio"`
}

// ExpenseTypeItem 费用类型库条目
type ExpenseTypeItem struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	VoucherHint   string `json:"voucherHint"`
	ComplianceTip string `json:"complianceTip,omitempty"`
}

// ExpenseTypesModel 费用类型库出参
type ExpenseTypesModel struct {
	List []ExpenseTypeItem `json:"list"`
}

// BankImportInp 银行流水导入入参
type BankImportInp struct {
	MemberId   uint64
	CsvContent string
	Ip         string
}

// BankImportModel 银行流水导入出参
type BankImportModel struct {
	Imported  int                    `json:"imported"`
	Failed    int                    `json:"failed"`
	Matched   int                    `json:"matched"`
	Unmatched int                    `json:"unmatched"`
	Errors    []IncomeImportRowError `json:"errors,omitempty"`
}

// IncomeDeleteInp 删除收入入参
type IncomeDeleteInp struct {
	MemberId uint64
	Id       uint64
	Ip       string
}

// ExpenseDeleteInp 删除费用入参
type ExpenseDeleteInp struct {
	MemberId uint64
	Id       uint64
	Ip       string
}

// IncomeImportPreviewRow CSV 预览行
type IncomeImportPreviewRow struct {
	Row         int     `json:"row"`
	OccurredAt  string  `json:"occurredAt"`
	Platform    string  `json:"platform"`
	Category    string  `json:"category"`
	GrossAmount float64 `json:"grossAmount"`
	PlatformFee float64 `json:"platformFee"`
	Valid       bool    `json:"valid"`
	Error       string  `json:"error,omitempty"`
}

// IncomeImportPreviewModel CSV 预览出参
type IncomeImportPreviewModel struct {
	Rows         []IncomeImportPreviewRow `json:"rows"`
	ValidCount   int                    `json:"validCount"`
	InvalidCount int                    `json:"invalidCount"`
}

// BankUnmatchedInp 未匹配流水查询入参
type BankUnmatchedInp struct {
	MemberId uint64
	Month    string
}

// BankUnmatchedItem 未匹配流水项
type BankUnmatchedItem struct {
	Id          uint64  `json:"id"`
	OccurredAt  string  `json:"occurredAt"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
}

// BankUnmatchedModel 未匹配流水出参
type BankUnmatchedModel struct {
	List  []BankUnmatchedItem `json:"list"`
	Count int                 `json:"count"`
}

// LedgerVoucherListInp 会计分录列表入参
type LedgerVoucherListInp struct {
	MemberId uint64
	Period   string
}

// LedgerVoucherItem 会计分录项
type LedgerVoucherItem struct {
	Id            uint64  `json:"id"`
	OccurredAt    string  `json:"occurredAt"`
	Summary       string  `json:"summary"`
	DebitAccount  string  `json:"debitAccount"`
	CreditAccount string  `json:"creditAccount"`
	Amount        float64 `json:"amount"`
}

// LedgerVoucherListModel 会计分录列表出参
type LedgerVoucherListModel struct {
	List []LedgerVoucherItem `json:"list"`
}
