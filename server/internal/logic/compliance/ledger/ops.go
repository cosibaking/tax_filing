package ledger

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/tax"
	"xygo/internal/model/input/compliancein"
	"xygo/utility"
)

var accountLabels = map[string]string{
	accountBank:    "银行存款",
	accountRevenue: "主营业务收入",
	accountExpense: "管理费/推广费",
}

// DeleteIncome 软删除收入条目并回滚凭证与利润汇总
func (s *sComplianceLedger) DeleteIncome(ctx context.Context, in *compliancein.IncomeDeleteInp) error {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return err
	}
	if in.Id == 0 {
		return gerror.New("请指定收入条目")
	}

	var row struct {
		Id         uint64  `json:"id"`
		GrossAmount float64 `json:"gross_amount"`
		OccurredAt uint64  `json:"occurred_at"`
	}
	err = g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("id", in.Id).
		Where("opc_id", opcId).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return gerror.Wrap(err, "查询收入失败")
	}
	if row.Id == 0 {
		return gerror.New("收入条目不存在")
	}

	now := uint64(utility.NowUnix())
	_, err = g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("id", row.Id).
		Data(g.Map{"deleted": 1, "update_time": now}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除收入失败")
	}

	_, _ = g.DB().Model(tableLedgerVoucher).Ctx(ctx).
		Where("opc_id", opcId).
		Where("ref_type", refTypeIncome).
		Where("ref_id", row.Id).
		Delete()

	_, _ = g.DB().Model(tableBankTransaction).Ctx(ctx).
		Where("opc_id", opcId).
		Where("income_id", row.Id).
		Where("deleted", 0).
		Data(g.Map{"matched": 0, "income_id": 0, "update_time": now}).
		Update()

	year, month := yearMonthFromUnix(row.OccurredAt)
	_ = refreshProfitSummary(ctx, opcId, year, month)
	_ = tax.RefreshOpcPeriodTaxAmounts(ctx, opcId, year, month)
	_ = audit.WriteAudit(ctx, "income_entry", row.Id, "income.delete", in.MemberId, "member", row, nil, in.Ip)
	return nil
}

// DeleteExpense 软删除费用条目并回滚凭证与利润汇总
func (s *sComplianceLedger) DeleteExpense(ctx context.Context, in *compliancein.ExpenseDeleteInp) error {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return err
	}
	if in.Id == 0 {
		return gerror.New("请指定费用条目")
	}

	var row struct {
		Id         uint64  `json:"id"`
		Amount     float64 `json:"amount"`
		OccurredAt uint64  `json:"occurred_at"`
	}
	err = g.DB().Model(tableExpenseEntry).Ctx(ctx).
		Where("id", in.Id).
		Where("opc_id", opcId).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return gerror.Wrap(err, "查询费用失败")
	}
	if row.Id == 0 {
		return gerror.New("费用条目不存在")
	}

	now := uint64(utility.NowUnix())
	_, err = g.DB().Model(tableExpenseEntry).Ctx(ctx).
		Where("id", row.Id).
		Data(g.Map{"deleted": 1, "update_time": now}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除费用失败")
	}

	_, _ = g.DB().Model(tableLedgerVoucher).Ctx(ctx).
		Where("opc_id", opcId).
		Where("ref_type", refTypeExpense).
		Where("ref_id", row.Id).
		Delete()

	year, month := yearMonthFromUnix(row.OccurredAt)
	_ = refreshProfitSummary(ctx, opcId, year, month)
	_ = tax.RefreshOpcPeriodTaxAmounts(ctx, opcId, year, month)
	_ = audit.WriteAudit(ctx, "expense_entry", row.Id, "expense.delete", in.MemberId, "member", row, nil, in.Ip)
	return nil
}

// PreviewIncomeImport CSV 导入预览（不落库）
func (s *sComplianceLedger) PreviewIncomeImport(ctx context.Context, in *compliancein.IncomeImportInp) (*compliancein.IncomeImportPreviewModel, error) {
	if _, err := s.requireActiveOpc(ctx, in.MemberId); err != nil {
		return nil, err
	}
	content := strings.TrimSpace(in.CsvContent)
	if content == "" {
		return nil, gerror.New("CSV 内容为空")
	}

	reader := csv.NewReader(strings.NewReader(content))
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, gerror.Wrap(err, "解析 CSV 失败")
	}
	if len(records) < 2 {
		return nil, gerror.New("CSV 须包含表头与至少一行数据")
	}

	header := normalizeCSVHeader(records[0])
	col := mapCSVColumns(header, []string{"platform", "category", "gross_amount", "platform_fee", "occurred_at", "remark"})
	if col["platform"] < 0 || col["category"] < 0 || col["gross_amount"] < 0 || col["occurred_at"] < 0 {
		return nil, gerror.New("CSV 表头须包含 platform,category,gross_amount,occurred_at")
	}

	out := &compliancein.IncomeImportPreviewModel{Rows: make([]compliancein.IncomeImportPreviewRow, 0)}
	for i := 1; i < len(records); i++ {
		row := records[i]
		if isEmptyCSVRow(row) {
			continue
		}
		rowNum := i + 1
		item, rowErr := parseIncomeCSVRow(row, col)
		preview := compliancein.IncomeImportPreviewRow{
			Row:         rowNum,
			Platform:    cell(row, col["platform"]),
			Category:    cell(row, col["category"]),
			GrossAmount: 0,
			PlatformFee: 0,
			Valid:       rowErr == nil,
		}
		if rowErr != nil {
			preview.Error = rowErr.Error()
			out.InvalidCount++
		} else {
			preview.OccurredAt = formatOccurredAt(item.occurredAt)
			preview.GrossAmount = item.gross
			preview.PlatformFee = item.fee
			out.ValidCount++
		}
		out.Rows = append(out.Rows, preview)
	}
	return out, nil
}

// ListBankUnmatched 未匹配银行流水列表
func (s *sComplianceLedger) ListBankUnmatched(ctx context.Context, in *compliancein.BankUnmatchedInp) (*compliancein.BankUnmatchedModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}

	m := g.DB().Model(tableBankTransaction).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		Where("matched", 0)
	if in.Month != "" {
		start, end, monthErr := monthRange(in.Month)
		if monthErr != nil {
			return nil, monthErr
		}
		m = m.WhereBetween("occurred_at", start, end)
	}

	count, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计未匹配流水失败")
	}

	var rows []struct {
		Id          uint64  `json:"id"`
		OccurredAt  uint64  `json:"occurred_at"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	err = m.OrderDesc("occurred_at").Limit(100).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询未匹配流水失败")
	}

	list := make([]compliancein.BankUnmatchedItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, compliancein.BankUnmatchedItem{
			Id:          r.Id,
			OccurredAt:  formatOccurredAt(r.OccurredAt),
			Amount:      r.Amount,
			Description: r.Description,
		})
	}
	return &compliancein.BankUnmatchedModel{List: list, Count: count}, nil
}

// ListLedgerVouchers 会计分录列表
func (s *sComplianceLedger) ListLedgerVouchers(ctx context.Context, in *compliancein.LedgerVoucherListInp) (*compliancein.LedgerVoucherListModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	period := strings.TrimSpace(in.Period)
	if period == "" {
		period = periodFromUnix(uint64(utility.NowUnix()))
	}

	var rows []struct {
		Id            uint64  `json:"id"`
		Period        string  `json:"period"`
		DebitAccount  string  `json:"debit_account"`
		CreditAccount string  `json:"credit_account"`
		Amount        float64 `json:"amount"`
		RefType       string  `json:"ref_type"`
		CreateTime    uint64  `json:"create_time"`
	}
	m := g.DB().Model(tableLedgerVoucher).Ctx(ctx).Where("opc_id", opcId)
	if len(period) == 4 {
		m = m.WhereLike("period", period+"-%")
	} else {
		m = m.Where("period", period)
	}
	err = m.OrderDesc("create_time").OrderDesc("id").Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询会计分录失败")
	}

	list := make([]compliancein.LedgerVoucherItem, 0, len(rows))
	for _, r := range rows {
		summary := voucherSummary(r.RefType, r.DebitAccount, r.CreditAccount)
		list = append(list, compliancein.LedgerVoucherItem{
			Id:            r.Id,
			OccurredAt:    formatOccurredAt(r.CreateTime),
			Summary:       summary,
			DebitAccount:  accountLabel(r.DebitAccount),
			CreditAccount: accountLabel(r.CreditAccount),
			Amount:        r.Amount,
		})
	}
	return &compliancein.LedgerVoucherListModel{List: list}, nil
}

func accountLabel(code string) string {
	if label, ok := accountLabels[code]; ok {
		return label
	}
	return code
}

func voucherSummary(refType, debit, credit string) string {
	switch refType {
	case refTypeIncome:
		return fmt.Sprintf("确认收入：借 %s / 贷 %s", accountLabel(debit), accountLabel(credit))
	case refTypeExpense:
		return fmt.Sprintf("确认费用：借 %s / 贷 %s", accountLabel(debit), accountLabel(credit))
	default:
		return fmt.Sprintf("借 %s / 贷 %s", accountLabel(debit), accountLabel(credit))
	}
}
