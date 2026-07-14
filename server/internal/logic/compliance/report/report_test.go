package report

import (
	"context"
	"encoding/json"
	"errors"
	"math"
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

func TestBuildStructuredRejectsNonFiniteCompletenessRate(t *testing.T) {
	for _, value := range []any{math.NaN(), math.Inf(1), float32(math.Inf(-1)), json.Number("NaN"), "Inf"} {
		got, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": value}})
		if err != nil {
			t.Fatal(err)
		}
		if got.Summary.CompletenessRate != 0 || got.Summary.Conclusion != ConclusionAttention || !strings.Contains(got.Summary.DataNotice, "数据不足") {
			t.Fatalf("rate(%v) summary=%+v", value, got.Summary)
		}
	}
}

func TestBuildStructuredRiskDefaultsAndDueDateOrdering(t *testing.T) {
	got, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 100}, Risks: []map[string]any{
		{"code": "no-date", "severity": "high"},
		{"code": "dated", "severity": "high", "dueDate": "2026-08-01"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if got.Anomalies[0].Code != "dated" || got.Anomalies[1].Code != "no-date" {
		t.Fatalf("due date order=%+v", got.Anomalies)
	}
	if !got.Anomalies[1].RequiresManualReview {
		t.Fatalf("high risk should require manual review by default: %+v", got.Anomalies[1])
	}
}

func TestBuildStructuredConclusionForCompleteData(t *testing.T) {
	medium, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 100}, Risks: []map[string]any{{"severity": "medium"}}})
	if err != nil {
		t.Fatal(err)
	}
	normal, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 100}})
	if err != nil {
		t.Fatal(err)
	}
	if medium.Summary.Conclusion != ConclusionAttention || normal.Summary.Conclusion != ConclusionNormal {
		t.Fatalf("medium=%+v normal=%+v", medium.Summary, normal.Summary)
	}
}

func TestRenderContentContainsRequiredSections(t *testing.T) {
	structured, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 75}, Risks: []map[string]any{{"severity": "medium"}}})
	if err != nil {
		t.Fatal(err)
	}
	content := RenderContent("2026-07", structured)
	for _, required := range []string{"2026-07 月度经营合规体检", "总体结论：需关注", "资料完整度：75.0%", "风险分布：高风险 0 项、中风险 1 项、低风险 0 项", "资料提示：", structured.Summary.DataNotice, "免责声明："} {
		if !strings.Contains(content, required) {
			t.Fatalf("content missing %q: %s", required, content)
		}
	}
}

func TestStructuredAndResultJSONUseStableLowerCamelFields(t *testing.T) {
	result, err := Generate(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 100}, Risks: []map[string]any{{"severity": "high"}}})
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	structured, ok := decoded["structuredReport"].(map[string]any)
	if !ok || structured["schemaVersion"] != float64(StructuredSchemaVersion) {
		t.Fatalf("json=%s", data)
	}
	for _, field := range []string{"summary", "categories", "anomalies"} {
		if _, ok := structured[field]; !ok {
			t.Fatalf("structured field %q missing: %s", field, data)
		}
	}
	summary := structured["summary"].(map[string]any)
	for _, field := range []string{"conclusion", "completenessRate", "highCount", "mediumCount", "lowCount", "dataNotice"} {
		if _, ok := summary[field]; !ok {
			t.Fatalf("summary field %q missing: %s", field, data)
		}
	}
	anomaly := structured["anomalies"].([]any)[0].(map[string]any)
	for _, field := range []string{"categoryCode", "requiredMaterials", "dueDate", "requiresManualReview", "ruleVersion"} {
		if _, ok := anomaly[field]; !ok {
			t.Fatalf("anomaly field %q missing: %s", field, data)
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

func TestApplyStructuredJSONMarksEmptySnapshotLegacy(t *testing.T) {
	item := MonthlyReport{StructuredReport: &StructuredReport{SchemaVersion: 99}}
	if err := item.ApplyStructuredJSON(" \n\t "); err != nil {
		t.Fatal(err)
	}
	if item.StructuredReport != nil || !item.Legacy {
		t.Fatalf("item=%+v", item)
	}
}

func TestApplyStructuredJSONRoundTripsValidSnapshot(t *testing.T) {
	want, err := BuildStructured(Input{PeriodKey: "2026-07", Statistics: map[string]any{"ok": true}, Completeness: map[string]any{"rate": 100}})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	item := MonthlyReport{}
	if err := item.ApplyStructuredJSON(string(raw)); err != nil {
		t.Fatal(err)
	}
	if item.Legacy || item.StructuredReport == nil || !reflect.DeepEqual(*item.StructuredReport, want) {
		t.Fatalf("item=%+v want=%+v", item, want)
	}
}

func TestApplyStructuredJSONRejectsInvalidJSONWithoutLeakingPartialData(t *testing.T) {
	item := MonthlyReport{StructuredReport: &StructuredReport{SchemaVersion: 99}}
	err := item.ApplyStructuredJSON(`{"schemaVersion":1,"summary":`)
	if err == nil || !strings.Contains(err.Error(), "结构化报告") {
		t.Fatalf("err=%v", err)
	}
	if item.StructuredReport != nil || !item.Legacy {
		t.Fatalf("item=%+v", item)
	}
}

func TestApplyStructuredJSONRejectsZeroSchemaVersion(t *testing.T) {
	item := MonthlyReport{}
	err := item.ApplyStructuredJSON(`{"schemaVersion":0}`)
	if err == nil || !strings.Contains(err.Error(), "schemaVersion") {
		t.Fatalf("err=%v", err)
	}
	if item.StructuredReport != nil || !item.Legacy {
		t.Fatalf("item=%+v", item)
	}
}

func TestNormalizeVersionConflict(t *testing.T) {
	for _, message := range []string{"Duplicate entry for key uk_report_version", "UNIQUE constraint failed: xy_compliance_monthly_report"} {
		err := normalizeVersionConflict(errors.New(message))
		if err == nil || err.Error() != "报告版本生成冲突，请重试" {
			t.Fatalf("message=%q err=%v", message, err)
		}
	}
	want := errors.New("connection lost")
	if got := normalizeVersionConflict(want); !errors.Is(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
