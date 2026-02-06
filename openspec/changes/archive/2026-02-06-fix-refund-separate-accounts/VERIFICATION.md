# 修复验证指南

## 编译验证

✅ **已通过编译验证**

```bash
# 编译核心代码
go build ./internal/...

# 编译主程序
go build -o charonoms.exe ./cmd/server/main.go
```

结果：编译成功，无错误

## 代码审查清单

✅ **变量命名**
- 使用 `refundItemsForRefundSeparate` 避免与冲回流程中的变量冲突
- 使用 `childSeparates` 存储子订单的售卖分账分布

✅ **SQL查询**
- 查询退费子订单：JOIN childorders和goods表获取完整信息
- 查询售卖分账分布：排除已冲回的分账（parent_id不为NULL）
- 按payment_id升序排序，确保分账顺序一致

✅ **逻辑正确性**
- 为每个退费子订单独立查询其售卖分账分布
- 按照售卖分账的payment_id和payment_type生成退费类分账
- 使用min(剩余退费金额, 售卖分账金额)计算每次分配金额
- 退费类分账金额为负数，type=2

✅ **事务完整性**
- 整个流程在同一个事务中执行
- 查询和插入操作使用同一个tx对象

✅ **注释清晰**
- 说明了为什么重新查询（确保数据最新且避免变量冲突）
- 说明了按照售卖分账分布生成退费类分账的逻辑
- 说明了计算逻辑（min函数）

## 功能测试指南

### 准备测试数据

```sql
-- 1. 创建订单和子订单
INSERT INTO orders (id, student_id, amount_received, status)
VALUES (1, 1, 600, 30);

INSERT INTO childorders (id, parentsid, goodsid, amount_received, status)
VALUES
    (1, 1, 1, 100, 30),
    (2, 1, 2, 200, 30),
    (3, 1, 3, 300, 30);

-- 2. 创建收款
INSERT INTO payment_collection (id, order_id, student_id, payment_amount, status)
VALUES
    (1, 1, 1, 400, 20),
    (2, 1, 1, 200, 20);

-- 3. 生成初始分账
INSERT INTO separate_account (uid, orders_id, childorders_id, payment_id, payment_type, goods_id, goods_name, separate_amount, type)
VALUES
    (1, 1, 1, 1, 0, 1, '商品1', 100, 0),
    (1, 1, 2, 1, 0, 2, '商品2', 200, 0),
    (1, 1, 3, 1, 0, 3, '商品3', 100, 0),
    (1, 1, 3, 2, 0, 3, '商品3', 200, 0);

-- 4. 创建退费订单
INSERT INTO refund_order (id, order_id, student_id, refund_amount, status, submit_time)
VALUES (1, 1, 1, 500, 0, NOW());

INSERT INTO refund_order_item (id, refund_order_id, childorder_id, refund_amount)
VALUES
    (1, 1, 2, 200),
    (2, 1, 3, 300);

INSERT INTO refund_payment (id, refund_order_id, payment_id, payment_type, refund_amount)
VALUES
    (1, 1, 1, 0, 400),
    (2, 1, 2, 0, 100);

-- 5. 创建审批流程
-- (根据实际审批流程创建)
```

### 执行测试

1. **触发退费审批通过**
   ```bash
   # 通过API或直接调用审批通过接口
   curl -X POST http://localhost:8080/api/approval/approve/{node_case_user_id}
   ```

2. **验证分账明细**
   ```sql
   -- 查看所有分账明细
   SELECT id, childorders_id, payment_id, separate_amount, type, parent_id
   FROM separate_account
   WHERE orders_id = 1
   ORDER BY id;
   ```

   **预期结果**:
   - 冲回记录（type=1, 负金额, parent_id不为NULL）: 4条
   - 新的售卖分账（type=0, 正金额）: 3条
     - 子订单2 ← 收款1: 200元
     - 子订单3 ← 收款1: 200元
     - 子订单3 ← 收款2: 100元
     - 子订单1 ← 收款2: 100元
   - 退费类分账（type=2, 负金额）: 3条
     - 子订单2 ← 收款1: -200元
     - 子订单3 ← 收款1: -200元
     - 子订单3 ← 收款2: -100元

3. **验证子订单状态**
   ```sql
   SELECT id, amount_received, status
   FROM childorders
   WHERE parentsid = 1;
   ```

   **预期结果**:
   - 子订单1: status=30 (已支付, 净分账100元)
   - 子订单2: status=10 (未支付, 净分账0元)
   - 子订单3: status=10 (未支付, 净分账0元)

4. **验证订单状态**
   ```sql
   SELECT id, amount_received, status
   FROM orders
   WHERE id = 1;
   ```

   **预期结果**:
   - 订单1: status=30 (部分支付, 净收款100元)

### 验证退费类分账按照售卖分账分布生成

运行测试脚本：
```bash
go run scripts/test_refund_separate_fix.go
```

**验证点**:
1. 子订单2的退费类分账只有收款1（因为售卖分账只有收款1）
2. 子订单3的退费类分账有收款1和收款2（因为售卖分账有两个收款）
3. 退费类分账的金额与售卖分账的金额对应
4. 退费类分账总额 = 退费金额

## 回归测试清单

### 不需要冲回的场景

**条件**: 所有收款在退费子订单上的分账总额 >= 退费金额

**验证点**:
- ✓ 不执行冲回流程（无冲回记录）
- ✓ 保留原有售卖分账
- ✓ 生成退费类分账（按照原有售卖分账分布）
- ✓ 子订单状态正确更新

### 需要冲回的场景

**条件**: 存在任何收款在退费子订单上的分账总额 < 退费金额

**验证点**:
- ✓ 执行冲回流程（生成冲回记录）
- ✓ 重新生成售卖分账（第一批+第二批）
- ✓ 生成退费类分账（按照新的售卖分账分布）
- ✓ 子订单状态正确更新

### 边界情况

1. **退费金额 = 子订单金额**
   - 验证退费类分账总额正确
   - 验证子订单状态为未支付

2. **部分退费**
   - 验证退费类分账只影响退费子订单
   - 验证未退费子订单的售卖分账不受影响

3. **多个收款**
   - 验证退费类分账按照售卖分账分布到多个收款
   - 验证每个收款的退费金额正确

4. **淘宝收款和常规收款混合**
   - 验证payment_type字段正确
   - 验证两种类型收款的退费类分账都正确生成

## 性能验证

**关注点**:
- 每个退费子订单增加一次数据库查询（查询售卖分账分布）
- 对于有N个退费子订单的退费订单，增加N次查询

**建议**:
- 监控退费审批通过的执行时间
- 如果性能有问题，可以考虑批量查询所有退费子订单的售卖分账

**当前性能评估**:
- 通常退费子订单数量不多（< 10个）
- 每次查询很简单（索引命中）
- 在可接受范围内

## 验证结果

- [ ] 编译验证通过
- [ ] 代码审查通过
- [ ] 功能测试通过（不需要冲回场景）
- [ ] 功能测试通过（需要冲回场景）
- [ ] 回归测试通过
- [ ] 性能验证通过

## 问题记录

（如有问题请记录在此）

---

**验证人**: ___________
**验证日期**: ___________
**验证结果**: □ 通过  □ 不通过
