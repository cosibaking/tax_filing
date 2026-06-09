-- M9 Phase 2：用工状态 + 社保指引 + 咨询工单

-- 1. OPC 主体扩展用工状态
ALTER TABLE `xy_opc_entity`
  ADD COLUMN `employment_status` varchar(20) NOT NULL DEFAULT 'unknown'
    COMMENT '用工状态：unknown/no_employee/has_employee' AFTER `status`,
  ADD COLUMN `employment_confirmed_at` bigint unsigned NOT NULL DEFAULT 0
    COMMENT '用工状态确认时间（Unix秒）' AFTER `employment_status`;

-- 2. 社保咨询工单
CREATE TABLE IF NOT EXISTS `xy_compliance_social_consult` (
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `member_id`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '会员ID',
    `opc_id`       bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
    `category`     varchar(32)     NOT NULL DEFAULT 'founder' COMMENT 'founder/employee/other',
    `question`     text            NOT NULL COMMENT '咨询问题',
    `region_code`  varchar(12)     NOT NULL DEFAULT '' COMMENT '属地编码',
    `status`       varchar(20)     NOT NULL DEFAULT 'open' COMMENT 'open/replied/closed',
    `reply`        text            COMMENT '顾问回复',
    `replied_by`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '回复顾问ID',
    `replied_at`   bigint unsigned NOT NULL DEFAULT 0 COMMENT '回复时间（Unix秒）',
    `closed_at`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '关闭时间（Unix秒）',
    `deleted`      tinyint         NOT NULL DEFAULT 0 COMMENT '软删除',
    `create_time`  bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（Unix秒）',
    `update_time`  bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（Unix秒）',
    PRIMARY KEY (`id`),
    KEY `idx_social_consult_member` (`member_id`),
    KEY `idx_social_consult_opc_status` (`opc_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='社保咨询工单';

-- 3. 社保指引文章
CREATE TABLE IF NOT EXISTS `xy_compliance_social_guide` (
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `slug`        varchar(64)     NOT NULL DEFAULT '' COMMENT '文章标识',
    `title`       varchar(128)    NOT NULL DEFAULT '' COMMENT '标题',
    `summary`     varchar(512)    NOT NULL DEFAULT '' COMMENT '摘要',
    `content_md`  mediumtext      NOT NULL COMMENT 'Markdown正文',
    `audience`    varchar(32)     NOT NULL DEFAULT 'founder' COMMENT 'founder/employer',
    `sort`        int             NOT NULL DEFAULT 0 COMMENT '排序',
    `enabled`     tinyint         NOT NULL DEFAULT 1 COMMENT '是否启用',
    `create_time` bigint unsigned NOT NULL DEFAULT 0,
    `update_time` bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_social_guide_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='社保指引文章';

-- 4. 指引种子数据
INSERT INTO `xy_compliance_social_guide`
    (`slug`, `title`, `summary`, `content_md`, `audience`, `sort`, `enabled`, `create_time`, `update_time`)
SELECT * FROM (
    SELECT
        'flexible_employment' AS slug,
        '灵活就业参保（推荐全职主播）' AS title,
        '一人 OPC 不强制为法人缴纳职工社保，全职主播可通过灵活就业身份缴纳养老、医疗保险。' AS summary,
        '## 什么是灵活就业参保\n\n灵活就业人员是指未与用人单位建立劳动关系、以个人身份参加职工基本养老保险和职工基本医疗保险的人员。\n\n## 适合谁\n\n- 已注册 OPC 或个体户，**无其他雇员**\n- 全职从事直播经营，需要保持养老、医疗连续性\n\n## 办理路径（一般流程）\n\n1. 携带身份证到户籍地或经营地社保经办机构窗口办理；\n2. 或通过当地「掌上12333」、政务服务网等线上渠道登记；\n3. 选择缴费基数档次，按月自行缴纳。\n\n## 注意事项\n\n- 各地缴费基数上下限、比例不同，请以当地社保局公布为准；\n- 灵活就业一般只含养老+医疗，不含工伤、失业、生育（以当地政策为准）；\n- 如有疑问，可提交「社保咨询」由顾问结合您的注册地答复。' AS content_md,
        'founder' AS audience,
        10 AS sort,
        1 AS enabled,
        UNIX_TIMESTAMP() AS create_time,
        UNIX_TIMESTAMP() AS update_time
    UNION ALL
    SELECT
        'founder_vs_employee',
        '创始人参保 vs 雇员参保',
        '区分法人本人社保与员工社保，避免混淆申报义务。',
        '## 两类场景\n\n### A. 创始人/法人本人\n\n一人有限责任公司**不强制**为唯一股东兼法人缴纳职工社保。建议通过灵活就业或城乡居民医保保持保障连续性。\n\n### B. 真实雇员\n\n若 OPC 聘用助理、剪辑等员工并签订劳动合同，须依法：\n\n1. 按月发放工资并入账；\n2. 通过自然人扣缴端申报工资个税；\n3. 在社保系统为员工申报五险；\n4. 在公积金中心申报住房公积金（如适用）。\n\n## 合规提醒\n\n- 不得通过虚设员工、虚假工资套取社保优惠；\n- 不得将打赏收入伪装为工资或分红以规避申报义务。',
        'founder',
        20,
        1,
        UNIX_TIMESTAMP(),
        UNIX_TIMESTAMP()
    UNION ALL
    SELECT
        'faq_q3',
        '常见问题：注册了公司还要交社保吗？',
        '源自实操手册 Q3 的标准答复。',
        '## 问：我当 OPC，还要交社保吗？\n\n**答：** 一人 OPC 不强制为法人缴纳职工社保。若您全职从事直播经营，建议通过**灵活就业**或**城乡居民医保**等方式保持养老、医疗连续性。\n\n若您聘用员工，须依法为员工缴纳社保，并按时申报工资个税。\n\n> 本内容仅供参考，具体以当地社保、税务部门规定为准。',
        'founder',
        30,
        1,
        UNIX_TIMESTAMP(),
        UNIX_TIMESTAMP()
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM `xy_compliance_social_guide` LIMIT 1);
