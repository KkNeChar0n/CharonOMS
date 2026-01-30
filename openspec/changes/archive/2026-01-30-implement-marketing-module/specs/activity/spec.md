# activity 规范增量

## 修改需求

### 需求： 活动详情查询

活动详情的 `details` 字段MUST仅用于满折活动（type=2）的折扣规则配置。其他类型的活动不使用此字段。

**变更类型**: 修改现有需求 `activity-detail`

- **标识符**: `activity-detail-structure`
- **优先级**: P0
- **关联**: `activity-detail`

#### 场景：获取满折活动详情

**前置条件**:
- 活动存在
- 关联的模板类型为满折（type=2）
- 活动有配置折扣规则

**操作**:
```http
GET /api/activities/1
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 200 OK
- 响应包含 `details` 字段，其中包含折扣规则
- `details` 数组按 `threshold_amount` 升序排序

---

## 新增需求

### 需求： 活动编辑权限控制

系统MUST只允许编辑禁用状态（status=1）的活动。对于启用状态的活动，系统MUST拒绝编辑操作。

- **标识符**: `activity-edit-permission`
- **优先级**: P0

#### 场景：编辑启用状态的活动失败

**前置条件**:
- 活动存在且状态为启用（status=0）
- 用户拥有 `edit_activity` 权限

**操作**:
```http
PUT /api/activities/1
Content-Type: application/json
Authorization: Bearer {token}

{
  "template_id": 1,
  "name": "修改后的名称",
  "start_time": "2026-03-01T00:00:00Z",
  "end_time": "2026-03-31T23:59:59Z",
  "details": []
}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：`{"error": "活动启用中，无法编辑"}`
- 活动未被修改

**验证规则**:
- 在更新前MUST检查活动的 `status` 字段
- 如果 `status = 0`（启用），MUST拒绝更新
- 错误消息MUST为："活动启用中，无法编辑"

---

### 需求： 活动时间范围验证

系统MUST在创建或更新活动时，验证开始时间早于结束时间。

- **标识符**: `activity-time-validation`
- **优先级**: P0

#### 场景：结束时间早于开始时间

**前置条件**:
- 用户拥有 `add_activity` 权限
- 活动模板存在且启用

**操作**:
```http
POST /api/activities
Content-Type: application/json
Authorization: Bearer {token}

{
  "template_id": 1,
  "name": "测试活动",
  "start_time": "2026-03-31T23:59:59Z",
  "end_time": "2026-03-01T00:00:00Z",
  "details": []
}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 活动未创建

**验证规则**:
- 开始时间MUST早于结束时间
- 时间字段不能为空

---

### 需求： 活动模板状态验证

系统MUST在创建活动时，验证关联的活动模板存在且状态为启用。

- **标识符**: `activity-template-status-validation`
- **优先级**: P0

#### 场景：使用禁用的模板创建活动

**前置条件**:
- 活动模板存在但状态为禁用（status=1）
- 用户拥有 `add_activity` 权限

**操作**:
```http
POST /api/activities
Content-Type: application/json
Authorization: Bearer {token}

{
  "template_id": 1,
  "name": "测试活动",
  "start_time": "2026-03-01T00:00:00Z",
  "end_time": "2026-03-31T23:59:59Z",
  "details": []
}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：`{"error": "活动模板未启用"}`
- 活动未创建

**验证规则**:
- 创建前MUST查询模板并检查其状态
- 如果模板不存在或 `status != 0`，MUST拒绝创建

---

### 需求： 按日期范围查询活动

系统MUST支持根据预计付款时间查询在该时间范围内启用的活动，MUST检测同类型活动冲突。

- **标识符**: `activity-query-by-date-range`
- **优先级**: P0

#### 场景：查询时间范围内的活动（无冲突）

**前置条件**:
- 存在多个启用的活动
- 活动时间范围覆盖查询时间
- 各活动类型不重复

**操作**:
```http
GET /api/activities/by-date-range?payment_time=2026-03-15T12:00:00Z
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 200 OK
- 响应包含 `has_duplicate: false` 和活动列表
- 仅满折活动（template_type=2）包含 `details` 字段
- 活动按 `id` 升序排序

**验证规则**:
- 查询条件：`start_time <= payment_time AND end_time >= payment_time AND status = 0`
- MUST检测同一时间范围内相同类型的活动

#### 场景：查询时间范围内的活动（有冲突）

**前置条件**:
- 存在两个相同类型的启用活动
- 活动时间范围覆盖查询时间

**操作**:
```http
GET /api/activities/by-date-range?payment_time=2026-03-15T12:00:00Z
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 200 OK
- 响应格式：`{"has_duplicate": true, "duplicate_type": 1, "type_name": "满减", "activities": []}`
- `type_name` MUST为类型的中文名称：1=满减，2=满折，3=满赠，4=换购

#### 场景：缺少必需参数

**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/activities/by-date-range
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：`{"error": "预计付款时间不能为空"}`

---

### 需求： 活动与订单关联

系统MUST支持订单关联活动，通过 `orders_activity` 表记录关联关系。

- **标识符**: `activity-order-association`
- **优先级**: P1

#### 场景：订单关联活动

**前置条件**:
- 订单存在
- 活动存在且启用

**操作**:
订单创建或更新时，传入活动ID列表：`{"activity_ids": [1, 2]}`

**预期结果**:
- 在 `orders_activity` 表中创建关联记录
- 每个活动ID创建一条记录

**验证规则**:
- 活动ID MUST存在
- 活动MUST为启用状态
- 不允许重复关联（UNIQUE KEY约束）
- 订单删除时，关联记录MUST级联删除

---

### 需求： 活动列表查询增强

系统MUST在查询活动列表时，包含关联的模板信息。

- **标识符**: `activity-list-with-template`
- **优先级**: P0

#### 场景：查询活动列表

**前置条件**:
- 用户已登录
- 用户拥有 `view_activity` 权限

**操作**:
```http
GET /api/activities
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 200 OK
- 每个活动MUST包含：template_name, template_type 字段
- 活动按 `id` 降序排序

**验证规则**:
- MUST使用 JOIN 查询 `activity_template` 表获取模板信息
