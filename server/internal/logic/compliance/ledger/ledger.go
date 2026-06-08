package ledger

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/logic/compliance/shared"
	"xygo/internal/logic/compliance/tax"
	"xygo/internal/model/input/compliancein"
)

type sComplianceLedger struct{}

func New() *sComplianceLedger {
	return &sComplianceLedger{}
}

// GetProfit 月/季/年利润表
func (s *sComplianceLedger) GetProfit(ctx context.Context, in *compliancein.ProfitQueryInp) (*compliancein.ProfitSummaryModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}

	year := in.Year
	if year <= 0 {
		year = gtime.Now().Year()
	}
	period := in.Period
	if period == "" {
		period = "monthly"
	}

	if err := tax.EnsureOpcYearTasks(ctx, opc.Id, year); err != nil {
		return nil, err
	}

	switch period {
	case "yearly":
		return s.yearlyProfit(ctx, opc.Id, year)
	case "quarterly":
		q := in.Quarter
		if q <= 0 {
			q = (int(gtime.Now().Month())-1)/3 + 1
		}
		return s.quarterlyProfit(ctx, opc.Id, year, q)
	default:
		month := in.Month
		if month <= 0 {
			month = int(gtime.Now().Month())
		}
		return s.monthlyProfit(ctx, opc.Id, year, month)
	}
}

func (s *sComplianceLedger) monthlyProfit(ctx context.Context, opcId uint64, year, month int) (*compliancein.ProfitSummaryModel, error) {
	row, err := loadProfitRow(ctx, opcId, year, month)
	if err != nil {
		return nil, err
	}
	item := compliancein.ProfitPeriodItem{
		Label:            fmt.Sprintf("%04d-%02d", year, month),
		Year:             year,
		Month:            month,
		Revenue:          row.Revenue,
		Cost:             row.Cost,
		Profit:           row.Profit,
		CumulativeProfit: row.CumulativeProfit,
	}
	return &compliancein.ProfitSummaryModel{
		Period:  "monthly",
		Year:    year,
		Items:   []compliancein.ProfitPeriodItem{item},
		Revenue: row.Revenue,
		Cost:    row.Cost,
		Profit:  row.Profit,
	}, nil
}

func (s *sComplianceLedger) quarterlyProfit(ctx context.Context, opcId uint64, year, quarter int) (*compliancein.ProfitSummaryModel, error) {
	startMonth := (quarter-1)*3 + 1
	items := make([]compliancein.ProfitPeriodItem, 0, 3)
	var totalRevenue, totalCost, totalProfit float64
	for m := startMonth; m < startMonth+3; m++ {
		row, err := loadProfitRow(ctx, opcId, year, m)
		if err != nil {
			return nil, err
		}
		items = append(items, compliancein.ProfitPeriodItem{
			Label:   fmt.Sprintf("%04d-%02d", year, m),
			Year:    year,
			Month:   m,
			Revenue: row.Revenue,
			Cost:    row.Cost,
			Profit:  row.Profit,
		})
		totalRevenue += row.Revenue
		totalCost += row.Cost
		totalProfit += row.Profit
	}
	return &compliancein.ProfitSummaryModel{
		Period:  "quarterly",
		Year:    year,
		Items:   items,
		Revenue: round2(totalRevenue),
		Cost:    round2(totalCost),
		Profit:  round2(totalProfit),
	}, nil
}

func (s *sComplianceLedger) yearlyProfit(ctx context.Context, opcId uint64, year int) (*compliancein.ProfitSummaryModel, error) {
	var rows []profitRow
	err := g.DB().Model(tableProfitSummary).Ctx(ctx).
		Where("opc_id", opcId).
		Where("year", year).
		OrderAsc("month").
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询年度利润失败")
	}

	items := make([]compliancein.ProfitPeriodItem, 0, len(rows))
	var totalRevenue, totalCost, totalProfit float64
	for _, row := range rows {
		items = append(items, compliancein.ProfitPeriodItem{
			Label:            fmt.Sprintf("%04d-%02d", year, row.Month),
			Year:             year,
			Month:            row.Month,
			Revenue:          row.Revenue,
			Cost:             row.Cost,
			Profit:           row.Profit,
			CumulativeProfit: row.CumulativeProfit,
		})
		totalRevenue += row.Revenue
		totalCost += row.Cost
		totalProfit += row.Profit
	}
	return &compliancein.ProfitSummaryModel{
		Period:  "yearly",
		Year:    year,
		Items:   items,
		Revenue: round2(totalRevenue),
		Cost:    round2(totalCost),
		Profit:  round2(totalProfit),
	}, nil
}

type profitRow struct {
	Month            int     `json:"month"`
	Revenue          float64 `json:"revenue"`
	Cost             float64 `json:"cost"`
	Profit           float64 `json:"profit"`
	CumulativeProfit float64 `json:"cumulative_profit"`
}

func loadProfitRow(ctx context.Context, opcId uint64, year, month int) (profitRow, error) {
	var row profitRow
	err := g.DB().Model(tableProfitSummary).Ctx(ctx).
		Where("opc_id", opcId).
		Where("year", year).
		Where("month", month).
		Scan(&row)
	if err != nil {
		return profitRow{}, gerror.Wrap(err, "查询利润汇总失败")
	}
	if row.Month == 0 {
		return profitRow{Month: month}, nil
	}
	return row, nil
}
