package ledger

import (
	"context"
	"database/sql"
	"errors"
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

	// 申报任务补齐失败不阻断利润查询
	if err := tax.EnsureOpcYearTasks(ctx, opc.Id, year); err != nil {
		g.Log().Warningf(ctx, "[ledger] 补齐申报任务失败 opcId=%d year=%d: %v", opc.Id, year, err)
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
	row, err := calcMonthProfit(ctx, opcId, year, month)
	if err != nil {
		return nil, err
	}
	_ = refreshProfitSummary(ctx, opcId, year, month)
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
		Period:           "monthly",
		Year:             year,
		Items:            []compliancein.ProfitPeriodItem{item},
		Revenue:          row.Revenue,
		Cost:             row.Cost,
		Profit:           row.Profit,
		CumulativeProfit: row.CumulativeProfit,
	}, nil
}

func (s *sComplianceLedger) quarterlyProfit(ctx context.Context, opcId uint64, year, quarter int) (*compliancein.ProfitSummaryModel, error) {
	startMonth := (quarter-1)*3 + 1
	items := make([]compliancein.ProfitPeriodItem, 0, 3)
	var totalRevenue, totalCost, totalProfit float64
	for m := startMonth; m < startMonth+3; m++ {
		row, err := calcMonthProfit(ctx, opcId, year, m)
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
	items := make([]compliancein.ProfitPeriodItem, 0, 12)
	var totalRevenue, totalCost, totalProfit float64
	for m := 1; m <= 12; m++ {
		row, err := calcMonthProfit(ctx, opcId, year, m)
		if err != nil {
			return nil, err
		}
		_ = refreshProfitSummary(ctx, opcId, year, m)
		if row.Revenue == 0 && row.Cost == 0 && row.Profit == 0 {
			continue
		}
		items = append(items, compliancein.ProfitPeriodItem{
			Label:            fmt.Sprintf("%04d-%02d", year, m),
			Year:             year,
			Month:            m,
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

// calcMonthProfit 与收入台账同源：按 monthRange 实时汇总收入/费用
func calcMonthProfit(ctx context.Context, opcId uint64, year, month int) (profitRow, error) {
	period := fmt.Sprintf("%04d-%02d", year, month)
	start, end, err := monthRange(period)
	if err != nil {
		return profitRow{}, err
	}
	revenue, err := sumIncome(ctx, opcId, start, end)
	if err != nil {
		return profitRow{}, err
	}
	cost, err := sumExpense(ctx, opcId, start, end)
	if err != nil {
		return profitRow{}, err
	}
	profit := round2(revenue - cost)
	cumulative, err := calcCumulativeProfit(ctx, opcId, year, month, profit)
	if err != nil {
		return profitRow{}, err
	}
	return profitRow{
		Month:            month,
		Revenue:          revenue,
		Cost:             cost,
		Profit:           profit,
		CumulativeProfit: cumulative,
	}, nil
}

func calcCumulativeProfit(ctx context.Context, opcId uint64, year, month int, currentProfit float64) (float64, error) {
	var prior float64
	for m := 1; m < month; m++ {
		p, err := calcMonthProfitSimple(ctx, opcId, year, m)
		if err != nil {
			return 0, err
		}
		prior += p.Profit
	}
	return round2(prior + currentProfit), nil
}

func calcMonthProfitSimple(ctx context.Context, opcId uint64, year, month int) (profitRow, error) {
	period := fmt.Sprintf("%04d-%02d", year, month)
	start, end, err := monthRange(period)
	if err != nil {
		return profitRow{}, err
	}
	revenue, err := sumIncome(ctx, opcId, start, end)
	if err != nil {
		return profitRow{}, err
	}
	cost, err := sumExpense(ctx, opcId, start, end)
	if err != nil {
		return profitRow{}, err
	}
	return profitRow{
		Month:  month,
		Revenue: revenue,
		Cost:    cost,
		Profit:  round2(revenue - cost),
	}, nil
}

func loadProfitRow(ctx context.Context, opcId uint64, year, month int) (profitRow, error) {
	row, err := calcMonthProfit(ctx, opcId, year, month)
	if err != nil {
		return profitRow{}, err
	}
	_ = refreshProfitSummary(ctx, opcId, year, month)

	var cached profitRow
	err = g.DB().Model(tableProfitSummary).Ctx(ctx).
		Where("opc_id", opcId).
		Where("year", year).
		Where("month", month).
		Scan(&cached)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return row, nil
		}
		return profitRow{}, gerror.Wrap(err, "查询利润汇总失败")
	}
	if cached.Month == 0 {
		return row, nil
	}
	return row, nil
}
