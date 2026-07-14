package member

import (
	"context"
	"fmt"

	api "xygo/api/member"
	compliancedocument "xygo/internal/logic/compliance/document"
	"xygo/internal/logic/compliance/profile"
	"xygo/internal/logic/compliance/shared"
	compliancetask "xygo/internal/logic/compliance/task"
	"xygo/internal/service"
)

func (c *ControllerV1) ComplianceDocumentCreate(ctx context.Context, req *api.ComplianceDocumentCreateReq) (*api.ComplianceDocumentCreateRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	opc, err := shared.LoadOpcByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if opc == nil {
		return nil, fmt.Errorf("请先完成企业主体建档")
	}
	item, err := service.ComplianceDocument().Register(ctx, compliancedocument.RegisterInput{OpcEntityID: opc.Id, MemberID: memberID, AttachmentID: req.AttachmentId, PeriodKey: req.PeriodKey, FileHash: req.FileHash})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceDocumentCreateRes{Document: item}, nil
}

func (c *ControllerV1) ComplianceDocumentList(ctx context.Context, req *api.ComplianceDocumentListReq) (*api.ComplianceDocumentListRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	items, err := service.ComplianceDocument().List(ctx, memberID, req.PeriodKey)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceDocumentListRes{List: items}, nil
}

func (c *ControllerV1) ComplianceDocumentConfirm(ctx context.Context, req *api.ComplianceDocumentConfirmReq) (*api.ComplianceDocumentConfirmRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	item, err := service.ComplianceDocument().Confirm(ctx, memberID, req.Id, req.DocumentType, req.Fields)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceDocumentConfirmRes{Document: item}, nil
}

func (c *ControllerV1) ComplianceProfileGet(ctx context.Context, _ *api.ComplianceProfileGetReq) (*api.ComplianceProfileGetRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	opc, err := shared.LoadOpcByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if opc == nil {
		return &api.ComplianceProfileGetRes{}, nil
	}
	item, err := service.ComplianceProfile().GetForMember(ctx, memberID, opc.Id)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceProfileGetRes{Profile: item}, nil
}

func (c *ControllerV1) ComplianceProfileSave(ctx context.Context, req *api.ComplianceProfileSaveReq) (*api.ComplianceProfileSaveRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	opc, err := shared.LoadOpcByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if opc == nil {
		return nil, fmt.Errorf("请先完成企业主体建档")
	}
	item, changed, err := service.ComplianceProfile().Save(ctx, profile.SaveInput{
		MemberID: memberID, OpcID: opc.Id, Source: profile.SourceUser, Data: req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceProfileSaveRes{Profile: item, Changed: changed}, nil
}

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

func (c *ControllerV1) ComplianceRiskScan(ctx context.Context, req *api.ComplianceRiskScanReq) (*api.ComplianceRiskScanRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	opc, err := shared.LoadOpcByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if opc == nil {
		return nil, fmt.Errorf("请先完成企业主体建档")
	}
	items, err := service.ComplianceRisk().ScanAndSave(ctx, opc.Id, memberID, req.PeriodKey, req.Facts)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceRiskScanRes{List: items}, nil
}

func (c *ControllerV1) ComplianceRiskList(ctx context.Context, req *api.ComplianceRiskListReq) (*api.ComplianceRiskListRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	items, err := service.ComplianceRisk().List(ctx, memberID, req.PeriodKey, req.Status)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceRiskListRes{List: items}, nil
}

func (c *ControllerV1) ComplianceRiskAction(ctx context.Context, req *api.ComplianceRiskActionReq) (*api.ComplianceRiskActionRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	item, err := service.ComplianceRisk().Decide(ctx, memberID, req.Id, req.Action, req.Note)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceRiskActionRes{Risk: item}, nil
}

func (c *ControllerV1) ComplianceReportCreate(ctx context.Context, req *api.ComplianceReportCreateReq) (*api.ComplianceReportCreateRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	opc, err := shared.LoadOpcByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if opc == nil {
		return nil, fmt.Errorf("请先完成企业主体建档")
	}
	item, err := service.ComplianceReport().Create(ctx, opc.Id, memberID, req.Data)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceReportCreateRes{Report: item}, nil
}
func (c *ControllerV1) ComplianceReportList(ctx context.Context, req *api.ComplianceReportListReq) (*api.ComplianceReportListRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	items, err := service.ComplianceReport().List(ctx, memberID, req.PeriodKey)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceReportListRes{List: items}, nil
}
func (c *ControllerV1) ComplianceReportPublish(ctx context.Context, req *api.ComplianceReportPublishReq) (*api.ComplianceReportPublishRes, error) {
	memberID, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	item, err := service.ComplianceReport().Publish(ctx, memberID, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceReportPublishRes{Report: item}, nil
}
