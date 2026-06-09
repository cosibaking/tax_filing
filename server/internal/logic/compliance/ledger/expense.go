package ledger

import (
	"context"
	"encoding/csv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/tax"
	"xygo/internal/model/input/compliancein"
	"xygo/utility"
)

// ListExpense 费用台账列表
func (s *sComplianceLedger) ListExpense(ctx context.Context, in *compliancein.ExpenseListInp) (*compliancein.ExpenseListModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	page, pageSize := normalizePage(in.Page, in.PageSize)

	m := g.DB().Model(tableExpenseEntry+" ee").Ctx(ctx).
		LeftJoin(tableExpenseType+" et", "et.code = ee.category AND et.status = 1").
		Where("ee.opc_id", opcId).
		Where("ee.deleted", 0)
	if in.Category != "" {
		m = m.Where("ee.category", strings.ToLower(strings.TrimSpace(in.Category)))
	}
	var monthStart, monthEnd uint64
	if in.Month != "" {
		monthStart, monthEnd, err = monthRange(in.Month)
		if err != nil {
			return nil, err
		}
		m = m.WhereBetween("ee.occurred_at", monthStart, monthEnd)
	}

	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询费用总数失败")
	}

	var rows []struct {
		Id           uint64  `json:"id"`
		Category     string  `json:"category"`
		CategoryName string  `json:"category_name"`
		Amount       float64 `json:"amount"`
		InvoiceType  string  `json:"invoice_type"`
		AttachmentId uint64  `json:"attachment_id"`
		Description  string  `json:"description"`
		WarningFlag  int     `json:"warning_flag"`
		OccurredAt   uint64  `json:"occurred_at"`
	}
	err = m.Fields(
		"ee.id", "ee.category", "et.name AS category_name",
		"ee.amount", "ee.invoice_type", "ee.attachment_id",
		"ee.description", "ee.warning_flag", "ee.occurred_at",
	).OrderDesc("ee.occurred_at").OrderDesc("ee.id").
		Page(page, pageSize).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询费用列表失败")
	}

	list := make([]compliancein.ExpenseItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, compliancein.ExpenseItem{
			Id:           r.Id,
			Category:     r.Category,
			CategoryName: r.CategoryName,
			Amount:       r.Amount,
			InvoiceType:  r.InvoiceType,
			AttachmentId: r.AttachmentId,
			Description:  r.Description,
			WarningFlag:  r.WarningFlag == 1,
			OccurredAt:   formatOccurredAt(r.OccurredAt),
		})
	}

	summary := compliancein.ExpenseSummary{}
	sumM := g.DB().Model(tableExpenseEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0)
	if in.Month != "" {
		sumM = sumM.WhereBetween("occurred_at", monthStart, monthEnd)
	}
	var sumRow struct {
		Total float64 `json:"total"`
	}
	_ = sumM.Fields("COALESCE(SUM(amount), 0) AS total").Scan(&sumRow)
	summary.TotalAmount = sumRow.Total

	if in.Month != "" {
		incomeNet, _, _ := s.monthTotals(ctx, opcId, monthStart, monthEnd)
		summary.CostRatio, _ = evaluateCostRatioCtx(ctx, incomeNet, sumRow.Total, 0)
	}

	return &compliancein.ExpenseListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Summary:  summary,
	}, nil
}

// CreateExpense 创建费用条目
func (s *sComplianceLedger) CreateExpense(ctx context.Context, in *compliancein.ExpenseCreateInp) (*compliancein.ExpenseCreateModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	category := strings.ToLower(strings.TrimSpace(in.Category))
	if category == "" {
		return nil, gerror.New("请选择费用类型")
	}
	if !s.expenseTypeExists(ctx, category) {
		return nil, gerror.New("费用类型无效")
	}
	if in.Amount <= 0 {
		return nil, gerror.New("费用金额须大于0")
	}
	invoiceType := strings.ToLower(strings.TrimSpace(in.InvoiceType))
	if invoiceType == "" {
		invoiceType = "general"
	}
	if !validInvoiceTypes[invoiceType] {
		return nil, gerror.New("发票类型无效，允许：general/special/receipt/none")
	}
	if err = checkKeywordBlacklist(in.Description); err != nil {
		return nil, err
	}
	occurredAt, err := parseOccurredAt(in.OccurredAt)
	if err != nil {
		return nil, err
	}
	if in.AttachmentId > 0 {
		if err = validateLedgerAttachment(ctx, in.AttachmentId); err != nil {
			return nil, err
		}
	}

	monthStart, monthEnd := monthStartEndFromUnix(occurredAt)
	incomeNet, expenseTotal, err := s.monthTotals(ctx, opcId, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}
	costRatio, warning := evaluateCostRatioCtx(ctx, incomeNet, expenseTotal, in.Amount)
	if warning && !in.ConfirmWarning {
		threshold := costRatioThreshold(ctx)
		return nil, gerror.Newf("当月成本占收入比例已达 %.0f%%，超过 %.0f%% 阈值，请确认风险后设置 confirmWarning=true 重试",
			costRatio*100, threshold*100)
	}

	now := uint64(utility.NowUnix())
	warningFlag := 0
	warningMsg := ""
	if warning {
		warningFlag = 1
		warningMsg = "当月成本占收入比例偏高，请核实费用合理性"
	}

	data := g.Map{
		"opc_id":         opcId,
		"category":       category,
		"amount":         roundMoney(in.Amount),
		"invoice_type":   invoiceType,
		"attachment_id":  in.AttachmentId,
		"description":    strings.TrimSpace(in.Description),
		"warning_flag":   warningFlag,
		"occurred_at":    occurredAt,
		"deleted":        0,
		"create_time":    now,
		"update_time":    now,
	}

	result, err := g.DB().Model(tableExpenseEntry).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "创建费用失败")
	}
	id, _ := result.LastInsertId()
	_ = audit.WriteAudit(ctx, "expense_entry", uint64(id), "expense.create", in.MemberId, "member", nil, data, in.Ip)
	_ = RecordExpenseVoucher(ctx, opcId, uint64(id), in.Amount, occurredAt)
	year, month := yearMonthFromUnix(occurredAt)
	_ = refreshProfitSummary(ctx, opcId, year, month)
	_ = tax.RefreshOpcPeriodTaxAmounts(ctx, opcId, year, month)

	return &compliancein.ExpenseCreateModel{
		Id:           uint64(id),
		Category:     category,
		Amount:       in.Amount,
		InvoiceType:  invoiceType,
		AttachmentId: in.AttachmentId,
		Description:  strings.TrimSpace(in.Description),
		WarningFlag:  warningFlag == 1,
		WarningMsg:   warningMsg,
		OccurredAt:   formatOccurredAt(occurredAt),
		CostRatio:    costRatio,
	}, nil
}

// ListExpenseTypes 费用类型库
func (s *sComplianceLedger) ListExpenseTypes(ctx context.Context) (*compliancein.ExpenseTypesModel, error) {
	var rows []struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		VoucherHint   string `json:"voucher_hint"`
		ComplianceTip string `json:"compliance_tip"`
	}
	err := g.DB().Model(tableExpenseType).Ctx(ctx).
		Where("status", 1).
		OrderAsc("sort").
		OrderAsc("id").
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询费用类型库失败")
	}
	list := make([]compliancein.ExpenseTypeItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, compliancein.ExpenseTypeItem{
			Code:          r.Code,
			Name:          r.Name,
			VoucherHint:   r.VoucherHint,
			ComplianceTip: r.ComplianceTip,
		})
	}
	return &compliancein.ExpenseTypesModel{List: list}, nil
}

// ImportBankStatement 银行流水导入与对账（stub）
func (s *sComplianceLedger) ImportBankStatement(ctx context.Context, in *compliancein.BankImportInp) (*compliancein.BankImportModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
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
	col := mapCSVColumns(header, []string{"occurred_at", "amount", "description"})
	if col["occurred_at"] < 0 || col["amount"] < 0 {
		return nil, gerror.New("CSV 表头须包含 occurred_at,amount（description 可选）")
	}

	out := &compliancein.BankImportModel{}
	now := uint64(utility.NowUnix())

	for i := 1; i < len(records); i++ {
		row := records[i]
		if isEmptyCSVRow(row) {
			continue
		}
		rowNum := i + 1

		occurredAt, parseErr := parseOccurredAt(cell(row, col["occurred_at"]))
		if parseErr != nil {
			out.Failed++
			out.Errors = append(out.Errors, compliancein.IncomeImportRowError{Row: rowNum, Message: parseErr.Error()})
			continue
		}
		amount, parseErr := parseMoney(cell(row, col["amount"]))
		if parseErr != nil || amount == 0 {
			out.Failed++
			out.Errors = append(out.Errors, compliancein.IncomeImportRowError{Row: rowNum, Message: "amount 无效"})
			continue
		}
		desc := ""
		if col["description"] >= 0 {
			desc = strings.TrimSpace(cell(row, col["description"]))
		}

		matched, incomeId := s.reconcileBankTx(ctx, opcId, occurredAt, amount)

		data := g.Map{
			"opc_id":       opcId,
			"occurred_at":  occurredAt,
			"amount":       roundMoney(amount),
			"description":  desc,
			"matched":      boolToInt(matched),
			"income_id":    incomeId,
			"deleted":      0,
			"create_time":  now,
			"update_time":  now,
		}
		result, insertErr := g.DB().Model(tableBankTransaction).Ctx(ctx).Data(data).Insert()
		if insertErr != nil {
			out.Failed++
			out.Errors = append(out.Errors, compliancein.IncomeImportRowError{Row: rowNum, Message: "写入失败"})
			continue
		}
		txId, _ := result.LastInsertId()
		_ = audit.WriteAudit(ctx, "bank_transaction", uint64(txId), "bank.import", in.MemberId, "member", nil, data, in.Ip)
		out.Imported++
		if matched {
			out.Matched++
		} else {
			out.Unmatched++
		}
	}
	return out, nil
}

func (s *sComplianceLedger) reconcileBankTx(ctx context.Context, opcId uint64, occurredAt uint64, amount float64) (matched bool, incomeId uint64) {
	if amount <= 0 {
		return false, 0
	}
	window := uint64(3 * 86400)
	from := occurredAt
	if occurredAt > window {
		from = occurredAt - window
	}
	to := occurredAt + window

	var row struct {
		Id uint64 `json:"id"`
	}
	_ = g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0).
		WhereBetween("occurred_at", from, to).
		Where("net_amount", roundMoney(amount)).
		OrderDesc("occurred_at").
		Limit(1).
		Scan(&row)
	if row.Id == 0 {
		return false, 0
	}
	return true, row.Id
}

func (s *sComplianceLedger) expenseTypeExists(ctx context.Context, code string) bool {
	count, err := g.DB().Model(tableExpenseType).Ctx(ctx).
		Where("code", code).
		Where("status", 1).
		Count()
	return err == nil && count > 0
}

func validateLedgerAttachment(ctx context.Context, fileId uint64) error {
	var att struct {
		Id       uint64 `json:"id"`
		Mimetype string `json:"mimetype"`
	}
	err := g.DB().Model("xy_sys_attachment").Ctx(ctx).Where("id", fileId).Scan(&att)
	if err != nil {
		return gerror.Wrap(err, "查询附件失败")
	}
	if att.Id == 0 {
		return gerror.New("附件不存在")
	}
	allowed := map[string]bool{
		"image/jpeg": true, "image/png": true, "application/pdf": true,
	}
	if !allowed[att.Mimetype] {
		return gerror.New("附件格式不支持，仅允许 jpg/png/pdf")
	}
	return nil
}

func evaluateCostRatioCtx(ctx context.Context, incomeNet, expenseTotal, extraAmount float64) (ratio float64, warning bool) {
	totalExpense := expenseTotal + extraAmount
	threshold := costRatioThreshold(ctx)
	if incomeNet <= 0 {
		if totalExpense > 0 {
			return 1, true
		}
		return 0, false
	}
	ratio = roundMoney(totalExpense / incomeNet)
	return ratio, ratio > threshold
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
