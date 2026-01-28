# attribute 规范

## 目的
待定 - 由归档变更 add-goods-management 创建。归档后请更新目的。
## 需求
### 需求：获取属性列表

系统必须提供获取所有属性列表的接口，含值数量统计字段。

- **接口**：`GET /api/attributes`
- **认证**：需要 JWT 认证

#### 场景：获取全部属性列表
- SQL: LEFT JOIN attribute_value，GROUP BY 获取 value_count
- 返回字段：id, name, classify, status, create_time, update_time, value_count
- 按 id DESC 排序
- 返回 `{"attributes": [...]}` 状态码 200

---

### 需求：新增属性

系统必须提供新增属性的接口，验证名称和分类字段的合法性。

- **接口**：`POST /api/attributes`
- **认证**：需要 JWT 认证

#### 场景：成功新增属性
- classify=0 表示属性，classify=1 表示规格
- 插入记录，status=0
- 返回 `{"message": "属性创建成功", "attribute_id": int}` 状态码 201

#### 场景：名称或分类为空
- 返回 `{"error": "名称和分类不能为空"}` 状态码 400

#### 场景：分类值非法
- classify 不是 0 或 1
- 返回 `{"error": "分类值必须为0或1"}` 状态码 400

---

### 需求：编辑属性

系统必须提供编辑属性信息的接口。

- **接口**：`PUT /api/attributes/:id`
- **认证**：需要 JWT 认证

#### 场景：成功编辑
- 返回 `{"message": "属性更新成功"}` 状态码 200

#### 场景：属性不存在
- 返回 `{"error": "属性不存在"}` 状态码 404

#### 场景：验证失败
- 同新增场景的验证错误

---

### 需求：更新属性状态

系统必须提供更新属性启用/禁用状态的接口。

- **接口**：`PUT /api/attributes/:id/status`
- **认证**：需要 JWT 认证

#### 场景：成功更新状态
- 返回 `{"message": "状态更新成功"}` 状态码 200

#### 场景：状态值非法
- 返回 `{"error": "状态值必须为0或1"}` 状态码 400

#### 场景：属性不存在
- 返回 `{"error": "属性不存在"}` 状态码 404

---

### 需求：获取属性值列表

系统必须提供获取指定属性的值列表接口。

- **接口**：`GET /api/attributes/:id/values`
- **认证**：需要 JWT 认证

#### 场景：获取指定属性的值
- 按 attributeid 查询 attribute_value 表
- 返回 `{"values": [{id, name, attributeid}, ...]}` 按 id 升序
- 状态码 200

---

### 需求：保存属性值（全量替换）

系统必须提供保存属性值的接口，采用全量替换策略（DELETE + INSERT）。

- **接口**：`POST /api/attributes/:id/values`
- **认证**：需要 JWT 认证

#### 场景：成功保存属性值
- 先 DELETE 该属性下所有旧值
- 再逐条 INSERT 新值
- 返回 `{"message": "属性值保存成功"}` 状态码 200

#### 场景：值数组为空
- 返回 `{"error": "至少需要填入一条属性值"}` 状态码 400

#### 场景：值中有空字符串
- 返回 `{"error": "属性值不能为空"}` 状态码 400

---

### 需求：获取启用属性（含值列表）

系统必须提供获取启用属性及其值列表的接口，用于商品表单的属性/规格选择。

- **接口**：`GET /api/attributes/active`
- **认证**：需要 JWT 认证

#### 场景：获取启用属性及其值
- 查询 status=0 的属性
- 每个属性嵌套其 values 数组
- 返回格式：
```json
{
  "attributes": [
    {
      "id": 1,
      "name": "颜色",
      "classify": 0,
      "values": [
        {"id": 10, "name": "红色"},
        {"id": 11, "name": "蓝色"}
      ]
    }
  ]
}
```
- 状态码 200

---

