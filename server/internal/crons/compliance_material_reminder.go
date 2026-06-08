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
	cronlib.Register(&ComplianceMaterialReminderTask{})
}

// ComplianceMaterialReminderTask F-65：每月1-5日催收材料上传
type ComplianceMaterialReminderTask struct{}

func (t *ComplianceMaterialReminderTask) GetName() string { return "compliance_material_reminder" }

func (t *ComplianceMaterialReminderTask) Execute(ctx context.Context, params []string) (string, error) {
	now := time.Now()
	day := now.Day()
	if day < 1 || day > 5 {
		return fmt.Sprintf("skip: today is day %d, only runs on 1st-5th", day), nil
	}

	// 注册资料未提交的 OPC 会员
	var pendingRegs []struct {
		MemberId uint64 `json:"member_id"`
	}
	err := g.DB().Model("xy_opc_entity").Ctx(ctx).
		Where("deleted", 0).
		WhereIn("status", []string{"pending", "materials"}).
		Where("materials_submitted_at", 0).
		Fields("member_id").
		Scan(&pendingRegs)
	if err != nil {
		return "", fmt.Errorf("query pending materials failed: %v", err)
	}

	// 已激活 OPC：上月无收入台账记录
	prev := now.AddDate(0, -1, 0)
	period := fmt.Sprintf("%04d-%02d", prev.Year(), int(prev.Month()))
	start, end := monthUnixRange(prev.Year(), int(prev.Month()))

	var activeOpc []struct {
		Id       uint64 `json:"id"`
		MemberId uint64 `json:"member_id"`
	}
	_ = g.DB().Model("xy_opc_entity").Ctx(ctx).
		Where("status", "active").
		Where("deleted", 0).
		Fields("id, member_id").
		Scan(&activeOpc)

	sent := 0
	seen := make(map[uint64]bool)

	for _, row := range pendingRegs {
		if seen[row.MemberId] {
			continue
		}
		seen[row.MemberId] = true
		title := "请尽快提交OPC注册资料"
		content := "您尚未完成 OPC 注册资料提交，请登录合规中心尽快上传，以免影响设立进度。"
		if err = notice.SendToMember(ctx, row.MemberId, title, content); err == nil {
			sent++
		}
	}

	for _, opc := range activeOpc {
		if seen[opc.MemberId] {
			continue
		}
		count, _ := g.DB().Model("xy_income_entry").Ctx(ctx).
			Where("opc_id", opc.Id).
			Where("deleted", 0).
			WhereBetween("occurred_at", start, end).
			Count()
		if count > 0 {
			continue
		}
		seen[opc.MemberId] = true
		title := fmt.Sprintf("请上传%s收入台账", period)
		content := fmt.Sprintf("您尚未录入 %s 月收入台账，请在每月5日前完成上传，以便顾问记账与申报。", period)
		if err = notice.SendToMember(ctx, opc.MemberId, title, content); err == nil {
			sent++
		}
	}

	output := fmt.Sprintf("material reminders sent=%d (day=%d)", sent, day)
	g.Log().Infof(ctx, "[cron:compliance_material_reminder] %s", output)
	return output, nil
}

func monthUnixRange(year, month int) (uint64, uint64) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return uint64(start.Unix()), uint64(end.Unix())
}
