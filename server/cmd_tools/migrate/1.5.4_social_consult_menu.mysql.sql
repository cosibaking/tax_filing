-- Migration: 1.5.4
-- Description: 确保顾问端「社保咨询」菜单在合规服务目录下，并启用小红点

SET @complianceDirId = (SELECT id FROM `xy_admin_menu` WHERE name = 'Compliance' AND type = 1 LIMIT 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @complianceDirId, 2, '社保咨询', 'ComplianceSocialConsults', 'social-consults', '/compliance/social-consults/index', 'compliance_social_consults', 'ri:question-answer-line', 0, 1, '', '', '["GET /admin/compliance/social-consults"]', 0, 0, 1, '', '', 0, 0, 35, 1, '社保咨询工单处理', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE @complianceDirId IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.name = 'ComplianceSocialConsults' AND t.type = 2);

UPDATE `xy_admin_menu`
SET `parent_id` = @complianceDirId,
    `show_badge` = 1,
    `sort` = 35,
    `status` = 1,
    `update_time` = UNIX_TIMESTAMP()
WHERE `name` = 'ComplianceSocialConsults'
  AND `type` = 2
  AND @complianceDirId IS NOT NULL;

SET @socialConsultsPageId = (SELECT id FROM `xy_admin_menu` WHERE name = 'ComplianceSocialConsults' AND type = 2 LIMIT 1);

INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `resource`, `icon`, `hidden`, `keep_alive`, `redirect`, `frame_src`, `perms`, `is_frame`, `affix`, `show_badge`, `badge_text`, `active_path`, `hide_tab`, `is_full_page`, `sort`, `status`, `remark`, `created_by`, `updated_by`, `create_time`, `update_time`)
SELECT @socialConsultsPageId, 3, '回复咨询', 'reply', '', '', 'compliance_social_consults', '', 0, 0, '', '', '["PATCH /admin/compliance/social-consults"]', 0, 0, 0, '', '', 0, 0, 1, 1, '', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE @socialConsultsPageId IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM `xy_admin_menu` t WHERE t.parent_id = @socialConsultsPageId AND t.type = 3 AND t.name = 'reply');
