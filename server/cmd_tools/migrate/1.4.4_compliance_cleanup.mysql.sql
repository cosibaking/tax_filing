-- Migration: 1.4.4
-- Description: 清理未使用表/字段，补充签约姓名落库字段

-- 删除未使用的会计科目表（凭证使用硬编码科目代码）
DROP TABLE IF EXISTS `xy_ledger_account`;

-- 订单：移除未实现的合同附件字段（幂等）
SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'xy_service_order' AND COLUMN_NAME = 'contract_file_id');
SET @sql := IF(@exist > 0, 'ALTER TABLE `xy_service_order` DROP COLUMN `contract_file_id`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 订单：补充签约姓名字段（幂等）
SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'xy_service_order' AND COLUMN_NAME = 'legal_name');
SET @sql := IF(@exist = 0, 'ALTER TABLE `xy_service_order` ADD COLUMN `legal_name` varchar(64) NOT NULL DEFAULT '''' COMMENT ''电子签约姓名'' AFTER `signed_at`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 对账单：移除未实现的 PDF 字段（MVP 以 JSON 摘要为准，幂等）
SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'xy_monthly_statement' AND COLUMN_NAME = 'pdf_file_id');
SET @sql := IF(@exist > 0, 'ALTER TABLE `xy_monthly_statement` DROP COLUMN `pdf_file_id`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
