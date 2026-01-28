# 任务清单：实现合同管理模块

## Phase 1: 数据库Schema验证和准备
- [x] 验证contract表是否存在
  - [x] 检查数据库中是否已有contract表
  - [x] 如不存在，准备创建脚本
  - [x] 验证字段完整性（id, name, student_id, type, signature_form, contract_amount, signatory, initiating_party, initiator, status, payment_status, termination_agreement, create_time）
- [x] 验证外键约束
  - [x] 检查是否存在student_id外键约束
  - [x] 如不存在，添加外键约束
  - [x] 验证外键引用student表的id字段
- [x] 准备测试数据（可选）
  - [x] 插入3-5条测试合同数据
  - [x] 测试数据包含不同类型、签署形式、状态
  - [x] 关联不同的学生

## Phase 2: Domain层实现
- [x] 创建合同实体模型
  - [x] 编写 `internal/domain/contract/entity/contract.go`
  - [x] 定义Contract结构体（包含所有字段）
  - [x] 添加GORM标签和JSON标签
  - [x] 定义TableName()方法

- [x] 创建仓储接口
  - [x] 编写 `internal/domain/contract/repository/contract_repository.go`
  - [x] 定义ContractRepository接口
  - [x] 方法清单：
    - [x] GetContractList() - 获取合同列表（含学生信息）
    - [x] GetContractByID() - 根据ID查询合同详情
    - [x] CreateContract() - 创建合同
    - [x] RevokeContract() - 撤销合同（更新status为98）
    - [x] TerminateContract() - 中止合作（更新status为99，保存termination_agreement）
    - [x] ContractExists() - 检查合同是否存在
    - [x] GetContractStatus() - 获取合同当前状态

## Phase 3: Infrastructure层实现
- [x] 创建仓储实现
  - [x] 编写 `internal/infrastructure/persistence/mysql/contract/contract_repository_impl.go`
  - [x] 实现ContractRepositoryImpl结构体
  - [x] 实现NewContractRepository构造函数

- [x] 实现查询方法
  - [x] 实现GetContractList（使用LEFT JOIN查询）
    - [x] JOIN student表获取学生姓名
    - [x] 返回完整的合同信息列表
    - [x] 按create_time倒序排列
  - [x] 实现GetContractByID（包含学生信息）

- [x] 实现写入方法
  - [x] 实现CreateContract
    - [x] 插入contract记录
    - [x] 设置status=0, payment_status=0
    - [x] 设置发起人（从context获取）
    - [x] 设置创建时间
  - [x] 实现RevokeContract
    - [x] 验证当前status必须为0
    - [x] 更新status为98
  - [x] 实现TerminateContract
    - [x] 验证当前status必须为50
    - [x] 更新status为99
    - [x] 保存termination_agreement路径

- [x] 实现辅助方法
  - [x] 实现ContractExists（验证合同是否存在）
  - [x] 实现GetContractStatus（获取合同当前状态）

## Phase 4: Application层实现
- [x] 创建DTO定义
  - [x] 编写 `internal/application/service/contract/dto.go`
  - [x] ContractListResponse - 合同列表响应
  - [x] ContractDetailDTO - 合同详情
  - [x] CreateContractRequest - 创建合同请求
  - [x] RevokeContractRequest - 撤销合同请求（可能不需要额外字段）
  - [x] TerminateContractRequest - 中止合作请求（包含termination_agreement）

- [x] 创建业务服务
  - [x] 编写 `internal/application/service/contract/contract_service.go`
  - [x] 定义ContractService结构体
  - [x] 实现NewContractService构造函数

- [x] 实现业务方法
  - [x] GetContractList - 获取合同列表
  - [x] GetContractByID - 获取合同详情
  - [x] CreateContract - 创建合同
    - [x] 验证必填字段（name, student_id, type, signature_form, contract_amount）
    - [x] 验证type值（0或1）
    - [x] 验证signature_form值（0或1）
    - [x] 调用仓储创建合同
  - [x] RevokeContract - 撤销合同
    - [x] 检查合同是否存在
    - [x] 验证当前status为0（待审核）
    - [x] 调用仓储更新状态
  - [x] TerminateContract - 中止合作
    - [x] 检查合同是否存在
    - [x] 验证当前status为50（已通过）
    - [x] 验证termination_agreement不为空
    - [x] 调用仓储更新状态和文件路径

## Phase 5: Interface层实现
- [x] 创建HTTP处理器
  - [x] 编写 `internal/interfaces/http/handler/contract/contract_handler.go`
  - [x] 定义ContractHandler结构体
  - [x] 实现NewContractHandler构造函数

- [x] 实现HTTP处理方法
  - [x] GetContracts - GET /api/contracts
    - [x] 调用业务服务
    - [x] 返回格式：`{"contracts": [...]}`
  - [x] GetContractByID - GET /api/contracts/:id
    - [x] 解析路径参数
    - [x] 调用业务服务
    - [x] 返回单个合同详情
  - [x] CreateContract - POST /api/contracts
    - [x] 绑定JSON请求体
    - [x] 从JWT token获取当前用户（initiator）
    - [x] 调用业务服务
    - [x] 返回201和合同ID
    - [x] 返回消息："合同新增成功"
  - [x] RevokeContract - PUT /api/contracts/:id/revoke
    - [x] 解析路径参数
    - [x] 调用业务服务
    - [x] 返回成功消息
  - [x] TerminateContract - PUT /api/contracts/:id/terminate
    - [x] 解析路径参数
    - [x] 绑定JSON请求体（termination_agreement）
    - [x] 调用业务服务
    - [x] 返回成功消息

- [x] 更新路由注册
  - [x] 修改 `internal/interfaces/http/router/router.go`
  - [x] 导入contract相关包
  - [x] 初始化ContractRepository
  - [x] 初始化ContractService
  - [x] 初始化ContractHandler
  - [x] 替换占位符路由为真实路由
    - [x] GET /contracts
    - [x] GET /contracts/:id
    - [x] POST /contracts
    - [x] PUT /contracts/:id/revoke
    - [x] PUT /contracts/:id/terminate

## Phase 6: 文件上传功能实现（中止合作需要）
- [x] 实现文件上传处理
  - [x] 编写 `internal/infrastructure/storage/file_upload.go`
  - [x] 实现文件保存到服务器
  - [x] 返回文件路径
  - [x] 配置文件存储目录

- [x] 集成到中止合作接口
  - [x] 在TerminateContract接口中处理文件上传
  - [x] 验证文件类型和大小
  - [x] 保存文件并获取路径
  - [x] 将路径存储到termination_agreement字段

## Phase 7: 测试验证

- [x] 单元测试（可选）
  - [x] 测试ContractService业务逻辑
  - [x] 测试字段验证
  - [x] 测试状态流转规则

- [x] 集成测试
  - [x] 测试创建合同（基本字段）
  - [x] 测试查询合同列表
  - [x] 测试查询合同详情
  - [x] 测试撤销待审核合同
  - [x] 测试撤销非待审核合同（应失败）
  - [x] 测试中止已通过合同
  - [x] 测试中止非已通过合同（应失败）

- [x] API测试（使用curl或Postman）
  - [x] POST /api/contracts - 创建合同
  - [x] GET /api/contracts - 获取列表
  - [x] GET /api/contracts/1 - 获取详情
  - [x] PUT /api/contracts/1/revoke - 撤销合同
  - [x] PUT /api/contracts/2/terminate - 中止合作

- [x] 前端联调测试
  - [x] 启动服务器
  - [x] 浏览器访问合同管理页面
  - [x] 验证合同列表正常显示
  - [x] 测试新增合同功能
    - [x] 验证可以选择学生
    - [x] 验证必填字段校验
    - [x] 验证提交成功
  - [x] 测试查看详情功能
    - [x] 验证详情显示完整
  - [x] 测试撤销功能
    - [x] 验证仅待审核状态显示撤销按钮
    - [x] 验证撤销成功
  - [x] 测试中止合作功能
    - [x] 验证仅已通过状态显示中止按钮
    - [x] 验证文件上传
    - [x] 验证中止成功
  - [x] 测试筛选功能（ID、学生ID、学生姓名、类型、状态）
  - [x] 测试分页功能

## Phase 8: 文档更新

- [x] 更新项目文档
  - [x] 更新 `REFACTORING_STATUS.md`
    - [x] 将合同管理模块标记为100%完成
    - [x] 更新API接口统计（5个接口）
    - [x] 更新完成度百分比（30% → 35%）

- [x] 更新README.md（可选）
  - [x] 添加合同管理API文档
  - [x] 列出5个接口的路径、方法、参数、响应
  - [x] 添加使用示例

- [x] 创建API文档（可选）
  - [x] 编写 `docs/api/contract.md`
  - [x] 详细说明每个接口的用法
  - [x] 提供请求示例和响应示例
  - [x] 说明错误码和异常情况

## 验收标准

- [x] 所有5个API接口实现并通过测试
- [x] 数据库表验证完成，包含正确的索引和约束
- [x] 合同列表查询返回完整的学生信息
- [x] 创建合同正确设置初始状态和发起人
- [x] 撤销合同正确验证状态并更新
- [x] 中止合作正确验证状态、处理文件上传并更新
- [x] 前端页面可以正常使用所有功能
- [x] 代码遵循DDD四层架构
- [x] 代码通过lint检查
- [x] 响应格式与原Python项目完全一致

## 注意事项

1. **路由顺序**：`/contracts/:id` 必须在 `/contracts/:id/revoke` 和 `/contracts/:id/terminate` 之前注册，避免路由冲突
2. **状态验证**：撤销和中止操作必须严格验证当前status，避免非法状态流转
3. **响应格式**：严格按照原项目格式返回，字段名使用下划线（student_name而非studentName）
4. **字段名匹配**：前端使用特定字段名，DTO的JSON标签需要匹配
5. **文件上传**：中止合作需要处理文件上传，确保文件安全保存
6. **发起人字段**：从JWT token中获取当前用户名，自动填充initiator字段
7. **字符集**：确保数据库连接使用utf8mb4字符集
8. **错误处理**：使用统一的错误处理机制（response包）
9. **日志记录**：关键操作添加日志，便于调试和追踪
10. **前端兼容**：返回消息需要与前端判断条件完全匹配
