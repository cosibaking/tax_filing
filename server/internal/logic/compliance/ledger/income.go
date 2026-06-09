package ledger

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/tax"
	"xygo/internal/model/input/compliancein"
	"xygo/utility"
)

// ListIncome 收入台账列表
func (s *sComplianceLedger) ListIncome(ctx context.Context, in *compliancein.IncomeListInp) (*compliancein.IncomeListModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	page, pageSize := normalizePage(in.Page, in.PageSize)

	m := g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0)
	if in.Platform != "" {
		m = m.Where("platform", strings.ToLower(strings.TrimSpace(in.Platform)))
	}
	if in.Category != "" {
		m = m.Where("category", strings.ToLower(strings.TrimSpace(in.Category)))
	}
	if in.Month != "" {
		start, end, monthErr := monthRange(in.Month)
		if monthErr != nil {
			return nil, monthErr
		}
		m = m.WhereBetween("occurred_at", start, end)
	}

	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询收入总数失败")
	}

	var rows []struct {
		Id          uint64  `json:"id"`
		Platform    string  `json:"platform"`
		Category    string  `json:"category"`
		GrossAmount float64 `json:"gross_amount"`
		PlatformFee float64 `json:"platform_fee"`
		NetAmount   float64 `json:"net_amount"`
		OccurredAt  uint64  `json:"occurred_at"`
		Source      string  `json:"source"`
		Remark      string  `json:"remark"`
	}
	err = m.OrderDesc("occurred_at").OrderDesc("id").
		Page(page, pageSize).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询收入列表失败")
	}

	list := make([]compliancein.IncomeItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, compliancein.IncomeItem{
			Id:          r.Id,
			Platform:    r.Platform,
			Category:    r.Category,
			GrossAmount: r.GrossAmount,
			PlatformFee: r.PlatformFee,
			NetAmount:   r.NetAmount,
			OccurredAt:  formatOccurredAt(r.OccurredAt),
			Source:      r.Source,
			Remark:      r.Remark,
		})
	}

	summary := compliancein.IncomeSummary{}
	sumM := g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).
		Where("deleted", 0)
	if in.Month != "" {
		start, end, _ := monthRange(in.Month)
		sumM = sumM.WhereBetween("occurred_at", start, end)
	}
	var sumRow struct {
		TotalGross float64 `json:"total_gross"`
		TotalFee   float64 `json:"total_fee"`
		TotalNet   float64 `json:"total_net"`
	}
	_ = sumM.Fields(
		"COALESCE(SUM(gross_amount), 0) AS total_gross",
		"COALESCE(SUM(platform_fee), 0) AS total_fee",
		"COALESCE(SUM(net_amount), 0) AS total_net",
	).Scan(&sumRow)
	summary.TotalGross = sumRow.TotalGross
	summary.TotalFee = sumRow.TotalFee
	summary.TotalNet = sumRow.TotalNet

	return &compliancein.IncomeListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Summary:  summary,
	}, nil
}

// CreateIncome 创建收入条目
func (s *sComplianceLedger) CreateIncome(ctx context.Context, in *compliancein.IncomeCreateInp) (*compliancein.IncomeCreateModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	if err = validatePlatform(in.Platform); err != nil {
		return nil, err
	}
	if err = validateIncomeCategory(in.Category); err != nil {
		return nil, err
	}
	if in.GrossAmount <= 0 {
		return nil, gerror.New("含税收入须大于0")
	}
	if in.PlatformFee < 0 {
		return nil, gerror.New("平台服务费不能为负数")
	}
	if in.PlatformFee > in.GrossAmount {
		return nil, gerror.New("平台服务费不能大于含税收入")
	}
	occurredAt, err := parseOccurredAt(in.OccurredAt)
	if err != nil {
		return nil, err
	}

	netAmount := calcNet(in.GrossAmount, in.PlatformFee)
	now := uint64(utility.NowUnix())
	data := g.Map{
		"opc_id":        opcId,
		"platform":      strings.ToLower(strings.TrimSpace(in.Platform)),
		"category":      strings.ToLower(strings.TrimSpace(in.Category)),
		"gross_amount":  roundMoney(in.GrossAmount),
		"platform_fee":  roundMoney(in.PlatformFee),
		"net_amount":    netAmount,
		"occurred_at":   occurredAt,
		"source":        sourceManual,
		"remark":        strings.TrimSpace(in.Remark),
		"deleted":       0,
		"create_time":   now,
		"update_time":   now,
	}

	result, err := g.DB().Model(tableIncomeEntry).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "创建收入失败")
	}
	id, _ := result.LastInsertId()

	_ = audit.WriteAudit(ctx, "income_entry", uint64(id), "income.create", in.MemberId, "member", nil, data, in.Ip)
	_ = RecordIncomeVoucher(ctx, opcId, uint64(id), in.GrossAmount, occurredAt)
	year, month := yearMonthFromUnix(occurredAt)
	_ = refreshProfitSummary(ctx, opcId, year, month)
	_ = tax.RefreshOpcPeriodTaxAmounts(ctx, opcId, year, month)

	return &compliancein.IncomeCreateModel{
		Id:          uint64(id),
		Platform:    data["platform"].(string),
		Category:    data["category"].(string),
		GrossAmount: in.GrossAmount,
		PlatformFee: in.PlatformFee,
		NetAmount:   netAmount,
		OccurredAt:  formatOccurredAt(occurredAt),
		Source:      sourceManual,
	}, nil
}

// ImportIncomeCSV CSV 导入收入
func (s *sComplianceLedger) ImportIncomeCSV(ctx context.Context, in *compliancein.IncomeImportInp) (*compliancein.IncomeImportModel, error) {
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
	col := mapCSVColumns(header, []string{"platform", "category", "gross_amount", "platform_fee", "occurred_at", "remark"})
	if col["platform"] < 0 || col["category"] < 0 || col["gross_amount"] < 0 || col["occurred_at"] < 0 {
		return nil, gerror.New("CSV 表头须包含 platform,category,gross_amount,occurred_at（platform_fee/remark 可选）")
	}

	out := &compliancein.IncomeImportModel{}
	now := uint64(utility.NowUnix())

	for i := 1; i < len(records); i++ {
		row := records[i]
		if isEmptyCSVRow(row) {
			continue
		}
		rowNum := i + 1
		item, rowErr := parseIncomeCSVRow(row, col)
		if rowErr != nil {
			out.Failed++
			out.Errors = append(out.Errors, compliancein.IncomeImportRowError{Row: rowNum, Message: rowErr.Error()})
			continue
		}

		data := g.Map{
			"opc_id":        opcId,
			"platform":      item.platform,
			"category":      item.category,
			"gross_amount":  item.gross,
			"platform_fee":  item.fee,
			"net_amount":    calcNet(item.gross, item.fee),
			"occurred_at":   item.occurredAt,
			"source":        sourceCSV,
			"remark":        item.remark,
			"deleted":       0,
			"create_time":   now,
			"update_time":   now,
		}
		result, insertErr := g.DB().Model(tableIncomeEntry).Ctx(ctx).Data(data).Insert()
		if insertErr != nil {
			out.Failed++
			out.Errors = append(out.Errors, compliancein.IncomeImportRowError{Row: rowNum, Message: "写入失败"})
			continue
		}
		id, _ := result.LastInsertId()
		_ = audit.WriteAudit(ctx, "income_entry", uint64(id), "income.import", in.MemberId, "member", nil, data, in.Ip)
		_ = RecordIncomeVoucher(ctx, opcId, uint64(id), item.gross, item.occurredAt)
		out.Imported++
	}
	return out, nil
}

type incomeCSVItem struct {
	platform   string
	category   string
	gross      float64
	fee        float64
	occurredAt uint64
	remark     string
}

func parseIncomeCSVRow(row []string, col map[string]int) (*incomeCSVItem, error) {
	platform := strings.ToLower(strings.TrimSpace(cell(row, col["platform"])))
	category := strings.ToLower(strings.TrimSpace(cell(row, col["category"])))
	if err := validatePlatform(platform); err != nil {
		return nil, err
	}
	if err := validateIncomeCategory(category); err != nil {
		return nil, err
	}
	gross, err := parseMoney(cell(row, col["gross_amount"]))
	if err != nil || gross <= 0 {
		return nil, gerror.New("gross_amount 无效")
	}
	fee := 0.0
	if col["platform_fee"] >= 0 {
		fee, err = parseMoney(cell(row, col["platform_fee"]))
		if err != nil || fee < 0 {
			return nil, gerror.New("platform_fee 无效")
		}
	}
	if fee > gross {
		return nil, gerror.New("platform_fee 不能大于 gross_amount")
	}
	occurredAt, err := parseOccurredAt(cell(row, col["occurred_at"]))
	if err != nil {
		return nil, err
	}
	remark := ""
	if col["remark"] >= 0 {
		remark = strings.TrimSpace(cell(row, col["remark"]))
	}
	return &incomeCSVItem{
		platform: platform, category: category,
		gross: gross, fee: fee, occurredAt: occurredAt, remark: remark,
	}, nil
}

func normalizeCSVHeader(header []string) []string {
	out := make([]string, len(header))
	for i, h := range header {
		out[i] = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(h, "\ufeff", "")))
	}
	return out
}

func mapCSVColumns(header []string, names []string) map[string]int {
	col := make(map[string]int, len(names))
	for _, n := range names {
		col[n] = -1
	}
	aliases := map[string][]string{
		"platform":     {"platform", "平台"},
		"category":     {"category", "收入类型", "类型"},
		"gross_amount": {"gross_amount", "gross", "含税收入", "金额"},
		"platform_fee": {"platform_fee", "fee", "平台服务费", "服务费"},
		"occurred_at":  {"occurred_at", "date", "发生日期", "日期"},
		"remark":       {"remark", "备注"},
	}
	for i, h := range header {
		for key, alts := range aliases {
			for _, a := range alts {
				if h == a {
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

func parseMoney(s string) (float64, error) {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func isEmptyCSVRow(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}

func monthStartEndFromUnix(ts uint64) (uint64, uint64) {
	t := gtime.NewFromTimeStamp(int64(ts))
	start := gtime.NewFromStr(t.Format("Y-m") + "-01 00:00:00")
	end := start.AddDate(0, 1, 0).Add(-1)
	return uint64(start.Unix()), uint64(end.Unix())
}

func readCSVContent(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
