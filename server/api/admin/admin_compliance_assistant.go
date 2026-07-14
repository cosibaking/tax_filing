package admin

import (
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/ruleengine"
)

type ComplianceRuleVersionSimulateReq struct {
	g.Meta `path:"/admin/compliance/rule-versions/{id}/simulate" method:"post" tags:"合规助手" summary:"试算规则版本"`
	Id     uint64         `p:"id" in:"path" v:"required|min:1#请指定规则版本"`
	Facts  map[string]any `json:"facts" v:"required#请提供试算事实"`
}

type ComplianceRuleVersionSimulateRes struct {
	Matched bool `json:"matched"`
}

type ComplianceRuleVersionTransitionReq struct {
	g.Meta `path:"/admin/compliance/rule-versions/{id}/transition" method:"post" tags:"合规助手" summary:"流转规则版本"`
	Id     uint64 `p:"id" in:"path" v:"required|min:1#请指定规则版本"`
	Action string `json:"action" v:"required|in:submit_review,publish,retire#请选择操作|操作无效"`
}

type ComplianceRuleVersionTransitionRes struct {
	Version *ruleengine.RuleVersion `json:"version"`
}

type ComplianceAssistantWorkspaceReq struct {
	g.Meta   `path:"/admin/compliance/assistant/workspace" method:"get" tags:"合规助手" summary:"合规助手管理工作台"`
	Section  string `p:"section" v:"required|in:rules,tasks,documents,risks,reports,tickets#请选择工作区|工作区无效"`
	Status   string `p:"status"`
	Page     int    `p:"page" d:"1"`
	PageSize int    `p:"pageSize" d:"20"`
}
type ComplianceAssistantWorkspaceRes struct {
	List     []map[string]any `json:"list"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}
