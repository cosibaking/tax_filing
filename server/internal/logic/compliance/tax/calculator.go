package tax

import "math"

const (
	VatDivisor            = 1.01
	VatExemptMonthlySales = 100000.0
	VatRate               = 0.01
	SurchargeRate         = 0.06
	CitMicroRate          = 0.05
)

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}

// ExTaxAmount 含税转不含税
func ExTaxAmount(grossAmount float64) float64 {
	if grossAmount <= 0 {
		return 0
	}
	return roundMoney(grossAmount / VatDivisor)
}

// CalcMonthlyVAT 月增值税：月不含税 < 10万免税，否则不含税 × 1%
func CalcMonthlyVAT(monthlyGross float64) float64 {
	if monthlyGross <= 0 {
		return 0
	}
	exTax := ExTaxAmount(monthlyGross)
	if exTax < VatExemptMonthlySales {
		return 0
	}
	return roundMoney(exTax * VatRate)
}

// CalcSurcharge 附加税 = 增值税 × 6%
func CalcSurcharge(vat float64) float64 {
	if vat <= 0 {
		return 0
	}
	return roundMoney(vat * SurchargeRate)
}

// CalcQuarterlyCIT 企税季度预缴 = max(0, 季度累计利润) × 5%
func CalcQuarterlyCIT(quarterlyProfit float64) float64 {
	if quarterlyProfit <= 0 {
		return 0
	}
	return roundMoney(quarterlyProfit * CitMicroRate)
}

// CalcAnnualCIT 企税年度汇算利润税负（小微5%）
func CalcAnnualCIT(annualProfit float64) float64 {
	return CalcQuarterlyCIT(annualProfit)
}
