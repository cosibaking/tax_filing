-- Migration: 1.5.7
-- Description: 月度体检报告增加结构化快照字段

SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'xy_compliance_monthly_report' AND COLUMN_NAME = 'structured_json');
SET @sql := IF(@exist = 0, 'ALTER TABLE `xy_compliance_monthly_report` ADD COLUMN `structured_json` JSON NULL AFTER `risk_snapshot_json`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
