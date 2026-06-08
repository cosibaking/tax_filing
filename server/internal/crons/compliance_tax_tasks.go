package crons

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	cronlib "xygo/internal/library/cron"
	"xygo/internal/logic/compliance/tax"
)

func init() {
	cronlib.Register(&ComplianceTaxTasksTask{})
}

// ComplianceTaxTasksTask 每月生成增值税/附加税及季度/年度企税申报任务
type ComplianceTaxTasksTask struct{}

func (t *ComplianceTaxTasksTask) GetName() string { return "compliance_tax_tasks" }

func (t *ComplianceTaxTasksTask) Execute(ctx context.Context, params []string) (string, error) {
	now := time.Now()
	created, skipped, err := tax.GenerateTasks(ctx, now)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("tax tasks created=%d skipped=%d month=%04d-%02d", created, skipped, now.Year(), int(now.Month()))
	g.Log().Infof(ctx, "[cron:compliance_tax_tasks] %s", output)
	return output, nil
}
