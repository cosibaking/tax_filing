package compliancein

// DashboardOverviewModel 合规工作台概览
type DashboardOverviewModel struct {
	ActiveCustomers int                    `json:"activeCustomers"`
	PendingOpcTasks int                    `json:"pendingOpcTasks"`
	PendingFilings  int                    `json:"pendingFilings"`
	OverdueFilings  int                    `json:"overdueFilings"`
	DraftStatements int                    `json:"draftStatements"`
	OpcTasks        []DashboardOpcTaskItem `json:"opcTasks"`
	FilingTasks     []DashboardFilingItem  `json:"filingTasks"`
}

// DashboardOpcTaskItem 待办 OPC 任务
type DashboardOpcTaskItem struct {
	OpcId                uint64 `json:"opcId"`
	MemberName           string `json:"memberName"`
	ProposedName         string `json:"proposedName"`
	Status               string `json:"status"`
	StatusLabel          string `json:"statusLabel"`
	MaterialsSubmittedAt string `json:"materialsSubmittedAt,omitempty"`
	DaysSinceSubmit      int    `json:"daysSinceSubmit"`
}

// DashboardFilingItem 待办申报任务
type DashboardFilingItem struct {
	Id               uint64  `json:"id"`
	MemberName       string  `json:"memberName"`
	CompanyName      string  `json:"companyName"`
	TaxType          string  `json:"taxType"`
	TaxTypeLabel     string  `json:"taxTypeLabel"`
	Period           string  `json:"period"`
	DueDate          string  `json:"dueDate"`
	Status           string  `json:"status"`
	CalculatedAmount float64 `json:"calculatedAmount"`
}
