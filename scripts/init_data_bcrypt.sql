-- Initialize data for CharonOMS (with bcrypt encrypted passwords)
-- This script should be run after the database schema has been created

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Disable foreign key checks temporarily
SET FOREIGN_KEY_CHECKS = 0;

-- Clear existing menu and permission data (keep users and roles)
DELETE FROM `role_permissions`;
DELETE FROM `permissions`;
DELETE FROM `menu`;

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;

-- Insert menus (aligned with original Python Flask project zhixinstudentsaas database)
INSERT INTO `menu` (`id`, `name`, `route`, `parent_id`, `sort_order`, `status`) VALUES
-- Level 1 menus (parent_id = NULL, using 100+ IDs to avoid conflicts)
(101, '学生管理', '', NULL, 1, 0),
(102, '教练管理', '', NULL, 2, 0),
(103, '订单管理', '', NULL, 3, 0),
(104, '商品管理', '', NULL, 4, 0),
(105, '审批流管理', '', NULL, 5, 0),
(106, '营销管理', '', NULL, 6, 0),
(107, '合同管理', '', NULL, 7, 0),
(108, '财务管理', '', NULL, 8, 0),
(109, '系统设置', '', NULL, 99, 0),

-- Level 2 menus - Student Management (parent_id=101)
(5, '学生管理', 'students', 101, 1, 0),

-- Level 2 menus - Coach Management (parent_id=102)
(6, '教练管理', 'coaches', 102, 1, 0),

-- Level 2 menus - Order Management (parent_id=103)
(7, '订单管理', 'orders', 103, 1, 0),
(20, '子订单管理', 'childorders', 103, 2, 0),
(29, '退款订单', 'refund_orders', 103, 3, 0),
(37, '子退费订单', 'refund_childorders', 103, 4, 0),

-- Level 2 menus - Goods Management (parent_id=104)
(16, '品牌管理', 'brands', 104, 1, 0),
(17, '属性管理', 'attributes', 104, 2, 0),
(18, '类型管理', 'classifies', 104, 3, 0),
(19, '商品管理', 'goods', 104, 4, 0),

-- Level 2 menus - Approval Flow Management (parent_id=105)
(31, '审批流类型', 'approval_flow_type', 105, 1, 0),
(32, '审批流模板', 'approval_flow_template', 105, 2, 0),
(34, '审批流管理', 'approval_flow_management', 105, 3, 0),

-- Level 2 menus - Marketing Management (parent_id=106)
(22, '活动模板', 'activity_template', 106, 1, 0),
(23, '活动管理', 'activity_management', 106, 2, 0),

-- Level 2 menus - Contract Management (parent_id=107)
(25, '合同管理', 'contract_management', 107, 1, 0),

-- Level 2 menus - Finance Management (parent_id=108)
(27, '收款管理', 'payment_collection', 108, 1, 0),
(28, '分账明细', 'separate_account', 108, 2, 0),
(38, '退费管理', 'refund_management', 108, 4, 0),
(39, '退费明细', 'refund_payment_detail', 108, 5, 0),

-- Level 2 menus - System Settings (parent_id=109)
(8, '账号管理', 'accounts', 109, 1, 0),
(33, '权限管理', 'permissions', 109, 10, 0),
(40, '角色管理', 'roles', 109, 102, 0),
(41, '菜单管理', 'menu_management', 109, 103, 0)
ON DUPLICATE KEY UPDATE name=VALUES(name), route=VALUES(route), parent_id=VALUES(parent_id), sort_order=VALUES(sort_order), status=VALUES(status);

-- Insert permissions linked to menus
INSERT INTO `permissions` (`id`, `name`, `action_id`, `menu_id`, `status`) VALUES
-- System Settings Permissions
(1, '查看账号', 'view_account', 8, 0),
(2, '新增账号', 'add_account', 8, 0),
(3, '编辑账号', 'edit_account', 8, 0),
(4, '删除账号', 'delete_account', 8, 0),

(5, '查看角色', 'view_role', 40, 0),
(6, '新增角色', 'add_role', 40, 0),
(7, '编辑角色', 'edit_role', 40, 0),
(8, '删除角色', 'delete_role', 40, 0),

(9, '查看权限', 'view_permission', 33, 0),
(10, '编辑权限', 'edit_permission', 33, 0),

(11, '查看菜单', 'view_menu', 41, 0),
(12, '编辑菜单', 'edit_menu', 41, 0),

-- Student Management Permissions
(21, '查看学生', 'view_student', 5, 0),
(22, '新增学生', 'add_student', 5, 0),
(23, '编辑学生', 'edit_student', 5, 0),
(24, '删除学生', 'delete_student', 5, 0),

-- Coach Management Permissions
(31, '查看教练', 'view_coach', 6, 0),
(32, '新增教练', 'add_coach', 6, 0),
(33, '编辑教练', 'edit_coach', 6, 0),
(34, '删除教练', 'delete_coach', 6, 0),

-- Order Management Permissions
(41, '查看订单', 'view_order', 7, 0),
(42, '新增订单', 'add_order', 7, 0),
(43, '编辑订单', 'edit_order', 7, 0),
(44, '删除订单', 'delete_order', 7, 0),

-- Goods Management Permissions
(51, '查看商品', 'view_goods', 19, 0),
(52, '新增商品', 'add_goods', 19, 0),
(53, '编辑商品', 'edit_goods', 19, 0),
(54, '删除商品', 'delete_goods', 19, 0),

(61, '查看品牌', 'view_brand', 16, 0),
(62, '新增品牌', 'add_brand', 16, 0),

(71, '查看属性', 'view_attribute', 17, 0),
(72, '新增属性', 'add_attribute', 17, 0),

(81, '查看类型', 'view_classify', 18, 0),
(82, '新增类型', 'add_classify', 18, 0),

-- Marketing Management Permissions
(91, '查看活动', 'view_activity', 23, 0),
(92, '新增活动', 'add_activity', 23, 0),
(93, '编辑活动', 'edit_activity', 23, 0),
(94, '删除活动', 'delete_activity', 23, 0),

-- Contract Management Permissions
(101, '查看合同', 'view_contract', 25, 0),
(102, '新增合同', 'add_contract', 25, 0),
(103, '编辑合同', 'edit_contract', 25, 0),

-- Finance Management Permissions
(111, '查看收款', 'view_payment', 27, 0),
(112, '新增收款', 'add_payment', 27, 0),

(121, '查看退费', 'view_refund', 38, 0),
(122, '处理退费', 'process_refund', 38, 0),

-- Approval Flow Management Permissions
(131, '查看审批', 'view_approval', 34, 0),
(132, '处理审批', 'process_approval', 34, 0)
ON DUPLICATE KEY UPDATE name=VALUES(name), action_id=VALUES(action_id), menu_id=VALUES(menu_id), status=VALUES(status);

-- Assign all permissions to super admin role
INSERT INTO `role_permissions` (`role_id`, `permissions_id`)
SELECT 1, id FROM `permissions` WHERE status = 0
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);
