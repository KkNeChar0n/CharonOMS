# 任务清单：营销管理模块迁移

## 阶段 1：数据库表创建（30分钟）

### 1.1 创建activity表
- [ ] 编写 `scripts/create_activity_tables.sql`
  - activity 表结构
  - activity_detail 表结构
  - 添加索引

**SQL脚本**：
```sql
CREATE TABLE `activity` (...);
CREATE TABLE `activity_detail` (...);
```

### 1.2 执行SQL脚本
- [ ] 在数据库中执行脚本
- [ ] 验证表已创建
- [ ] 验证索引已创建

**验证**：`SHOW TABLES LIKE 'activity%';`

## 阶段 2：基础设施和数据模型（2-3小时）

### 2.1 创建领域实体
- [ ] 创建 `internal/domain/marketing/entity/activity_template.go`
  - ActivityTemplate 实体
  - ActivityTemplateClassify 实体
  - ActivityTemplateGoods 实体
  - TableName 方法
- [ ] 创建 `internal/domain/marketing/entity/activity.go`
  - Activity 实体
  - ActivityDetail 实体
  - TableName 方法

**验证**：运行 `go build` 确保编译通过

### 2.2 定义仓储接口
- [ ] 创建 `internal/domain/marketing/repository/template_repository.go`
  - List 方法（支持筛选）
  - GetActiveList 方法（获取启用的模板）
  - GetByID 方法（包含关联）
  - Create 方法（模板+关联）
  - Update 方法（模板+关联）
  - UpdateStatus 方法
- [ ] 创建 `internal/domain/marketing/repository/activity_repository.go`
  - List 方法（支持筛选）
  - GetByID 方法（包含细节）
  - Create 方法（活动+细节）
  - Update 方法（活动+细节）
  - UpdateStatus 方法

**验证**：接口定义清晰，方法签名完整

## 阶段 3：基础设施层实现（3-4小时）

### 3.1 实现活动模板仓储
- [ ] 创建 `internal/infrastructure/persistence/mysql/marketing/template_repository_impl.go`
  - 实现 TemplateRepository 接口
  - List 方法：支持 id、name、type、status 筛选
  - GetActiveList 方法：只返回 status=0 的模板
  - GetByID 方法：根据 select_type 动态 Preload Classifies 或 Goods
  - Create 方法：使用事务创建模板和关联（根据 select_type）
  - Update 方法：使用事务更新模板，先删除旧关联再创建新关联
  - UpdateStatus 方法

**验证**：
- 测试按分类和按商品两种模式的创建和查询
- 测试更新时关联数据的正确替换

### 3.2 实现活动仓储
- [ ] 创建 `internal/infrastructure/persistence/mysql/marketing/activity_repository_impl.go`
  - 实现 ActivityRepository 接口
  - List 方法：支持 id、name、template_id、status 筛选，Preload Template
  - GetByID 方法：Preload Template 和 Details（Details 需要 Preload Goods）
  - Create 方法：使用事务创建活动和细节
  - Update 方法：使用事务更新活动，先删除旧细节再创建新细节
  - UpdateStatus 方法

**验证**：
- 测试活动列表的 template_name 和 template_type 是否正确返回
- 测试活动详情的细节数据是否完整

## 阶段 4：应用服务层（2-3小时）

### 4.1 实现活动模板服务
- [ ] 创建 `internal/application/service/marketing/template_service.go`
  - GetTemplateList 方法（支持筛选参数）
  - GetActiveTemplateList 方法
  - GetTemplateByID 方法
  - CreateTemplate 方法（验证请求参数）
  - UpdateTemplate 方法（验证状态和参数）
  - UpdateTemplateStatus 方法
- [ ] 定义请求和响应 DTO
  - TemplateListRequest
  - CreateTemplateRequest
  - UpdateTemplateRequest
  - UpdateStatusRequest

**验证**：
- 测试 select_type 的验证逻辑
- 测试启用状态下禁止编辑的业务规则

### 4.2 实现活动服务
- [ ] 创建 `internal/application/service/marketing/activity_service.go`
  - GetActivityList 方法（支持筛选参数）
  - GetActivityByID 方法
  - CreateActivity 方法（验证请求参数，包括时间范围、模板状态）
  - UpdateActivity 方法（验证状态和参数）
  - UpdateActivityStatus 方法
- [ ] 定义请求和响应 DTO
  - ActivityListRequest
  - CreateActivityRequest
  - UpdateActivityRequest

**验证**：
- 测试时间范围验证（start_time < end_time）
- 测试模板状态验证（只能使用启用的模板）
- 测试启用状态下禁止编辑的业务规则

## 阶段 5：HTTP 接口层（2-3小时）

### 5.1 实现活动模板 Handler
- [ ] 创建 `internal/interfaces/http/handler/marketing/template_handler.go`
  - GetTemplates 方法（GET /api/activity-templates）
  - GetActiveTemplates 方法（GET /api/activity-templates/active）
  - GetTemplateDetail 方法（GET /api/activity-templates/:id）
  - CreateTemplate 方法（POST /api/activity-templates）
  - UpdateTemplate 方法（PUT /api/activity-templates/:id）
  - UpdateTemplateStatus 方法（PUT /api/activity-templates/:id/status）

**响应格式**：
- 列表接口：`{"activity_templates": [...]}`
- 详情接口：直接返回对象
- 错误响应：`{"error": "错误信息"}`

### 5.2 实现活动 Handler
- [ ] 创建 `internal/interfaces/http/handler/marketing/activity_handler.go`
  - GetActivities 方法（GET /api/activities）
  - GetActivityDetail 方法（GET /api/activities/:id）
  - CreateActivity 方法（POST /api/activities）
  - UpdateActivity 方法（PUT /api/activities/:id）
  - UpdateActivityStatus 方法（PUT /api/activities/:id/status）

**响应格式**：
- 列表接口：`{"activities": [...]}`
- 详情接口：直接返回对象

**验证**：
- 使用 curl 或 Postman 测试所有接口
- 验证响应格式与 Python 版本一致

## 阶段 6：路由注册和依赖注入（1小时）

### 6.1 注册路由
- [ ] 修改 `internal/interfaces/http/router/router.go`
  - 注入活动模板和活动的仓储、服务、Handler
  - 注册活动模板相关路由（6个）
  - 注册活动相关路由（5个）
  - 所有路由使用 JWTAuth 中间件保护

**依赖关系**：
```
仓储实现 → 服务 → Handler → 路由注册
```

**验证**：
- 启动服务，检查日志确认路由已注册
- 测试未登录访问返回 401

## 阶段 7：权限集成（2小时）

### 7.1 查询菜单ID
- [ ] 查询数据库获取活动模板、活动管理的菜单ID
  ```sql
  SELECT id, name FROM menu WHERE name IN ('活动模板', '活动管理');
  ```

### 7.2 创建权限播种脚本
- [ ] 创建 `scripts/seed_marketing_permissions.sql`
  - 插入活动模板权限（view, add, edit, enable, disable）
  - 插入活动权限（view, add, edit, enable, disable）

**权限清单**：
- `view_activity_template` - 查看活动模板
- `add_activity_template` - 新增活动模板
- `edit_activity_template` - 编辑活动模板
- `enable_activity_template` - 启用活动模板
- `disable_activity_template` - 禁用活动模板
- `view_activity` - 查看活动
- `add_activity` - 新增活动
- `edit_activity` - 编辑活动
- `enable_activity` - 启用活动
- `disable_activity` - 禁用活动

### 7.3 执行权限播种
- [ ] 执行 SQL 脚本插入权限数据
- [ ] 验证权限已正确创建并关联菜单

**验证**：
- 查询 permissions 表确认权限已创建
- 登录超级管理员账号，验证所有按钮可见
- 创建测试角色，分配部分权限，验证按钮权限控制

## 阶段 8：前后端联调和测试（2-3小时）

### 8.1 活动模板管理测试
- [ ] 测试查询所有活动模板
- [ ] 测试按类型、名称筛选
- [ ] 测试获取启用的模板列表
- [ ] 测试获取模板详情（按分类和按商品）
- [ ] 测试新增活动模板（按分类）
- [ ] 测试新增活动模板（按商品）
- [ ] 测试更新活动模板
- [ ] 测试启用/禁用操作
- [ ] 验证权限控制（新增、编辑、启用/禁用按钮）
- [ ] 验证启用状态下禁止编辑

### 8.2 活动管理测试
- [ ] 测试查询所有活动
- [ ] 测试按模板、名称筛选
- [ ] 测试获取活动详情（包含细节）
- [ ] 测试新增活动（包含活动细节）
- [ ] 测试更新活动
- [ ] 测试启用/禁用操作
- [ ] 验证权限控制（新增、编辑、启用/禁用按钮）
- [ ] 验证启用状态下禁止编辑
- [ ] 验证时间范围验证
- [ ] 验证模板状态验证（只能使用启用的模板）

### 8.3 异常场景测试
- [ ] 测试未登录访问（应返回 401）
- [ ] 测试无权限访问（按钮不可见）
- [ ] 测试创建时参数缺失（应返回错误）
- [ ] 测试更新启用状态的模板/活动（应返回错误）
- [ ] 测试时间范围无效（应返回错误）
- [ ] 测试使用禁用的模板创建活动（应返回错误）

## 阶段 9：文档和归档（1小时）

### 9.1 代码提交
- [ ] 提交所有代码到 Git
- [ ] 编写清晰的提交信息

**提交信息格式**：
```
feat: 实现营销管理模块迁移

主要功能：
1. 活动模板管理（查询、详情、新增、更新、状态更新）
2. 活动管理（查询、详情、新增、更新、状态更新）
3. 权限集成（10个权限，支持按钮级控制）
4. 数据库表创建（activity、activity_detail）

技术细节：
- DDD 四层架构
- 支持按分类和按商品两种模板类型
- 活动与活动细节的事务处理
- 时间范围和状态验证

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

### 9.2 归档提案
- [ ] 使用 `openspec-cn archive migrate-marketing-module` 归档提案
- [ ] 验证归档成功

## 总计时间估算
- 数据库表创建：30分钟
- 基础设施和数据模型：2-3小时
- 基础设施层实现：3-4小时
- 应用服务层：2-3小时
- HTTP 接口层：2-3小时
- 路由注册和依赖注入：1小时
- 权限集成：2小时
- 前后端联调和测试：2-3小时
- 文档和归档：1小时

**总计：16-21小时**

## 依赖关系
- 阶段 1 → 阶段 2：表创建后才能定义实体
- 阶段 2 → 阶段 3：实体定义完成后才能实现仓储
- 阶段 3 → 阶段 4：仓储实现完成后才能实现服务
- 阶段 4 → 阶段 5：服务实现完成后才能实现 Handler
- 阶段 5 → 阶段 6：Handler 实现完成后才能注册路由
- 阶段 6 → 阶段 7：路由注册后才能集成权限
- 阶段 7 → 阶段 8：权限集成后才能进行完整测试

## 并行工作机会
- 阶段 2.1 和 2.2 可以并行（实体和接口定义）
- 阶段 3.1 和 3.2 可以并行（模板和活动仓储）
- 阶段 4.1 和 4.2 可以并行（模板和活动服务）
- 阶段 5.1 和 5.2 可以并行（模板和活动 Handler）
