-- CharonOMS 测试数据初始化脚本
-- 注意：请先创建数据库 charonoms 并选择使用

USE charonoms;

-- ===== 1. 创建角色表 =====
CREATE TABLE IF NOT EXISTS `role` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '角色描述',
  `is_super_admin` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否超级管理员 0-否 1-是',
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态 0-正常 1-禁用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ===== 2. 创建用户账号表 =====
CREATE TABLE IF NOT EXISTS `useraccount` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码',
  `name` VARCHAR(100) DEFAULT NULL COMMENT '姓名',
  `phone` VARCHAR(20) DEFAULT NULL UNIQUE COMMENT '手机号',
  `role_id` INT UNSIGNED NOT NULL COMMENT '角色ID',
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态 0-正常 1-禁用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户账号表';

-- ===== 3. 创建菜单表 =====
CREATE TABLE IF NOT EXISTS `menu` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL COMMENT '菜单名称',
  `route` VARCHAR(100) DEFAULT NULL COMMENT '路由路径',
  `parent_id` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级菜单ID',
  `level` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '层级 0-一级菜单 1-二级菜单',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态 0-启用 1-禁用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';

-- ===== 4. 创建权限表 =====
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
  `action_id` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限标识',
  `menu_id` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属菜单ID',
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态 0-启用 1-禁用',
  INDEX `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- ===== 5. 创建角色权限关联表 =====
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `role_id` INT UNSIGNED NOT NULL COMMENT '角色ID',
  `permissions_id` INT UNSIGNED NOT NULL COMMENT '权限ID',
  UNIQUE KEY `idx_role_permission` (`role_id`, `permissions_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- ===== 插入测试数据 =====

-- 插入角色数据
INSERT INTO `role` (`id`, `name`, `description`, `is_super_admin`, `status`) VALUES
(1, '超级管理员', '拥有系统所有权限', 1, 0),
(2, '普通管理员', '普通管理员权限', 0, 0),
(3, '操作员', '基础操作权限', 0, 0);

-- 插入测试用户（密码都是: password）
INSERT INTO `useraccount` (`id`, `username`, `password`, `name`, `phone`, `role_id`, `status`) VALUES
(1, 'admin', 'password', '超级管理员', '13800000001', 1, 0),
(2, 'manager', 'password', '普通管理员', '13800000002', 2, 0),
(3, 'operator', 'password', '操作员', '13800000003', 3, 0);

-- 插入菜单数据
INSERT INTO `menu` (`id`, `name`, `route`, `parent_id`, `level`, `sort`, `status`) VALUES
-- 一级菜单
(1, '系统管理', '', 0, 0, 1, 0),
(2, '学生管理', '', 0, 0, 2, 0),
(3, '教练管理', '', 0, 0, 3, 0),
(4, '商品管理', '', 0, 0, 4, 0),
(5, '订单管理', '', 0, 0, 5, 0),
(6, '活动管理', '', 0, 0, 6, 0),
(7, '财务管理', '', 0, 0, 7, 0),
(8, '审批管理', '', 0, 0, 8, 0),

-- 二级菜单 - 系统管理
(11, '账号管理', 'accounts', 1, 1, 1, 0),
(12, '角色管理', 'roles', 1, 1, 2, 0),
(13, '权限管理', 'permissions', 1, 1, 3, 0),
(14, '菜单管理', 'menu-management', 1, 1, 4, 0),

-- 二级菜单 - 学生管理
(21, '学生列表', 'students', 2, 1, 1, 0),

-- 二级菜单 - 教练管理
(31, '教练列表', 'coaches', 3, 1, 1, 0),

-- 二级菜单 - 商品管理
(41, '品牌管理', 'brands', 4, 1, 1, 0),
(42, '分类管理', 'classifies', 4, 1, 2, 0),
(43, '属性管理', 'attributes', 4, 1, 3, 0),
(44, '商品列表', 'goods', 4, 1, 4, 0),

-- 二级菜单 - 订单管理
(51, '订单列表', 'orders', 5, 1, 1, 0),
(52, '子订单', 'childorders', 5, 1, 2, 0),

-- 二级菜单 - 活动管理
(61, '活动模板', 'activity-templates', 6, 1, 1, 0),
(62, '活动列表', 'activities', 6, 1, 2, 0),

-- 二级菜单 - 财务管理
(71, '合同管理', 'contracts', 7, 1, 1, 0),
(72, '收款管理', 'payment-collections', 7, 1, 2, 0),
(73, '退款管理', 'refund-orders', 7, 1, 3, 0),

-- 二级菜单 - 审批管理
(81, '我发起的', 'approval-initiated', 8, 1, 1, 0),
(82, '待我审批', 'approval-pending', 8, 1, 2, 0),
(83, '已审批', 'approval-completed', 8, 1, 3, 0);

-- 插入权限数据
INSERT INTO `permissions` (`id`, `name`, `action_id`, `menu_id`, `status`) VALUES
-- 账号管理权限
(1, '查看账号', 'view_account', 11, 0),
(2, '添加账号', 'add_account', 11, 0),
(3, '编辑账号', 'edit_account', 11, 0),
(4, '删除账号', 'delete_account', 11, 0),

-- 角色管理权限
(11, '查看角色', 'view_role', 12, 0),
(12, '添加角色', 'add_role', 12, 0),
(13, '编辑角色', 'edit_role', 12, 0),
(14, '删除角色', 'delete_role', 12, 0),

-- 权限管理权限
(21, '查看权限', 'view_permission', 13, 0),
(22, '编辑权限', 'edit_permission', 13, 0),

-- 菜单管理权限
(31, '查看菜单', 'view_menu', 14, 0),
(32, '编辑菜单', 'edit_menu', 14, 0),

-- 学生管理权限
(41, '查看学生', 'view_student', 21, 0),
(42, '添加学生', 'add_student', 21, 0),
(43, '编辑学生', 'edit_student', 21, 0),
(44, '删除学生', 'delete_student', 21, 0),

-- 教练管理权限
(51, '查看教练', 'view_coach', 31, 0),
(52, '添加教练', 'add_coach', 31, 0),
(53, '编辑教练', 'edit_coach', 31, 0),
(54, '删除教练', 'delete_coach', 31, 0),

-- 品牌管理权限
(61, '查看品牌', 'view_brand', 41, 0),
(62, '添加品牌', 'add_brand', 41, 0),
(63, '编辑品牌', 'edit_brand', 41, 0),
(64, '删除品牌', 'delete_brand', 41, 0),

-- 分类管理权限
(71, '查看分类', 'view_classify', 42, 0),
(72, '添加分类', 'add_classify', 42, 0),
(73, '编辑分类', 'edit_classify', 42, 0),
(74, '删除分类', 'delete_classify', 42, 0),

-- 属性管理权限
(81, '查看属性', 'view_attribute', 43, 0),
(82, '添加属性', 'add_attribute', 43, 0),
(83, '编辑属性', 'edit_attribute', 43, 0),
(84, '删除属性', 'delete_attribute', 43, 0),

-- 商品管理权限
(91, '查看商品', 'view_goods', 44, 0),
(92, '添加商品', 'add_goods', 44, 0),
(93, '编辑商品', 'edit_goods', 44, 0),
(94, '删除商品', 'delete_goods', 44, 0),

-- 订单管理权限
(101, '查看订单', 'view_order', 51, 0),
(102, '添加订单', 'add_order', 51, 0),
(103, '编辑订单', 'edit_order', 51, 0),
(104, '删除订单', 'delete_order', 51, 0);

-- 插入角色权限关联（普通管理员拥有学生、教练、商品、订单的查看和编辑权限）
INSERT INTO `role_permissions` (`role_id`, `permissions_id`) VALUES
-- 普通管理员权限
(2, 41), (2, 42), (2, 43),  -- 学生查看、添加、编辑
(2, 51), (2, 52), (2, 53),  -- 教练查看、添加、编辑
(2, 61), (2, 62), (2, 63),  -- 品牌查看、添加、编辑
(2, 71), (2, 72), (2, 73),  -- 分类查看、添加、编辑
(2, 81), (2, 82), (2, 83),  -- 属性查看、添加、编辑
(2, 91), (2, 92), (2, 93),  -- 商品查看、添加、编辑
(2, 101), (2, 102), (2, 103),  -- 订单查看、添加、编辑

-- 操作员权限（仅查看）
(3, 41), (3, 51), (3, 61), (3, 71), (3, 81), (3, 91), (3, 101);

-- 输出提示
SELECT '数据库初始化完成！' AS message;
SELECT '测试账号如下：' AS message;
SELECT username, password, name, '角色' AS role_name FROM useraccount JOIN role ON useraccount.role_id = role.id;
