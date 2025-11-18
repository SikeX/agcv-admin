-- 数据库迁移SQL - 添加psid字段
-- 执行日期: 2024
-- 说明: 为agvc_nbq_setting表添加psid字段

-- 为agvc_nbq_setting表添加psid字段
ALTER TABLE `agvc_nbq_setting` ADD COLUMN `psid` INT DEFAULT NULL COMMENT '电站编号';

-- 为psid字段创建索引，提高查询性能
CREATE INDEX `idx_psid` ON `agvc_nbq_setting`(`psid`);

-- 注意：agvc_nbq_his不是数据库表，而是用于InfluxDB查询的结构体，不需要ALTER TABLE
-- InfluxDB中的数据通过tags进行过滤，psid已在代码中作为tag添加
