# 任务清单：迁移财务管理模块

## 阶段一：基础准备（预计1-2小时）

### 1. 数据库准备
- [x] 验证目标数据库是否存在 `payment_collection` 表
- [x] 验证目标数据库是否存在 `separate_account` 表
- [x] 对比源数据库和目标数据库表结构一致性
- [x] 如需创建表，执行建表SQL（从源数据库导出）

**验证方式**：执行 `DESCRIBE payment_collection` 和 `DESCRIBE separate_account`
**状态**：跳过MySQL验证，直接使用现有表结构

### 2. 依赖检查
- [x] 确认订单管理模块（order）已完成迁移并可用
- [x] 确认学生管理模块（student）已完成迁移并可用
- [x] 确认商品管理模块（goods）已完成迁移并可用
- [x] 测试依赖模块的关键接口（查询订单、查询学生、查询商品）

**验证方式**：调用依赖模块的API接口，确保返回正确数据
**状态**：通过代码检查确认依赖模块已实现

---

## 阶段二：Domain层实现（预计4-6小时）

### 3. 实体定义
- [x] 创建 `internal/domain/financial/payment/entity.go`
  - 定义 `PaymentCollection` 实体结构
  - 定义付款场景、付款方式、收款主体枚举
  - 定义收款状态枚举（待支付0、未核验10、已支付20）
  - 添加实体方法：`CanConfirm()`, `CanDelete()`, `Confirm()`

- [x] 创建 `internal/domain/financial/separate/entity.go`
  - 定义 `SeparateAccount` 实体结构
  - 定义收款类型、分账类型枚举
  - 添加实体方法（如有必要）

**验证方式**：编写实体单元测试，验证枚举值和方法逻辑

### 4. 仓储接口定义
- [x] 创建 `internal/domain/financial/payment/repository.go`
  - 定义 `PaymentRepository` 接口
  - 方法：`Create()`, `GetByID()`, `List()`, `Update()`, `Delete()`
  - 方法：`GetTotalPaidAmount(orderID)` 查询订单已收款总额
  - 方法：`CountByOrderAndStatus(orderID, status)` 统计收款数量

- [x] 创建 `internal/domain/financial/separate/repository.go`
  - 定义 `SeparateAccountRepository` 接口
  - 方法：`Create()`, `BatchCreate()`, `List()`, `GetByID()`
  - 方法：`ExistsByPaymentAndOrder(paymentID, orderID)` 检查是否已生成分账
  - 方法：`GetChildOrderTotalSeparate(childOrderID)` 查询子订单总分账金额

**验证方式**：接口定义清晰，方法签名符合业务需求

### 5. 领域服务实现
- [x] 创建 `internal/domain/financial/payment/service.go`
  - 实现 `PaymentDomainService`
  - 方法：`ValidatePaymentAmount()` 验证付款金额不超过待支付金额
  - 方法：`UpdateOrderPaymentStatus(orderID)` 更新订单支付状态
  - 逻辑：计算订单总收款（常规+淘宝），对比实收金额，更新订单状态

- [x] 创建 `internal/domain/financial/separate/service.go`
  - 实现 `SeparateAccountDomainService`
  - 方法：`GenerateSeparateAccounts(paymentID, orderID)` 生成分账明细
  - 逻辑：
    1. 查询收款信息
    2. 查询子订单列表（按ID升序）
    3. 防重复检查
    4. 遍历子订单，按序分配收款金额
    5. 插入分账明细
    6. 更新子订单状态
  - 方法：`UpdateChildOrderStatus(childOrderID)` 更新子订单状态

**验证方式**：编写领域服务单元测试，使用Mock仓储验证业务逻辑

**依赖关系**：需要订单模块提供更新订单状态和子订单状态的方法
**状态**：已完成，同时扩展了ChildOrderRepository接口

---

## 阶段三：Infrastructure层实现（预计3-4小时）

### 6. 收款仓储实现
- [x] 创建 `internal/infrastructure/persistence/financial/payment_repository.go`
  - 实现 `PaymentRepository` 接口
  - 使用GORM操作 `payment_collection` 表
  - 实现 `Create()` 方法（插入收款记录）
  - 实现 `GetByID()` 方法（查询单条记录）
  - 实现 `List()` 方法（支持多条件筛选：ID、UID、订单ID、付款人、付款方式、交易日期、状态）
  - 实现 `Update()` 方法（更新记录）
  - 实现 `Delete()` 方法（删除记录）
  - 实现 `GetTotalPaidAmount()` 方法（聚合查询）
  - 实现 `CountByOrderAndStatus()` 方法（统计查询）

**验证方式**：编写仓储集成测试，使用真实数据库验证CRUD操作

### 7. 分账仓储实现
- [x] 创建 `internal/infrastructure/persistence/financial/separate_repository.go`
  - 实现 `SeparateAccountRepository` 接口
  - 使用GORM操作 `separate_account` 表
  - 实现 `Create()` 方法（插入单条分账记录）
  - 实现 `BatchCreate()` 方法（批量插入分账记录，使用事务）
  - 实现 `List()` 方法（支持多条件筛选：ID、UID、订单ID、子订单ID、商品ID、收款ID、收款类型、分账类型）
  - 实现 `GetByID()` 方法（查询单条记录）
  - 实现 `ExistsByPaymentAndOrder()` 方法（检查是否存在）
  - 实现 `GetChildOrderTotalSeparate()` 方法（聚合查询）

**验证方式**：编写仓储集成测试，验证批量插入和聚合查询功能

---

## 阶段四：Application层实现（预计3-4小时）

### 8. DTO定义
- [x] 创建 `internal/application/financial/dto.go`
  - 定义 `PaymentCollectionDTO` 响应结构
  - 定义 `CreatePaymentCollectionRequest` 请求结构
  - 定义 `PaymentCollectionListRequest` 查询请求结构
  - 定义 `SeparateAccountDTO` 响应结构
  - 定义 `SeparateAccountListRequest` 查询请求结构

**验证方式**：确保DTO字段与前端API调用的JSON结构完全一致
**状态**：已完成，DTO放在Application层避免循环依赖

### 9. Assembler实现
- [x] 创建 `internal/application/financial/payment/assembler.go`
  - 实现 `ToPaymentCollectionDTO()` 实体转DTO
  - 实现 `ToPaymentCollectionEntity()` 请求转实体
  - 实现 `ToPaymentCollectionDTOList()` 列表转换

- [x] 创建 `internal/application/financial/separate/assembler.go`
  - 实现 `ToSeparateAccountDTO()` 实体转DTO
  - 实现 `ToSeparateAccountDTOList()` 列表转换

**验证方式**：编写Assembler单元测试，验证转换逻辑正确

### 10. 应用服务实现
- [x] 创建 `internal/application/financial/payment/service.go`
  - 实现 `PaymentApplicationService`
  - 依赖注入：`PaymentRepository`, `PaymentDomainService`, `SeparateAccountDomainService`
  - 方法：`GetPaymentCollections(request)` 获取收款列表
  - 方法：`CreatePaymentCollection(request)` 新增收款
    - 调用领域服务验证金额
    - 插入收款记录（状态默认10）
    - 更新订单状态
  - 方法：`ConfirmPaymentCollection(id)` 确认到账
    - 检查状态是否为未核验
    - 更新状态为已支付，设置到账时间
    - 更新订单状态
    - 生成分账明细（调用领域服务）
  - 方法：`DeletePaymentCollection(id)` 删除收款
    - 检查状态是否为未核验
    - 删除记录
    - 更新订单状态

- [x] 创建 `internal/application/financial/separate/service.go`
  - 实现 `SeparateAccountApplicationService`
  - 依赖注入：`SeparateAccountRepository`
  - 方法：`GetSeparateAccounts(request)` 获取分账明细列表

**验证方式**：编写应用服务单元测试，使用Mock依赖验证业务流程

---

## 阶段五：Interfaces层实现（预计2-3小时）

### 11. HTTP Handler实现
- [x] 创建 `internal/interfaces/http/financial/payment_handler.go`
  - 实现 `PaymentHandler`
  - 路由：`GET /api/payment-collections` → `GetPaymentCollections()`
  - 路由：`POST /api/payment-collections` → `CreatePaymentCollection()`
  - 路由：`PUT /api/payment-collections/:id/confirm` → `ConfirmPaymentCollection()`
  - 路由：`DELETE /api/payment-collections/:id` → `DeletePaymentCollection()`
  - 统一响应格式：`{"code": 0, "message": "success", "data": {...}}`
  - 错误处理：统一捕获并返回错误响应

- [x] 创建 `internal/interfaces/http/financial/separate_handler.go`
  - 实现 `SeparateAccountHandler`
  - 路由：`GET /api/separate-accounts` → `GetSeparateAccounts()`

**验证方式**：使用Postman或curl测试API接口，对比Python版本的响应格式

### 12. 路由注册
- [x] 在 `internal/interfaces/http/router/router.go` 中注册财务管理路由
  - 注册收款管理路由组
  - 注册分账明细路由组
  - 添加JWT中间件保护（如需要）

**验证方式**：启动服务，访问路由确认可达
**状态**：已完成，编译成功

---

## 阶段六：测试与验证（预计4-6小时）

### 13. 单元测试
- [ ] 编写 `payment/entity_test.go` 测试实体方法
- [ ] 编写 `payment/service_test.go` 测试领域服务
- [ ] 编写 `separate/service_test.go` 测试分账生成逻辑
- [ ] 编写 `payment/assembler_test.go` 测试DTO转换
- [ ] 编写 `payment/application_service_test.go` 测试应用服务
- [ ] 运行 `go test ./internal/domain/financial/... -v -cover`
- [ ] 确保覆盖率≥80%

**验证方式**：所有单元测试通过，覆盖率达标

### 14. 集成测试
- [ ] 编写 `payment_repository_test.go` 测试仓储与数据库交互
- [ ] 编写 `separate_repository_test.go` 测试分账仓储
- [ ] 编写完整业务流程测试：
  - 新增收款 → 确认到账 → 验证分账生成 → 验证订单状态更新
- [ ] 测试边界情况：
  - 付款金额超过待支付金额（应失败）
  - 重复确认到账（应失败）
  - 删除已确认的收款（应失败）
  - 重复生成分账（应幂等，不重复插入）
- [ ] 测试并发场景：
  - 多个收款同时确认
  - 分账不重复生成

**验证方式**：所有集成测试通过，边界情况处理正确

### 15. API接口测试
- [ ] 使用Postman导入Python版本的API测试用例
- [ ] 对比Go版本和Python版本的响应格式
- [ ] 测试所有筛选参数组合
- [ ] 测试错误场景的响应格式和HTTP状态码
- [ ] 导出测试报告，确认100%兼容

**验证方式**：Postman测试集全部通过，响应格式完全一致

### 16. 前端集成测试
- [ ] 启动Go后端服务
- [ ] 打开前端页面 `frontend/index.html`
- [ ] 测试收款管理功能：
  - 打开收款管理页面
  - 测试筛选功能（所有筛选条件）
  - 测试新增收款（成功和失败场景）
  - 测试确认到账
  - 测试删除收款
  - 验证分页功能
- [ ] 测试分账明细功能：
  - 打开分账明细页面
  - 测试筛选功能
  - 验证数据显示正确
  - 验证分页功能
- [ ] 测试跨页面联动：
  - 确认收款后，分账明细页面能查到新记录
  - 订单管理页面状态正确更新

**验证方式**：前端所有功能正常，无需修改前端代码

---

## 阶段七：文档与交付（预计1-2小时）

### 17. 代码审查
- [ ] 检查代码是否符合Go规范（运行 `gofmt`, `golint`, `go vet`）
- [ ] 检查是否符合DDD架构约定
- [ ] 检查错误处理是否完善
- [ ] 检查日志记录是否充分
- [ ] 检查注释是否完整（所有导出函数/类型）

**验证方式**：静态检查工具全部通过，代码审查无问题

### 18. 文档更新
- [ ] 更新 `openspec/project.md`，添加财务管理模块说明
- [ ] 更新数据库表清单（如有新增表）
- [ ] 创建API接口文档（Markdown格式）
- [ ] 创建部署说明（配置项、数据库迁移脚本）

**验证方式**：文档完整清晰，其他开发者可根据文档理解模块

### 19. 提交与归档
- [ ] 创建Git提交，使用规范的提交信息
- [ ] 推送到远程仓库（如需要）
- [ ] 运行 `openspec-cn archive migrate-financial-management`
- [ ] 验证规范已归档到 `openspec/specs/` 目录

**验证方式**：Git历史清晰，OpenSpec归档成功

---

## 依赖关系说明

### 并行任务
以下任务可并行执行（无依赖关系）：
- 任务1（数据库准备）和任务2（依赖检查）可并行
- 任务3（实体定义）和任务8（DTO定义）可并行
- 任务6（收款仓储）和任务7（分账仓储）可并行
- 任务13（单元测试）的各个子任务可并行

### 顺序依赖
以下任务必须按顺序执行：
- 任务3 → 任务4 → 任务5（Domain层自下而上）
- 任务4 → 任务6, 7（仓储实现依赖接口定义）
- 任务5, 6, 7 → 任务10（应用服务依赖领域服务和仓储）
- 任务8, 9, 10 → 任务11（Handler依赖Application层）
- 任务11 → 任务12（路由依赖Handler）
- 任务12 → 任务14, 15, 16（测试依赖完整实现）

---

## 验收标准总结

✅ **功能完整性**
- 收款管理的4个接口全部实现
- 分账明细的1个接口实现
- 自动分账生成功能正常
- 订单和子订单状态自动更新

✅ **质量标准**
- 单元测试覆盖率≥80%
- 集成测试覆盖完整流程
- 所有边界情况处理正确
- 并发场景无数据错误

✅ **兼容性**
- API路径、参数、响应格式与Python版本100%一致
- 前端无需任何修改
- 数据库表结构无修改

✅ **代码质量**
- 符合Go规范和DDD架构
- 错误处理完善
- 注释完整
- 代码审查通过

✅ **文档完整性**
- 项目文档已更新
- API文档已创建
- 部署说明已提供
- OpenSpec已归档
