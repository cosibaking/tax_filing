package compliancein

// TaxCalendarInp 申报日历入参
type TaxCalendarInp struct {
	MemberId uint64
	Year     int
	Month    int
}

// TaxCalendarItem 申报日历项
type TaxCalendarItem struct {
	Id               uint64  `json:"id"`
	TaxType          string  `json:"taxType"`
	TaxTypeLabel     string  `json:"taxTypeLabel"`
	Period           string  `json:"period"`
	DueDate          string  `json:"dueDate"`
	Status           string  `json:"status"`
	CalculatedAmount float64 `json:"calculatedAmount"`
	FiledAmount      float64 `json:"filedAmount,omitempty"`
	AllChecklistDone bool    `json:"allChecklistDone"`
}

// TaxCalendarModel 申报日历出参
type TaxCalendarModel struct {
	Year     int               `json:"year"`
	Month    int               `json:"month,omitempty"`
	Items    []TaxCalendarItem `json:"items"`
	Tasks    []TaxCalendarItem `json:"tasks"`
	DueDates []string          `json:"dueDates,omitempty"`
	NextDueDate    string      `json:"nextDueDate,omitempty"`
	DaysUntilDue   int         `json:"daysUntilDue,omitempty"`
}

// TaxTaskDetailModel 申报任务详情
type TaxTaskDetailModel struct {
	Id               uint64          `json:"id"`
	TaxType          string          `json:"taxType"`
	TaxTypeLabel     string          `json:"taxTypeLabel"`
	Period           string          `json:"period"`
	DueDate          string          `json:"dueDate"`
	Status           string          `json:"status"`
	CalculatedAmount float64         `json:"calculatedAmount"`
	FiledAmount      float64         `json:"filedAmount,omitempty"`
	Detail           TaxTaskCalcDetail `json:"detail,omitempty"`
}

// TaxTaskCalcDetail 税额计算明细
type TaxTaskCalcDetail struct {
	RevenueExTax float64 `json:"revenueExTax"`
	Vat          float64 `json:"vat"`
	Surcharge    float64 `json:"surcharge"`
	Cit          float64 `json:"cit"`
	Total        float64 `json:"total"`
}

// TaxTaskDetailInp 申报任务详情入参
type TaxTaskDetailInp struct {
	MemberId uint64
	TaskId   uint64
}

// MemberStatementListInp 会员对账单列表入参
type MemberStatementListInp struct {
	MemberId uint64
	Page     int
	PageSize int
}

// MemberStatementItem 会员对账单项
type MemberStatementItem struct {
	Id               uint64  `json:"id"`
	Period           string  `json:"period"`
	Revenue          float64 `json:"revenue"`
	Cost             float64 `json:"cost"`
	Profit           float64 `json:"profit"`
	PrepaidTax       float64 `json:"prepaidTax"`
	CumulativeProfit float64 `json:"cumulativeProfit"`
	FilingStatus     string  `json:"filingStatus,omitempty"`
	Status           string  `json:"status"`
	PdfUrl           string  `json:"pdfUrl,omitempty"`
}

// StatementPdfInp 对账单 PDF 入参
type StatementPdfInp struct {
	MemberId    uint64
	StatementId uint64
}

// StatementPdfModel 对账单 PDF 出参
type StatementPdfModel struct {
	Url string `json:"url"`
}

// MemberStatementListModel 会员对账单列表出参
type MemberStatementListModel struct {
	List  []MemberStatementItem `json:"list"`
	Total int                   `json:"total"`
}

// TaxChecklistInp 自查清单入参
type TaxChecklistInp struct {
	MemberId uint64
	TaskId   uint64
	Period   string
}

// ChecklistItem 自查清单项
type ChecklistItem struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Checked bool   `json:"checked"`
	Na      bool   `json:"na,omitempty"`
	Hint    string `json:"hint,omitempty"`
}

// TaxChecklistModel 自查清单出参
type TaxChecklistModel struct {
	TaskId      uint64          `json:"taskId"`
	Period      string          `json:"period"`
	TaxType     string          `json:"taxType"`
	Items       []ChecklistItem `json:"items"`
	AllComplete bool            `json:"allComplete"`
	ReadOnly    bool            `json:"readOnly"`
}

// FilingListInp 申报任务列表入参
type FilingListInp struct {
	Page     int
	PageSize int
	Status   string
	TaxType  string
	Period   string
	Query    string
}

// FilingTaskItem 申报任务列表项
type FilingTaskItem struct {
	Id               uint64  `json:"id"`
	OpcId            uint64  `json:"opcId"`
	MemberId         uint64  `json:"memberId"`
	CompanyName      string  `json:"companyName"`
	MemberName       string  `json:"memberName"`
	EmploymentStatus string  `json:"employmentStatus,omitempty"`
	TaxType          string  `json:"taxType"`
	TaxTypeLabel     string  `json:"taxTypeLabel"`
	Period           string  `json:"period"`
	DueDate          string  `json:"dueDate"`
	Status           string  `json:"status"`
	CalculatedAmount float64 `json:"calculatedAmount"`
	FiledAmount      float64 `json:"filedAmount,omitempty"`
	AllChecklistDone bool    `json:"allChecklistDone"`
}

// FilingListModel 申报任务列表出参
type FilingListModel struct {
	List     []FilingTaskItem `json:"list"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	Total    int              `json:"total"`
}

// FilingMarkInp 标记已申报入参
type FilingMarkInp struct {
	TaskId         uint64
	FiledAmount    float64
	ReportedIncome float64
	ReceiptFileId  uint64
	Checklist      []ChecklistItem
	AdminId        uint64
	Ip             string
}

// FilingMarkModel 标记已申报出参
type FilingMarkModel struct {
	Id     uint64 `json:"id"`
	Status string `json:"status"`
}

// StatementSendInp 生成对账单入参
type StatementSendInp struct {
	OpcId   uint64
	Year    int
	Month   int
	AdminId uint64
	Ip      string
}

// StatementSendModel 生成对账单出参
type StatementSendModel struct {
	Id     uint64 `json:"id"`
	OpcId  uint64 `json:"opcId"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	SentAt string `json:"sentAt"`
}

// StatementListInp 对账单管理列表入参
type StatementListInp struct {
	Page     int
	PageSize int
	Period   string
	Status   string
	Query    string
}

// AdminStatementItem 对账单管理列表项
type AdminStatementItem struct {
	Id             uint64  `json:"id"`
	MemberId       uint64  `json:"memberId"`
	MemberName     string  `json:"memberName"`
	MemberPhone    string  `json:"memberPhone"`
	OpcCompanyName string  `json:"opcCompanyName"`
	Period         string  `json:"period"`
	Revenue        float64 `json:"revenue"`
	Cost           float64 `json:"cost"`
	Profit         float64 `json:"profit"`
	PrepaidTax     float64 `json:"prepaidTax"`
	Status         string  `json:"status"`
	SentAt         string  `json:"sentAt,omitempty"`
	PdfUrl         string  `json:"pdfUrl,omitempty"`
}

// StatementListModel 对账单列表出参
type StatementListModel struct {
	List     []AdminStatementItem `json:"list"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int                  `json:"total"`
}

// StatementGenerateInp 批量生成对账单入参
type StatementGenerateInp struct {
	Period    string
	MemberIds []uint64
	AdminId   uint64
	Ip        string
}

// StatementGenerateModel 批量生成出参
type StatementGenerateModel struct {
	Generated int `json:"generated"`
}

// StatementNotifyInp 批量发送对账单通知入参
type StatementNotifyInp struct {
	Ids     []uint64
	AdminId uint64
	Ip      string
}

// StatementNotifyModel 批量发送出参
type StatementNotifyModel struct {
	Sent int `json:"sent"`
}
