package admin

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"

	api "xygo/api/admin"
	"xygo/internal/library/contexts"
	"xygo/internal/logic/compliance/ruleengine"
	"xygo/internal/service"
)

func (c *ControllerV1) ComplianceRuleVersionSimulate(ctx context.Context, req *api.ComplianceRuleVersionSimulateReq) (*api.ComplianceRuleVersionSimulateRes, error) {
	matched, err := service.ComplianceRuleVersion().Simulate(ctx, req.Id, req.Facts)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceRuleVersionSimulateRes{Matched: matched}, nil
}

func (c *ControllerV1) ComplianceAssistantWorkspace(ctx context.Context, req *api.ComplianceAssistantWorkspaceReq) (*api.ComplianceAssistantWorkspaceRes, error) {
	tables := map[string]string{"rules": "xy_compliance_rule_version", "tasks": "xy_enterprise_compliance_task", "documents": "xy_business_document", "risks": "xy_risk_event", "reports": "xy_compliance_monthly_report", "tickets": "xy_service_ticket"}
	table := tables[req.Section]
	if table == "" {
		return nil, errors.New("工作区无效")
	}
	page, size := req.Page, req.PageSize
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	model := g.DB().Model(table).Ctx(ctx)
	if req.Section != "rules" {
		model = model.Where("deleted", 0)
	}
	if req.Status != "" {
		if req.Section == "documents" {
			model = model.Where("process_status", req.Status)
		} else {
			model = model.Where("status", req.Status)
		}
	}
	total, err := model.Count()
	if err != nil {
		return nil, err
	}
	var list []map[string]any
	if err := model.OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, err
	}
	return &api.ComplianceAssistantWorkspaceRes{List: list, Total: total, Page: page, PageSize: size}, nil
}

func (c *ControllerV1) ComplianceRuleVersionTransition(ctx context.Context, req *api.ComplianceRuleVersionTransitionReq) (*api.ComplianceRuleVersionTransitionRes, error) {
	version, err := service.ComplianceRuleVersion().Transition(ctx, ruleengine.TransitionInput{
		VersionID: req.Id, Action: req.Action, OperatorID: contexts.GetUserId(ctx),
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceRuleVersionTransitionRes{Version: version}, nil
}
