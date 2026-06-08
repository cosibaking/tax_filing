package statement

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/notice"
	"xygo/internal/logic/compliance/shared"
	"xygo/internal/logic/compliance/tax"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
	"xygo/utility"
)

type sComplianceStatement struct{}

func init() {
	service.RegisterComplianceStatement(New())
}

func New() *sComplianceStatement {
	return &sComplianceStatement{}
}

// SendStatement 生成并发送月度对账单
func (s *sComplianceStatement) SendStatement(ctx context.Context, in *compliancein.StatementSendInp) (*compliancein.StatementSendModel, error) {
	if in.OpcId == 0 {
		return nil, gerror.New("请指定OPC主体")
	}
	if in.Year <= 0 || in.Month <= 0 || in.Month > 12 {
		return nil, gerror.New("请指定有效的年月")
	}

	opc, err := shared.LoadOpcById(ctx, in.OpcId)
	if err != nil {
		return nil, err
	}
	if opc.Status != shared.OpcStatusActive {
		return nil, gerror.New("OPC主体尚未激活，无法生成对账单")
	}

	if err := tax.EnsurePeriodTasks(ctx, opc.Id, in.Year, in.Month); err != nil {
		return nil, err
	}

	summary, err := buildAdminSummary(ctx, opc.Id, in.Year, in.Month)
	if err != nil {
		return nil, err
	}
	summaryJSON, _ := gjson.EncodeString(summary)
	now := utility.NowUnix()

	var existing struct {
		Id uint64 `json:"id"`
	}
	_ = g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).
		Where("opc_id", opc.Id).
		Where("year", in.Year).
		Where("month", in.Month).
		Where("deleted", 0).
		Scan(&existing)

	var stmtId uint64
	if existing.Id > 0 {
		_, err = g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).
			Where("id", existing.Id).
			Data(g.Map{
				"summary":     summaryJSON,
				"sent_at":     now,
				"update_time": now,
			}).
			Update()
		stmtId = existing.Id
	} else {
		result, insErr := g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).Data(g.Map{
			"opc_id":      opc.Id,
			"year":        in.Year,
			"month":       in.Month,
			"summary":     summaryJSON,
			"sent_at":     now,
			"deleted":     0,
			"create_time": now,
			"update_time": now,
		}).Insert()
		if insErr != nil {
			return nil, gerror.Wrap(insErr, "生成对账单失败")
		}
		id, _ := result.LastInsertId()
		stmtId = uint64(id)
	}
	if err != nil {
		return nil, err
	}

	_ = audit.WriteAudit(ctx, "monthly_statement", stmtId, "send", in.AdminId, "admin",
		nil, g.Map{"opc_id": opc.Id, "year": in.Year, "month": in.Month}, in.Ip)

	return &compliancein.StatementSendModel{
		Id:     stmtId,
		OpcId:  opc.Id,
		Year:   in.Year,
		Month:  in.Month,
		SentAt: utility.UnixToGTime(now).Format("Y-m-d H:i:s"),
	}, nil
}

// ListStatements 管理端对账单列表
func (s *sComplianceStatement) ListStatements(ctx context.Context, in *compliancein.StatementListInp) (*compliancein.StatementListModel, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	m := g.DB().Model(shared.TableMonthlyStmt+" s").Ctx(ctx).
		LeftJoin(shared.TableOpcEntity+" o", "o.id = s.opc_id AND o.deleted = 0").
		LeftJoin(shared.TableMember+" m", "m.id = o.member_id").
		Where("s.deleted", 0)

	if period := strings.TrimSpace(in.Period); period != "" {
		parts := strings.Split(period, "-")
		if len(parts) == 2 {
			year, _ := strconv.Atoi(parts[0])
			month, _ := strconv.Atoi(parts[1])
			if year > 0 && month > 0 {
				m = m.Where("s.year", year).Where("s.month", month)
			}
		}
	}
	if q := strings.TrimSpace(in.Query); q != "" {
		like := "%" + q + "%"
		m = m.Where("(o.company_name LIKE ? OR m.nickname LIKE ? OR m.mobile LIKE ?)", like, like, like)
	}
	switch in.Status {
	case "draft":
		m = m.Where("s.sent_at", 0)
	case "sent", "completed":
		m = m.Where("s.sent_at > ?", 0)
	}

	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计对账单失败")
	}

	var rows []struct {
		Id          uint64 `json:"id"`
		OpcId       uint64 `json:"opc_id"`
		MemberId    uint64 `json:"member_id"`
		Year        int    `json:"year"`
		Month       int    `json:"month"`
		Summary     string `json:"summary"`
		SentAt      uint64 `json:"sent_at"`
		CompanyName string `json:"company_name"`
		Nickname    string `json:"nickname"`
		Mobile      string `json:"mobile"`
	}
	err = m.Fields(
		"s.id", "s.opc_id", "o.member_id", "s.year", "s.month", "s.summary",
		"s.sent_at", "o.company_name", "m.nickname", "m.mobile",
	).OrderDesc("s.year").OrderDesc("s.month").Page(page, pageSize).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询对账单失败")
	}

	list := make([]compliancein.AdminStatementItem, 0, len(rows))
	for _, row := range rows {
		revenue, cost, profit, prepaidTax := parseSummaryAmounts(row.Summary)
		status := statementStatus(row.SentAt)
		item := compliancein.AdminStatementItem{
			Id:             row.Id,
			MemberId:       row.MemberId,
			MemberName:     row.Nickname,
			MemberPhone:    maskMobile(row.Mobile),
			OpcCompanyName: row.CompanyName,
			Period:         fmt.Sprintf("%04d-%02d", row.Year, row.Month),
			Revenue:        revenue,
			Cost:           cost,
			Profit:         profit,
			PrepaidTax:     prepaidTax,
			Status:         status,
		}
		if row.SentAt > 0 {
			item.SentAt = utility.UnixToGTime(int64(row.SentAt)).Format("Y-m-d H:i:s")
		}
		list = append(list, item)
	}

	return &compliancein.StatementListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// GenerateStatements 按月份批量生成对账单
func (s *sComplianceStatement) GenerateStatements(ctx context.Context, in *compliancein.StatementGenerateInp) (*compliancein.StatementGenerateModel, error) {
	year, month, err := parsePeriod(in.Period)
	if err != nil {
		return nil, err
	}
	if len(in.MemberIds) > 0 {
		created, genErr := s.generateForMembers(ctx, year, month, in.MemberIds)
		if genErr != nil {
			return nil, genErr
		}
		return &compliancein.StatementGenerateModel{Generated: created}, nil
	}
	created, _, err := GenerateForPeriod(ctx, year, month)
	if err != nil {
		return nil, err
	}
	return &compliancein.StatementGenerateModel{Generated: created}, nil
}

// NotifyStatements 批量发送对账单通知
func (s *sComplianceStatement) NotifyStatements(ctx context.Context, in *compliancein.StatementNotifyInp) (*compliancein.StatementNotifyModel, error) {
	if len(in.Ids) == 0 {
		return nil, gerror.New("请选择对账单")
	}

	sent := 0
	now := utility.NowUnix()
	for _, id := range in.Ids {
		var row struct {
			Id       uint64 `json:"id"`
			OpcId    uint64 `json:"opc_id"`
			Year     int    `json:"year"`
			Month    int    `json:"month"`
			Summary  string `json:"summary"`
			MemberId uint64 `json:"member_id"`
		}
		err := g.DB().Model(shared.TableMonthlyStmt+" s").Ctx(ctx).
			LeftJoin(shared.TableOpcEntity+" o", "o.id = s.opc_id AND o.deleted = 0").
			Fields("s.id", "s.opc_id", "s.year", "s.month", "s.summary", "o.member_id").
			Where("s.id", id).
			Where("s.deleted", 0).
			Scan(&row)
		if err != nil || row.Id == 0 {
			continue
		}

		revenue, _, profit, _ := parseSummaryAmounts(row.Summary)
		period := fmt.Sprintf("%04d-%02d", row.Year, row.Month)
		_ = notice.SendToMember(ctx, row.MemberId,
			fmt.Sprintf("%s 月度对账单", period),
			fmt.Sprintf("您的 %s 月度对账单已发送，收入 %.2f 元、利润 %.2f 元，请登录会员中心查看。", period, revenue, profit))

		_, _ = g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).
			Where("id", row.Id).
			Data(g.Map{"sent_at": now, "update_time": now}).
			Update()

		_ = audit.WriteAudit(ctx, "monthly_statement", row.Id, "notify", in.AdminId, "admin",
			nil, g.Map{"period": period}, in.Ip)
		sent++
	}

	return &compliancein.StatementNotifyModel{Sent: sent}, nil
}

func parsePeriod(period string) (year, month int, err error) {
	parts := strings.Split(strings.TrimSpace(period), "-")
	if len(parts) != 2 {
		return 0, 0, gerror.New("月份格式应为 YYYY-MM")
	}
	year, err = strconv.Atoi(parts[0])
	if err != nil || year <= 0 {
		return 0, 0, gerror.New("月份格式应为 YYYY-MM")
	}
	month, err = strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, gerror.New("月份格式应为 YYYY-MM")
	}
	return year, month, nil
}

func (s *sComplianceStatement) generateForMembers(ctx context.Context, year, month int, memberIds []uint64) (int, error) {
	created := 0
	for _, memberId := range memberIds {
		var opc struct {
			Id uint64 `json:"id"`
		}
		_ = g.DB().Model(shared.TableOpcEntity).Ctx(ctx).
			Where("member_id", memberId).
			Where("status", shared.OpcStatusActive).
			Where("deleted", 0).
			Fields("id").
			Scan(&opc)
		if opc.Id == 0 {
			continue
		}
		exists, _ := g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).
			Where("opc_id", opc.Id).
			Where("year", year).
			Where("month", month).
			Where("deleted", 0).
			Count()
		if exists > 0 {
			continue
		}
		summary, sumErr := buildSummary(ctx, opc.Id, year, month, fmt.Sprintf("%04d-%02d", year, month))
		if sumErr != nil {
			return created, sumErr
		}
		summaryJSON, _ := gjson.Encode(summary)
		now := utility.NowUnix()
		_, insErr := g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).Data(g.Map{
			"opc_id":      opc.Id,
			"year":        year,
			"month":       month,
			"summary":     string(summaryJSON),
			"create_time": now,
			"update_time": now,
		}).Insert()
		if insErr != nil {
			return created, gerror.Wrap(insErr, "写入对账单失败")
		}
		created++
	}
	return created, nil
}

func parseSummaryAmounts(raw string) (revenue, cost, profit, prepaidTax float64) {
	if raw == "" {
		return
	}
	var summary map[string]float64
	if err := gjson.DecodeTo(raw, &summary); err != nil {
		return
	}
	revenue = summary["revenue"]
	cost = summary["cost"]
	profit = summary["profit"]
	prepaidTax = summary["taxAmount"]
	if profit == 0 && revenue > 0 {
		profit = revenue - cost
	}
	return
}

func statementStatus(sentAt uint64) string {
	if sentAt == 0 {
		return "draft"
	}
	return "completed"
}

func maskMobile(mobile string) string {
	if len(mobile) < 7 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[len(mobile)-4:]
}

func buildAdminSummary(ctx context.Context, opcId uint64, year, month int) (g.Map, error) {
	var profit struct {
		Revenue float64 `json:"revenue"`
		Cost    float64 `json:"cost"`
		Profit  float64 `json:"profit"`
	}
	_ = g.DB().Model(shared.TableProfitSummary).Ctx(ctx).
		Where("opc_id", opcId).
		Where("year", year).
		Where("month", month).
		Scan(&profit)

	incomeCount, _ := g.DB().Model(shared.TableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		Count()
	expenseCount, _ := g.DB().Model(shared.TableExpenseEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		Count()
	pendingTasks, _ := g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		Where("period", fmt.Sprintf("%04d-%02d", year, month)).
		WhereIn("status", g.Slice{"pending", "overdue"}).
		Count()

	return g.Map{
		"period":         fmt.Sprintf("%04d-%02d", year, month),
		"revenue":        profit.Revenue,
		"cost":           profit.Cost,
		"profit":         profit.Profit,
		"incomeEntries":  incomeCount,
		"expenseEntries": expenseCount,
		"pendingFilings": pendingTasks,
		"generatedAt":    utility.UnixToGTime(utility.NowUnix()).Format("Y-m-d H:i:s"),
	}, nil
}
