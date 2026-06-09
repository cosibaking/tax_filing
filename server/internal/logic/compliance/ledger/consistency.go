package ledger

import (
	"context"
	"math"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/model/input/compliancein"
)

const consistencyTolerance = 0.01 // 1%

// GetIncomeConsistency 收入一致性比对（台账 vs 银行 vs 平台报送）
func (s *sComplianceLedger) GetIncomeConsistency(ctx context.Context, in *compliancein.IncomeConsistencyInp) (*compliancein.IncomeConsistencyModel, error) {
	opcId, err := s.requireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	month := in.Month
	if month == "" {
		month = gtime.Now().Format("Y-m")
	}
	start, end, err := monthRange(month)
	if err != nil {
		return nil, err
	}

	out := &compliancein.IncomeConsistencyModel{Month: month}

	// 台账汇总
	var ledgerRow struct {
		Gross float64 `json:"gross"`
		Net   float64 `json:"net"`
	}
	_ = g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).Where("deleted", 0).
		WhereBetween("occurred_at", start, end).
		Fields("COALESCE(SUM(gross_amount),0) AS gross", "COALESCE(SUM(net_amount),0) AS net").
		Scan(&ledgerRow)
	out.LedgerGross = ledgerRow.Gross
	out.LedgerNet = ledgerRow.Net

	// 银行入账
	var bankRows []struct {
		Amount  float64 `json:"amount"`
		Matched int     `json:"matched"`
	}
	_ = g.DB().Model(tableBankTransaction).Ctx(ctx).
		Where("opc_id", opcId).Where("deleted", 0).
		WhereBetween("occurred_at", start, end).
		Where("amount > 0").
		Fields("amount", "matched").
		Scan(&bankRows)
	for _, b := range bankRows {
		out.BankInflow += b.Amount
		if b.Matched == 1 {
			out.BankMatched += b.Amount
		} else {
			out.BankUnmatched += b.Amount
		}
	}
	out.BankInflow = roundMoney(out.BankInflow)
	out.BankMatched = roundMoney(out.BankMatched)
	out.BankUnmatched = roundMoney(out.BankUnmatched)

	// 平台报送快照（OCR/手工）
	var snapRows []struct {
		Platform string  `json:"platform"`
		Net      float64 `json:"net"`
	}
	_ = g.DB().Model(tablePlatformIncomeSnapshot).Ctx(ctx).
		Where("opc_id", opcId).
		Where("period", month).
		Fields("platform", "net_amount AS net").
		Scan(&snapRows)
	platformMap := map[string]float64{}
	for _, snap := range snapRows {
		platformMap[snap.Platform] += snap.Net
		out.PlatformReported += snap.Net
	}
	out.PlatformReported = roundMoney(out.PlatformReported)

	// 分平台台账
	var platLedger []struct {
		Platform string  `json:"platform"`
		Net      float64 `json:"net"`
	}
	_ = g.DB().Model(tableIncomeEntry).Ctx(ctx).
		Where("opc_id", opcId).Where("deleted", 0).
		WhereBetween("occurred_at", start, end).
		Fields("platform", "COALESCE(SUM(net_amount),0) AS net").
		Group("platform").
		Scan(&platLedger)

	platformsSeen := map[string]bool{}
	for _, pl := range platLedger {
		platformsSeen[pl.Platform] = true
		reported := platformMap[pl.Platform]
		item := compliancein.IncomeConsistencyPlatformItem{
			Platform:         pl.Platform,
			LedgerNet:        roundMoney(pl.Net),
			PlatformReported: roundMoney(reported),
			Variance:         roundMoney(pl.Net - reported),
		}
		item.Status = consistencyStatus(pl.Net, reported)
		out.Platforms = append(out.Platforms, item)
	}
	for p, reported := range platformMap {
		if platformsSeen[p] {
			continue
		}
		out.Platforms = append(out.Platforms, compliancein.IncomeConsistencyPlatformItem{
			Platform:         p,
			PlatformReported: roundMoney(reported),
			Variance:         roundMoney(-reported),
			Status:           "danger",
		})
	}

	out.VarianceLedgerBank = roundMoney(out.LedgerNet - out.BankMatched)
	out.VarianceLedgerPlatform = roundMoney(out.LedgerNet - out.PlatformReported)
	base := math.Max(out.LedgerNet, math.Max(out.BankMatched, out.PlatformReported))
	if base > 0 {
		out.VarianceRate = roundMoney(math.Abs(out.VarianceLedgerBank) / base)
		if out.PlatformReported > 0 {
			pr := math.Abs(out.VarianceLedgerPlatform) / base
			if pr > out.VarianceRate {
				out.VarianceRate = roundMoney(pr)
			}
		}
	}

	out.Status = overallConsistencyStatus(out)
	out.IncomeMatchPassed = out.Status == "ok"
	out.Hints = buildConsistencyHints(out)
	return out, nil
}

func consistencyStatus(ledger, reported float64) string {
	if ledger == 0 && reported == 0 {
		return "ok"
	}
	base := math.Max(ledger, reported)
	if base == 0 {
		return "warning"
	}
	rate := math.Abs(ledger-reported) / base
	if rate <= consistencyTolerance {
		return "ok"
	}
	if rate <= 0.05 {
		return "warning"
	}
	return "danger"
}

func overallConsistencyStatus(out *compliancein.IncomeConsistencyModel) string {
	if out.LedgerNet == 0 && out.BankMatched == 0 && out.PlatformReported == 0 {
		return "warning"
	}
	if out.VarianceRate <= consistencyTolerance {
		return "ok"
	}
	if out.VarianceRate <= 0.05 {
		return "warning"
	}
	return "danger"
}

func buildConsistencyHints(out *compliancein.IncomeConsistencyModel) []string {
	var hints []string
	if out.LedgerNet == 0 {
		hints = append(hints, "本月尚未录入收入台账")
	}
	if out.BankInflow == 0 {
		hints = append(hints, "本月尚无银行流水，建议导入对公/个人卡流水")
	}
	if out.PlatformReported == 0 {
		hints = append(hints, "本月无平台报送快照，可通过 OCR 导入创作者中心截图")
	}
	if out.BankUnmatched > 0 {
		hints = append(hints, "存在未匹配银行入账，请核对是否遗漏收入记录")
	}
	if out.VarianceLedgerBank != 0 && out.Status != "ok" {
		hints = append(hints, "台账实收与银行已匹配金额存在偏差，请核对")
	}
	if out.VarianceLedgerPlatform != 0 && out.PlatformReported > 0 && out.Status != "ok" {
		hints = append(hints, "台账与平台 OCR 报送数据存在偏差，请核对")
	}
	if out.Status == "ok" && len(hints) == 0 {
		hints = append(hints, "台账、银行与平台数据偏差在 1% 以内")
	}
	return hints
}
