-- 1.5.2 收入 OCR / MCN 分成 / 平台报送快照
ALTER TABLE `xy_income_entry`
    ADD COLUMN `settlement_type` varchar(32) NOT NULL DEFAULT 'personal' COMMENT '结算方式：personal/mcn_public' AFTER `source`,
    ADD COLUMN `mcn_name` varchar(128) NOT NULL DEFAULT '' COMMENT 'MCN 机构名称' AFTER `settlement_type`,
    ADD COLUMN `mcn_split_ratio` decimal(5,4) NOT NULL DEFAULT 0.0000 COMMENT 'MCN 分成比例 0-1' AFTER `mcn_name`,
    ADD COLUMN `mcn_share_amount` decimal(14,2) NOT NULL DEFAULT 0.00 COMMENT 'MCN 分成金额' AFTER `mcn_split_ratio`,
    ADD COLUMN `gross_before_split` decimal(14,2) NOT NULL DEFAULT 0.00 COMMENT '分成前平台结算额' AFTER `mcn_share_amount`,
    ADD COLUMN `attachment_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OCR 截图附件ID' AFTER `gross_before_split`;

CREATE TABLE IF NOT EXISTS `xy_platform_income_snapshot` (
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`        bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `platform`      varchar(32)     NOT NULL DEFAULT '' COMMENT '平台',
    `period`        varchar(7)      NOT NULL DEFAULT '' COMMENT '月份 YYYY-MM',
    `gross_amount`  decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '平台报送/识别含税收入',
    `platform_fee`  decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '平台服务费',
    `net_amount`    decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '实收',
    `source`        varchar(32)     NOT NULL DEFAULT 'ocr' COMMENT '来源：ocr/manual',
    `attachment_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '凭证附件',
    `remark`        varchar(255)    NOT NULL DEFAULT '',
    `create_time`   bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_platform_snapshot_opc_period` (`opc_id`, `period`),
    KEY `idx_platform_snapshot_platform` (`platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='平台收入报送快照（OCR/手工）';
