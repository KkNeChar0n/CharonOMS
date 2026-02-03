# 提案：迁移财务管理模块（收款管理与分账明细）

## 概述

从 ZhixinStudentSaaS 项目迁移财务管理模块到 CharonOMS，包括**收款管理**和**分账明细**两个核心功能。使用 Golang 重写后端业务逻辑，保持前端代码和数据库结构不变，确保平滑迁移。

## 为什么

### 业务价值

财务管理是教育培训机构的核心业务环节，直接关系到资金流转和财务对账的准确性。收款管理和分账明细功能能够：

1. **准确记录资金流水**：详细记录每笔收款的来源、金额、时间、付款方式等信息
2. **自动化分账处理**：收款确认后自动按子订单分配，减少人工操作和错误
3. **实时状态同步**：订单和子订单状态根据收款和分账情况自动更新，保持数据一致性
4. **支持财务对账**：提供完整的收款和分账明细，便于财务核对和审计

### 迁移必要性

1. **完整业务闭环**：订单管理已完成迁移，但缺少收款和分账模块，业务流程不完整
2. **性能提升**：Python版本处理大量收款记录时性能不佳，Go版本可大幅提升并发处理能力
3. **代码维护性**：Python版本业务逻辑都在一个文件中，难以维护，Go版本采用DDD架构，职责清晰
4. **技术栈统一**：项目整体向Go迁移，保持技术栈一致性

### 技术优势

- **类型安全**：Go的静态类型系统可避免运行时类型错误
- **高性能**：Go的并发模型和编译型语言特性，处理高并发收款场景更高效
- **可测试性**：DDD架构便于编写单元测试和集成测试
- **可扩展性**：模块化设计便于后续扩展淘宝收款、退费管理等功能

## 背景

CharonOMS 项目已完成以下模块的迁移：
- 认证与授权（auth）
- 权限管理（rbac）
- 基础字典（basic）
- 学生管理（student）
- 教练管理（coach）
- 订单管理（order）
- 合同管理（contract）
- 审批流管理（approval）
- 商品管理（goods）
- 活动管理（activity）

财务管理是订单管理的下游模块，负责：
1. **收款记录管理**：记录订单的每笔收款，包括常规收款和淘宝收款
2. **分账明细生成**：收款确认到账后，自动按子订单顺序分配收款金额
3. **订单与子订单状态同步**：根据收款和分账情况自动更新订单状态

## 目标

1. **后端 Go 化**：使用 Golang 实现所有后端 API 和业务逻辑
2. **前端不变**：保持前端 HTML/CSS/JavaScript 代码完全不动
3. **数据库不变**：使用现有数据库表结构，不做任何修改
4. **API 兼容**：确保接口路径、请求参数、响应格式与 Python 版本完全一致
5. **业务一致**：实现相同的业务逻辑和验证规则
6. **遵循 DDD 架构**：按照项目约定的四层架构（Interfaces → Application → Domain → Infrastructure）组织代码

## 范围

### 包含功能

#### 1. 收款管理（payment_collection）
- 获取收款列表（支持多条件筛选）
- 新增收款记录（含待支付金额校验）
- 确认到账（状态变更 + 分账生成）
- 删除收款（仅未核验状态可删除）
- 订单支付状态更新（自动计算并同步）

#### 2. 分账明细（separate_account）
- 获取分账明细列表（支持多条件筛选）
- 自动生成分账（收款确认到账后触发）
- 按子订单顺序分配收款金额
- 子订单状态更新（根据分账金额自动计算）

### 不包含功能

- 淘宝收款管理（taobao_payment）
- 退费管理（refund_payment_details）
- 收款修改功能
- 分账手动调整功能
- 财务报表功能
- 对账功能

## 技术方案

### 1. 数据库表

使用现有表结构，无需修改：

#### payment_collection（收款表）
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键 |
| order_id | INT | 订单ID |
| student_id | INT | 学生UID |
| payment_scenario | INT | 付款场景：0=线上，1=线下 |
| payment_method | INT | 付款方式：0=微信，1=支付宝，2=优利，3=零零购，9=对公 |
| payment_amount | DECIMAL(10,2) | 付款金额 |
| payer | VARCHAR(100) | 付款人 |
| payee_entity | INT | 收款主体：0=北京，1=西安 |
| trading_hours | DATETIME | 交易时间 |
| arrival_time | DATETIME | 到账时间 |
| merchant_order | VARCHAR(100) | 商户订单号 |
| status | INT | 状态：0=待支付，10=未核验，20=已支付 |
| create_time | TIMESTAMP | 创建时间 |

#### separate_account（分账明细表）
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键 |
| uid | INT | 学生UID |
| orders_id | INT | 订单ID |
| childorders_id | INT | 子订单ID |
| payment_id | INT | 收款ID |
| payment_type | INT | 收款类型：0=常规，1=淘宝 |
| goods_id | INT | 商品ID |
| goods_name | VARCHAR(100) | 商品名称 |
| separate_amount | DECIMAL(10,2) | 分账金额 |
| type | INT | 类型：0=售卖，1=冲回，2=退费 |
| create_time | TIMESTAMP | 创建时间 |

### 2. 模块结构

按照 DDD 四层架构组织代码：

```
internal/
├── domain/
│   └── financial/           # 财务领域
│       ├── payment/         # 收款子领域
│       │   ├── entity.go            # PaymentCollection实体
│       │   ├── repository.go        # 仓储接口
│       │   └── service.go           # 领域服务（状态更新逻辑）
│       └── separate/        # 分账子领域
│           ├── entity.go            # SeparateAccount实体
│           ├── repository.go        # 仓储接口
│           └── service.go           # 领域服务（分账生成逻辑）
│
├── application/
│   └── financial/
│       ├── payment/
│       │   ├── service.go           # 收款应用服务
│       │   └── assembler.go         # DTO转换器
│       └── separate/
│           ├── service.go           # 分账应用服务
│           └── assembler.go         # DTO转换器
│
├── infrastructure/
│   └── persistence/
│       └── financial/
│           ├── payment_repository.go    # 收款仓储实现
│           └── separate_repository.go   # 分账仓储实现
│
└── interfaces/
    └── http/
        └── financial/
            ├── payment_handler.go       # 收款接口处理器
            ├── separate_handler.go      # 分账接口处理器
            └── dto.go                   # 数据传输对象
```

### 3. API 接口清单

#### 收款管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/payment-collections | 获取收款列表 |
| POST | /api/payment-collections | 新增收款 |
| PUT | /api/payment-collections/:id/confirm | 确认到账 |
| DELETE | /api/payment-collections/:id | 删除收款 |

#### 分账明细
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/separate-accounts | 获取分账明细列表 |

### 4. 核心业务逻辑

#### 4.1 新增收款验证
```go
// 1. 查询订单实收金额
orderAmount := getOrderActualAmount(orderID)

// 2. 计算已收款总额（状态为10或20的常规收款 + 状态为30的淘宝收款）
paidAmount := getTotalPaidAmount(orderID)

// 3. 计算待支付金额
unpaidAmount := orderAmount - paidAmount

// 4. 校验：付款金额不能超过待支付金额
if paymentAmount > unpaidAmount {
    return errors.New("付款金额不能超过待支付金额")
}

// 5. 插入收款记录，状态默认为10（未核验）
collection := PaymentCollection{
    Status: StatusUnverified, // 10
    ...
}
```

#### 4.2 确认到账流程
```go
// 1. 检查状态（必须是未核验10）
if collection.Status != StatusUnverified {
    return errors.New("只能确认未核验的收款")
}

// 2. 更新状态为20（已支付），设置到账时间
collection.Status = StatusPaid
collection.ArrivalTime = time.Now()

// 3. 更新订单支付状态
updateOrderPaymentStatus(collection.OrderID)

// 4. 生成分账明细
generateSeparateAccounts(collection.ID, collection.OrderID)
```

#### 4.3 分账生成算法
```go
// 1. 获取收款信息
collection := getPaymentCollection(paymentID)

// 2. 查询子订单列表（按ID升序）
childOrders := getChildOrders(orderID)

// 3. 防重复检查
if existsSeparateAccount(paymentID, orderID) {
    return nil
}

// 4. 遍历子订单，按序分配收款
remainingAmount := collection.PaymentAmount

for _, child := range childOrders {
    // 计算子订单已分配金额
    allocatedAmount := getChildOrderAllocatedAmount(child.ID)

    // 计算还需金额
    neededAmount := child.ActualAmount - allocatedAmount

    if neededAmount <= 0 {
        continue
    }

    // 分配金额 = min(剩余收款, 还需金额)
    separateAmount := min(remainingAmount, neededAmount)

    // 插入分账明细
    insertSeparateAccount(SeparateAccount{
        UID: collection.StudentID,
        OrdersID: collection.OrderID,
        ChildOrdersID: child.ID,
        PaymentID: paymentID,
        PaymentType: 0, // 常规收款
        GoodsID: child.GoodsID,
        GoodsName: child.GoodsName,
        SeparateAmount: separateAmount,
        Type: 0, // 售卖
    })

    // 更新剩余金额
    remainingAmount -= separateAmount

    // 更新子订单状态
    updateChildOrderStatus(child.ID)

    if remainingAmount <= 0 {
        break
    }
}
```

#### 4.4 订单状态更新逻辑
```go
// 计算订单总收款金额
totalPaid := getTotalPaidAmount(orderID) // 常规收款 + 淘宝收款

// 获取订单实收金额
actualAmount := getOrderActualAmount(orderID)

// 根据收款情况更新状态
if totalPaid == 0 {
    orderStatus = StatusUnpaid // 20-未支付
} else if totalPaid < actualAmount {
    orderStatus = StatusPartiallyPaid // 30-部分支付
} else {
    orderStatus = StatusPaid // 40-已支付
}
```

#### 4.5 子订单状态更新逻辑
```go
// 计算子订单总分账金额
totalSeparate := getChildOrderTotalSeparate(childOrderID)

// 获取子订单实收金额
actualAmount := getChildOrderActualAmount(childOrderID)

// 根据分账情况更新状态
if totalSeparate == 0 {
    childOrderStatus = StatusUnpaid // 10-未支付
} else if totalSeparate < actualAmount {
    childOrderStatus = StatusPartiallyPaid // 20-部分支付
} else {
    childOrderStatus = StatusPaid // 30-已支付
}
```

### 5. 响应格式

保持与 Python 版本一致：

#### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "collections": [...]
  }
}
```

#### 错误响应
```json
{
  "code": 1,
  "message": "付款金额不能超过待支付金额",
  "data": null
}
```

### 6. 前端保持不变

- HTML结构：`frontend/index.html` 的收款管理和分账明细页面部分
- JavaScript逻辑：`frontend/app.js` 的相关方法
- CSS样式：`frontend/styles.css`

前端调用的API路径、请求参数、响应格式保持完全一致，无需任何修改。

## 依赖关系

### 前置依赖
- 订单管理模块（order）：收款关联订单
- 学生管理模块（student）：收款关联学生
- 商品管理模块（goods）：分账记录商品信息

### 被依赖
- 无（当前阶段）

## 风险与挑战

### 1. 业务逻辑复杂性
- **风险**：分账生成逻辑涉及多表查询和事务处理，逻辑复杂
- **应对**：
  - 编写详细的单元测试，覆盖各种边界情况
  - 使用GORM事务确保数据一致性
  - 添加幂等性检查，防止重复生成分账

### 2. 金额计算精度
- **风险**：浮点数运算可能导致精度丢失
- **应对**：
  - 使用`decimal.Decimal`类型处理金额
  - 数据库使用`DECIMAL(10,2)`类型存储
  - 金额比较时使用精确比较方法

### 3. API兼容性
- **风险**：前端不修改，后端必须100%兼容原API
- **应对**：
  - 详细对比Python版本的请求/响应格式
  - 编写集成测试验证API兼容性
  - 使用Postman导出测试用例对比

### 4. 并发安全
- **风险**：多个收款同时确认可能导致分账错误
- **应对**：
  - 在生成分账时使用数据库悲观锁或乐观锁
  - 添加重复检查机制
  - 记录详细日志便于问题排查

## 成功标准

1. ✅ 所有API接口路径、参数、响应格式与Python版本一致
2. ✅ 前端页面无需修改即可正常工作
3. ✅ 收款流程：新增→确认→分账，全流程正常
4. ✅ 订单和子订单状态自动更新准确
5. ✅ 金额计算精确无误差
6. ✅ 单元测试覆盖率达到80%以上
7. ✅ 集成测试验证完整业务流程
8. ✅ 代码遵循DDD架构规范
9. ✅ 通过前端手工测试验证

## 后续扩展

本次迁移完成后，未来可扩展：
- 淘宝收款管理
- 退费管理
- 财务报表功能
- 对账功能
- 收款修改和作废
- 分账手动调整

## 验收检查清单

- [ ] 所有API接口实现并测试通过
- [ ] 前端无需修改即可调用新后端
- [ ] 数据库表结构无修改
- [ ] 新增收款功能正常（含金额校验）
- [ ] 确认到账功能正常（含分账生成）
- [ ] 删除收款功能正常（含状态校验）
- [ ] 分账自动生成正确（按子订单顺序）
- [ ] 订单状态自动更新准确
- [ ] 子订单状态自动更新准确
- [ ] 金额计算精确（无精度丢失）
- [ ] 并发场景下分账不重复
- [ ] 单元测试覆盖率≥80%
- [ ] 集成测试覆盖完整业务流程
- [ ] 代码符合Go规范和DDD架构
- [ ] 代码审查通过

## 参考资料

- 源项目：`D:\claude space\ZhixinStudentSaaS`
- 目标项目：`D:\claude space\CharonOMS`
- 项目约定：`openspec/project.md`
