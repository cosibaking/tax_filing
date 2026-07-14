package crons

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	cronlib "xygo/internal/library/cron"
	"xygo/internal/logic/compliance/notice"
)

func init() {
	cronlib.Register(&ComplianceAssistantReminderTask{})
	cronlib.Register(&ComplianceAssistantRiskTask{})
	cronlib.Register(&ComplianceAssistantReportTask{})
}

type ComplianceAssistantReminderTask struct{}

func (*ComplianceAssistantReminderTask) GetName() string { return "compliance_assistant_reminder" }
func (*ComplianceAssistantReminderTask) Execute(ctx context.Context, _ []string) (string, error) {
	if !g.Cfg().MustGet(ctx, "complianceAssistant.enabled", false).Bool() || !g.Cfg().MustGet(ctx, "complianceAssistant.taskCronEnabled", false).Bool() {
		return "skip: disabled", nil
	}
	now := uint64(time.Now().Unix())
	var rows []struct {
		ID, TaskID, MemberID uint64
		Title                string
		Attempts             int
	}
	err := g.DB().Model("xy_task_reminder r").Ctx(ctx).LeftJoin("xy_enterprise_compliance_task t", "t.id=r.task_id").Fields("r.id,r.task_id,r.attempts,t.member_id,t.title").Where("r.status", "pending").WhereLTE("r.planned_at", now).Where("t.deleted", 0).Limit(200).Scan(&rows)
	if err != nil {
		return "", err
	}
	sent, failed := 0, 0
	for _, row := range rows {
		err = notice.SendToMember(ctx, row.MemberID, "合规任务提醒", fmt.Sprintf("%s 即将到期，请及时处理。", row.Title))
		if err == nil {
			_, err = g.DB().Model("xy_task_reminder").Ctx(ctx).Where("id", row.ID).Data(g.Map{"status": "sent", "sent_at": now, "update_time": now}).Update()
			sent++
		} else {
			attempts, status := row.Attempts+1, "pending"
			if attempts >= 3 {
				status = "manual_attention"
			}
			_, _ = g.DB().Model("xy_task_reminder").Ctx(ctx).Where("id", row.ID).Data(g.Map{"status": status, "attempts": attempts, "last_error": err.Error(), "update_time": now}).Update()
			failed++
		}
	}
	return fmt.Sprintf("sent=%d failed=%d", sent, failed), nil
}

type ComplianceAssistantRiskTask struct{}

func (*ComplianceAssistantRiskTask) GetName() string { return "compliance_assistant_risk_scan" }
func (*ComplianceAssistantRiskTask) Execute(ctx context.Context, _ []string) (string, error) {
	if !g.Cfg().MustGet(ctx, "complianceAssistant.enabled", false).Bool() || !g.Cfg().MustGet(ctx, "complianceAssistant.riskCronEnabled", false).Bool() {
		return "skip: disabled", nil
	}
	return "ready: risk scans require confirmed period facts", nil
}

type ComplianceAssistantReportTask struct{}

func (*ComplianceAssistantReportTask) GetName() string { return "compliance_assistant_monthly_report" }
func (*ComplianceAssistantReportTask) Execute(ctx context.Context, _ []string) (string, error) {
	if !g.Cfg().MustGet(ctx, "complianceAssistant.enabled", false).Bool() || !g.Cfg().MustGet(ctx, "complianceAssistant.reportCronEnabled", false).Bool() {
		return "skip: disabled", nil
	}
	return "ready: reports are generated from confirmed data", nil
}
