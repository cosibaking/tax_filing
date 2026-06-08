package tax

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
	"xygo/utility"
)

const (
	taskStatusPending = "pending"
	taskStatusOverdue = "overdue"
)

type sComplianceTax struct{}

func init() {
	service.RegisterComplianceTax(New())
}

func New() *sComplianceTax {
	return &sComplianceTax{}
}

// GetCalendar 申报日历
func (s *sComplianceTax) GetCalendar(ctx context.Context, in *compliancein.TaxCalendarInp) (*compliancein.TaxCalendarModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}

	year := in.Year
	if year <= 0 {
		year = gtime.Now().Year()
	}
	if err := EnsureOpcYearTasks(ctx, opc.Id, year); err != nil {
		return nil, err
	}
	_ = refreshOverdueStatus(ctx, opc.Id)

	var rows []taskRow
	err = g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("opc_id", opc.Id).
		Where("deleted", 0).
		WhereLike("period", fmt.Sprintf("%d%%", year)).
		OrderAsc("due_date").
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询申报日历失败")
	}

	items := make([]compliancein.TaxCalendarItem, 0, len(rows))
	dueSet := make(map[string]struct{})
	var nextDue uint64
	for _, row := range rows {
		item := toCalendarItem(row)
		items = append(items, item)
		if row.DueDate > 0 {
			dueStr := formatUnix(row.DueDate)
			dueSet[dueStr] = struct{}{}
			if nextDue == 0 || row.DueDate < nextDue {
				if row.Status == taskStatusPending || row.Status == taskStatusOverdue {
					nextDue = row.DueDate
				}
			}
		}
	}
	dueDates := make([]string, 0, len(dueSet))
	for d := range dueSet {
		dueDates = append(dueDates, d)
	}
	month := int(gtime.Now().Month())
	if in.Month > 0 {
		month = in.Month
	}
	out := &compliancein.TaxCalendarModel{
		Year:  year,
		Month: month,
		Items: items,
		Tasks: items,
		DueDates: dueDates,
	}
	if nextDue > 0 {
		out.NextDueDate = formatUnix(nextDue)
		now := uint64(utility.NowUnix())
		if nextDue > now {
			out.DaysUntilDue = int((nextDue - now) / 86400)
		}
	}
	return out, nil
}

// GetChecklist 报税前自查清单（会员只读）
func (s *sComplianceTax) GetChecklist(ctx context.Context, in *compliancein.TaxChecklistInp) (*compliancein.TaxChecklistModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}

	items := defaultChecklistItems()
	taskId := in.TaskId
	var task taskRow

	if taskId > 0 {
		err = g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
			Where("id", taskId).
			Where("opc_id", opc.Id).
			Where("deleted", 0).
			Scan(&task)
		if err != nil {
			return nil, gerror.Wrap(err, "查询申报任务失败")
		}
		if task.Id == 0 {
			return nil, gerror.New("申报任务不存在")
		}
	} else if in.Period != "" {
		_ = g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
			Where("opc_id", opc.Id).
			Where("deleted", 0).
			Where("period", in.Period).
			WhereIn("status", g.Slice{taskStatusPending, taskStatusOverdue}).
			OrderAsc("due_date").
			Limit(1).
			Scan(&task)
		taskId = task.Id
	} else {
		_ = g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
			Where("opc_id", opc.Id).
			Where("deleted", 0).
			WhereIn("status", g.Slice{taskStatusPending, taskStatusOverdue}).
			OrderAsc("due_date").
			Limit(1).
			Scan(&task)
		taskId = task.Id
	}

	if task.Id > 0 && task.Checklist != "" {
		items = mergeChecklist(items, task.Checklist)
	}

	return &compliancein.TaxChecklistModel{
		TaskId:      taskId,
		Period:      task.Period,
		TaxType:     task.TaxType,
		Items:       items,
		AllComplete: checklistComplete(items),
		ReadOnly:    true,
	}, nil
}

// GetTaskDetail 申报任务详情（含税额计算明细）
func (s *sComplianceTax) GetTaskDetail(ctx context.Context, in *compliancein.TaxTaskDetailInp) (*compliancein.TaxTaskDetailModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	if in.TaskId == 0 {
		return nil, gerror.New("请指定申报任务")
	}

	var task taskRow
	err = g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("id", in.TaskId).
		Where("opc_id", opc.Id).
		Where("deleted", 0).
		Scan(&task)
	if err != nil {
		return nil, gerror.Wrap(err, "查询申报任务失败")
	}
	if task.Id == 0 {
		return nil, gerror.New("申报任务不存在")
	}

	detail := buildTaskCalcDetail(ctx, opc.Id, task)
	return &compliancein.TaxTaskDetailModel{
		Id:               task.Id,
		TaxType:          task.TaxType,
		TaxTypeLabel:     taxTypeLabel(task.TaxType),
		Period:           task.Period,
		DueDate:          formatUnix(task.DueDate),
		Status:           task.Status,
		CalculatedAmount: task.CalculatedAmount,
		FiledAmount:      task.FiledAmount,
		Detail:           detail,
	}, nil
}

func buildTaskCalcDetail(ctx context.Context, opcId uint64, task taskRow) compliancein.TaxTaskCalcDetail {
	detail := compliancein.TaxTaskCalcDetail{Total: task.CalculatedAmount}
	switch task.TaxType {
	case taxTypeVAT:
		revenue, _ := sumMonthlyGrossByPeriod(ctx, opcId, task.Period)
		detail.RevenueExTax = ExTaxAmount(revenue)
		detail.Vat = task.CalculatedAmount
	case taxTypeSurcharge:
		var vatTask struct {
			CalculatedAmount float64 `json:"calculated_amount"`
		}
		_ = g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
			Where("opc_id", opcId).
			Where("tax_type", taxTypeVAT).
			Where("period", task.Period).
			Where("deleted", 0).
			Scan(&vatTask)
		detail.Vat = vatTask.CalculatedAmount
		detail.Surcharge = task.CalculatedAmount
	case taxTypeCITQuarterly, taxTypeCITAnnual:
		detail.Cit = task.CalculatedAmount
	}
	if detail.Total == 0 {
		detail.Total = detail.Vat + detail.Surcharge + detail.Cit
	}
	return detail
}

type taskRow struct {
	Id               uint64  `json:"id"`
	OpcId            uint64  `json:"opc_id"`
	TaxType          string  `json:"tax_type"`
	Period           string  `json:"period"`
	DueDate          uint64  `json:"due_date"`
	Status           string  `json:"status"`
	CalculatedAmount float64 `json:"calculated_amount"`
	FiledAmount      float64 `json:"filed_amount"`
	Checklist        string  `json:"checklist"`
}

func toCalendarItem(row taskRow) compliancein.TaxCalendarItem {
	items := defaultChecklistItems()
	if row.Checklist != "" {
		items = mergeChecklist(items, row.Checklist)
	}
	return compliancein.TaxCalendarItem{
		Id:               row.Id,
		TaxType:          row.TaxType,
		TaxTypeLabel:     taxTypeLabel(row.TaxType),
		Period:           row.Period,
		DueDate:          formatUnix(row.DueDate),
		Status:           row.Status,
		CalculatedAmount: row.CalculatedAmount,
		FiledAmount:      row.FiledAmount,
		AllChecklistDone: checklistComplete(items),
	}
}

// DefaultChecklistItems 报税前自查清单九项
func DefaultChecklistItems() []compliancein.ChecklistItem {
	return defaultChecklistItems()
}

func defaultChecklistItems() []compliancein.ChecklistItem {
	return []compliancein.ChecklistItem{
		{Key: "income_match", Label: "收入流水是否与平台数据一致？"},
		{Key: "invoice_filed", Label: "是否有该申报未申报的发票？"},
		{Key: "cost_booked", Label: "成本费用发票是否已入账？"},
		{Key: "payroll_tax", Label: "OPC是否有员工要报工资个税？"},
		{Key: "vat_filed", Label: "本季度增值税申报了吗？"},
		{Key: "cit_prepaid", Label: "企业所得税预缴了吗？"},
		{Key: "social_insurance", Label: "有员工的话，社保申报了吗？"},
		{Key: "other_platform", Label: "主播是否还有其他平台收入要合并计算？"},
		{Key: "prior_correction", Label: "上一期申报是否有错误要更正？"},
	}
}

func mergeChecklist(defaults []compliancein.ChecklistItem, raw string) []compliancein.ChecklistItem {
	var saved []compliancein.ChecklistItem
	if err := gjson.DecodeTo(raw, &saved); err != nil || len(saved) == 0 {
		return defaults
	}
	byKey := make(map[string]compliancein.ChecklistItem, len(saved))
	for _, item := range saved {
		byKey[item.Key] = item
	}
	out := make([]compliancein.ChecklistItem, 0, len(defaults))
	for _, d := range defaults {
		if saved, ok := byKey[d.Key]; ok {
			out = append(out, saved)
		} else {
			out = append(out, d)
		}
	}
	return out
}

func checklistComplete(items []compliancein.ChecklistItem) bool {
	for _, item := range items {
		if !item.Checked {
			return false
		}
	}
	return len(items) > 0
}

func taxTypeLabel(taxType string) string {
	labels := map[string]string{
		taxTypeVAT:          "增值税",
		taxTypeSurcharge:    "附加税费",
		taxTypeCITQuarterly: "企业所得税（季度预缴）",
		taxTypeCITAnnual:    "企业所得税（年度汇算）",
	}
	if l, ok := labels[taxType]; ok {
		return l
	}
	return taxType
}

func formatUnix(ts uint64) string {
	if ts == 0 {
		return ""
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m-d")
}
