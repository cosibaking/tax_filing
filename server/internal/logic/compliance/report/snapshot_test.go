package report

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

type snapshotLoaderStub struct {
	wantOpc, wantMember uint64
	wantPeriod          string
	result              Input
	calls               int
}

func (s *snapshotLoaderStub) Load(_ context.Context, opcID, memberID uint64, period string) (Input, error) {
	s.calls++
	if opcID != s.wantOpc || memberID != s.wantMember || period != s.wantPeriod {
		panic("snapshot identity or period scope changed")
	}
	return s.result, nil
}

type createRepositoryStub struct {
	saved  Input
	result Result
}

func (r *createRepositoryStub) SaveDraft(_ context.Context, opcID, memberID uint64, in Input, result Result) (*MonthlyReport, error) {
	r.saved, r.result = in, result
	return &MonthlyReport{ID: 1, OpcEntityID: opcID, MemberID: memberID, PeriodKey: in.PeriodKey, StructuredReport: &result.Structured}, nil
}
func (*createRepositoryStub) Get(context.Context, uint64, uint64) (*MonthlyReport, error) {
	return nil, nil
}
func (*createRepositoryStub) List(context.Context, uint64, string) ([]MonthlyReport, error) {
	return nil, nil
}
func (*createRepositoryStub) Publish(context.Context, uint64, uint64) (*MonthlyReport, error) {
	return nil, nil
}

func TestCreateIgnoresMaliciousClientSnapshotAndUsesScopedServerSnapshot(t *testing.T) {
	trusted := Input{PeriodKey: "2026-07", Statistics: map[string]any{"trusted": true}, Completeness: map[string]any{"rate": 100}, Risks: []map[string]any{{
		"code": "SERVER_RISK", "categoryCode": "tax", "severity": "medium", "title": "服务端风险",
		"facts": "服务端事实", "basis": "服务端依据", "impact": "服务端影响", "recommendation": "服务端建议", "ruleVersion": "server-v1",
	}}}
	loader := &snapshotLoaderStub{wantOpc: 81, wantMember: 29, wantPeriod: "2026-07", result: trusted}
	repo := &createRepositoryStub{}
	svc := newServiceWithSnapshotLoader(repo, nil, loader)

	malicious := Input{PeriodKey: "2026-07", Statistics: map[string]any{"revenue": 999999}, Completeness: map[string]any{"rate": 0}, Risks: []map[string]any{{
		"code": "CLIENT_FORGED", "severity": "high", "facts": "伪造事实", "basis": "伪造依据",
	}}}
	item, err := svc.Create(context.Background(), 81, 29, malicious)
	if err != nil {
		t.Fatal(err)
	}
	if loader.calls != 1 || !reflect.DeepEqual(repo.saved, trusted) {
		t.Fatalf("loader calls=%d saved=%+v", loader.calls, repo.saved)
	}
	encoded := repo.result.Content + " " + repo.result.Structured.Anomalies[0].Facts + " " + repo.result.Structured.Anomalies[0].Basis
	if strings.Contains(encoded, "伪造") || repo.result.Structured.Anomalies[0].Code != "SERVER_RISK" || item.MemberID != 29 {
		t.Fatalf("client data reached report: %+v", repo.result.Structured)
	}
}

func TestTrustedRiskProducesDetailedServerOwnedAnomaly(t *testing.T) {
	raw := trustedRisk(riskSnapshotRow{RiskCode: "ZERO_WITH_CASHFLOW", Severity: "high", Summary: "零申报期间存在经营流水", RuleVersionID: 42}, map[string]any{"revenue": 12000})
	result, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 80}, Risks: []map[string]any{raw}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Anomalies) != 1 {
		t.Fatalf("anomalies=%+v", result.Anomalies)
	}
	a := result.Anomalies[0]
	if a.CategoryCode != "tax" || a.Severity != "high" || a.Facts == "" || !strings.Contains(a.Facts, "revenue=12000") || a.Basis == "" || a.Impact == "" || a.Recommendation == "" || len(a.RequiredMaterials) == 0 || a.RuleVersion != "rule-version-42" {
		t.Fatalf("incomplete trusted anomaly=%+v", a)
	}
}

func TestMissingDocumentsRiskIsDetailedAndDeterministic(t *testing.T) {
	raw := missingDocumentsRisk([]string{"bank_statement", "contract", "tax_receipt"})
	got, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 40}, Risks: []map[string]any{raw}})
	if err != nil {
		t.Fatal(err)
	}
	a := got.Anomalies[0]
	if a.Code != "DOCUMENTS_INCOMPLETE" || a.CategoryCode != "documents" || a.Severity != "medium" || !strings.Contains(a.Facts, "银行流水") || a.Basis == "" || a.Impact == "" || a.Recommendation == "" || len(a.RequiredMaterials) != 3 || a.RuleVersion != "report-snapshot-v1" {
		t.Fatalf("anomaly=%+v", a)
	}
}

func TestTrustedEmptySourcesAreExplicitlyInsufficient(t *testing.T) {
	got, err := BuildStructured(Input{PeriodKey: "2026-07",
		Statistics:   map[string]any{"trustedSnapshot": true, "sourceAvailable": false, "transactionCount": 0},
		Completeness: map[string]any{"trustedSnapshot": true, "rate": 0, "confirmedDocumentCount": 0},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Summary.Conclusion != ConclusionAttention || !strings.Contains(got.Summary.DataNotice, "数据不足") {
		t.Fatalf("trusted empty snapshot must not be normal: %+v", got)
	}
	for _, category := range got.Categories {
		if category.Status != "insufficient" {
			t.Fatalf("empty trusted category %s must be insufficient: %+v", category.Code, got.Categories)
		}
	}
}

func TestTrustedSingleCategoryRiskKeepsOtherUnavailableCategoriesInsufficient(t *testing.T) {
	got, err := BuildStructured(Input{PeriodKey: "2026-07",
		Statistics:   map[string]any{"trustedSnapshot": true, "sourceAvailable": false, "transactionCount": 0},
		Completeness: map[string]any{"trustedSnapshot": true, "rate": 0, "confirmedDocumentCount": 0},
		Risks:        []map[string]any{{"code": "ZERO_WITH_CASHFLOW", "categoryCode": "tax", "severity": "high", "title": "服务端税务风险", "facts": "服务端事实", "basis": "服务端依据", "impact": "服务端影响", "recommendation": "服务端建议", "ruleVersion": "v1"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, category := range got.Categories {
		if category.Code == "tax" {
			if category.Status != "high" {
				t.Fatalf("tax category=%+v", category)
			}
			continue
		}
		if category.Status != "insufficient" {
			t.Fatalf("unavailable category %s incorrectly normal: %+v", category.Code, got.Categories)
		}
	}
}

func TestPeriodUnixRangeUsesHalfOpenMonthBoundary(t *testing.T) {
	start, end, err := periodUnixRange("2026-07")
	if err != nil || start >= end {
		t.Fatalf("start=%d end=%d err=%v", start, end, err)
	}
	_, _, err = periodUnixRange("2026-13")
	if err == nil {
		t.Fatal("invalid month accepted")
	}
}
