# 设计文档：系统设置迁移与权限种植

## 编辑守卫设计

### 业务背景

角色和菜单是RBAC权限体系的核心配置项。一旦角色启用（status=0），用户已通过该角色获得系统访问权限。此时若允许随意修改角色名称、绑定权限或修改菜单配置，可能导致：
- 权限静默变化，用户未知情失去或获得访问能力
- 审计追踪困难

### 守卫规则

| 操作 | 条件 | Python行为 | 对齐策略 |
|------|------|-----------|---------|
| 编辑角色(name/comment) | status=1(禁用)时可行 | 启用时返回400 | Service层编辑前校验status |
| 修改角色权限 | status=1(禁用)时可行 | 启用时返回400 | Service层修改权限前校验status |
| 编辑菜单(name/sort_order) | status=1(禁用)时可行 | 启用时返回400 | Service层编辑前校验status |

### 实现层次

守卫逻辑放在Service层，Repository层不感知此规则：
```
Handler → Service(校验status) → Repository(执行操作)
```

## sort_order唯一约束设计

菜单编辑时，同parent_id下的sort_order不能重复。校验逻辑：
1. 查询同parent_id下所有菜单（排除当前编辑的菜单本身）
2. 检查是否存在sort_order与目标值相同的记录
3. 冲突时返回400：`"同级菜单中排序已存在"`

## 新建角色默认状态

角色entity中Status字段设置默认值为1（禁用），Service层创建角色时不设置status，让entity默认生效。

## 数据种植脚本设计

**文件位置**：`scripts/seed_permissions.sql`

**设计要点**：
- 使用 `INSERT IGNORE INTO` 语法，遇到主键冲突自动跳过
- 先种植menu表（权限表外键依赖），后种植permissions表
- 包含完整的43条权限和22条菜单记录
- 时间字段使用固定值，便于审计追踪

**执行方式**：
```bash
mysql -u <user> -p <db> < scripts/seed_permissions.sql
```

## 数据关联关系

```
menu (一级)
  └── menu (二级)
        └── permissions (权限条目)
              └── role_permissions (角色-权限关联)
                    └── role (角色)
                          └── useraccount (用户账号)
```

种植脚本仅处理menu和permissions表，role_permissions由运维通过UI配置。
