package ruleengine

import "testing"

func TestEvaluateAllExpression(t *testing.T) {
	expr := Expression{All: []Node{
		{Condition: &Condition{Field: "profile.region", Op: "eq", Value: "CN-BJ"}},
		{Condition: &Condition{Field: "profile.employeeCount", Op: "lte", Value: float64(5)}},
	}}

	matched, err := Evaluate(expr, map[string]any{
		"profile": map[string]any{"region": "CN-BJ", "employeeCount": float64(1)},
	})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if !matched {
		t.Fatal("expected expression to match")
	}
}

func TestEvaluateAnyAndIn(t *testing.T) {
	expr := Expression{Any: []Node{
		{Condition: &Condition{Field: "profile.taxpayerType", Op: "in", Value: []any{"small_scale", "general"}}},
		{Condition: &Condition{Field: "profile.invoiceEnabled", Op: "eq", Value: true}},
	}}

	matched, err := Evaluate(expr, map[string]any{
		"profile": map[string]any{"taxpayerType": "small_scale", "invoiceEnabled": false},
	})
	if err != nil || !matched {
		t.Fatalf("expected any expression to match, matched=%v err=%v", matched, err)
	}
}

func TestEvaluateExistsAndNotEqual(t *testing.T) {
	expr := Expression{All: []Node{
		{Condition: &Condition{Field: "period.revenue", Op: "exists", Value: true}},
		{Condition: &Condition{Field: "profile.vatPeriod", Op: "ne", Value: "monthly"}},
	}}

	matched, err := Evaluate(expr, map[string]any{
		"profile": map[string]any{"vatPeriod": "quarterly"},
		"period":  map[string]any{"revenue": float64(0)},
	})
	if err != nil || !matched {
		t.Fatalf("zero is an existing value, matched=%v err=%v", matched, err)
	}
}

func TestEvaluateRejectsUnknownField(t *testing.T) {
	expr := Expression{All: []Node{{Condition: &Condition{
		Field: "database.password", Op: "exists", Value: true,
	}}}}

	if _, err := Evaluate(expr, map[string]any{}); err == nil {
		t.Fatal("expected unknown field to be rejected")
	}
}

func TestEvaluateRejectsUnknownOperator(t *testing.T) {
	expr := Expression{All: []Node{{Condition: &Condition{
		Field: "profile.region", Op: "exec", Value: "CN-BJ",
	}}}}

	if _, err := Evaluate(expr, map[string]any{}); err == nil {
		t.Fatal("expected unknown operator to be rejected")
	}
}

func TestEvaluateNumericComparisons(t *testing.T) {
	data := map[string]any{"period": map[string]any{"rollingSales": float64(95)}}
	for _, tc := range []struct {
		op    string
		value float64
		want  bool
	}{
		{op: "gt", value: 90, want: true},
		{op: "gte", value: 95, want: true},
		{op: "lt", value: 100, want: true},
		{op: "lte", value: 94, want: false},
	} {
		expr := Expression{All: []Node{{Condition: &Condition{
			Field: "period.rollingSales", Op: tc.op, Value: tc.value,
		}}}}
		got, err := Evaluate(expr, data)
		if err != nil || got != tc.want {
			t.Errorf("op=%s got=%v want=%v err=%v", tc.op, got, tc.want, err)
		}
	}
}
