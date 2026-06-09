package excelutil

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/xuri/excelize/v2"
)

const (
	SheetSummary = "申报汇总"
	SheetExpense = "报销明细"
)

// FilingSummary 申报汇总
type FilingSummary struct {
	Period          string
	Revenue         float64
	ExpenseTotal    float64
	Profit          float64
	VatAmount       float64
	CitAmount       float64
	SurchargeAmount float64
	Remark          string
}

// FilingExpenseRow 报销明细行
type FilingExpenseRow struct {
	RowNum      int
	OccurredAt  string
	Category    string
	Amount      float64
	InvoiceType string
	Description string
}

// ParseFilingWorkbook 解析报税 Excel（双 Sheet）
func ParseFilingWorkbook(content []byte) (*FilingSummary, []FilingExpenseRow, error) {
	if len(content) == 0 {
		return nil, nil, gerror.New("Excel 文件为空")
	}
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return nil, nil, gerror.Wrap(err, "解析 Excel 失败")
	}
	defer f.Close()

	summarySheet := pickSheet(f, []string{SheetSummary, "summary", "申报"})
	expenseSheet := pickSheet(f, []string{SheetExpense, "expense", "报销"})

	summary, err := parseSummarySheet(f, summarySheet)
	if err != nil {
		return nil, nil, err
	}
	expenses, err := parseExpenseSheet(f, expenseSheet)
	if err != nil {
		return nil, nil, err
	}
	return summary, expenses, nil
}

// GenerateFilingTemplate 生成报税 Excel 模板
func GenerateFilingTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	_, _ = f.NewSheet(SheetSummary)
	_, _ = f.NewSheet(SheetExpense)
	_ = f.DeleteSheet("Sheet1")

	summaryHeaders := []string{"申报期间", "营业收入", "成本费用合计", "利润总额", "增值税", "企业所得税", "附加税", "备注"}
	for i, h := range summaryHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(SheetSummary, cell, h)
	}
	example := []interface{}{"2026-06", 100000, 30000, 70000, 990, 1750, 59, "示例数据，请删除后填写"}
	for i, v := range example {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		_ = f.SetCellValue(SheetSummary, cell, v)
	}

	expenseHeaders := []string{"发生日期", "费用类型", "金额", "发票类型", "费用说明"}
	for i, h := range expenseHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(SheetExpense, cell, h)
	}
	examples := [][]interface{}{
		{"2026-06-05", "设备", 5000, "专票", "直播设备采购"},
		{"2026-06-12", "网费", 200, "普票", "办公宽带"},
	}
	for r, row := range examples {
		for c, v := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			_ = f.SetCellValue(SheetExpense, cell, v)
		}
	}

	f.SetActiveSheet(0)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, gerror.Wrap(err, "生成模板失败")
	}
	return buf.Bytes(), nil
}

func pickSheet(f *excelize.File, candidates []string) string {
	sheets := f.GetSheetList()
	for _, c := range candidates {
		for _, s := range sheets {
			if strings.EqualFold(strings.TrimSpace(s), strings.TrimSpace(c)) {
				return s
			}
		}
	}
	if len(sheets) > 0 {
		return sheets[0]
	}
	return SheetSummary
}

func parseSummarySheet(f *excelize.File, sheet string) (*FilingSummary, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, gerror.Wrap(err, "读取申报汇总失败")
	}
	if len(rows) < 2 {
		return nil, gerror.New("申报汇总须包含表头与至少一行数据")
	}
	header := normalizeHeader(rows[0])
	col := mapHeaderCols(header, map[string][]string{
		"period":           {"申报期间", "period", "期间"},
		"revenue":          {"营业收入", "revenue", "收入"},
		"expense_total":    {"成本费用合计", "expense_total", "成本费用", "费用合计"},
		"profit":           {"利润总额", "profit", "利润"},
		"vat_amount":       {"增值税", "vat_amount", "vat"},
		"cit_amount":       {"企业所得税", "cit_amount", "cit", "所得税"},
		"surcharge_amount": {"附加税", "surcharge_amount", "surcharge"},
		"remark":           {"备注", "remark"},
	})
	if col["period"] < 0 {
		return nil, gerror.New("申报汇总缺少「申报期间」列")
	}

	dataRow := rows[1]
	for i := 1; i < len(rows); i++ {
		if !isEmptyRow(rows[i]) {
			dataRow = rows[i]
			break
		}
	}

	summary := &FilingSummary{
		Period: strings.TrimSpace(cell(dataRow, col["period"])),
		Remark: strings.TrimSpace(cell(dataRow, col["remark"])),
	}
	if summary.Period == "" {
		return nil, gerror.New("申报期间不能为空")
	}

	var parseErr error
	summary.Revenue, parseErr = parseAmount(cell(dataRow, col["revenue"]))
	if parseErr != nil {
		return nil, gerror.Wrap(parseErr, "营业收入格式无效")
	}
	summary.ExpenseTotal, _ = parseAmount(cell(dataRow, col["expense_total"]))
	summary.Profit, _ = parseAmount(cell(dataRow, col["profit"]))
	summary.VatAmount, _ = parseAmount(cell(dataRow, col["vat_amount"]))
	summary.CitAmount, _ = parseAmount(cell(dataRow, col["cit_amount"]))
	summary.SurchargeAmount, _ = parseAmount(cell(dataRow, col["surcharge_amount"]))
	return summary, nil
}

func parseExpenseSheet(f *excelize.File, sheet string) ([]FilingExpenseRow, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, gerror.Wrap(err, "读取报销明细失败")
	}
	if len(rows) < 2 {
		return []FilingExpenseRow{}, nil
	}
	header := normalizeHeader(rows[0])
	col := mapHeaderCols(header, map[string][]string{
		"occurred_at":  {"发生日期", "occurred_at", "日期", "date"},
		"category":     {"费用类型", "category", "类型"},
		"amount":       {"金额", "amount"},
		"invoice_type": {"发票类型", "invoice_type", "发票"},
		"description":  {"费用说明", "description", "说明", "备注"},
	})
	if col["occurred_at"] < 0 || col["category"] < 0 || col["amount"] < 0 {
		return nil, gerror.New("报销明细表头须包含：发生日期、费用类型、金额")
	}

	out := make([]FilingExpenseRow, 0)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if isEmptyRow(row) {
			continue
		}
		amount, amountErr := parseAmount(cell(row, col["amount"]))
		item := FilingExpenseRow{
			RowNum:      i + 1,
			OccurredAt:  strings.TrimSpace(cell(row, col["occurred_at"])),
			Category:    strings.TrimSpace(cell(row, col["category"])),
			Amount:      amount,
			InvoiceType: strings.TrimSpace(cell(row, col["invoice_type"])),
			Description: strings.TrimSpace(cell(row, col["description"])),
		}
		if amountErr != nil {
			item.Amount = 0
		}
		out = append(out, item)
	}
	return out, nil
}

func normalizeHeader(row []string) []string {
	out := make([]string, len(row))
	for i, h := range row {
		out[i] = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(h, "\ufeff", "")))
	}
	return out
}

func mapHeaderCols(header []string, aliases map[string][]string) map[string]int {
	col := make(map[string]int, len(aliases))
	for k := range aliases {
		col[k] = -1
	}
	for i, h := range header {
		for key, alts := range aliases {
			for _, a := range alts {
				if h == strings.ToLower(a) {
					col[key] = i
				}
			}
		}
	}
	return col
}

func cell(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func isEmptyRow(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}

func parseAmount(s string) (float64, error) {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	s = strings.ReplaceAll(s, "¥", "")
	s = strings.ReplaceAll(s, "元", "")
	if s == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("金额格式无效: %s", s)
	}
	return v, nil
}
