-- 营销管理模块权限播种脚本

-- 活动模板权限 (menu_id: 22)
INSERT INTO `permissions` (`menu_id`, `action_id`, `name`, `status`) VALUES
(22, 'view_activity_template', '查看活动模板', 0),
(22, 'add_activity_template', '新增活动模板', 0),
(22, 'edit_activity_template', '编辑活动模板', 0),
(22, 'enable_activity_template', '启用活动模板', 0),
(22, 'disable_activity_template', '禁用活动模板', 0);

-- 活动管理权限 (menu_id: 23)
INSERT INTO `permissions` (`menu_id`, `action_id`, `name`, `status`) VALUES
(23, 'view_activity', '查看活动', 0),
(23, 'add_activity', '新增活动', 0),
(23, 'edit_activity', '编辑活动', 0),
(23, 'enable_activity', '启用活动', 0),
(23, 'disable_activity', '禁用活动', 0);
