package ledger

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/library/platformocr"
	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/tax"
	"xygo/internal/model/input/compliancein"
	"xygo/utility"
)

const (
	sourceOCR                  = "ocr"
	tablePlatformIncomeSnapshot = "xy_platform_income_snapshot"
)

// PreviewIncomeOCR 平台流水 OCR 预览
func (s *sComplianceLedger) PreviewIncomeOCR(ctx context.Context, in *compliancein.IncomeOCRPreviewInp) (*compliancein.IncomeOCRPreviewModel, error) {
	if _, err := s.requireActiveOpc(ctx, in.MemberId); err != nil {
		return nil, err
	}
	text, ocrSource, err := platformocr.ExtractText(ctx, in.ImageBytes, in.Filename, in.OcrText)
	if err != nil {
		return nil, err
	}
	platform := platformocr.GuessPlatformFromText(text, in.Platform)
	if err = validatePlatform(platform); err != nil {
		return nil, err
	}

	parsed := platformocr.ParseStatementText(text, platform)
	out := &compliancein.IncomeOCRPreviewModel{
		Platform:       platform,
		OcrSource:      ocrSource,
		RawTextPreview: truncateOCRText(text, 500),
		Rows:           make([]compliancein.IncomeImportPreviewRow, 0, len(parsed)),
	}
	for i, row := range parsed {
		preview := compliancein.IncomeImportPreviewRow{
			Row:         i + 1,
			OccurredAt:  row.OccurredAt,
			Platform:    platform,
			Category:    row.Category,
			GrossAmount: row.GrossAmount,
			PlatformFee: row.PlatformFee,
		}
		if err = validateIncomeCategory(row.Category); err != nil {
			preview.Valid = false
			preview.Error = err.Error()
			out.InvalidCount++
		} else if _, err = parseOccurredAt(row.OccurredAt); err != nil {
			preview.Valid = false
			preview.Error = err.Error()
			out.InvalidCount++
		} else if row.GrossAmount <= 0 {
			preview.Valid = false
			preview.Error = "金额无效"
			out.InvalidCount++
		} else {
			preview.Valid = true
			out.ValidCount++
			out.TotalGross += row.GrossAmount
			out.TotalNet += calcNet(row.GrossAmount, row.PlatformFee)
		}
		out.Rows = append(out.Rows, preview)
	}
	if out.ValidCount == 0 {
		return nil, gerror.New("未能从截图/文本中解析出有效收入行，请检查内容或手动修正后重试")
	}
	return out, nil
}

// ImportIncomeOCR OCR 确认导入
func (s *sComplianceLedger) ImportIncomeOCR(ctx context.Context, in *compliancein.IncomeOCRImportInp) (*compliancein.IncomeImportModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	preview, err := s.PreviewIncomeOCR(ctx, &compliancein.IncomeOCRPreviewInp{
		MemberId:     in.MemberId,
		Platform:     in.Platform,
		OcrText:      in.OcrText,
		ImageBytes:   in.ImageBytes,
		Filename:     in.Filename,
		AttachmentId: in.AttachmentId,
	})
	if err != nil {
		return nil, err
	}
	if preview.InvalidCount > 0 {
		return nil, gerror.New("存在无效行，请修正 OCR 文本后重试")
	}

	out := &compliancein.IncomeImportModel{}
	now := uint64(utility.NowUnix())
	platform := preview.Platform
	type periodSum struct {
		gross, fee, net float64
	}
	periodTotals := map[string]*periodSum{}

	for _, row := range preview.Rows {
		if !row.Valid {
			continue
		}
		occurredAt, _ := parseOccurredAt(row.OccurredAt)
		net := calcNet(row.GrossAmount, row.PlatformFee)
		data := g.Map{
			"opc_id":        opcId,
			"platform":      platform,
			"category":      row.Category,
			"gross_amount":  roundMoney(row.GrossAmount),
			"platform_fee":  roundMoney(row.PlatformFee),
			"net_amount":    net,
			"occurred_at":   occurredAt,
			"source":        sourceOCR,
			"remark":        "平台流水 OCR 导入",
			"attachment_id": in.AttachmentId,
			"deleted":       0,
			"create_time":   now,
			"update_time":   now,
		}
		result, insertErr := g.DB().Model(tableIncomeEntry).Ctx(ctx).Data(data).Insert()
		if insertErr != nil {
			out.Failed++
			out.Errors = append(out.Errors, compliancein.IncomeImportRowError{Row: row.Row, Message: "写入失败"})
			continue
		}
		id, _ := result.LastInsertId()
		_ = audit.WriteAudit(ctx, "income_entry", uint64(id), "income.ocr_import", in.MemberId, "member", nil, data, in.Ip)
		_ = RecordIncomeVoucher(ctx, opcId, uint64(id), row.GrossAmount, occurredAt)
		year, month := yearMonthFromUnix(occurredAt)
		_ = refreshProfitSummary(ctx, opcId, year, month)
		_ = tax.RefreshOpcPeriodTaxAmounts(ctx, opcId, year, month)
		period := formatPeriod(year, month)
		if periodTotals[period] == nil {
			periodTotals[period] = &periodSum{}
		}
		periodTotals[period].gross += row.GrossAmount
		periodTotals[period].fee += row.PlatformFee
		periodTotals[period].net += net
		out.Imported++
	}

	// 写入平台报送快照（供一致性比对）
	for period, sum := range periodTotals {
		_ = upsertPlatformSnapshot(ctx, opcId, platform, period, sum.gross, sum.fee, sum.net, in.AttachmentId, now)
	}

	return out, nil
}

func upsertPlatformSnapshot(ctx context.Context, opcId uint64, platform, period string, gross, fee, net float64, attachmentId, now uint64) error {
	var existing struct{ Id uint64 `json:"id"` }
	_ = g.DB().Model(tablePlatformIncomeSnapshot).Ctx(ctx).
		Where("opc_id", opcId).
		Where("platform", platform).
		Where("period", period).
		Where("source", sourceOCR).
		Scan(&existing)
	data := g.Map{
		"gross_amount":  roundMoney(gross),
		"platform_fee":  roundMoney(fee),
		"net_amount":    roundMoney(net),
		"attachment_id": attachmentId,
		"update_time":   now,
	}
	if existing.Id > 0 {
		_, err := g.DB().Model(tablePlatformIncomeSnapshot).Ctx(ctx).Where("id", existing.Id).Data(data).Update()
		return err
	}
	data["opc_id"] = opcId
	data["platform"] = platform
	data["period"] = period
	data["source"] = sourceOCR
	data["remark"] = "OCR 导入汇总"
	data["create_time"] = now
	_, err := g.DB().Model(tablePlatformIncomeSnapshot).Ctx(ctx).Data(data).Insert()
	return err
}

func truncateOCRText(s string, n int) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

func formatPeriod(year, month int) string {
	return strconv.Itoa(year) + "-" + pad2(month)
}

func pad2(v int) string {
	if v < 10 {
		return "0" + strconv.Itoa(v)
	}
	return strconv.Itoa(v)
}
