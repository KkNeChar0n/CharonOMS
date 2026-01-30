# 实施任务清单

## 任务依赖关系

```
T1: 准备工作
  ├── T2: 活动模板领域层
  │   ├── T3: 活动模板基础设施层
  │   ├── T4: 活动模板应用层
  │   └── T5: 活动模板接口层
  │
  ├── T6: 活动领域层
  │   ├── T7: 活动基础设施层
  │   ├── T8: 活动应用层
  │   └── T9: 活动接口层
  │
  └── T10: 集成和测试
      └── T11: 文档和部署
```

## 任务列表

### T1: 准备工作

**描述**: 创建项目结构，准备测试数据

**优先级**: P0

**估时**: 0.5天

**任务**:
1. 创建目录结构
   ```
   internal/domain/activity_template/
   internal/domain/activity_template/entity/
   internal/domain/activity_template/repository/
   internal/domain/activity/
   internal/domain/activity/entity/
   internal/domain/activity/repository/
   internal/application/activity_template/
   internal/application/activity/
   internal/infrastructure/persistence/
   internal/interfaces/http/handler/activity_template/
   internal/interfaces/http/handler/activity/
   internal/interfaces/http/dto/
   ```

2. 准备测试数据SQL脚本
   - 基础分类数据
   - 基础商品数据
   - 活动模板测试数据
   - 活动测试数据

3. 确认数据库连接配置正确

**验证标准**:
- 所有目录创建成功
- 测试数据脚本可以正常执行
- 数据库连接正常

**依赖**: 无

---

### T2: 活动模板领域层实现

**描述**: 实现活动模板的领域实体和仓储接口

**优先级**: P0

**估时**: 1天

**任务**:
1. 实现 `ActivityTemplate` 实体 (`internal/domain/activity_template/entity/activity_template.go`)
   - 定义结构体字段
   - GORM 标签映射
   - 验证方法

2. 实现 `ActivityTemplateGoods` 实体 (`internal/domain/activity_template/entity/activity_template_goods.go`)
   - 定义结构体字段
   - GORM 标签映射

3. 定义 `ActivityTemplateRepository` 接口 (`internal/domain/activity_template/repository/repository.go`)
   - Create 方法
   - Update 方法
   - Delete 方法
   - FindByID 方法
   - List 方法
   - FindActiveTemplates 方法
   - CountRelatedActivities 方法

**验证标准**:
- 实体结构定义正确，字段映射到数据库表
- 仓储接口定义清晰，方法签名正确
- 代码通过 `go build`

**依赖**: T1

**被依赖**: T3, T4

---

### T3: 活动模板基础设施层实现

**描述**: 实现活动模板仓储的 GORM 实现

**优先级**: P0

**估时**: 1.5天

**任务**:
1. 实现 `GormActivityTemplateRepository` (`internal/infrastructure/persistence/activity_template_repository.go`)
   - 实现所有接口方法
   - 使用 GORM 进行数据库操作
   - 实现事务支持

2. 编写单元测试 (`internal/infrastructure/persistence/activity_template_repository_test.go`)
   - 测试 Create 方法
   - 测试 Update 方法（包括关联更新）
   - 测试 Delete 方法
   - 测试 FindByID 方法（包括关联查询）
   - 测试 List 方法
   - 测试 FindActiveTemplates 方法
   - 测试 CountRelatedActivities 方法

**验证标准**:
- 所有仓储方法正确实现
- 事务正确处理
- 单元测试覆盖率 > 80%
- 所有测试通过

**依赖**: T2

**被依赖**: T4

---

### T4: 活动模板应用层实现

**描述**: 实现活动模板的应用服务和 DTO 转换

**优先级**: P0

**估时**: 1.5天

**任务**:
1. 定义 DTOs (`internal/interfaces/http/dto/activity_template_dto.go`)
   - CreateTemplateDTO
   - UpdateTemplateDTO
   - TemplateDTO
   - TemplateDetailDTO
   - ClassifyRelationDTO
   - GoodsRelationDTO

2. 实现 Assembler (`internal/application/activity_template/assembler.go`)
   - DTO -> Entity 转换
   - Entity -> DTO 转换
   - 处理关联数据转换

3. 实现 Service (`internal/application/activity_template/service.go`)
   - CreateTemplate 方法
   - UpdateTemplate 方法
   - DeleteTemplate 方法
   - GetTemplate 方法
   - ListTemplates 方法
   - ListActiveTemplates 方法
   - UpdateTemplateStatus 方法

4. 编写单元测试 (`internal/application/activity_template/service_test.go`)
   - 使用 Mock Repository
   - 测试所有业务场景
   - 测试错误处理

**验证标准**:
- DTO 定义与 Python 版本一致
- Service 业务逻辑正确
- 单元测试覆盖率 > 80%
- 所有测试通过

**依赖**: T2, T3

**被依赖**: T5

---

### T5: 活动模板接口层实现

**描述**: 实现活动模板的 HTTP 处理器和路由

**优先级**: P0

**估时**: 1天

**任务**:
1. 实现 Handler (`internal/interfaces/http/handler/activity_template/handler.go`)
   - GetTemplates 方法
   - GetTemplate 方法
   - CreateTemplate 方法
   - UpdateTemplate 方法
   - DeleteTemplate 方法
   - UpdateTemplateStatus 方法
   - GetActiveTemplates 方法

2. 注册路由 (`internal/interfaces/http/router/activity_template_routes.go`)
   - GET /api/activity-templates
   - GET /api/activity-templates/:id
   - POST /api/activity-templates
   - PUT /api/activity-templates/:id
   - DELETE /api/activity-templates/:id
   - PUT /api/activity-templates/:id/status
   - GET /api/activity-templates/active

3. 编写集成测试 (`internal/interfaces/http/handler/activity_template/handler_test.go`)
   - 测试所有 API 接口
   - 验证响应格式
   - 验证HTTP状态码

**验证标准**:
- API 路由与 Python 版本完全一致
- 响应格式与 Python 版本一致
- HTTP 状态码正确
- 错误消息与 Python 版本一致
- 集成测试全部通过

**依赖**: T4

**被依赖**: T10

---

### T6: 活动领域层实现

**描述**: 实现活动的领域实体和仓储接口

**优先级**: P0

**估时**: 1天

**任务**:
1. 实现 `Activity` 实体 (`internal/domain/activity/entity/activity.go`)
   - 定义结构体字段
   - GORM 标签映射
   - 验证方法

2. 实现 `ActivityDetail` 实体 (`internal/domain/activity/entity/activity_detail.go`)
   - 定义结构体字段
   - GORM 标签映射

3. 定义 `ActivityRepository` 接口 (`internal/domain/activity/repository/repository.go`)
   - Create 方法
   - Update 方法
   - Delete 方法
   - FindByID 方法
   - List 方法
   - FindByDateRange 方法

**验证标准**:
- 实体结构定义正确，字段映射到数据库表
- 仓储接口定义清晰，方法签名正确
- 代码通过 `go build`

**依赖**: T1

**被依赖**: T7, T8

---

### T7: 活动基础设施层实现

**描述**: 实现活动仓储的 GORM 实现

**优先级**: P0

**估时**: 1.5天

**任务**:
1. 实现 `GormActivityRepository` (`internal/infrastructure/persistence/activity_repository.go`)
   - 实现所有接口方法
   - 使用 GORM 进行数据库操作
   - 实现事务支持
   - 实现复杂的 JOIN 查询（关联模板信息）

2. 编写单元测试 (`internal/infrastructure/persistence/activity_repository_test.go`)
   - 测试 Create 方法
   - 测试 Update 方法（包括详情更新）
   - 测试 Delete 方法
   - 测试 FindByID 方法（包括关联查询）
   - 测试 List 方法（包括模板信息）
   - 测试 FindByDateRange 方法

**验证标准**:
- 所有仓储方法正确实现
- JOIN 查询正确
- 事务正确处理
- 单元测试覆盖率 > 80%
- 所有测试通过

**依赖**: T6

**被依赖**: T8

---

### T8: 活动应用层实现

**描述**: 实现活动的应用服务和 DTO 转换

**优先级**: P0

**估时**: 2天

**任务**:
1. 定义 DTOs (`internal/interfaces/http/dto/activity_dto.go`)
   - CreateActivityDTO
   - UpdateActivityDTO
   - ActivityDTO
   - ActivityDetailDTO
   - ActivityDetailWithTemplateDTO
   - ActivitiesByDateRangeDTO

2. 实现 Assembler (`internal/application/activity/assembler.go`)
   - DTO -> Entity 转换
   - Entity -> DTO 转换
   - 处理详情数据转换

3. 实现 Service (`internal/application/activity/service.go`)
   - CreateActivity 方法
   - UpdateActivity 方法
   - DeleteActivity 方法
   - GetActivity 方法
   - ListActivities 方法
   - GetActivitiesByDateRange 方法（包括冲突检测）
   - UpdateActivityStatus 方法

4. 编写单元测试 (`internal/application/activity/service_test.go`)
   - 使用 Mock Repository
   - 测试所有业务场景
   - 测试冲突检测逻辑
   - 测试错误处理

**验证标准**:
- DTO 定义与 Python 版本一致
- Service 业务逻辑正确
- 冲突检测逻辑正确
- 单元测试覆盖率 > 80%
- 所有测试通过

**依赖**: T6, T7

**被依赖**: T9

---

### T9: 活动接口层实现

**描述**: 实现活动的 HTTP 处理器和路由

**优先级**: P0

**估时**: 1天

**任务**:
1. 实现 Handler (`internal/interfaces/http/handler/activity/handler.go`)
   - GetActivities 方法
   - GetActivity 方法
   - CreateActivity 方法
   - UpdateActivity 方法
   - DeleteActivity 方法
   - UpdateActivityStatus 方法
   - GetActivitiesByDateRange 方法

2. 注册路由 (`internal/interfaces/http/router/activity_routes.go`)
   - GET /api/activities
   - GET /api/activities/:id
   - POST /api/activities
   - PUT /api/activities/:id
   - DELETE /api/activities/:id
   - PUT /api/activities/:id/status
   - GET /api/activities/by-date-range

3. 编写集成测试 (`internal/interfaces/http/handler/activity/handler_test.go`)
   - 测试所有 API 接口
   - 验证响应格式
   - 验证HTTP状态码
   - 测试冲突检测

**验证标准**:
- API 路由与 Python 版本完全一致
- 响应格式与 Python 版本一致
- HTTP 状态码正确
- 错误消息与 Python 版本一致
- 冲突检测响应格式正确
- 集成测试全部通过

**依赖**: T8

**被依赖**: T10

---

### T10: 集成和端到端测试

**描述**: 完整的集成测试和兼容性验证

**优先级**: P0

**估时**: 1.5天

**任务**:
1. 编写端到端测试
   - 完整的活动模板生命周期测试
   - 完整的活动生命周期测试
   - 跨模块集成测试（活动模板 -> 活动）

2. API 兼容性验证
   - 对比 Python 版本的所有 API 响应
   - 验证请求参数处理一致性
   - 验证错误消息一致性
   - 验证时间格式一致性

3. 性能测试
   - 测试查询性能
   - 测试批量创建性能
   - 识别慢查询

4. 边界条件测试
   - 空数据测试
   - 大数据量测试
   - 并发访问测试

**验证标准**:
- 所有端到端测试通过
- API 响应与 Python 版本完全一致
- 性能满足要求（响应时间 < 200ms）
- 边界条件处理正确

**依赖**: T5, T9

**被依赖**: T11

---

### T11: 文档和部署

**描述**: 更新文档，准备部署

**优先级**: P1

**估时**: 0.5天

**任务**:
1. 更新 API 文档
   - 记录所有新增的 API 接口
   - 添加请求/响应示例

2. 更新规范文档
   - 更新 `openspec/specs/activity-template/spec.md` 的目的说明
   - 更新 `openspec/specs/activity/spec.md` 的目的说明
   - 标记需求为已实现

3. 编写迁移指南
   - 数据库准备步骤
   - 配置修改说明
   - 部署步骤

4. 代码审查
   - 检查代码规范
   - 检查注释完整性
   - 检查测试覆盖率

**验证标准**:
- API 文档完整准确
- 规范文档更新完成
- 迁移指南清晰可执行
- 代码通过审查

**依赖**: T10

---

## 任务执行顺序

### 第一阶段：活动模板模块（5天）
1. T1: 准备工作 (0.5天)
2. T2: 活动模板领域层 (1天)
3. T3: 活动模板基础设施层 (1.5天)
4. T4: 活动模板应用层 (1.5天)
5. T5: 活动模板接口层 (1天)

**里程碑**: 活动模板所有 API 接口正常工作

### 第二阶段：活动模块（4.5天）
6. T6: 活动领域层 (1天)
7. T7: 活动基础设施层 (1.5天)
8. T8: 活动应用层 (2天)
9. T9: 活动接口层 (1天)

**里程碑**: 活动所有 API 接口正常工作

### 第三阶段：集成和部署（2天）
10. T10: 集成和端到端测试 (1.5天)
11. T11: 文档和部署 (0.5天)

**里程碑**: 营销管理模块完整实现并可部署

## 总估时: 11.5天

## 可并行执行的任务

- T2 和 T6 可以并行（两个领域层独立）
- T3 和 T7 可以在各自领域层完成后并行
- T4 和 T8 可以在各自基础设施层完成后并行
- T5 和 T9 可以在各自应用层完成后并行

采用并行执行，总时间可以缩短到约 **8天**。

## 风险缓解

每个任务完成后都应该：
1. 运行所有单元测试
2. 运行已有的集成测试
3. 进行代码审查
4. 及时提交代码

如果某个任务出现问题：
1. 及时沟通和记录
2. 调整后续任务的依赖
3. 必要时调整估时
