# 规范增量：分类管理 (Classify)

## 新增需求

### 需求：获取分类列表

系统必须提供获取所有分类列表的接口，含父级名称信息，按级别和ID排序。

- **接口**：`GET /api/classifies`
- **认证**：需要 JWT 认证

#### 场景：获取全部分类列表
- 调用 `GET /api/classifies`
- SQL: LEFT JOIN classify 获取 parent_name
- 返回字段：id, name, level, parent_id, parent_name, status
- 按 level ASC, id DESC 排序
- 返回 `{"classifies": [...]}` 状态码 200

#### 场景：一级分类（无父级）
- level=0 的分类，parent_id 和 parent_name 均为 null

---

### 需求：获取父级分类

系统必须提供获取启用一级分类列表的接口，用于二级分类的父级选择下拉。

- **接口**：`GET /api/classifies/parents`
- **认证**：需要 JWT 认证

#### 场景：获取启用的一级分类
- 查询 status=0 AND level=0
- 返回 `{"parents": [{id, name}, ...]}` 按 id 升序
- 状态码 200

---

### 需求：新增分类

系统必须提供新增分类的接口，验证名称、级别和层级唯一性约束。

- **接口**：`POST /api/classifies`
- **认证**：需要 JWT 认证

#### 场景：成功新增一级分类
- level=0，parent_id 可为空
- 名称在一级分类中唯一
- 插入记录，status=0
- 返回 `{"message": "类型创建成功", "classify_id": int}` 状态码 201

#### 场景：成功新增二级分类
- level=1，必须包含 parent_id
- 名称在同一 parentid 下唯一
- 返回 `{"message": "类型创建成功", "classify_id": int}` 状态码 201

#### 场景：名称或级别为空
- 返回 `{"error": "名称和级别不能为空"}` 状态码 400

#### 场景：级别值非法
- level 不是 0 或 1
- 返回 `{"error": "级别值必须为0或1"}` 状态码 400

#### 场景：二级分类无父级
- level=1 但缺少 parent_id
- 返回 `{"error": "二级类型必须选择父级类型"}` 状态码 400

#### 场景：一级分类名称重复
- 返回 `{"error": "该一级类型名称已存在"}` 状态码 400

#### 场景：同父级二级分类名称重复
- 返回 `{"error": "该父级类型下已存在同名的二级类型"}` 状态码 400

---

### 需求：编辑分类

系统必须提供编辑分类信息的接口，验证规则同新增，排除自身ID。

- **接口**：`PUT /api/classifies/:id`
- **认证**：需要 JWT 认证

#### 场景：成功编辑
- 验证规则同新增，排除自身 ID
- 返回 `{"message": "类型更新成功"}` 状态码 200

#### 场景：分类不存在
- 返回 `{"error": "类型不存在"}` 状态码 404

#### 场景：其他验证失败
- 同新增场景的错误响应

---

### 需求：更新分类状态

系统必须提供更新分类启用/禁用状态的接口。

- **接口**：`PUT /api/classifies/:id/status`
- **认证**：需要 JWT 认证

#### 场景：成功更新状态
- 返回 `{"message": "状态更新成功"}` 状态码 200

#### 场景：状态值非法
- status 不是 0 或 1
- 返回 `{"error": "状态值必须为0或1"}` 状态码 400

#### 场景：分类不存在
- 返回 `{"error": "类型不存在"}` 状态码 404

---

### 需求：获取启用分类

系统必须提供获取启用二级分类列表的接口，用于商品表单下拉选择。

- **接口**：`GET /api/classifies/active`
- **认证**：需要 JWT 认证

#### 场景：获取启用的二级分类
- 查询 status=0 AND level=1
- 返回 `{"classifies": [{id, name}, ...]}` 按 id 升序
- 状态码 200

---

## 数据库表

**classify**
| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 主键，自增 |
| name | varchar(50) | 分类名称 |
| level | int | 级别：0=一级, 1=二级 |
| parentid | int(nullable) | 父级分类ID（level=1时必填） |
| status | int | 0=启用, 1=禁用 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

> 注意：classify 表的时间字段为 `created_at`/`updated_at`，与其他表（`create_time`/`update_time`）不同。
