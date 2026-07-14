# 数据库迁移工具

## 使用方法

```bash
cd server

# 交互式菜单
gf run tools.go

# 直接执行
gf run tools.go --args "migrate up"
gf run tools.go --args "migrate status"
gf run tools.go --args "migrate history"
```

## 迁移文件命名规则

文件放在 `cmd_tools/migrate/` 目录下，命名格式：

```
{版本号}_{描述}.sql           -- 通用 SQL（两种数据库都能执行）
{版本号}_{描述}.mysql.sql     -- 仅 MySQL
{版本号}_{描述}.pgsql.sql     -- 仅 PostgreSQL
```

工具根据 `config.yaml` 中的数据库类型自动选择：优先匹配带数据库后缀的文件，找不到则用通用 `.sql`。

版本号使用语义化版本（如 `1.2.6`），工具按版本号排序依次执行。

`1.5.6_compliance_assistant_core` 只新增经营合规助手表和索引。生产回滚时关闭
`complianceAssistant` 相关功能开关并保留新增数据，不执行 `DROP`、`TRUNCATE` 或批量删除。

---

## 各类数据库变更的 SQL 写法

### 1. 新建表

**PostgreSQL:**
```sql
CREATE TABLE IF NOT EXISTS xy_example (
    id bigserial PRIMARY KEY,
    name varchar(64) NOT NULL DEFAULT '',
    status smallint NOT NULL DEFAULT 1,
    created_at bigint
);
COMMENT ON TABLE xy_example IS '示例表';
COMMENT ON COLUMN xy_example.name IS '名称';
```

**MySQL:**
```sql
CREATE TABLE IF NOT EXISTS `xy_example` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态',
    `created_at` bigint UNSIGNED DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='示例表';
```

### 2. 添加字段

**PostgreSQL:**
```sql
ALTER TABLE xy_example ADD COLUMN IF NOT EXISTS remark varchar(255) NOT NULL DEFAULT '';
COMMENT ON COLUMN xy_example.remark IS '备注';
```

**MySQL（不支持 IF NOT EXISTS，需存储过程）:**
```sql
DELIMITER $$
DROP PROCEDURE IF EXISTS _add_col$$
CREATE PROCEDURE _add_col(IN p_tbl VARCHAR(64), IN p_col VARCHAR(64), IN p_def TEXT)
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = DATABASE() AND table_name = p_tbl AND column_name = p_col
    ) THEN
        SET @sql = CONCAT('ALTER TABLE `', p_tbl, '` ADD COLUMN `', p_col, '` ', p_def);
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END$$
DELIMITER ;

CALL _add_col('xy_example', 'remark', "varchar(255) NOT NULL DEFAULT '' COMMENT '备注'");

DROP PROCEDURE IF EXISTS _add_col;
```

### 3. 修改字段类型/默认值

**PostgreSQL:**
```sql
ALTER TABLE xy_example ALTER COLUMN name TYPE varchar(128);
ALTER TABLE xy_example ALTER COLUMN status SET DEFAULT 0;
```

**MySQL:**
```sql
ALTER TABLE `xy_example` MODIFY COLUMN `name` varchar(128) NOT NULL DEFAULT '' COMMENT '名称';
ALTER TABLE `xy_example` ALTER COLUMN `status` SET DEFAULT 0;
```

### 4. 添加索引

**PostgreSQL:**
```sql
CREATE INDEX IF NOT EXISTS idx_example_name ON xy_example USING btree (name);
CREATE UNIQUE INDEX IF NOT EXISTS uk_example_code ON xy_example USING btree (code);
```

**MySQL:**
```sql
-- MySQL 无 IF NOT EXISTS，用存储过程
DELIMITER $$
DROP PROCEDURE IF EXISTS _add_idx$$
CREATE PROCEDURE _add_idx(IN p_tbl VARCHAR(64), IN p_idx VARCHAR(64), IN p_def TEXT)
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.statistics
        WHERE table_schema = DATABASE() AND table_name = p_tbl AND index_name = p_idx
    ) THEN
        SET @sql = CONCAT('ALTER TABLE `', p_tbl, '` ADD INDEX `', p_idx, '` ', p_def);
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END$$
DELIMITER ;

CALL _add_idx('xy_example', 'idx_example_name', '(`name`)');

DROP PROCEDURE IF EXISTS _add_idx;
```

### 5. 删除字段

**PostgreSQL:**
```sql
ALTER TABLE xy_example DROP COLUMN IF EXISTS old_field;
```

**MySQL:**
```sql
DELIMITER $$
DROP PROCEDURE IF EXISTS _drop_col$$
CREATE PROCEDURE _drop_col(IN p_tbl VARCHAR(64), IN p_col VARCHAR(64))
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = DATABASE() AND table_name = p_tbl AND column_name = p_col
    ) THEN
        SET @sql = CONCAT('ALTER TABLE `', p_tbl, '` DROP COLUMN `', p_col, '`');
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END$$
DELIMITER ;

CALL _drop_col('xy_example', 'old_field');

DROP PROCEDURE IF EXISTS _drop_col;
```

### 6. 删除表

```sql
-- 两种数据库通用
DROP TABLE IF EXISTS xy_example;
```

### 7. 插入种子数据（幂等）

**PostgreSQL:**
```sql
INSERT INTO xy_example (id, name, status) VALUES (1, '默认', 1)
ON CONFLICT (id) DO NOTHING;
```

**MySQL:**
```sql
INSERT IGNORE INTO `xy_example` (`id`, `name`, `status`) VALUES (1, '默认', 1);
```

### 8. 插入菜单记录（幂等）

**PostgreSQL:**
```sql
INSERT INTO xy_admin_menu (parent_id, type, title, name, path, component, icon, sort, status, create_time, update_time)
SELECT 0, 2, '新功能', 'NewFeature', 'new-feature', '/new-feature/index', 'ri:star-line', 50, 1,
       EXTRACT(EPOCH FROM NOW())::int, EXTRACT(EPOCH FROM NOW())::int
WHERE NOT EXISTS (SELECT 1 FROM xy_admin_menu WHERE name = 'NewFeature');
```

**MySQL:**
```sql
INSERT INTO `xy_admin_menu` (`parent_id`, `type`, `title`, `name`, `path`, `component`, `icon`, `sort`, `status`, `create_time`, `update_time`)
SELECT 0, 2, '新功能', 'NewFeature', 'new-feature', '/new-feature/index', 'ri:star-line', 50, 1,
       UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `xy_admin_menu` WHERE `name` = 'NewFeature');
```

---

## 注意事项

1. **所有 SQL 必须幂等** — 重复执行不报错、不产生重复数据
2. **每个迁移文件对应一个版本号** — 同版本号只会执行一次
3. **执行失败会停止** — 修复问题后重新运行，已成功的不会重复执行
4. **先在测试环境验证** — 迁移涉及数据结构变更，务必先测试
