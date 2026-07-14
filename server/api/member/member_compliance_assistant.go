package member

import (
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/profile"
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
