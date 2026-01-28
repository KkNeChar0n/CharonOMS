# Change: 实现合同管理模块

## Why

合同管理是教育SaaS系统的核心业务模块之一，用于管理学生与机构之间的培训合同。目前Go项目中尚未实现，导致：

1. **业务流程受阻**：无法完成从学生入学、签约到订单支付的完整业务闭环
2. **数据完整性缺失**：合同是连接学生、订单、收款的关键实体，缺失导致业务数据不完整
3. **用户体验问题**：前端菜单已有"合同管理"入口，但点击后返回"功能开发中"的占位符提示
4. **重构进度滞后**：根据REFACTORING_STATUS.md，合同管理被标记为P0高优先级，需要尽快实现

合同管理模块需要支持完整的生命周期管理，包括新增、撤销、中止合作等状态流转。

## What Changes

### 1. 数据库Schema验证

**contract表（合同主表）**：
- `id` INT PRIMARY KEY - 合同ID
- `name` VARCHAR(200) - 合同名称
- `student_id` INT - 学生ID（外键关联student表）
- `type` TINYINT - 合同类型（0=首报，1=续报）
- `signature_form` TINYINT - 签署形式（0=线上签署，1=线下签署）
- `contract_amount` DECIMAL(10,2) - 合同金额
- `signatory` VARCHAR(100) - 签署方
- `initiating_party` VARCHAR(100) - 发起方（可为空）
- `initiator` VARCHAR(50) - 发起人（创建合同的用户）
- `status` TINYINT - 合同状态（0=待审核，50=已通过，98=已作废，99=协议中止）
- `payment_status` TINYINT - 付款状态（0=未付款，10=部分付款，30=已付款）
- `termination_agreement` VARCHAR(255) - 中止协议文件路径（协议中止时填充）
- `create_time` TIMESTAMP - 创建时间

### 2. API接口实现（5个）

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 获取合同列表 | GET | `/api/contracts` | 返回所有合同（含学生信息） |
| 获取合同详情 | GET | `/api/contracts/:id` | 返回单个合同的完整信息 |
| 新增合同 | POST | `/api/contracts` | 创建新合同，初始状态为待审核 |
| 撤销合同 | PUT | `/api/contracts/:id/revoke` | 撤销待审核状态的合同 |
| 中止合作 | PUT | `/api/contracts/:id/terminate` | 中止已通过状态的合同 |

### 3. DDD四层架构实现

**Domain层**：
- `internal/domain/contract/entity/contract.go` - 合同实体模型
- `internal/domain/contract/repository/contract_repository.go` - 仓储接口

**Infrastructure层**：
- `internal/infrastructure/persistence/mysql/contract/contract_repository_impl.go` - 仓储实现
- 实现GORM查询，包含LEFT JOIN查询（学生信息）

**Application层**：
- `internal/application/service/contract/contract_service.go` - 业务逻辑
- `internal/application/service/contract/dto.go` - 请求/响应DTO

**Interface层**：
- `internal/interfaces/http/handler/contract/contract_handler.go` - HTTP处理器
- 更新 `internal/interfaces/http/router/router.go` - 路由注册

### 4. 业务逻辑实现

**验证规则**：
- 必填字段：合同名称、学生ID、合同类型、签署形式、合同金额
- 合同类型：只能是0（首报）或1（续报）
- 签署形式：只能是0（线上签署）或1（线下签署）
- 合同状态：新建时默认为0（待审核）
- 付款状态：新建时默认为0（未付款）

**状态流转规则**：
- 待审核(0) → 已作废(98)：通过"撤销"操作
- 已通过(50) → 协议中止(99)：通过"中止合作"操作
- 撤销操作：仅在status=0时允许
- 中止合作：仅在status=50时允许，需提供termination_agreement

**关联关系**：
- 合同与学生：多对一（必选）
- 合同与订单：一对多（通过订单表的contract_id外键，待订单模块实现）

### 5. 查询优化

**列表查询**：
- 使用LEFT JOIN查询关联的学生信息
- 支持前端过滤：ID、学生ID、学生姓名、合同类型、合同状态、付款状态
- 按创建时间倒序排列

**响应格式**（与原项目保持一致）：
```json
{
  "contracts": [
    {
      "id": 1,
      "name": "1001张三首报合同",
      "student_id": 1,
      "student_name": "张三",
      "type": 0,
      "signature_form": 1,
      "contract_amount": 12000.00,
      "signatory": "张三家长",
      "initiating_party": "",
      "initiator": "admin",
      "status": 0,
      "payment_status": 0,
      "termination_agreement": "",
      "create_time": "2026-01-28 10:00:00"
    }
  ]
}
```

### 6. 付款状态字段处理

**当前阶段**：
- 数据库表中保留 `payment_status` 字段
- API响应中返回该字段
- 不实现付款状态的更新逻辑
- 等待后续订单收款模块实现后，由收款模块自动更新该字段

## Impact

### 影响的规范
- 新增 `specs/contract/spec.md` - 合同管理功能规范

### 影响的代码

**新增文件**：
- `internal/domain/contract/` - 合同领域层（实体、仓储接口）
- `internal/infrastructure/persistence/mysql/contract/` - 合同数据访问层
- `internal/application/service/contract/` - 合同应用服务层
- `internal/interfaces/http/handler/contract/` - 合同HTTP处理器

**修改文件**：
- `internal/interfaces/http/router/router.go` - 注册合同管理路由，替换占位符
- `REFACTORING_STATUS.md` - 更新合同管理模块完成状态

### 依赖关系

**依赖的模块**（已实现）：
- ✅ 学生管理（student） - `GET /api/students`
- ✅ 认证系统 - JWT中间件和用户信息获取

**被依赖的模块**（待实现）：
- ❌ 订单管理 - 需要合同ID字段（contract_id）
- ❌ 收款管理 - 需要合同信息进行收款关联

### 破坏性变更

**路由变更**：
- 将 `GET /api/contracts` 从占位符改为真实实现
- 新增 `/api/contracts/:id` 详情接口
- 新增 `/api/contracts/:id/revoke` 撤销接口
- 新增 `/api/contracts/:id/terminate` 中止接口

这些变更是功能增强，不会破坏现有功能。

## Migration Path

### Phase 1: 数据库准备（如需要）
1. 验证contract表是否已存在，如不存在则创建
2. 验证字段完整性（所有必需字段存在）
3. 验证外键约束：`student_id` 引用 `student.id`
4. 可选：导入测试数据

### Phase 2: 后端实现
1. 实现Domain层（实体和仓储接口）
2. 实现Infrastructure层（GORM仓储实现）
3. 实现Application层（业务服务和DTO）
4. 实现Interface层（HTTP处理器）
5. 更新路由配置

### Phase 3: 测试验证
1. 单元测试：业务逻辑验证
2. 集成测试：API接口测试
3. 前端联调：验证与前端的集成

### Phase 4: 文档更新
1. 更新 `REFACTORING_STATUS.md` 标记合同管理为已完成
2. 更新 `README.md` 添加合同管理API文档

## Validation

### 验证标准

- [ ] 数据库表验证完成，包含正确的字段和索引
- [ ] 5个API接口全部实现并通过测试
- [ ] 合同列表返回格式与原项目一致
- [ ] 新增合同功能正常，初始状态正确
- [ ] 撤销合同功能正常，状态验证正确
- [ ] 中止合作功能正常，状态验证正确
- [ ] 前端页面可以正常显示合同列表
- [ ] 前端可以执行新增、撤销、中止操作

### 测试场景

#### 1. 基本CRUD操作
- 创建合同（所有必填字段）
- 查询合同列表，验证学生信息正确关联
- 查询合同详情，验证所有字段完整
- 验证不同合同类型和签署形式

#### 2. 状态流转
- 新增合同后验证status=0, payment_status=0
- 撤销待审核合同，验证status变为98
- 尝试撤销非待审核合同，验证返回错误
- 中止已通过合同，验证status变为99
- 尝试中止非已通过合同，验证返回错误

#### 3. 关联关系
- 查询合同时验证学生姓名正确显示
- 查询不存在的学生ID，验证外键约束生效

#### 4. 边界条件
- 尝试创建合同（缺少必填字段） - 应返回400错误
- 尝试撤销不存在的合同 - 应返回404错误
- 尝试中止合作但不提供termination_agreement - 应返回400错误

#### 5. 前端集成
- 前端访问合同管理页面显示合同列表
- 前端筛选：按ID、学生ID、学生姓名、状态过滤
- 前端分页：验证分页功能正常
- 前端新增：打开新增对话框，选择学生，提交成功
- 前端详情：查看合同详情，操作按钮根据状态显示
- 前端撤销：点击撤销按钮，确认后成功撤销
- 前端中止：点击中止按钮，上传文件，成功中止

## Notes

1. **与原Python项目的差异**：
   - 原项目中没有"编辑合同"功能，Go版本也不实现
   - 原项目中没有"删除合同"功能，Go版本也不实现
   - 保持业务逻辑完全一致，确保前端无缝兼容

2. **付款状态处理**：
   - 当前阶段仅保留字段，不实现更新逻辑
   - 未来订单收款模块会自动更新该字段
   - 合同列表和详情中正常显示该字段值

3. **文件上传处理**：
   - 中止合作需要上传termination_agreement文件
   - 文件上传逻辑需要实现（保存到服务器，存储路径）
   - 建议使用统一的文件上传服务

4. **响应格式**：严格按照原项目的响应格式（直接返回contracts数组），与前端代码保持兼容

5. **状态字段说明**：
   - status: 0=待审核, 50=已通过, 98=已作废, 99=协议中止
   - payment_status: 0=未付款, 10=部分付款, 30=已付款
   - type: 0=首报, 1=续报
   - signature_form: 0=线上签署, 1=线下签署

6. **数据库字符集**：确保使用utf8mb4字符集，支持中文合同名称和签署方名称

7. **发起人字段**：从JWT token中获取当前登录用户名，自动填充initiator字段
