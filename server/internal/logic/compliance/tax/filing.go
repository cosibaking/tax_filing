package tax

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/rules"
	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
	"xygo/utility"
)

const taskStatusFiled = "filed"

func init() {
	service.RegisterComplianceFiling(New())
}

// ListTasks 管理端申报任务列表
func (s *sComplianceTax) ListTasks(ctx context.Context, in *compliancein.FilingListInp) (*compliancein.FilingListModel, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	m := g.DB().Model(shared.TableTaxFilingTask+" t").Ctx(ctx).
		LeftJoin(shared.TableOpcEntity+" o", "o.id = t.opc_id AND o.deleted = 0").
		LeftJoin(shared.TableMember+" m", "m.id = o.member_id").
		Where("t.deleted", 0)

	if in.Status != "" {
		m = m.Where("t.status", in.Status)
	}
	if in.TaxType != "" {
		m = m.Where("t.tax_type", in.TaxType)
	}
	if q := strings.TrimSpace(in.Query); q != "" {
		like := "%" + q + "%"
		m = m.Where("(o.company_name LIKE ? OR m.nickname LIKE ? OR m.mobile LIKE ?)", like, like, like)
	}
	if period := strings.TrimSpace(in.Period); period != "" {
		m = m.Where("t.period", period)
	}

	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计申报任务失败")
	}

	var rows []filingListRow
	err = m.Fields(
		"t.id", "t.opc_id", "t.tax_type", "t.period", "t.due_date", "t.status",
		"t.calculated_amount", "t.filed_amount", "t.checklist",
		"o.company_name", "o.member_id", "m.nickname AS member_name",
	).OrderDesc("t.due_date").Page(page, pageSize).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询申报任务失败")
	}

	list := make([]compliancein.FilingTaskItem, 0, len(rows))
	for _, row := range rows {
		items := defaultChecklistItems()
		if row.Checklist != "" {
			items = mergeChecklist(items, row.Checklist)
		}
		list = append(list, compliancein.FilingTaskItem{
			Id:               row.Id,
			OpcId:            row.OpcId,
			MemberId:         row.MemberId,
			CompanyName:      row.CompanyName,
			MemberName:       row.MemberName,
			TaxType:          row.TaxType,
			TaxTypeLabel:     taxTypeLabel(row.TaxType),
			Period:           row.Period,
			DueDate:          formatUnix(row.DueDate),
			Status:           row.Status,
			CalculatedAmount: row.CalculatedAmount,
			FiledAmount:      row.FiledAmount,
			AllChecklistDone: checklistComplete(items),
		})
	}

	return &compliancein.FilingListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// MarkFiled 顾问标记已申报（自查清单 + 合规规则引擎 F-68）
func (s *sComplianceTax) MarkFiled(ctx context.Context, in *compliancein.FilingMarkInp) (*compliancein.FilingMarkModel, error) {
	if in.TaskId == 0 {
		return nil, gerror.New("请指定申报任务")
	}

	var task taskRow
	err := g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("id", in.TaskId).
		Where("deleted", 0).
		Scan(&task)
	if err != nil {
		return nil, gerror.Wrap(err, "查询申报任务失败")
	}
	if task.Id == 0 {
		return nil, gerror.New("申报任务不存在")
	}
	if task.Status == taskStatusFiled {
		return nil, gerror.New("该任务已标记申报完成")
	}

	items := in.Checklist
	if len(items) == 0 {
		items = defaultChecklistItems()
		if task.Checklist != "" {
			items = mergeChecklist(items, task.Checklist)
		}
	}
	if !checklistComplete(items) {
		return nil, gerror.New("自查清单未全部勾选，无法标记已申报")
	}
	if in.ReceiptFileId == 0 {
		return nil, gerror.New("请上传申报回执 PDF")
	}

	filedAmount := in.FiledAmount
	if filedAmount <= 0 {
		filedAmount = task.CalculatedAmount
	}
	if err = rules.ValidateFiling(ctx, task.OpcId, task.Period, filedAmount, task.CalculatedAmount, in.ReportedIncome); err != nil {
		return nil, err
	}

	checklistJSON, _ := json.Marshal(items)
	now := utility.NowUnix()

	_, err = g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("id", task.Id).
		Data(g.Map{
			"status":          taskStatusFiled,
			"filed_amount":    filedAmount,
			"receipt_file_id": in.ReceiptFileId,
			"checklist":       string(checklistJSON),
			"update_time":     now,
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新申报任务失败")
	}

	_ = audit.WriteAudit(ctx, "tax_filing_task", task.Id, "filing.mark_filed", in.AdminId, "admin",
		g.Map{"status": task.Status}, g.Map{"status": taskStatusFiled, "filed_amount": filedAmount}, in.Ip)

	return &compliancein.FilingMarkModel{Id: task.Id, Status: taskStatusFiled}, nil
}

type filingListRow struct {
	Id               uint64  `json:"id"`
	OpcId            uint64  `json:"opc_id"`
	TaxType          string  `json:"tax_type"`
	Period           string  `json:"period"`
	DueDate          uint64  `json:"due_date"`
	Status           string  `json:"status"`
	CalculatedAmount float64 `json:"calculated_amount"`
	FiledAmount      float64 `json:"filed_amount"`
	Checklist        string  `json:"checklist"`
	MemberId         uint64  `json:"member_id"`
	CompanyName      string  `json:"company_name"`
	MemberName       string  `json:"member_name"`
}
