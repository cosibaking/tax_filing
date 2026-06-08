package tax

import (
	"math"
	"testing"
)

func TestExTaxAmount(t *testing.T) {
	got := ExTaxAmount(101000)
	want := 100000.0
	if math.Abs(got-want) > 0.01 {
		t.Fatalf("ExTaxAmount(101000)=%v, want %v", got, want)
	}
	if ExTaxAmount(0) != 0 {
		t.Fatal("ExTaxAmount(0) should be 0")
	}
}

func TestCalcMonthlyVAT_ExemptUnder100k(t *testing.T) {
	// 月含税 10万 → 不含税约 9.9万，低于 10万免税线
	vat := CalcMonthlyVAT(100000)
	if vat != 0 {
		t.Fatalf("vat=%v, want 0 for exempt monthly sales", vat)
	}
}

func TestCalcMonthlyVAT_AtThreshold(t *testing.T) {
	// 月含税刚好使不含税 = 10万 → 不免税
	gross := VatExemptMonthlySales * VatDivisor
	vat := CalcMonthlyVAT(gross)
	expected := roundMoney(VatExemptMonthlySales * VatRate)
	if math.Abs(vat-expected) > 0.01 {
		t.Fatalf("vat=%v, want %v at threshold", vat, expected)
	}
}

func TestCalcMonthlyVAT_AboveThreshold(t *testing.T) {
	monthlyGross := 1200000.0
	vat := CalcMonthlyVAT(monthlyGross)
	exTax := ExTaxAmount(monthlyGross)
	expected := roundMoney(exTax * VatRate)
	if math.Abs(vat-expected) > 0.01 {
		t.Fatalf("vat=%v, want %v", vat, expected)
	}
}

func TestCalcSurcharge(t *testing.T) {
	if CalcSurcharge(0) != 0 {
		t.Fatal("surcharge should be 0 when vat is 0")
	}
	vat := 1000.0
	surcharge := CalcSurcharge(vat)
	if math.Abs(surcharge-60) > 0.01 {
		t.Fatalf("surcharge=%v, want 60", surcharge)
	}
}

func TestCalcQuarterlyCIT(t *testing.T) {
	if CalcQuarterlyCIT(0) != 0 {
		t.Fatal("CIT should be 0 for zero profit")
	}
	if CalcQuarterlyCIT(-1000) != 0 {
		t.Fatal("CIT should be 0 for negative profit")
	}
	cit := CalcQuarterlyCIT(200000)
	if math.Abs(cit-10000) > 0.01 {
		t.Fatalf("cit=%v, want 10000", cit)
	}
}

func TestCalcAnnualCIT(t *testing.T) {
	cit := CalcAnnualCIT(500000)
	if math.Abs(cit-25000) > 0.01 {
		t.Fatalf("annual cit=%v, want 25000", cit)
	}
}
