# separate-account 规范增量

## 修改需求

### 需求：系统必须在退费审批通过后正确生成退费类分账明细

系统必须修改退费审批通过后的分账明细生成逻辑，确保退费类分账明细的生成与原Python版本保持一致。无论是否执行冲回流程，都必须生成退费类分账明细。

#### 场景：退费审批通过后生成退费类分账明细

**前置条件**：
- 退费订单已提交并进入审批流程
- 审批流程已全部通过
- 订单存在收款记录（已到账）
- 退费订单包含退费子订单

**触发时机**：
- 审批流程最后一个节点通过后自动触发

**处理逻辑**：

1. **更新退费订单状态**
   - 更新退费订单状态为已通过(10)
   - 更新退费子订单状态为已通过(10)
   - 更新退费补充信息状态为已通过(10)

2. **获取收款信息**
   ```sql
   -- 常规收款
   SELECT id as payment_id, payment_amount, 0 as payment_type
   FROM payment_collection
   WHERE order_id = ? AND status IN (10, 20)
   ORDER BY id ASC

   -- 淘宝收款
   SELECT id as payment_id, payment_amount, 1 as payment_type
   FROM taobao_payment
   WHERE order_id = ? AND status = 30
   ORDER BY id ASC
   ```

3. **检查是否需要冲回**
   - 获取退费子订单ID列表
   - 获取退费收款分配列表（refund_payment表）
   - 对每个退费收款，计算该收款在被退费子订单上的分账总额：
     ```sql
     SELECT COALESCE(SUM(separate_amount), 0) as total_separate
     FROM separate_account
     WHERE payment_id = ? AND payment_type = ?
         AND childorders_id IN (退费子订单ID列表)
         AND type = 0  -- 售卖类
         AND id NOT IN (
             SELECT parent_id FROM separate_account
             WHERE payment_id = ? AND payment_type = ? AND parent_id IS NOT NULL
         )  -- 排除已冲回的
     ```
   - 如果任何收款的分账总额 < 退费金额，则需要冲回

4. **如果需要冲回，执行冲回和重新分账**

   a. **冲回所有未冲回的售卖类分账明细**
   ```sql
   -- 查询所有未冲回的售卖类分账
   SELECT id, uid, orders_id, childorders_id, payment_id, payment_type,
          goods_id, goods_name, separate_amount
   FROM separate_account
   WHERE orders_id = ? AND type = 0
       AND id NOT IN (
           SELECT parent_id FROM separate_account
           WHERE orders_id = ? AND parent_id IS NOT NULL
       )

   -- 为每条原分账生成冲回记录
   INSERT INTO separate_account
   (uid, orders_id, childorders_id, payment_id, payment_type,
    goods_id, goods_name, separate_amount, type, parent_id)
   VALUES (原分账的值..., -原分账金额, 1, 原分账ID)
   ```

   b. **第一批售卖分账：用退费收款分配给退费子订单**
   - 获取退费收款分配列表（refund_payment表）
   - 获取退费子订单列表（refund_order_item表，关联商品信息）
   - 为退费子订单记录剩余需求
   - 按退费收款顺序，为每个退费子订单分配金额：
     - 本次分配金额 = min(退费收款剩余, 退费子订单剩余需求)
     - 插入售卖类分账明细（type=0）
     - 更新剩余金额

   c. **第二批售卖分账：用剩余收款分配给剩余子订单需求**
   - 计算每个收款的剩余金额 = 总收款金额 - 退费金额
   - 获取所有子订单
   - 计算每个子订单的剩余需求 = 实收金额 - 退费金额
   - 按收款顺序，为每个子订单分配剩余金额：
     - 本次分配金额 = min(收款剩余, 子订单剩余需求)
     - 插入售卖类分账明细（type=0）
     - 更新剩余金额

5. **生成退费类分账明细（无论是否冲回都执行）**

   **重要**：此步骤必须在冲回流程之外独立执行，确保无论是否执行冲回，都会生成退费类分账明细。

   ```sql
   -- 查询所有退费子订单和退费金额
   SELECT roi.childorder_id, roi.refund_amount, co.goodsid, g.name as goods_name
   FROM refund_order_item roi
   INNER JOIN childorders co ON roi.childorder_id = co.id
   LEFT JOIN goods g ON co.goodsid = g.id
   WHERE roi.refund_order_id = ?
   ORDER BY roi.childorder_id ASC
   ```

   为每个退费子订单生成退费类分账明细：

   a. **查询该子订单当前的售卖分账分布**（按收款ID升序）：
   ```sql
   SELECT payment_id, payment_type, separate_amount
   FROM separate_account
   WHERE childorders_id = ? AND type = 0
       AND id NOT IN (
           SELECT parent_id FROM separate_account
           WHERE childorders_id = ? AND parent_id IS NOT NULL
       )
   ORDER BY payment_id ASC
   ```

   b. **按照售卖分账的分布生成退费类分账**：
   - 遍历该子订单的售卖分账记录
   - 从每个售卖分账中扣除对应金额，直到退费金额分配完毕
   - 每个售卖分账对应生成一条退费类分账明细：
     ```sql
     INSERT INTO separate_account
     (uid, orders_id, childorders_id, payment_id, payment_type,
      goods_id, goods_name, separate_amount, type)
     VALUES (学生ID, 订单ID, 子订单ID, 收款ID, 收款类型,
             商品ID, 商品名称, -分配的退费金额, 2)
     ```
   - 金额为负数，类型为2（退费类）

   **示例**：
   - 子订单3的售卖分账：收款1有200元，收款2有100元
   - 子订单3退费300元
   - 生成退费类分账：
     - 收款1扣除200元（-200元，type=2）
     - 收款2扣除100元（-100元，type=2）

   **注意**：
   - 如果执行了冲回，使用重新生成的售卖分账
   - 如果未执行冲回，使用原有的售卖分账
   - 退费类分账的分布必须与售卖分账的分布一致

6. **更新子订单状态**
   - 查询订单的所有子订单
   - 对每个子订单计算净分账金额：
     ```sql
     SELECT COALESCE(SUM(separate_amount), 0) as net_allocated
     FROM separate_account
     WHERE childorders_id = ?
         AND (
             (type = 0 AND id NOT IN (
                 SELECT parent_id FROM separate_account
                 WHERE childorders_id = ? AND parent_id IS NOT NULL
             ))  -- 未冲回的售卖类
             OR type = 2  -- 退费类
         )
     ```
   - 根据净分账金额确定子订单状态：
     - 净分账 <= 0：未支付(10)
     - 0 < 净分账 < 实收金额：部分支付(20)
     - 净分账 >= 实收金额：已支付(30)
   - 更新子订单状态

7. **更新订单状态**
   - 计算总收款金额（常规+淘宝）
   - 计算总退费金额（常规+淘宝，已通过的退费订单）
   - 净收款 = 总收款 - 总退费
   - 根据净收款确定订单状态：
     - 净收款 <= 0：未支付(20)
     - 净收款 >= 应收金额：已支付(40)
     - 其他：部分支付(30)
   - 更新订单状态

**输出**：
- 无返回值（内部调用）

**异常情况**：
- 退费订单不存在：抛出异常
- 订单不存在：抛出异常
- 无收款信息：跳过分账逻辑，仅更新状态
- 数据库错误：抛出异常，回滚事务

**注意事项**：
- **事务性**：整个流程必须在同一个事务中执行，确保原子性
- **退费类分账独立性**：退费类分账明细的生成必须独立于冲回流程，无论是否执行冲回都要生成
- **金额精度**：使用 `DECIMAL(10,2)` 类型，避免浮点数精度问题
- **分配顺序**：
  - 收款按ID升序
  - 子订单按ID升序
- **简化处理**：退费金额全部分配到第一个收款，与原Python版本保持一致

**与原需求的差异**：
- 原需求未详细描述退费类分账明细的生成逻辑
- 本次修改补充了完整的退费类分账明细生成流程
- 强调退费类分账必须在冲回流程外独立执行

**实现要点**：
1. 避免变量名冲突：退费类分账使用独立的变量名
2. 数据查询：重新查询退费子订单和收款信息，确保数据最新
3. 逻辑清晰：退费类分账明细生成在代码结构上独立于冲回流程
4. 保持一致：与原Python版本的逻辑保持完全一致
