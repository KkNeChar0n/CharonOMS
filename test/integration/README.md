# 集成测试说明

## 概述

本目录包含订单管理模块的集成测试，测试完整的API端到端流程。

## 前提条件

1. MySQL数据库运行中
2. 数据库中已有测试数据（学生、商品等基础数据）
3. 修改 `order_test.go` 中的数据库连接字符串

## 运行测试

### 修改数据库连接

在 `order_test.go` 文件的 `setupTestEnv` 函数中，修改数据库连接字符串：

```go
dsn := "root:your_password@tcp(localhost:3306)/your_database?charset=utf8mb4&parseTime=True&loc=Local"
```

### 运行所有集成测试

```bash
go test ./test/integration/... -v
```

### 运行特定测试

```bash
go test ./test/integration/... -v -run TestOrderAPI_CreateOrder
```

## 测试覆盖

### 订单管理测试

- **TestOrderAPI_CreateOrder**: 测试创建订单
- **TestOrderAPI_GetOrders**: 测试获取订单列表
- **TestOrderAPI_UpdateOrder**: 测试更新订单
- **TestOrderAPI_SubmitOrder**: 测试提交订单
- **TestOrderAPI_CancelOrder**: 测试作废订单
- **TestOrderAPI_GetOrderGoods**: 测试获取订单商品列表

### 子订单管理测试

- **TestOrderAPI_GetChildOrders**: 测试获取子订单列表

### 商品相关测试

- **TestOrderAPI_GetActiveGoodsForOrder**: 测试获取启用商品列表
- **TestOrderAPI_CalculateDiscount**: 测试计算订单优惠

### 错误场景测试

- **TestOrderAPI_SubmitNonDraftOrder**: 测试重复提交订单（应该失败）
- **TestOrderAPI_UpdateNonDraftOrder**: 测试更新非草稿订单（应该失败）

## 注意事项

1. 集成测试会在数据库中创建真实数据，测试结束后会自动清理
2. 确保测试数据库中存在必要的基础数据（如学生ID=1，商品ID=1,2等）
3. 如果测试失败，可能需要手动清理测试数据
4. 建议在专门的测试数据库中运行集成测试，避免影响生产数据

## 测试数据清理

测试使用 `defer cleanupTestData(t, orderID)` 自动清理创建的订单数据，包括：
- orders 表中的订单记录
- childorders 表中的子订单记录
- orders_activity 表中的订单活动关联记录

如果测试中断导致数据未清理，可以手动执行：

```sql
-- 查找测试创建的订单（通常是最新的几条）
SELECT * FROM orders ORDER BY create_time DESC LIMIT 10;

-- 手动删除测试数据
DELETE FROM orders_activity WHERE orderid = ?;
DELETE FROM childorders WHERE parentsid = ?;
DELETE FROM orders WHERE id = ?;
```
