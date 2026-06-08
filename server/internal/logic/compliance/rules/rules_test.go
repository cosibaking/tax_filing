package rules

import "testing"

func TestCheckLowReportedIncome_blocks(t *testing.T) {
	err := CheckLowReportedIncome(10000, 8000)
	if err == nil {
		t.Fatal("expected block when declared < ledger")
	}
}

func TestCheckLowReportedIncome_allows_equal(t *testing.T) {
	if err := CheckLowReportedIncome(10000, 10000); err != nil {
		t.Fatalf("equal amounts should pass: %v", err)
	}
}

func TestCheckProhibitedIncomeCategory(t *testing.T) {
	if err := CheckProhibitedIncomeCategory("loan"); err == nil {
		t.Fatal("loan category should be blocked")
	}
	if err := CheckProhibitedIncomeCategory("tip"); err != nil {
		t.Fatalf("tip should be allowed: %v", err)
	}
}

func TestCheckIncomeRemark(t *testing.T) {
	if err := CheckIncomeRemark("朋友借款还款"); err == nil {
		t.Fatal("remark with 借款 should be blocked")
	}
}

func TestPeriodRange_month(t *testing.T) {
	start, end, err := PeriodRange("2026-06")
	if err != nil {
		t.Fatalf("periodRange() error = %v", err)
	}
	if start == 0 || end <= start {
		t.Fatalf("invalid range: %d - %d", start, end)
	}
}
