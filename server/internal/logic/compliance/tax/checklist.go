package tax

import (
	"context"

	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
)

var employmentNAKeys = map[string]bool{
	"payroll_tax":      true,
	"social_insurance": true,
}

func enrichChecklistForOpc(ctx context.Context, opcId uint64, items []compliancein.ChecklistItem) []compliancein.ChecklistItem {
	info, err := shared.LoadEmploymentStatus(ctx, opcId)
	if err != nil || info.Status != shared.EmploymentNoEmployee {
		return items
	}
	return enrichChecklistByEmployment(items)
}

func enrichChecklistByEmployment(items []compliancein.ChecklistItem) []compliancein.ChecklistItem {
	out := make([]compliancein.ChecklistItem, 0, len(items))
	for _, item := range items {
		if employmentNAKeys[item.Key] {
			item.Checked = true
			item.Na = true
			item.Hint = "不适用：无雇员"
		}
		out = append(out, item)
	}
	return out
}
