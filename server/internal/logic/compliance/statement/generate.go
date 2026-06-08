package statement

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/notice"
	"xygo/internal/logic/compliance/rules"
	"xygo/internal/logic/compliance/shared"
	"xygo/internal/logic/compliance/tax"
	"xygo/utility"
)

// GenerateForPeriod 为所有 active OPC 生成指定年月对账单（F-61）
func GenerateForPeriod(ctx context.Context, year, month int) (created, skipped int, err error) {
	period := fmt.Sprintf("%04d-%02d", year, month)

	var opcs []struct {
		Id       uint64 `json:"id"`
		MemberId uint64 `json:"member_id"`
	}
	err = g.DB().Model(shared.TableOpcEntity).Ctx(ctx).
		Where("status", "active").
		Where("deleted", 0).
		Fields("id, member_id").
		Scan(&opcs)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "查询OPC主体失败")
	}

	now := uint64(utility.NowUnix())
	for _, opc := range opcs {
		exists, _ := g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).
			Where("opc_id", opc.Id).
			Where("year", year).
			Where("month", month).
			Where("deleted", 0).
			Count()
		if exists > 0 {
			skipped++
			continue
		}

		summary, sumErr := buildSummary(ctx, opc.Id, year, month, period)
		if sumErr != nil {
			return created, skipped, sumErr
		}
		summaryJSON, _ := gjson.Encode(summary)

		_, insErr := g.DB().Model(shared.TableMonthlyStmt).Ctx(ctx).Data(g.Map{
			"opc_id":      opc.Id,
			"year":        year,
			"month":       month,
			"summary":     string(summaryJSON),
			"create_time": now,
			"update_time": now,
		}).Insert()
		if insErr != nil {
			return created, skipped, gerror.Wrap(insErr, "写入对账单失败")
		}
		created++

		_ = notice.SendToMember(ctx, opc.MemberId,
			fmt.Sprintf("%s 月度对账单已生成", period),
			fmt.Sprintf("您的 %s 月度对账单已生成，收入 %.2f 元、利润 %.2f 元，请登录会员中心查看。",
				period, summary.Revenue, summary.Profit))
	}
	return created, skipped, nil
}

type statementSummary struct {
	Revenue           float64 `json:"revenue"`
	Cost              float64 `json:"cost"`
	Profit            float64 `json:"profit"`
	TaxAmount         float64 `json:"taxAmount"`
	CumulativeProfit  float64 `json:"cumulativeProfit"`
	FilingStatus      string  `json:"filingStatus"`
}

func buildSummary(ctx context.Context, opcId uint64, year, month int, period string) (*statementSummary, error) {
	var profitRow struct {
		Revenue          float64 `json:"revenue"`
		Cost             float64 `json:"cost"`
		Profit           float64 `json:"profit"`
		CumulativeProfit float64 `json:"cumulative_profit"`
	}
	_ = g.DB().Model(shared.TableProfitSummary).Ctx(ctx).
		Where("opc_id", opcId).
		Where("year", year).
		Where("month", month).
		Scan(&profitRow)

	if profitRow.Revenue == 0 && profitRow.Cost == 0 {
		income, _ := rules.SumLedgerIncome(ctx, opcId, period)
		profitRow.Revenue = income
		profitRow.Profit = income
	}

	vat := tax.CalcMonthlyVAT(profitRow.Revenue)
	surcharge := tax.CalcSurcharge(vat)
	taxAmount := vat + surcharge

	filingStatus := "pending"
	filedCount, _ := g.DB().Model(shared.TableTaxFilingTask).Ctx(ctx).
		Where("opc_id", opcId).
		Where("period", period).
		Where("status", "filed").
		Where("deleted", 0).
		Count()
	if filedCount > 0 {
		filingStatus = "filed"
	}

	return &statementSummary{
		Revenue:          profitRow.Revenue,
		Cost:             profitRow.Cost,
		Profit:           profitRow.Profit,
		TaxAmount:        taxAmount,
		CumulativeProfit: profitRow.CumulativeProfit,
		FilingStatus:     filingStatus,
	}, nil
}

// PreviousMonth 返回上一个月的年月
func PreviousMonth(now time.Time) (year, month int) {
	prev := now.AddDate(0, -1, 0)
	return prev.Year(), int(prev.Month())
}
