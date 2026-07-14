-- Migration: 1.5.7
-- Description: 月度体检报告增加结构化快照字段

ALTER TABLE xy_compliance_monthly_report ADD COLUMN IF NOT EXISTS structured_json jsonb NULL;
