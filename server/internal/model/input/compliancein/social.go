package compliancein

// EmploymentStatusInp 设置用工状态入参
type EmploymentStatusInp struct {
	MemberId         uint64
	EmploymentStatus string
	Ip               string
}

// EmploymentStatusModel 用工状态出参
type EmploymentStatusModel struct {
	EmploymentStatus        string `json:"employmentStatus"`
	EmploymentConfirmedAt   string `json:"employmentConfirmedAt,omitempty"`
	CanChangeToNoEmployee   bool   `json:"canChangeToNoEmployee"`
}

// SocialGuideItem 社保指引列表项
type SocialGuideItem struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// SocialGuideListModel 指引列表出参
type SocialGuideListModel struct {
	List []SocialGuideItem `json:"list"`
}

// SocialGuideDetailInp 指引详情入参
type SocialGuideDetailInp struct {
	Slug string
}

// SocialGuideDetailModel 指引详情出参
type SocialGuideDetailModel struct {
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	ContentMd  string `json:"contentMd"`
	Audience   string `json:"audience"`
}

// SocialConsultCreateInp 创建咨询入参
type SocialConsultCreateInp struct {
	MemberId   uint64
	Category   string
	Question   string
	RegionCode string
	Ip         string
}

// SocialConsultCreateModel 创建咨询出参
type SocialConsultCreateModel struct {
	Id uint64 `json:"id"`
}

// SocialConsultListInp 咨询列表入参（会员）
type SocialConsultListInp struct {
	MemberId uint64
	Page     int
	PageSize int
}

// SocialConsultItem 咨询项
type SocialConsultItem struct {
	Id         uint64 `json:"id"`
	Category   string `json:"category"`
	Question   string `json:"question"`
	RegionCode string `json:"regionCode,omitempty"`
	Status     string `json:"status"`
	Reply      string `json:"reply,omitempty"`
	RepliedAt  string `json:"repliedAt,omitempty"`
	CreatedAt  string `json:"createdAt"`
}

// SocialConsultListModel 咨询列表出参
type SocialConsultListModel struct {
	List     []SocialConsultItem `json:"list"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	Total    int                 `json:"total"`
}

// AdminSocialConsultListInp 管理端咨询列表入参
type AdminSocialConsultListInp struct {
	Page     int
	PageSize int
	Status   string
	Query    string
}

// AdminSocialConsultItem 管理端咨询项
type AdminSocialConsultItem struct {
	Id          uint64 `json:"id"`
	MemberId    uint64 `json:"memberId"`
	MemberName  string `json:"memberName"`
	OpcId       uint64 `json:"opcId"`
	CompanyName string `json:"companyName"`
	Category    string `json:"category"`
	Question    string `json:"question"`
	RegionCode  string `json:"regionCode,omitempty"`
	Status      string `json:"status"`
	Reply       string `json:"reply,omitempty"`
	RepliedAt   string `json:"repliedAt,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

// AdminSocialConsultListModel 管理端咨询列表出参
type AdminSocialConsultListModel struct {
	List     []AdminSocialConsultItem `json:"list"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
	Total    int                      `json:"total"`
}

// AdminSocialConsultReplyInp 顾问回复入参
type AdminSocialConsultReplyInp struct {
	Id      uint64
	Reply   string
	Action  string // reply | close
	AdminId uint64
	Ip      string
}

// AdminSocialConsultReplyModel 顾问回复出参
type AdminSocialConsultReplyModel struct {
	Id     uint64 `json:"id"`
	Status string `json:"status"`
}
