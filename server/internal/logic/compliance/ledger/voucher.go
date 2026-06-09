package ledger

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/utility"
)

const (
	accountBank       = "1001"
	accountRevenue    = "6001"
	accountExpense    = "6401"
	refTypeIncome     = "income"
	refTypeExpense    = "expense"
)

// RecordIncomeVoucher 确认收入时自动生成凭证：借银行存款 / 贷主营业务收入
func RecordIncomeVoucher(ctx context.Context, opcId, refId uint64, amount float64, occurredAt uint64) error {
	if opcId == 0 || amount <= 0 {
		return nil
	}
	period := periodFromUnix(occurredAt)
	now := utility.NowUnix()

	count, err := g.DB().Model(tableLedgerVoucher).Ctx(ctx).
		Where("opc_id", opcId).
		Where("ref_type", refTypeIncome).
		Where("ref_id", refId).
		Count()
	if err != nil {
		return gerror.Wrap(err, "查询收入凭证失败")
	}
	if count > 0 {
		year, month := yearMonthFromUnix(occurredAt)
		return refreshProfitSummary(ctx, opcId, year, month)
	}

	_, err = g.DB().Model(tableLedgerVoucher).Ctx(ctx).Data(g.Map{
		"opc_id":         opcId,
		"period":         period,
		"debit_account":  accountBank,
		"credit_account": accountRevenue,
		"amount":         amount,
		"ref_type":       refTypeIncome,
		"ref_id":         refId,
		"create_time":    now,
	}).Insert()
	if err != nil {
		return gerror.Wrap(err, "写入收入凭证失败")
	}

	year, month := yearMonthFromUnix(occurredAt)
	return refreshProfitSummary(ctx, opcId, year, month)
}

// RecordExpenseVoucher 确认费用时自动生成凭证：借管理费/推广费 / 贷银行存款
func RecordExpenseVoucher(ctx context.Context, opcId, refId uint64, amount float64, occurredAt uint64) error {
	if opcId == 0 || amount <= 0 {
		return nil
	}
	period := periodFromUnix(occurredAt)
	now := utility.NowUnix()

	count, err := g.DB().Model(tableLedgerVoucher).Ctx(ctx).
		Where("opc_id", opcId).
		Where("ref_type", refTypeExpense).
		Where("ref_id", refId).
		Count()
	if err != nil {
		return gerror.Wrap(err, "查询费用凭证失败")
	}
	if count > 0 {
		year, month := yearMonthFromUnix(occurredAt)
		return refreshProfitSummary(ctx, opcId, year, month)
	}

	_, err = g.DB().Model(tableLedgerVoucher).Ctx(ctx).Data(g.Map{
		"opc_id":         opcId,
		"period":         period,
		"debit_account":  accountExpense,
		"credit_account": accountBank,
		"amount":         amount,
		"ref_type":       refTypeExpense,
		"ref_id":         refId,
		"create_time":    now,
	}).Insert()
	if err != nil {
		return gerror.Wrap(err, "写入费用凭证失败")
	}

	year, month := yearMonthFromUnix(occurredAt)
	return refreshProfitSummary(ctx, opcId, year, month)
}

func refreshProfitSummary(ctx context.Context, opcId uint64, year, month int) error {
	if opcId == 0 || year <= 0 || month <= 0 {
		return nil
	}

	startTs, endTs, err := monthRange(fmt.Sprintf("%04d-%02d", year, month))
	if err != nil {
		return err
	}
	revenue, err := sumIncome(ctx, opcId, startTs, endTs)
	if err != nil {
		return err
	}
	cost, err := sumExpense(ctx, opcId, startTs, endTs)
	if err != nil {
		return err
	}
	profit := round2(revenue - cost)
	cumulative, err := sumYearProfitBefore(ctx, opcId, year, month, profit)
	if err != nil {
		return err
	}

	now := utility.NowUnix()
	_, err = g.DB().Model(tableProfitSummary).Ctx(ctx).
		Data(g.Map{
			"opc_id":             opcId,
			"year":               year,
			"month":              month,
			"revenue":            revenue,
			"cost":               cost,
			"profit":             profit,
			"cumulative_profit":  cumulative,
			"update_time":        now,
		}).
		OnDuplicate(g.Map{
			"revenue":           revenue,
			"cost":              cost,
			"profit":            profit,
			"cumulative_profit": cumulative,
			"update_time":       now,
		}).
		Insert()
	if err != nil {
		return gerror.Wrap(err, "更新利润汇总失败")
	}
	return nil
}

func sumIncome(ctx context.Context, opcId uint64, startTs, endTs uint64) (float64, error) {
	var row struct {
		Total float64 `json:"total"`
	}
	err := g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		WhereBetween("occurred_at", startTs, endTs).
		Fields("COALESCE(SUM(gross_amount), 0) AS total").
		Scan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, gerror.Wrap(err, "汇总收入失败")
	}
	return round2(row.Total), nil
}

func sumExpense(ctx context.Context, opcId uint64, startTs, endTs uint64) (float64, error) {
	var row struct {
		Total float64 `json:"total"`
	}
	err := g.DB().Model(tableExpenseEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		WhereBetween("occurred_at", startTs, endTs).
		Fields("COALESCE(SUM(amount), 0) AS total").
		Scan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, gerror.Wrap(err, "汇总费用失败")
	}
	return round2(row.Total), nil
}

func sumYearProfitBefore(ctx context.Context, opcId uint64, year, month int, currentProfit float64) (float64, error) {
	var row struct {
		Total float64 `json:"total"`
	}
	err := g.DB().Model(tableProfitSummary).Ctx(ctx).
		Where("opc_id", opcId).
		Where("year", year).
		Where("month < ?", month).
		Fields("COALESCE(SUM(profit), 0) AS total").
		Scan(&row)
	if err != nil {
		return 0, gerror.Wrap(err, "汇总累计利润失败")
	}
	return round2(row.Total + currentProfit), nil
}

func periodFromUnix(ts uint64) string {
	if ts == 0 {
		ts = uint64(utility.NowUnix())
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m")
}

func yearMonthFromUnix(ts uint64) (year, month int) {
	if ts == 0 {
		ts = uint64(utility.NowUnix())
	}
	t := utility.UnixToGTime(int64(ts))
	return t.Year(), int(t.Month())
}

func monthRangeUnix(year, month int) (startTs, endTs uint64) {
	start := gtime.NewFromStr(fmt.Sprintf("%04d-%02d-01 00:00:00", year, month))
	end := start.AddDate(0, 1, 0).Add(-1)
	return uint64(start.Unix()), uint64(end.Unix())
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}
