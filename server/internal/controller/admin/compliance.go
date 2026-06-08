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
	"context"
	"fmt"

	"github.com/gogf/gf/v2/net/ghttp"



	api "xygo/api/admin"

	"xygo/internal/library/contexts"

	"xygo/internal/model/input/compliancein"

	"xygo/internal/service"

)



// ComplianceOpcTaskList OPC 任务列表

func (c *ControllerV1) ComplianceOpcTaskList(ctx context.Context, req *api.ComplianceOpcTaskListReq) (res *api.ComplianceOpcTaskListRes, err error) {

	out, err := service.ComplianceOpc().ListTasks(ctx, &compliancein.OpcTaskListInp{

		Page:     req.Page,

		PageSize: req.PageSize,

		Status:   req.Status,

		Query:    req.Query,

	})

	if err != nil {

		return nil, err

	}

	return &api.ComplianceOpcTaskListRes{OpcTaskListModel: out}, nil

}



// ComplianceOpcTaskDetail OPC 任务详情

func (c *ControllerV1) ComplianceOpcTaskDetail(ctx context.Context, req *api.ComplianceOpcTaskDetailReq) (res *api.ComplianceOpcTaskDetailRes, err error) {

	out, err := service.ComplianceOpc().GetTaskDetail(ctx, &compliancein.OpcTaskDetailInp{

		OpcId: req.OpcId,

	})

	if err != nil {

		return nil, err

	}

	return &api.ComplianceOpcTaskDetailRes{OpcTaskDetailModel: out}, nil

}



// ComplianceOpcTaskAction 推进 OPC 任务

func (c *ControllerV1) ComplianceOpcTaskAction(ctx context.Context, req *api.ComplianceOpcTaskActionReq) (res *api.ComplianceOpcTaskActionRes, err error) {

	ip := ""

	if r := ghttp.RequestFromCtx(ctx); r != nil {

		ip = r.GetClientIp()

	}



	out, err := service.ComplianceOpc().AdvanceTask(ctx, &compliancein.OpcTaskActionInp{

		OpcId:             req.OpcId,

		Action:            req.Action,

		Note:              req.Note,

		AdminId:           contexts.GetUserId(ctx),

		Ip:                ip,

		CompanyName:       req.CompanyName,

		CreditCode:        req.CreditCode,

		EstablishedAt:     req.EstablishedAt,

		LicenseFileId:     req.LicenseFileId,

		TaxActivatedAt:    req.TaxActivatedAt,

		BankAccount:       req.BankAccount,

		BankReceiptFileId: req.BankReceiptFileId,

	})

	if err != nil {

		return nil, err

	}

	return &api.ComplianceOpcTaskActionRes{OpcTaskActionModel: out}, nil

}



// ComplianceCustomerList 合规客户列表

func (c *ControllerV1) ComplianceCustomerList(ctx context.Context, req *api.ComplianceCustomerListReq) (res *api.ComplianceCustomerListRes, err error) {

	out, err := service.ComplianceOpc().ListCustomers(ctx, &compliancein.ComplianceCustomerListInp{

		Page:      req.Page,

		PageSize:  req.PageSize,

		Query:     req.Query,
		OpcStatus: req.OpcStatus,

	})

	if err != nil {

		return nil, err

	}

	return &api.ComplianceCustomerListRes{ComplianceCustomerListModel: out}, nil

}

// ComplianceFilingList 申报任务列表
func (c *ControllerV1) ComplianceFilingList(ctx context.Context, req *api.ComplianceFilingListReq) (res *api.ComplianceFilingListRes, err error) {
	out, err := service.ComplianceFiling().ListTasks(ctx, &compliancein.FilingListInp{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		TaxType:  req.TaxType,
		Period:   req.Period,
		Query:    req.Query,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceFilingListRes{FilingListModel: out}, nil
}

// ComplianceFilingMark 标记已申报
func (c *ControllerV1) ComplianceFilingMark(ctx context.Context, req *api.ComplianceFilingMarkReq) (res *api.ComplianceFilingMarkRes, err error) {
	ip := ""
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
	}
	out, err := service.ComplianceFiling().MarkFiled(ctx, &compliancein.FilingMarkInp{
		TaskId:         req.TaskId,
		FiledAmount:    req.FiledAmount,
		ReportedIncome: req.ReportedIncome,
		ReceiptFileId:  req.ReceiptFileId,
		Checklist:      req.Checklist,
		AdminId:        contexts.GetUserId(ctx),
		Ip:             ip,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceFilingMarkRes{FilingMarkModel: out}, nil
}

// ComplianceStatementList 对账单管理列表
func (c *ControllerV1) ComplianceStatementList(ctx context.Context, req *api.ComplianceStatementListReq) (res *api.ComplianceStatementListRes, err error) {
	out, err := service.ComplianceStatement().ListStatements(ctx, &compliancein.StatementListInp{
		Page:     req.Page,
		PageSize: req.PageSize,
		Period:   req.Period,
		Status:   req.Status,
		Query:    req.Query,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceStatementListRes{StatementListModel: out}, nil
}

// ComplianceStatementGenerate 批量生成对账单
func (c *ControllerV1) ComplianceStatementGenerate(ctx context.Context, req *api.ComplianceStatementGenerateReq) (res *api.ComplianceStatementGenerateRes, err error) {
	ip := ""
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
	}
	out, err := service.ComplianceStatement().GenerateStatements(ctx, &compliancein.StatementGenerateInp{
		Period:    req.Period,
		MemberIds: req.MemberIds,
		AdminId:   contexts.GetUserId(ctx),
		Ip:        ip,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceStatementGenerateRes{StatementGenerateModel: out}, nil
}

// ComplianceStatementNotify 批量发送对账单通知
func (c *ControllerV1) ComplianceStatementNotify(ctx context.Context, req *api.ComplianceStatementNotifyReq) (res *api.ComplianceStatementNotifyRes, err error) {
	ip := ""
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
	}
	out, err := service.ComplianceStatement().NotifyStatements(ctx, &compliancein.StatementNotifyInp{
		Ids:     req.Ids,
		AdminId: contexts.GetUserId(ctx),
		Ip:      ip,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceStatementNotifyRes{StatementNotifyModel: out}, nil
}

// ComplianceStatementSend 单主体生成并发送月度对账单
func (c *ControllerV1) ComplianceStatementSend(ctx context.Context, req *api.ComplianceStatementSendReq) (res *api.ComplianceStatementSendRes, err error) {
	ip := ""
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
	}
	out, err := service.ComplianceStatement().SendStatement(ctx, &compliancein.StatementSendInp{
		OpcId:   req.OpcId,
		Year:    req.Year,
		Month:   req.Month,
		AdminId: contexts.GetUserId(ctx),
		Ip:      ip,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceStatementSendRes{StatementSendModel: out}, nil
}

// ComplianceAuditExport 导出会员合规留痕 CSV
func (c *ControllerV1) ComplianceAuditExport(ctx context.Context, req *api.ComplianceAuditExportReq) (res *api.ComplianceAuditExportRes, err error) {
	data, filename, err := service.ComplianceAudit().ExportCSV(ctx, req.MemberId)
	if err != nil {
		return nil, err
	}
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
		r.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		r.Response.Write(data)
	}
	return &api.ComplianceAuditExportRes{}, nil
}

// ComplianceDashboard 合规工作台概览
func (c *ControllerV1) ComplianceDashboard(ctx context.Context, req *api.ComplianceDashboardReq) (res *api.ComplianceDashboardRes, err error) {
	out, err := service.ComplianceDashboard().GetOverview(ctx)
	if err != nil {
		return nil, err
	}
	return &api.ComplianceDashboardRes{DashboardOverviewModel: out}, nil
}

