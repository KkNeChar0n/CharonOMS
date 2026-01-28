# Change: 实现教练管理模块

## Why

教练管理是教育SaaS系统的核心业务模块，与学生管理同等重要。目前Go项目中尚未实现，导致：

1. **业务流程受阻**：无法录入教练信息，学生-教练关联功能不完整
2. **数据完整性缺失**：教练是系统的核心实体，与学生、学科、课程等模块都有关联关系
3. **用户体验问题**：前端菜单已有"教练管理"入口，但点击后返回"功能开发中"的占位符提示
4. **重构进度滞后**：根据REFACTORING_STATUS.md，教练管理被标记为P0高优先级，需要尽快实现

教练管理模块与已实现的学生管理模块设计类似，可以复用大部分架构模式。

## What Changes

### 1. 数据库Schema

**coach表（教练主表）**：
- `id` INT PRIMARY KEY - 教练ID
- `name` VARCHAR(100) - 教练姓名
- `sex_id` INT - 性别ID（外键关联sex表）
- `subject_id` INT - 学科ID（外键关联subject表）
- `phone` VARCHAR(20) - 联系电话
- `status` TINYINT - 状态（0=启用，1=禁用）
- `created_at` TIMESTAMP - 创建时间
- `updated_at` TIMESTAMP - 更新时间

**student_coach表（学生与教练的多对多关系表）**：
- 已在学生管理模块中创建
- 需要补充coach外键约束（如果尚未添加）

### 2. API接口实现（6个）

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 获取教练列表 | GET | `/api/coaches` | 返回所有教练（含性别、学科信息） |
| 获取启用教练 | GET | `/api/coaches/active` | 仅返回启用状态的教练列表 |
| 新增教练 | POST | `/api/coaches` | 创建新教练，可选关联学生 |
| 更新教练信息 | PUT | `/api/coaches/:id` | 更新教练基本信息 |
| 更新教练状态 | PUT | `/api/coaches/:id/status` | 切换启用/禁用状态 |
| 删除教练 | DELETE | `/api/coaches/:id` | 删除教练（级联删除student_coach关联） |

### 3. DDD四层架构实现

**Domain层**：
- `internal/domain/coach/entity/coach.go` - 教练实体模型
- `internal/domain/coach/repository/coach_repository.go` - 仓储接口

**Infrastructure层**：
- `internal/infrastructure/persistence/mysql/coach/coach_repository_impl.go` - 仓储实现
- 实现GORM查询，包含JOIN查询（性别、学科）

**Application层**：
- `internal/application/service/coach/coach_service.go` - 业务逻辑
- `internal/application/service/coach/dto.go` - 请求/响应DTO

**Interface层**：
- `internal/interfaces/http/handler/coach/coach_handler.go` - HTTP处理器
- 更新 `internal/interfaces/http/router/router.go` - 路由注册

### 4. 业务逻辑实现

**验证规则**：
- 必填字段：教练姓名、性别ID、学科ID、联系电话
- 状态值：只能是0（启用）或1（禁用）

**关联关系**：
- 教练与性别：多对一（必选）
- 教练与学科：多对一（必选）
- 教练与学生：多对多（可选，通过student_coach表）

**删除约束**：
- 删除教练时级联删除student_coach表中的关联记录
- 使用事务保证数据一致性

### 5. 查询优化

**列表查询**：
- 使用LEFT JOIN查询关联的性别、学科信息
- 前端过滤支持：ID、姓名（模糊）、性别、学科、状态

**响应格式**（与原项目保持一致）：
```json
{
  "coaches": [
    {
      "id": 1,
      "coach_name": "李老师",
      "sex_id": 2,
      "sex": "女",
      "subject_id": 1,
      "subject": "数学",
      "phone": "13800138001",
      "status": 0
    }
  ]
}
```

### 6. 与学生管理模块的集成

**更新学生管理模块**：
- 修改 `GET /api/coaches/active` 从临时占位符改为真实查询
- 确保学生列表查询中的教练关联正常工作

## Impact

### 影响的规范
- 新增 `specs/coach/spec.md` - 教练管理功能规范

### 影响的代码

**新增文件**：
- `internal/domain/coach/` - 教练领域层（实体、仓储接口）
- `internal/infrastructure/persistence/mysql/coach/` - 教练数据访问层
- `internal/application/service/coach/` - 教练应用服务层
- `internal/interfaces/http/handler/coach/` - 教练HTTP处理器

**修改文件**：
- `internal/interfaces/http/router/router.go` - 注册教练管理路由，替换占位符
- `REFACTORING_STATUS.md` - 更新教练管理模块完成状态

### 依赖关系

**依赖的模块**（已实现）：
- ✅ 性别管理（sex） - `GET /api/sexes`
- ✅ 学科管理（subject） - `GET /api/subjects/active`
- ✅ 学生管理（student） - `GET /api/students/active`（用于新增教练时关联学生）
- ✅ student_coach表 - 已在学生管理模块中创建

**被依赖的模块**（待实现）：
- 课程管理 - 需要教练信息
- 排课管理 - 需要教练信息
- 学生列表 - 需要显示关联的教练信息（已实现基础功能）

### 破坏性变更

**路由变更**：
- 将 `GET /api/coaches/active` 从临时占位符改为真实实现
- 其他教练路由从占位符改为真实实现

这些变更是功能增强，不会破坏现有功能。

## Migration Path

### Phase 1: 数据库准备（如需要）
1. 验证coach表是否已存在，如不存在则创建
2. 验证student_coach表的coach外键约束
3. 可选：导入测试数据

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
1. 更新 `REFACTORING_STATUS.md` 标记教练管理为已完成
2. 更新 `README.md` 添加教练管理API文档

## Validation

### 验证标准

- [ ] 数据库表创建成功，包含正确的字段和索引
- [ ] 6个API接口全部实现并通过测试
- [ ] 教练列表返回格式与原项目一致
- [ ] 新增教练可以正确关联学生（多对多）
- [ ] 删除教练时正确级联删除关联关系
- [ ] 教练状态切换功能正常
- [ ] 前端页面可以正常显示教练列表
- [ ] 前端可以执行新增、编辑、删除操作

### 测试场景

#### 1. 基本CRUD操作
- 创建教练（必填字段）
- 创建教练并关联学生
- 查询教练列表，验证关联信息正确显示
- 更新教练信息
- 删除教练（自动清理student_coach关联）

#### 2. 状态管理
- 切换教练状态为禁用
- 查询启用教练列表，验证已禁用的教练不出现
- 切换教练状态为启用

#### 3. 关联关系
- 查询教练时验证性别和学科信息正确显示
- 新增教练时关联学生，验证student_coach表记录正确创建
- 删除教练后，验证student_coach表记录自动清理

#### 4. 边界条件
- 尝试创建教练（缺少必填字段） - 应返回400错误
- 尝试更新不存在的教练 - 应返回404错误
- 尝试设置非法状态值 - 应返回400错误

#### 5. 前端集成
- 前端访问 `http://localhost:5001/#/coaches` 显示教练列表
- 前端筛选：按ID、姓名、性别、学科、状态过滤
- 前端分页：验证分页功能正常
- 前端新增：打开新增对话框，选择关联学生，提交成功
- 前端编辑：打开编辑对话框，提交成功
- 前端删除：删除确认对话框，删除成功

#### 6. 与学生管理的集成
- 学生列表查询显示关联的教练名称
- 新增学生时可以选择教练
- 删除教练后，学生列表中的教练名称正确更新

## Notes

1. **与学生管理模块的相似性**：教练管理与学生管理在架构和实现上高度相似，可以参考学生管理模块的实现模式。

2. **student_coach表复用**：该表已在学生管理模块中创建，教练管理模块直接复用，需确保外键约束完整。

3. **响应格式**：严格按照原项目的响应格式（直接返回coaches数组），与前端代码保持兼容。

4. **数据库字符集**：确保使用utf8mb4字符集，支持中文姓名和特殊字符。

5. **状态字段验证**：参考学生管理模块的实现，使用指针类型（`*int`）处理status=0的验证问题。

6. **前端字段名称**：前端使用`coach_name`而非`name`，DTO需要匹配。
