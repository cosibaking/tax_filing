-- 会员报税 Excel 申报提交记录
CREATE TABLE IF NOT EXISTS `xy_tax_filing_submission` (
    `id`               bigint unsigned NOT NULL AUTO_INCREMENT,
    `opc_id`           bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `period`           varchar(16)     NOT NULL DEFAULT '' COMMENT '申报期间（如 2026-06）',
    `revenue`          decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '营业收入',
    `expense_total`    decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '成本费用合计',
    `profit`           decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '利润总额',
    `vat_amount`       decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '增值税',
    `cit_amount`       decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '企业所得税',
    `surcharge_amount` decimal(14,2)   NOT NULL DEFAULT 0.00 COMMENT '附加税',
    `remark`           text            COMMENT '备注',
    `expense_imported` int             NOT NULL DEFAULT 0 COMMENT '导入报销条数',
    `excel_file_id`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '上传的Excel附件ID',
    `status`           varchar(16)     NOT NULL DEFAULT 'submitted' COMMENT '状态：submitted',
    `deleted`          tinyint         NOT NULL DEFAULT 0 COMMENT '软删除',
    `create_time`      bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`      bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_tax_submission_opc_period` (`opc_id`, `period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='会员报税Excel申报提交';
