-- Correct test data initialization script matching actual schema
USE charonoms;

-- Insert roles (using 'comment' field)
INSERT INTO `role` (`id`, `name`, `comment`, `is_super_admin`, `status`) VALUES
(1, '超级管理员', '拥有系统所有权限', 1, 0),
(2, '普通管理员', '普通管理员权限', 0, 0),
(3, '操作员', '基础操作权限', 0, 0)
ON DUPLICATE KEY UPDATE name=VALUES(name), comment=VALUES(comment), is_super_admin=VALUES(is_super_admin), status=VALUES(status);

-- Update existing admin user's role
UPDATE `useraccount` SET `role_id` = 1 WHERE `username` = 'admin';

-- Insert additional test users
INSERT INTO `useraccount` (`id`, `username`, `password`, `name`, `phone`, `role_id`, `status`) VALUES
(2, 'manager', 'password', '普通管理员', '13800000002', 2, 0),
(3, 'operator', 'password', '操作员', '13800000003', 3, 0)
ON DUPLICATE KEY UPDATE password=VALUES(password), role_id=VALUES(role_id);

-- Insert menus (using 'sort_order' field, no 'level' field)
INSERT INTO `menu` (`id`, `name`, `route`, `parent_id`, `sort_order`, `status`) VALUES
-- Level 1 menus
(1, '系统管理', '', 0, 1, 0),
(2, '学生管理', '', 0, 2, 0),
(3, '教练管理', '', 0, 3, 0),
(4, '商品管理', '', 0, 4, 0),
(5, '订单管理', '', 0, 5, 0),
(6, '活动管理', '', 0, 6, 0),
(7, '财务管理', '', 0, 7, 0),
(8, '审批管理', '', 0, 8, 0),

-- Level 2 menus - System Management
(11, '账号管理', 'accounts', 1, 1, 0),
(12, '角色管理', 'roles', 1, 2, 0),
(13, '权限管理', 'permissions', 1, 3, 0),
(14, '菜单管理', 'menu-management', 1, 4, 0),

-- Level 2 menus - Student Management
(21, '学生列表', 'students', 2, 1, 0),

-- Level 2 menus - Coach Management
(31, '教练列表', 'coaches', 3, 1, 0),

-- Level 2 menus - Product Management
(41, '品牌管理', 'brands', 4, 1, 0),
(42, '分类管理', 'classifies', 4, 2, 0),
(43, '属性管理', 'attributes', 4, 3, 0),
(44, '商品列表', 'goods', 4, 4, 0),

-- Level 2 menus - Order Management
(51, '订单列表', 'orders', 5, 1, 0),
(52, '子订单', 'childorders', 5, 2, 0),

-- Level 2 menus - Activity Management
(61, '活动模板', 'activity-templates', 6, 1, 0),
(62, '活动列表', 'activities', 6, 2, 0),

-- Level 2 menus - Finance Management
(71, '合同管理', 'contracts', 7, 1, 0),
(72, '收款管理', 'payment-collections', 7, 2, 0),
(73, '退款管理', 'refund-orders', 7, 3, 0),

-- Level 2 menus - Approval Management
(81, '我发起的', 'approval-initiated', 8, 1, 0),
(82, '待我审批', 'approval-pending', 8, 2, 0),
(83, '已审批', 'approval-completed', 8, 3, 0)
ON DUPLICATE KEY UPDATE name=VALUES(name), route=VALUES(route), parent_id=VALUES(parent_id), sort_order=VALUES(sort_order), status=VALUES(status);

-- Insert permissions
INSERT INTO `permissions` (`id`, `name`, `action_id`, `menu_id`, `status`) VALUES
-- System Management Permissions
(1, '查看账号', 'view_account', 11, 0),
(2, '新增账号', 'add_account', 11, 0),
(3, '编辑账号', 'edit_account', 11, 0),
(4, '删除账号', 'delete_account', 11, 0),

(5, '查看角色', 'view_role', 12, 0),
(6, '新增角色', 'add_role', 12, 0),
(7, '编辑角色', 'edit_role', 12, 0),
(8, '删除角色', 'delete_role', 12, 0),

(9, '查看权限', 'view_permission', 13, 0),
(10, '编辑权限', 'edit_permission', 13, 0),

(11, '查看菜单', 'view_menu', 14, 0),
(12, '编辑菜单', 'edit_menu', 14, 0),

-- Student Management Permissions
(21, '查看学生', 'view_student', 21, 0),
(22, '新增学生', 'add_student', 21, 0),
(23, '编辑学生', 'edit_student', 21, 0),
(24, '删除学生', 'delete_student', 21, 0),
(25, '启用学生', 'enable_student', 21, 0),
(26, '禁用学生', 'disable_student', 21, 0),

-- Coach Management Permissions
(31, '查看教练', 'view_coach', 31, 0),
(32, '新增教练', 'add_coach', 31, 0),
(33, '编辑教练', 'edit_coach', 31, 0),
(34, '删除教练', 'delete_coach', 31, 0),
(35, '启用教练', 'enable_coach', 31, 0),
(36, '禁用教练', 'disable_coach', 31, 0),

-- Product Management Permissions
(41, '查看品牌', 'view_brand', 41, 0),
(42, '新增品牌', 'add_brand', 41, 0),
(43, '编辑品牌', 'edit_brand', 41, 0),
(44, '删除品牌', 'delete_brand', 41, 0),

(51, '查看分类', 'view_classify', 42, 0),
(52, '新增分类', 'add_classify', 42, 0),
(53, '编辑分类', 'edit_classify', 42, 0),
(54, '删除分类', 'delete_classify', 42, 0),

(61, '查看属性', 'view_attribute', 43, 0),
(62, '新增属性', 'add_attribute', 43, 0),
(63, '编辑属性', 'edit_attribute', 43, 0),
(64, '删除属性', 'delete_attribute', 43, 0),

(71, '查看商品', 'view_goods', 44, 0),
(72, '新增商品', 'add_goods', 44, 0),
(73, '编辑商品', 'edit_goods', 44, 0),
(74, '删除商品', 'delete_goods', 44, 0),

-- Order Management Permissions
(81, '查看订单', 'view_order', 51, 0),
(82, '新增订单', 'add_order', 51, 0),
(83, '编辑订单', 'edit_order', 51, 0),
(84, '删除订单', 'delete_order', 51, 0),

(91, '查看子订单', 'view_childorder', 52, 0),
(92, '新增子订单', 'add_childorder', 52, 0),
(93, '编辑子订单', 'edit_childorder', 52, 0),
(94, '删除子订单', 'delete_childorder', 52, 0),

-- Activity Management Permissions
(101, '查看活动模板', 'view_activity_template', 61, 0),
(102, '新增活动模板', 'add_activity_template', 61, 0),
(103, '编辑活动模板', 'edit_activity_template', 61, 0),
(104, '删除活动模板', 'delete_activity_template', 61, 0),

(111, '查看活动', 'view_activity', 62, 0),
(112, '新增活动', 'add_activity', 62, 0),
(113, '编辑活动', 'edit_activity', 62, 0),
(114, '删除活动', 'delete_activity', 62, 0),

-- Finance Management Permissions
(121, '查看合同', 'view_contract', 71, 0),
(122, '新增合同', 'add_contract', 71, 0),
(123, '编辑合同', 'edit_contract', 71, 0),
(124, '删除合同', 'delete_contract', 71, 0),

(131, '查看收款', 'view_payment', 72, 0),
(132, '新增收款', 'add_payment', 72, 0),
(133, '编辑收款', 'edit_payment', 72, 0),
(134, '删除收款', 'delete_payment', 72, 0),

(141, '查看退款', 'view_refund', 73, 0),
(142, '新增退款', 'add_refund', 73, 0),
(143, '编辑退款', 'edit_refund', 73, 0),
(144, '删除退款', 'delete_refund', 73, 0),

-- Approval Management Permissions
(151, '查看我发起的', 'view_my_approval', 81, 0),
(152, '查看待审批', 'view_pending_approval', 82, 0),
(153, '审批', 'approve', 82, 0),
(154, '查看已审批', 'view_approved', 83, 0)
ON DUPLICATE KEY UPDATE name=VALUES(name), action_id=VALUES(action_id), menu_id=VALUES(menu_id), status=VALUES(status);

-- Assign all permissions to super admin role (role_id=1)
INSERT INTO `role_permissions` (`role_id`, `permissions_id`)
SELECT 1, id FROM `permissions`
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- Assign most permissions to regular admin (role_id=2) - exclude delete operations
INSERT INTO `role_permissions` (`role_id`, `permissions_id`)
SELECT 2, id FROM `permissions` WHERE action_id NOT LIKE '%delete%'
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- Assign view-only permissions to operator (role_id=3)
INSERT INTO `role_permissions` (`role_id`, `permissions_id`)
SELECT 3, id FROM `permissions` WHERE action_id LIKE 'view%'
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

SELECT '✓ Data initialization completed successfully!' AS Result;
SELECT CONCAT('✓ Inserted ', COUNT(*), ' roles') AS Result FROM role;
SELECT CONCAT('✓ Inserted ', COUNT(*), ' users') AS Result FROM useraccount;
SELECT CONCAT('✓ Inserted ', COUNT(*), ' menus') AS Result FROM menu;
SELECT CONCAT('✓ Inserted ', COUNT(*), ' permissions') AS Result FROM permissions;
SELECT CONCAT('✓ Assigned ', COUNT(*), ' role-permission mappings') AS Result FROM role_permissions;
