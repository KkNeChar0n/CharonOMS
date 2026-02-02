# Change: 迁移订单管理功能

## Why

从原Python项目(ZhixinStudentSaaS)迁移订单管理功能到新的Go版本(CharonOMS)。订单管理是教育培训机构管理系统的核心业务模块,负责管理课程订单、商品购买、收款、退款等核心业务流程。

当前阶段优先迁移：
1. **订单管理**：订单的创建、编辑、提交、作废等生命周期管理
2. **子订单管理**：订单商品明细及优惠分摊

## What Changes

### 新增功能
- **订单实体和仓储**：创建Order和ChildOrder领域实体及其仓储接口
- **订单状态管理**：实现订单状态流转逻辑（草稿→未支付→部分支付→已支付→退费中→已作废）
- **订单CRUD接口**：
  - 创建订单（支持多商品、优惠分摊、活动关联）
  - 查询订单列表（含学生、金额、状态信息）
  - 获取订单商品详情
  - 编辑草稿订单
  - 提交订单（草稿→未支付）
  - 作废草稿订单
- **子订单管理**：
  - 查询子订单列表
  - 子订单与订单的级联操作
- **商品相关接口**：
  - 获取启用商品列表（用于订单选择）
  - 获取商品总价计算
- **订单优惠计算**：基于活动模板的满折优惠计算及子订单优惠分摊

### 技术约束
- **语言迁移**：Python Flask → Go + Gin
- **数据库不变**：使用现有数据库表结构，不修改数据结构
- **API兼容**：保持与Python版本相同的API路径和响应格式
- **架构模式**：采用DDD分层架构（Interface → Application → Domain → Infrastructure）

### 暂不迁移
- 支付管理（payment_collection, taobao_payment）
- 分账管理（separate_account）
- 退款管理（refund_order, refund_order_item）
- 这些功能将在后续阶段迁移

## Impact

### 影响的规范
- **order**（新增）：订单管理核心功能
- **goods**（扩展）：新增订单相关的商品查询接口
- **activity**（依赖）：订单需要活动数据进行优惠计算
- **student**（依赖）：订单关联学生信息

### 影响的代码
- `internal/domain/order/`：新增订单领域模块
  - `entity/order.go`：订单实体
  - `entity/childorder.go`：子订单实体
  - `repository/order_repository.go`：订单仓储接口
  - `service/order_service.go`：订单领域服务
- `internal/application/order/`：新增订单应用服务
  - `order_service.go`：订单业务编排
  - `assembler.go`：DTO转换器
- `internal/interfaces/http/order/`：新增订单HTTP接口
  - `order_handler.go`：订单接口处理器
  - `dto.go`：请求响应DTO
- `internal/infrastructure/persistence/order/`：新增订单持久化
  - `order_repository_impl.go`：订单仓储实现
- `cmd/server/router.go`：新增订单路由配置

### 数据库表
使用现有表结构：
- `orders`：订单主表
- `childorders`：子订单表
- `orders_activity`：订单活动关联表
- `goods`：商品表（已存在）
- `student`：学生表（已存在）
- `activity`：活动表（已存在）

### 迁移风险
- **API兼容性**：必须与Python版本API完全兼容，确保前端无缝切换
- **金额计算精度**：订单金额涉及decimal类型，需确保计算精度
- **事务一致性**：订单创建/编辑涉及多表操作，需使用事务保证一致性
- **状态流转逻辑**：必须严格遵循原系统的状态流转规则
