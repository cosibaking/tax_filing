package ledger

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/logic/compliance/shared"
	"xygo/utility"
)

const (
	tableOpcEntity       = "xy_opc_entity"
	tableServiceOrder    = "xy_service_order"
	tableIncomeEntry     = "xy_income_entry"
	tableExpenseEntry    = "xy_expense_entry"
	tableExpenseType     = "xy_expense_type"
	tableBankTransaction = "xy_bank_transaction"
	tableProfitSummary   = "xy_profit_summary"
	tableLedgerVoucher   = "xy_ledger_voucher"

	sourceManual = "manual"
	sourceCSV    = "csv"
)

var (
	validPlatforms = map[string]bool{
		"douyin": true, "kuaishou": true, "bilibili": true,
		"channels": true, "wechat": true, "xiaohongshu": true,
		"taobao": true, "alipay": true, "offline": true, "other": true,
	}
	validIncomeCategories = map[string]bool{
		"tip": true, "commission": true, "ad": true,
		"service_fee": true, "product_sales": true,
		"slot_fee": true, "offline": true, "other": true,
	}
	forbiddenIncomeCategories = map[string]bool{
		"loan": true, "gift": true,
	}
	validInvoiceTypes = map[string]bool{
		"general": true, "special": true, "receipt": true, "none": true,
	}
)

func (s *sComplianceLedger) requireActiveOpc(ctx context.Context, memberId uint64) (uint64, error) {
	opc, err := shared.RequireActiveOpc(ctx, memberId)
	if err != nil {
		return 0, err
	}
	return opc.Id, nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return page, pageSize
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}

func calcNet(gross, fee float64) float64 {
	return roundMoney(gross - fee)
}

func parseOccurredAt(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, gerror.New("请填写发生日期")
	}
	layouts := []string{"2006-01-02", "2006-01-02 15:04:05", "2006/01/02"}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return uint64(t.Unix()), nil
		}
	}
	t, err := gtime.StrToTime(s)
	if err != nil {
		return 0, gerror.New("发生日期格式无效，请使用 YYYY-MM-DD")
	}
	return uint64(t.Unix()), nil
}

func formatOccurredAt(ts uint64) string {
	if ts == 0 {
		return ""
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m-d")
}

func monthRange(month string) (uint64, uint64, error) {
	month = strings.TrimSpace(month)
	if month == "" {
		return 0, 0, nil
	}
	t, err := time.ParseInLocation("2006-01", month, time.Local)
	if err != nil {
		return 0, 0, gerror.New("月份格式无效，请使用 YYYY-MM")
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return uint64(start.Unix()), uint64(end.Unix()), nil
}
