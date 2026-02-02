# 审批流系统设计文档

## 上下文

审批流系统是从Python Flask项目(ZhixinStudentSaaS)迁移到Go语言的CharonOMS项目。这是一个通用的、基于模板的多层级审批框架，需要支持会签和或签两种审批模式，并能与业务流程（如退费）深度集成。

### 背景
- 原项目采用Python Flask + MySQL实现
- 数据库表结构已存在且不可修改
- 需要保持API完全兼容以支持前端无缝迁移
- 目标项目采用DDD架构和Go语言

### 约束
1. **数据库约束**：表结构已固定，不可修改
2. **API约束**：必须保持与原Python项目的API完全兼容
3. **架构约束**：必须遵循DDD四层架构（Domain、Application、Infrastructure、Interfaces）
4. **权限约束**：必须集成现有RBAC权限系统

### 利益相关者
- 前端开发团队：依赖API兼容性
- 业务团队：依赖审批流程的正确性
- 运维团队：关注系统性能和稳定性

## 目标与非目标

### 目标
1. **完整迁移**：实现原Python项目的全部15个审批流API
2. **业务正确性**：准确实现会签/或签逻辑，确保审批流转的正确性
3. **数据一致性**：通过事务管理保证多表操作的原子性
4. **可扩展性**：设计通用框架，便于未来扩展到其他审批场景
5. **性能优化**：利用Go的并发特性提升性能

### 非目标
1. **UI迁移**：前端保持原有Vue.js实现
2. **退费集成**：本次不包含与退费模块的集成（退费模块尚未迁移）
3. **表结构优化**：不修改现有数据库表结构
4. **实时通知**：不实现WebSocket等实时通知功能（保持原有轮询方式）

## 技术决策

### 决策1：DDD架构分层

**选择**：采用DDD四层架构

**理由**：
- 符合项目整体架构规范
- 清晰的职责分离，便于维护和测试
- 核心业务逻辑（会签/或签）封装在Domain层，与基础设施解耦

**分层职责**：
```
┌─────────────────────────────────────────┐
│  Interfaces Layer (HTTP Handler/DTO)   │  ← 处理HTTP请求/响应
├─────────────────────────────────────────┤
│  Application Layer (Service/Assembler) │  ← 编排业务用例，DTO转换
├─────────────────────────────────────────┤
│  Domain Layer (Entity/Service/Repo)    │  ← 核心业务逻辑
├─────────────────────────────────────────┤
│  Infrastructure Layer (Persistence)     │  ← 数据库访问实现
└─────────────────────────────────────────┘
```

**替代方案**：
- **MVC三层架构**：不考虑，不符合项目规范
- **微服务架构**：过度设计，当前系统规模不需要

### 决策2：会签/或签逻辑实现

**选择**：在Domain Service中实现审批流转逻辑

**实现方案**：

**会签节点(type=0)**：
```go
// 判断节点是否通过
func (s *ApprovalFlowService) IsCountersignNodePassed(nodeCase *ApprovalNodeCase) bool {
    // 所有人都审批 && 所有人都通过 → 节点通过
    allApproved := true
    allPassed := true

    for _, user := range nodeCase.Users {
        if user.Result == nil {
            allApproved = false
            break
        }
        if *user.Result != 0 {
            allPassed = false
        }
    }

    return allApproved && allPassed
}

// 判断节点是否驳回
func (s *ApprovalFlowService) IsCountersignNodeRejected(nodeCase *ApprovalNodeCase) bool {
    // 任意一人驳回 → 节点驳回
    for _, user := range nodeCase.Users {
        if user.Result != nil && *user.Result == 1 {
            return true
        }
    }
    return false
}
```

**或签节点(type=1)**：
```go
// 判断节点是否通过
func (s *ApprovalFlowService) IsOrSignNodePassed(nodeCase *ApprovalNodeCase) bool {
    // 任意一人通过 → 节点通过
    for _, user := range nodeCase.Users {
        if user.Result != nil && *user.Result == 0 {
            return true
        }
    }
    return false
}

// 判断节点是否驳回
func (s *ApprovalFlowService) IsOrSignNodeRejected(nodeCase *ApprovalNodeCase) bool {
    // 所有人都审批 && 所有人都驳回 → 节点驳回
    allApproved := true
    allRejected := true

    for _, user := range nodeCase.Users {
        if user.Result == nil {
            allApproved = false
            break
        }
        if *user.Result != 1 {
            allRejected = false
        }
    }

    return allApproved && allRejected
}
```

**理由**：
- 业务逻辑清晰，易于理解和测试
- 封装在Domain Service中，符合DDD设计原则
- 便于单元测试（不依赖数据库）

**替代方案**：
- **在Repository中实现**：不符合DDD职责划分，Repository只负责数据访问
- **在Application Service中实现**：会导致业务逻辑分散

### 决策3：事务管理策略

**选择**：在Repository层使用GORM事务管理

**实现方案**：
```go
// 示例：创建模板的事务管理
func (r *ApprovalFlowTemplateRepository) Create(template *entity.ApprovalFlowTemplate) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建模板
        if err := tx.Create(&template).Error; err != nil {
            return err
        }

        // 2. 创建节点
        for _, node := range template.Nodes {
            node.TemplateID = template.ID
            if err := tx.Create(&node).Error; err != nil {
                return err
            }

            // 3. 创建节点审批人员
            for _, approver := range node.Approvers {
                if err := tx.Create(&approver).Error; err != nil {
                    return err
                }
            }
        }

        // 4. 创建抄送人员
        for _, copyUser := range template.CopyUsers {
            if err := tx.Create(&copyUser).Error; err != nil {
                return err
            }
        }

        return nil
    })
}
```

**理由**：
- GORM提供了简洁的事务API
- 自动回滚机制保证数据一致性
- 符合Go语言习惯用法

**关键事务场景**：
1. **创建模板**：模板 + 节点 + 审批人员 + 抄送人员
2. **创建审批流实例**：审批流 + 第一个节点实例 + 审批人员记录
3. **审批操作**：更新审批结果 + 更新节点状态 + 流转/完成操作

### 决策4：并发控制

**选择**：依赖MySQL的行锁机制 + 乐观并发控制

**实现方案**：
```go
// 使用FOR UPDATE锁定待更新的记录
func (r *ApprovalNodeCaseUserRepository) UpdateResult(id int, result int) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        var user entity.ApprovalNodeCaseUser

        // 锁定记录
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&user, id).Error; err != nil {
            return err
        }

        // 检查状态
        if user.Result != nil {
            return errors.New("该审批已处理")
        }

        // 更新结果
        now := time.Now()
        return tx.Model(&user).Updates(map[string]interface{}{
            "result":      result,
            "handle_time": now,
        }).Error
    })
}
```

**理由**：
- 防止同一审批人重复提交
- 利用数据库的ACID特性保证一致性
- 性能开销小

**替代方案**：
- **分布式锁（Redis）**：过度设计，当前场景不需要
- **无锁设计**：风险高，可能导致数据不一致

### 决策5：权限集成方式

**选择**：使用现有的RBAC中间件

**实现方案**：
```go
// 路由注册时应用权限中间件
approvalGroup := r.Group("/api/approval-flow-templates")
approvalGroup.Use(middleware.JWTAuth(), middleware.PermissionCheck())
{
    approvalGroup.GET("", handler.GetApprovalFlowTemplates)
    approvalGroup.POST("", handler.CreateApprovalFlowTemplate)
    // ...
}
```

**权限定义**：
```sql
-- 审批流类型管理
INSERT INTO permissions (name, comment) VALUES
('enable_approval_type', '启用审批流类型'),
('disable_approval_type', '禁用审批流类型');

-- 审批流模板管理
INSERT INTO permissions (name, comment) VALUES
('view_approval_template', '查看审批流模板'),
('add_approval_template', '新增审批流模板'),
('enable_approval_template', '启用审批流模板'),
('disable_approval_template', '禁用审批流模板');

-- 审批流实例管理
INSERT INTO permissions (name, comment) VALUES
('view_approval_flow', '查看审批流'),
('create_approval_flow', '创建审批流'),
('cancel_approval_flow', '撤销审批流'),
('approve_flow', '审批通过'),
('reject_flow', '审批驳回');
```

**理由**：
- 复用现有权限系统，保持一致性
- 细粒度权限控制，符合最小权限原则
- 超级管理员自动拥有所有权限

### 决策6：DTO设计策略

**选择**：为每个API响应定义独立的DTO

**实现方案**：
```go
// 模板列表响应DTO
type ApprovalFlowTemplateListResponse struct {
    ID                  int       `json:"id"`
    Name                string    `json:"name"`
    ApprovalFlowTypeID  int       `json:"approval_flow_type_id"`
    FlowTypeName        string    `json:"flow_type_name"`
    Creator             string    `json:"creator"`
    Status              int       `json:"status"`
    CreateTime          time.Time `json:"create_time"`
    UpdateTime          time.Time `json:"update_time"`
}

// 模板详情响应DTO
type ApprovalFlowTemplateDetailResponse struct {
    ID                  int                         `json:"id"`
    Name                string                      `json:"name"`
    ApprovalFlowTypeID  int                         `json:"approval_flow_type_id"`
    FlowTypeName        string                      `json:"flow_type_name"`
    Creator             string                      `json:"creator"`
    Status              int                         `json:"status"`
    CreateTime          time.Time                   `json:"create_time"`
    UpdateTime          time.Time                   `json:"update_time"`
    Nodes               []ApprovalFlowTemplateNode  `json:"nodes"`
    CopyUsers           []UserAccountSimple         `json:"copy_users"`
}
```

**理由**：
- 与原Python项目的响应格式完全一致
- 避免Entity直接暴露给外部
- 便于独立调整响应格式

## 数据模型设计

### ER关系图

```
approval_flow_type (审批流类型)
    ↓ 1:N
approval_flow_template (审批流模板)
    ↓ 1:N
approval_flow_template_node (审批节点模板)
    ↓ N:M (通过approval_node_useraccount)
useraccount (用户账号)

approval_flow_template → approval_copy_useraccount → useraccount (抄送配置)

approval_flow_template
    ↓ 1:N
approval_flow_management (审批流实例)
    ↓ 1:N
approval_node_case (审批节点实例)
    ↓ N:M (通过approval_node_case_user)
useraccount (审批人员)

approval_flow_management → approval_copy_useraccount_case → useraccount (抄送记录)
```

### 关键字段说明

**状态字段统一定义**：
- 审批流类型/模板：0=启用，1=禁用
- 审批流实例：0=待审批，10=已通过，20=已驳回，99=已撤销
- 节点类型：0=会签，1=或签
- 审批结果：NULL=未处理，0=通过，1=驳回

## 风险与权衡

### 风险1：并发审批的数据一致性

**风险描述**：多个审批人同时审批同一节点时，可能导致数据不一致。

**缓解措施**：
- 使用数据库行锁（FOR UPDATE）
- 在更新前检查记录状态
- 使用事务保证原子性

**监控指标**：
- 审批操作失败率
- 并发冲突次数

### 风险2：会签/或签逻辑的复杂性

**风险描述**：会签和或签的判断逻辑较复杂，可能出现边界情况处理不当。

**缓解措施**：
- 在Domain Service中集中实现逻辑
- 编写充分的单元测试覆盖所有场景
- 参考原Python项目的实现和测试用例

**测试覆盖**：
- 会签节点：所有人通过、部分人通过、任意人驳回
- 或签节点：任意人通过、所有人驳回、部分人驳回

### 风险3：API兼容性破坏

**风险描述**：Go和Python的JSON序列化可能有差异，导致前端调用失败。

**缓解措施**：
- 严格定义DTO的JSON标签
- 对比原项目的响应格式
- 进行前后端联调测试

**验证方法**：
- 使用Postman对比响应格式
- 时间字段统一使用ISO 8601格式
- NULL值处理保持一致

### 风险4：性能下降

**风险描述**：从Python迁移到Go后，性能可能不如预期。

**缓解措施**：
- 优化数据库查询（使用索引、避免N+1查询）
- 使用连接池管理数据库连接
- 对关键操作进行性能测试

**性能目标**：
- 审批流创建：< 200ms
- 审批操作：< 100ms
- 列表查询：< 50ms

## 迁移计划

### 阶段1：基础设施准备（1天）
- 创建目录结构
- 定义Entity和Repository接口
- 准备测试环境

### 阶段2：数据访问层实现（2-3天）
- 实现所有Repository
- 编写Repository集成测试
- 验证数据访问正确性

### 阶段3：业务逻辑实现（2天）
- 实现Domain Service（会签/或签逻辑）
- 实现Application Service
- 编写单元测试

### 阶段4：接口层实现（2天）
- 实现HTTP Handler
- 实现DTO和Assembler
- 注册路由和中间件

### 阶段5：集成测试（1-2天）
- 端到端测试
- API兼容性测试
- 性能测试

### 阶段6：部署和验证（1天）
- 部署到测试环境
- 前端联调
- 数据验证

### 回滚计划

如果迁移失败，回滚步骤：
1. 停止Go服务
2. 恢复Python服务
3. 数据库无需回滚（表结构未变）

## 未解决问题

### 问题1：退费流程集成时机

**问题描述**：审批流系统需要与退费流程集成，但退费模块尚未迁移。

**临时方案**：
- 本次实现审批流核心功能
- 预留业务集成接口（如审批完成的回调）
- 退费模块迁移时再集成

**待决定**：
- 审批完成后的业务处理如何触发？
- 是否需要事件驱动架构？

### 问题2：审批流程可视化

**问题描述**：前端需要展示审批流程数轴，需要后端提供完整的节点信息。

**当前方案**：
- 在审批流详情API中返回所有节点信息
- 包含节点状态、完成时间等

**待优化**：
- 是否需要实时更新审批进度？
- 是否需要WebSocket推送？

### 问题3：审批流统计和报表

**问题描述**：业务可能需要审批流的统计数据（如平均审批时长、驳回率等）。

**当前方案**：
- 暂不实现统计功能
- 记录关键时间字段（create_time、handle_time、complete_time）

**待评估**：
- 是否需要额外的统计表？
- 是否需要定时任务聚合数据？

## 参考资料

- [原Python项目审批流实现](D:\claude space\ZhixinStudentSaaS\backend\app.py) - 行6751-8265
- [数据库表结构](D:\claude space\CharonOMS\charonoms_backup_20260127.sql)
- [CharonOMS项目规范](D:\claude space\CharonOMS\openspec\project.md)
- [DDD设计原则](https://domain-driven-design.org/)
