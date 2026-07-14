package admin

import (
	"context"

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

func (c *ControllerV1) ComplianceRuleVersionTransition(ctx context.Context, req *api.ComplianceRuleVersionTransitionReq) (*api.ComplianceRuleVersionTransitionRes, error) {
	version, err := service.ComplianceRuleVersion().Transition(ctx, ruleengine.TransitionInput{
		VersionID: req.Id, Action: req.Action, OperatorID: contexts.GetUserId(ctx),
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceRuleVersionTransitionRes{Version: version}, nil
}
