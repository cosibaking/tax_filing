-- Migration: 1.5.3
-- Description: 修复基础表 id 缺少 AUTO_INCREMENT（导致注册、菜单插入等失败）

ALTER TABLE `xy_admin_chat_message` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_chat_session` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_chat_session_member` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_dept` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_dept_closure` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_field_perm` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_login_log` MODIFY `id` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_menu` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_notice` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_notice_read` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_operation_log` MODIFY `id` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_post` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_role` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_admin_user` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_captcha` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_cms_changelog` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_cms_doc` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_cms_doc_category` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_demo_article` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_demo_category` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '会员ID';
ALTER TABLE `xy_member_checkin` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_group` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_login_log` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_menu` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_money_log` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_notice` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_notice_read` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_oauth` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_member_score_log` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sms_log` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sms_template` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sms_variable` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sys_attachment` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sys_config` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sys_cron` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sys_cron_group` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sys_cron_log` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sys_gen_codes` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_sys_gen_codes_column` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_test_category` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_test_code` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_test_codec` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;
ALTER TABLE `xy_test_order` MODIFY `id` bigint unsigned NOT NULL AUTO_INCREMENT;

-- 补充 xy_member 唯一索引（若已存在则跳过）
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'xy_member' AND index_name = 'uk_username');
SET @sql := IF(@exist = 0, 'ALTER TABLE `xy_member` ADD UNIQUE KEY `uk_username` (`username`)', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @exist := (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'xy_member' AND index_name = 'uk_mobile');
SET @sql := IF(@exist = 0, 'ALTER TABLE `xy_member` ADD UNIQUE KEY `uk_mobile` (`mobile`)', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
