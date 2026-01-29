# 活动模板管理规范

## 新增需求

### 需求：活动模板查询

系统MUST支持查询活动模板列表，MUST支持按ID、名称、类型、状态筛选。

- **标识符**: `activity-template-list`
- **优先级**: P0

#### 场景：查询所有活动模板
**前置条件**:
- 用户已登录
- 用户拥有 `view_activity_template` 权限

**操作**:
```http
GET /api/activity-templates
Authorization: Bearer {token}
```

**预期结果**:
- 返回所有活动模板列表
- 响应格式：
  ```json
  {
    "activity_templates": [
      {
        "id": 1,
        "name": "春季促销满减活动",
        "type": 1,
        "select_type": 1,
        "status": 0,
        "create_time": "2026-01-29T10:00:00Z",
        "update_time": "2026-01-29T10:00:00Z"
      }
    ]
  }
  ```

#### 场景：按活动类型筛选
**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/activity-templates?type=1
```

**预期结果**:
- 仅返回类型为满减（type=1）的活动模板

#### 场景：按名称模糊搜索
**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/activity-templates?name=促销
```

**预期结果**:
- 返回名称包含"促销"的活动模板

---

### 需求：获取启用的活动模板列表

系统MUST提供获取启用状态活动模板的接口，供活动管理使用。

- **标识符**: `activity-template-active-list`
- **优先级**: P0

#### 场景：获取启用的活动模板
**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/activity-templates/active
```

**预期结果**:
- 仅返回状态为启用（status=0）的活动模板
- 用于活动管理页面的模板选择下拉列表

---

### 需求：活动模板详情查询

系统MUST支持查询活动模板的详细信息，MUST包含关联的商品或分类信息。

- **标识符**: `activity-template-detail`
- **优先级**: P0

#### 场景：获取按分类选择的活动模板详情
**前置条件**:
- 活动模板存在，select_type=1（按分类）
- 用户拥有 `view_activity_template` 权限

**操作**:
```http
GET /api/activity-templates/1
```

**预期结果**:
- 返回模板详细信息
- 包含关联的分类列表
- 响应格式：
  ```json
  {
    "id": 1,
    "name": "春季促销满减活动",
    "type": 1,
    "select_type": 1,
    "status": 0,
    "classifies": [
      {
        "id": 1,
        "template_id": 1,
        "classify_id": 101,
        "classify_name": "运动装备"
      }
    ]
  }
  ```

#### 场景：获取按商品选择的活动模板详情
**前置条件**:
- 活动模板存在，select_type=2（按商品）

**操作**:
```http
GET /api/activity-templates/2
```

**预期结果**:
- 返回模板详细信息
- 包含关联的商品列表
- 响应格式：
  ```json
  {
    "id": 2,
    "name": "指定商品满折活动",
    "type": 2,
    "select_type": 2,
    "status": 0,
    "goods": [
      {
        "id": 1,
        "template_id": 2,
        "goods_id": 201,
        "goods_name": "篮球"
      }
    ]
  }
  ```

---

### 需求：活动模板创建

系统MUST支持创建活动模板，MUST包括基本信息和关联配置。

- **标识符**: `activity-template-create`
- **优先级**: P0

#### 场景：创建按分类选择的活动模板
**前置条件**:
- 用户拥有 `add_activity_template` 权限
- 分类存在

**操作**:
```http
POST /api/activity-templates
Content-Type: application/json

{
  "name": "春季促销满减活动",
  "type": 1,
  "select_type": 1,
  "classify_ids": [101, 102]
}
```

**预期结果**:
- 活动模板创建成功
- 关联的分类记录创建成功
- 返回创建的模板ID
- 模板状态默认为禁用（status=1）

**验证规则**:
- 模板名称不能为空
- 活动类型MUST为 1/2/3（满减/满折/满赠）
- 选择方式MUST为 1/2（按分类/按商品）
- 当 select_type=1 时，classify_ids 不能为空
- 当 select_type=2 时，goods_ids 不能为空

#### 场景：创建按商品选择的活动模板
**前置条件**:
- 用户拥有 `add_activity_template` 权限
- 商品存在

**操作**:
```http
POST /api/activity-templates
Content-Type: application/json

{
  "name": "指定商品满折活动",
  "type": 2,
  "select_type": 2,
  "goods_ids": [201, 202, 203]
}
```

**预期结果**:
- 活动模板创建成功
- 关联的商品记录创建成功

#### 场景：创建失败 - 缺少关联配置
**前置条件**:
- 用户已登录

**操作**:
```http
POST /api/activity-templates
Content-Type: application/json

{
  "name": "测试活动",
  "type": 1,
  "select_type": 1
}
```

**预期结果**:
- 返回错误：按分类选择时MUST提供分类ID列表
- 模板未创建

---

### 需求：活动模板更新

系统MUST支持更新活动模板的基本信息和关联配置，只有禁用状态的模板才能编辑。

- **标识符**: `activity-template-update`
- **优先级**: P0

#### 场景：更新活动模板基本信息
**前置条件**:
- 活动模板存在且状态为禁用
- 用户拥有 `edit_activity_template` 权限

**操作**:
```http
PUT /api/activity-templates/1
Content-Type: application/json

{
  "name": "春季大促销满减活动",
  "type": 1,
  "select_type": 1,
  "classify_ids": [101, 102, 103]
}
```

**预期结果**:
- 活动模板基本信息更新成功
- 关联的分类记录更新成功（删除旧关联，创建新关联）

#### 场景：更新失败 - 模板已启用
**前置条件**:
- 活动模板状态为启用（status=0）

**操作**:
```http
PUT /api/activity-templates/1
Content-Type: application/json

{
  "name": "新名称"
}
```

**预期结果**:
- 返回错误：活动模板启用中，无法编辑
- 模板未更新

---

### 需求：活动模板状态更新

系统MUST支持更新活动模板的状态（启用/禁用）。

- **标识符**: `activity-template-status-update`
- **优先级**: P0

#### 场景：启用活动模板
**前置条件**:
- 活动模板存在且状态为禁用
- 用户拥有 `enable_activity_template` 权限

**操作**:
```http
PUT /api/activity-templates/1/status
Content-Type: application/json

{
  "status": 0
}
```

**预期结果**:
- 模板状态更新为启用

#### 场景：禁用活动模板
**前置条件**:
- 活动模板存在且状态为启用
- 用户拥有 `disable_activity_template` 权限

**操作**:
```http
PUT /api/activity-templates/1/status
Content-Type: application/json

{
  "status": 1
}
```

**预期结果**:
- 模板状态更新为禁用

## 跨功能需求

### 性能
- 活动模板列表查询MUST在 100ms 内返回（100 条数据以内）
- 模板详情查询（包含关联）MUST在 50ms 内返回

### 安全
- 所有接口MUST需要 JWT 认证
- 创建、编辑、状态更新操作MUST需要权限验证
- 超级管理员拥有所有权限

### 数据一致性
- 创建/更新活动模板时MUST使用事务，确保模板和关联数据同时成功或同时失败
- 更新关联时MUST先删除旧关联再创建新关联

### 业务约束
- 活动类型枚举：1-满减、2-满折、3-满赠
- 选择方式枚举：1-按分类、2-按商品
- 新创建的模板默认状态为禁用（status=1）
- 只有禁用状态的模板才能编辑

## 依赖关系
- 依赖 RBAC 权限系统（`openspec/specs/rbac`）
- 依赖商品管理模块（`openspec/specs/goods`）
- 依赖分类管理模块（`openspec/specs/classify`）
