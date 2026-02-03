# 规范：分账明细管理

## 新增需求

### 需求：系统必须支持分账明细的自动生成和查询

系统必须在收款确认到账后自动生成分账明细，将收款金额按子订单顺序分配，并支持分账明细的查询功能。

#### 场景：分账明细自动生成

**前置条件**：
- 收款记录已确认到账（status=20）
- 订单存在且包含子订单
- 该收款尚未生成分账明细

**触发时机**：
- 收款确认到账后自动触发（由收款管理模块调用）

**输入**：
- `payment_id` (int): 收款记录ID
- `order_id` (int): 订单ID

**处理逻辑**：

1. **查询收款信息**
   ```sql
   SELECT id, student_id, payment_amount
   FROM payment_collection
   WHERE id = ?
   ```

2. **查询子订单列表**（按ID升序，确保分配顺序固定）
   ```sql
   SELECT id, goods_id, goods_name, actual_amount
   FROM childorders
   WHERE orders_id = ?
   ORDER BY id ASC
   ```

3. **防重复检查**
   ```sql
   SELECT COUNT(*) as count
   FROM separate_account
   WHERE payment_id = ? AND orders_id = ? AND payment_type = 0
   ```
   如果 count > 0，说明已生成分账，直接返回（幂等性保证）

4. **遍历子订单，按序分配收款金额**
   ```
   初始化 remaining_amount = 收款金额

   FOR EACH child_order IN child_orders:
       // 4.1 计算子订单已分配金额
       allocated_amount = SELECT SUM(separate_amount)
                          FROM separate_account
                          WHERE childorders_id = child_order.id

       // 4.2 计算子订单还需金额
       needed_amount = child_order.actual_amount - allocated_amount

       IF needed_amount <= 0:
           CONTINUE  // 子订单已满额，跳过

       // 4.3 计算本次分账金额（取剩余收款和还需金额的较小值）
       separate_amount = MIN(remaining_amount, needed_amount)

       // 4.4 插入分账明细
       INSERT INTO separate_account (
           uid, orders_id, childorders_id, payment_id, payment_type,
           goods_id, goods_name, separate_amount, type
       ) VALUES (
           student_id, order_id, child_order.id, payment_id, 0,
           child_order.goods_id, child_order.goods_name, separate_amount, 0
       )

       // 4.5 更新剩余收款金额
       remaining_amount -= separate_amount

       // 4.6 更新子订单状态（详见"子订单支付状态自动更新"需求）
       CALL UpdateChildOrderStatus(child_order.id)

       IF remaining_amount <= 0:
           BREAK  // 收款已分配完，结束循环
   ```

5. **返回成功**

**输出**：
- 无返回值（内部调用）

**异常情况**：
- 收款记录不存在：抛出异常
- 订单不存在：抛出异常
- 子订单列表为空：记录警告日志，不生成分账
- 数据库错误：抛出异常，调用方回滚事务

**注意事项**：
- **分配顺序**：严格按子订单ID升序分配，确保"先进先出"原则
- **幂等性**：通过防重复检查，确保同一收款不会重复生成分账
- **原子性**：必须在收款确认的同一事务中执行
- **金额精度**：使用 `DECIMAL(10,2)` 类型，避免浮点数精度问题

---

#### 场景：子订单支付状态自动更新

**触发时机**：
- 分账明细生成后，针对每个子订单自动触发

**输入**：
- `childorder_id` (int): 子订单ID

**处理逻辑**：
1. 查询子订单的实收金额（`actual_amount`）
2. 计算子订单的总分账金额：
   ```sql
   SELECT COALESCE(SUM(separate_amount), 0)
   FROM separate_account
   WHERE childorders_id = ? AND type = 0
   ```
   （只统计售卖类型，type=0，不包括冲回和退费）

3. 根据分账情况确定子订单状态：
   - 如果 `total_separate == 0`：子订单状态设置为 10（未支付）
   - 如果 `0 < total_separate < actual_amount`：子订单状态设置为 20（部分支付）
   - 如果 `total_separate >= actual_amount`：子订单状态设置为 30（已支付）

4. 更新 `childorders` 表的 `status` 字段

**输出**：
- 无返回值（内部调用）

**注意事项**：
- 此逻辑由分账管理模块内部调用，不对外暴露接口
- 必须在生成分账的同一事务中执行
- 金额比较时使用精确比较（避免浮点数精度问题）

---

#### 场景：用户查询分账明细列表

**前置条件**：
- 用户已登录系统

**输入**（所有参数可选，支持多条件组合查询）：
- `id` (int): 分账明细ID
- `uid` (int): 学生UID
- `orders_id` (int): 订单ID
- `childorders_id` (int): 子订单ID
- `goods_id` (int): 商品ID
- `payment_id` (int): 收款ID
- `payment_type` (int): 收款类型，0=常规收款，1=淘宝收款
- `type` (int): 分账类型，0=售卖，1=冲回，2=退费
- `page` (int, 默认1): 页码
- `page_size` (int, 默认20): 每页记录数

**处理逻辑**：
1. 构建查询条件：
   - 如提供 `id`，精确匹配
   - 如提供 `uid`，精确匹配
   - 如提供 `orders_id`，精确匹配
   - 如提供 `childorders_id`，精确匹配
   - 如提供 `goods_id`，精确匹配
   - 如提供 `payment_id`，精确匹配
   - 如提供 `payment_type`，精确匹配
   - 如提供 `type`，精确匹配
2. 执行分页查询，按 `create_time DESC` 排序
3. 返回结果列表

**输出**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "separate_accounts": [
      {
        "id": 1,
        "uid": 789,
        "orders_id": 456,
        "childorders_id": 123,
        "payment_id": 100,
        "payment_type": 0,
        "goods_id": 50,
        "goods_name": "数学课程",
        "separate_amount": 300.00,
        "type": 0,
        "create_time": "2024-01-15 10:35:00"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

**异常情况**：
- 数据库错误：返回 `{"code": 1, "message": "系统错误"}`

---

## 数据模型

### 实体：SeparateAccount（分账明细）

**表名**：`separate_account`

**字段**：

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY, AUTO_INCREMENT | 主键 |
| uid | INT | NOT NULL | 学生UID |
| orders_id | INT | NOT NULL | 订单ID |
| childorders_id | INT | NOT NULL | 子订单ID |
| payment_id | INT | NOT NULL | 收款ID |
| payment_type | INT | NOT NULL | 收款类型：0=常规收款，1=淘宝收款 |
| goods_id | INT | NOT NULL | 商品ID |
| goods_name | VARCHAR(100) | NOT NULL | 商品名称 |
| separate_amount | DECIMAL(10,2) | NOT NULL | 分账金额 |
| type | INT | NOT NULL, DEFAULT 0 | 类型：0=售卖，1=冲回，2=退费 |
| create_time | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**索引**：
- PRIMARY KEY: `id`
- INDEX: `orders_id`
- INDEX: `childorders_id`
- INDEX: `payment_id`
- INDEX: `uid`

**外键**：
- 无（为保持与源系统一致，不使用外键约束）

**联合唯一索引**（用于防重复）：
- UNIQUE KEY: `(payment_id, orders_id, childorders_id)` （可选，增强幂等性）

---

## API接口

### GET /api/separate-accounts

获取分账明细列表。

**请求参数**（Query String）：
- `id` (int, 可选): 分账明细ID
- `uid` (int, 可选): 学生UID
- `orders_id` (int, 可选): 订单ID
- `childorders_id` (int, 可选): 子订单ID
- `goods_id` (int, 可选): 商品ID
- `payment_id` (int, 可选): 收款ID
- `payment_type` (int, 可选): 收款类型
- `type` (int, 可选): 分账类型
- `page` (int, 可选, 默认1): 页码
- `page_size` (int, 可选, 默认20): 每页记录数

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "separate_accounts": [...],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 业务规则

1. **分配顺序**：收款金额按子订单ID升序分配，遵循"先进先出"原则
2. **分配策略**：每个子订单分配金额 = MIN(剩余收款金额, 子订单还需金额)
3. **幂等性**：同一收款不能重复生成分账，通过防重复检查保证
4. **原子性**：分账生成必须与收款确认在同一事务中执行
5. **类型固定**：自动生成的分账明细，`type` 固定为0（售卖）
6. **收款类型**：常规收款生成的分账，`payment_type` 为0；淘宝收款为1（本次不实现）
7. **级联更新**：分账生成后，自动更新子订单支付状态

---

## 算法示例

### 示例1：单笔收款完全分配

**订单信息**：
- 子订单1：实收500元，已分账0元
- 子订单2：实收300元，已分账0元
- 子订单3：实收200元，已分账0元

**收款信息**：
- 收款金额：1000元

**分配结果**：
- 子订单1：分账500元（状态更新为30-已支付）
- 子订单2：分账300元（状态更新为30-已支付）
- 子订单3：分账200元（状态更新为30-已支付）
- 剩余收款：0元

---

### 示例2：单笔收款部分分配

**订单信息**：
- 子订单1：实收500元，已分账0元
- 子订单2：实收300元，已分账0元
- 子订单3：实收200元，已分账0元

**收款信息**：
- 收款金额：600元

**分配结果**：
- 子订单1：分账500元（状态更新为30-已支付）
- 子订单2：分账100元（状态更新为20-部分支付）
- 子订单3：分账0元（状态保持10-未支付）
- 剩余收款：0元

---

### 示例3：多笔收款累计分配

**订单信息**（初始状态）：
- 子订单1：实收500元，已分账0元
- 子订单2：实收300元，已分账0元

**第一笔收款**：300元
- 子订单1：分账300元（状态20-部分支付）
- 子订单2：分账0元（状态10-未支付）

**第二笔收款**：400元
- 子订单1：分账200元（累计500元，状态30-已支付）
- 子订单2：分账200元（状态20-部分支付）

**第三笔收款**：100元
- 子订单2：分账100元（累计300元，状态30-已支付）

---

## 依赖关系

**依赖模块**：
- 收款管理模块（payment-collection）：获取收款信息
- 订单管理模块（order）：查询子订单列表、更新子订单状态
- 商品管理模块（goods）：分账记录中包含商品信息

**被依赖模块**：
- 收款管理模块（payment-collection）：收款确认到账时调用分账生成

---

## 非功能需求

1. **性能**：分账生成时间 < 200ms（10个子订单）
2. **并发**：同一收款的分账生成使用数据库锁，防止并发重复生成
3. **精度**：金额计算使用 `DECIMAL(10,2)` 类型，避免浮点数精度问题
4. **事务**：分账生成与收款确认在同一事务中，失败时自动回滚
5. **日志**：记录分账生成的详细日志，包括分配策略和结果

---

## 测试用例

### 测试用例1：单笔收款完全覆盖

**前置条件**：
- 订单实收1000元
- 子订单1实收500元，已分账0元
- 子订单2实收300元，已分账0元
- 子订单3实收200元，已分账0元

**操作**：确认收款1000元

**预期结果**：
- 生成3条分账明细：500元、300元、200元
- 3个子订单状态均为30（已支付）

---

### 测试用例2：单笔收款部分覆盖

**前置条件**：
- 订单实收1000元
- 子订单1实收500元，已分账0元
- 子订单2实收300元，已分账0元
- 子订单3实收200元，已分账0元

**操作**：确认收款600元

**预期结果**：
- 生成2条分账明细：500元、100元
- 子订单1状态30（已支付）
- 子订单2状态20（部分支付）
- 子订单3状态10（未支付）

---

### 测试用例3：多笔收款累计覆盖

**前置条件**：
- 订单实收800元
- 子订单1实收500元，已分账0元
- 子订单2实收300元，已分账0元

**操作**：
1. 确认收款300元
2. 确认收款400元
3. 确认收款100元

**预期结果**：
- 第一笔：生成1条分账300元（子订单1部分支付）
- 第二笔：生成2条分账200元+200元（子订单1已支付，子订单2部分支付）
- 第三笔：生成1条分账100元（子订单2已支付）
- 总计5条分账明细

---

### 测试用例4：防重复生成

**前置条件**：
- 收款已确认，分账已生成

**操作**：再次调用分账生成逻辑（模拟并发）

**预期结果**：
- 防重复检查生效，不生成新的分账明细
- 返回成功（幂等性）

---

### 测试用例5：子订单已满额

**前置条件**：
- 订单实收800元
- 子订单1实收500元，已分账500元（已支付）
- 子订单2实收300元，已分账0元

**操作**：确认收款400元

**预期结果**：
- 跳过子订单1（已满额）
- 生成1条分账300元给子订单2
- 子订单2状态30（已支付）
- 剩余100元未分配（警告：订单已满额但仍有剩余收款）

---

## 边界情况处理

1. **收款金额大于订单总额**：允许创建收款，但分账时只分配到订单满额为止，剩余金额不分配（记录警告日志）
2. **子订单列表为空**：记录警告日志，不生成分账，返回成功
3. **子订单已全部满额**：记录警告日志，不生成分账，返回成功
4. **并发生成分账**：使用数据库行锁或防重复检查，确保幂等性
5. **金额精度问题**：使用 `DECIMAL` 类型和精确比较，避免浮点数误差

---

## 参考资料

- 源系统实现：`D:\claude space\ZhixinStudentSaaS\backend\app.py` (行号: 3588-3812, 5249-5325)
- 前端页面：`D:\claude space\ZhixinStudentSaaS\frontend\index.html` (行号: 2114-2222)
- 项目约定：`openspec/project.md`
- 关联规范：`payment-collection` 收款管理规范
