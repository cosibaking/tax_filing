package crons

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	cronlib "xygo/internal/library/cron"
	"xygo/internal/logic/compliance/statement"
)

func init() {
	cronlib.Register(&ComplianceStatementGenerateTask{})
}

// ComplianceStatementGenerateTask F-61：每月16-20日生成上月对账单
type ComplianceStatementGenerateTask struct{}

func (t *ComplianceStatementGenerateTask) GetName() string { return "compliance_statement_generate" }

func (t *ComplianceStatementGenerateTask) Execute(ctx context.Context, params []string) (string, error) {
	now := time.Now()
	day := now.Day()
	if day < 16 || day > 20 {
		return fmt.Sprintf("skip: today is day %d, only runs on 16th-20th", day), nil
	}

	year, month := statement.PreviousMonth(now)
	created, skipped, err := statement.GenerateForPeriod(ctx, year, month)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("statements created=%d skipped=%d period=%04d-%02d", created, skipped, year, month)
	g.Log().Infof(ctx, "[cron:compliance_statement_generate] %s", output)
	return output, nil
}
