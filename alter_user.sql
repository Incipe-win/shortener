-- 新增注册用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `username` VARCHAR(64) NOT NULL UNIQUE COMMENT '用户名',
    `password_hash` VARCHAR(128) NOT NULL COMMENT '密码哈希(bcrypt)',
    `is_del` tinyint UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除：0正常1删除',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_username` (`username`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = '注册用户表';

-- 为 short_url_map 表新增 user_id 字段
ALTER TABLE `short_url_map`
    ADD COLUMN `user_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '创建者用户ID(未注册为NULL)' AFTER `create_by`;

ALTER TABLE `short_url_map`
    ADD INDEX `idx_user_id` (`user_id`);
