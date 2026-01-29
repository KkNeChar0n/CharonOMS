-- 活动模板表
CREATE TABLE IF NOT EXISTS `activity_template` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `name` varchar(100) NOT NULL COMMENT '模板名称',
  `type` tinyint NOT NULL COMMENT '模板类型：1-满减 2-满折 3-满赠',
  `select_type` tinyint NOT NULL COMMENT '选择类型：1-按分类 2-按商品',
  `status` tinyint DEFAULT '1' COMMENT '状态：0-启用 1-禁用',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动模板表';

-- 活动模板与分类关联表
CREATE TABLE IF NOT EXISTS `activity_template_classify` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `template_id` int NOT NULL COMMENT '模板ID',
  `classify_id` int NOT NULL COMMENT '分类ID',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_classify_id` (`classify_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动模板分类关联表';

-- 活动模板与商品关联表
CREATE TABLE IF NOT EXISTS `activity_template_goods` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `template_id` int NOT NULL COMMENT '模板ID',
  `goods_id` int NOT NULL COMMENT '商品ID',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_goods_id` (`goods_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动模板商品关联表';

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
