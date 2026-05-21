-- +goose Up
-- 配置表
CREATE TABLE IF NOT EXISTS `configs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `group` VARCHAR(100) NOT NULL COMMENT '配置分组',
    `key` VARCHAR(100) NOT NULL COMMENT '配置键',
    `value` TEXT COMMENT '配置值',
    `type` VARCHAR(20) NOT NULL DEFAULT 'string' COMMENT '配置类型：string, number, boolean, json',
    `label` VARCHAR(200) COMMENT '配置标签/名称',
    `help_text` TEXT COMMENT '帮助文本',
    `sort_order` INT DEFAULT 0 COMMENT '排序',
    `is_system` BOOLEAN DEFAULT FALSE COMMENT '是否系统配置（不可删除）',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_group_key` (`group`, `key`),
    INDEX `idx_group` (`group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- +goose Down
DROP TABLE IF EXISTS `configs`;
