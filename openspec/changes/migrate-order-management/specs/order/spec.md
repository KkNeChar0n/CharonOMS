# order 规范增量

## ADDED Requirements

### 需求：获取订单列表

系统必须提供获取订单列表的接口，包含学生信息、金额、状态等完整信息。

- **接口**：`GET /api/orders`
- **认证**：需要 JWT 认证

#### 场景：获取全部订单列表

- **WHEN** 用户请求订单列表
- **THEN** 系统SHALL：
  - 查询orders表并LEFT JOIN student表获取学生姓名
  - 返回字段：id, uid(student_id), student_name, expected_payment_time, amount_receivable, discount_amount, amount_received, create_time, status
  - 按create_time DESC排序
  - 返回格式：`{"orders": [{id, uid, student_name, expected_payment_time, amount_receivable, discount_amount, amount_received, create_time, status}]}`
  - HTTP状态码：200

---

### 需求：创建订单

系统必须提供创建订单的接口，支持多商品、优惠分摊、活动关联。

- **接口**：`POST /api/orders`
- **认证**：需要 JWT 认证
- **请求体**：
  ```json
  {
    "student_id": int,
    "goods_list": [{"goods_id": int, "total_price": float, "price": float}],
    "expected_payment_time": "YYYY-MM-DD",
    "activity_ids": [int],
    "discount_amount": float,
    "child_discounts": {"goods_id": float}
  }
  ```

#### 场景：成功创建订单

- **WHEN** 用户提交有效的创建订单请求
- **THEN** 系统SHALL：
  - 验证student_id非空
  - 验证goods_list至少包含一个商品
  - 计算应收金额 = SUM(goods.total_price)
  - 计算实收金额 = SUM(goods.price) - discount_amount
  - 插入orders记录，状态设为10（草稿）
  - 为每个商品创建childorders记录，子订单状态设为0（初始）
  - 每个子订单的实收金额 = goods.price - child_discounts[goods_id]
  - 如果提供activity_ids，插入orders_activity关联记录
  - 使用数据库事务保证原子性
  - 返回：`{"message": "订单创建成功", "order_id": int}`
  - HTTP状态码：201

#### 场景：学生ID缺失

- **WHEN** 请求中student_id为空
- **THEN** 系统SHALL：
  - 返回：`{"error": "学生ID不能为空"}`
  - HTTP状态码：400

#### 场景：商品列表为空

- **WHEN** 请求中goods_list为空或不存在
- **THEN** 系统SHALL：
  - 返回：`{"error": "必须至少选择一个商品"}`
  - HTTP状态码：400

---

### 需求：获取订单商品列表

系统必须提供获取订单中商品详情的接口，包含商品名称、品牌、分类、属性等信息。

- **接口**：`GET /api/orders/:id/goods`
- **认证**：需要 JWT 认证

#### 场景：获取订单商品

- **WHEN** 用户请求指定订单的商品列表
- **THEN** 系统SHALL：
  - 查询childorders表
  - LEFT JOIN goods表获取商品信息
  - LEFT JOIN brand表获取品牌名称
  - LEFT JOIN classify表获取分类名称
  - 对每个商品查询goods_attributevalue关联的属性信息
  - 构建attributes数组：`[{"attr_name": "属性名", "value_name": "值"}]`
  - 返回格式：`{"goods": [{id, goodsid, goods_name, isgroup, brand_name, classify_name, amount_receivable, discount_amount, amount_received, attributes}]}`
  - 按childorders.id升序排序
  - HTTP状态码：200

---

### 需求：编辑订单

系统必须提供编辑订单的接口，仅允许编辑草稿状态的订单。

- **接口**：`PUT /api/orders/:id`
- **认证**：需要 JWT 认证
- **请求体**：同创建订单请求

#### 场景：成功编辑草稿订单

- **WHEN** 用户提交编辑请求，且订单状态为10（草稿）
- **THEN** 系统SHALL：
  - 验证订单存在
  - 验证订单状态为10（草稿）
  - 重新计算应收金额和实收金额
  - 更新orders表的金额、预计付款时间、优惠金额
  - 删除旧的childorders记录（WHERE parentsid = order_id）
  - 删除旧的orders_activity记录（WHERE orders_id = order_id）
  - 创建新的childorders和orders_activity记录
  - 使用数据库事务保证原子性
  - 返回：`{"message": "订单更新成功"}`
  - HTTP状态码：200

#### 场景：订单不存在

- **WHEN** 请求的订单ID不存在
- **THEN** 系统SHALL：
  - 返回：`{"error": "订单不存在"}`
  - HTTP状态码：404

#### 场景：订单不是草稿状态

- **WHEN** 订单状态不为10（草稿）
- **THEN** 系统SHALL：
  - 返回：`{"error": "只能编辑草稿状态的订单"}`
  - HTTP状态码：400

---

### 需求：提交订单

系统必须提供提交订单的接口，将草稿订单变更为未支付状态。

- **接口**：`PUT /api/orders/:id/submit`
- **认证**：需要 JWT 认证

#### 场景：成功提交草稿订单

- **WHEN** 用户提交草稿状态的订单
- **THEN** 系统SHALL：
  - 验证订单存在
  - 验证订单状态为10（草稿）
  - 更新订单状态为20（未支付）
  - 更新所有关联子订单状态为10（未支付）
  - 使用数据库事务保证原子性
  - 返回：`{"message": "订单提交成功"}`
  - HTTP状态码：200

#### 场景：订单不是草稿状态

- **WHEN** 订单状态不为10（草稿）
- **THEN** 系统SHALL：
  - 返回：`{"error": "只能提交草稿状态的订单"}`
  - HTTP状态码：400

---

### 需求：作废订单

系统必须提供作废订单的接口，将草稿订单标记为已作废。

- **接口**：`PUT /api/orders/:id/cancel`
- **认证**：需要 JWT 认证

#### 场景：成功作废草稿订单

- **WHEN** 用户作废草稿状态的订单
- **THEN** 系统SHALL：
  - 验证订单存在
  - 验证订单状态为10（草稿）
  - 更新订单状态为99（已作废）
  - 更新所有关联子订单状态为99（已作废）
  - 使用数据库事务保证原子性
  - 返回：`{"message": "订单已作废"}`
  - HTTP状态码：200

#### 场景：订单不是草稿状态

- **WHEN** 订单状态不为10（草稿）
- **THEN** 系统SHALL：
  - 返回：`{"error": "只能作废草稿状态的订单"}`
  - HTTP状态码：400

---

### 需求：获取子订单列表

系统必须提供获取子订单列表的接口，包含商品信息。

- **接口**：`GET /api/childorders`
- **认证**：需要 JWT 认证

#### 场景：获取全部子订单列表

- **WHEN** 用户请求子订单列表
- **THEN** 系统SHALL：
  - 查询childorders表
  - LEFT JOIN goods表获取商品名称
  - 返回字段：id, parentsid, goodsid, goods_name, amount_receivable, discount_amount, amount_received, status, create_time
  - 按id DESC排序
  - 返回格式：`{"childorders": [{id, parentsid, goodsid, goods_name, amount_receivable, discount_amount, amount_received, status, create_time}]}`
  - HTTP状态码：200

---

### 需求：获取启用商品列表（用于订单）

系统必须提供获取启用状态商品的接口，计算组合商品的总价，用于订单选择商品。

- **接口**：`GET /api/goods/active-for-order`
- **认证**：需要 JWT 认证

#### 场景：获取启用商品

- **WHEN** 用户请求启用商品列表
- **THEN** 系统SHALL：
  - 查询status=0（启用）的商品
  - LEFT JOIN brand和classify获取名称
  - 对每个商品查询属性信息并构建attributes字符串
  - 计算total_price：
    - 如果isgroup=1（单商品），total_price = price
    - 如果isgroup=0（组合商品），total_price = SUM(子商品price)，通过goods_goods表查询
  - 返回字段：id, name, brandid, classifyid, isgroup, price, brand_name, classify_name, attributes, total_price
  - 返回格式：`{"goods": [...]}`
  - HTTP状态码：200

---

### 需求：获取商品总价

系统必须提供获取单个商品总价的接口，用于订单金额计算。

- **接口**：`GET /api/goods/:id/total-price`
- **认证**：需要 JWT 认证

#### 场景：获取单商品总价

- **WHEN** 用户请求单商品的总价
- **THEN** 系统SHALL：
  - 查询商品基本信息
  - 如果isgroup=1，total_price = price
  - 返回：`{"goods_id": int, "price": float, "total_price": float, "isgroup": int}`
  - HTTP状态码：200

#### 场景：获取组合商品总价

- **WHEN** 用户请求组合商品的总价
- **THEN** 系统SHALL：
  - 查询商品基本信息
  - 如果isgroup=0，通过goods_goods表查询所有子商品
  - 计算total_price = SUM(子商品price)
  - 返回：`{"goods_id": int, "price": float, "total_price": float, "isgroup": int}`
  - HTTP状态码：200

#### 场景：商品不存在

- **WHEN** 请求的商品ID不存在
- **THEN** 系统SHALL：
  - 返回：`{"error": "商品不存在"}`
  - HTTP状态码：404

---

### 需求：计算订单优惠

系统必须提供计算订单优惠的接口，基于活动规则计算满折优惠并分摊到子订单。

- **接口**：`POST /api/orders/calculate-discount`
- **认证**：需要 JWT 认证
- **请求体**：
  ```json
  {
    "goods_list": [{"goods_id": int, "price": float, "total_price": float}],
    "activity_ids": [int]
  }
  ```

#### 场景：计算满折优惠

- **WHEN** 用户提交优惠计算请求
- **THEN** 系统SHALL：
  - 查询activity_ids对应的活动模板
  - 筛选type=2（满折）的活动
  - 对每个活动：
    - 如果select_type=2（按商品），查询activity_template_goods获取参与商品列表
    - 如果select_type≠2（按分类），查询goods表获取商品的classifyid，匹配activity_template_goods中的分类
    - 统计参与该活动的商品数量
    - 查询activity_template_goods中threshold_amount（档位）和discount_percentage（折扣率）
    - 匹配最高满足的档位
    - 计算优惠金额 = 参与商品标准售价之和 × (1 - 折扣率)
  - 合并所有活动的优惠金额
  - 将总优惠平均分摊到各商品：child_discounts[goods_id] = 总优惠 / 商品数量
  - 返回：`{"total_discount": float, "child_discounts": {"goods_id": float}}`
  - HTTP状态码：200

#### 场景：无满足条件的活动

- **WHEN** 没有活动满足条件或activity_ids为空
- **THEN** 系统SHALL：
  - 返回：`{"total_discount": 0, "child_discounts": {}}`
  - HTTP状态码：200

---

### 需求：订单状态流转规则

系统必须严格遵循订单状态流转规则，保证业务一致性。

#### 场景：订单状态定义

- **订单状态码**：
  - 10：草稿（Draft）- 订单初始状态
  - 20：未支付（Unpaid）- 订单已提交
  - 30：部分支付（PartialPaid）- 部分收款
  - 40：已支付（Paid）- 全额收款
  - 50：退费中（Refunding）- 退款申请已提交
  - 99：已作废（Cancelled）- 订单已取消

#### 场景：子订单状态定义

- **子订单状态码**：
  - 0：初始（Initial）- 子订单创建时的初始状态
  - 10：未支付（Unpaid）- 分账金额为0
  - 20：部分支付（PartialPaid）- 0 < 分账金额 < 实收金额
  - 30：已支付（Paid）- 分账金额 >= 实收金额
  - 99：已作废（Cancelled）- 关联订单已作废

#### 场景：允许的状态流转

- **WHEN** 订单状态为10（草稿）
- **THEN** 系统SHALL允许：
  - 编辑订单（更新商品、金额、优惠）
  - 提交订单（变更为20-未支付）
  - 作废订单（变更为99-已作废）

- **WHEN** 订单状态为20/30/40（未支付/部分支付/已支付）
- **THEN** 系统SHALL：
  - 不允许编辑订单
  - 不允许作废订单
  - 允许申请退款（变更为50-退费中）

- **WHEN** 订单状态为50（退费中）或99（已作废）
- **THEN** 系统SHALL：
  - 不允许任何修改操作

---

### 需求：订单金额计算规则

系统必须正确计算订单金额，保证数据一致性。

#### 场景：订单金额计算

- **WHEN** 创建或编辑订单
- **THEN** 系统SHALL：
  - 应收金额(amount_receivable) = SUM(所有商品的total_price)
  - 实收金额(amount_received) = SUM(所有商品的price) - 订单优惠金额(discount_amount)
  - 子订单应收金额 = 商品total_price
  - 子订单实收金额 = 商品price - 子订单优惠金额(child_discounts[goods_id])
  - 验证：SUM(子订单实收金额) = 订单实收金额

#### 场景：金额精度处理

- **WHEN** 进行金额计算
- **THEN** 系统SHALL：
  - 所有金额保留两位小数
  - 使用四舍五入规则
  - 金额比较时允许0.01的误差范围

---

### 需求：订单事务一致性

系统必须使用数据库事务保证订单操作的原子性。

#### 场景：创建订单事务

- **WHEN** 创建订单
- **THEN** 系统SHALL在同一事务中：
  - INSERT orders
  - INSERT childorders（批量）
  - INSERT orders_activity（批量）
  - 任一操作失败则全部回滚

#### 场景：编辑订单事务

- **WHEN** 编辑订单
- **THEN** 系统SHALL在同一事务中：
  - UPDATE orders
  - DELETE childorders（旧数据）
  - DELETE orders_activity（旧数据）
  - INSERT childorders（新数据）
  - INSERT orders_activity（新数据）
  - 任一操作失败则全部回滚

#### 场景：状态变更事务

- **WHEN** 提交或作废订单
- **THEN** 系统SHALL在同一事务中：
  - UPDATE orders状态
  - UPDATE childorders状态（批量）
  - 任一操作失败则全部回滚
