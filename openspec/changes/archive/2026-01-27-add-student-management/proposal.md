# Change: 实现学生管理模块

## Why

学生管理是教育SaaS系统的核心业务模块，目前Go项目中尚未实现。学生数据是订单、合同、收款等业务流程的基础，缺少学生管理模块将导致：

1. **业务流程受阻**：无法录入学生信息，后续的订单创建、合同签订、收款管理等功能无法开展
2. **数据完整性缺失**：学生是系统的核心实体，与教练、订单、合同、收款等模块都有关联关系
3. **用户体验问题**：前端菜单已有"学生管理"入口，但点击后返回"功能开发中"的占位符提示
4. **重构进度滞后**：根据REFACTORING_STATUS.md，学生管理被标记为P0高优先级，需要尽快实现

从原Python Flask项目搬运学生管理功能是继续项目重构的关键步骤。

## What Changes

### 1. 数据库Schema创建

创建学生相关表结构：

**student表（学生主表）**：
- `id` INT PRIMARY KEY - 学生ID
- `name` VARCHAR(100) - 学生姓名
- `sex_id` INT - 性别ID（外键关联sex表）
- `grade_id` INT - 年级ID（外键关联grade表）
- `phone` VARCHAR(20) - 联系电话
- `status` TINYINT - 状态（0=启用，1=禁用）
- `created_at` TIMESTAMP - 创建时间
- `updated_at` TIMESTAMP - 更新时间

**student_coach表（学生与教练的多对多关系表）**：
- `id` INT PRIMARY KEY
- `student_id` INT - 学生ID（外键）
- `coach_id` INT - 教练ID（外键）
- `created_at` TIMESTAMP - 创建时间

### 2. API接口实现（6个）

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 获取学生列表 | GET | `/api/students` | 返回所有学生（含性别、年级、教练信息） |
| 新增学生 | POST | `/api/students` | 创建新学生，可选关联教练 |
| 更新学生信息 | PUT | `/api/students/:id` | 更新学生基本信息 |
| 更新学生状态 | PUT | `/api/students/:id/status` | 切换启用/禁用状态 |
| 删除学生 | DELETE | `/api/students/:id` | 删除学生（需检查关联订单） |
| 获取启用学生 | GET | `/api/students/active` | 仅返回启用状态的学生列表 |

### 3. DDD四层架构实现

**Domain层**：
- `internal/domain/student/entity/student.go` - 学生实体模型
- `internal/domain/student/repository/student_repository.go` - 仓储接口

**Infrastructure层**：
- `internal/infrastructure/persistence/mysql/student/student_repository_impl.go` - 仓储实现
- 实现GORM查询，包含JOIN查询（性别、年级、教练）

**Application层**：
- `internal/application/service/student/student_service.go` - 业务逻辑
- `internal/application/service/student/dto.go` - 请求/响应DTO

**Interface层**：
- `internal/interfaces/http/handler/student/student_handler.go` - HTTP处理器
- `internal/interfaces/http/router/router.go` - 路由注册

### 4. 业务逻辑实现

**验证规则**：
- 必填字段：学生姓名、性别ID、年级ID、联系电话
- 状态值：只能是0（启用）或1（禁用）

**关联关系**：
- 学生与性别：多对一（必选）
- 学生与年级：多对一（必选）
- 学生与教练：多对多（可选，通过student_coach表）

**删除约束**：
- 删除前必须检查是否有关联订单，有则阻止删除
- 删除学生时级联删除student_coach表中的关联记录

### 5. 查询优化

**列表查询**：
- 使用LEFT JOIN查询关联的性别、年级、教练信息
- 使用GROUP_CONCAT聚合教练名称（逗号分隔）
- 前端过滤支持：ID、姓名（模糊）、年级、状态

**响应格式**（与原项目保持一致）：
```json
{
  "students": [
    {
      "id": 1,
      "student_name": "张三",
      "sex_id": 1,
      "sex": "男",
      "grade_id": 2,
      "grade": "初中",
      "phone": "13800138000",
      "status": 0,
      "coach_names": "李老师, 王老师"
    }
  ]
}
```

### 6. 数据库迁移脚本

创建Schema迁移脚本：
- `scripts/migrations/001_create_student_tables.sql` - 创建student和student_coach表
- 包含索引、外键约束、初始化数据

## Impact

### 影响的规范
- 新增 `specs/student/spec.md` - 学生管理功能规范

### 影响的代码

**新增文件**：
- `internal/domain/student/` - 学生领域层（实体、仓储接口）
- `internal/infrastructure/persistence/mysql/student/` - 学生数据访问层
- `internal/application/service/student/` - 学生应用服务层
- `internal/interfaces/http/handler/student/` - 学生HTTP处理器
- `scripts/migrations/001_create_student_tables.sql` - 数据库迁移

**修改文件**：
- `internal/interfaces/http/router/router.go` - 注册学生管理路由
- 删除占位符路由 `/students`

### 依赖关系

**依赖的模块**（已实现）：
- ✅ 性别管理（sex） - `GET /api/sexes`
- ✅ 年级管理（grade） - `GET /api/grades/active`
- ❌ 教练管理（coach） - 尚未实现，但不影响学生管理的基本功能

**被依赖的模块**（待实现）：
- 订单管理 - 需要学生信息
- 合同管理 - 需要学生信息
- 收款管理 - 需要学生信息

### 破坏性变更

本变更不涉及破坏性变更，是纯粹的功能新增。

## Migration Path

### Phase 1: 数据库准备
1. 执行数据库迁移脚本创建student和student_coach表
2. 验证表结构和索引
3. 可选：导入测试数据

### Phase 2: 后端实现
1. 实现Domain层（实体和仓储接口）
2. 实现Infrastructure层（GORM仓储实现）
3. 实现Application层（业务服务和DTO）
4. 实现Interface层（HTTP处理器）
5. 注册路由

### Phase 3: 测试验证
1. 单元测试：业务逻辑验证
2. 集成测试：API接口测试
3. 前端联调：验证与前端的集成

### Phase 4: 文档更新
1. 更新 `REFACTORING_STATUS.md` 标记学生管理为已完成
2. 更新 `README.md` 添加学生管理API文档

## Validation

### 验证标准

- [ ] 数据库表创建成功，包含正确的字段和索引
- [ ] 6个API接口全部实现并通过测试
- [ ] 学生列表返回格式与原项目一致
- [ ] 新增学生可以正确关联教练（多对多）
- [ ] 删除学生时正确检查订单关联
- [ ] 学生状态切换功能正常
- [ ] 前端页面可以正常显示学生列表
- [ ] 前端可以执行新增、编辑、删除操作

### 测试场景

#### 1. 基本CRUD操作
- 创建学生（必填字段）
- 创建学生并关联教练
- 查询学生列表，验证关联信息正确显示
- 更新学生信息
- 删除学生（无订单关联）

#### 2. 状态管理
- 切换学生状态为禁用
- 查询启用学生列表，验证已禁用的学生不出现
- 切换学生状态为启用

#### 3. 关联关系
- 查询学生时验证性别和年级信息正确显示
- 查询学生时验证教练名称正确聚合
- 删除教练后，验证student_coach表记录自动清理

#### 4. 边界条件
- 尝试创建学生（缺少必填字段） - 应返回400错误
- 尝试删除有订单的学生 - 应返回错误提示
- 尝试更新不存在的学生 - 应返回404错误
- 尝试设置非法状态值 - 应返回400错误

#### 5. 前端集成
- 前端访问 `http://localhost:5001/#/students` 显示学生列表
- 前端筛选：按ID、姓名、年级、状态过滤
- 前端分页：验证分页功能正常
- 前端新增：打开新增对话框，提交成功
- 前端编辑：打开编辑对话框，提交成功
- 前端删除：删除确认对话框，删除成功

## Notes

1. **教练管理依赖**：虽然学生可以关联教练，但教练管理模块尚未实现。建议先实现学生管理的核心功能，教练关联功能可以后续补充。

2. **订单依赖检查**：删除学生时的订单检查功能依赖订单表存在，如果订单表未创建，暂时可以跳过此检查，待订单模块实现后再补充。

3. **响应格式**：严格按照原项目的响应格式（直接返回students数组），与前端代码保持兼容。

4. **数据库字符集**：确保使用utf8mb4字符集，支持中文姓名和特殊字符。

5. **性能优化**：学生列表查询使用LEFT JOIN和GROUP BY，可能影响性能。如果学生数量较大，考虑添加索引或分页查询。
