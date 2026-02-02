# 提案：迁移审批流功能

## 变更ID
`migrate-approval-workflow`

## 背景与目的

从 ZhixinStudentSaaS (Python Flask) 迁移审批流功能到 CharonOMS (Go)。审批流系统是业务流程管理的核心组件，提供灵活的、基于模板的多层级审批流程框架，支持会签/或签等高级审批特性，并可与退费等业务流程深度集成。

## 迁移范围

### 审批流类型管理 (2接口)
- `GET /api/approval-flow-types` — 获取审批流类型列表（支持ID、名称、状态筛选）
- `PUT /api/approval-flow-types/:id/status` — 更新审批流类型状态（启用/禁用）

### 审批流模板管理 (4接口)
- `GET /api/approval-flow-templates` — 获取审批流模板列表（支持ID、类型、名称、状态筛选，含关联类型名称）
- `GET /api/approval-flow-templates/:id` — 获取模板详情（含节点列表、审批人员、抄送人员）
- `POST /api/approval-flow-templates` — 创建审批流模板（含节点配置、审批人员、抄送人员）
- `PUT /api/approval-flow-templates/:id/status` — 更新模板状态（启用时自动禁用同类型其他模板）

### 审批流实例管理 (6接口)
- `GET /api/approval-flows/initiated` — 获取当前用户发起的审批流（支持ID、类型、状态筛选）
- `GET /api/approval-flows/pending` — 获取待当前用户审批的任务列表
- `GET /api/approval-flows/completed` — 获取当前用户已处理的审批列表
- `GET /api/approval-flows/copied` — 获取抄送给当前用户的通知列表
- `GET /api/approval-flows/:id/detail` — 获取审批流详情（含完整流程节点、当前用户审批记录）
- `POST /api/approval-flows/create-from-template` — 从模板创建审批流实例

### 审批操作 (3接口)
- `PUT /api/approval-flows/:id/cancel` — 撤销审批流（仅发起者且待审批状态可撤销）
- `POST /api/approval-flows/approve` — 审批通过（会签/或签逻辑，自动流转节点）
- `POST /api/approval-flows/reject` — 审批驳回（会签/或签逻辑，更新审批流状态）

**总计：15个API接口**

## 数据库表（已存在）

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `approval_flow_type` | 审批流类型 | id, name, status, create_time, update_time |
| `approval_flow_template` | 审批流模板 | id, name, approval_flow_type_id, creator, status |
| `approval_flow_template_node` | 审批节点模板 | id, template_id, name, sort, type(0=会签,1=或签) |
| `approval_node_useraccount` | 审批节点人员配置 | id, node_id, useraccount_id |
| `approval_copy_useraccount` | 抄送人员配置 | id, approval_flow_template_id, useraccount_id |
| `approval_flow_management` | 审批流实例主记录 | id, template_id, type_id, step, create_user, status(0=待审批,10=已通过,20=已驳回,99=已撤销) |
| `approval_node_case` | 审批节点实例 | id, node_id, flow_id, type, sort, result |
| `approval_node_case_user` | 审批人员记录 | id, node_case_id, useraccount_id, result, handle_time |
| `approval_copy_useraccount_case` | 抄送记录 | id, flow_id, useraccount_id, copy_info |

**注意**：数据库表结构不可修改，必须严格按照现有表结构实现。

## 关键业务规则

### 模板管理规则
1. **单一启用原则**：同一审批流类型，仅能有一个模板处于启用状态(status=0)
2. **启用时自动禁用**：启用一个模板时，系统自动禁用同类型的其他模板
3. **节点配置要求**：每个模板至少一个节点，每个节点必须配置至少一个审批人

### 审批节点类型
- **会签节点(type=0)**：所有审批人都必须审批通过才能进入下一节点；任意一人驳回则整个审批流驳回
- **或签节点(type=1)**：任意一个审批人通过即可进入下一节点；所有人都驳回才驳回审批流

### 审批流转逻辑
1. **通过逻辑**：
   - 会签节点：所有人通过 → 节点通过 → 进入下一节点
   - 或签节点：任意人通过 → 节点通过 → 进入下一节点（自动删除同节点其他待审批记录）
   - 最后一个节点通过 → 审批流状态变为"已通过"(status=10) → 创建抄送记录

2. **驳回逻辑**：
   - 会签节点：任意人驳回 → 节点驳回 → 审批流状态变为"已驳回"(status=20)（自动删除同节点其他待审批记录）
   - 或签节点：所有人驳回 → 节点驳回 → 审批流状态变为"已驳回"(status=20)

3. **撤销逻辑**：
   - 仅发起者可撤销
   - 仅待审批(status=0)状态可撤销
   - 撤销后审批流状态变为"已撤销"(status=99)

### 状态定义
- **审批流状态**：0=待审批，10=已通过，20=已驳回，99=已撤销
- **节点结果**：NULL=审批中，0=已通过，1=已驳回
- **人员审批结果**：NULL=待审批，0=通过，1=驳回
- **模板/类型状态**：0=启用，1=禁用

## 技术实现要点

### DDD 四层架构
1. **Domain层**（internal/domain/approval/）
   - entity/：定义所有审批流相关实体
   - repository/：定义仓储接口
   - service/：实现审批流转的核心业务逻辑（会签/或签判断）

2. **Application层**（internal/application/）
   - service/：协调领域对象，实现用例
   - assembler/：DTO与Entity转换

3. **Infrastructure层**（internal/infrastructure/persistence/approval/）
   - 实现repository接口
   - 使用GORM进行数据库操作
   - 事务管理保证数据一致性

4. **Interfaces层**（internal/interfaces/http/handler/approval/）
   - HTTP Handler处理请求
   - DTO定义
   - 路由注册

### 关键技术难点
1. **会签/或签逻辑**：需要准确实现两种审批模式的判断和流转
2. **事务管理**：审批操作涉及多表更新，需要确保原子性
3. **权限控制**：基于RBAC的操作权限控制
4. **并发控制**：多个审批人同时审批时的并发处理

## API兼容性

必须与原Python项目的API保持完全兼容：
- 请求参数格式一致
- 响应JSON结构一致
- HTTP状态码一致
- 错误消息格式一致

响应格式：
```json
{
  "code": 0,
  "message": "success",
  "data": {...}
}
```

## 验收标准

- [ ] 全部15个接口实现并通过测试
- [ ] 会签/或签逻辑正确实现
- [ ] 审批流转逻辑正确（通过、驳回、撤销）
- [ ] 模板启用/禁用的互斥逻辑正确
- [ ] 事务管理确保数据一致性
- [ ] 响应格式与原Python项目完全一致
- [ ] DDD四层架构实现
- [ ] 权限控制正确（基于RBAC）
- [ ] 路由注册到Gin框架
- [ ] 单元测试覆盖率 > 70%

## 影响范围

- **新增规范**：approval（审批流管理）
- **新增代码**：
  - internal/domain/approval/
  - internal/application/service/approval/
  - internal/application/assembler/approval/
  - internal/infrastructure/persistence/approval/
  - internal/interfaces/http/handler/approval/
  - internal/interfaces/http/dto/approval/
- **路由注册**：cmd/server/main.go 或 internal/interfaces/http/router/
- **权限数据**：需要添加审批流相关权限到permissions表

## 后续扩展

审批流系统设计为通用框架，未来可以轻松扩展到其他业务场景：
- 订单取消审批
- 费用报销审批
- 合同审批
- 权限变更审批

本次迁移专注于审批流核心功能的实现，不包含与退费流程的集成（退费模块尚未迁移）。业务集成将在退费模块迁移时一并实现。
