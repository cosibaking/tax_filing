-- Migration: 1.4.2
-- Description: 禁用 Art Design Pro 演示菜单，保留 SaaS 基础功能与合规服务

-- 1. 仪表盘仅保留工作台
UPDATE `xy_admin_menu`
SET `status` = 0, `update_time` = UNIX_TIMESTAMP()
WHERE `name` IN ('Analysis', 'Ecommerce');

-- 2. 禁用演示/模板类顶级目录
UPDATE `xy_admin_menu`
SET `status` = 0, `update_time` = UNIX_TIMESTAMP()
WHERE `name` IN (
    'Template', 'Widgets', 'Examples', 'Cms',
    'Result', 'Exception', 'Develop', 'ChangeLog', 'Nested'
);

-- 3. 级联禁用已禁用菜单的子项（最多 4 层）
UPDATE `xy_admin_menu` c
INNER JOIN `xy_admin_menu` p ON c.parent_id = p.id
SET c.status = 0, c.update_time = UNIX_TIMESTAMP()
WHERE p.status = 0 AND c.status = 1;

UPDATE `xy_admin_menu` c
INNER JOIN `xy_admin_menu` p ON c.parent_id = p.id
SET c.status = 0, c.update_time = UNIX_TIMESTAMP()
WHERE p.status = 0 AND c.status = 1;

UPDATE `xy_admin_menu` c
INNER JOIN `xy_admin_menu` p ON c.parent_id = p.id
SET c.status = 0, c.update_time = UNIX_TIMESTAMP()
WHERE p.status = 0 AND c.status = 1;

UPDATE `xy_admin_menu` c
INNER JOIN `xy_admin_menu` p ON c.parent_id = p.id
SET c.status = 0, c.update_time = UNIX_TIMESTAMP()
WHERE p.status = 0 AND c.status = 1;

-- 4. 对账单管理补充按钮权限（列表/批量生成/发送）
SET @statementsPageId = (SELECT id FROM `xy_admin_menu` WHERE name = 'ComplianceStatements' AND type = 2 LIMIT 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @statementsPageId, 3, '批量生成', 'generate', '', '', 'compliance_statements', '', 0, 0, '', '', '["POST /admin/compliance/statements/generate"]', 0, 0, 0, '', '', 0, 0, 2, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE @statementsPageId IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.parent_id = @statementsPageId AND t.type = 3 AND t.name = 'generate');
