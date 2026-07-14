package report

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type failingNarrator struct{}

func (failingNarrator) Narrate(context.Context, Input) (string, string, error) {
	return "", "", errors.New("model down")
}

func TestGenerateMarksMissingData(t *testing.T) {
	result, err := Generate(Input{PeriodKey: "2026-07"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Content, "数据不足") {
		t.Fatalf("content=%s", result.Content)
	}
}

func TestBuildStructuredEmptyDataNeedsAttentionAndFixedCategories(t *testing.T) {
	got, err := BuildStructured(Input{PeriodKey: "2026-07"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Summary.Conclusion != ConclusionAttention || got.Summary.DataNotice == "" {
		t.Fatalf("summary=%+v", got.Summary)
	}
	wantCodes := []string{"business", "documents", "tax", "funds", "employment", "annual"}
	if len(got.Categories) != len(wantCodes) {
		t.Fatalf("categories=%+v", got.Categories)
	}
	for i, code := range wantCodes {
		if got.Categories[i].Code != code || got.Categories[i].Status == "normal" || got.Categories[i].Summary == "" || len(got.Categories[i].Checks) == 0 {
			t.Fatalf("category[%d]=%+v", i, got.Categories[i])
		}
	}
}

func TestBuildStructuredNormalizesAndSortsRisks(t *testing.T) {
	in := Input{PeriodKey: "2026-07", Statistics: map[string]any{"revenue": 1}, Completeness: map[string]any{"rate": 88}, Risks: []map[string]any{
		{"code": "late", "categoryCode": "tax", "severity": "high", "dueDate": "2026-08-20"},
		{"code": "low", "categoryCode": "business", "severity": "low"},
		{"code": "early", "categoryCode": "tax", "severity": "HIGH", "dueDate": "2026-08-01"},
	}}
	got, err := BuildStructured(in)
	if err != nil {
		t.Fatal(err)
	}
	if got.Summary.HighCount != 2 || got.Summary.LowCount != 1 || got.Summary.Conclusion != ConclusionUrgent {
		t.Fatalf("summary=%+v", got.Summary)
	}
	if got.Anomalies[0].Code != "early" || got.Anomalies[1].Code != "late" || got.Anomalies[2].Code != "low" {
		t.Fatalf("unexpected order: %+v", got.Anomalies)
	}
	for _, anomaly := range got.Anomalies {
		if anomaly.Facts == "" || anomaly.Basis == "" || anomaly.Impact == "" || anomaly.Recommendation == "" || anomaly.RequiredMaterials == nil {
			t.Fatalf("unsafe anomaly defaults: %+v", anomaly)
		}
	}
}

func TestBuildStructuredRejectsInvalidPeriod(t *testing.T) {
	if _, err := BuildStructured(Input{PeriodKey: "2026-13"}); err == nil {
		t.Fatal("invalid calendar month must be rejected")
	}
}

func TestBuildStructuredClampsCompletenessRate(t *testing.T) {
	for _, tc := range []struct {
		value any
		want  float64
	}{{int64(-2), 0}, {float32(120.5), 100}, {"72.5", 72.5}} {
		got, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": tc.value}})
		if err != nil {
			t.Fatal(err)
		}
		if got.Summary.CompletenessRate != tc.want {
			t.Fatalf("rate(%v)=%v want %v", tc.value, got.Summary.CompletenessRate, tc.want)
		}
	}
}

func TestBuildStructuredSafelyNormalizesUnknownSeverityAndCategory(t *testing.T) {
	got, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 100}, Risks: []map[string]any{{"severity": "critical", "categoryCode": "unknown", "requiresManualReview": false}}})
	if err != nil {
		t.Fatal(err)
	}
	anomaly := got.Anomalies[0]
	if anomaly.Severity != "medium" || anomaly.CategoryCode != "documents" || anomaly.RequiresManualReview {
		t.Fatalf("anomaly=%+v", anomaly)
	}
}

type successfulNarrator struct{}

func (successfulNarrator) Narrate(context.Context, Input) (string, string, error) {
	return "AI润色内容", "test-model", nil
}

func TestNarratorOnlyReplacesPresentationNotStructuredFacts(t *testing.T) {
	in := Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 90}, Risks: []map[string]any{{"code": "h1", "severity": "high", "categoryCode": "tax"}}}
	want, err := BuildStructured(in)
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewService(nil, successfulNarrator{}).Build(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if got.Content != "AI润色内容" || got.AIModel != "test-model" || !reflect.DeepEqual(got.Structured, want) {
		t.Fatalf("result=%+v want structured=%+v", got, want)
	}
}

func TestNarratorFailureFallsBackToDeterministicReport(t *testing.T) {
	svc := NewService(nil, failingNarrator{})
	result, err := svc.Build(context.Background(), Input{PeriodKey: "2026-07"})
	if err != nil {
		t.Fatal(err)
	}
	if result.AIModel != "" || !strings.Contains(result.Content, "数据不足") {
		t.Fatalf("result=%+v", result)
	}
}

func TestPublishedReportIsImmutable(t *testing.T) {
	if err := ValidateMutation("published"); err == nil {
		t.Fatal("published report must be immutable")
	}
}
