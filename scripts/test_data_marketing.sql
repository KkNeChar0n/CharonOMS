-- 营销管理模块测试数据
-- 日期: 2026-01-30

SET NAMES utf8mb4;
USE charonoms;

-- 清理现有测试数据
DELETE FROM orders_activity;
DELETE FROM activity_detail;
DELETE FROM activity;
DELETE FROM activity_template_goods;
DELETE FROM activity_template;

-- 插入活动模板测试数据
INSERT INTO activity_template (id, name, type, select_type, status, create_time, update_time) VALUES
(1, '春季满减促销', 1, 1, 0, NOW(), NOW()),
(2, '夏季满折促销', 2, 2, 0, NOW(), NOW()),
(3, '秋季满赠促销', 3, 1, 1, NOW(), NOW());

-- 插入活动模板商品关联（使用实际存在的分类和商品数据）
-- 模板1 按分类选择（使用实际存在的分类ID）
INSERT INTO activity_template_goods (template_id, classify_id) VALUES
(1, 15),
(1, 16);

-- 模板2 按商品选择（使用实际存在的商品ID）
INSERT INTO activity_template_goods (template_id, goods_id) VALUES
(2, 14),
(2, 15),
(2, 16);

-- 插入活动测试数据
INSERT INTO activity (id, name, template_id, start_time, end_time, status, create_time) VALUES
(1, '春季大促销', 1, '2026-03-01 00:00:00', '2026-03-31 23:59:59', 0, NOW()),
(2, '夏季满折活动', 2, '2026-06-01 00:00:00', '2026-06-30 23:59:59', 0, NOW()),
(3, '测试已禁用活动', 1, '2026-01-01 00:00:00', '2026-01-31 23:59:59', 1, NOW());

-- 插入活动详情（满折规则）
-- 注意：满折类型的discount_value使用百分比形式
-- 95表示9.5折（顾客付95%），90表示9折（顾客付90%），85表示8.5折（顾客付85%）
INSERT INTO activity_detail (activity_id, threshold_amount, discount_value) VALUES
(2, 100.00, 95.00),
(2, 200.00, 90.00),
(2, 500.00, 85.00);

SELECT '营销管理模块测试数据插入成功' AS status;
