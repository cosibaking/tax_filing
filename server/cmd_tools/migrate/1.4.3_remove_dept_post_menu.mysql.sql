-- Migration: 1.4.3
-- Description: 移除权限管理下的部门/岗位菜单（MVP 不需要）

UPDATE `xy_admin_menu`
SET `status` = 0, `update_time` = UNIX_TIMESTAMP()
WHERE `name` IN ('Dept', 'Post') AND `type` = 2;

-- 级联禁用按钮权限
UPDATE `xy_admin_menu` c
INNER JOIN `xy_admin_menu` p ON c.parent_id = p.id
SET c.status = 0, c.update_time = UNIX_TIMESTAMP()
WHERE p.name IN ('Dept', 'Post') AND c.status = 1;
