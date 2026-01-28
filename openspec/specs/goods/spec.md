# goods 规范

## 目的
待定 - 由归档变更 add-goods-management 创建。归档后请更新目的。
## 需求
### 需求：获取商品列表

系统必须提供获取商品列表的接口，含品牌名称、分类名称和属性字符串，支持分类和状态筛选。

- **接口**：`GET /api/goods`
- **认证**：需要 JWT 认证

#### 场景：获取全部商品列表
- SQL: LEFT JOIN brand + classify 获取 brand_name, classify_name
- 对每条商品额外查询 goods_attributevalue 获取属性信息
- 构建 attributes 字符串：`"属性名:值1/值2,属性名2:值3"`
  - 按 attributeid 分组
  - 同一属性的多个值用 `/` 分隔
  - 不同属性用 `,` 分隔
  - 属性(classify=0)和规格(classify=1)均包含
- 按 id DESC 排序
- 返回 `{"goods": [{id, name, brandid, classifyid, isgroup, price, status, brand_name, classify_name, attributes, attributes_full}, ...]}` 状态码 200

#### 场景：按分类筛选
- 提供 classifyid 参数
- 仅返回该分类下的商品

#### 场景：按状态筛选
- 提供 status 参数（0=启用，1=禁用）
- 仅返回该状态的商品

---

### 需求：新增商品

系统必须提供新增商品的接口，支持单商品和组合商品，验证必填字段和组合商品子商品约束。

- **接口**：`POST /api/goods`
- **认证**：需要 JWT 认证

#### 场景：成功新增单商品
- isgroup=1（或不提供，默认1）
- 插入 goods 记录，status=0
- 若有 attributevalue_ids，插入 goods_attributevalue 关联
- 返回 `{"message": "商品添加成功", "id": int}` 状态码 201

#### 场景：成功新增组合商品
- isgroup=0
- 必须提供非空的 included_goods_ids
- 插入 goods + goods_attributevalue + goods_goods（子商品关联）
- 返回 `{"message": "商品添加成功", "id": int}` 状态码 201

#### 场景：必填字段缺失
- name、brandid、classifyid、price 任一为空
- 返回 `{"error": "商品信息不完整"}` 状态码 400

#### 场景：组合商品无子商品
- isgroup=0 但 included_goods_ids 为空或不提供
- 返回 `{"error": "组合商品必须至少包含一个子商品"}` 状态码 400

---

### 需求：获取商品详情

系统必须提供获取商品详情的接口，返回属性值ID数组和子商品ID数组。

- **接口**：`GET /api/goods/:id`
- **认证**：需要 JWT 认证

#### 场景：获取单商品详情
- 返回商品基本信息 + brand_name + classify_name
- attributevalue_ids：该商品关联的属性值 ID 数组
- included_goods_ids：空数组（单商品无子商品）
- 返回 `{"goods": {id, name, brandid, classifyid, isgroup, price, status, brand_name, classify_name, attributevalue_ids, included_goods_ids}}` 状态码 200

#### 场景：获取组合商品详情
- 同上，但 included_goods_ids 包含子商品 ID 数组（按 goodsid 排序）

#### 场景：商品不存在
- 返回 `{"error": "商品不存在"}` 状态码 404

---

### 需求：编辑商品

系统必须提供编辑商品信息的接口，全量替换属性绑定和子商品关联。

- **接口**：`PUT /api/goods/:id`
- **认证**：需要 JWT 认证

#### 场景：成功编辑单商品
- 更新 name, brandid, classifyid, price
- 全量替换 goods_attributevalue（DELETE + INSERT）
- 返回 `{"message": "商品信息更新成功"}` 状态码 200

#### 场景：成功编辑组合商品
- 同上，额外全量替换 goods_goods（DELETE + INSERT）
- isgroup=0 时必须提供非空的 included_goods_ids

#### 场景：商品不存在
- 返回 `{"error": "商品不存在"}` 状态码 404

#### 场景：必填字段缺失
- 返回 `{"error": "商品信息不完整"}` 状态码 400

---

### 需求：更新商品状态

系统必须提供更新商品启用/禁用状态的接口。

- **接口**：`PUT /api/goods/:id/status`
- **认证**：需要 JWT 认证

#### 场景：成功更新状态
- 返回 `{"message": "商品状态更新成功"}` 状态码 200

#### 场景：状态为空
- 返回 `{"error": "状态不能为空"}` 状态码 400

#### 场景：商品不存在
- 返回 `{"error": "商品不存在"}` 状态码 404

---

### 需求：获取订单可用商品

系统必须提供获取订单可用商品列表的接口，仅返回启用商品，含计算后的总价字段。

- **接口**：`GET /api/goods/active-for-order`
- **认证**：需要 JWT 认证

#### 场景：获取所有启用商品
- 查询 status=0 的所有商品
- 含 brand_name, classify_name
- total_price 计算：
  - 单商品（isgroup=1）：total_price = price
  - 组合商品（isgroup=0）：total_price = SUM(子商品 price)
- 返回 `{"goods": [{..., total_price}, ...]}` 状态码 200

---

### 需求：获取组合商品子商品列表

系统必须提供获取组合商品子商品列表的接口，通过goods_goods关联表查询。

- **接口**：`GET /api/goods/:id/included-goods`
- **认证**：需要 JWT 认证

#### 场景：获取子商品信息
- 通过 goods_goods 表查询 parentsid = :id 的子商品
- 返回子商品的基本信息（id, name, price 等）
- 返回 `{"included_goods": [...]}` 状态码 200

#### 场景：非组合商品或无子商品
- 返回空列表 `{"included_goods": []}`

---

### 需求：计算商品总价

系统必须提供计算商品总价的接口，单商品返回自身价格，组合商品返回子商品价格之和。

- **接口**：`GET /api/goods/:id/total-price`
- **认证**：需要 JWT 认证

#### 场景：计算单商品总价
- isgroup=1
- total_price = price
- 返回 `{"goods_id": int, "price": float, "total_price": float, "isgroup": int}` 状态码 200

#### 场景：计算组合商品总价
- isgroup=0
- total_price = SUM(所有子商品 price)
- 返回同上格式

---

### 需求：获取可用于组合的商品

系统必须提供获取可用于组合的商品列表接口，仅返回启用的单商品。

- **接口**：`GET /api/goods/available-for-combo`
- **认证**：需要 JWT 认证

#### 场景：获取可选子商品列表
- 条件：status=0 AND isgroup=1（启用的单商品）
- 若提供 exclude_id，排除该 ID
- 含 brand_name, classify_name
- 返回 `{"goods": [{id, name, brandid, classifyid, price, brand_name, classify_name}, ...]}` 状态码 200

---

