-- Smart-Shortener Gateway: DDL 变更脚本
-- 为 short_url_map 表新增 AI 分析和安全检查相关字段

ALTER TABLE `short_url_map`
    ADD COLUMN `ai_summary`   TEXT         DEFAULT NULL COMMENT 'AI生成的页面摘要' AFTER `surl`,
    ADD COLUMN `ai_keywords`  VARCHAR(512) DEFAULT NULL COMMENT 'AI提取的关键词(JSON数组)' AFTER `ai_summary`,
    ADD COLUMN `ai_slug`      VARCHAR(128) DEFAULT NULL COMMENT 'AI生成的语义化短链' AFTER `ai_keywords`,
    ADD COLUMN `risk_level`   VARCHAR(16)  DEFAULT 'pending' COMMENT '安全等级:safe/warning/danger/pending' AFTER `ai_slug`,
    ADD COLUMN `risk_reason`  VARCHAR(512) DEFAULT NULL COMMENT '风险原因' AFTER `risk_level`;

ALTER TABLE `short_url_map`
    ADD INDEX `idx_ai_slug` (`ai_slug`),
    ADD INDEX `idx_risk_level` (`risk_level`);
