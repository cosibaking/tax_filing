package member

import (
	"github.com/gogf/gf/v2/frame/g"

	compliancedocument "xygo/internal/logic/compliance/document"
	"xygo/internal/logic/compliance/profile"
	compliancereport "xygo/internal/logic/compliance/report"
	compliancerisk "xygo/internal/logic/compliance/risk"
	compliancetask "xygo/internal/logic/compliance/task"
)

type ComplianceProfileGetReq struct {
	g.Meta `path:"/compliance/profile" method:"get" tags:"会员合规助手" summary:"获取企业画像"`
}

type ComplianceProfileGetRes struct {
	Profile *profile.Profile `json:"profile"`
}

type ComplianceProfileSaveReq struct {
	g.Meta `path:"/compliance/profile" method:"put" tags:"会员合规助手" summary:"保存企业画像"`
	Data   profile.Data `json:"data" v:"required#请填写企业画像"`
}

type ComplianceProfileSaveRes struct {
	Profile *profile.Profile `json:"profile"`
	Changed bool             `json:"changed"`
}

type ComplianceDocumentCreateReq struct {
	g.Meta       `path:"/compliance/documents" method:"post" tags:"会员合规助手" summary:"登记经营资料"`
	AttachmentId uint64 `json:"attachmentId" v:"required|min:1#请上传附件"`
	PeriodKey    string `json:"periodKey" v:"required#请选择资料期间"`
	FileHash     string `json:"fileHash"`
}
type ComplianceDocumentCreateRes struct {
	Document *compliancedocument.BusinessDocument `json:"document"`
}
type ComplianceDocumentListReq struct {
	g.Meta    `path:"/compliance/documents" method:"get" tags:"会员合规助手" summary:"经营资料列表"`
	PeriodKey string `p:"periodKey"`
}
type ComplianceDocumentListRes struct {
	List []compliancedocument.BusinessDocument `json:"list"`
}
type ComplianceDocumentConfirmReq struct {
	g.Meta       `path:"/compliance/documents/{id}/confirm" method:"put" tags:"会员合规助手" summary:"确认经营资料"`
	Id           uint64         `p:"id" in:"path" v:"required|min:1#请指定资料"`
	DocumentType string         `json:"documentType" v:"required#请选择资料类型"`
	Fields       map[string]any `json:"fields"`
}
type ComplianceDocumentConfirmRes struct {
	Document *compliancedocument.BusinessDocument `json:"document"`
}

type ComplianceTaskListReq struct {
	g.Meta    `path:"/compliance/tasks" method:"get" tags:"会员合规助手" summary:"合规任务列表"`
	Status    string `p:"status" json:"status"`
	TaskType  string `p:"taskType" json:"taskType"`
	PeriodKey string `p:"periodKey" json:"periodKey"`
	Page      int    `p:"page" json:"page" d:"1"`
	PageSize  int    `p:"pageSize" json:"pageSize" d:"20"`
}

type ComplianceTaskListRes struct {
	List     []compliancetask.Task `json:"list"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

type ComplianceTaskDetailReq struct {
	g.Meta `path:"/compliance/tasks/{id}" method:"get" tags:"会员合规助手" summary:"合规任务详情"`
	Id     uint64 `p:"id" in:"path" v:"required|min:1#请指定任务"`
}

type ComplianceTaskDetailRes struct {
	Task *compliancetask.Task `json:"task"`
}

type ComplianceTaskActionReq struct {
	g.Meta `path:"/compliance/tasks/{id}/action" method:"post" tags:"会员合规助手" summary:"流转合规任务"`
	Id     uint64 `p:"id" in:"path" v:"required|min:1#请指定任务"`
	Action string `json:"action" v:"required|in:start,submit,reject,complete,cancel#请选择操作|操作无效"`
	Note   string `json:"note"`
}

type ComplianceTaskActionRes struct {
	Task *compliancetask.Task `json:"task"`
}

type ComplianceRiskScanReq struct {
	g.Meta    `path:"/compliance/risks/scan" method:"post" tags:"会员合规助手" summary:"执行合规风险扫描"`
	PeriodKey string               `json:"periodKey" v:"required#请选择扫描期间"`
	Facts     compliancerisk.Facts `json:"facts" v:"required#请提供扫描事实"`
}
type ComplianceRiskScanRes struct {
	List []compliancerisk.Event `json:"list"`
}

type ComplianceRiskListReq struct {
	g.Meta    `path:"/compliance/risks" method:"get" tags:"会员合规助手" summary:"合规风险列表"`
	PeriodKey string `p:"periodKey"`
	Status    string `p:"status"`
}
type ComplianceRiskListRes struct {
	List []compliancerisk.Event `json:"list"`
}

type ComplianceRiskActionReq struct {
	g.Meta `path:"/compliance/risks/{id}/action" method:"post" tags:"会员合规助手" summary:"处置合规风险"`
	Id     uint64 `p:"id" in:"path" v:"required|min:1#请指定风险事件"`
	Action string `json:"action" v:"required|in:confirm,dismiss,resolve#请选择操作|操作无效"`
	Note   string `json:"note"`
}
type ComplianceRiskActionRes struct {
	Risk *compliancerisk.Event `json:"risk"`
}

type ComplianceReportCreateReq struct {
	g.Meta `path:"/compliance/reports" method:"post" tags:"会员合规助手" summary:"生成月度体检报告"`
	Data   compliancereport.Input `json:"data" v:"required#请提供报告数据"`
}
type ComplianceReportCreateRes struct {
	Report *compliancereport.MonthlyReport `json:"report"`
}
type ComplianceReportListReq struct {
	g.Meta    `path:"/compliance/reports" method:"get" tags:"会员合规助手" summary:"月度体检报告列表"`
	PeriodKey string `p:"periodKey"`
}
type ComplianceReportListRes struct {
	List []compliancereport.MonthlyReport `json:"list"`
}
type ComplianceReportPublishReq struct {
	g.Meta `path:"/compliance/reports/{id}/publish" method:"post" tags:"会员合规助手" summary:"发布月度体检报告"`
	Id     uint64 `p:"id" in:"path" v:"required|min:1#请指定报告"`
}
type ComplianceReportPublishRes struct {
	Report *compliancereport.MonthlyReport `json:"report"`
}
