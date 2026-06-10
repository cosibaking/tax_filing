package dashboard

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
	"xygo/utility"
)

type sComplianceDashboard struct{}

func init() {
	service.RegisterComplianceDashboard(New())
}

func New() *sComplianceDashboard {
	return &sComplianceDashboard{}
}

// GetOverview 合规服务工作台概览
func (s *sComplianceDashboard) GetOverview(ctx context.Context) (*compliancein.DashboardOverviewModel, error) {
	out := &compliancein.DashboardOverviewModel{}

	activeCustomers, err := g.DB().Model(shared.TableServiceOrder).Ctx(ctx).
		Where("status", shared.OrderStatusActive).
		Where("deleted", 0).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计签约客户失败")
	}
	out.ActiveCustomers = activeCustomers

	pendingOpc, err := g.DB().Model(shared.TableOpcEntity).Ctx(ctx).
		Where("deleted", 0).
		WhereNot("status", shared.OpcStatusActive).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计OPC任务失败")
	}
	out.PendingOpcTasks = pendingOpc

	pendingFilings, err := g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("deleted", 0).
		Where("status", "pending").
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计待申报任务失败")
	}
	out.PendingFilings = pendingFilings

	overdueFilings, err := g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("deleted", 0).
		Where("status", "overdue").
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计逾期申报失败")
	}
	out.OverdueFilings = overdueFilings

	draftStatements, err := g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).
		Where("deleted", 0).
		Where("sent_at", 0).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计草稿对账单失败")
	}
	out.DraftStatements = draftStatements

	pendingSocialConsults, err := g.DB().Model("xy_compliance_social_consult").Ctx(ctx).
		Where("deleted", 0).
		Where("status", "open").
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计待回复社保咨询失败")
	}
	out.PendingSocialConsults = pendingSocialConsults

	out.OpcTasks, err = s.loadPendingOpcTasks(ctx)
	if err != nil {
		return nil, err
	}
	out.FilingTasks, err = s.loadPendingFilings(ctx)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *sComplianceDashboard) loadPendingOpcTasks(ctx context.Context) ([]compliancein.DashboardOpcTaskItem, error) {
	var rows []struct {
		Id                   uint64 `json:"id"`
		ProposedNames        string `json:"proposed_names"`
		Status               string `json:"status"`
		MaterialsSubmittedAt uint64 `json:"materials_submitted_at"`
		Nickname             string `json:"nickname"`
	}
	err := g.DB().Model(shared.TableOpcEntity+" oe").Ctx(ctx).
		LeftJoin(shared.TableMember+" m", "m.id = oe.member_id").
		Where("oe.deleted", 0).
		WhereNot("oe.status", shared.OpcStatusActive).
		Fields("oe.id", "oe.proposed_names", "oe.status", "oe.materials_submitted_at", "m.nickname").
		OrderDesc("oe.materials_submitted_at").
		OrderDesc("oe.update_time").
		Limit(5).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询OPC待办失败")
	}

	list := make([]compliancein.DashboardOpcTaskItem, 0, len(rows))
	for _, row := range rows {
		proposedName := ""
		var names []string
		if row.ProposedNames != "" {
			_ = gjson.DecodeTo(row.ProposedNames, &names)
			if len(names) > 0 {
				proposedName = names[0]
			}
		}
		item := compliancein.DashboardOpcTaskItem{
			OpcId:        row.Id,
			MemberName:   row.Nickname,
			ProposedName: proposedName,
			Status:       row.Status,
			StatusLabel:  opcStatusLabel(row.Status),
		}
		if row.MaterialsSubmittedAt > 0 {
			item.MaterialsSubmittedAt = formatUnix(row.MaterialsSubmittedAt)
			item.DaysSinceSubmit = daysSince(row.MaterialsSubmittedAt)
		}
		list = append(list, item)
	}
	return list, nil
}

func (s *sComplianceDashboard) loadPendingFilings(ctx context.Context) ([]compliancein.DashboardFilingItem, error) {
	var rows []struct {
		Id               uint64  `json:"id"`
		TaxType          string  `json:"tax_type"`
		Period           string  `json:"period"`
		DueDate          uint64  `json:"due_date"`
		Status           string  `json:"status"`
		CalculatedAmount float64 `json:"calculated_amount"`
		CompanyName      string  `json:"company_name"`
		MemberName       string  `json:"member_name"`
	}
	err := g.DB().Model(shared.TableTaxFilingTask+" t").Ctx(ctx).
		LeftJoin(shared.TableOpcEntity+" o", "o.id = t.opc_id AND o.deleted = 0").
		LeftJoin(shared.TableMember+" m", "m.id = o.member_id").
		Where("t.deleted", 0).
		WhereIn("t.status", g.Slice{"pending", "overdue"}).
		Fields(
			"t.id", "t.tax_type", "t.period", "t.due_date", "t.status", "t.calculated_amount",
			"o.company_name", "m.nickname AS member_name",
		).
		OrderAsc("t.due_date").
		Limit(5).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询申报待办失败")
	}

	list := make([]compliancein.DashboardFilingItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, compliancein.DashboardFilingItem{
			Id:               row.Id,
			MemberName:       row.MemberName,
			CompanyName:      row.CompanyName,
			TaxType:          row.TaxType,
			TaxTypeLabel:     taxTypeLabel(row.TaxType),
			Period:           row.Period,
			DueDate:          formatUnix(row.DueDate),
			Status:           row.Status,
			CalculatedAmount: row.CalculatedAmount,
		})
	}
	return list, nil
}

func opcStatusLabel(status string) string {
	labels := map[string]string{
		"pending":          "待提交资料",
		"materials":        "资料待补正",
		"materials_review": "资料审核中",
		"registering":      "工商注册中",
		"tax":              "税务登记中",
		"bank":             "银行开户中",
		"active":           "已激活",
	}
	if v, ok := labels[status]; ok {
		return v
	}
	return status
}

func taxTypeLabel(taxType string) string {
	labels := map[string]string{
		"vat":           "增值税",
		"cit_quarterly": "企税季度预缴",
		"cit_annual":    "企税年度汇算",
		"surcharge":     "附加税费",
	}
	if v, ok := labels[taxType]; ok {
		return v
	}
	return taxType
}

func formatUnix(ts uint64) string {
	if ts == 0 {
		return ""
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m-d")
}

func daysSince(ts uint64) int {
	if ts == 0 {
		return 0
	}
	now := utility.NowUnix()
	if int64(now) <= int64(ts) {
		return 0
	}
	return int((int64(now)-int64(ts))/86400) + 1
}
