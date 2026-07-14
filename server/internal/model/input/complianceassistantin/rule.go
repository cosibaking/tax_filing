package complianceassistantin

import "xygo/internal/logic/compliance/ruleengine"

type RuleSimulationInp struct {
	Expression ruleengine.Expression `json:"expression"`
	Facts      map[string]any        `json:"facts"`
}

type RuleSimulationModel struct {
	Matched bool `json:"matched"`
}
