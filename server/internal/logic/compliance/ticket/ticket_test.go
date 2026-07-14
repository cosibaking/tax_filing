package ticket

import "testing"

func TestTransition(t *testing.T) {
	if got, err := Transition("submitted", "accept"); err != nil || got != "accepted" {
		t.Fatalf("got=%s err=%v", got, err)
	}
	if _, err := Transition("completed", "accept"); err == nil {
		t.Fatal("completed ticket must be immutable")
	}
}
func TestLicensedServiceRequiresProvider(t *testing.T) {
	if err := ValidateProvider("tax_filing", ""); err == nil {
		t.Fatal("provider required")
	}
	if err := ValidateProvider("compliance_review", ""); err != nil {
		t.Fatal(err)
	}
}
