// +----------------------------------------------------------------------

// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]

// +----------------------------------------------------------------------

// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.

// +----------------------------------------------------------------------

// | Licensed ( https://opensource.org/licenses/MIT )

// +----------------------------------------------------------------------

// | Author: 喜羊羊 <751300685@qq.com>

// +----------------------------------------------------------------------



package admin



import (

	"github.com/gogf/gf/v2/frame/g"



	"xygo/internal/model/input/compliancein"

)



// ComplianceOpcTaskListReq OPC 任务列表

type ComplianceOpcTaskListReq struct {

	g.Meta   `path:"/admin/compliance/opc-tasks" method:"get" tags:"合规服务" summary:"OPC任务列表"`

	Page     int    `p:"page" json:"page" d:"1"`

	PageSize int    `p:"pageSize" json:"pageSize" d:"20"`

	Status   string `p:"status" json:"status"`

	Query    string `p:"q" json:"q"`

}



// ComplianceOpcTaskListRes OPC 任务列表响应

type ComplianceOpcTaskListRes struct {

	*compliancein.OpcTaskListModel

}



// ComplianceOpcTaskDetailReq OPC 任务详情

type ComplianceOpcTaskDetailReq struct {

	g.Meta `path:"/admin/compliance/opc-tasks/{opcId}" method:"get" tags:"合规服务" summary:"OPC任务详情"`

	OpcId  uint64 `json:"opcId" v:"required|min:1#请指定任务ID"`

}



// ComplianceOpcTaskDetailRes OPC 任务详情响应

type ComplianceOpcTaskDetailRes struct {

	*compliancein.OpcTaskDetailModel

}



// ComplianceOpcTaskActionReq 推进 OPC 任务

type ComplianceOpcTaskActionReq struct {

	g.Meta `path:"/admin/compliance/opc-tasks" method:"patch" tags:"合规服务" summary:"推进OPC任务进度"`

	OpcId  uint64 `json:"opcId" v:"required|min:1#请指定任务ID"`

	Action string `json:"action" v:"required|in:approve_materials,reject_materials,issue_license,complete_tax,complete_bank#请指定操作|操作无效"`

	Note   string `json:"note"`

	// issue_license

	CompanyName   string `json:"companyName"`

	CreditCode    string `json:"creditCode"`

	EstablishedAt string `json:"establishedAt"`

	LicenseFileId uint64 `json:"licenseFileId"`

	// complete_tax

	TaxActivatedAt string `json:"taxActivatedAt"`

	// complete_bank

	BankAccount       string `json:"bankAccount"`

	BankReceiptFileId uint64 `json:"bankReceiptFileId"`

}



// ComplianceOpcTaskActionRes 推进任务响应

type ComplianceOpcTaskActionRes struct {

	*compliancein.OpcTaskActionModel

}



// ComplianceCustomerListReq 合规客户列表

type ComplianceCustomerListReq struct {

	g.Meta   `path:"/admin/compliance/customers" method:"get" tags:"合规服务" summary:"合规客户列表"`

	Page     int    `p:"page" json:"page" d:"1"`

	PageSize int    `p:"pageSize" json:"pageSize" d:"20"`

	Query     string `p:"q" json:"q"`
	OpcStatus string `p:"opcStatus" json:"opcStatus"`

}



// ComplianceCustomerListRes 合规客户列表响应

type ComplianceCustomerListRes struct {

	*compliancein.ComplianceCustomerListModel

}

// ComplianceFilingListReq 申报任务列表
type ComplianceFilingListReq struct {
	g.Meta   `path:"/admin/compliance/filing" method:"get" tags:"合规服务" summary:"申报任务列表"`
	Page     int    `p:"page" json:"page" d:"1"`
	PageSize int    `p:"pageSize" json:"pageSize" d:"20"`
	Status   string `p:"status" json:"status"`
	TaxType  string `p:"taxType" json:"taxType"`
	Period   string `p:"period" json:"period"`
	Query    string `p:"q" json:"q"`
}

type ComplianceFilingListRes struct {
	*compliancein.FilingListModel
}

// ComplianceFilingMarkReq 标记已申报
type ComplianceFilingMarkReq struct {
	g.Meta         `path:"/admin/compliance/filing" method:"patch" tags:"合规服务" summary:"标记已申报"`
	TaskId         uint64                      `json:"taskId" v:"required|min:1#请指定任务"`
	FiledAmount    float64                     `json:"filedAmount"`
	ReportedIncome float64                     `json:"reportedIncome"`
	ReceiptFileId  uint64                      `json:"receiptFileId"`
	Checklist      []compliancein.ChecklistItem `json:"checklist"`
}

type ComplianceFilingMarkRes struct {
	*compliancein.FilingMarkModel
}

// ComplianceStatementSendReq 单主体生成并发送月度对账单（保留兼容）
type ComplianceStatementSendReq struct {
	g.Meta `path:"/admin/compliance/statements/send-one" method:"post" tags:"合规服务" summary:"单主体生成并发送对账单"`
	OpcId  uint64 `json:"opcId" v:"required|min:1#请指定OPC主体"`
	Year   int    `json:"year" v:"required|min:2000#请指定年份"`
	Month  int    `json:"month" v:"required|between:1,12#请指定月份"`
}

type ComplianceStatementSendRes struct {
	*compliancein.StatementSendModel
}

// ComplianceStatementListReq 对账单管理列表
type ComplianceStatementListReq struct {
	g.Meta   `path:"/admin/compliance/statements" method:"get" tags:"合规服务" summary:"对账单管理列表"`
	Page     int    `p:"page" json:"page" d:"1"`
	PageSize int    `p:"pageSize" json:"pageSize" d:"20"`
	Period   string `p:"period" json:"period"`
	Status   string `p:"status" json:"status"`
	Query    string `p:"q" json:"q"`
}

type ComplianceStatementListRes struct {
	*compliancein.StatementListModel
}

// ComplianceStatementGenerateReq 批量生成对账单
type ComplianceStatementGenerateReq struct {
	g.Meta    `path:"/admin/compliance/statements/generate" method:"post" tags:"合规服务" summary:"批量生成对账单"`
	Period    string   `json:"period" v:"required#请指定月份"`
	MemberIds []uint64 `json:"memberIds"`
}

type ComplianceStatementGenerateRes struct {
	*compliancein.StatementGenerateModel
}

// ComplianceStatementNotifyReq 批量发送对账单通知
type ComplianceStatementNotifyReq struct {
	g.Meta `path:"/admin/compliance/statements/send" method:"post" tags:"合规服务" summary:"批量发送对账单通知"`
	Ids    []uint64 `json:"ids" v:"required#请选择对账单"`
}

type ComplianceStatementNotifyRes struct {
	*compliancein.StatementNotifyModel
}

// ComplianceStatementPdfReq 对账单 PDF 下载地址
type ComplianceStatementPdfReq struct {
	g.Meta `path:"/admin/compliance/statements/{id}/pdf" method:"get" tags:"合规服务" summary:"对账单PDF下载地址"`
	Id     uint64 `p:"id" in:"path" json:"id" v:"required|min:1#请指定对账单|对账单无效"`
}

type ComplianceStatementPdfRes struct {
	Url string `json:"url"`
}

// ComplianceAuditExportReq 合规留痕 CSV 导出
type ComplianceAuditExportReq struct {
	g.Meta   `path:"/admin/compliance/audit/export" method:"get" tags:"合规服务" summary:"导出合规留痕CSV"`
	MemberId uint64 `p:"memberId" json:"memberId" v:"required|min:1#请指定会员ID"`
}

type ComplianceAuditExportRes struct{}

// ComplianceDashboardReq 合规工作台概览
type ComplianceDashboardReq struct {
	g.Meta `path:"/admin/compliance/dashboard" method:"get" tags:"合规服务" summary:"合规工作台概览"`
}

type ComplianceDashboardRes struct {
	*compliancein.DashboardOverviewModel
}

// ComplianceSocialConsultListReq 社保咨询工单列表
type ComplianceSocialConsultListReq struct {
	g.Meta   `path:"/admin/compliance/social-consults" method:"get" tags:"合规服务" summary:"社保咨询工单列表"`
	Page     int    `p:"page" json:"page" d:"1"`
	PageSize int    `p:"pageSize" json:"pageSize" d:"20"`
	Status   string `p:"status" json:"status"`
	Query    string `p:"q" json:"q"`
}

type ComplianceSocialConsultListRes struct {
	*compliancein.AdminSocialConsultListModel
}

// ComplianceSocialConsultReplyReq 回复或关闭社保咨询
type ComplianceSocialConsultReplyReq struct {
	g.Meta `path:"/admin/compliance/social-consults/{id}" method:"patch" tags:"合规服务" summary:"回复或关闭社保咨询"`
	Id     uint64 `p:"id" in:"path" v:"required|min:1#请指定工单"`
	Reply  string `json:"reply"`
	Action string `json:"action" d:"reply" v:"in:reply,close#操作无效"`
}

type ComplianceSocialConsultReplyRes struct {
	*compliancein.AdminSocialConsultReplyModel
}

