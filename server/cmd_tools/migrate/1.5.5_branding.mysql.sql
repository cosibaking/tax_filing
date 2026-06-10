-- Migration: 1.5.5
-- Description: 更新站点品牌为「金税管家」并替换 LOGO

UPDATE `xy_sys_config`
SET `value` = '金税管家',
    `update_time` = UNIX_TIMESTAMP()
WHERE `key` = 'site_name';

UPDATE `xy_sys_config`
SET `value` = '/attachment/upload/20260610/jinshui-logo.png',
    `update_time` = UNIX_TIMESTAMP()
WHERE `key` = 'site_logo';
