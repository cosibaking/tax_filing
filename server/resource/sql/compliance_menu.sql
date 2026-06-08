-- 合规服务模块 菜单 SQL
-- ======= 挂载模式：创建顶级目录(type=1) + 页面(type=2) + 按钮(type=3) =======

-- 1. 创建顶级目录
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT 0, 1, '合规服务', 'Compliance', '/compliance', '/index/index', '', 'ri:shield-check-line', 0, 0, '', '', '', 0, 0, 0, '', '', 0, 0, 90, 1, 'OPC合规服务管理', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'Compliance' AND t.type = 1);

SET @complianceDirId = (SELECT id FROM `xy_admin_menu` WHERE name = 'Compliance' AND type = 1 LIMIT 1);

-- 2. 合规客户
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @complianceDirId, 2, '合规客户', 'ComplianceCustomers', 'customers', '/compliance/customers/index', 'compliance_customers', 'ri:team-line', 0, 1, '', '', '["GET /admin/compliance/customers"]', 0, 0, 0, '', '', 0, 0, 10, 1, '合规客户列表', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceCustomers' AND t.type = 2);

SET @customersPageId = (SELECT id FROM `xy_admin_menu` WHERE name = 'ComplianceCustomers' AND type = 2 LIMIT 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @customersPageId, 3, '导出留痕', 'ComplianceCustomersExport', '', '', 'compliance_customers', '', 0, 0, '', '', '["GET /admin/compliance/audit/export"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.parent_id = @customersPageId AND t.type = 3 AND t.name = 'ComplianceCustomersExport');

-- 3. OPC 任务
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @complianceDirId, 2, 'OPC任务', 'ComplianceOpcTasks', 'opc-tasks', '/compliance/opc-tasks/index', 'compliance_opc_tasks', 'ri:file-list-3-line', 0, 1, '', '', '["GET /admin/compliance/opc-tasks"]', 0, 0, 0, '', '', 0, 0, 20, 1, 'OPC落地任务队列', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceOpcTasks' AND t.type = 2);

SET @opcTasksPageId = (SELECT id FROM `xy_admin_menu` WHERE name = 'ComplianceOpcTasks' AND type = 2 LIMIT 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @opcTasksPageId, 3, '查看详情', 'ComplianceOpcTasksView', '', '', 'compliance_opc_tasks', '', 0, 0, '', '', '["GET /admin/compliance/opc-tasks"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.parent_id = @opcTasksPageId AND t.type = 3 AND t.name = 'ComplianceOpcTasksView');

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @opcTasksPageId, 3, '推进进度', 'ComplianceOpcTasksAdvance', '', '', 'compliance_opc_tasks', '', 0, 0, '', '', '["PATCH /admin/compliance/opc-tasks"]', 0, 0, 0, '', '', 0, 0, 2, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.parent_id = @opcTasksPageId AND t.type = 3 AND t.name = 'ComplianceOpcTasksAdvance');

-- 4. 申报工作台
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @complianceDirId, 2, '申报工作台', 'ComplianceFiling', 'filing', '/compliance/filing/index', 'compliance_filing', 'ri:calendar-check-line', 0, 1, '', '', '["GET /admin/compliance/filing"]', 0, 0, 0, '', '', 0, 0, 30, 1, '税务申报任务处理', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceFiling' AND t.type = 2);

SET @filingPageId = (SELECT id FROM `xy_admin_menu` WHERE name = 'ComplianceFiling' AND type = 2 LIMIT 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @filingPageId, 3, '标记已申报', 'ComplianceFilingFiled', '', '', 'compliance_filing', '', 0, 0, '', '', '["PATCH /admin/compliance/filing"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.parent_id = @filingPageId AND t.type = 3 AND t.name = 'ComplianceFilingFiled');

-- 5. 对账单管理
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @complianceDirId, 2, '对账单管理', 'ComplianceStatements', 'statements', '/compliance/statements/index', 'compliance_statements', 'ri:file-chart-line', 0, 1, '', '', '["GET /admin/compliance/statements"]', 0, 0, 0, '', '', 0, 0, 40, 1, '月度对账单生成与发送', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceStatements' AND t.type = 2);

SET @statementsPageId = (SELECT id FROM `xy_admin_menu` WHERE name = 'ComplianceStatements' AND type = 2 LIMIT 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @statementsPageId, 3, '生成发送', 'ComplianceStatementsSend', '', '', 'compliance_statements', '', 0, 0, '', '', '["POST /admin/compliance/statements/send"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.parent_id = @statementsPageId AND t.type = 3 AND t.name = 'ComplianceStatementsSend');
