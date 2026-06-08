package ledger

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var keywordBlacklist = []string{"咨询费", "服务费", "借款", "赠与"}

func costRatioThreshold(ctx context.Context) float64 {
	v, err := g.Cfg().Get(ctx, "compliance.costRatioThreshold")
	if err != nil || v.IsEmpty() {
		return 0.8
	}
	f := v.Float64()
	if f <= 0 || f > 1 {
		return 0.8
	}
	return f
}

func validateIncomeCategory(category string) error {
	category = strings.TrimSpace(strings.ToLower(category))
	if category == "" {
		return gerror.New("请选择收入类型")
	}
	if forbiddenIncomeCategories[category] {
		return gerror.New("不允许录入借款或赠与类收入")
	}
	if !validIncomeCategories[category] {
		return gerror.New("收入类型无效，允许：tip/commission/ad/slot_fee/offline/other")
	}
	return nil
}

func validatePlatform(platform string) error {
	platform = strings.TrimSpace(strings.ToLower(platform))
	if platform == "" {
		return gerror.New("请选择平台")
	}
	if !validPlatforms[platform] {
		return gerror.New("平台无效，允许：douyin/kuaishou/bilibili/wechat/xiaohongshu")
	}
	return nil
}

func checkKeywordBlacklist(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	for _, kw := range keywordBlacklist {
		if strings.Contains(text, kw) {
			return gerror.New("费用说明含禁止关键词：" + kw)
		}
	}
	return nil
}

func (s *sComplianceLedger) monthTotals(ctx context.Context, opcId uint64, monthStart, monthEnd uint64) (incomeNet, expenseTotal float64, err error) {
	if monthStart == 0 || monthEnd == 0 {
		return 0, 0, nil
	}
	var incomeRow struct {
		Total float64 `json:"total"`
	}
	err = g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		WhereBetween("occurred_at", monthStart, monthEnd).
		Fields("COALESCE(SUM(net_amount), 0) AS total").
		Scan(&incomeRow)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "统计收入失败")
	}

	var expenseRow struct {
		Total float64 `json:"total"`
	}
	err = g.DB().Model(tableExpenseEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		WhereBetween("occurred_at", monthStart, monthEnd).
		Fields("COALESCE(SUM(amount), 0) AS total").
		Scan(&expenseRow)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "统计费用失败")
	}
	return incomeRow.Total, expenseRow.Total, nil
}

