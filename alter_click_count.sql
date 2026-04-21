-- 新增 click_count 字段，用于记录短链接点击次数
ALTER TABLE `short_url_map`
  ADD COLUMN `click_count` bigint unsigned NOT NULL DEFAULT 0 COMMENT '点击次数' AFTER `risk_reason`;
