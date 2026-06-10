package compliancein

// MaterialsConfirmations 资料提交授权确认
type MaterialsConfirmations struct {
	Truthful      bool `json:"truthful"`
	UsageConsent  bool `json:"usageConsent"`
	OpcLimitAck   bool `json:"opcLimitAck"`
}

// MaterialsSubmitInp OPC 注册资料提交入参
type MaterialsSubmitInp struct {
	MemberId            uint64
	ProposedNames       []string
	RegisteredCapital   float64
	CapitalTermYears    int
	BusinessTermType    string
	BusinessTermEnd     string
	BusinessScope       string
	RegisterProvince    string
	RegisterCity        string
	RegisterDistrict    string
	RegisterAddress     string
	AddressProofFileId  uint64
	LegalPersonName     string
	IdCard              string
	IdCardValidFrom     string
	IdCardValidTo       string
	Ethnicity           string
	HouseholdAddress    string
	ResidentialAddress  string
	Phone               string
	Email               string
	IdCardFrontFileId   uint64
	IdCardBackFileId    uint64
	EsignAuthorized     bool
	Confirmations       MaterialsConfirmations
	Ip                  string
}

// MaterialsSubmitModel 资料提交出参
type MaterialsSubmitModel struct {
	OpcId  uint64 `json:"opcId"`
	Status string `json:"status"`
}

// OpcSummaryModel OPC 主体摘要
type OpcSummaryModel struct {
	OpcId                 uint64 `json:"opcId,omitempty"`
	Status                string `json:"status"`
	CompanyName           string `json:"companyName,omitempty"`
	CreditCode            string `json:"creditCode,omitempty"`
	PlanTier              string `json:"planTier,omitempty"`
	IsActive              bool   `json:"isActive"`
	EmploymentStatus      string `json:"employmentStatus,omitempty"`
	EmploymentConfirmedAt string `json:"employmentConfirmedAt,omitempty"`
}

// ProgressStepItem 时间轴节点
type ProgressStepItem struct {
	Key    string `json:"key"`
	Label  string `json:"label"`
	Status string `json:"status"` // done / current / pending
	Date   string `json:"date,omitempty"`
}

// MaterialsReadonlySummary 只读资料摘要（脱敏）
type MaterialsReadonlySummary struct {
	LegalPersonName string `json:"legalPersonName,omitempty"`
	Phone           string `json:"phone,omitempty"`
	ProposedName    string `json:"proposedName,omitempty"`
}

// OpcProgressModel 进度时间轴出参
type OpcProgressModel struct {
	OpcStatus          string                    `json:"opcStatus"`
	CompanyName        string                    `json:"companyName,omitempty"`
	CreditCode         string                    `json:"creditCode,omitempty"`
	EstimatedSlaDays   int                       `json:"estimatedSlaDays"`
	Steps              []ProgressStepItem        `json:"steps"`
	MaterialsReadonly  *MaterialsReadonlySummary `json:"materialsReadonly,omitempty"`
	MaterialsSubmitted bool                      `json:"materialsSubmitted"`
	MaterialsEditable  bool                      `json:"materialsEditable"`
	RejectNote         string                    `json:"rejectNote,omitempty"`
	BankAccountMasked  string                    `json:"bankAccountMasked,omitempty"`
	PlanTier           string                    `json:"planTier,omitempty"`
	PlanName           string                    `json:"planName,omitempty"`
	PlanAmount         float64                   `json:"planAmount,omitempty"`
	SignedAt           string                    `json:"signedAt,omitempty"`
}

// MaterialsEntityOverview 已提交资料实体卡片（脱敏）
type MaterialsEntityOverview struct {
	OpcId              uint64 `json:"opcId"`
	CompanyNameMasked  string `json:"companyNameMasked"`
	LegalPersonSummary string `json:"legalPersonSummary"`
	Status             string `json:"status"`
	StatusLabel        string `json:"statusLabel"`
}

// MaterialsOverviewModel 资料概览列表（按 OPC 实体）
type MaterialsOverviewModel struct {
	Entities []MaterialsEntityOverview `json:"entities"`
}

// MaterialsDetailSection 资料详情分组（脱敏）
type MaterialsDetailSection struct {
	Key    string                  `json:"key"`
	Title  string                  `json:"title"`
	Fields []MaterialsSectionField `json:"fields"`
}

// MaterialsEntityDetailModel 单个 OPC 实体完整资料（脱敏）
type MaterialsEntityDetailModel struct {
	OpcId       uint64                   `json:"opcId"`
	Status      string                   `json:"status"`
	StatusLabel string                   `json:"statusLabel"`
	Sections    []MaterialsDetailSection `json:"sections"`
}

// MaterialsSectionField 资料分组字段
type MaterialsSectionField struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// MaterialsSectionDetailModel 资料分组详情（脱敏）
type MaterialsSectionDetailModel struct {
	Key         string                  `json:"key"`
	Title       string                  `json:"title"`
	Status      string                  `json:"status"`
	StatusLabel string                  `json:"statusLabel"`
	Fields      []MaterialsSectionField `json:"fields"`
}

// MaterialsRevealInp 密码验证后查看完整资料
type MaterialsRevealInp struct {
	MemberId uint64
	OpcId    uint64
	Section  string
	Password string
}

// MaterialsRevealSectionModel 完整资料分组
type MaterialsRevealSectionModel struct {
	Key         string                      `json:"key"`
	Title       string                      `json:"title"`
	Fields      []MaterialsSectionField     `json:"fields"`
	Attachments []MaterialsAttachmentReveal `json:"attachments,omitempty"`
}

// MaterialsRevealAllModel 完整资料（全部明细分组）
type MaterialsRevealAllModel struct {
	OpcId    uint64                        `json:"opcId"`
	Sections []MaterialsRevealSectionModel `json:"sections"`
}

// MaterialsAttachmentReveal 完整资料附件项
type MaterialsAttachmentReveal struct {
	Label     string `json:"label"`
	FileId    uint64 `json:"fileId"`
	FileName  string `json:"fileName"`
	MimeType  string `json:"mimeType"`
	AccessUrl string `json:"accessUrl"`
}

// MaterialsRevealModel 完整资料（明文）
type MaterialsRevealModel struct {
	Key         string                      `json:"key"`
	Fields      []MaterialsSectionField     `json:"fields"`
	Attachments []MaterialsAttachmentReveal `json:"attachments,omitempty"`
}

// BankReceiptInp 开户回执提交入参
type BankReceiptInp struct {
	MemberId         uint64
	BankName         string
	BankReceiptFileId uint64
	Ip               string
}

// BankReceiptModel 开户回执出参
type BankReceiptModel struct {
	OpcId  uint64 `json:"opcId"`
	Status string `json:"status"`
}

// OpcTaskListInp 顾问 OPC 任务列表入参
type OpcTaskListInp struct {
	Page     int
	PageSize int
	Status   string
	Query    string
}

// OpcTaskListItem 任务列表项
type OpcTaskListItem struct {
	OpcId                uint64 `json:"opcId"`
	MemberId             uint64 `json:"memberId"`
	MemberName           string `json:"memberName"`
	MemberPhoneMasked    string `json:"memberPhoneMasked"`
	ProposedName         string `json:"proposedName"`
	LegalPersonName      string `json:"legalPersonName"`
	Status               string `json:"status"`
	StatusLabel          string `json:"statusLabel"`
	MaterialsSubmittedAt string `json:"materialsSubmittedAt,omitempty"`
	DaysSinceSubmit      int    `json:"daysSinceSubmit"`
	PlanTier             string `json:"planTier,omitempty"`
}

// OpcTaskListModel 任务列表出参
type OpcTaskListModel struct {
	List     []OpcTaskListItem `json:"list"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	Total    int               `json:"total"`
}

// OpcTaskDetailInp 任务详情入参
type OpcTaskDetailInp struct {
	OpcId uint64
}

// OpcTaskDetailModel 任务详情（脱敏）
type OpcTaskDetailModel struct {
	OpcId                uint64               `json:"opcId"`
	MemberId             uint64               `json:"memberId"`
	MemberName           string               `json:"memberName"`
	MemberPhone          string               `json:"memberPhone"`
	PlanTier             string               `json:"planTier,omitempty"`
	Status               string               `json:"status"`
	StatusLabel          string               `json:"statusLabel"`
	RejectNote           string               `json:"rejectNote,omitempty"`
	MaterialsSubmittedAt string               `json:"materialsSubmittedAt,omitempty"`
	Company              OpcCompanyDetail     `json:"company"`
	LegalPerson OpcLegalPersonDetail   `json:"legalPerson"`
	Business    OpcBusinessResult      `json:"business"`
	Tax         OpcTaxDetail           `json:"tax"`
	Bank        OpcBankDetail          `json:"bank"`
	AllowedActions []string              `json:"allowedActions"`
	ProgressLogs   []OpcProgressLogItem  `json:"progressLogs"`
}

type OpcCompanyDetail struct {
	ProposedNames      []string `json:"proposedNames"`
	RegisteredCapital  float64  `json:"registeredCapital"`
	CapitalTermYears   int      `json:"capitalTermYears"`
	BusinessTermType   string   `json:"businessTermType"`
	BusinessTermEnd    string   `json:"businessTermEnd,omitempty"`
	BusinessScope      string   `json:"businessScope"`
	RegisterProvince   string   `json:"registerProvince"`
	RegisterCity       string   `json:"registerCity"`
	RegisterDistrict   string   `json:"registerDistrict"`
	RegisterAddress    string   `json:"registerAddress"`
	AddressProofFileId uint64   `json:"addressProofFileId,omitempty"`
}

type OpcLegalPersonDetail struct {
	LegalPersonName    string `json:"legalPersonName"`
	IdCardMasked       string `json:"idCardMasked"`
	IdCardValidFrom    string `json:"idCardValidFrom,omitempty"`
	IdCardValidTo      string `json:"idCardValidTo,omitempty"`
	Ethnicity          string `json:"ethnicity,omitempty"`
	HouseholdAddress   string `json:"householdAddress"`
	ResidentialAddress string `json:"residentialAddress"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	IdCardFrontFileId  uint64 `json:"idCardFrontFileId,omitempty"`
	IdCardBackFileId   uint64 `json:"idCardBackFileId,omitempty"`
	EsignAuthorized    bool   `json:"esignAuthorized"`
}

type OpcBusinessResult struct {
	CompanyName   string `json:"companyName,omitempty"`
	CreditCode    string `json:"creditCode,omitempty"`
	EstablishedAt string `json:"establishedAt,omitempty"`
	LicenseFileId uint64 `json:"licenseFileId,omitempty"`
}

type OpcTaxDetail struct {
	TaxpayerType   string `json:"taxpayerType"`
	TaxActivatedAt string `json:"taxActivatedAt,omitempty"`
}

type OpcBankDetail struct {
	BankName          string `json:"bankName,omitempty"`
	BankAccountMasked string `json:"bankAccountMasked,omitempty"`
	BankReceiptFileId uint64 `json:"bankReceiptFileId,omitempty"`
}

type OpcProgressLogItem struct {
	Step      string `json:"step"`
	Status    string `json:"status"`
	Note      string `json:"note,omitempty"`
	OperatedBy uint64 `json:"operatedBy,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// OpcTaskActionInp 顾问推进任务入参
type OpcTaskActionInp struct {
	OpcId       uint64
	Action      string
	Note        string
	AdminId     uint64
	Ip          string
	// issue_license payload
	CompanyName   string
	CreditCode    string
	EstablishedAt string
	LicenseFileId uint64
	// complete_tax payload
	TaxActivatedAt string
	// complete_bank payload
	BankAccount    string
	BankReceiptFileId uint64
}

// OpcTaskActionModel 顾问推进出参
type OpcTaskActionModel struct {
	OpcId  uint64 `json:"opcId"`
	Status string `json:"status"`
}

// ComplianceCustomerItem 合规客户列表项
type ComplianceCustomerItem struct {
	MemberId          uint64 `json:"memberId"`
	MemberName        string `json:"memberName"`
	MemberPhoneMasked string `json:"memberPhoneMasked"`
	PlanTier          string `json:"planTier,omitempty"`
	OrderStatus       string `json:"orderStatus,omitempty"`
	OpcId             uint64 `json:"opcId,omitempty"`
	OpcStatus         string `json:"opcStatus,omitempty"`
	OpcStatusLabel    string `json:"opcStatusLabel,omitempty"`
	OpcCompanyName    string `json:"opcCompanyName,omitempty"`
	SignedAt          string `json:"signedAt,omitempty"`
}

// ComplianceCustomerListInp 合规客户列表入参
type ComplianceCustomerListInp struct {
	Page      int
	PageSize  int
	Query     string
	OpcStatus string
}

// ComplianceCustomerListModel 合规客户列表出参
type ComplianceCustomerListModel struct {
	List     []ComplianceCustomerItem `json:"list"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
	Total    int                      `json:"total"`
}
