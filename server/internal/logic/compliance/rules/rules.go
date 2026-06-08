package rules

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	tableIncomeEntry = "xy_income_entry"

	prohibitedIncomeCategories = "loan,gift,借款,赠与"
)

var prohibitedKeywords = []string{"借款", "赠与", "loan", "gift"}

// ValidateFiling 申报前合规校验（F-68）
func ValidateFiling(ctx context.Context, opcId uint64, period string, filedAmount, calculatedAmount, reportedIncome float64) error {
	ledgerIncome, err := SumLedgerIncome(ctx, opcId, period)
	if err != nil {
		return err
	}
	if reportedIncome > 0 {
		if err := CheckLowReportedIncome(ledgerIncome, reportedIncome); err != nil {
			return err
		}
		return nil
	}
	if filedAmount > 0 && calculatedAmount > 0 && filedAmount+0.01 < calculatedAmount {
		return gerror.New("申报税额低于系统测算值，存在低报收入风险，已阻断")
	}
	if filedAmount > 0 && ledgerIncome > 0 && filedAmount < ledgerIncome {
		return CheckLowReportedIncome(ledgerIncome, filedAmount)
	}
	return nil
}

// CheckLowReportedIncome 低报收入：申报值低于台账汇总则阻断
func CheckLowReportedIncome(ledgerIncome, declaredAmount float64) error {
	if ledgerIncome <= 0 || declaredAmount <= 0 {
		return nil
	}
	if declaredAmount+0.01 < ledgerIncome {
		return gerror.Newf("申报金额 %.2f 低于台账收入汇总 %.2f，禁止隐瞒收入", declaredAmount, ledgerIncome)
	}
	return nil
}

// CheckProhibitedIncomeCategory 禁止借款/赠与类收入
func CheckProhibitedIncomeCategory(category string) error {
	cat := strings.ToLower(strings.TrimSpace(category))
	for _, banned := range strings.Split(prohibitedIncomeCategories, ",") {
		if cat == strings.TrimSpace(banned) {
			return gerror.New("禁止选择借款/赠与类收入，请如实申报")
		}
	}
	return nil
}

// CheckIncomeRemark 备注关键词拦截
func CheckIncomeRemark(remark string) error {
	text := strings.TrimSpace(remark)
	if text == "" {
		return nil
	}
	lower := strings.ToLower(text)
	for _, kw := range prohibitedKeywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return gerror.Newf("收入备注含禁止关键词「%s」，请核实", kw)
		}
	}
	return nil
}

// SumLedgerIncome 汇总指定期间台账收入（net_amount）
func SumLedgerIncome(ctx context.Context, opcId uint64, period string) (float64, error) {
	start, end, err := periodRange(period)
	if err != nil {
		return 0, err
	}
	var row struct {
		Total float64 `json:"total"`
	}
	err = g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Fields("COALESCE(SUM(net_amount), 0) AS total").
		Where("opc_id", opcId).
		Where("deleted", 0).
		WhereBetween("occurred_at", start, end).
		Scan(&row)
	if err != nil {
		return 0, gerror.Wrap(err, "汇总台账收入失败")
	}
	return row.Total, nil
}
