-- Migration: 1.4.0
-- Description: OPC 合规 MVP 数据表、种子数据与后台菜单

-- ============================================================
-- 一、合规诊断
-- ============================================================
CREATE TABLE IF NOT EXISTS `xy_compliance_diagnosis` (
    `id`                    bigint unsigned NOT NULL AUTO_INCREMENT,
    `member_id`             bigint unsigned NOT NULL DEFAULT 0 COMMENT '会员ID（匿名为0）',
    `platforms`             json            DEFAULT NULL COMMENT '直播平台列表',
    `monthly_income_range`  varchar(32)     NOT NULL DEFAULT '' COMMENT '月收入区间',
    `annual_cost_estimate`  decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '年成本估算',
    `existing_entity`       varchar(32)     NOT NULL DEFAULT '' COMMENT '现有主体：none/individual/opc/other',
    `has_filed_tax`         varchar(16)     NOT NULL DEFAULT '' COMMENT '报税状态：yes/no/unsure',
    `tax_bureau_contact`    tinyint         NOT NULL DEFAULT 0 COMMENT '是否被税局联系：1=是 0=否',
    `notes`                 text            COMMENT '补充说明',
    `cost_breakdown`        json            DEFAULT NULL COMMENT '成本细项',
    `recommended_plan`      varchar(32)     NOT NULL DEFAULT '' COMMENT '推荐方案：opc/individual/labor/transitional',
    `tax_comparison`        json            DEFAULT NULL COMMENT '三方案税负对比结果',
    `deleted`               tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`           bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`           bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_compliance_diagnosis_member_id` (`member_id`),
    KEY `idx_compliance_diagnosis_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规诊断记录';

-- ============================================================
-- 二、服务套餐与订单
-- ============================================================
CREATE TABLE IF NOT EXISTS `xy_service_plan` (
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`           varchar(64)     NOT NULL DEFAULT '' COMMENT '套餐名称',
    `tier`           varchar(32)     NOT NULL DEFAULT '' COMMENT '档位：basic/advanced/premium',
    `monthly_price`  decimal(10,2)   NOT NULL DEFAULT 0.00 COMMENT '月费（元）',
    `price_display`  varchar(32)     NOT NULL DEFAULT '' COMMENT '展示价格（如面议）',
    `features`       json            DEFAULT NULL COMMENT '包含项 JSON',
    `sort`           int             NOT NULL DEFAULT 0 COMMENT '排序',
    `status`         tinyint         NOT NULL DEFAULT 1 COMMENT '状态：1=启用 0=禁用',
    `create_time`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_service_plan_tier` (`tier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规服务套餐';

CREATE TABLE IF NOT EXISTS `xy_service_order` (
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `member_id`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '会员ID',
    `plan_id`           bigint unsigned NOT NULL DEFAULT 0 COMMENT '套餐ID',
    `diagnosis_id`      bigint unsigned NOT NULL DEFAULT 0 COMMENT '关联诊断ID',
    `status`            varchar(32)     NOT NULL DEFAULT 'pending' COMMENT '状态：pending/active/cancelled',
    `amount`            decimal(10,2)   NOT NULL DEFAULT 0.00 COMMENT '订单金额',
    `signed_at`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '签约时间（Unix秒）',
    `legal_name`        varchar(64)     NOT NULL DEFAULT '' COMMENT '电子签约姓名',
    `deleted`           tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`       bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`       bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_service_order_member_id` (`member_id`),
    KEY `idx_service_order_plan_id` (`plan_id`),
    KEY `idx_service_order_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规服务订单';

-- ============================================================
-- 三、合规同意留痕
-- ============================================================
CREATE TABLE IF NOT EXISTS `xy_compliance_consent` (
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `member_id`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '会员ID',
    `order_id`          bigint unsigned NOT NULL DEFAULT 0 COMMENT '关联订单ID',
    `type`              varchar(32)     NOT NULL DEFAULT '' COMMENT '类型：terms/privacy/risk_disclosure/plan_confirm/contract_sign',
    `document_version`  varchar(32)     NOT NULL DEFAULT '' COMMENT '文档版本号',
    `ip`                varchar(45)     NOT NULL DEFAULT '' COMMENT '同意时IP',
    `user_agent`        text            NOT NULL COMMENT '同意时UA',
    `agreed_at`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '同意时间（Unix秒）',
    `create_time`       bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_compliance_consent_member_id` (`member_id`),
    KEY `idx_compliance_consent_order_id` (`order_id`),
    KEY `idx_compliance_consent_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规同意留痕';

-- ============================================================
-- 四、OPC 主体（§5.1 全量字段）
-- ============================================================
CREATE TABLE IF NOT EXISTS `xy_opc_entity` (
    `id`                      bigint unsigned NOT NULL AUTO_INCREMENT,
    `member_id`               bigint unsigned NOT NULL DEFAULT 0 COMMENT '会员ID',
    -- 拟设公司
    `proposed_names`          json            DEFAULT NULL COMMENT '备选公司名称（1~3条）',
    `registered_capital`      decimal(12,2)   NOT NULL DEFAULT 0.00 COMMENT '注册资本（万元）',
    `capital_term_years`      int             NOT NULL DEFAULT 0 COMMENT '认缴期限（年）',
    `business_term_type`      varchar(16)     NOT NULL DEFAULT '' COMMENT '营业期限：long_term/fixed',
    `business_term_end`       date            DEFAULT NULL COMMENT '固定期限截止日',
    `business_scope`          text            COMMENT '经营范围',
    `register_province`       varchar(32)     NOT NULL DEFAULT '' COMMENT '注册省',
    `register_city`           varchar(32)     NOT NULL DEFAULT '' COMMENT '注册市',
    `register_district`       varchar(32)     NOT NULL DEFAULT '' COMMENT '注册区',
    `register_address`        varchar(255)    NOT NULL DEFAULT '' COMMENT '详细注册地址',
    `address_proof_file_id`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '地址证明附件ID',
    -- 法人/股东
    `legal_person_name`       varchar(64)     NOT NULL DEFAULT '' COMMENT '法人姓名',
    `id_card_encrypted`       text            COMMENT '身份证号（AES-GCM加密）',
    `id_card_valid_from`      date            DEFAULT NULL COMMENT '身份证有效期起',
    `id_card_valid_to`        date            DEFAULT NULL COMMENT '身份证有效期止',
    `id_card_front_file_id`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '身份证正面附件ID',
    `id_card_back_file_id`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '身份证反面附件ID',
    `ethnicity`               varchar(16)     NOT NULL DEFAULT '' COMMENT '民族',
    `household_address`       varchar(255)    NOT NULL DEFAULT '' COMMENT '户籍地址',
    `residential_address`     varchar(255)    NOT NULL DEFAULT '' COMMENT '现居住地址',
    `phone`                   varchar(20)     NOT NULL DEFAULT '' COMMENT '联系手机',
    `email`                   varchar(128)    NOT NULL DEFAULT '' COMMENT '邮箱',
    `esign_authorized`        tinyint         NOT NULL DEFAULT 0 COMMENT '电子签授权：1=是 0=否',
    -- 工商结果（顾问填）
    `company_name`            varchar(128)    NOT NULL DEFAULT '' COMMENT '核准公司名称',
    `credit_code`             varchar(32)     NOT NULL DEFAULT '' COMMENT '统一社会信用代码',
    `established_at`          date            DEFAULT NULL COMMENT '成立日期',
    `license_file_id`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '营业执照附件ID',
    -- 税务
    `taxpayer_type`           varchar(32)     NOT NULL DEFAULT 'small_scale' COMMENT '纳税人类型',
    `tax_activated_at`        bigint unsigned NOT NULL DEFAULT 0 COMMENT '电子税务局激活时间（Unix秒）',
    -- 银行
    `bank_name`               varchar(64)     NOT NULL DEFAULT '' COMMENT '开户银行',
    `bank_account_enc`        text            COMMENT '对公账号（AES-GCM加密）',
    `bank_receipt_file_id`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '开户回执附件ID',
    -- 流程
    `status`                  varchar(32)     NOT NULL DEFAULT 'pending' COMMENT '状态：pending/materials/materials_review/registering/tax/bank/active',
    `materials_submitted_at`  bigint unsigned NOT NULL DEFAULT 0 COMMENT '资料首次提交时间（Unix秒）',
    `materials_approved_at`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '资料审核通过时间（Unix秒）',
    `deleted`                 tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`             bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`             bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_opc_entity_member_id` (`member_id`),
    KEY `idx_opc_entity_status` (`status`),
    KEY `idx_opc_entity_credit_code` (`credit_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='OPC 主体';

CREATE TABLE IF NOT EXISTS `xy_opc_progress_log` (
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`       bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `step`         varchar(32)     NOT NULL DEFAULT '' COMMENT '环节：materials/business/tax/bank/complete',
    `status`       varchar(32)     NOT NULL DEFAULT '' COMMENT '子状态',
    `note`         text            COMMENT '顾问备注或驳回原因',
    `operated_by`  bigint unsigned NOT NULL DEFAULT 0 COMMENT '操作管理员ID',
    `create_time`  bigint unsigned NOT NULL DEFAULT 0 COMMENT '操作时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_opc_progress_log_opc_id` (`opc_id`),
    KEY `idx_opc_progress_log_step` (`step`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='OPC 进度日志';

-- ============================================================
-- 五、收入与费用台账
-- ============================================================
CREATE TABLE IF NOT EXISTS `xy_expense_type` (
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT,
    `code`            varchar(32)     NOT NULL DEFAULT '' COMMENT '类型编码',
    `name`            varchar(64)     NOT NULL DEFAULT '' COMMENT '类型名称',
    `voucher_hint`    text            NOT NULL COMMENT '所需凭证说明',
    `compliance_tip`  varchar(255)    NOT NULL DEFAULT '' COMMENT '合规提示',
    `sort`            int             NOT NULL DEFAULT 0 COMMENT '排序',
    `status`          tinyint         NOT NULL DEFAULT 1 COMMENT '状态：1=启用 0=禁用',
    `create_time`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_expense_type_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='费用类型库';

CREATE TABLE IF NOT EXISTS `xy_income_entry` (
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`        bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `platform`      varchar(32)     NOT NULL DEFAULT '' COMMENT '平台：douyin/kuaishou/bilibili/wechat/xiaohongshu/other',
    `category`      varchar(32)     NOT NULL DEFAULT '' COMMENT '收入类型：tip/commission/ad/slot_fee/offline/other',
    `gross_amount`  decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '含税收入',
    `platform_fee`  decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '平台服务费',
    `net_amount`    decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '实收金额',
    `occurred_at`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '发生时间（Unix秒）',
    `source`        varchar(32)     NOT NULL DEFAULT 'manual' COMMENT '来源：manual/csv',
    `remark`        varchar(255)    NOT NULL DEFAULT '' COMMENT '备注',
    `deleted`       tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_income_entry_opc_id` (`opc_id`),
    KEY `idx_income_entry_occurred_at` (`occurred_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='收入台账';

CREATE TABLE IF NOT EXISTS `xy_expense_entry` (
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`          bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `category`        varchar(32)     NOT NULL DEFAULT '' COMMENT '费用类型编码',
    `amount`          decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '金额',
    `invoice_type`    varchar(32)     NOT NULL DEFAULT '' COMMENT '发票类型',
    `attachment_id`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '凭证附件ID',
    `description`     text            COMMENT '费用说明',
    `warning_flag`    tinyint         NOT NULL DEFAULT 0 COMMENT '预警标记：1=是 0=否',
    `occurred_at`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '发生时间（Unix秒）',
    `deleted`         tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_expense_entry_opc_id` (`opc_id`),
    KEY `idx_expense_entry_occurred_at` (`occurred_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='费用台账';

CREATE TABLE IF NOT EXISTS `xy_bank_transaction` (
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`        bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `occurred_at`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '发生时间（Unix秒）',
    `amount`        decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '金额',
    `description`   text            COMMENT '摘要',
    `matched`       tinyint         NOT NULL DEFAULT 0 COMMENT '是否已匹配：1=是 0=否',
    `income_id`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '匹配的收入ID',
    `deleted`       tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_bank_tx_opc_id` (`opc_id`),
    KEY `idx_bank_tx_occurred_at` (`occurred_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='银行流水';

-- ============================================================
-- 六、记账与申报
-- ============================================================
CREATE TABLE IF NOT EXISTS `xy_ledger_voucher` (
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`          bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `period`          varchar(7)      NOT NULL DEFAULT '' COMMENT '会计期间（如 2026-06）',
    `debit_account`   varchar(16)     NOT NULL DEFAULT '' COMMENT '借方科目代码',
    `credit_account`  varchar(16)     NOT NULL DEFAULT '' COMMENT '贷方科目代码',
    `amount`          decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '金额',
    `ref_type`        varchar(32)     NOT NULL DEFAULT '' COMMENT '来源类型：income/expense',
    `ref_id`          bigint unsigned NOT NULL DEFAULT 0 COMMENT '来源记录ID',
    `create_time`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_ledger_voucher_opc_period` (`opc_id`, `period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='会计凭证';

CREATE TABLE IF NOT EXISTS `xy_profit_summary` (
    `id`                  bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`              bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `year`                int             NOT NULL DEFAULT 0 COMMENT '年份',
    `month`               int             NOT NULL DEFAULT 0 COMMENT '月份',
    `revenue`             decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '收入',
    `cost`                decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '成本',
    `profit`              decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '利润',
    `cumulative_profit`   decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '累计年度利润',
    `create_time`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_profit_summary_opc_period` (`opc_id`, `year`, `month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='利润汇总';

CREATE TABLE IF NOT EXISTS `xy_tax_filing_task` (
    `id`                  bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`              bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `tax_type`            varchar(32)     NOT NULL DEFAULT '' COMMENT '税种：vat/cit_quarterly/cit_annual/surcharge',
    `period`              varchar(16)     NOT NULL DEFAULT '' COMMENT '申报期间',
    `due_date`            bigint unsigned NOT NULL DEFAULT 0 COMMENT '截止日（Unix秒）',
    `status`              varchar(16)     NOT NULL DEFAULT 'pending' COMMENT '状态：pending/filed/overdue',
    `calculated_amount`   decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '系统计算税额',
    `filed_amount`        decimal(14,2)   DEFAULT NULL COMMENT '实际申报税额',
    `receipt_file_id`     bigint unsigned NOT NULL DEFAULT 0 COMMENT '申报回执附件ID',
    `checklist`           json            DEFAULT NULL COMMENT '报税前自查清单',
    `deleted`             tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`         bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_tax_filing_task_opc_period` (`opc_id`, `period`),
    KEY `idx_tax_filing_task_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='税务申报任务';

-- ============================================================
-- 七、对账单与审计
-- ============================================================
CREATE TABLE IF NOT EXISTS `xy_monthly_statement` (
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`        bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `year`          int             NOT NULL DEFAULT 0 COMMENT '年份',
    `month`         int             NOT NULL DEFAULT 0 COMMENT '月份',
    `summary`       json            DEFAULT NULL COMMENT '对账单摘要',
    `sent_at`       bigint unsigned NOT NULL DEFAULT 0 COMMENT '发送时间（Unix秒）',
    `deleted`       tinyint         NOT NULL DEFAULT 0 COMMENT '软删除：1=是 0=否',
    `create_time`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_monthly_statement_opc_period` (`opc_id`, `year`, `month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='月度对账单';

CREATE TABLE IF NOT EXISTS `xy_compliance_audit_log` (
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT,
    `entity_type`    varchar(32)     NOT NULL DEFAULT '' COMMENT '实体类型',
    `entity_id`      bigint unsigned NOT NULL DEFAULT 0 COMMENT '实体ID',
    `action`         varchar(64)     NOT NULL DEFAULT '' COMMENT '操作动作',
    `operator_id`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '操作人ID',
    `operator_type`  varchar(16)     NOT NULL DEFAULT '' COMMENT '操作人类型：member/admin',
    `before`         json            DEFAULT NULL COMMENT '变更前快照',
    `after`          json            DEFAULT NULL COMMENT '变更后快照',
    `ip`             varchar(45)     NOT NULL DEFAULT '' COMMENT '操作IP',
    `create_time`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '操作时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_compliance_audit_entity` (`entity_type`, `entity_id`),
    KEY `idx_compliance_audit_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规审计日志';

-- ============================================================
-- 八、种子数据
-- ============================================================
INSERT IGNORE INTO `xy_service_plan` (`name`, `tier`, `monthly_price`, `price_display`, `features`, `sort`, `status`, `create_time`, `update_time`)
VALUES
('基础套餐', 'basic', 299.00, '¥299/月', '["合规诊断", "OPC设立代办", "月度记账", "申报提醒"]', 10, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('进阶套餐', 'advanced', 799.00, '¥799/月', '["基础套餐全部", "专属顾问", "银行开户协助", "月度对账单"]', 20, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('尊享套餐', 'premium', 2999.00, '面议', '["进阶套餐全部", "多平台收入对账", "税务筹划咨询", "优先申报处理"]', 30, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

INSERT IGNORE INTO `xy_expense_type` (`code`, `name`, `voucher_hint`, `compliance_tip`, `sort`, `status`, `create_time`, `update_time`)
VALUES
('equipment', '设备', '购置设备的发票（普票/专票）', '设备采购须与实际经营相关', 10, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('network', '网费', '网络服务发票', '家庭宽带不可全额列支，须按经营比例分摊', 20, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('venue', '场地', '租赁合同 + 租金发票', '租赁地址须与注册地址或实际经营地一致', 30, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('promotion', '投流', '平台推广费发票（如 dou+）', '投流费用须与直播推广直接相关', 40, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('outsource', '外包', '外包服务合同 + 发票', '外包内容须为真实业务，禁止虚构咨询费', 50, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('travel', '差旅', '车票/机票/住宿发票', '差旅须与业务活动直接相关', 60, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('other', '其他', '费用说明 + 相关凭证', '其他费用须有合理商业目的', 70, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- ============================================================
-- 九、后台菜单（合规模块）
-- ============================================================
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT 0, 1, '合规服务', 'Compliance', '/compliance', '/index/index', '', 'ri:shield-check-line', 0, 0, '', '', '', 0, 0, 0, '', '', 0, 0, 90, 1, 'OPC合规服务管理', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'Compliance' AND t.type = 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 2, '合规客户', 'ComplianceCustomers', 'customers', '/compliance/customers/index', 'compliance_customers', 'ri:team-line', 0, 1, '', '', '["GET /admin/compliance/customers"]', 0, 0, 0, '', '', 0, 0, 10, 1, '合规客户列表', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'Compliance' AND p.type = 1
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceCustomers' AND t.type = 2);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 2, 'OPC任务', 'ComplianceOpcTasks', 'opc-tasks', '/compliance/opc-tasks/index', 'compliance_opc_tasks', 'ri:file-list-3-line', 0, 1, '', '', '["GET /admin/compliance/opc-tasks"]', 0, 0, 0, '', '', 0, 0, 20, 1, 'OPC落地任务队列', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'Compliance' AND p.type = 1
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceOpcTasks' AND t.type = 2);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 2, '申报工作台', 'ComplianceFiling', 'filing', '/compliance/filing/index', 'compliance_filing', 'ri:calendar-check-line', 0, 1, '', '', '["GET /admin/compliance/filing"]', 0, 0, 0, '', '', 0, 0, 30, 1, '税务申报任务处理', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'Compliance' AND p.type = 1
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceFiling' AND t.type = 2);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 2, '对账单管理', 'ComplianceStatements', 'statements', '/compliance/statements/index', 'compliance_statements', 'ri:file-chart-line', 0, 1, '', '', '["GET /admin/compliance/statements"]', 0, 0, 0, '', '', 0, 0, 40, 1, '月度对账单生成与发送', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'Compliance' AND p.type = 1
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceStatements' AND t.type = 2);

-- OPC 任务按钮
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 3, '查看详情', 'view', '', '', 'compliance_opc_tasks', '', 0, 0, '', '', '["GET /admin/compliance/opc-tasks"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'ComplianceOpcTasks' AND p.type = 2
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` sub WHERE sub.parent_id = p.id AND sub.type = 3 AND sub.name = 'view');

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 3, '推进进度', 'advance', '', '', 'compliance_opc_tasks', '', 0, 0, '', '', '["PATCH /admin/compliance/opc-tasks"]', 0, 0, 0, '', '', 0, 0, 2, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'ComplianceOpcTasks' AND p.type = 2
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` sub WHERE sub.parent_id = p.id AND sub.type = 3 AND sub.name = 'advance');

-- 申报工作台按钮
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 3, '标记已申报', 'filed', '', '', 'compliance_filing', '', 0, 0, '', '', '["PATCH /admin/compliance/filing"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'ComplianceFiling' AND p.type = 2
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` sub WHERE sub.parent_id = p.id AND sub.type = 3 AND sub.name = 'filed');

-- 对账单按钮
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 3, '生成发送', 'send', '', '', 'compliance_statements', '', 0, 0, '', '', '["POST /admin/compliance/statements/send"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'ComplianceStatements' AND p.type = 2
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` sub WHERE sub.parent_id = p.id AND sub.type = 3 AND sub.name = 'send');

-- 合规客户按钮
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT p.id, 3, '导出留痕', 'export', '', '', 'compliance_customers', '', 0, 0, '', '', '["GET /admin/compliance/audit/export"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `xy_admin_menu` p WHERE p.name = 'ComplianceCustomers' AND p.type = 2
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` sub WHERE sub.parent_id = p.id AND sub.type = 3 AND sub.name = 'export');
