-- AlterTable: OPC 落地扩展字段
ALTER TABLE `opc_entities`
  ADD COLUMN `proposed_names` JSON NULL AFTER `status`,
  ADD COLUMN `registered_capital` DECIMAL(12, 2) NULL AFTER `proposed_names`,
  ADD COLUMN `capital_term_years` INT NULL AFTER `registered_capital`,
  ADD COLUMN `business_term_type` VARCHAR(16) NULL AFTER `capital_term_years`,
  ADD COLUMN `business_term_end` DATE NULL AFTER `business_term_type`,
  ADD COLUMN `register_province` VARCHAR(32) NULL AFTER `business_scope`,
  ADD COLUMN `register_city` VARCHAR(32) NULL AFTER `register_province`,
  ADD COLUMN `register_district` VARCHAR(32) NULL AFTER `register_city`,
  ADD COLUMN `address_proof_file_id` BIGINT NULL AFTER `register_address`,
  ADD COLUMN `legal_person_name` VARCHAR(64) NULL AFTER `address_proof_file_id`,
  ADD COLUMN `id_card_valid_from` DATE NULL AFTER `id_card_encrypted`,
  ADD COLUMN `id_card_valid_to` DATE NULL AFTER `id_card_valid_from`,
  ADD COLUMN `id_card_front_file_id` BIGINT NULL AFTER `id_card_valid_to`,
  ADD COLUMN `id_card_back_file_id` BIGINT NULL AFTER `id_card_front_file_id`,
  ADD COLUMN `ethnicity` VARCHAR(16) NULL AFTER `id_card_back_file_id`,
  ADD COLUMN `household_address` VARCHAR(255) NULL AFTER `ethnicity`,
  ADD COLUMN `residential_address` VARCHAR(255) NULL AFTER `household_address`,
  ADD COLUMN `esign_authorized` BOOLEAN NOT NULL DEFAULT false AFTER `email`,
  ADD COLUMN `established_at` DATE NULL AFTER `esign_authorized`,
  ADD COLUMN `taxpayer_type` VARCHAR(32) NOT NULL DEFAULT 'small_scale' AFTER `established_at`,
  ADD COLUMN `tax_activated_at` DATETIME(3) NULL AFTER `taxpayer_type`,
  ADD COLUMN `bank_name` VARCHAR(128) NULL AFTER `tax_activated_at`,
  ADD COLUMN `bank_receipt_file_id` BIGINT NULL AFTER `bank_account_enc`,
  ADD COLUMN `materials_submitted_at` DATETIME(3) NULL AFTER `license_file_id`,
  ADD COLUMN `materials_approved_at` DATETIME(3) NULL AFTER `materials_submitted_at`,
  ADD COLUMN `reject_note` TEXT NULL AFTER `materials_approved_at`;

-- 历史数据：已提交资料但未细分状态的，归入 materials_review
UPDATE `opc_entities`
SET `status` = 'materials_review'
WHERE `status` = 'materials'
  AND (`id_card_encrypted` IS NOT NULL OR `phone` IS NOT NULL);
