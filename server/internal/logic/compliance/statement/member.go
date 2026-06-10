package statement

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
)

// ListMemberStatements 会员端对账单列表
func (s *sComplianceStatement) ListMemberStatements(ctx context.Context, in *compliancein.MemberStatementListInp) (*compliancein.MemberStatementListModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	m := g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).
		Where("opc_id", opc.Id).
		Where("deleted", 0)

	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计对账单失败")
	}

	var rows []struct {
		Id      uint64 `json:"id"`
		Year    int    `json:"year"`
		Month   int    `json:"month"`
		Summary string `json:"summary"`
		SentAt  uint64 `json:"sent_at"`
	}
	err = m.OrderDesc("year").OrderDesc("month").Page(page, pageSize).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询对账单失败")
	}

	list := make([]compliancein.MemberStatementItem, 0, len(rows))
	for _, row := range rows {
		revenue, cost, profit, prepaidTax, cumulative, filingStatus := parseMemberSummary(row.Summary)
		status := "draft"
		if row.SentAt > 0 {
			status = "completed"
		}
		list = append(list, compliancein.MemberStatementItem{
			Id:                row.Id,
			Period:            fmt.Sprintf("%04d-%02d", row.Year, row.Month),
			Revenue:           revenue,
			Cost:              cost,
			Profit:            profit,
			PrepaidTax:        prepaidTax,
			CumulativeProfit:  cumulative,
			FilingStatus:      filingStatus,
			Status:            status,
			PdfUrl:            statementPdfUrlIfExists(row.Id),
		})
	}

	return &compliancein.MemberStatementListModel{
		List:  list,
		Total: total,
	}, nil
}

func parseMemberSummary(raw string) (revenue, cost, profit, prepaidTax, cumulative float64, filingStatus string) {
	if raw == "" {
		return
	}
	var summary struct {
		Revenue          float64 `json:"revenue"`
		Cost             float64 `json:"cost"`
		Profit           float64 `json:"profit"`
		TaxAmount        float64 `json:"taxAmount"`
		CumulativeProfit float64 `json:"cumulativeProfit"`
		FilingStatus     string  `json:"filingStatus"`
	}
	if err := gjson.DecodeTo(raw, &summary); err != nil {
		revenue, cost, profit, prepaidTax = parseSummaryAmounts(raw)
		return
	}
	revenue = summary.Revenue
	cost = summary.Cost
	profit = summary.Profit
	prepaidTax = summary.TaxAmount
	cumulative = summary.CumulativeProfit
	filingStatus = summary.FilingStatus
	if profit == 0 && revenue > 0 {
		profit = revenue - cost
	}
	return
}
