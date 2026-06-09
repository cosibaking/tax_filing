package tax

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
)

var employmentNAKeys = map[string]bool{
	"payroll_tax":      true,
	"social_insurance": true,
}

func enrichChecklistForOpc(ctx context.Context, memberId, opcId uint64, items []compliancein.ChecklistItem, period string) []compliancein.ChecklistItem {
	info, err := shared.LoadEmploymentStatus(ctx, opcId)
	if err == nil && info.Status == shared.EmploymentNoEmployee {
		items = enrichChecklistByEmployment(items)
	}
	month := strings.TrimSpace(period)
	if len(month) >= 7 {
		month = month[:7]
	} else {
		month = gtime.Now().Format("Y-m")
	}
	items = enrichChecklistWithConsistency(ctx, memberId, items, month)
	return items
}

func enrichChecklistWithConsistency(ctx context.Context, memberId uint64, items []compliancein.ChecklistItem, month string) []compliancein.ChecklistItem {
	if memberId == 0 {
		return items
	}
	report, err := service.ComplianceLedger().GetIncomeConsistency(ctx, &compliancein.IncomeConsistencyInp{
		MemberId: memberId,
		Month:    month,
	})
	out := make([]compliancein.ChecklistItem, 0, len(items))
	for _, item := range items {
		if item.Key == "income_match" && err == nil && report != nil {
			item.Checked = report.IncomeMatchPassed
			if len(report.Hints) > 0 {
				item.Hint = report.Hints[0]
			}
			if report.Status == "ok" {
				item.Hint = "系统自动比对：台账与银行/平台数据一致（偏差≤1%）"
			}
		}
		out = append(out, item)
	}
	return out
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
