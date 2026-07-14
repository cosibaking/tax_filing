package member

import (
	"context"

	api "xygo/api/member"
	compliancetask "xygo/internal/logic/compliance/task"
	"xygo/internal/service"
)

func (c *ControllerV1) ComplianceTaskList(ctx context.Context, req *api.ComplianceTaskListReq) (*api.ComplianceTaskListRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	items, total, err := service.ComplianceTask().List(ctx, compliancetask.Query{
		MemberID: memberID, Status: req.Status, TaskType: req.TaskType, PeriodKey: req.PeriodKey,
		Page: req.Page, PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceTaskListRes{List: items, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (c *ControllerV1) ComplianceTaskDetail(ctx context.Context, req *api.ComplianceTaskDetailReq) (*api.ComplianceTaskDetailRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	item, err := service.ComplianceTask().GetForMember(ctx, memberID, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceTaskDetailRes{Task: item}, nil
}

func (c *ControllerV1) ComplianceTaskAction(ctx context.Context, req *api.ComplianceTaskActionReq) (*api.ComplianceTaskActionRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	item, err := service.ComplianceTask().ApplyEvent(ctx, memberID, req.Id, req.Action, req.Note)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceTaskActionRes{Task: item}, nil
}
