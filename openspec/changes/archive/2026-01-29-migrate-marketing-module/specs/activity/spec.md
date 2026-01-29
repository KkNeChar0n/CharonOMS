# 活动管理规范

## 新增需求

### 需求：活动查询

系统MUST支持查询活动列表，MUST支持按ID、名称、模板、状态筛选。

- **标识符**: `activity-list`
- **优先级**: P0

#### 场景：查询所有活动
**前置条件**:
- 用户已登录
- 用户拥有 `view_activity` 权限

**操作**:
```http
GET /api/activities
Authorization: Bearer {token}
```

**预期结果**:
- 返回所有活动列表
- 响应包含关联的模板名称和类型
- 响应格式：
  ```json
  {
    "activities": [
      {
        "id": 1,
        "template_id": 1,
        "template_name": "春季促销",
        "template_type": 1,
        "name": "春季大促销活动",
        "start_time": "2026-03-01T00:00:00Z",
        "end_time": "2026-03-31T23:59:59Z",
        "status": 0,
        "create_time": "2026-01-29T10:00:00Z",
        "update_time": "2026-01-29T10:00:00Z"
      }
    ]
  }
  ```

#### 场景：按模板筛选活动
**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/activities?template_id=1
```

**预期结果**:
- 仅返回关联指定模板的活动

#### 场景：按名称模糊搜索
**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/activities?name=促销
```

**预期结果**:
- 返回名称包含"促销"的活动

---

### 需求：活动详情查询

系统MUST支持查询活动的详细信息，MUST包含活动细节配置。

- **标识符**: `activity-detail`
- **优先级**: P0

#### 场景：获取活动详情
**前置条件**:
- 活动存在
- 用户拥有 `view_activity` 权限

**操作**:
```http
GET /api/activities/1
```

**预期结果**:
- 返回活动详细信息
- 包含关联的模板信息
- 包含所有活动细节（参与商品和折扣）
- 响应格式：
  ```json
  {
    "id": 1,
    "template_id": 1,
    "template_name": "春季促销",
    "template_type": 1,
    "name": "春季大促销活动",
    "start_time": "2026-03-01T00:00:00Z",
    "end_time": "2026-03-31T23:59:59Z",
    "status": 0,
    "details": [
      {
        "id": 1,
        "activity_id": 1,
        "goods_id": 101,
        "goods_name": "篮球",
        "discount": 50.00
      },
      {
        "id": 2,
        "activity_id": 1,
        "goods_id": 102,
        "goods_name": "足球",
        "discount": 30.00
      }
    ]
  }
  ```

---

### 需求：活动创建

系统MUST支持创建活动，MUST包括基本信息和活动细节配置。

- **标识符**: `activity-create`
- **优先级**: P0

#### 场景：创建活动
**前置条件**:
- 用户拥有 `add_activity` 权限
- 活动模板存在且状态为启用

**操作**:
```http
POST /api/activities
Content-Type: application/json

{
  "template_id": 1,
  "name": "春季大促销活动",
  "start_time": "2026-03-01T00:00:00Z",
  "end_time": "2026-03-31T23:59:59Z",
  "details": [
    {
      "goods_id": 101,
      "discount": 50.00
    },
    {
      "goods_id": 102,
      "discount": 30.00
    }
  ]
}
```

**预期结果**:
- 活动创建成功
- 所有活动细节创建成功
- 返回创建的活动ID
- 活动状态默认为禁用（status=1）

**验证规则**:
- 活动名称不能为空
- 活动模板MUST存在且状态为启用
- 开始时间MUST早于结束时间
- 活动细节中的商品MUST存在
- 活动细节中的折扣值MUST大于0
- MUST至少包含一个活动细节

#### 场景：创建失败 - 时间范围无效
**前置条件**:
- 用户已登录

**操作**:
```http
POST /api/activities
Content-Type: application/json

{
  "template_id": 1,
  "name": "测试活动",
  "start_time": "2026-03-31T23:59:59Z",
  "end_time": "2026-03-01T00:00:00Z",
  "details": [
    {
      "goods_id": 101,
      "discount": 50.00
    }
  ]
}
```

**预期结果**:
- 返回错误：开始时间MUST早于结束时间
- 活动未创建

#### 场景：创建失败 - 模板未启用
**前置条件**:
- 活动模板状态为禁用

**操作**:
```http
POST /api/activities
Content-Type: application/json

{
  "template_id": 1,
  "name": "测试活动",
  "start_time": "2026-03-01T00:00:00Z",
  "end_time": "2026-03-31T23:59:59Z",
  "details": []
}
```

**预期结果**:
- 返回错误：活动模板未启用
- 活动未创建

---

### 需求：活动更新

系统MUST支持更新活动的基本信息和活动细节，只有禁用状态的活动才能编辑。

- **标识符**: `activity-update`
- **优先级**: P0

#### 场景：更新活动基本信息
**前置条件**:
- 活动存在且状态为禁用
- 用户拥有 `edit_activity` 权限

**操作**:
```http
PUT /api/activities/1
Content-Type: application/json

{
  "template_id": 1,
  "name": "春季大促销活动（延期）",
  "start_time": "2026-03-01T00:00:00Z",
  "end_time": "2026-04-15T23:59:59Z",
  "details": [
    {
      "goods_id": 101,
      "discount": 60.00
    }
  ]
}
```

**预期结果**:
- 活动基本信息更新成功
- 活动细节更新成功（删除旧细节，创建新细节）

#### 场景：更新失败 - 活动已启用
**前置条件**:
- 活动状态为启用（status=0）

**操作**:
```http
PUT /api/activities/1
Content-Type: application/json

{
  "name": "新名称"
}
```

**预期结果**:
- 返回错误：活动启用中，无法编辑
- 活动未更新

---

### 需求：活动状态更新

系统MUST支持更新活动的状态（启用/禁用）。

- **标识符**: `activity-status-update`
- **优先级**: P0

#### 场景：启用活动
**前置条件**:
- 活动存在且状态为禁用
- 用户拥有 `enable_activity` 权限

**操作**:
```http
PUT /api/activities/1/status
Content-Type: application/json

{
  "status": 0
}
```

**预期结果**:
- 活动状态更新为启用

#### 场景：禁用活动
**前置条件**:
- 活动存在且状态为启用
- 用户拥有 `disable_activity` 权限

**操作**:
```http
PUT /api/activities/1/status
Content-Type: application/json

{
  "status": 1
}
```

**预期结果**:
- 活动状态更新为禁用

## 跨功能需求

### 性能
- 活动列表查询MUST在 100ms 内返回（100 条数据以内）
- 活动详情查询（包含细节）MUST在 50ms 内返回

### 安全
- 所有接口MUST需要 JWT 认证
- 创建、编辑、状态更新操作MUST需要权限验证
- 超级管理员拥有所有权限

### 数据一致性
- 创建/更新活动时MUST使用事务，确保活动和细节数据同时成功或同时失败
- 更新细节时MUST先删除旧细节再创建新细节

### 业务约束
- 新创建的活动默认状态为禁用（status=1）
- 只有禁用状态的活动才能编辑
- 开始时间MUST早于结束时间
- 只有启用的活动模板才能用于创建活动
- 活动细节中的折扣值MUST大于0

## 依赖关系
- 依赖 RBAC 权限系统（`openspec/specs/rbac`）
- 依赖活动模板管理（`openspec/changes/migrate-marketing-module/specs/activity-template`）
- 依赖商品管理模块（`openspec/specs/goods`）
