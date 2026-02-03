# 规范：收款管理

## 新增需求

### 需求：系统必须支持收款记录的创建、查询、确认和删除

系统必须提供收款管理功能，用于记录订单的收款信息，包括付款金额、付款方式、付款人、交易时间等，并支持收款的确认到账和删除操作。

#### 场景：用户创建收款记录

**前置条件**：
- 用户已登录系统
- 订单已存在且状态正常
- 学生已存在

**输入**：
- `order_id` (int, 必填): 订单ID
- `student_id` (int, 必填): 学生UID
- `payment_scenario` (int, 必填): 付款场景，0=线上，1=线下
- `payment_method` (int, 必填): 付款方式，0=微信，1=支付宝，2=优利支付，3=零零购支付，9=对公转账
- `payment_amount` (decimal, 必填): 付款金额，精度为2位小数
- `payer` (string, 可选): 付款人姓名，最大长度100字符
- `payee_entity` (int, 必填): 收款主体，0=北京，1=西安
- `merchant_order` (string, 条件必填): 商户订单号，线下场景时必填，最大长度100字符
- `trading_hours` (datetime, 可选): 交易时间，格式为"YYYY-MM-DD HH:mm:ss"

**处理逻辑**：
1. 验证必填字段不为空
2. 查询订单的实收金额（`actual_amount`）
3. 计算订单已收款总额：
   - 常规收款：`SUM(payment_amount)` WHERE `order_id = ? AND status IN (10, 20)`
   - 淘宝收款：`SUM(payment_amount)` WHERE `order_id = ? AND status = 30`
   - 总收款 = 常规收款 + 淘宝收款
4. 计算待支付金额：`unpaid_amount = actual_amount - total_paid`
5. 验证：`payment_amount <= unpaid_amount`，否则返回错误"付款金额不能超过待支付金额"
6. 插入收款记录到 `payment_collection` 表：
   - `status` 默认设置为 10（未核验）
   - `create_time` 自动设置为当前时间
   - 其他字段使用输入值
7. 调用订单状态更新逻辑（详见"订单支付状态自动更新"需求）
8. 返回成功响应

**输出**：
```json
{
  "code": 0,
  "message": "收款记录创建成功",
  "data": {
    "id": 123
  }
}
```

**异常情况**：
- 订单不存在：返回 `{"code": 1, "message": "订单不存在"}`
- 付款金额超过待支付金额：返回 `{"code": 1, "message": "付款金额不能超过待支付金额"}`
- 数据库错误：返回 `{"code": 1, "message": "系统错误"}`

---

#### 场景：用户查询收款列表

**前置条件**：
- 用户已登录系统

**输入**（所有参数可选，支持多条件组合查询）：
- `id` (int): 收款记录ID
- `student_id` (int): 学生UID
- `order_id` (int): 订单ID
- `payer` (string): 付款人姓名（模糊匹配）
- `payment_method` (int): 付款方式
- `trading_date` (date): 交易日期，格式为"YYYY-MM-DD"
- `status` (int): 状态，0=待支付，10=未核验，20=已支付
- `page` (int, 默认1): 页码
- `page_size` (int, 默认20): 每页记录数

**处理逻辑**：
1. 构建查询条件：
   - 如提供 `id`，精确匹配
   - 如提供 `student_id`，精确匹配
   - 如提供 `order_id`，精确匹配
   - 如提供 `payer`，使用 `LIKE '%payer%'` 模糊匹配
   - 如提供 `payment_method`，精确匹配
   - 如提供 `trading_date`，匹配交易时间的日期部分
   - 如提供 `status`，精确匹配
2. 执行分页查询，按 `create_time DESC` 排序
3. 关联查询学生姓名（LEFT JOIN `student` 表）
4. 返回结果列表

**输出**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "collections": [
      {
        "id": 123,
        "order_id": 456,
        "student_id": 789,
        "student_name": "张三",
        "payment_scenario": 0,
        "payment_method": 0,
        "payment_amount": 500.00,
        "payer": "李四",
        "payee_entity": 0,
        "trading_hours": "2024-01-15 10:30:00",
        "arrival_time": "2024-01-15 10:35:00",
        "merchant_order": "202401150001",
        "status": 20,
        "create_time": "2024-01-15 10:28:00"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

**异常情况**：
- 数据库错误：返回 `{"code": 1, "message": "系统错误"}`

---

#### 场景：用户确认收款到账

**前置条件**：
- 用户已登录系统
- 收款记录存在且状态为10（未核验）

**输入**：
- `id` (int, 必填): 收款记录ID（通过URL路径参数传递）

**处理逻辑**：
1. 查询收款记录，验证记录是否存在
2. 验证收款状态是否为10（未核验），否则返回错误"只能确认未核验的收款"
3. 更新收款记录：
   - `status` 设置为 20（已支付）
   - `arrival_time` 设置为当前时间
4. 调用订单状态更新逻辑（详见"订单支付状态自动更新"需求）
5. 调用分账生成逻辑（详见 `separate-account` 规范的"分账明细自动生成"需求）
6. 返回成功响应

**输出**：
```json
{
  "code": 0,
  "message": "确认到账成功",
  "data": null
}
```

**异常情况**：
- 收款记录不存在：返回 `{"code": 1, "message": "收款记录不存在"}`
- 收款状态不是未核验：返回 `{"code": 1, "message": "只能确认未核验的收款"}`
- 分账生成失败：回滚事务，返回 `{"code": 1, "message": "分账生成失败"}`
- 数据库错误：回滚事务，返回 `{"code": 1, "message": "系统错误"}`

**注意事项**：
- 确认到账、更新订单状态、生成分账明细必须在同一事务中执行，确保原子性
- 如分账生成失败，必须回滚收款状态的更新

---

#### 场景：用户删除收款记录

**前置条件**：
- 用户已登录系统
- 收款记录存在且状态为10（未核验）

**输入**：
- `id` (int, 必填): 收款记录ID（通过URL路径参数传递）

**处理逻辑**：
1. 查询收款记录，验证记录是否存在
2. 验证收款状态是否为10（未核验），否则返回错误"只能删除未核验的收款"
3. 删除收款记录（物理删除）
4. 调用订单状态更新逻辑（详见"订单支付状态自动更新"需求）
5. 返回成功响应

**输出**：
```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

**异常情况**：
- 收款记录不存在：返回 `{"code": 1, "message": "收款记录不存在"}`
- 收款状态不是未核验：返回 `{"code": 1, "message": "只能删除未核验的收款"}`
- 数据库错误：返回 `{"code": 1, "message": "系统错误"}`

**注意事项**：
- 只有未核验（status=10）的收款才能删除
- 已确认到账（status=20）的收款不能删除，防止已生成的分账明细数据不一致

---

### 需求：订单支付状态必须根据收款情况自动更新

当收款记录发生变化（新增、确认、删除）时，系统必须自动计算订单的总收款金额，并根据收款情况更新订单的支付状态。

#### 场景：订单支付状态自动更新

**触发时机**：
- 新增收款记录后
- 确认收款到账后
- 删除收款记录后

**输入**：
- `order_id` (int): 订单ID

**处理逻辑**：
1. 查询订单的实收金额（`actual_amount`）
2. 计算常规收款总额：
   ```sql
   SELECT COALESCE(SUM(payment_amount), 0)
   FROM payment_collection
   WHERE order_id = ? AND status IN (10, 20)
   ```
3. 计算淘宝收款总额（本次不实现淘宝收款，暂不统计）：
   ```sql
   SELECT COALESCE(SUM(payment_amount), 0)
   FROM taobao_payment
   WHERE order_id = ? AND status = 30
   ```
4. 计算总收款金额：`total_paid = regular_paid + taobao_paid`
5. 根据收款情况确定订单状态：
   - 如果 `total_paid == 0`：订单状态设置为 20（未支付）
   - 如果 `0 < total_paid < actual_amount`：订单状态设置为 30（部分支付）
   - 如果 `total_paid >= actual_amount`：订单状态设置为 40（已支付）
6. 更新 `orders` 表的 `status` 字段

**输出**：
- 无返回值（内部调用）

**注意事项**：
- 此逻辑由收款管理模块内部调用，不对外暴露接口
- 必须在收款记录变更的同一事务中执行
- 金额比较时使用精确比较（避免浮点数精度问题）

---

## 数据模型

### 实体：PaymentCollection（收款记录）

**表名**：`payment_collection`

**字段**：

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY, AUTO_INCREMENT | 主键 |
| order_id | INT | NOT NULL | 关联订单ID |
| student_id | INT | NOT NULL | 学生UID |
| payment_scenario | INT | NOT NULL | 付款场景：0=线上，1=线下 |
| payment_method | INT | NOT NULL | 付款方式：0=微信，1=支付宝，2=优利，3=零零购，9=对公 |
| payment_amount | DECIMAL(10,2) | NOT NULL | 付款金额 |
| payer | VARCHAR(100) | NULL | 付款人 |
| payee_entity | INT | NOT NULL | 收款主体：0=北京，1=西安 |
| trading_hours | DATETIME | NULL | 交易时间 |
| arrival_time | DATETIME | NULL | 到账时间 |
| merchant_order | VARCHAR(100) | NULL | 商户订单号 |
| status | INT | NOT NULL, DEFAULT 10 | 状态：0=待支付，10=未核验，20=已支付 |
| create_time | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**索引**：
- PRIMARY KEY: `id`
- INDEX: `order_id`
- INDEX: `student_id`
- INDEX: `status`

**外键**：
- 无（为保持与源系统一致，不使用外键约束）

---

## API接口

### GET /api/payment-collections

获取收款列表。

**请求参数**（Query String）：
- `id` (int, 可选): 收款ID
- `student_id` (int, 可选): 学生UID
- `order_id` (int, 可选): 订单ID
- `payer` (string, 可选): 付款人
- `payment_method` (int, 可选): 付款方式
- `trading_date` (string, 可选): 交易日期，格式"YYYY-MM-DD"
- `status` (int, 可选): 状态
- `page` (int, 可选, 默认1): 页码
- `page_size` (int, 可选, 默认20): 每页记录数

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "collections": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

### POST /api/payment-collections

新增收款记录。

**请求体**（JSON）：
```json
{
  "order_id": 456,
  "student_id": 789,
  "payment_scenario": 0,
  "payment_method": 0,
  "payment_amount": 500.00,
  "payer": "李四",
  "payee_entity": 0,
  "merchant_order": "202401150001",
  "trading_hours": "2024-01-15 10:30:00"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "收款记录创建成功",
  "data": {
    "id": 123
  }
}
```

---

### PUT /api/payment-collections/:id/confirm

确认收款到账。

**路径参数**：
- `id` (int): 收款记录ID

**请求体**：无

**响应**：
```json
{
  "code": 0,
  "message": "确认到账成功",
  "data": null
}
```

---

### DELETE /api/payment-collections/:id

删除收款记录。

**路径参数**：
- `id` (int): 收款记录ID

**请求体**：无

**响应**：
```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

---

## 业务规则

1. **金额验证**：新增收款时，付款金额不能超过订单的待支付金额
2. **状态约束**：
   - 新增收款时，状态默认为10（未核验）
   - 只有未核验（status=10）的收款才能确认到账
   - 只有未核验（status=10）的收款才能删除
3. **原子性**：确认到账操作必须与分账生成在同一事务中执行
4. **幂等性**：重复确认到账应返回错误，不能重复生成分账
5. **级联更新**：收款记录变化时，自动更新订单支付状态

---

## 依赖关系

**依赖模块**：
- 订单管理模块（order）：查询订单实收金额、更新订单状态
- 学生管理模块（student）：查询学生信息
- 分账管理模块（separate-account）：生成分账明细

**被依赖模块**：
- 无

---

## 非功能需求

1. **性能**：收款列表查询响应时间 < 500ms（10000条数据）
2. **并发**：支持多个收款同时确认，使用数据库锁防止重复生成分账
3. **精度**：金额计算使用 `DECIMAL(10,2)` 类型，避免浮点数精度问题
4. **事务**：确认到账操作必须在事务中执行，失败时自动回滚
5. **日志**：关键操作（确认到账、删除）记录详细日志，包括操作人、操作时间、影响记录

---

## 测试用例

### 测试用例1：新增收款成功

**前置条件**：订单实收1000元，已收款0元

**操作**：新增收款500元

**预期结果**：
- 收款记录创建成功，status=10
- 订单状态更新为30（部分支付）

---

### 测试用例2：新增收款失败（金额超限）

**前置条件**：订单实收1000元，已收款800元

**操作**：新增收款300元

**预期结果**：
- 返回错误"付款金额不能超过待支付金额"
- 收款记录未创建

---

### 测试用例3：确认到账成功

**前置条件**：收款记录存在，status=10

**操作**：确认收款到账

**预期结果**：
- 收款状态更新为20
- arrival_time设置为当前时间
- 自动生成分账明细
- 订单状态更新

---

### 测试用例4：确认到账失败（状态不符）

**前置条件**：收款记录存在，status=20

**操作**：确认收款到账

**预期结果**：
- 返回错误"只能确认未核验的收款"
- 收款状态不变

---

### 测试用例5：删除收款成功

**前置条件**：收款记录存在，status=10

**操作**：删除收款

**预期结果**：
- 收款记录被删除
- 订单状态更新

---

### 测试用例6：删除收款失败（状态不符）

**前置条件**：收款记录存在，status=20

**操作**：删除收款

**预期结果**：
- 返回错误"只能删除未核验的收款"
- 收款记录不变

---

## 参考资料

- 源系统实现：`D:\claude space\ZhixinStudentSaaS\backend\app.py` (行号: 4304-4640)
- 前端页面：`D:\claude space\ZhixinStudentSaaS\frontend\index.html` (行号: 1642-2113)
- 项目约定：`openspec/project.md`
