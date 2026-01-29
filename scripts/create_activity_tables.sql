-- 活动表
CREATE TABLE IF NOT EXISTS `activity` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `template_id` int NOT NULL COMMENT '关联的活动模板ID',
  `name` varchar(100) NOT NULL COMMENT '活动名称',
  `start_time` timestamp NOT NULL COMMENT '开始时间',
  `end_time` timestamp NOT NULL COMMENT '结束时间',
  `status` int DEFAULT '1' COMMENT '状态：0-启用 1-禁用',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_status` (`status`),
  KEY `idx_time` (`start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动表';

-- 活动细节表
CREATE TABLE IF NOT EXISTS `activity_detail` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '细节ID',
  `activity_id` int NOT NULL COMMENT '关联的活动ID',
  `goods_id` int NOT NULL COMMENT '商品ID',
  `discount` decimal(10,2) NOT NULL COMMENT '折扣值',
  PRIMARY KEY (`id`),
  KEY `idx_activity_id` (`activity_id`),
  KEY `idx_goods_id` (`goods_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动细节表';
