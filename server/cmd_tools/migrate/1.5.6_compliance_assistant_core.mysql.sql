-- OPC 公司经营合规助手核心表。
-- 回滚策略：关闭 complianceAssistant 功能开关并保留数据；禁止在线 DROP/TRUNCATE。

CREATE TABLE IF NOT EXISTS `xy_enterprise_profile` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `opc_entity_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'OPC主体ID',
  `version` int unsigned NOT NULL DEFAULT 1 COMMENT '画像版本',
  `profile_json` json NOT NULL COMMENT '标准企业画像',
  `source` varchar(32) NOT NULL DEFAULT 'user' COMMENT 'user/advisor/import',
  `status` varchar(16) NOT NULL DEFAULT 'draft' COMMENT 'draft/active/archived',
  `change_reason` varchar(500) NOT NULL DEFAULT '',
  `confirmed_by` bigint unsigned NOT NULL DEFAULT 0,
  `confirmed_at` bigint unsigned NOT NULL DEFAULT 0,
  `deleted` tinyint NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_profile_version` (`opc_entity_id`,`version`),
  KEY `idx_profile_active` (`opc_entity_id`,`status`,`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='企业画像版本';

CREATE TABLE IF NOT EXISTS `xy_compliance_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `rule_code` varchar(64) NOT NULL DEFAULT '',
  `rule_type` varchar(20) NOT NULL DEFAULT 'task' COMMENT 'task/risk/document/reminder',
  `name` varchar(128) NOT NULL DEFAULT '',
  `enabled` tinyint NOT NULL DEFAULT 1,
  `deleted` tinyint NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rule_code` (`rule_code`),
  KEY `idx_rule_type_enabled` (`rule_type`,`enabled`,`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规规则定义';

CREATE TABLE IF NOT EXISTS `xy_compliance_rule_version` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `rule_id` bigint unsigned NOT NULL DEFAULT 0,
  `version` int unsigned NOT NULL DEFAULT 1,
  `region_code` varchar(16) NOT NULL DEFAULT 'CN-BJ',
  `condition_json` json NOT NULL,
  `action_json` json NOT NULL,
  `status` varchar(16) NOT NULL DEFAULT 'draft' COMMENT 'draft/reviewing/published/retired',
  `effective_from` bigint unsigned NOT NULL DEFAULT 0,
  `effective_to` bigint unsigned NOT NULL DEFAULT 0,
  `created_by` bigint unsigned NOT NULL DEFAULT 0,
  `reviewed_by` bigint unsigned NOT NULL DEFAULT 0,
  `published_at` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rule_version` (`rule_id`,`version`),
  KEY `idx_rule_version_active` (`rule_id`,`status`,`effective_from`,`effective_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='不可变合规规则版本';

CREATE TABLE IF NOT EXISTS `xy_enterprise_compliance_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `opc_entity_id` bigint unsigned NOT NULL DEFAULT 0,
  `member_id` bigint unsigned NOT NULL DEFAULT 0,
  `rule_version_id` bigint unsigned NOT NULL DEFAULT 0,
  `risk_code` varchar(64) NOT NULL DEFAULT '',
  `period_key` varchar(16) NOT NULL DEFAULT '',
  `trigger_key` varchar(64) NOT NULL DEFAULT 'periodic',
  `task_type` varchar(32) NOT NULL DEFAULT 'general',
  `title` varchar(160) NOT NULL DEFAULT '',
  `description` text,
  `snapshot_json` json NOT NULL,
  `material_json` json NOT NULL,
  `status` varchar(24) NOT NULL DEFAULT 'not_started',
  `priority` varchar(16) NOT NULL DEFAULT 'medium',
  `due_at` bigint unsigned NOT NULL DEFAULT 0,
  `completed_at` bigint unsigned NOT NULL DEFAULT 0,
  `completion_note` varchar(1000) NOT NULL DEFAULT '',
  `deleted` tinyint NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_idempotency` (`opc_entity_id`,`rule_version_id`,`period_key`,`trigger_key`),
  KEY `idx_task_member_status_due` (`member_id`,`status`,`due_at`),
  KEY `idx_task_opc_period` (`opc_entity_id`,`period_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='企业通用合规任务';

CREATE TABLE IF NOT EXISTS `xy_task_reminder` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_id` bigint unsigned NOT NULL DEFAULT 0,
  `channel` varchar(20) NOT NULL DEFAULT 'in_app',
  `planned_at` bigint unsigned NOT NULL DEFAULT 0,
  `status` varchar(24) NOT NULL DEFAULT 'pending',
  `attempts` int unsigned NOT NULL DEFAULT 0,
  `last_error` varchar(1000) NOT NULL DEFAULT '',
  `sent_at` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reminder_idempotency` (`task_id`,`channel`,`planned_at`),
  KEY `idx_reminder_dispatch` (`status`,`planned_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规任务提醒';

CREATE TABLE IF NOT EXISTS `xy_business_document` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `opc_entity_id` bigint unsigned NOT NULL DEFAULT 0,
  `member_id` bigint unsigned NOT NULL DEFAULT 0,
  `attachment_id` bigint unsigned NOT NULL DEFAULT 0,
  `period_key` varchar(16) NOT NULL DEFAULT '',
  `document_type` varchar(32) NOT NULL DEFAULT 'unknown',
  `process_status` varchar(24) NOT NULL DEFAULT 'uploaded',
  `file_hash` char(64) NOT NULL DEFAULT '',
  `extracted_json` json,
  `confidence` decimal(5,4) NOT NULL DEFAULT 0,
  `confirmed_by` bigint unsigned NOT NULL DEFAULT 0,
  `confirmed_at` bigint unsigned NOT NULL DEFAULT 0,
  `deleted` tinyint NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_document_opc_period` (`opc_entity_id`,`period_key`,`document_type`),
  KEY `idx_document_hash` (`opc_entity_id`,`file_hash`),
  KEY `idx_document_review` (`process_status`,`confidence`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='经营资料索引';

CREATE TABLE IF NOT EXISTS `xy_document_relation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint unsigned NOT NULL DEFAULT 0,
  `relation_type` varchar(24) NOT NULL DEFAULT '',
  `relation_id` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_document_relation` (`document_id`,`relation_type`,`relation_id`),
  KEY `idx_relation_target` (`relation_type`,`relation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='经营资料业务关联';

CREATE TABLE IF NOT EXISTS `xy_risk_event` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `opc_entity_id` bigint unsigned NOT NULL DEFAULT 0,
  `member_id` bigint unsigned NOT NULL DEFAULT 0,
  `rule_version_id` bigint unsigned NOT NULL DEFAULT 0,
  `period_key` varchar(16) NOT NULL DEFAULT '',
  `severity` varchar(16) NOT NULL DEFAULT 'medium',
  `status` varchar(16) NOT NULL DEFAULT 'open',
  `summary` varchar(500) NOT NULL DEFAULT '',
  `evidence_json` json NOT NULL,
  `resolution_note` varchar(1000) NOT NULL DEFAULT '',
  `first_hit_at` bigint unsigned NOT NULL DEFAULT 0,
  `last_hit_at` bigint unsigned NOT NULL DEFAULT 0,
  `resolved_at` bigint unsigned NOT NULL DEFAULT 0,
  `deleted` tinyint NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_risk_period` (`opc_entity_id`,`risk_code`,`period_key`),
  KEY `idx_risk_member_status` (`member_id`,`status`,`severity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='合规风险事件';

CREATE TABLE IF NOT EXISTS `xy_compliance_monthly_report` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `opc_entity_id` bigint unsigned NOT NULL DEFAULT 0,
  `member_id` bigint unsigned NOT NULL DEFAULT 0,
  `period_key` varchar(16) NOT NULL DEFAULT '',
  `version` int unsigned NOT NULL DEFAULT 1,
  `status` varchar(16) NOT NULL DEFAULT 'draft',
  `statistics_json` json NOT NULL,
  `completeness_json` json NOT NULL,
  `risk_snapshot_json` json NOT NULL,
  `content` mediumtext,
  `rule_versions_json` json NOT NULL,
  `ai_model` varchar(64) NOT NULL DEFAULT '',
  `knowledge_version` varchar(64) NOT NULL DEFAULT '',
  `published_by` bigint unsigned NOT NULL DEFAULT 0,
  `published_at` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_report_version` (`opc_entity_id`,`period_key`,`version`),
  KEY `idx_report_member_period` (`member_id`,`period_key`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='月度经营合规体检';

CREATE TABLE IF NOT EXISTS `xy_service_ticket` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `opc_entity_id` bigint unsigned NOT NULL DEFAULT 0,
  `member_id` bigint unsigned NOT NULL DEFAULT 0,
  `ticket_type` varchar(32) NOT NULL DEFAULT 'compliance_review',
  `status` varchar(24) NOT NULL DEFAULT 'submitted',
  `title` varchar(160) NOT NULL DEFAULT '',
  `description` text,
  `task_id` bigint unsigned NOT NULL DEFAULT 0,
  `risk_event_id` bigint unsigned NOT NULL DEFAULT 0,
  `assignee_id` bigint unsigned NOT NULL DEFAULT 0,
  `provider_name` varchar(160) NOT NULL DEFAULT '',
  `sla_due_at` bigint unsigned NOT NULL DEFAULT 0,
  `completed_at` bigint unsigned NOT NULL DEFAULT 0,
  `deleted` tinyint NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  `update_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_member_status` (`member_id`,`status`),
  KEY `idx_ticket_assignee_status` (`assignee_id`,`status`,`sla_due_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='通用合规服务工单';

CREATE TABLE IF NOT EXISTS `xy_service_ticket_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint unsigned NOT NULL DEFAULT 0,
  `from_status` varchar(24) NOT NULL DEFAULT '',
  `to_status` varchar(24) NOT NULL DEFAULT '',
  `content` text,
  `attachment_id` bigint unsigned NOT NULL DEFAULT 0,
  `operator_type` varchar(16) NOT NULL DEFAULT 'member',
  `operator_id` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_record` (`ticket_id`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='工单不可变处理记录';

CREATE TABLE IF NOT EXISTS `xy_ai_audit_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `opc_entity_id` bigint unsigned NOT NULL DEFAULT 0,
  `scene` varchar(32) NOT NULL DEFAULT '',
  `prompt_version` varchar(32) NOT NULL DEFAULT '',
  `model_name` varchar(64) NOT NULL DEFAULT '',
  `input_digest` varchar(128) NOT NULL DEFAULT '',
  `output_digest` varchar(128) NOT NULL DEFAULT '',
  `duration_ms` int unsigned NOT NULL DEFAULT 0,
  `status` varchar(16) NOT NULL DEFAULT 'success',
  `error_message` varchar(500) NOT NULL DEFAULT '',
  `confirmed_by` bigint unsigned NOT NULL DEFAULT 0,
  `confirmed_at` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ai_audit_opc_scene` (`opc_entity_id`,`scene`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI调用审计摘要';
