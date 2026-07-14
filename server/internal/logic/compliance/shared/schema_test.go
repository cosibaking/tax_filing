package shared

import (
	"os"
	"strings"
	"testing"
)

func TestComplianceAssistantMigrationContainsCoreTables(t *testing.T) {
	body, err := os.ReadFile("../../../../cmd_tools/migrate/1.5.6_compliance_assistant_core.mysql.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	sql := string(body)
	for _, table := range []string{
		"xy_enterprise_profile",
		"xy_compliance_rule",
		"xy_compliance_rule_version",
		"xy_enterprise_compliance_task",
		"xy_task_reminder",
		"xy_business_document",
		"xy_document_relation",
		"xy_risk_event",
		"xy_compliance_monthly_report",
		"xy_service_ticket",
		"xy_service_ticket_record",
		"xy_ai_audit_log",
	} {
		needle := "CREATE TABLE IF NOT EXISTS `" + table + "`"
		if !strings.Contains(sql, needle) {
			t.Errorf("migration missing %s", table)
		}
	}

	for _, index := range []string{
		"uk_profile_version",
		"uk_rule_code",
		"uk_rule_version",
		"uk_task_idempotency",
		"uk_reminder_idempotency",
		"uk_report_version",
	} {
		if !strings.Contains(sql, index) {
			t.Errorf("migration missing unique index %s", index)
		}
	}
}
