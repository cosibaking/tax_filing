package tax

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/library/excelutil"
	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
	"xygo/utility"
)

var expenseKeywordBlacklist = []string{"咨询费", "服务费", "借款", "赠与"}

const tableTaxFilingSubmission = "xy_tax_filing_submission"

var periodPattern = regexp.MustCompile(`^\d{4}-\d{2}$`)

// GenerateFilingTemplate 生成报税 Excel 模板
func (s *sComplianceTax) GenerateFilingTemplate(ctx context.Context) ([]byte, error) {
	_ = ctx
	return excelutil.GenerateFilingTemplate()
}

// PreviewTaxFilingImport 报税 Excel 导入预览
func (s *sComplianceTax) PreviewTaxFilingImport(ctx context.Context, in *compliancein.TaxFilingImportInp) (*compliancein.TaxFilingImportPreviewModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	summary, expenseRows, err := excelutil.ParseFilingWorkbook(in.FileContent)
	if err != nil {
		return nil, err
	}

	codeToName, nameToCode, err := loadExpenseCategoryMaps(ctx)
	if err != nil {
		return nil, err
	}

	out := &compliancein.TaxFilingImportPreviewModel{
		Summary:     previewSummary(summary),
		ExpenseRows: make([]compliancein.TaxFilingExpensePreviewRow, 0, len(expenseRows)),
	}

	for _, row := range expenseRows {
		preview := compliancein.TaxFilingExpensePreviewRow{
			Row:         row.RowNum,
			OccurredAt:  row.OccurredAt,
			Category:    row.Category,
			Amount:      row.Amount,
			InvoiceType: row.InvoiceType,
			Description: row.Description,
		}
		code, name, rowErr := resolveExpenseCategory(row.Category, codeToName, nameToCode)
		preview.Category = code
		preview.CategoryName = name
		if rowErr != nil {
			preview.Valid = false
			preview.Error = rowErr.Error()
			out.InvalidExpenseCount++
		} else if _, parseErr := parseExpenseOccurredAt(row.OccurredAt); parseErr != nil {
			preview.Valid = false
			preview.Error = parseErr.Error()
			out.InvalidExpenseCount++
		} else if row.Amount <= 0 {
			preview.Valid = false
			preview.Error = "金额须大于0"
			out.InvalidExpenseCount++
		} else if err = checkExpenseKeyword(row.Description); err != nil {
			preview.Valid = false
			preview.Error = err.Error()
			out.InvalidExpenseCount++
		} else {
			preview.Valid = true
			out.ValidExpenseCount++
		}
		out.ExpenseRows = append(out.ExpenseRows, preview)
	}
	_ = opc
	return out, nil
}

// ImportTaxFilingExcel 报税 Excel 导入
func (s *sComplianceTax) ImportTaxFilingExcel(ctx context.Context, in *compliancein.TaxFilingImportInp) (*compliancein.TaxFilingImportModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	preview, err := s.PreviewTaxFilingImport(ctx, in)
	if err != nil {
		return nil, err
	}
	if !preview.Summary.Valid {
		return nil, gerror.New(preview.Summary.Error)
	}
	if preview.InvalidExpenseCount > 0 {
		return nil, gerror.Newf("存在 %d 条无效报销明细，请修正后重试", preview.InvalidExpenseCount)
	}

	now := uint64(utility.NowUnix())
	imported := 0
	failed := 0

	for _, row := range preview.ExpenseRows {
		if !row.Valid {
			failed++
			continue
		}
		_, createErr := service.ComplianceLedger().CreateExpense(ctx, &compliancein.ExpenseCreateInp{
			MemberId:       in.MemberId,
			Category:       row.Category,
			Amount:         row.Amount,
			InvoiceType:    normalizeInvoiceType(row.InvoiceType),
			Description:    row.Description,
			OccurredAt:     row.OccurredAt,
			ConfirmWarning: true,
			Ip:             in.Ip,
		})
		if createErr != nil {
			failed++
			continue
		}
		imported++
	}

	subData := g.Map{
		"opc_id":           opc.Id,
		"period":           preview.Summary.Period,
		"revenue":          preview.Summary.Revenue,
		"expense_total":    preview.Summary.ExpenseTotal,
		"profit":           preview.Summary.Profit,
		"vat_amount":       preview.Summary.VatAmount,
		"cit_amount":       preview.Summary.CitAmount,
		"surcharge_amount": preview.Summary.SurchargeAmount,
		"remark":           preview.Summary.Remark,
		"expense_imported": imported,
		"excel_file_id":    in.ExcelFileId,
		"status":           "submitted",
		"deleted":          0,
		"create_time":      now,
		"update_time":      now,
	}
	result, err := g.DB().Model(tableTaxFilingSubmission).Ctx(ctx).Data(subData).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "保存申报记录失败")
	}
	subId, _ := result.LastInsertId()
	_ = audit.WriteAudit(ctx, "tax_filing_submission", uint64(subId), "filing.import", in.MemberId, "member", nil, subData, in.Ip)

	return &compliancein.TaxFilingImportModel{
		SubmissionId:    uint64(subId),
		Period:          preview.Summary.Period,
		ExpenseImported: imported,
		ExpenseFailed:   failed,
	}, nil
}

// ListTaxFilingSubmissions 申报提交历史
func (s *sComplianceTax) ListTaxFilingSubmissions(ctx context.Context, in *compliancein.TaxFilingSubmissionListInp) (*compliancein.TaxFilingSubmissionListModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	m := g.DB().Model(tableTaxFilingSubmission).Ctx(ctx).
		Where("opc_id", opc.Id).
		Where("deleted", 0)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询申报记录失败")
	}

	var rows []struct {
		Id              uint64  `json:"id"`
		Period          string  `json:"period"`
		Revenue         float64 `json:"revenue"`
		ExpenseTotal    float64 `json:"expense_total"`
		Profit          float64 `json:"profit"`
		VatAmount       float64 `json:"vat_amount"`
		CitAmount       float64 `json:"cit_amount"`
		SurchargeAmount float64 `json:"surcharge_amount"`
		ExpenseImported int     `json:"expense_imported"`
		Remark          string  `json:"remark"`
		CreateTime      uint64  `json:"create_time"`
	}
	err = m.OrderDesc("create_time").Page(page, pageSize).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询申报记录失败")
	}

	list := make([]compliancein.TaxFilingSubmissionItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, compliancein.TaxFilingSubmissionItem{
			Id:              r.Id,
			Period:          r.Period,
			Revenue:         r.Revenue,
			ExpenseTotal:    r.ExpenseTotal,
			Profit:          r.Profit,
			VatAmount:       r.VatAmount,
			CitAmount:       r.CitAmount,
			SurchargeAmount: r.SurchargeAmount,
			ExpenseImported: r.ExpenseImported,
			Remark:          r.Remark,
			CreatedAt:       formatUnix(r.CreateTime),
		})
	}
	return &compliancein.TaxFilingSubmissionListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func previewSummary(summary *excelutil.FilingSummary) compliancein.TaxFilingSummaryPreview {
	out := compliancein.TaxFilingSummaryPreview{
		Period:          summary.Period,
		Revenue:         summary.Revenue,
		ExpenseTotal:    summary.ExpenseTotal,
		Profit:          summary.Profit,
		VatAmount:       summary.VatAmount,
		CitAmount:       summary.CitAmount,
		SurchargeAmount: summary.SurchargeAmount,
		Remark:          summary.Remark,
		Valid:           true,
	}
	if !periodPattern.MatchString(summary.Period) {
		out.Valid = false
		out.Error = "申报期间格式须为 YYYY-MM"
	}
	return out
}

func loadExpenseCategoryMaps(ctx context.Context) (codeToName, nameToCode map[string]string, err error) {
	var rows []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
	err = g.DB().Model("xy_expense_type").Ctx(ctx).
		Where("status", 1).
		Scan(&rows)
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询费用类型失败")
	}
	codeToName = make(map[string]string, len(rows))
	nameToCode = make(map[string]string, len(rows))
	for _, r := range rows {
		code := strings.ToLower(strings.TrimSpace(r.Code))
		name := strings.TrimSpace(r.Name)
		codeToName[code] = name
		nameToCode[name] = code
	}
	return codeToName, nameToCode, nil
}

func resolveExpenseCategory(input string, codeToName, nameToCode map[string]string) (code, name string, err error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", "", gerror.New("费用类型不能为空")
	}
	if c, ok := nameToCode[trimmed]; ok {
		return c, trimmed, nil
	}
	lower := strings.ToLower(trimmed)
	if n, ok := codeToName[lower]; ok {
		return lower, n, nil
	}
	return "", "", gerror.Newf("费用类型无效: %s", input)
}

func parseExpenseOccurredAt(s string) (uint64, error) {
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

func checkExpenseKeyword(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	for _, kw := range expenseKeywordBlacklist {
		if strings.Contains(text, kw) {
			return gerror.New("费用说明含禁止关键词：" + kw)
		}
	}
	return nil
}

func normalizeInvoiceType(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	aliases := map[string]string{
		"专票": "special", "增值税专用发票": "special", "special": "special",
		"普票": "general", "增值税普通发票": "general", "general": "general",
		"收据": "receipt", "receipt": "receipt",
		"无票": "none", "none": "none", "": "general",
	}
	if v, ok := aliases[raw]; ok {
		return v
	}
	return "general"
}
