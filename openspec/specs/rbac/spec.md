# Specification: Role-Based Access Control (RBAC)

基于角色的访问控制模块，提供角色、权限和菜单的管理功能。

## Requirements

### Requirement: 获取菜单列表
系统 SHALL 根据用户角色返回可访问的菜单列表。

#### Scenario: 超级管理员获取所有菜单
- **WHEN** 超级管理员请求菜单列表
- **THEN** 系统返回所有启用状态的菜单（树形结构）
- **AND** 响应格式为 `{"code": 0, "message": "success", "data": {"menus": [...]}}`

#### Scenario: 普通用户获取授权菜单
- **WHEN** 非超级管理员用户请求菜单列表
- **THEN** 系统返回该用户角色关联的权限对应的菜单
- **AND** 菜单按树形结构组织（parent_id关联）

#### Scenario: 未登录用户访问
- **WHEN** 未提供有效JWT token的用户请求菜单
- **THEN** 系统返回401错误
- **AND** 错误消息为 "未授权"

#### Scenario: 菜单树形结构
- **WHEN** 返回菜单列表
- **THEN** 顶层菜单的parent_id为NULL
- **AND** 子菜单的parent_id指向父菜单的id
- **AND** 菜单按sort_order字段排序

### Requirement: 获取角色列表
系统 SHALL 允许查看所有角色信息。

#### Scenario: 获取角色列表成功
- **WHEN** 已认证用户请求角色列表
- **THEN** 系统返回所有角色（包括ID、名称、备注、状态）
- **AND** 响应格式为 `{"code": 0, "message": "success", "data": {"roles": [...]}}`

#### Scenario: 角色包含权限数量
- **WHEN** 返回角色列表
- **THEN** 每个角色包含关联的权限数量统计

### Requirement: 创建角色
系统 SHALL 允许创建新角色。

#### Scenario: 创建角色成功
- **WHEN** 提供有效的角色信息（名称、备注）
- **THEN** 系统创建角色并返回角色ID
- **AND** 响应格式为 `{"code": 0, "message": "角色创建成功", "data": {"role_id": <id>}}`

#### Scenario: 角色名称重复
- **WHEN** 创建的角色名称已存在
- **THEN** 系统返回400错误
- **AND** 错误消息为 "角色名称已存在"

#### Scenario: 缺少必填字段
- **WHEN** 请求缺少角色名称
- **THEN** 系统返回400错误
- **AND** 错误消息提示缺失的字段

### Requirement: 更新角色信息
系统 SHALL 允许更新角色的基本信息。

#### Scenario: 更新角色成功
- **WHEN** 提供有效的角色ID和更新信息
- **THEN** 系统更新角色信息
- **AND** 响应格式为 `{"code": 0, "message": "角色更新成功"}`

#### Scenario: 角色不存在
- **WHEN** 提供的角色ID不存在
- **THEN** 系统返回404错误
- **AND** 错误消息为 "角色不存在"

#### Scenario: 不允许修改超级管理员角色
- **WHEN** 尝试修改is_super_admin=1的角色
- **THEN** 系统返回403错误
- **AND** 错误消息为 "不允许修改超级管理员角色"

### Requirement: 更新角色状态
系统 SHALL 允许启用或禁用角色。

#### Scenario: 禁用角色成功
- **WHEN** 将角色status设置为1（禁用）
- **THEN** 系统更新角色状态
- **AND** 该角色的用户在下次同步角色时检测到变更

#### Scenario: 启用角色成功
- **WHEN** 将角色status设置为0（启用）
- **THEN** 系统更新角色状态

#### Scenario: 状态值无效
- **WHEN** 提供的status值不是0或1
- **THEN** 系统返回400错误
- **AND** 错误消息为 "状态值无效"

### Requirement: 获取角色权限
系统 SHALL 允许查看角色关联的权限列表。

#### Scenario: 获取角色权限成功
- **WHEN** 提供有效的角色ID
- **THEN** 系统返回该角色关联的所有权限
- **AND** 权限按菜单分组，包含菜单信息和操作权限

#### Scenario: 角色无权限
- **WHEN** 角色未关联任何权限
- **THEN** 系统返回空权限列表

### Requirement: 更新角色权限
系统 SHALL 允许批量更新角色的权限绑定。

#### Scenario: 更新权限成功
- **WHEN** 提供角色ID和权限ID列表
- **THEN** 系统删除该角色的所有现有权限
- **AND** 系统创建新的角色-权限关联
- **AND** 响应格式为 `{"code": 0, "message": "权限更新成功"}`

#### Scenario: 权限ID无效
- **WHEN** 提供的权限ID列表包含不存在的ID
- **THEN** 系统返回400错误
- **AND** 错误消息为 "权限ID无效"

#### Scenario: 清空角色权限
- **WHEN** 提供空的权限ID列表
- **THEN** 系统删除该角色的所有权限关联

### Requirement: 获取权限列表
系统 SHALL 允许查看所有权限信息。

#### Scenario: 获取所有权限
- **WHEN** 请求权限列表（不带过滤条件）
- **THEN** 系统返回所有权限

#### Scenario: 按状态过滤权限
- **WHEN** 请求权限列表时指定status参数
- **THEN** 系统返回指定状态的权限
- **AND** 支持URL参数 `?status=0`（启用）或 `?status=1`（禁用）

#### Scenario: 权限包含菜单信息
- **WHEN** 返回权限列表
- **THEN** 每个权限包含关联的菜单名称和路由信息

### Requirement: 更新权限状态
系统 SHALL 允许启用或禁用权限。

#### Scenario: 禁用权限成功
- **WHEN** 将权限status设置为1（禁用）
- **THEN** 系统更新权限状态
- **AND** 拥有该权限的角色用户将失去对应操作权限

#### Scenario: 启用权限成功
- **WHEN** 将权限status设置为0（启用）
- **THEN** 系统更新权限状态

### Requirement: 获取权限树
系统 SHALL 以树形结构返回权限列表，按菜单层级组织。

#### Scenario: 返回权限树成功
- **WHEN** 请求权限树
- **THEN** 系统返回按菜单层级组织的权限树
- **AND** 每个菜单节点包含其下的所有权限

#### Scenario: 权限树包含菜单层级
- **WHEN** 返回权限树
- **THEN** 顶层为一级菜单
- **AND** 子菜单嵌套在父菜单下
- **AND** 权限嵌套在对应的菜单下

### Requirement: 获取菜单管理列表
系统 SHALL 允许查看所有菜单（用于管理界面）。

#### Scenario: 获取所有菜单
- **WHEN** 请求菜单管理列表
- **THEN** 系统返回所有菜单（包括启用和禁用状态）
- **AND** 菜单按树形结构组织

#### Scenario: 菜单包含完整信息
- **WHEN** 返回菜单管理列表
- **THEN** 每个菜单包含ID、名称、路由、父菜单ID、排序、状态

### Requirement: 更新菜单信息
系统 SHALL 允许更新菜单的基本信息。

#### Scenario: 更新菜单成功
- **WHEN** 提供有效的菜单ID和更新信息
- **THEN** 系统更新菜单信息
- **AND** 响应格式为 `{"code": 0, "message": "菜单更新成功"}`

#### Scenario: 菜单不存在
- **WHEN** 提供的菜单ID不存在
- **THEN** 系统返回404错误
- **AND** 错误消息为 "菜单不存在"

#### Scenario: 父菜单ID无效
- **WHEN** 设置的parent_id不存在（除NULL外）
- **THEN** 系统返回400错误
- **AND** 错误消息为 "父菜单不存在"

### Requirement: 更新菜单状态
系统 SHALL 允许启用或禁用菜单。

#### Scenario: 禁用菜单成功
- **WHEN** 将菜单status设置为1（禁用）
- **THEN** 系统更新菜单状态
- **AND** 该菜单不再出现在用户的菜单列表中

#### Scenario: 启用菜单成功
- **WHEN** 将菜单status设置为0（启用）
- **THEN** 系统更新菜单状态

#### Scenario: 禁用父菜单影响子菜单
- **WHEN** 禁用父菜单
- **THEN** 其所有子菜单也不再显示（即使子菜单status为0）

### Requirement: 超级管理员权限
系统 SHALL 为超级管理员提供无限制访问权限。

#### Scenario: 超级管理员跳过权限检查
- **WHEN** 用户的is_super_admin标识为1
- **THEN** 系统跳过所有权限验证
- **AND** 允许访问所有功能和数据

#### Scenario: 超级管理员获取所有菜单
- **WHEN** 超级管理员请求菜单列表
- **THEN** 系统返回所有启用菜单，不检查权限关联

## API Endpoints

### GET /api/menus
获取用户菜单列表（需要认证）

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "menus": [
      {
        "id": 1,
        "name": "系统管理",
        "route": "",
        "parent_id": null,
        "sort_order": 1,
        "children": [
          {
            "id": 2,
            "name": "用户管理",
            "route": "/users",
            "parent_id": 1,
            "sort_order": 1
          }
        ]
      }
    ]
  }
}
```

### GET /api/roles
获取角色列表（需要认证）

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "roles": [
      {
        "id": 1,
        "name": "超级管理员",
        "comment": "拥有系统所有权限",
        "is_super_admin": 1,
        "status": 0,
        "permission_count": 0
      }
    ]
  }
}
```

### POST /api/roles
创建角色（需要认证）

**Request:**
```json
{
  "name": "教务管理员",
  "comment": "管理学生和教练信息"
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "角色创建成功",
  "data": {
    "role_id": 2
  }
}
```

### PUT /api/roles/:id
更新角色信息（需要认证）

**Request:**
```json
{
  "name": "教务管理员",
  "comment": "管理学生、教练和订单信息"
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "角色更新成功"
}
```

### PUT /api/roles/:id/status
更新角色状态（需要认证）

**Request:**
```json
{
  "status": 1
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "状态更新成功"
}
```

### GET /api/roles/:id/permissions
获取角色权限（需要认证）

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "permissions": [
      {
        "id": 1,
        "name": "查看用户",
        "menu_id": 2,
        "menu_name": "用户管理",
        "action_id": 1,
        "status": 0
      }
    ]
  }
}
```

### PUT /api/roles/:id/permissions
更新角色权限（需要认证）

**Request:**
```json
{
  "permission_ids": [1, 2, 3]
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "权限更新成功"
}
```

### GET /api/permissions
获取权限列表（需要认证）

**Query Parameters:**
- `status`: 可选，过滤状态（0=启用，1=禁用）

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "permissions": [
      {
        "id": 1,
        "name": "查看用户",
        "menu_id": 2,
        "menu_name": "用户管理",
        "menu_route": "/users",
        "action_id": 1,
        "status": 0
      }
    ]
  }
}
```

### PUT /api/permissions/:id/status
更新权限状态（需要认证）

**Request:**
```json
{
  "status": 1
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "状态更新成功"
}
```

### GET /api/permissions/tree
获取权限树（需要认证）

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tree": [
      {
        "id": 1,
        "name": "系统管理",
        "children": [
          {
            "id": 2,
            "name": "用户管理",
            "permissions": [
              {"id": 1, "name": "查看用户"},
              {"id": 2, "name": "编辑用户"}
            ]
          }
        ]
      }
    ]
  }
}
```

### GET /api/menu-management
获取菜单管理列表（需要认证）

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "menus": [
      {
        "id": 1,
        "name": "系统管理",
        "route": "",
        "parent_id": null,
        "sort_order": 1,
        "status": 0
      }
    ]
  }
}
```

### PUT /api/menu-management/:id
更新菜单信息（需要认证）

**Request:**
```json
{
  "name": "系统设置",
  "route": "/settings",
  "parent_id": null,
  "sort_order": 1
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "菜单更新成功"
}
```

### PUT /api/menu-management/:id/status
更新菜单状态（需要认证）

**Request:**
```json
{
  "status": 1
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "状态更新成功"
}
```

## Data Models

### Role
```go
type Role struct {
    ID           uint      `gorm:"column:id;primaryKey"`
    Name         string    `gorm:"column:name;size:50;not null"`
    Comment      string    `gorm:"column:comment;size:200"`
    Status       int       `gorm:"column:status;default:0"`
    IsSuperAdmin int       `gorm:"column:is_super_admin;default:0"`
    CreateTime   time.Time `gorm:"column:create_time;autoCreateTime"`
    UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime"`
}
```

### Menu
```go
type Menu struct {
    ID        uint   `gorm:"column:id;primaryKey"`
    Name      string `gorm:"column:name;size:50;not null"`
    Route     string `gorm:"column:route;size:100"`
    ParentID  *uint  `gorm:"column:parent_id"`
    SortOrder int    `gorm:"column:sort_order;default:0"`
    Status    int    `gorm:"column:status;default:0"`
}
```

### Permission
```go
type Permission struct {
    ID         uint      `gorm:"column:id;primaryKey"`
    Name       string    `gorm:"column:name;size:50;not null"`
    MenuID     uint      `gorm:"column:menu_id;not null"`
    ActionID   uint      `gorm:"column:action_id;not null"`
    Status     int       `gorm:"column:status;default:0"`
    CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
    UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}
```

### RolePermission
```go
type RolePermission struct {
    ID           uint      `gorm:"column:id;primaryKey"`
    RoleID       uint      `gorm:"column:role_id;not null"`
    PermissionID uint      `gorm:"column:permissions_id;not null"`
    CreateTime   time.Time `gorm:"column:create_time;autoCreateTime"`
}
```

## Business Rules

1. **菜单层级**: 顶层菜单parent_id为NULL，不能使用0
2. **状态字段**: 0=启用，1=禁用
3. **超级管理员**: is_super_admin=1的角色跳过所有权限检查
4. **角色-权限绑定**: 通过role_permissions中间表实现多对多关系
5. **菜单显示**: 只显示启用状态（status=0）的菜单
6. **权限过滤**: 普通用户只能看到其角色关联的权限对应的菜单
7. **级联影响**: 禁用父菜单会导致子菜单也不显示

## Security Considerations

1. **权限验证**: 所有RBAC管理接口都需要JWT认证
2. **超级管理员保护**: 不允许修改或删除超级管理员角色
3. **数据隔离**: 用户只能看到自己权限范围内的数据
4. **审计日志**: 建议记录角色和权限的变更操作（未来功能）
