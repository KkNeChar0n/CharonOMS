# 实施任务清单

## 1. 准备工作

- [ ] 1.1 阅读ZhixinStudentSaaS项目的审批流实现代码（app.py: 6751-8265行）
- [ ] 1.2 分析数据库表结构和关系（charonoms_backup_20260127.sql）
- [ ] 1.3 确认现有RBAC权限系统集成方式
- [ ] 1.4 创建DDD四层目录结构

## 2. Domain层实现（核心业务逻辑）

- [ ] 2.1 定义审批流实体（internal/domain/approval/entity/）
  - [ ] 2.1.1 ApprovalFlowType（审批流类型）
  - [ ] 2.1.2 ApprovalFlowTemplate（审批流模板）
  - [ ] 2.1.3 ApprovalFlowTemplateNode（审批节点模板）
  - [ ] 2.1.4 ApprovalNodeUserAccount（审批节点人员）
  - [ ] 2.1.5 ApprovalCopyUserAccount（抄送人员配置）
  - [ ] 2.1.6 ApprovalFlowManagement（审批流实例）
  - [ ] 2.1.7 ApprovalNodeCase（审批节点实例）
  - [ ] 2.1.8 ApprovalNodeCaseUser（审批人员记录）
  - [ ] 2.1.9 ApprovalCopyUserAccountCase（抄送记录）

- [ ] 2.2 定义Repository接口（internal/domain/approval/repository/）
  - [ ] 2.2.1 ApprovalFlowTypeRepository
  - [ ] 2.2.2 ApprovalFlowTemplateRepository
  - [ ] 2.2.3 ApprovalFlowManagementRepository
  - [ ] 2.2.4 ApprovalNodeCaseRepository

- [ ] 2.3 实现Domain Service（internal/domain/approval/service/）
  - [ ] 2.3.1 ApprovalFlowService - 审批流转核心逻辑
    - [ ] 会签节点判断逻辑（所有人通过/任意人驳回）
    - [ ] 或签节点判断逻辑（任意人通过/所有人驳回）
    - [ ] 节点流转逻辑（进入下一节点）
    - [ ] 审批流完成逻辑（创建抄送记录）

## 3. Infrastructure层实现（数据访问）

- [ ] 3.1 实现ApprovalFlowTypeRepository（internal/infrastructure/persistence/approval/）
  - [ ] 3.1.1 GetList - 获取审批流类型列表
  - [ ] 3.1.2 GetByID - 根据ID查询
  - [ ] 3.1.3 UpdateStatus - 更新状态

- [ ] 3.2 实现ApprovalFlowTemplateRepository
  - [ ] 3.2.1 GetList - 获取模板列表（含关联类型）
  - [ ] 3.2.2 GetByID - 获取模板详情
  - [ ] 3.2.3 GetDetailByID - 获取完整详情（含节点、人员）
  - [ ] 3.2.4 Create - 创建模板（事务：模板、节点、人员、抄送）
  - [ ] 3.2.5 UpdateStatus - 更新状态
  - [ ] 3.2.6 DisableSameTypeTemplates - 禁用同类型其他模板

- [ ] 3.3 实现ApprovalFlowManagementRepository
  - [ ] 3.3.1 GetInitiatedFlows - 获取用户发起的审批流
  - [ ] 3.3.2 GetPendingFlows - 获取待审批任务
  - [ ] 3.3.3 GetCompletedFlows - 获取已处理任务
  - [ ] 3.3.4 GetCopiedFlows - 获取抄送通知
  - [ ] 3.3.5 GetDetailByID - 获取审批流详情
  - [ ] 3.3.6 CreateFromTemplate - 从模板创建实例
  - [ ] 3.3.7 UpdateStatus - 更新审批流状态
  - [ ] 3.3.8 Cancel - 撤销审批流

- [ ] 3.4 实现ApprovalNodeCaseRepository
  - [ ] 3.4.1 GetByFlowIDAndStep - 获取当前节点实例
  - [ ] 3.4.2 CreateNextNode - 创建下一节点实例
  - [ ] 3.4.3 UpdateNodeResult - 更新节点结果
  - [ ] 3.4.4 GetNodeUsers - 获取节点审批人员
  - [ ] 3.4.5 UpdateUserResult - 更新人员审批结果
  - [ ] 3.4.6 DeletePendingUsers - 删除待审批人员（会签驳回/或签通过）
  - [ ] 3.4.7 CreateCopyRecords - 创建抄送记录

## 4. Application层实现（应用服务）

- [ ] 4.1 实现DTO定义（internal/interfaces/http/dto/approval/）
  - [ ] 4.1.1 审批流类型DTO
  - [ ] 4.1.2 审批流模板DTO（含节点、人员）
  - [ ] 4.1.3 审批流实例DTO
  - [ ] 4.1.4 审批操作DTO

- [ ] 4.2 实现Assembler（internal/application/assembler/approval/）
  - [ ] 4.2.1 ApprovalFlowTypeAssembler
  - [ ] 4.2.2 ApprovalFlowTemplateAssembler
  - [ ] 4.2.3 ApprovalFlowManagementAssembler

- [ ] 4.3 实现Application Service（internal/application/service/approval/）
  - [ ] 4.3.1 ApprovalFlowTypeService
    - [ ] GetList - 查询类型列表
    - [ ] UpdateStatus - 更新状态
  - [ ] 4.3.2 ApprovalFlowTemplateService
    - [ ] GetList - 查询模板列表
    - [ ] GetDetail - 查询模板详情
    - [ ] Create - 创建模板
    - [ ] UpdateStatus - 更新状态（含同类型互斥）
  - [ ] 4.3.3 ApprovalFlowManagementService
    - [ ] GetInitiatedFlows - 我发起的
    - [ ] GetPendingFlows - 待我审批
    - [ ] GetCompletedFlows - 处理完成
    - [ ] GetCopiedFlows - 抄送我的
    - [ ] GetDetail - 审批流详情
    - [ ] CreateFromTemplate - 创建审批流
    - [ ] Cancel - 撤销审批流
    - [ ] Approve - 审批通过（调用Domain Service）
    - [ ] Reject - 审批驳回（调用Domain Service）

## 5. Interfaces层实现（HTTP接口）

- [ ] 5.1 实现审批流类型Handler（internal/interfaces/http/handler/approval/）
  - [ ] 5.1.1 GetApprovalFlowTypes - GET /api/approval-flow-types
  - [ ] 5.1.2 UpdateApprovalFlowTypeStatus - PUT /api/approval-flow-types/:id/status

- [ ] 5.2 实现审批流模板Handler
  - [ ] 5.2.1 GetApprovalFlowTemplates - GET /api/approval-flow-templates
  - [ ] 5.2.2 GetApprovalFlowTemplateDetail - GET /api/approval-flow-templates/:id
  - [ ] 5.2.3 CreateApprovalFlowTemplate - POST /api/approval-flow-templates
  - [ ] 5.2.4 UpdateApprovalFlowTemplateStatus - PUT /api/approval-flow-templates/:id/status

- [ ] 5.3 实现审批流管理Handler
  - [ ] 5.3.1 GetInitiatedFlows - GET /api/approval-flows/initiated
  - [ ] 5.3.2 GetPendingFlows - GET /api/approval-flows/pending
  - [ ] 5.3.3 GetCompletedFlows - GET /api/approval-flows/completed
  - [ ] 5.3.4 GetCopiedFlows - GET /api/approval-flows/copied
  - [ ] 5.3.5 GetApprovalFlowDetail - GET /api/approval-flows/:id/detail
  - [ ] 5.3.6 CreateFromTemplate - POST /api/approval-flows/create-from-template

- [ ] 5.4 实现审批操作Handler
  - [ ] 5.4.1 CancelApprovalFlow - PUT /api/approval-flows/:id/cancel
  - [ ] 5.4.2 ApproveFlow - POST /api/approval-flows/approve
  - [ ] 5.4.3 RejectFlow - POST /api/approval-flows/reject

## 6. 路由和中间件配置

- [ ] 6.1 注册审批流路由（internal/interfaces/http/router/approval.go）
- [ ] 6.2 应用JWT认证中间件
- [ ] 6.3 应用RBAC权限中间件
- [ ] 6.4 集成到主路由（cmd/server/main.go）

## 7. 权限数据初始化

- [ ] 7.1 创建审批流类型管理权限
  - [ ] enable_approval_type - 启用审批流类型
  - [ ] disable_approval_type - 禁用审批流类型

- [ ] 7.2 创建审批流模板管理权限
  - [ ] view_approval_template - 查看审批流模板
  - [ ] add_approval_template - 新增审批流模板
  - [ ] enable_approval_template - 启用审批流模板
  - [ ] disable_approval_template - 禁用审批流模板

- [ ] 7.3 创建审批流实例管理权限
  - [ ] view_approval_flow - 查看审批流
  - [ ] create_approval_flow - 创建审批流
  - [ ] cancel_approval_flow - 撤销审批流
  - [ ] approve_flow - 审批通过
  - [ ] reject_flow - 审批驳回

- [ ] 7.4 创建SQL初始化脚本（scripts/init_approval_permissions.sql）

## 8. 测试

- [ ] 8.1 单元测试
  - [ ] 8.1.1 Domain Service单元测试（会签/或签逻辑）
  - [ ] 8.1.2 Application Service单元测试

- [ ] 8.2 集成测试
  - [ ] 8.2.1 审批流类型管理集成测试
  - [ ] 8.2.2 审批流模板管理集成测试
  - [ ] 8.2.3 审批流实例管理集成测试
  - [ ] 8.2.4 审批操作集成测试（通过/驳回/撤销）

- [ ] 8.3 端到端测试
  - [ ] 8.3.1 创建模板→创建审批流→多级审批→通过流程
  - [ ] 8.3.2 创建模板→创建审批流→审批驳回流程
  - [ ] 8.3.3 创建模板→创建审批流→撤销流程
  - [ ] 8.3.4 会签节点测试（多人审批）
  - [ ] 8.3.5 或签节点测试（任意一人审批）

## 9. 文档和验证

- [ ] 9.1 编写API文档
- [ ] 9.2 更新项目README
- [ ] 9.3 与原Python项目进行API兼容性测试
- [ ] 9.4 性能测试（审批流创建和流转性能）
- [ ] 9.5 运行openspec validate验证规范

## 10. 部署准备

- [ ] 10.1 确认数据库表结构（已存在）
- [ ] 10.2 执行权限数据初始化脚本
- [ ] 10.3 编译测试
- [ ] 10.4 前端联调测试

## 依赖关系说明

- 任务2（Domain层）是基础，必须先完成
- 任务3（Infrastructure层）依赖任务2的Repository接口定义
- 任务4（Application层）依赖任务2和任务3
- 任务5（Interfaces层）依赖任务4
- 任务6（路由配置）依赖任务5
- 任务8（测试）可以与开发并行，但最终验证需要等待全部实现完成
