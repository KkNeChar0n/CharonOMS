# RBAC规范增量变更

## MODIFIED Requirements

### Requirement: 菜单路由命名规范

系统 SHALL 使用snake_case格式命名所有菜单路由，与原Python Flask项目保持完全一致。

**修改原因**：统一前后端路由命名规范，确保前端页面能够正确跳转到对应的后端API路由。

#### Scenario: 菜单管理路由正确命名
- **GIVEN** 系统有菜单管理功能
- **WHEN** 用户访问菜单管理页面
- **THEN** 前端路由为 `/#/menu_management`
- **AND** 后端API路由为 `/api/menu_management`
- **AND** 路由使用snake_case格式（menu_management）而非kebab-case（menu-management）

#### Scenario: 活动模板路由正确命名
- **GIVEN** 系统有活动模板功能
- **WHEN** 用户访问活动模板页面
- **THEN** 前端路由为 `/#/activity_template`
- **AND** 后端API路由为 `/api/activity_template`
- **AND** 路由使用snake_case格式

#### Scenario: 所有路由遵循snake_case规范
- **GIVEN** 系统有多个功能模块
- **WHEN** 定义新的菜单路由
- **THEN** 路由名称 MUST 使用snake_case格式
- **AND** 禁止使用kebab-case、camelCase或PascalCase格式

### Requirement: 完整的财务管理菜单

系统 SHALL 提供完整的财务管理菜单，包括退款和分账功能，与原项目保持一致。

#### Scenario: 退款子订单菜单
- **GIVEN** 用户有财务管理权限
- **WHEN** 查看财务管理菜单
- **THEN** 能看到"退款子订单"菜单项
- **AND** 路由为 `refund_childorders`
- **AND** 点击后能查看退款订单明细

#### Scenario: 退款管理菜单
- **GIVEN** 用户有财务管理权限
- **WHEN** 查看财务管理菜单
- **THEN** 能看到"退款管理"菜单项
- **AND** 路由为 `refund_management`
- **AND** 点击后能管理常规退费和淘宝退费

#### Scenario: 退款支付详情菜单
- **GIVEN** 用户有财务管理权限
- **WHEN** 查看财务管理菜单
- **THEN** 能看到"退款支付详情"菜单项
- **AND** 路由为 `refund_payment_detail`
- **AND** 点击后能查看退款支付明细

#### Scenario: 分账明细菜单
- **GIVEN** 用户有财务管理权限
- **WHEN** 查看财务管理菜单
- **THEN** 能看到"分账明细"菜单项
- **AND** 路由为 `separate_account`
- **AND** 点击后能查看分账管理信息

### Requirement: 审批流程管理结构

系统 SHALL 采用审批流程管理模式，提供审批流类型、模板和管理三个核心功能。

**修改原因**：将简单的三层审批（我发起的/待我审批/已审批）改为更灵活的审批流程管理设计。

#### Scenario: 审批流类型管理
- **GIVEN** 用户有审批管理权限
- **WHEN** 访问审批管理菜单
- **THEN** 能看到"审批流类型"菜单项
- **AND** 路由为 `approval_flow_type`
- **AND** 能够定义和管理不同类型的审批流程（如请假审批、报销审批等）

#### Scenario: 审批流模板配置
- **GIVEN** 用户有审批管理权限
- **WHEN** 访问审批管理菜单
- **THEN** 能看到"审批流模板"菜单项
- **AND** 路由为 `approval_flow_template`
- **AND** 能够配置审批流程的模板（定义审批节点、审批人等）

#### Scenario: 审批流程统一管理
- **GIVEN** 用户有审批管理权限
- **WHEN** 访问审批流管理页面
- **THEN** 路由为 `approval_flow_management`
- **AND** 能查看我发起的审批流程
- **AND** 能查看待我审批的流程
- **AND** 能查看已审批的流程历史

### Requirement: 菜单中文名称统一

系统 SHALL 统一使用"管理"后缀命名业务功能菜单，而不是"列表"。

#### Scenario: 学生功能命名
- **GIVEN** 系统有学生管理功能
- **WHEN** 查看学生管理菜单
- **THEN** 菜单名称为"学生管理"而非"学生列表"

#### Scenario: 教练功能命名
- **GIVEN** 系统有教练管理功能
- **WHEN** 查看教练管理菜单
- **THEN** 菜单名称为"教练管理"而非"教练列表"

#### Scenario: 订单功能命名
- **GIVEN** 系统有订单管理功能
- **WHEN** 查看订单管理菜单
- **THEN** 菜单名称为"订单管理"而非"订单列表"

#### Scenario: 活动功能命名
- **GIVEN** 系统有活动管理功能
- **WHEN** 查看活动管理菜单
- **THEN** 菜单名称为"活动管理"而非"活动列表"

### Requirement: 菜单权限检查机制

系统 SHALL 根据用户角色和权限动态过滤菜单显示。

#### Scenario: 超级管理员查看所有菜单
- **GIVEN** 用户是超级管理员（is_super=1）
- **WHEN** 获取用户菜单树
- **THEN** 返回所有状态为启用的菜单
- **AND** 共8个一级菜单
- **AND** 共24个二级菜单

#### Scenario: 普通用户根据权限查看菜单
- **GIVEN** 用户是普通角色（is_super=0）
- **AND** 用户角色拥有部分菜单权限
- **WHEN** 获取用户菜单树
- **THEN** 只返回用户有权限的菜单
- **AND** 无权限的菜单不出现在菜单树中

#### Scenario: 权限变更后菜单立即更新
- **GIVEN** 用户原本有某个菜单权限
- **WHEN** 管理员移除该用户角色的菜单权限
- **THEN** 用户下次登录时看不到该菜单
- **AND** 直接访问该菜单路由返回403 Forbidden

## ADDED Requirements

### Requirement: 菜单ID编号规则

系统 SHALL 按功能模块划分菜单ID编号区间，便于管理和扩展。

#### Scenario: 系统管理模块ID区间
- **GIVEN** 系统管理模块有多个功能
- **WHEN** 分配菜单ID
- **THEN** 用户管理ID为1
- **AND** 角色管理ID为4
- **AND** 菜单管理ID为13
- **AND** 权限管理ID为14
- **AND** 所有系统管理功能ID在1-14区间内

#### Scenario: 商品管理模块ID区间
- **GIVEN** 商品管理模块有多个功能
- **WHEN** 分配菜单ID
- **THEN** 商品类别ID为15
- **AND** 商品管理ID为16
- **AND** 套餐管理ID为17
- **AND** 课程管理ID为18
- **AND** 班级管理ID为19
- **AND** 所有商品管理功能ID在15-19区间内

#### Scenario: 财务管理模块ID区间
- **GIVEN** 财务管理模块有多个功能
- **WHEN** 分配菜单ID
- **THEN** 合同管理ID为24
- **AND** 收款管理ID为25
- **AND** 退款相关功能ID在26-28区间
- **AND** 分账明细ID为39

#### Scenario: 审批管理模块ID区间
- **GIVEN** 审批管理模块有多个功能
- **WHEN** 分配菜单ID
- **THEN** 审批流类型ID为30
- **AND** 审批流模板ID为31
- **AND** 审批流管理ID为32
- **AND** 所有审批管理功能ID在30-34区间内
