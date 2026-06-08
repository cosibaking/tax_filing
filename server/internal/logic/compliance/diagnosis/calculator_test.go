package diagnosis

import (
	"math"
	"strings"
	"testing"
)

func TestCalcLaborTax_MillionIncome(t *testing.T) {
	tax := CalcLaborTax(1000000)
	rate := tax / 1000000
	if rate < 0.25 || rate > 0.29 {
		t.Fatalf("CalcLaborTax(1M) rate=%.4f, want ~0.27", rate)
	}
}

func TestCalcOpcTax_VatExemptUnder100kMonthly(t *testing.T) {
	annualIncome := 1000000.0 // 月销约8.33万，低于10万免税
	annualCost := 200000.0

	vat, surcharge, cit, total := CalcOpcTax(annualIncome, annualCost)
	if vat != 0 {
		t.Fatalf("vat=%v, want 0 for monthly sales < 100k", vat)
	}
	if surcharge != 0 {
		t.Fatalf("surcharge=%v, want 0 when vat is 0", surcharge)
	}

	exTaxIncome := annualIncome / 1.01
	expectedCit := (exTaxIncome - annualCost) * citMicroRate
	if math.Abs(cit-expectedCit) > 1 {
		t.Fatalf("cit=%v, want ~%v", cit, expectedCit)
	}
	if math.Abs(total-cit) > 0.01 {
		t.Fatalf("total=%v, want cit only %v", total, cit)
	}
}

func TestCalcOpcTax_VatAbove100kMonthly(t *testing.T) {
	annualIncome := 1500000.0 // 月销12.5万
	annualCost := 300000.0

	vat, surcharge, cit, total := CalcOpcTax(annualIncome, annualCost)
	exTaxIncome := annualIncome / 1.01
	expectedVat := exTaxIncome * vatRate
	if math.Abs(vat-expectedVat) > 1 {
		t.Fatalf("vat=%v, want ~%v", vat, expectedVat)
	}
	expectedSurcharge := expectedVat * surchargeRate
	if math.Abs(surcharge-expectedSurcharge) > 0.5 {
		t.Fatalf("surcharge=%v, want ~%v", surcharge, expectedSurcharge)
	}
	expectedCit := (exTaxIncome - annualCost) * citMicroRate
	if math.Abs(cit-expectedCit) > 1 {
		t.Fatalf("cit=%v, want ~%v", cit, expectedCit)
	}
	if math.Abs(total-(vat+surcharge+cit)) > 0.5 {
		t.Fatalf("total=%v, want sum %v", total, vat+surcharge+cit)
	}
}

func TestCompareTaxSchemes_IncludesNoTaxWarning(t *testing.T) {
	comparison, plan, _ := CompareTaxSchemes(500000, 100000, false)
	if plan == "" {
		t.Fatal("recommended plan should not be empty")
	}
	if len(comparison.Items) != 4 {
		t.Fatalf("items len=%d, want 4", len(comparison.Items))
	}
	if !comparison.Items[0].Warning || comparison.Items[0].Plan != planNone {
		t.Fatalf("first item should be no_tax warning, got %+v", comparison.Items[0])
	}
}

func TestRecommendPlan_TaxBureauContact(t *testing.T) {
	plan, reasons := RecommendPlan(200000, 50000, true)
	if plan != planOpc {
		t.Fatalf("plan=%s, want opc when tax bureau contacted", plan)
	}
	if len(reasons) == 0 {
		t.Fatal("reasons should not be empty")
	}
}

func TestRecommendPlan_HighIncome(t *testing.T) {
	plan, _ := RecommendPlan(1500000, 200000, false)
	if plan != planOpc {
		t.Fatalf("plan=%s, want opc for income >= 1M", plan)
	}
}

func TestRecommendPlan_AlwaysReturnsValidPlan(t *testing.T) {
	valid := map[string]bool{
		planLabor: true, planIndividual: true, planOpc: true, planTransitional: true,
	}
	cases := []struct{ income, cost float64 }{
		{100000, 0},
		{300000, 50000},
		{250000, 100000},
		{1500000, 200000},
	}
	for _, c := range cases {
		plan, reasons := RecommendPlan(c.income, c.cost, false)
		if !valid[plan] {
			t.Fatalf("invalid plan %s for income=%v cost=%v", plan, c.income, c.cost)
		}
		if len(reasons) == 0 {
			t.Fatalf("empty reasons for income=%v cost=%v", c.income, c.cost)
		}
	}
}

func TestBuildAssumptionHints_VatExempt(t *testing.T) {
	hints := BuildAssumptionHints(300000, 0)
	if len(hints) < 4 {
		t.Fatalf("hints len=%d, want at least 4", len(hints))
	}
	found := false
	for _, h := range hints {
		if strings.Contains(h, "10 万元") || strings.Contains(h, "10万元") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("hints should mention VAT exemption threshold, got %v", hints)
	}
}


func TestCalcIndividualTax_Progressive(t *testing.T) {
	tax := CalcIndividualTax(500000, 100000)
	if tax <= 0 {
		t.Fatal("individual tax should be positive")
	}
	// 应税所得40万，适用30%档
	if tax < 70000 || tax > 120000 {
		t.Fatalf("tax=%v, want reasonable progressive result", tax)
	}
}

func TestAnnualIncomeFromRange(t *testing.T) {
	cases := map[string]float64{
		"0-2万":  120000,
		"2-5万":  420000,
		"5-15万": 1200000,
		"15万+":  2400000,
	}
	for input, want := range cases {
		if got := AnnualIncomeFromRange(input); got != want {
			t.Fatalf("AnnualIncomeFromRange(%q)=%v, want %v", input, got, want)
		}
	}
}
