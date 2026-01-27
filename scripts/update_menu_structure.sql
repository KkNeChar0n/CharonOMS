-- 更新菜单结构以匹配原项目 ZhixinStudentSaaS
-- 执行前请备份数据库！

USE charonoms;

-- ========================================
-- 第一步：更新一级菜单的名称和顺序
-- ========================================

-- 更新"系统管理"为"系统设置"，并调整到最后
UPDATE menu SET name = '系统设置', sort_order = 99 WHERE id = 1;

-- 更新"学生管理"的顺序
UPDATE menu SET sort_order = 1 WHERE id = 2;

-- 更新"教练管理"的顺序
UPDATE menu SET sort_order = 2 WHERE id = 3;

-- 更新"订单管理"的顺序
UPDATE menu SET sort_order = 3 WHERE id = 5;

-- 更新"商品管理"的顺序
UPDATE menu SET sort_order = 4 WHERE id = 4;

-- 更新"活动管理"为"营销管理"，并调整顺序
UPDATE menu SET name = '营销管理', sort_order = 6 WHERE id = 6;

-- 更新"财务管理"的顺序
UPDATE menu SET sort_order = 8 WHERE id = 7;

-- 更新"审批管理"为"审批流管理"，并调整顺序
UPDATE menu SET name = '审批流管理', sort_order = 5 WHERE id = 8;
