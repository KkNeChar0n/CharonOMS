-- 恢复 charonoms 数据库缺失的表
-- 日期: 2026-01-30
-- 说明: 根据 zhixinstudentsaas 数据库创建 charonoms 中缺失的表

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE charonoms;

/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;

-- 1. 创建 activity_template 表（营销活动模板）
DROP TABLE IF EXISTS `activity_template`;
CREATE TABLE `activity_template` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `name` varchar(100) NOT NULL COMMENT '模板名称',
  `type` int NOT NULL COMMENT '活动类型：1=满减 2=满折 3=满赠',
  `select_type` int NOT NULL COMMENT '选择方式：1=按分类 2=按商品',
  `status` int DEFAULT '0' COMMENT '状态：0=启用 1=禁用',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动模板表';

-- 2. 创建 activity 表（营销活动）
DROP TABLE IF EXISTS `activity`;
CREATE TABLE `activity` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `name` varchar(100) NOT NULL COMMENT '活动名称',
  `template_id` int NOT NULL COMMENT '关联的活动模板ID',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `status` int DEFAULT '0' COMMENT '状态：0=启用 1=禁用',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `template_id` (`template_id`),
  KEY `idx_time_status` (`start_time`,`end_time`,`status`),
  CONSTRAINT `activity_ibfk_1` FOREIGN KEY (`template_id`) REFERENCES `activity_template` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='营销活动表';

-- 3. 创建 activity_detail 表（活动详情）
DROP TABLE IF EXISTS `activity_detail`;
CREATE TABLE `activity_detail` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `activity_id` int NOT NULL COMMENT '活动ID',
  `threshold_amount` decimal(10,2) NOT NULL COMMENT '门槛金额（满XX元）',
  `discount_value` decimal(10,2) NOT NULL COMMENT '优惠值（减/折/赠的具体值）',
  PRIMARY KEY (`id`),
  KEY `activity_id` (`activity_id`),
  CONSTRAINT `activity_detail_ibfk_1` FOREIGN KEY (`activity_id`) REFERENCES `activity` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动详情表（满减满折规则）';

-- 4. 创建 activity_template_goods 表（活动模板商品关联）
DROP TABLE IF EXISTS `activity_template_goods`;
CREATE TABLE `activity_template_goods` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `template_id` int NOT NULL COMMENT '模板ID',
  `goods_id` int DEFAULT NULL COMMENT '商品ID（select_type=2时使用）',
  `classify_id` int DEFAULT NULL COMMENT '分类ID（select_type=1时使用）',
  PRIMARY KEY (`id`),
  KEY `template_id` (`template_id`),
  KEY `goods_id` (`goods_id`),
  KEY `classify_id` (`classify_id`),
  CONSTRAINT `activity_template_goods_ibfk_1` FOREIGN KEY (`template_id`) REFERENCES `activity_template` (`id`) ON DELETE CASCADE,
  CONSTRAINT `activity_template_goods_ibfk_2` FOREIGN KEY (`goods_id`) REFERENCES `goods` (`id`) ON DELETE CASCADE,
  CONSTRAINT `activity_template_goods_ibfk_3` FOREIGN KEY (`classify_id`) REFERENCES `classify` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动模板商品关联表';

-- 5. 创建 orders 表（订单）
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `student_id` int NOT NULL COMMENT '学生ID',
  `expected_payment_time` datetime DEFAULT NULL COMMENT '预计付款时间',
  `amount_receivable` decimal(10,2) NOT NULL COMMENT '应收金额',
  `amount_received` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '实收金额',
  `discount_amount` decimal(10,2) DEFAULT '0.00' COMMENT '优惠金额',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `status` tinyint NOT NULL DEFAULT '10' COMMENT '订单状态: 10=草稿, 20=审核中, 30=已通过, 40=已驳回, 99=已作废',
  PRIMARY KEY (`id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_status` (`status`),
  KEY `idx_create_time` (`create_time`),
  CONSTRAINT `fk_order_student` FOREIGN KEY (`student_id`) REFERENCES `student` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 6. 创建 childorders 表（子订单）
DROP TABLE IF EXISTS `childorders`;
CREATE TABLE `childorders` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '子订单ID',
  `parentsid` int NOT NULL COMMENT '父订单ID',
  `goodsid` int NOT NULL COMMENT '商品ID',
  `amount_receivable` decimal(10,2) NOT NULL COMMENT '应收金额（商品标价）',
  `amount_received` decimal(10,2) NOT NULL COMMENT '实收金额（标准售价）',
  `discount_amount` decimal(10,2) DEFAULT '0.00' COMMENT '优惠金额',
  `status` int DEFAULT '10' COMMENT '状态：10=草稿 20=审核中 30=已通过 40=已驳回 99=已作废',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `parentsid` (`parentsid`),
  KEY `goodsid` (`goodsid`),
  CONSTRAINT `childorders_ibfk_1` FOREIGN KEY (`parentsid`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `childorders_ibfk_2` FOREIGN KEY (`goodsid`) REFERENCES `goods` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='子产品订单表';

-- 7. 创建 orders_activity 表（订单活动关联）
DROP TABLE IF EXISTS `orders_activity`;
CREATE TABLE `orders_activity` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `orders_id` int NOT NULL COMMENT '订单ID',
  `activity_id` int NOT NULL COMMENT '活动ID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orders_activity` (`orders_id`,`activity_id`),
  KEY `idx_orders_id` (`orders_id`),
  KEY `idx_activity_id` (`activity_id`),
  CONSTRAINT `orders_activity_ibfk_1` FOREIGN KEY (`orders_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `orders_activity_ibfk_2` FOREIGN KEY (`activity_id`) REFERENCES `activity` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单活动关联表';

/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;

-- 验证表创建
SELECT 'CharonOMS 缺失表已创建完成！' AS status;
SELECT '以下是新创建的表：' AS info;
SHOW TABLES LIKE 'activity%';
SHOW TABLES LIKE 'orders%';
SHOW TABLES LIKE 'childorders';
