-- Migration: 1.4.1
-- Description: OPC 合规 Cron 任务种子（M6）

INSERT INTO `xy_sys_cron` (`group_id`, `title`, `name`, `params`, `pattern`, `policy`, `count`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
SELECT 1, '合规材料催收', 'compliance_material_reminder', '', '0 0 9 1-5 * *', 2, 1, 50, '每月1-5日提醒上传材料（F-65）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_sys_cron` t WHERE t.name = 'compliance_material_reminder');

INSERT INTO `xy_sys_cron` (`group_id`, `title`, `name`, `params`, `pattern`, `policy`, `count`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
SELECT 1, '月度对账单生成', 'compliance_statement_generate', '', '0 0 2 16-20 * *', 2, 1, 51, '每月16-20日生成上月对账单（F-61）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_sys_cron` t WHERE t.name = 'compliance_statement_generate');

INSERT INTO `xy_sys_cron` (`group_id`, `title`, `name`, `params`, `pattern`, `policy`, `count`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
SELECT 1, '税务申报任务生成', 'compliance_tax_tasks', '', '0 0 1 1 * *', 2, 1, 52, '每月1日生成增值税/附加税及季度企税任务', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM `xy_sys_cron` t WHERE t.name = 'compliance_tax_tasks');
