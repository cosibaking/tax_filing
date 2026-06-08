package diagnosis

import (
	"fmt"
	"math"

	"xygo/internal/model/input/compliancein"
)

const (
	planNone       = "none"
	planLabor      = "labor"
	planIndividual = "individual"
	planOpc        = "opc"
	planTransitional = "transitional"

	vatExemptMonthlySales = 100000.0
	vatRate               = 0.01
	surchargeRate         = 0.06
	citMicroRate          = 0.05
	citMicroProfitCap = 3000000.0
)

type taxBracket struct {
	upper       float64
	rate        float64
	quickOffset float64
}

var laborBrackets = []taxBracket{
	{36000, 0.03, 0},
	{144000, 0.10, 2520},
	{300000, 0.20, 16920},
	{420000, 0.25, 31920},
	{660000, 0.30, 52920},
	{960000, 0.35, 85920},
	{math.MaxFloat64, 0.45, 181920},
}

var businessBrackets = []taxBracket{
	{30000, 0.05, 0},
	{90000, 0.10, 1500},
	{300000, 0.20, 10500},
	{500000, 0.30, 40500},
	{math.MaxFloat64, 0.35, 65500},
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}

func roundRate(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func calcProgressiveTax(taxable float64, brackets []taxBracket) float64 {
	if taxable <= 0 {
		return 0
	}
	for _, b := range brackets {
		if taxable <= b.upper {
			return roundMoney(taxable*b.rate - b.quickOffset)
		}
	}
	return 0
}

// AnnualIncomeFromRange 月收入区间转年收入估算（万元区间取中位）
func AnnualIncomeFromRange(monthlyIncomeRange string) float64 {
	switch monthlyIncomeRange {
	case "0-2万":
		return 120000
	case "2-5万":
		return 420000
	case "5-15万":
		return 1200000
	case "15万+":
		return 2400000
	default:
		return 0
	}
}

// CalcLaborTax 劳务报酬简化年估（超额累进，年入100万约27%）
func CalcLaborTax(annualIncome float64) float64 {
	if annualIncome <= 0 {
		return 0
	}
	// 简化年估：以年收入为应纳税所得额（不含专项扣除），贴近产品文档「年入100万约27%」
	return calcProgressiveTax(annualIncome, laborBrackets)
}

// CalcIndividualTax 个体户查账征收（经营所得5级累进）
func CalcIndividualTax(annualIncome, annualCost float64) float64 {
	taxable := annualIncome - annualCost
	return calcProgressiveTax(taxable, businessBrackets)
}

// CalcOpcTax OPC 小微税负：增值税(1%, 月10万免税) + 附加税(增值税×6%) + 企税(利润×5%)
func CalcOpcTax(annualIncome, annualCost float64) (vat, surcharge, cit, total float64) {
	if annualIncome <= 0 {
		return 0, 0, 0, 0
	}

	exTaxIncome := annualIncome / 1.01
	monthlySales := annualIncome / 12

	if monthlySales < vatExemptMonthlySales {
		vat = 0
	} else {
		vat = exTaxIncome * vatRate
	}
	surcharge = vat * surchargeRate

	profit := exTaxIncome - annualCost
	if profit > 0 && profit <= citMicroProfitCap {
		cit = profit * citMicroRate
	} else if profit > citMicroProfitCap {
		cit = citMicroProfitCap*citMicroRate + (profit-citMicroProfitCap)*0.25
	}

	total = roundMoney(vat + surcharge + cit)
	vat = roundMoney(vat)
	surcharge = roundMoney(surcharge)
	cit = roundMoney(cit)
	return
}

// NoTaxWarning 不报税高危提示
func NoTaxWarning() compliancein.TaxComparisonItem {
	return compliancein.TaxComparisonItem{
		Plan:      planNone,
		Label:     "不报税",
		AnnualTax: 0,
		TaxRate:   0,
		Warning:   true,
		Note:      "潜在补税+罚款风险，强烈建议尽快合规",
	}
}

type planTaxSnapshot struct {
	plan string
	tax  float64
}

// RecommendPlan 综合税负测算与公开政策规则推荐方案（劳务/个体户/OPC 三选一）
func RecommendPlan(annualIncome, annualCost float64, taxBureauContact bool) (plan string, reasons []string) {
	if taxBureauContact {
		return planOpc, []string{
			"您已被税务机关联系，依据税收征管法应依法办理纳税申报，建议优先设立公司主体配合整改",
			"OPC 有限公司便于规范建账、申报与留存合规资料",
		}
	}

	if annualIncome <= 0 {
		return planLabor, []string{"请填写有效的年收入后再进行方案推荐"}
	}

	laborTax := CalcLaborTax(annualIncome)
	individualTax := CalcIndividualTax(annualIncome, annualCost)
	_, _, _, opcTax := CalcOpcTax(annualIncome, annualCost)

	candidates := []planTaxSnapshot{
		{planLabor, laborTax},
		{planIndividual, individualTax},
		{planOpc, opcTax},
	}
	minTax := laborTax
	plan = planLabor
	for _, c := range candidates[1:] {
		if c.tax < minTax {
			minTax = c.tax
			plan = c.plan
		}
	}

	costRatio := annualCost / annualIncome

	// 合规加权：在税负接近时优先更利于平台经济合规经营的主体形式
	switch {
	case annualIncome >= 1000000 && opcTax <= minTax*1.12:
		plan = planOpc
	case annualIncome >= 300000 && opcTax <= minTax*1.08:
		plan = planOpc
	case costRatio >= 0.35 && individualTax <= opcTax*1.05:
		plan = planIndividual
	case annualIncome < 200000 && laborTax < minTax*0.9:
		plan = planLabor
	case annualIncome < 300000 && costRatio < 0.15 && individualTax <= laborTax*1.1:
		plan = planTransitional
	default:
		// 保持税负最低方案
	}

	return plan, buildRecommendReasons(plan, annualIncome, annualCost, costRatio, laborTax, individualTax, opcTax)
}

func buildRecommendReasons(plan string, annualIncome, annualCost, costRatio, laborTax, individualTax, opcTax float64) []string {
	switch plan {
	case planLabor:
		return []string{
			fmt.Sprintf("在现行测算下，劳务报酬预估税负 %.0f 元/年，为三方案中较低", laborTax),
			"劳务报酬适用综合所得税率表（3%~45%超额累进），适合收入规模较小、结算链路简单的场景",
			"若收入持续增长或需平台对公结算，建议重新评估个体户或 OPC 方案",
		}
	case planIndividual, planTransitional:
		reasons := []string{
			fmt.Sprintf("个体户查账征收预估税负 %.0f 元/年，经营所得可扣除合法成本 %.0f 元", individualTax, annualCost),
			"经营所得适用 5%~35% 五级超额累进税率（个人所得税法经营所得）",
		}
		if costRatio >= 0.3 {
			reasons = append(reasons, "您填报的成本占比较高，查账征收下扣除空间更明显")
		}
		if plan == planTransitional {
			reasons = append(reasons, "当前收入规模可先以个体户过渡，待业务稳定后再评估升级 OPC")
		} else {
			reasons = append(reasons, "建议保留完整成本票据，确保税前扣除合规")
		}
		return reasons
	case planOpc:
		reasons := []string{
			fmt.Sprintf("OPC 小微公司预估综合税负 %.0f 元/年，含增值税、附加税及企业所得税", opcTax),
			"小规模纳税人月销售额 10 万元以下免征增值税（延续阶段性政策口径估算）",
			"年应纳税所得额不超过 300 万元部分，按小微企业 5% 优惠税率估算企业所得税",
		}
		if annualIncome >= 300000 {
			reasons = append(reasons, "便于与平台/MCN 对公结算，降低个人账户大额流水合规风险")
		}
		return reasons
	default:
		return []string{"请结合业务场景选择合法经营主体并按时申报"}
	}
}

// BuildAssumptionHints 生成税负计算器想定提示
func BuildAssumptionHints(annualIncome, annualCost float64) []string {
	hints := []string{
		"本结果为示意性测算，不构成税务鉴定或申报依据；实际纳税以主管税务机关核定为准。",
		fmt.Sprintf("测算基数：年收入 %.0f 元，年可扣除成本 %.0f 元。", annualIncome, annualCost),
	}

	monthlySales := annualIncome / 12
	if monthlySales < vatExemptMonthlySales {
		hints = append(hints, fmt.Sprintf(
			"OPC 增值税：折算月销售额约 %.1f 万元，低于 10 万元/月免征门槛，按免征估算。",
			monthlySales/10000,
		))
	} else {
		hints = append(hints, "OPC 增值税：月销售额超过 10 万元，按小规模 1% 征收率估算（含税收入÷1.01×1%）。")
	}

	hints = append(hints,
		"劳务报酬：按综合所得税率表超额累进估算，未模拟专项附加扣除及预扣预缴差异。",
		"个体户：按经营所得五级超额累进税率估算，成本须为与经营相关的合法票据。",
	)

	profit := annualIncome/1.01 - annualCost
	if profit > 0 && profit <= citMicroProfitCap {
		hints = append(hints, "企业所得税：按小微企业 5% 优惠税率估算（利润≤300 万元区间）。")
	} else if profit > citMicroProfitCap {
		hints = append(hints, "企业所得税：利润超过 300 万元部分按 25% 税率分段估算。")
	}

	return hints
}

// CompareTaxSchemes 三方案税负对比（含不报税高危列）
func CompareTaxSchemes(annualIncome, annualCost float64, taxBureauContact bool) (comparison compliancein.TaxComparison, recommendedPlan string, reasons []string) {
	comparison = compliancein.TaxComparison{
		AnnualIncome: roundMoney(annualIncome),
		AnnualCost:   roundMoney(annualCost),
		Items:        make([]compliancein.TaxComparisonItem, 0, 4),
	}

	recommendedPlan, reasons = RecommendPlan(annualIncome, annualCost, taxBureauContact)

	comparison.Items = append(comparison.Items, NoTaxWarning())

	laborTax := CalcLaborTax(annualIncome)
	comparison.Items = append(comparison.Items, compliancein.TaxComparisonItem{
		Plan:        planLabor,
		Label:       "纯劳务",
		AnnualTax:   laborTax,
		TaxRate:     roundRate(safeRate(laborTax, annualIncome)),
		Recommended: recommendedPlan == planLabor,
	})

	individualTax := CalcIndividualTax(annualIncome, annualCost)
	comparison.Items = append(comparison.Items, compliancein.TaxComparisonItem{
		Plan:        planIndividual,
		Label:       "个体户",
		AnnualTax:   individualTax,
		TaxRate:     roundRate(safeRate(individualTax, annualIncome)),
		Recommended: recommendedPlan == planIndividual || recommendedPlan == planTransitional,
	})

	_, _, _, opcTotal := CalcOpcTax(annualIncome, annualCost)
	comparison.Items = append(comparison.Items, compliancein.TaxComparisonItem{
		Plan:        planOpc,
		Label:       "OPC",
		AnnualTax:   opcTotal,
		TaxRate:     roundRate(safeRate(opcTotal, annualIncome)),
		Recommended: recommendedPlan == planOpc,
		Note:        "含增值税、附加税及小微企税（不含分红模拟）",
	})

	for i := range comparison.Items {
		if comparison.Items[i].Plan == recommendedPlan ||
			(recommendedPlan == planTransitional && comparison.Items[i].Plan == planIndividual) {
			comparison.Items[i].Recommended = true
		}
	}

	return comparison, recommendedPlan, reasons
}

func safeRate(tax, income float64) float64 {
	if income <= 0 {
		return 0
	}
	return tax / income
}
