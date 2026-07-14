package ruleengine

import (
	"fmt"
	"reflect"
	"strings"
)

type Condition struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Value any    `json:"value"`
}

type Node struct {
	Condition *Condition `json:"condition,omitempty"`
	All       []Node     `json:"all,omitempty"`
	Any       []Node     `json:"any,omitempty"`
}

type Expression struct {
	All []Node `json:"all,omitempty"`
	Any []Node `json:"any,omitempty"`
}

var allowedFields = map[string]struct{}{
	"profile.region":               {},
	"profile.entityType":           {},
	"profile.taxpayerType":         {},
	"profile.vatPeriod":            {},
	"profile.employeeCount":        {},
	"profile.invoiceEnabled":       {},
	"profile.hasRevenue":           {},
	"profile.hasPublicBankAccount": {},
	"profile.complexity":           {},
	"period.revenue":               {},
	"period.expense":               {},
	"period.invoiceAmount":         {},
	"period.receiptAmount":         {},
	"period.rollingSales":          {},
	"documents.completeness":       {},
	"documents.missingTypes":       {},
}

var allowedOperators = map[string]struct{}{
	"eq": {}, "ne": {}, "in": {}, "gt": {}, "gte": {}, "lt": {}, "lte": {}, "exists": {},
}

func Evaluate(expr Expression, data map[string]any) (bool, error) {
	if len(expr.All) > 0 && len(expr.Any) > 0 {
		return false, fmt.Errorf("expression cannot contain both all and any")
	}
	if len(expr.All) > 0 {
		return evaluateNodes(expr.All, data, true)
	}
	if len(expr.Any) > 0 {
		return evaluateNodes(expr.Any, data, false)
	}
	return false, fmt.Errorf("expression must contain all or any")
}

func evaluateNodes(nodes []Node, data map[string]any, requireAll bool) (bool, error) {
	if len(nodes) == 0 {
		return false, fmt.Errorf("expression group cannot be empty")
	}
	for _, node := range nodes {
		matched, err := evaluateNode(node, data)
		if err != nil {
			return false, err
		}
		if requireAll && !matched {
			return false, nil
		}
		if !requireAll && matched {
			return true, nil
		}
	}
	return requireAll, nil
}

func evaluateNode(node Node, data map[string]any) (bool, error) {
	kinds := 0
	if node.Condition != nil {
		kinds++
	}
	if len(node.All) > 0 {
		kinds++
	}
	if len(node.Any) > 0 {
		kinds++
	}
	if kinds != 1 {
		return false, fmt.Errorf("rule node must contain exactly one of condition, all or any")
	}
	if node.Condition != nil {
		return evaluateCondition(*node.Condition, data)
	}
	if len(node.All) > 0 {
		return evaluateNodes(node.All, data, true)
	}
	return evaluateNodes(node.Any, data, false)
}

func evaluateCondition(condition Condition, data map[string]any) (bool, error) {
	if _, ok := allowedFields[condition.Field]; !ok {
		return false, fmt.Errorf("field %q is not allowed", condition.Field)
	}
	if _, ok := allowedOperators[condition.Op]; !ok {
		return false, fmt.Errorf("operator %q is not allowed for field %q", condition.Op, condition.Field)
	}

	actual, exists := resolveField(data, condition.Field)
	if condition.Op == "exists" {
		want, ok := condition.Value.(bool)
		if !ok {
			return false, fmt.Errorf("operator exists requires boolean value for field %q", condition.Field)
		}
		return exists == want, nil
	}
	if !exists {
		return false, nil
	}

	switch condition.Op {
	case "eq":
		return valuesEqual(actual, condition.Value), nil
	case "ne":
		return !valuesEqual(actual, condition.Value), nil
	case "in":
		items, ok := condition.Value.([]any)
		if !ok {
			return false, fmt.Errorf("operator in requires array value for field %q", condition.Field)
		}
		for _, item := range items {
			if valuesEqual(actual, item) {
				return true, nil
			}
		}
		return false, nil
	default:
		left, leftOK := number(actual)
		right, rightOK := number(condition.Value)
		if !leftOK || !rightOK {
			return false, fmt.Errorf("operator %s requires numeric values for field %q", condition.Op, condition.Field)
		}
		switch condition.Op {
		case "gt":
			return left > right, nil
		case "gte":
			return left >= right, nil
		case "lt":
			return left < right, nil
		case "lte":
			return left <= right, nil
		}
	}
	return false, fmt.Errorf("unsupported operator %q", condition.Op)
}

func resolveField(data map[string]any, path string) (any, bool) {
	var current any = data
	for _, part := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = object[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func valuesEqual(left, right any) bool {
	if l, ok := number(left); ok {
		if r, ok := number(right); ok {
			return l == r
		}
	}
	return reflect.DeepEqual(left, right)
}

func number(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}
