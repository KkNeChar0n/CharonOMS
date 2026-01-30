-- 审批流权限初始化脚本
-- 创建日期: 2026-01-30

USE charonoms;

-- 审批流类型管理权限
INSERT INTO permissions (name, comment, status, create_time, update_time) VALUES
('enable_approval_type', '启用审批流类型', 0, NOW(), NOW()),
('disable_approval_type', '禁用审批流类型', 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE comment=VALUES(comment);

-- 审批流模板管理权限
INSERT INTO permissions (name, comment, status, create_time, update_time) VALUES
('view_approval_template', '查看审批流模板', 0, NOW(), NOW()),
('add_approval_template', '新增审批流模板', 0, NOW(), NOW()),
('enable_approval_template', '启用审批流模板', 0, NOW(), NOW()),
('disable_approval_template', '禁用审批流模板', 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE comment=VALUES(comment);

-- 审批流实例管理权限
INSERT INTO permissions (name, comment, status, create_time, update_time) VALUES
('view_approval_flow', '查看审批流', 0, NOW(), NOW()),
('create_approval_flow', '创建审批流', 0, NOW(), NOW()),
('cancel_approval_flow', '撤销审批流', 0, NOW(), NOW()),
('approve_flow', '审批通过', 0, NOW(), NOW()),
('reject_flow', '审批驳回', 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE comment=VALUES(comment);

-- 查询已创建的权限ID
SELECT id, name, comment FROM permissions
WHERE name IN (
    'enable_approval_type', 'disable_approval_type',
    'view_approval_template', 'add_approval_template', 'enable_approval_template', 'disable_approval_template',
    'view_approval_flow', 'create_approval_flow', 'cancel_approval_flow', 'approve_flow', 'reject_flow'
)
ORDER BY id;

-- 说明：
-- 如需将这些权限分配给角色，请执行：
-- INSERT INTO role_permissions (role_id, permission_id) VALUES ([role_id], [permission_id]);
-- 其中role_id可以从role表中查询：SELECT id, name FROM role;
