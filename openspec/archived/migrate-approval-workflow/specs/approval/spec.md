# approval 规范增量

## ADDED Requirements

### 需求：审批流类型查询

系统MUST支持查询审批流类型列表，MUST支持按ID、名称、状态筛选。

- **标识符**: `approval-flow-type-list`
- **优先级**: P0

#### 场景：查询所有审批流类型

**前置条件**:
- 用户已登录
- 用户拥有相关权限

**操作**:
```http
GET /api/approval-flow-types
Authorization: Bearer {token}
```

**预期结果**:
- 返回所有审批流类型列表
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "approval_flow_types": [
        {
          "id": 1,
          "name": "退费审批",
          "status": 0,
          "create_time": "2026-01-30T10:00:00Z",
          "update_time": "2026-01-30T10:00:00Z"
        }
      ]
    }
  }
  ```

#### 场景：按条件筛选类型

**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/approval-flow-types?name=退费&status=0
```

**预期结果**:
- 返回名称包含"退费"且状态为启用的类型

---

### 需求：审批流类型状态更新

系统MUST支持更新审批流类型的状态（启用/禁用）。

- **标识符**: `approval-flow-type-status-update`
- **优先级**: P0

#### 场景：启用审批流类型

**前置条件**:
- 用户已登录
- 用户拥有 `enable_approval_type` 权限

**操作**:
```http
PUT /api/approval-flow-types/1/status
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": 0
}
```

**预期结果**:
- 审批流类型状态更新为启用
- 返回成功响应

#### 场景：禁用审批流类型

**前置条件**:
- 用户已登录
- 用户拥有 `disable_approval_type` 权限

**操作**:
```http
PUT /api/approval-flow-types/1/status
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": 1
}
```

**预期结果**:
- 审批流类型状态更新为禁用
- 返回成功响应

---

### 需求：审批流模板查询

系统MUST支持查询审批流模板列表，MUST支持按ID、类型ID、名称、状态筛选，MUST包含关联的审批流类型名称。

- **标识符**: `approval-flow-template-list`
- **优先级**: P0

#### 场景：查询所有模板

**前置条件**:
- 用户已登录
- 用户拥有 `view_approval_template` 权限

**操作**:
```http
GET /api/approval-flow-templates
Authorization: Bearer {token}
```

**预期结果**:
- 返回所有审批流模板列表
- 包含关联的审批流类型名称
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "approval_flow_templates": [
        {
          "id": 1,
          "name": "标准退费审批流程",
          "approval_flow_type_id": 1,
          "flow_type_name": "退费审批",
          "creator": "admin",
          "status": 0,
          "create_time": "2026-01-30T10:00:00Z",
          "update_time": "2026-01-30T10:00:00Z"
        }
      ]
    }
  }
  ```

#### 场景：按类型筛选模板

**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/approval-flow-templates?approval_flow_type_id=1
```

**预期结果**:
- 仅返回指定类型的模板

---

### 需求：审批流模板详情查询

系统MUST支持查询审批流模板的完整详情，MUST包含节点列表、每个节点的审批人员、抄送人员配置。

- **标识符**: `approval-flow-template-detail`
- **优先级**: P0

#### 场景：获取模板完整详情

**前置条件**:
- 用户已登录
- 用户拥有 `view_approval_template` 权限

**操作**:
```http
GET /api/approval-flow-templates/1
Authorization: Bearer {token}
```

**预期结果**:
- 返回模板完整信息
- 包含按sort排序的节点列表
- 每个节点包含审批人员列表
- 包含抄送人员列表
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "template": {
        "id": 1,
        "name": "标准退费审批流程",
        "approval_flow_type_id": 1,
        "flow_type_name": "退费审批",
        "creator": "admin",
        "status": 0,
        "create_time": "2026-01-30T10:00:00Z",
        "update_time": "2026-01-30T10:00:00Z",
        "nodes": [
          {
            "id": 1,
            "name": "财务审批",
            "sort": 0,
            "type": 0,
            "approvers": [
              {
                "id": 2,
                "username": "finance_user"
              }
            ]
          },
          {
            "id": 2,
            "name": "总经理审批",
            "sort": 1,
            "type": 1,
            "approvers": [
              {
                "id": 3,
                "username": "manager1"
              },
              {
                "id": 4,
                "username": "manager2"
              }
            ]
          }
        ],
        "copy_users": [
          {
            "id": 5,
            "username": "hr_user"
          }
        ]
      }
    }
  }
  ```

---

### 需求：审批流模板创建

系统MUST支持创建审批流模板，MUST包含节点配置、审批人员、抄送人员，MUST验证节点和审批人员的完整性。

- **标识符**: `approval-flow-template-create`
- **优先级**: P0

#### 场景：创建审批流模板

**前置条件**:
- 用户已登录
- 用户拥有 `add_approval_template` 权限

**操作**:
```http
POST /api/approval-flow-templates
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "标准退费审批流程",
  "approval_flow_type_id": 1,
  "nodes": [
    {
      "name": "财务审批",
      "type": 0,
      "approvers": [2]
    },
    {
      "name": "总经理审批",
      "type": 1,
      "approvers": [3, 4]
    }
  ],
  "copy_users": [5]
}
```

**预期结果**:
- 模板创建成功
- 所有节点和审批人员关联创建成功
- 抄送人员配置创建成功
- 返回新创建的模板ID
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1
    }
  }
  ```

#### 场景：创建时验证节点完整性

**前置条件**:
- 用户已登录
- 用户拥有 `add_approval_template` 权限

**操作**:
```http
POST /api/approval-flow-templates
Content-Type: application/json

{
  "name": "测试模板",
  "approval_flow_type_id": 1,
  "nodes": []
}
```

**预期结果**:
- 返回错误响应
- 错误信息：模板至少需要一个节点

#### 场景：创建时验证审批人员完整性

**前置条件**:
- 用户已登录
- 用户拥有 `add_approval_template` 权限

**操作**:
```http
POST /api/approval-flow-templates
Content-Type: application/json

{
  "name": "测试模板",
  "approval_flow_type_id": 1,
  "nodes": [
    {
      "name": "审批节点",
      "type": 0,
      "approvers": []
    }
  ]
}
```

**预期结果**:
- 返回错误响应
- 错误信息：每个节点至少需要一个审批人

---

### 需求：审批流模板状态更新

系统MUST支持更新审批流模板状态，启用模板时MUST自动禁用同类型的其他模板，确保同一类型仅有一个启用的模板。

- **标识符**: `approval-flow-template-status-update`
- **优先级**: P0

#### 场景：启用模板（同类型互斥）

**前置条件**:
- 用户已登录
- 用户拥有 `enable_approval_template` 权限
- 同类型存在其他启用的模板

**操作**:
```http
PUT /api/approval-flow-templates/2/status
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": 0
}
```

**预期结果**:
- 模板2状态更新为启用
- 同类型其他启用的模板自动禁用
- 返回成功响应

#### 场景：禁用模板

**前置条件**:
- 用户已登录
- 用户拥有 `disable_approval_template` 权限

**操作**:
```http
PUT /api/approval-flow-templates/1/status
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": 1
}
```

**预期结果**:
- 模板状态更新为禁用
- 不影响其他模板
- 返回成功响应

---

### 需求：查询我发起的审批流

系统MUST支持查询当前用户发起的审批流列表，MUST支持按ID、审批流类型、状态筛选。

- **标识符**: `approval-flow-initiated-list`
- **优先级**: P0

#### 场景：查询我发起的所有审批流

**前置条件**:
- 用户已登录
- 用户拥有 `view_approval_flow` 权限

**操作**:
```http
GET /api/approval-flows/initiated
Authorization: Bearer {token}
```

**预期结果**:
- 返回当前用户发起的所有审批流
- 包含审批流类型名称
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "flows": [
        {
          "id": 1,
          "approval_flow_template_id": 1,
          "approval_flow_type_id": 1,
          "flow_type_name": "退费审批",
          "step": 0,
          "create_user": 1,
          "create_time": "2026-01-30T10:00:00Z",
          "status": 0,
          "complete_time": null
        }
      ]
    }
  }
  ```

#### 场景：按状态筛选我发起的审批流

**前置条件**:
- 用户已登录

**操作**:
```http
GET /api/approval-flows/initiated?status=10
```

**预期结果**:
- 仅返回状态为"已通过"的审批流

---

### 需求：查询待我审批的任务

系统MUST支持查询待当前用户审批的任务列表，MUST支持按审批流ID、审批流类型筛选。

- **标识符**: `approval-flow-pending-list`
- **优先级**: P0

#### 场景：查询待我审批的所有任务

**前置条件**:
- 用户已登录
- 用户拥有 `view_approval_flow` 权限

**操作**:
```http
GET /api/approval-flows/pending
Authorization: Bearer {token}
```

**预期结果**:
- 返回待当前用户审批且审批流状态为待审批的任务
- 包含审批流类型、发起人信息、节点信息
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "tasks": [
        {
          "id": 1,
          "approval_flow_id": 1,
          "approval_flow_type_id": 1,
          "flow_type_name": "退费审批",
          "create_user_name": "张三",
          "create_time": "2026-01-30T10:00:00Z",
          "node_type": 0,
          "node_sort": 0,
          "node_case_user_id": 1
        }
      ]
    }
  }
  ```

---

### 需求：查询处理完成的审批

系统MUST支持查询当前用户已处理的审批列表，MUST包含审批结果和处理时间。

- **标识符**: `approval-flow-completed-list`
- **优先级**: P0

#### 场景：查询处理完成的所有审批

**前置条件**:
- 用户已登录
- 用户拥有 `view_approval_flow` 权限

**操作**:
```http
GET /api/approval-flows/completed
Authorization: Bearer {token}
```

**预期结果**:
- 返回当前用户已处理的审批任务
- 包含审批结果（通过/驳回）和处理时间
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "tasks": [
        {
          "id": 1,
          "approval_flow_id": 1,
          "approval_flow_type_id": 1,
          "flow_type_name": "退费审批",
          "create_user_name": "张三",
          "create_time": "2026-01-30T10:00:00Z",
          "node_type": 0,
          "node_sort": 0,
          "result": 0,
          "handle_time": "2026-01-30T11:00:00Z"
        }
      ]
    }
  }
  ```

---

### 需求：查询抄送我的通知

系统MUST支持查询抄送给当前用户的审批通知列表。

- **标识符**: `approval-flow-copied-list`
- **优先级**: P0

#### 场景：查询抄送我的所有通知

**前置条件**:
- 用户已登录
- 用户拥有 `view_approval_flow` 权限

**操作**:
```http
GET /api/approval-flows/copied
Authorization: Bearer {token}
```

**预期结果**:
- 返回抄送给当前用户的通知列表
- 包含抄送信息
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "notifications": [
        {
          "id": 1,
          "approval_flow_id": 1,
          "approval_flow_type_id": 1,
          "flow_type_name": "退费审批",
          "create_user_name": "张三",
          "create_time": "2026-01-30T10:00:00Z",
          "complete_time": "2026-01-30T12:00:00Z",
          "copy_info": "学生李四的退费申请已通过"
        }
      ]
    }
  }
  ```

---

### 需求：查询审批流详情

系统MUST支持查询审批流的完整详情，MUST包含所有节点信息（用于展示审批流程数轴）、当前用户的审批记录。

- **标识符**: `approval-flow-detail`
- **优先级**: P0

#### 场景：获取审批流完整详情

**前置条件**:
- 用户已登录
- 用户拥有 `view_approval_flow` 权限

**操作**:
```http
GET /api/approval-flows/1/detail
Authorization: Bearer {token}
```

**预期结果**:
- 返回审批流完整信息
- 包含所有节点信息（按sort排序，显示节点状态）
- 包含当前用户在各节点的审批记录
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "flow": {
        "id": 1,
        "approval_flow_type_id": 1,
        "flow_type_name": "退费审批",
        "create_user": 1,
        "create_user_name": "张三",
        "create_time": "2026-01-30T10:00:00Z",
        "status": 0,
        "step": 0
      },
      "nodes": [
        {
          "node_id": 1,
          "name": "财务审批",
          "sort": 0,
          "type": 0,
          "result": null,
          "complete_time": null,
          "is_current_user_node": true
        },
        {
          "node_id": 2,
          "name": "总经理审批",
          "sort": 1,
          "type": 1,
          "result": null,
          "complete_time": null,
          "is_current_user_node": false
        }
      ],
      "user_approval_record": {
        "node_case_user_id": 1,
        "result": null,
        "handle_time": null
      }
    }
  }
  ```

---

### 需求：从模板创建审批流实例

系统MUST支持从启用的模板创建审批流实例，MUST创建第一个节点实例和对应的审批人员记录。

- **标识符**: `approval-flow-create-from-template`
- **优先级**: P0

#### 场景：成功创建审批流实例

**前置条件**:
- 用户已登录
- 用户拥有 `create_approval_flow` 权限
- 模板存在且启用

**操作**:
```http
POST /api/approval-flows/create-from-template
Authorization: Bearer {token}
Content-Type: application/json

{
  "template_id": 1
}
```

**预期结果**:
- 审批流实例创建成功
- 第一个节点实例创建成功
- 第一个节点的所有审批人员记录创建成功
- 返回新创建的审批流ID
- 响应格式：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "flow_id": 1
    }
  }
  ```

#### 场景：模板不存在或已禁用

**前置条件**:
- 用户已登录
- 模板不存在或状态为禁用

**操作**:
```http
POST /api/approval-flows/create-from-template
Content-Type: application/json

{
  "template_id": 999
}
```

**预期结果**:
- 返回错误响应
- 错误信息：模板不存在或已禁用

---

### 需求：撤销审批流

系统MUST支持撤销审批流，MUST验证仅发起者可撤销且仅待审批状态可撤销。

- **标识符**: `approval-flow-cancel`
- **优先级**: P0

#### 场景：成功撤销审批流

**前置条件**:
- 用户已登录
- 用户拥有 `cancel_approval_flow` 权限
- 用户是审批流的发起者
- 审批流状态为待审批

**操作**:
```http
PUT /api/approval-flows/1/cancel
Authorization: Bearer {token}
```

**预期结果**:
- 审批流状态更新为已撤销(99)
- 记录完成时间
- 返回成功响应

#### 场景：非发起者无法撤销

**前置条件**:
- 用户已登录
- 用户不是审批流的发起者

**操作**:
```http
PUT /api/approval-flows/1/cancel
Authorization: Bearer {token}
```

**预期结果**:
- 返回错误响应
- 错误信息：仅发起者可以撤销审批流

#### 场景：非待审批状态无法撤销

**前置条件**:
- 用户已登录
- 用户是审批流的发起者
- 审批流状态不是待审批

**操作**:
```http
PUT /api/approval-flows/1/cancel
Authorization: Bearer {token}
```

**预期结果**:
- 返回错误响应
- 错误信息：仅待审批状态可以撤销

---

### 需求：审批通过

系统MUST支持审批通过操作，MUST根据节点类型（会签/或签）正确处理审批逻辑，MUST自动流转到下一节点或完成审批流。

- **标识符**: `approval-flow-approve`
- **优先级**: P0

#### 场景：会签节点所有人通过后流转

**前置条件**:
- 用户已登录
- 用户拥有 `approve_flow` 权限
- 当前节点为会签节点(type=0)
- 当前用户是最后一个审批人

**操作**:
```http
POST /api/approval-flows/approve
Authorization: Bearer {token}
Content-Type: application/json

{
  "node_case_user_id": 1
}
```

**预期结果**:
- 当前用户审批结果更新为通过(0)
- 记录处理时间
- 节点结果更新为通过(0)
- 如果有下一节点：创建下一节点实例和审批人员记录，step+1
- 如果没有下一节点：审批流状态更新为已通过(10)，创建抄送记录
- 返回成功响应

#### 场景：或签节点任意人通过后流转

**前置条件**:
- 用户已登录
- 用户拥有 `approve_flow` 权限
- 当前节点为或签节点(type=1)

**操作**:
```http
POST /api/approval-flows/approve
Authorization: Bearer {token}
Content-Type: application/json

{
  "node_case_user_id": 2
}
```

**预期结果**:
- 当前用户审批结果更新为通过(0)
- 记录处理时间
- 删除同节点其他待审批人员记录
- 节点结果更新为通过(0)
- 如果有下一节点：创建下一节点实例和审批人员记录，step+1
- 如果没有下一节点：审批流状态更新为已通过(10)，创建抄送记录
- 返回成功响应

---

### 需求：审批驳回

系统MUST支持审批驳回操作，MUST根据节点类型（会签/或签）正确处理驳回逻辑，MUST更新审批流状态为已驳回。

- **标识符**: `approval-flow-reject`
- **优先级**: P0

#### 场景：会签节点任意人驳回

**前置条件**:
- 用户已登录
- 用户拥有 `reject_flow` 权限
- 当前节点为会签节点(type=0)

**操作**:
```http
POST /api/approval-flows/reject
Authorization: Bearer {token}
Content-Type: application/json

{
  "node_case_user_id": 1
}
```

**预期结果**:
- 当前用户审批结果更新为驳回(1)
- 记录处理时间
- 删除同节点其他待审批人员记录
- 节点结果更新为驳回(1)
- 审批流状态更新为已驳回(20)
- 记录审批流完成时间
- 返回成功响应

#### 场景：或签节点所有人驳回

**前置条件**:
- 用户已登录
- 用户拥有 `reject_flow` 权限
- 当前节点为或签节点(type=1)
- 当前用户是最后一个审批人

**操作**:
```http
POST /api/approval-flows/reject
Authorization: Bearer {token}
Content-Type: application/json

{
  "node_case_user_id": 2
}
```

**预期结果**:
- 当前用户审批结果更新为驳回(1)
- 记录处理时间
- 节点结果更新为驳回(1)
- 审批流状态更新为已驳回(20)
- 记录审批流完成时间
- 返回成功响应

---

### 需求：审批流权限控制

系统MUST基于RBAC权限系统控制审批流的操作权限。

- **标识符**: `approval-flow-permission-control`
- **优先级**: P0

#### 场景：超级管理员拥有所有权限

**前置条件**:
- 用户已登录
- 用户是超级管理员(is_super_admin=1)

**操作**:
- 执行任意审批流操作

**预期结果**:
- 所有操作均通过权限验证

#### 场景：普通用户需要对应权限

**前置条件**:
- 用户已登录
- 用户不是超级管理员
- 用户角色未授予对应权限

**操作**:
```http
POST /api/approval-flow-templates
Authorization: Bearer {token}
```

**预期结果**:
- 返回权限不足错误
- HTTP状态码：403

---

### 需求：审批流数据完整性

系统MUST通过数据库事务保证审批流相关操作的数据完整性。

- **标识符**: `approval-flow-data-integrity`
- **优先级**: P0

#### 场景：创建模板事务完整性

**前置条件**:
- 用户已登录
- 创建模板时节点或人员数据写入失败

**操作**:
- 创建审批流模板

**预期结果**:
- 整个创建操作回滚
- 模板、节点、人员、抄送记录均不创建
- 返回错误响应

#### 场景：审批操作事务完整性

**前置条件**:
- 用户已登录
- 审批操作过程中任一步骤失败

**操作**:
- 执行审批通过或驳回

**预期结果**:
- 整个审批操作回滚
- 相关状态均不更新
- 返回错误响应
