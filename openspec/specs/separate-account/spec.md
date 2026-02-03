# separate-account 规范

## 目的
待定 - 由归档变更 migrate-financial-management 创建。归档后请更新目的。
## 需求
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

