# 设计文档

## 问题分析

### 原Python版本逻辑(正确的实现)

```python
# 位置: app.py:3816-4300
def handle_refund_chargeback(cursor, refund_order_id, order_id):
    # 1. 获取退费相关信息和收款信息
    # 2. 检查是否需要冲回
    # 3. 如果需要冲回:
    #    a. 冲回所有未冲回的售卖类分账明细
    #    b. 第一批售卖分账:用退费收款分配给退费子订单
    #    c. 第二批售卖分账:用剩余收款分配给剩余子订单
    # 4. 生成退费类分账明细(无论是否冲回都执行)
    #    - 为每个退费子订单按收款顺序分配退费金额
    #    - 简化处理:退费金额全部分配到第一个收款
    #    - 插入负金额的退费类分账明细(type=2)
    # 5. 更新子订单状态
    # 6. 更新订单状态
```

关键点:
- **退费类分账明细生成在冲回流程之外** (第4130-4189行)
- 无论是否执行冲回,都要生成退费类分账明细
- 为每个退费子订单生成一条退费类分账(type=2,负金额)
- 简化处理:退费金额全部分配到第一个收款

### 当前Go版本的问题

```go
// 位置: approval_flow_service.go:282-840
func (s *ApprovalFlowService) processRefundApproval(refundOrderID int, orderID int) error {
    // ...
    // 3.3 执行冲回和重新分账
    if needChargeback {
        // 3.3a. 冲回所有未冲回的售卖类分账明细
        // 3.3b. 重新生成售卖类分账明细
        // ...
    }

    // 3.4 生成退费类分账明细(无论是否冲回都执行)  ← 注释是对的
    // 但实现有问题:

    // 第644行: 重复定义refundItemsList变量
    var refundItemsList []struct { ... }

    // 第665行: 重复定义allPayments变量
    allPayments := make([]struct { ... }, ...)

    // 第682-712行: 退费类分账生成逻辑
    // 这段逻辑结构上正确,但变量定义有问题
}
```

问题:
1. 第644-662行重复定义了`refundItemsList`和`allPayments`变量
2. 这些变量在冲回流程中已经定义(第467-480行和第533-551行)
3. 导致编译错误或逻辑混乱

## 解决方案

### 方案1: 删除重复变量定义,复用冲回流程中的变量

```go
// 3.4 生成退费类分账明细(无论是否冲回都执行)
// 直接使用冲回流程中查询的数据,或者在冲回流程外重新查询

// 如果执行了冲回,变量已存在
// 如果没有执行冲回,需要重新查询
var refundItemsListForRefundSeparate []struct {
    ChildOrderID int     `gorm:"column:childorder_id"`
    RefundAmount float64 `gorm:"column:refund_amount"`
    GoodsID      int     `gorm:"column:goodsid"`
    GoodsName    string  `gorm:"column:goods_name"`
}

// 重新查询退费子订单(确保数据最新)
if err := tx.Raw(`...`).Scan(&refundItemsListForRefundSeparate).Error; err != nil {
    return err
}

// 合并收款数据
allPaymentsForRefundSeparate := make([]struct {
    PaymentID   int
    PaymentType int
}, 0, len(regularPayments)+len(taobaoPayments))
// ...

// 生成退费类分账明细
for _, item := range refundItemsListForRefundSeparate {
    // ...
}
```

**优点**:
- 逻辑清晰,退费类分账独立于冲回流程
- 避免变量冲突
- 数据最新

**缺点**:
- 需要重新查询数据库

### 方案2: 在方法开始时统一查询所有需要的数据

```go
func (s *ApprovalFlowService) processRefundApproval(refundOrderID int, orderID int) error {
    // 1. 在开始时统一查询所有需要的数据
    refundItemsList := ...  // 查询退费子订单
    regularPayments := ...  // 查询常规收款
    taobaoPayments := ...   // 查询淘宝收款
    allPayments := ...      // 合并收款

    // 2. 更新状态
    // ...

    // 3. 检查是否需要冲回
    needChargeback := false
    // ...

    // 4. 如果需要冲回
    if needChargeback {
        // 使用统一查询的数据
    }

    // 5. 生成退费类分账明细(使用统一查询的数据)
    for _, item := range refundItemsList {
        // ...
    }

    // 6. 更新订单状态
    // ...
}
```

**优点**:
- 逻辑最清晰
- 数据统一管理
- 避免重复查询

**缺点**:
- 需要重构现有代码结构

## 推荐方案

**推荐方案1**,原因:
1. 修改最小,风险最低
2. 保持现有代码结构
3. 退费类分账逻辑独立,易于维护
4. 与原Python版本逻辑一致

## 实施细节

### 修改位置

文件: `internal/domain/approval/service/approval_flow_service.go`

方法: `processRefundApproval`

行号: 644-712

### 修改前代码

```go
// 3.4 生成退费类分账明细(无论是否冲回都执行)
// 获取所有退费子订单和退费金额
var refundItemsList []struct {
    ChildOrderID int     `gorm:"column:childorder_id"`
    RefundAmount float64 `gorm:"column:refund_amount"`
    GoodsID      int     `gorm:"column:goodsid"`
    GoodsName    string  `gorm:"column:goods_name"`
}
if err := tx.Raw(`...`).Scan(&refundItemsList).Error; err != nil {
    return err
}

// 合并所有收款用于退费类分账
allPayments := make([]struct {
    PaymentID   int
    PaymentType int
}, 0, len(regularPayments)+len(taobaoPayments))
// ...
```

### 修改后代码

```go
// 3.4 生成退费类分账明细(无论是否冲回都执行)
// 重新查询所有退费子订单和退费金额(确保数据最新且避免变量冲突)
var refundItemsForSeparate []struct {
    ChildOrderID int     `gorm:"column:childorder_id"`
    RefundAmount float64 `gorm:"column:refund_amount"`
    GoodsID      int     `gorm:"column:goodsid"`
    GoodsName    string  `gorm:"column:goods_name"`
}
if err := tx.Raw(`
    SELECT roi.childorder_id, roi.refund_amount, co.goodsid, g.name as goods_name
    FROM refund_order_item roi
    INNER JOIN childorders co ON roi.childorder_id = co.id
    LEFT JOIN goods g ON co.goodsid = g.id
    WHERE roi.refund_order_id = ?
    ORDER BY roi.childorder_id ASC
`, refundOrderID).Scan(&refundItemsForSeparate).Error; err != nil {
    return err
}

// 合并所有收款用于退费类分账
paymentsForRefundSeparate := make([]struct {
    PaymentID   int
    PaymentType int
}, 0, len(regularPayments)+len(taobaoPayments))
for _, p := range regularPayments {
    paymentsForRefundSeparate = append(paymentsForRefundSeparate, struct {
        PaymentID   int
        PaymentType int
    }{p.PaymentID, p.PaymentType})
}
for _, p := range taobaoPayments {
    paymentsForRefundSeparate = append(paymentsForRefundSeparate, struct {
        PaymentID   int
        PaymentType int
    }{p.PaymentID, p.PaymentType})
}

// 为每个退费子订单按收款顺序分配退费金额,生成退费类分账明细
for _, item := range refundItemsForSeparate {
    remainingRefund := item.RefundAmount

    // 按收款顺序分配退费金额(简化处理:退费金额全部分配到第一个收款)
    for _, payment := range paymentsForRefundSeparate {
        if remainingRefund <= 0 {
            break
        }

        allocatedRefund := remainingRefund

        // 插入退费类分账明细(负金额)
        if err := tx.Table("separate_account").Create(map[string]interface{}{
            "uid":             studentID,
            "orders_id":       orderID,
            "childorders_id":  item.ChildOrderID,
            "payment_id":      payment.PaymentID,
            "payment_type":    payment.PaymentType,
            "goods_id":        item.GoodsID,
            "goods_name":      item.GoodsName,
            "separate_amount": -allocatedRefund,
            "type":            2,
        }).Error; err != nil {
            return err
        }

        remainingRefund = 0 // 简化处理:退费金额全部分配到第一个收款
        break
    }
}
```

### 关键改动

1. **重命名变量**:
   - `refundItemsList` → `refundItemsForSeparate`
   - `allPayments` → `paymentsForRefundSeparate`

2. **重新查询数据**:
   - 退费子订单列表
   - 合并收款列表

3. **保持原逻辑**:
   - 为每个退费子订单生成退费类分账
   - 简化处理:退费金额全部分配到第一个收款
   - 插入负金额(type=2)

## 测试验证

### 测试场景1: 不需要冲回

**前置条件**:
- 订单有2个子订单:子订单A(500元),子订单B(500元)
- 有1个收款:收款1(1000元)
- 已生成分账:收款1→子订单A(500元),收款1→子订单B(500元)
- 创建退费订单:退费子订单A(200元)

**预期结果**:
- 不执行冲回流程(分账金额500元>=退费金额200元)
- 生成1条退费类分账:收款1→子订单A(-200元,type=2)
- 子订单A状态:部分支付(净分账300元)
- 子订单B状态:已支付(净分账500元)
- 订单状态:部分支付(净收款800元)

### 测试场景2: 需要冲回

**前置条件（购买时录单）**:
- **订单**:
  - id:1, 实收金额:600元
- **子订单**:
  - id:1, 实收金额:100元, 订单id:1
  - id:2, 实收金额:200元, 订单id:1
  - id:3, 实收金额:300元, 订单id:1
- **收款**:
  - id:1, 收款金额:400元
  - id:2, 收款金额:200元
- **初始分账明细**:
  - id:1, 子订单id:1, 收款id:1, 分账金额:100元, 类型:售卖
  - id:2, 子订单id:2, 收款id:1, 分账金额:200元, 类型:售卖
  - id:3, 子订单id:3, 收款id:1, 分账金额:100元, 类型:售卖
  - id:4, 子订单id:3, 收款id:2, 分账金额:200元, 类型:售卖

**退费操作**:
- **退费订单**:
  - id:1, 退费金额:500元
- **子退费订单（refund_order_item）**:
  - id:1, 子订单id:2, 退费金额:200元
  - id:2, 子订单id:3, 退费金额:300元
- **退费明细（refund_payment表，收款列表区填写）**:
  - id:1, 收款id:1, 退费金额:400元
  - id:2, 收款id:2, 退费金额:100元

**判断是否需要冲回**:
1. 拉取填入退费金额的收款关联的所有分账明细数据
2. 按收款id分组计算退费相关子订单的分账总额:
   - **收款1**在退费子订单上的分账总额: 300元（id:2的200元 + id:3的100元）
   - **收款2**在退费子订单上的分账总额: 200元（id:4的200元）
3. 比较分账总额与退费明细中填写的退费金额:
   - 收款1: 分账总额300元 < 退费金额400元 → **需要冲回**
   - 收款2: 分账总额200元 >= 退费金额100元 → 充足
4. **结论**: 存在收款1不足，需要对整个订单进行冲回

**预期结果**:

1. **冲回所有未冲回的售卖类分账明细**:
   - id:5, 子订单id:1, 收款id:1, 分账金额:-100元, 类型:冲回, parent_id:1
   - id:6, 子订单id:2, 收款id:1, 分账金额:-200元, 类型:冲回, parent_id:2
   - id:7, 子订单id:3, 收款id:1, 分账金额:-100元, 类型:冲回, parent_id:3
   - id:8, 子订单id:3, 收款id:2, 分账金额:-200元, 类型:冲回, parent_id:4

2. **第一批售卖分账: 将退费收款分配给退费子订单**:
   - 收款1退费400元 → 分配给子订单2(200元) + 子订单3(200元)
   - 收款2退费100元 → 分配给子订单3(100元)
   - 生成分账:
     - id:9, 子订单id:2, 收款id:1, 分账金额:200元, 类型:售卖
     - id:10, 子订单id:3, 收款id:1, 分账金额:200元, 类型:售卖
     - id:11, 子订单id:3, 收款id:2, 分账金额:100元, 类型:售卖

3. **第二批售卖分账: 将剩余收款分配给剩余子订单**:
   - 计算收款剩余:
     - 收款1剩余: 400元 - 400元 = 0元
     - 收款2剩余: 200元 - 100元 = 100元
   - 计算子订单剩余需求:
     - 子订单1剩余需求: 100元 - 0元 = 100元
     - 子订单2剩余需求: 200元 - 200元 = 0元
     - 子订单3剩余需求: 300元 - 200元 - 100元 = 0元
   - 生成分账:
     - id:12, 子订单id:1, 收款id:2, 分账金额:100元, 类型:售卖

4. **生成退费类分账（按照新的售卖分账分布）**:
   - 子订单2的售卖分账: 收款1有200元（id:9）
   - 子订单2退费200元 → 从收款1扣除200元:
     - id:13, 子订单id:2, 收款id:1, 分账金额:-200元, 类型:退费
   - 子订单3的售卖分账: 收款1有200元（id:10）, 收款2有100元（id:11）
   - 子订单3退费300元 → 从收款1扣除200元, 从收款2扣除100元:
     - id:14, 子订单id:3, 收款id:1, 分账金额:-200元, 类型:退费
     - id:15, 子订单id:3, 收款id:2, 分账金额:-100元, 类型:退费

5. **更新子订单状态（根据净分账金额）**:
   - 子订单1净分账: 100元 (id:12), 状态:已支付（100元=实收100元）
   - 子订单2净分账: 0元 (200元-200元), 状态:未支付
   - 子订单3净分账: 0元 (200元+100元-200元-100元), 状态:未支付

6. **更新订单状态**:
   - 总收款: 400元 + 200元 = 600元
   - 总退费: 400元 + 100元 = 500元
   - 净收款: 600元 - 500元 = 100元
   - 订单状态: 部分支付（100元 < 实收600元）

**关键验证点**:
- ✓ 判断冲回逻辑: 按收款比较分账总额与退费金额
- ✓ 冲回所有售卖分账（包括与退费无关的子订单1）
- ✓ 第一批分账优先满足退费子订单
- ✓ 第二批分账按顺序分配给剩余子订单
- ✓ 退费类分账按照子订单的售卖分账分布生成
- ✓ 净分账计算正确（售卖未冲回的 + 退费类）

## 参考资料

- 原Python版本代码: `D:\claude space\ZhixinStudentSaaS\backend\app.py` (第3816-4300行)
- 当前Go版本代码: `D:\claude space\CharonOMS\internal\domain\approval\service\approval_flow_service.go` (第282-840行)
- 分账明细规范: `D:\claude space\CharonOMS\openspec\specs\separate-account\spec.md`
