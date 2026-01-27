-- Menu Structure Migration Script
-- Migrates existing menu data to align with original Python Flask project
-- Date: 2026-01-27

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE charonoms;

-- Step 1: Backup existing menu data
CREATE TABLE IF NOT EXISTS menu_backup_20260127 AS SELECT * FROM menu;

-- Step 2: Update existing menu routes (kebab-case to snake_case)
UPDATE menu SET route = 'menu_management' WHERE route = 'menu-management';
UPDATE menu SET route = 'activity_template' WHERE route = 'activity-templates';
UPDATE menu SET route = 'activity_management' WHERE route = 'activities';
UPDATE menu SET route = 'contract_management' WHERE route = 'contracts';
UPDATE menu SET route = 'payment_collection' WHERE route = 'payment-collections';
UPDATE menu SET route = 'refund_orders' WHERE route = 'refund-orders';

-- Step 3: Update existing menu Chinese names
UPDATE menu SET name = '学生管理' WHERE route = 'students' AND name = '学生列表';
UPDATE menu SET name = '教练管理' WHERE route = 'coaches' AND name = '教练列表';
UPDATE menu SET name = '订单管理' WHERE route = 'orders' AND name = '订单列表';
UPDATE menu SET name = '活动管理' WHERE route = 'activity_management' AND name = '活动列表';

-- Step 4: Delete old approval management menus
DELETE FROM menu WHERE route IN ('approval-initiated', 'approval-pending', 'approval-completed');

-- Step 5: Restructure level 1 menus (update IDs and sort_order)
-- Note: This may require careful handling of foreign key constraints
-- Temporarily disable foreign key checks
SET FOREIGN_KEY_CHECKS = 0;

-- Update level 1 menu structure
UPDATE menu SET sort_order = 99 WHERE id = 1 AND name = '系统管理';
UPDATE menu SET sort_order = 1 WHERE id = 2 AND name = '学生管理';
UPDATE menu SET sort_order = 2 WHERE id = 3 AND name = '教练管理';
UPDATE menu SET sort_order = 3 WHERE id = 4 AND name = '订单管理';
UPDATE menu SET sort_order = 4 WHERE id = 5 AND name = '商品管理';
UPDATE menu SET sort_order = 5 WHERE id = 6 AND name = '审批管理';
UPDATE menu SET sort_order = 6 WHERE id = 7 AND name = '活动管理';
UPDATE menu SET sort_order = 8 WHERE id = 8 AND name = '财务管理';

-- Step 6: Insert new financial management menus (if not exist)
INSERT INTO menu (id, name, route, parent_id, sort_order, status)
SELECT 26, '退款子订单', 'refund_childorders', 8, 3, 0
WHERE NOT EXISTS (SELECT 1 FROM menu WHERE id = 26);

INSERT INTO menu (id, name, route, parent_id, sort_order, status)
SELECT 27, '退款管理', 'refund_management', 8, 4, 0
WHERE NOT EXISTS (SELECT 1 FROM menu WHERE id = 27);

INSERT INTO menu (id, name, route, parent_id, sort_order, status)
SELECT 28, '退款支付详情', 'refund_payment_detail', 8, 5, 0
WHERE NOT EXISTS (SELECT 1 FROM menu WHERE id = 28);

INSERT INTO menu (id, name, route, parent_id, sort_order, status)
SELECT 39, '分账明细', 'separate_account', 8, 6, 0
WHERE NOT EXISTS (SELECT 1 FROM menu WHERE id = 39);

-- Step 7: Insert new approval management menus (if not exist)
INSERT INTO menu (id, name, route, parent_id, sort_order, status)
SELECT 30, '审批流类型', 'approval_flow_type', 6, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM menu WHERE id = 30);

INSERT INTO menu (id, name, route, parent_id, sort_order, status)
SELECT 31, '审批流模板', 'approval_flow_template', 6, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM menu WHERE id = 31);

INSERT INTO menu (id, name, route, parent_id, sort_order, status)
SELECT 32, '审批流管理', 'approval_flow_management', 6, 3, 0
WHERE NOT EXISTS (SELECT 1 FROM menu WHERE id = 32);

-- Step 8: Update level 2 menu IDs and parent relationships
-- System Management
UPDATE menu SET id = 1, name = '用户管理' WHERE route = 'accounts' AND parent_id = 1;
UPDATE menu SET id = 4 WHERE route = 'roles' AND parent_id = 1;
UPDATE menu SET id = 13 WHERE route = 'menu_management' AND parent_id = 1;
UPDATE menu SET id = 14 WHERE route = 'permissions' AND parent_id = 1;

-- Student Management
UPDATE menu SET id = 2, parent_id = 2 WHERE route = 'students';

-- Coach Management
UPDATE menu SET id = 6, parent_id = 3 WHERE route = 'coaches';

-- Order Management
UPDATE menu SET id = 3, parent_id = 4 WHERE route = 'orders';
UPDATE menu SET id = 7, name = '学员订单', route = 'student_orders', parent_id = 4, sort_order = 2 WHERE id = 7 OR route = 'childorders';
UPDATE menu SET id = 29, parent_id = 4, sort_order = 3 WHERE route = 'refund_orders';

-- Goods Management
UPDATE menu SET id = 15, name = '商品类别', route = 'goods_category', parent_id = 5, sort_order = 1 WHERE route = 'brands' OR id = 15;
UPDATE menu SET id = 16, name = '商品管理', route = 'goods_management', parent_id = 5, sort_order = 2 WHERE route = 'goods' OR id = 16;
UPDATE menu SET id = 17, name = '套餐管理', route = 'package_management', parent_id = 5, sort_order = 3 WHERE route = 'classifies' OR id = 17;
UPDATE menu SET id = 18, name = '课程管理', route = 'course_management', parent_id = 5, sort_order = 4 WHERE route = 'attributes' OR id = 18;
UPDATE menu SET id = 19, name = '班级管理', route = 'class_management', parent_id = 5, sort_order = 5 WHERE id = 19;

-- Activity Management
UPDATE menu SET id = 21, parent_id = 7, sort_order = 1 WHERE route = 'activity_template';
UPDATE menu SET id = 22, parent_id = 7, sort_order = 2 WHERE route = 'activity_management';
UPDATE menu SET id = 23, name = '优惠券管理', route = 'coupon_management', parent_id = 7, sort_order = 3 WHERE id = 23;

-- Finance Management
UPDATE menu SET id = 24, parent_id = 8, sort_order = 1 WHERE route = 'contract_management';
UPDATE menu SET id = 25, parent_id = 8, sort_order = 2 WHERE route = 'payment_collection';

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;

-- Step 9: Verification queries (optional, for manual checking)
-- SELECT 'Total menus:', COUNT(*) FROM menu;
-- SELECT 'Level 1 menus:', COUNT(*) FROM menu WHERE parent_id IS NULL;
-- SELECT 'Level 2 menus:', COUNT(*) FROM menu WHERE parent_id IS NOT NULL;
-- SELECT * FROM menu ORDER BY parent_id, sort_order;
