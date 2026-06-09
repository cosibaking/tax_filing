package tax

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/rules"
	"xygo/utility"
)

const (
	tableTaxFilingTask = "xy_tax_filing_task"
	tableOpcEntity     = "xy_opc_entity"
	tableProfitSummary = "xy_profit_summary"

	taxTypeVAT           = "vat"
	taxTypeSurcharge     = "surcharge"
	taxTypeCITQuarterly  = "cit_quarterly"
	taxTypeCITAnnual     = "cit_annual"
)

// GenerateTasks 为所有 active OPC 生成当期税务申报任务
func GenerateTasks(ctx context.Context, now time.Time) (created, skipped int, err error) {
	var opcs []struct {
		Id uint64 `json:"id"`
	}
	err = g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("status", "active").
		Where("deleted", 0).
		Fields("id").
		Scan(&opcs)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "查询OPC主体失败")
	}

	prev := now.AddDate(0, -1, 0)
	monthPeriod := fmt.Sprintf("%04d-%02d", prev.Year(), int(prev.Month()))
	monthRevenuePeriod := monthPeriod

	for _, opc := range opcs {
		revenue, _ := sumMonthlyGrossByPeriod(ctx, opc.Id, monthRevenuePeriod)
		vat := CalcMonthlyVAT(revenue)
		surcharge := CalcSurcharge(vat)

		c, s, e := upsertTask(ctx, opc.Id, taxTypeVAT, monthPeriod, vat, dueDateNextMonth15(now))
		created += c
		skipped += s
		if e != nil {
			return created, skipped, e
		}
		c, s, e = upsertTask(ctx, opc.Id, taxTypeSurcharge, monthPeriod, surcharge, dueDateNextMonth15(now))
		created += c
		skipped += s
		if e != nil {
			return created, skipped, e
		}

		if isQuarterStart(now.Month()) {
			q, y := previousQuarter(now)
			qPeriod := fmt.Sprintf("Q%d-%d", q, y)
			profit := sumQuarterProfit(ctx, opc.Id, y, q)
			cit := CalcQuarterlyCIT(profit)
			c, s, e = upsertTask(ctx, opc.Id, taxTypeCITQuarterly, qPeriod, cit, dueDateMonthEnd(now))
			created += c
			skipped += s
			if e != nil {
				return created, skipped, e
			}
		}

		if now.Month() == time.January {
			annualPeriod := fmt.Sprintf("%d", now.Year()-1)
			annualProfit := sumAnnualProfit(ctx, opc.Id, now.Year()-1)
			annualCit := CalcQuarterlyCIT(annualProfit)
			c, s, e = upsertTask(ctx, opc.Id, taxTypeCITAnnual, annualPeriod, annualCit, dueDateMay31(now.Year()))
			created += c
			skipped += s
			if e != nil {
				return created, skipped, e
			}
		}
	}
	return created, skipped, nil
}

func upsertTask(ctx context.Context, opcId uint64, taxType, period string, amount float64, dueDate time.Time) (created, skipped int, err error) {
	var existing struct {
		Id     uint64 `json:"id"`
		Status string `json:"status"`
	}
	err = g.DB().Model(tableTaxFilingTask).Ctx(ctx).
		Where("opc_id", opcId).
		Where("tax_type", taxType).
		Where("period", period).
		Where("deleted", 0).
		Scan(&existing)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}

	now := uint64(utility.NowUnix())
	if existing.Id > 0 {
		if existing.Status == "filed" {
			return 0, 1, nil
		}
		_, err = g.DB().Model(tableTaxFilingTask).Ctx(ctx).
			Where("id", existing.Id).
			Data(g.Map{
				"calculated_amount": amount,
				"update_time":       now,
			}).Update()
		if err != nil {
			return 0, 0, gerror.Wrap(err, "更新申报任务税额失败")
		}
		return 0, 1, nil
	}

	_, err = g.DB().Model(tableTaxFilingTask).Ctx(ctx).Data(g.Map{
		"opc_id":            opcId,
		"tax_type":          taxType,
		"period":            period,
		"due_date":          uint64(dueDate.Unix()),
		"status":            "pending",
		"calculated_amount": amount,
		"create_time":       now,
		"update_time":       now,
	}).Insert()
	if err != nil {
		return 0, 0, gerror.Wrap(err, "创建申报任务失败")
	}
	return 1, 0, nil
}

// RefreshOpcPeriodTaxAmounts 台账变更后重算指定月份申报税额
func RefreshOpcPeriodTaxAmounts(ctx context.Context, opcId uint64, year, month int) error {
	if opcId == 0 || year <= 0 || month <= 0 {
		return nil
	}
	return EnsurePeriodTasks(ctx, opcId, year, month)
}

func dueDateNextMonth15(now time.Time) time.Time {
	y, m := now.Year(), int(now.Month())+1
	if m > 12 {
		m = 1
		y++
	}
	return time.Date(y, time.Month(m), 15, 23, 59, 59, 0, time.Local)
}

func dueDateMonthEnd(now time.Time) time.Time {
	y, m := now.Year(), now.Month()
	start := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	return start.AddDate(0, 1, 0).Add(-time.Second)
}

func dueDateMay31(year int) time.Time {
	return time.Date(year, 5, 31, 23, 59, 59, 0, time.Local)
}

func isQuarterStart(m time.Month) bool {
	return m == time.January || m == time.April || m == time.July || m == time.October
}

func previousQuarter(now time.Time) (q, year int) {
	m := int(now.Month())
	switch {
	case m <= 3:
		return 4, now.Year() - 1
	case m <= 6:
		return 1, now.Year()
	case m <= 9:
		return 2, now.Year()
	default:
		return 3, now.Year()
	}
}

func sumQuarterProfit(ctx context.Context, opcId uint64, year, quarter int) float64 {
	startMonth := (quarter-1)*3 + 1
	var total float64
	for m := startMonth; m < startMonth+3; m++ {
		var row struct {
			Profit float64 `json:"profit"`
		}
		_ = g.DB().Model(tableProfitSummary).Ctx(ctx).
			Where("opc_id", opcId).
			Where("year", year).
			Where("month", m).
			Scan(&row)
		total += row.Profit
	}
	return total
}

// EnsurePeriodTasks 为单个 OPC 补齐指定月份申报任务
func EnsurePeriodTasks(ctx context.Context, opcId uint64, year, month int) error {
	period := fmt.Sprintf("%04d-%02d", year, month)
	revenue, err := sumMonthlyGrossByPeriod(ctx, opcId, period)
	if err != nil {
		return err
	}
	vat := CalcMonthlyVAT(revenue)
	surcharge := CalcSurcharge(vat)
	due := dueDateNextMonth15(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0))
	if _, _, err = upsertTask(ctx, opcId, taxTypeVAT, period, vat, due); err != nil {
		return err
	}
	if _, _, err = upsertTask(ctx, opcId, taxTypeSurcharge, period, surcharge, due); err != nil {
		return err
	}
	if month%3 == 0 {
		q := month / 3
		qPeriod := fmt.Sprintf("Q%d-%d", q, year)
		profit := sumQuarterProfit(ctx, opcId, year, q)
		cit := CalcQuarterlyCIT(profit)
		qDue := dueDateMonthEnd(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0))
		if _, _, err = upsertTask(ctx, opcId, taxTypeCITQuarterly, qPeriod, cit, qDue); err != nil {
			return err
		}
	}
	return refreshOverdueStatus(ctx, opcId)
}

// EnsureOpcYearTasks 为单个 OPC 补齐指定年份申报任务（日历/利润预览用）
func EnsureOpcYearTasks(ctx context.Context, opcId uint64, year int) error {
	for m := 1; m <= 12; m++ {
		period := fmt.Sprintf("%04d-%02d", year, m)
		revenue, err := sumMonthlyGrossByPeriod(ctx, opcId, period)
		if err != nil {
			return err
		}
		vat := CalcMonthlyVAT(revenue)
		surcharge := CalcSurcharge(vat)
		due := dueDateNextMonth15(time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0))
		if _, _, err = upsertTask(ctx, opcId, taxTypeVAT, period, vat, due); err != nil {
			return err
		}
		if _, _, err = upsertTask(ctx, opcId, taxTypeSurcharge, period, surcharge, due); err != nil {
			return err
		}
		if m%3 == 0 {
			q := m / 3
			qPeriod := fmt.Sprintf("Q%d-%d", q, year)
			profit := sumQuarterProfit(ctx, opcId, year, q)
			cit := CalcQuarterlyCIT(profit)
			qDue := dueDateMonthEnd(time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0))
			if _, _, err = upsertTask(ctx, opcId, taxTypeCITQuarterly, qPeriod, cit, qDue); err != nil {
				return err
			}
		}
	}
	annualPeriod := fmt.Sprintf("%d", year)
	annualProfit := sumAnnualProfit(ctx, opcId, year)
	annualCit := CalcAnnualCIT(annualProfit)
	if _, _, err := upsertTask(ctx, opcId, taxTypeCITAnnual, annualPeriod, annualCit, dueDateMay31(year+1)); err != nil {
		return err
	}
	return refreshOverdueStatus(ctx, opcId)
}

func refreshOverdueStatus(ctx context.Context, opcId uint64) error {
	now := uint64(utility.NowUnix())
	_, err := g.DB().Model(tableTaxFilingTask).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		Where("status", "pending").
		Where("due_date < ?", now).
		Where("due_date > ?", 0).
		Data(g.Map{"status": "overdue", "update_time": utility.NowUnix()}).
		Update()
	return err
}

func sumMonthlyGrossByPeriod(ctx context.Context, opcId uint64, period string) (float64, error) {
	start, end, err := rules.PeriodRange(period)
	if err != nil {
		return 0, err
	}
	var row struct {
		Total float64 `json:"total"`
	}
	err = g.DB().Model("xy_income_entry").Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		WhereBetween("occurred_at", start, end).
		Fields("COALESCE(SUM(gross_amount), 0) AS total").
		Scan(&row)
	if err != nil {
		return 0, gerror.Wrap(err, "汇总月收入失败")
	}
	return row.Total, nil
}

func sumAnnualProfit(ctx context.Context, opcId uint64, year int) float64 {
	var row struct {
		Total float64 `json:"total"`
	}
	_ = g.DB().Model(tableProfitSummary).Ctx(ctx).
		Fields("COALESCE(SUM(profit), 0) AS total").
		Where("opc_id", opcId).
		Where("year", year).
		Scan(&row)
	return row.Total
}
