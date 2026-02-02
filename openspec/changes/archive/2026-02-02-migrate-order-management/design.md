# 订单管理迁移设计文档

## Context

### 背景
从原Python Flask项目(ZhixinStudentSaaS)迁移订单管理功能到新的Go版本(CharonOMS)。原系统使用单体应用架构，直接在Flask路由中编写业务逻辑和SQL查询；新系统采用DDD分层架构，需要按照领域驱动设计的原则重新组织代码。

### 约束条件
1. **数据库不可变**：必须使用现有的MySQL数据库表结构，不能修改字段名称、类型或表关系
2. **API兼容性**：必须保持与Python版本完全相同的API路径、请求参数和响应格式，确保前端无缝切换
3. **前端不变**：前端Vue.js代码不做任何修改，完全依赖后端API兼容
4. **语言切换**：从Python切换到Go，需要处理语言差异（如decimal精度、时间格式等）

### 利益相关者
- **开发团队**：需要清晰的架构指导和迁移路径
- **前端开发**：依赖API兼容性保证
- **业务用户**：期望功能无缝迁移，无感知切换

## Goals / Non-Goals

### Goals
1. **完整迁移订单基础功能**：
   - 订单CRUD（创建、查询、编辑、提交、作废）
   - 子订单管理
   - 订单商品关联
   - 订单优惠计算
   - 订单状态流转

2. **保证数据一致性**：
   - 订单创建使用数据库事务
   - 订单编辑时级联更新子订单和活动关联
   - 状态流转遵循业务规则约束

3. **保持API完全兼容**：
   - 响应格式：`{"orders": [...]}`, `{"childorders": [...]}`等
   - 字段命名：使用Python版本的snake_case命名（如amount_receivable）
   - 状态码：与Python版本保持一致

4. **采用DDD分层架构**：
   - 清晰的领域边界
   - 依赖倒置原则
   - 业务逻辑集中在领域层

### Non-Goals
1. **不迁移支付管理**：支付收款、支付验证功能后续迁移
2. **不迁移分账管理**：分账生成和更新逻辑后续迁移
3. **不迁移退款管理**：退款申请、审批流程后续迁移
4. **不优化数据库结构**：保持与原系统完全一致
5. **不修改业务逻辑**：严格按照原系统逻辑实现

## Decisions

### 决策1：DDD分层架构设计

**选择**：采用经典四层DDD架构
```
Interfaces Layer (HTTP Handler, DTO)
    ↓
Application Layer (Business Orchestration)
    ↓
Domain Layer (Entity, Repository Interface, Domain Service)
    ↓
Infrastructure Layer (Repository Implementation, Database)
```

**理由**：
- 符合项目现有架构模式（auth, student, coach等模块已采用）
- 清晰的职责分离，便于维护和测试
- 依赖倒置，领域层不依赖基础设施

**替代方案及理由**：
- **MVC架构**：过于简单，不适合复杂业务逻辑，领域逻辑容易分散
- **六边形架构**：过度设计，当前项目规模不需要如此复杂的端口适配器模式

### 决策2：订单聚合根设计

**选择**：Order作为聚合根，ChildOrder作为实体，通过Order管理ChildOrder的生命周期

**聚合边界**：
```
Order (聚合根)
├── ChildOrder (实体列表)
└── OrdersActivity (关联关系)
```

**理由**：
- 子订单不能独立于订单存在，必须通过订单管理
- 订单编辑时需要级联删除和重建子订单
- 订单状态变更需要级联更新子订单状态
- 事务边界清晰：一个订单的所有操作在一个事务内完成

**替代方案及理由**：
- **ChildOrder独立聚合**：不合理，子订单的生命周期完全依赖订单，独立管理会破坏业务一致性

### 决策3：金额类型处理

**选择**：使用 `float64` 存储金额，在DTO层进行精度控制

**理由**：
- MySQL中金额字段为DECIMAL(10,2)
- Go的float64可以满足两位小数精度要求
- GORM原生支持float64与DECIMAL的映射
- 前端显示时统一保留两位小数

**注意事项**：
- 避免直接比较float64相等性，使用差值判断
- 金额计算后四舍五入到两位小数
- 金额验证时允许0.01的误差范围

**替代方案及理由**：
- **decimal库（shopspring/decimal）**：更精确但增加复杂性，当前业务场景（最多两位小数）float64足够
- **int64（分为单位）**：需要大量转换代码，增加复杂度，且与现有数据库不兼容

### 决策4：状态流转管理

**选择**：在Order实体中实现状态流转验证方法

**实现方式**：
```go
func (o *Order) CanEdit() bool {
    return o.Status == StatusDraft
}

func (o *Order) CanSubmit() bool {
    return o.Status == StatusDraft
}

func (o *Order) CanCancel() bool {
    return o.Status == StatusDraft
}
```

**理由**：
- 状态流转规则是核心业务逻辑，应该在领域实体中实现
- 集中管理状态流转，避免在多个地方重复判断
- 便于单元测试和维护

**替代方案及理由**：
- **状态机库（looplab/fsm）**：过度设计，当前状态流转规则简单，不需要引入第三方库
- **应用层判断**：违反DDD原则，业务规则应该在领域层

### 决策5：优惠计算策略

**选择**：在Domain Service中实现优惠计算逻辑

**实现位置**：`internal/domain/order/service/discount_service.go`

**理由**：
- 优惠计算涉及Order、Activity、Goods多个实体的协作
- 逻辑复杂，不适合放在单一实体中
- 需要访问Activity仓储获取活动规则

**计算流程**：
1. 根据activity_ids查询活动模板和规则
2. 筛选满折类型活动（type=2）
3. 按活动分组方式（按商品/按分类）筛选参与商品
4. 匹配折扣档位（根据参与商品数量）
5. 计算总优惠金额
6. 平均分摊优惠到各子订单

**替代方案及理由**：
- **应用层计算**：不合理，优惠计算是核心业务逻辑，应该在领域层
- **Order实体方法**：职责过重，Order实体不应该知道Activity的细节

### 决策6：订单编辑实现策略

**选择**：删除旧子订单，重新创建新子订单（DELETE + INSERT）

**理由**：
- 与Python版本保持一致
- 实现简单，避免复杂的diff逻辑
- 子订单ID对前端不重要（前端展示使用订单ID）
- 数据量小（每个订单通常只有几个子订单），性能影响可忽略

**替代方案及理由**：
- **UPDATE现有子订单**：需要复杂的diff算法，增加出错风险，且与原系统不一致

### 决策7：事务管理策略

**选择**：在仓储层实现事务管理，应用层调用仓储方法

**实现方式**：
```go
// Repository接口
type OrderRepository interface {
    CreateOrder(ctx context.Context, order *Order, childOrders []*ChildOrder, activityIDs []int) error
    // CreateOrder内部使用GORM事务
}
```

**理由**：
- GORM提供了便捷的事务API（db.Transaction）
- 订单创建涉及多表操作（orders, childorders, orders_activity），必须保证原子性
- 仓储层统一管理事务，应用层不需要关心事务细节

**替代方案及理由**：
- **应用层管理事务**：需要传递db对象到仓储层，增加耦合
- **不使用事务**：风险太高，可能导致数据不一致

## Architecture

### 模块结构

```
internal/
├── domain/
│   └── order/
│       ├── entity/
│       │   ├── order.go              # 订单实体
│       │   └── childorder.go         # 子订单实体
│       ├── repository/
│       │   ├── order_repository.go   # 订单仓储接口
│       │   └── childorder_repository.go
│       └── service/
│           ├── order_service.go      # 订单领域服务
│           └── discount_service.go   # 优惠计算服务
├── application/
│   └── order/
│       ├── order_service.go          # 订单应用服务
│       └── assembler.go              # DTO转换器
├── interfaces/
│   └── http/
│       └── order/
│           ├── order_handler.go      # HTTP处理器
│           └── dto.go                # 请求响应DTO
└── infrastructure/
    └── persistence/
        └── order/
            ├── order_repository_impl.go    # 仓储实现
            ├── childorder_repository_impl.go
            └── models.go                   # GORM模型
```

### 数据流

#### 创建订单流程
```
HTTP Request (CreateOrderRequest)
    ↓
OrderHandler.CreateOrder()
    ↓
OrderService.CreateOrder() (Application)
    ├─→ DiscountService.CalculateDiscount() (Domain)
    ├─→ OrderService.AllocateChildDiscounts() (Domain)
    └─→ OrderRepository.CreateOrder() (Infrastructure)
        └─→ GORM Transaction
            ├─→ INSERT INTO orders
            ├─→ INSERT INTO childorders (batch)
            └─→ INSERT INTO orders_activity (batch)
    ↓
HTTP Response (OrderResponse)
```

#### 查询订单列表流程
```
HTTP Request
    ↓
OrderHandler.GetOrders()
    ↓
OrderService.GetOrders() (Application)
    ↓
OrderRepository.GetOrders() (Infrastructure)
    └─→ SELECT orders JOIN student
    ↓
Assembler.ToOrderDTOs() (Application)
    ↓
HTTP Response ({"orders": [...]})
```

### 关键类图

```
┌─────────────────┐
│     Order       │ (Aggregate Root)
├─────────────────┤
│ - ID            │
│ - StudentID     │
│ - ExpectedPaymentTime │
│ - AmountReceivable    │
│ - AmountReceived      │
│ - DiscountAmount      │
│ - Status        │
├─────────────────┤
│ + CanEdit()     │
│ + CanSubmit()   │
│ + CanCancel()   │
│ + ValidateAmounts() │
└─────────────────┘
        │
        │ 1
        │
        │ *
        ↓
┌─────────────────┐
│  ChildOrder     │ (Entity)
├─────────────────┤
│ - ID            │
│ - ParentsID     │
│ - GoodsID       │
│ - AmountReceivable │
│ - AmountReceived   │
│ - DiscountAmount   │
│ - Status        │
└─────────────────┘

┌──────────────────────┐
│ OrderRepository      │ (Interface)
├──────────────────────┤
│ + GetOrders()        │
│ + GetOrderByID()     │
│ + CreateOrder()      │
│ + UpdateOrder()      │
│ + GetOrderGoods()    │
└──────────────────────┘
         △
         │ implements
         │
┌──────────────────────┐
│ OrderRepositoryImpl  │ (Infrastructure)
├──────────────────────┤
│ - db *gorm.DB        │
└──────────────────────┘
```

## Risks / Trade-offs

### 风险1：金额精度问题
**风险**：使用float64可能导致精度丢失
**缓解措施**：
- 在DTO序列化时强制保留两位小数
- 金额比较时使用误差范围判断（0.01）
- 添加单元测试验证金额计算精度

### 风险2：API兼容性
**风险**：Go与Python的时间格式、null值处理不同，可能导致前端解析失败
**缓解措施**：
- 使用GORM的自定义序列化（json标签）
- 时间格式统一使用RFC3339（与Python datetime.isoformat()兼容）
- null值统一返回null而非0或空字符串
- 编写集成测试对比Python版本响应

### 风险3：状态流转逻辑遗漏
**风险**：Python版本可能存在未文档化的状态流转规则
**缓解措施**：
- 仔细review Python代码，找出所有状态判断逻辑
- 编写单元测试覆盖所有状态流转场景
- 前端集成测试，验证所有操作是否正常

### 风险4：性能问题
**风险**：查询订单列表时关联查询可能较慢
**缓解措施**：
- 使用GORM的Preload或Joins减少N+1查询
- 如有性能问题，后续添加索引（注意：当前阶段不修改数据库）
- 添加分页查询（如订单数量增长）

### Trade-off：DELETE + INSERT vs UPDATE
**选择**：DELETE + INSERT
**优点**：
- 实现简单，不需要diff算法
- 与Python版本一致
- 不会有遗留数据

**缺点**：
- 子订单ID会变化
- 多一次DELETE操作

**评估**：优点大于缺点，子订单ID不重要，性能影响可忽略

## Migration Plan

### 阶段1：领域层和基础设施层（3天）
1. 创建实体定义（Order, ChildOrder）
2. 创建仓储接口
3. 实现仓储（含事务管理）
4. 编写单元测试

### 阶段2：应用层和接口层（4天）
1. 创建应用服务（订单业务编排）
2. 创建HTTP Handler和DTO
3. 配置路由
4. 编写集成测试

### 阶段3：优惠计算和验证（2天）
1. 实现优惠计算服务
2. 实现优惠分摊逻辑
3. 对比Python版本API响应
4. 前端联调测试

### 阶段4：文档和部署（1天）
1. 更新API文档
2. 编写迁移说明
3. 部署到测试环境
4. 验收测试

### 回滚计划
如果迁移失败，可以：
1. 通过Nginx配置切换回Python服务
2. 保留Python版本作为fallback
3. 数据库不受影响（未修改结构）

## Open Questions

1. **订单编号规则**：原系统是否有订单编号生成规则？当前仅使用自增ID
   - **建议**：保持使用自增ID，如需编号规则后续添加

2. **分页查询**：订单列表是否需要分页？
   - **建议**：初期不实现，如订单数量>1000再添加

3. **软删除**：订单是否需要支持软删除（deleted_at）？
   - **建议**：当前阶段不支持，使用状态99（已作废）表示删除

4. **并发控制**：订单编辑是否需要乐观锁？
   - **建议**：初期不实现，如有并发冲突再添加版本号字段

5. **日志审计**：订单状态变更是否需要记录审计日志？
   - **建议**：后续添加，当前使用统一日志中间件记录请求
