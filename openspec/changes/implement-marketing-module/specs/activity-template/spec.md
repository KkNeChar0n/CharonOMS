# activity-template 规范增量

## 修改需求

### 需求：活动模板删除前检查

在删除活动模板时，系统MUST检查是否有关联的活动。如果存在关联活动，系统MUST禁止删除并返回明确的错误信息。

**变更类型**: 修改现有需求 `activity-template-delete`

- **标识符**: `activity-template-delete-validation`
- **优先级**: P0
- **关联**: `activity-template-delete`

#### 场景：删除有关联活动的模板

**前置条件**:
- 活动模板存在
- 存在关联该模板的活动记录
- 用户拥有 `delete_activity_template` 权限

**操作**:
```http
DELETE /api/activity-templates/1
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：
  ```json
  {
    "error": "该模板有关联活动，无法删除"
  }
  ```
- 活动模板未被删除
- 数据库中的数据保持不变

**验证规则**:
- 删除前MUST查询 `activity` 表中 `template_id` 匹配的记录数
- 如果记录数 > 0，MUST拒绝删除
- 错误消息MUST为："该模板有关联活动，无法删除"

---

## 新增需求

### 需求：活动模板编辑权限控制

系统MUST只允许编辑禁用状态（status=1）的活动模板。对于启用状态的模板，系统MUST拒绝编辑操作。

- **标识符**: `activity-template-edit-permission`
- **优先级**: P0

#### 场景：编辑启用状态的模板失败

**前置条件**:
- 活动模板存在且状态为启用（status=0）
- 用户拥有 `edit_activity_template` 权限

**操作**:
```http
PUT /api/activity-templates/1
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "修改后的名称",
  "type": 1,
  "select_type": 1,
  "classify_ids": [101]
}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：
  ```json
  {
    "error": "活动模板启用中，无法编辑"
  }
  ```
- 活动模板未被修改

**验证规则**:
- 在更新前MUST检查模板的 `status` 字段
- 如果 `status = 0`（启用），MUST拒绝更新
- 错误消息MUST为："活动模板启用中，无法编辑"

#### 场景：成功编辑禁用状态的模板

**前置条件**:
- 活动模板存在且状态为禁用（status=1）
- 用户拥有 `edit_activity_template` 权限

**操作**:
```http
PUT /api/activity-templates/1
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "修改后的名称",
  "type": 1,
  "select_type": 1,
  "classify_ids": [101, 102]
}
```

**预期结果**:
- 返回 HTTP 200 OK
- 响应格式：
  ```json
  {
    "message": "活动模板更新成功"
  }
  ```
- 活动模板基本信息已更新
- 旧的关联记录已删除
- 新的关联记录已创建

---

### 需求：活动模板关联配置验证

系统MUST在创建或更新活动模板时，根据 select_type 验证关联配置的完整性。

- **标识符**: `activity-template-relation-validation`
- **优先级**: P0

#### 场景：按分类选择但未提供分类ID

**前置条件**:
- 用户拥有 `add_activity_template` 权限

**操作**:
```http
POST /api/activity-templates
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "测试模板",
  "type": 1,
  "select_type": 1
}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：
  ```json
  {
    "error": "模板名称、活动类型和选择方式不能为空"
  }
  ```
- 活动模板未创建

**验证规则**:
- 当 `select_type = 1` 时，`classify_ids` 数组MUST不为空
- 当 `select_type = 2` 时，`goods_ids` 数组MUST不为空

#### 场景：按商品选择但未提供商品ID

**前置条件**:
- 用户拥有 `add_activity_template` 权限

**操作**:
```http
POST /api/activity-templates
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "测试模板",
  "type": 2,
  "select_type": 2
}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：
  ```json
  {
    "error": "模板名称、活动类型和选择方式不能为空"
  }
  ```
- 活动模板未创建

---

### 需求：活动模板关联数据查询

系统MUST在获取活动模板详情时，根据 select_type 返回对应的关联数据（分类或商品）。

- **标识符**: `activity-template-relation-query`
- **优先级**: P0

#### 场景：查询按分类选择的模板详情

**前置条件**:
- 活动模板存在，select_type=1
- 模板关联了分类

**操作**:
```http
GET /api/activity-templates/1
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 200 OK
- 响应包含 `classify_list` 字段
- 响应不包含 `goods_list` 字段
- 响应格式：
  ```json
  {
    "id": 1,
    "name": "春季促销",
    "type": 1,
    "select_type": 1,
    "status": 0,
    "create_time": "2026-01-30T10:00:00Z",
    "update_time": "2026-01-30T10:00:00Z",
    "classify_list": [
      {
        "classify_id": 101,
        "classify_name": "运动装备"
      }
    ]
  }
  ```

#### 场景：查询按商品选择的模板详情

**前置条件**:
- 活动模板存在，select_type=2
- 模板关联了商品

**操作**:
```http
GET /api/activity-templates/2
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 200 OK
- 响应包含 `goods_list` 字段
- 响应不包含 `classify_list` 字段
- 响应格式：
  ```json
  {
    "id": 2,
    "name": "指定商品促销",
    "type": 2,
    "select_type": 2,
    "status": 0,
    "create_time": "2026-01-30T10:00:00Z",
    "update_time": "2026-01-30T10:00:00Z",
    "goods_list": [
      {
        "goods_id": 201,
        "goods_name": "篮球",
        "price": 120.00,
        "brand_name": "耐克",
        "classify_name": "运动装备"
      }
    ]
  }
  ```

**验证规则**:
- MUST使用 JOIN 查询获取关联数据
- 按分类选择时，MUST查询 `classify` 表获取分类名称
- 按商品选择时，MUST查询 `goods`, `brand`, `classify` 表获取完整商品信息
