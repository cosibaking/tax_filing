-- CreateTable
CREATE TABLE `members` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `phone` VARCHAR(20) NOT NULL,
    `password_hash` VARCHAR(255) NULL,
    `name` VARCHAR(64) NULL,
    `email` VARCHAR(128) NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `deleted_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `members_phone_key`(`phone`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `admin_users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL,
    `role` VARCHAR(32) NOT NULL DEFAULT 'advisor',
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `deleted_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `admin_users_username_key`(`username`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `otp_codes` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `phone` VARCHAR(20) NOT NULL,
    `code` VARCHAR(8) NOT NULL,
    `purpose` VARCHAR(32) NOT NULL,
    `expires_at` DATETIME(3) NOT NULL,
    `used` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `otp_codes_phone_purpose_idx`(`phone`, `purpose`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sys_attachments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `member_id` BIGINT NULL,
    `file_name` VARCHAR(255) NOT NULL,
    `mime_type` VARCHAR(128) NOT NULL,
    `size` INTEGER NOT NULL,
    `storage_key` VARCHAR(512) NOT NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `deleted_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `member_notices` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `member_id` BIGINT NOT NULL,
    `type` VARCHAR(64) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `content` TEXT NOT NULL,
    `read` BOOLEAN NOT NULL DEFAULT false,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `member_notices_member_id_idx`(`member_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `compliance_audit_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `entity_type` VARCHAR(64) NOT NULL,
    `entity_id` BIGINT NOT NULL,
    `action` VARCHAR(64) NOT NULL,
    `operator_id` BIGINT NOT NULL,
    `operator_type` VARCHAR(16) NOT NULL,
    `before` JSON NULL,
    `after` JSON NULL,
    `ip` VARCHAR(45) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `compliance_audit_logs_entity_type_entity_id_idx`(`entity_type`, `entity_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `compliance_consents` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `member_id` BIGINT NOT NULL,
    `order_id` BIGINT NULL,
    `type` VARCHAR(32) NOT NULL,
    `document_version` VARCHAR(16) NOT NULL,
    `ip` VARCHAR(45) NULL,
    `user_agent` TEXT NULL,
    `agreed_at` DATETIME(3) NOT NULL,

    INDEX `compliance_consents_member_id_idx`(`member_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `compliance_diagnoses` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `member_id` BIGINT NULL,
    `platforms` JSON NOT NULL,
    `monthly_income_range` VARCHAR(32) NOT NULL,
    `annual_cost_estimate` DECIMAL(14, 2) NOT NULL,
    `existing_entity` VARCHAR(32) NOT NULL,
    `has_filed_tax` BOOLEAN NOT NULL,
    `tax_bureau_contact` BOOLEAN NOT NULL,
    `recommended_plan` VARCHAR(32) NOT NULL,
    `tax_comparison` JSON NOT NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `service_plans` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL,
    `tier` VARCHAR(32) NOT NULL,
    `monthly_price` DECIMAL(10, 2) NOT NULL,
    `features` JSON NOT NULL,
    `active` BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `service_orders` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `member_id` BIGINT NOT NULL,
    `plan_id` BIGINT NOT NULL,
    `diagnosis_id` BIGINT NULL,
    `status` VARCHAR(32) NOT NULL,
    `amount` DECIMAL(10, 2) NOT NULL,
    `signed_at` DATETIME(3) NULL,
    `contract_file_id` BIGINT NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `service_orders_member_id_idx`(`member_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `opc_entities` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `member_id` BIGINT NOT NULL,
    `company_name` VARCHAR(128) NULL,
    `credit_code` VARCHAR(32) NULL,
    `status` VARCHAR(32) NOT NULL,
    `business_scope` TEXT NULL,
    `register_address` VARCHAR(255) NULL,
    `id_card_encrypted` TEXT NULL,
    `phone` VARCHAR(20) NULL,
    `email` VARCHAR(128) NULL,
    `bank_account_enc` TEXT NULL,
    `license_file_id` BIGINT NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `opc_entities_member_id_key`(`member_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `opc_progress_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `step` VARCHAR(32) NOT NULL,
    `status` VARCHAR(32) NOT NULL,
    `note` TEXT NULL,
    `operated_by` BIGINT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `opc_progress_logs_opc_id_idx`(`opc_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `expense_categories` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(32) NOT NULL,
    `name` VARCHAR(64) NOT NULL,
    `voucher_hint` TEXT NOT NULL,
    `active` BOOLEAN NOT NULL DEFAULT true,

    UNIQUE INDEX `expense_categories_code_key`(`code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `income_entries` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `platform` VARCHAR(32) NOT NULL,
    `category` VARCHAR(32) NOT NULL,
    `gross_amount` DECIMAL(14, 2) NOT NULL,
    `platform_fee` DECIMAL(14, 2) NOT NULL,
    `net_amount` DECIMAL(14, 2) NOT NULL,
    `occurred_at` DATETIME(3) NOT NULL,
    `source` VARCHAR(16) NOT NULL DEFAULT 'manual',
    `note` TEXT NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `deleted_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `income_entries_opc_id_occurred_at_idx`(`opc_id`, `occurred_at`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `expense_entries` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `category` VARCHAR(32) NOT NULL,
    `amount` DECIMAL(14, 2) NOT NULL,
    `invoice_type` VARCHAR(16) NOT NULL,
    `attachment_id` BIGINT NULL,
    `description` TEXT NULL,
    `warning_flag` BOOLEAN NOT NULL DEFAULT false,
    `occurred_at` DATETIME(3) NOT NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `deleted_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `expense_entries_opc_id_occurred_at_idx`(`opc_id`, `occurred_at`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `bank_transactions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `occurred_at` DATETIME(3) NOT NULL,
    `amount` DECIMAL(14, 2) NOT NULL,
    `description` TEXT NULL,
    `matched` BOOLEAN NOT NULL DEFAULT false,
    `income_id` BIGINT NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `bank_transactions_opc_id_idx`(`opc_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ledger_accounts` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(16) NOT NULL,
    `name` VARCHAR(64) NOT NULL,
    `type` VARCHAR(16) NOT NULL,

    UNIQUE INDEX `ledger_accounts_code_key`(`code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ledger_vouchers` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `period` VARCHAR(7) NOT NULL,
    `debit_account` VARCHAR(16) NOT NULL,
    `credit_account` VARCHAR(16) NOT NULL,
    `amount` DECIMAL(14, 2) NOT NULL,
    `ref_type` VARCHAR(32) NOT NULL,
    `ref_id` BIGINT NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `ledger_vouchers_opc_id_period_idx`(`opc_id`, `period`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `profit_summaries` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `year` INTEGER NOT NULL,
    `month` INTEGER NOT NULL,
    `revenue` DECIMAL(14, 2) NOT NULL,
    `cost` DECIMAL(14, 2) NOT NULL,
    `profit` DECIMAL(14, 2) NOT NULL,
    `cumulative_profit` DECIMAL(14, 2) NOT NULL,

    UNIQUE INDEX `profit_summaries_opc_id_year_month_key`(`opc_id`, `year`, `month`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `tax_filing_tasks` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `tax_type` VARCHAR(32) NOT NULL,
    `period` VARCHAR(16) NOT NULL,
    `due_date` DATETIME(3) NOT NULL,
    `status` VARCHAR(16) NOT NULL,
    `calculated_amount` DECIMAL(14, 2) NOT NULL,
    `filed_amount` DECIMAL(14, 2) NULL,
    `receipt_file_id` BIGINT NULL,
    `checklist` JSON NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `tax_filing_tasks_opc_id_period_idx`(`opc_id`, `period`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `monthly_statements` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `opc_id` BIGINT NOT NULL,
    `year` INTEGER NOT NULL,
    `month` INTEGER NOT NULL,
    `summary` JSON NOT NULL,
    `pdf_file_id` BIGINT NULL,
    `sent_at` DATETIME(3) NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `monthly_statements_opc_id_year_month_key`(`opc_id`, `year`, `month`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `member_notices` ADD CONSTRAINT `member_notices_member_id_fkey` FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `compliance_consents` ADD CONSTRAINT `compliance_consents_member_id_fkey` FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `compliance_diagnoses` ADD CONSTRAINT `compliance_diagnoses_member_id_fkey` FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `service_orders` ADD CONSTRAINT `service_orders_member_id_fkey` FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `service_orders` ADD CONSTRAINT `service_orders_plan_id_fkey` FOREIGN KEY (`plan_id`) REFERENCES `service_plans`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `opc_entities` ADD CONSTRAINT `opc_entities_member_id_fkey` FOREIGN KEY (`member_id`) REFERENCES `members`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `opc_progress_logs` ADD CONSTRAINT `opc_progress_logs_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `income_entries` ADD CONSTRAINT `income_entries_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `expense_entries` ADD CONSTRAINT `expense_entries_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `bank_transactions` ADD CONSTRAINT `bank_transactions_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ledger_vouchers` ADD CONSTRAINT `ledger_vouchers_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `profit_summaries` ADD CONSTRAINT `profit_summaries_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `tax_filing_tasks` ADD CONSTRAINT `tax_filing_tasks_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `monthly_statements` ADD CONSTRAINT `monthly_statements_opc_id_fkey` FOREIGN KEY (`opc_id`) REFERENCES `opc_entities`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

