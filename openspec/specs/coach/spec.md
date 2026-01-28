# 规范：教练管理模块

## 新增需求

### 需求：教练基本信息管理

系统SHALL提供教练基本信息的完整管理功能，包括创建、查询、更新和删除操作。

#### 场景：创建教练

**GIVEN** 用户已登录系统
**WHEN** 用户提交创建教练请求，包含必填字段（姓名、性别ID、学科ID、联系电话）
**THEN** 系统SHALL：
- 验证所有必填字段是否存在
- 创建新的教练记录
- 返回201状态码和教练ID
- 如果提供了student_ids，SHALL同时创建student_coach关联记录
- 返回消息："教练添加成功"

**WHEN** 用户提交创建教练请求，缺少必填字段
**THEN** 系统SHALL：
- 返回400错误
- 提供清晰的错误信息说明缺少哪个字段

#### 场景：查询教练列表

**GIVEN** 系统中存在教练数据
**WHEN** 用户请求教练列表
**THEN** 系统SHALL：
- 返回所有教练的完整信息
- 包含关联的性别名称（通过sex_id关联sex表）
- 包含关联的学科名称（通过subject_id关联subject表）
- 响应格式为：`{"coaches": [...]}`
- 每个教练对象包含字段：id, coach_name, sex_id, sex, subject_id, subject, phone, status

#### 场景：查询启用状态的教练

**GIVEN** 系统中存在启用和禁用的教练
**WHEN** 用户请求启用教练列表（GET /api/coaches/active）
**THEN** 系统SHALL：
- 仅返回status=0的教练
- 返回简化格式：`{"coaches": [{"id": 1, "coach_name": "..."}]}`
- 不包含禁用（status=1）的教练

#### 场景：更新教练信息

**GIVEN** 系统中存在指定ID的教练
**WHEN** 用户提交更新教练请求
**THEN** 系统SHALL：
- 验证教练ID是否存在
- 验证必填字段是否完整
- 更新教练基本信息
- 返回消息："教练信息更新成功"

**WHEN** 用户尝试更新不存在的教练
**THEN** 系统SHALL：
- 返回404错误
- 提供错误信息说明教练不存在

#### 场景：删除教练

**GIVEN** 系统中存在指定ID的教练
**WHEN** 用户提交删除教练请求
**THEN** 系统SHALL：
- 删除student_coach表中的所有关联记录
- 删除coach表中的教练记录
- 使用数据库事务保证操作原子性
- 返回消息："教练删除成功"

### 需求：教练状态管理

系统SHALL提供教练启用/禁用状态的管理功能。

#### 场景：更新教练状态

**GIVEN** 系统中存在指定ID的教练
**WHEN** 用户提交更新状态请求（PUT /api/coaches/:id/status）
**THEN** 系统SHALL：
- 验证status值必须为0（启用）或1（禁用）
- 更新教练的status字段
- 返回成功消息

**WHEN** 用户提交非法状态值（既非0也非1）
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息说明status值必须为0或1

**WHEN** 用户提交status=0的请求
**THEN** 系统SHALL：
- 正确处理0值（不应因为0值而验证失败）
- 使用指针类型（*int）处理status字段

#### 场景：状态过滤查询

**GIVEN** 系统中存在不同状态的教练
**WHEN** 用户请求启用教练列表
**THEN** 系统SHALL仅返回status=0的教练

**WHEN** 用户请求完整教练列表
**THEN** 系统SHALL返回所有状态的教练（包括启用和禁用）

### 需求：教练与学生关联管理

系统SHALL支持教练与学生的多对多关系管理。

#### 场景：创建教练时关联学生

**GIVEN** 系统中存在学生数据
**WHEN** 用户创建教练时提供student_ids数组
**THEN** 系统SHALL：
- 创建教练记录
- 为每个student_id在student_coach表中创建关联记录
- 使用数据库事务保证操作原子性

**WHEN** 用户创建教练时不提供student_ids
**THEN** 系统SHALL：
- 仅创建教练记录
- 不创建student_coach关联记录

#### 场景：删除教练时清理关联

**GIVEN** 教练已关联一个或多个学生
**WHEN** 用户删除教练
**THEN** 系统SHALL：
- 先删除student_coach表中的所有关联记录（coach_id匹配）
- 再删除coach表中的教练记录
- 使用事务保证操作原子性

#### 场景：学生列表查询教练信息

**GIVEN** 学生已关联一个或多个教练
**WHEN** 用户查询学生列表
**THEN** 系统SHALL：
- 通过LEFT JOIN student_coach和coach表查询教练信息
- 使用GROUP_CONCAT将多个教练名称聚合为逗号分隔的字符串
- 在响应的coach_names字段中返回

**GIVEN** 学生未关联任何教练
**WHEN** 用户查询学生列表
**THEN** coach_names字段SHALL为空字符串或null

### 需求：数据验证

系统SHALL对所有输入数据执行严格验证。

#### 场景：必填字段验证

**WHEN** 创建或更新教练时
**THEN** 系统SHALL验证以下字段必须存在且非空：
- coach_name（教练姓名）
- sex_id（性别ID）
- subject_id（学科ID）
- phone（联系电话）

#### 场景：外键约束验证

**WHEN** 创建或更新教练时
**THEN** 系统SHALL：
- 验证sex_id在sex表中存在
- 验证subject_id在subject表中存在
- 如果提供student_ids，验证每个student_id在student表中存在
- 如果外键约束失败，返回400错误和明确的错误信息

#### 场景：数据类型验证

**WHEN** 接收到教练数据
**THEN** 系统SHALL验证：
- id为正整数
- sex_id为正整数
- subject_id为正整数
- status为0或1
- phone为字符串，长度不超过20
- coach_name为字符串，长度不超过100

### 需求：响应格式规范

系统SHALL返回符合原项目规范的响应格式。

#### 场景：列表查询响应

**WHEN** 查询教练列表（GET /api/coaches）
**THEN** 响应格式SHALL为：
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

#### 场景：启用教练查询响应

**WHEN** 查询启用教练（GET /api/coaches/active）
**THEN** 响应格式SHALL为：
```json
{
  "coaches": [
    {
      "id": 1,
      "coach_name": "李老师"
    }
  ]
}
```

#### 场景：创建成功响应

**WHEN** 成功创建教练
**THEN** 系统SHALL：
- 返回201状态码
- 响应体包含教练ID：`{"id": 1, "message": "教练添加成功"}`

#### 场景：更新/删除成功响应

**WHEN** 成功更新教练信息
**THEN** 系统SHALL：
- 返回200状态码
- 响应体：`{"message": "教练信息更新成功"}`

**WHEN** 成功删除教练
**THEN** 系统SHALL：
- 返回200状态码
- 响应体：`{"message": "教练删除成功"}`

#### 场景：错误响应

**WHEN** 操作失败
**THEN** 系统SHALL：
- 返回适当的HTTP状态码（400/404/500）
- 响应体包含错误信息：`{"code": 400, "message": "错误描述"}`

### 需求：数据库事务处理

系统SHALL使用数据库事务保证多表操作的原子性。

#### 场景：创建教练的事务处理

**WHEN** 创建教练并关联学生
**THEN** 系统SHALL：
- 在单个事务中执行以下操作：
  1. 插入coach记录
  2. 批量插入student_coach记录
- 如果任何步骤失败，SHALL回滚整个事务
- 保证数据一致性

#### 场景：删除教练的事务处理

**WHEN** 删除教练
**THEN** 系统SHALL：
- 在单个事务中执行以下操作：
  1. 删除student_coach关联记录
  2. 删除coach记录
- 如果任何步骤失败，SHALL回滚整个事务

### 需求：API路由规范

系统SHALL提供以下REST API端点。

#### 场景：路由定义

系统SHALL注册以下路由：
- GET /api/coaches - 获取教练列表
- GET /api/coaches/active - 获取启用教练（MUST在/:id之前注册）
- GET /api/coaches/:id - 获取教练详情（待实现）
- POST /api/coaches - 创建教练
- PUT /api/coaches/:id - 更新教练信息
- PUT /api/coaches/:id/status - 更新教练状态
- DELETE /api/coaches/:id - 删除教练

#### 场景：路由顺序

**WHEN** 注册路由
**THEN** 系统SHALL：
- 将 /coaches/active 路由注册在 /coaches/:id 之前
- 避免active被误识别为教练ID

### 需求：DDD架构实现

系统SHALL遵循领域驱动设计（DDD）的四层架构。

#### 场景：架构分层

系统SHALL按以下层次组织代码：

**Domain层**（领域层）：
- `internal/domain/coach/entity/coach.go` - 教练实体和值对象
- `internal/domain/coach/repository/coach_repository.go` - 仓储接口定义

**Infrastructure层**（基础设施层）：
- `internal/infrastructure/persistence/mysql/coach/coach_repository_impl.go` - 仓储实现
- 使用GORM进行数据库操作

**Application层**（应用层）：
- `internal/application/service/coach/coach_service.go` - 业务逻辑服务
- `internal/application/service/coach/dto.go` - 数据传输对象

**Interface层**（接口层）：
- `internal/interfaces/http/handler/coach/coach_handler.go` - HTTP处理器
- 路由注册在 `internal/interfaces/http/router/router.go`

#### 场景：依赖方向

**WHEN** 实现各层代码
**THEN** 系统SHALL：
- 保持依赖方向：Interface → Application → Domain ← Infrastructure
- Domain层不依赖其他层
- Infrastructure层实现Domain层定义的接口
- 使用依赖注入传递仓储实现到应用服务

### 需求：数据库Schema

系统SHALL创建或验证符合规范的数据库表结构。

#### 场景：coach表结构

系统SHALL验证或创建coach表，包含以下字段：
- `id` INT PRIMARY KEY AUTO_INCREMENT
- `name` VARCHAR(100) NOT NULL - 教练姓名
- `sex_id` INT NOT NULL - 性别ID
- `subject_id` INT NOT NULL - 学科ID
- `phone` VARCHAR(20) NOT NULL - 联系电话
- `status` TINYINT DEFAULT 0 - 状态（0=启用，1=禁用）
- `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
- `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

索引SHALL包括：
- PRIMARY KEY (id)
- INDEX idx_sex_id (sex_id)
- INDEX idx_subject_id (subject_id)
- INDEX idx_status (status)

外键约束SHALL包括：
- FOREIGN KEY (sex_id) REFERENCES sex(id)
- FOREIGN KEY (subject_id) REFERENCES subject(id)

#### 场景：student_coach关联表验证

系统SHALL验证student_coach表包含coach外键：
- FOREIGN KEY (coach_id) REFERENCES coach(id) ON DELETE CASCADE
- 如果外键不存在，SHALL添加该约束

### 需求：性能优化

系统SHALL优化查询性能。

#### 场景：列表查询优化

**WHEN** 查询教练列表
**THEN** 系统SHALL：
- 使用LEFT JOIN而非N+1查询获取关联数据
- 在常用过滤字段（sex_id, subject_id, status）上创建索引

#### 场景：启用教练查询优化

**WHEN** 查询启用教练
**THEN** 系统SHALL：
- 使用WHERE status=0条件过滤
- 利用status字段的索引提升查询速度
- 仅返回必要字段（id, coach_name）

### 需求：错误处理

系统SHALL提供统一的错误处理机制。

#### 场景：统一错误响应

**WHEN** 发生错误
**THEN** 系统SHALL：
- 使用response包统一处理
- 返回适当的HTTP状态码
- 提供清晰的中文错误消息
- 记录详细的错误日志（包括堆栈跟踪）

#### 场景：业务错误处理

**WHEN** 业务规则验证失败
**THEN** 系统SHALL：
- 返回400状态码
- 提供业务相关的错误消息
- 不暴露内部实现细节

#### 场景：系统错误处理

**WHEN** 发生系统错误（数据库连接失败等）
**THEN** 系统SHALL：
- 返回500状态码
- 记录详细错误日志
- 返回通用错误消息给客户端

### 需求：日志记录

系统SHALL记录关键操作日志。

#### 场景：操作日志

**WHEN** 执行教练CRUD操作
**THEN** 系统SHALL记录：
- 操作类型（创建/更新/删除）
- 操作用户ID（从JWT token获取）
- 教练ID
- 操作时间
- 操作结果（成功/失败）

#### 场景：错误日志

**WHEN** 发生错误
**THEN** 系统SHALL记录：
- 错误类型
- 错误消息
- 堆栈跟踪
- 请求上下文（URL、方法、参数）

### 需求：字符集支持

系统SHALL支持UTF-8字符集存储中文数据。

#### 场景：数据库字符集

**WHEN** 创建数据库表
**THEN** 系统SHALL：
- 使用utf8mb4字符集
- 使用utf8mb4_unicode_ci排序规则
- 确保中文姓名正确存储和查询

#### 场景：API响应编码

**WHEN** 返回API响应
**THEN** 系统SHALL：
- 设置Content-Type为application/json; charset=utf-8
- 确保中文字符正确编码

### 需求：与学生管理模块的集成

系统SHALL确保教练管理与学生管理模块正确集成。

#### 场景：学生列表显示教练信息

**GIVEN** 学生已关联教练
**WHEN** 查询学生列表
**THEN** 系统SHALL：
- 在响应中包含coach_names字段
- 显示所有关联教练的名称（逗号分隔）

#### 场景：新增学生时选择教练

**GIVEN** 系统中存在启用的教练
**WHEN** 用户新增学生
**THEN** 系统SHALL：
- 提供启用教练列表供选择
- 允许关联多个教练

#### 场景：删除教练后学生数据更新

**GIVEN** 教练已被删除
**WHEN** 查询学生列表
**THEN** 系统SHALL：
- 不显示已删除的教练名称
- coach_names字段正确反映当前关联

## 依赖关系

### 前置依赖（已实现）
- ✅ 性别管理（sex表和API）
- ✅ 学科管理（subject表和API）
- ✅ 学生管理（student表和API）
- ✅ student_coach表

### 后置依赖（待实现）
- ❌ 课程管理 - 需要教练信息
- ❌ 排课管理 - 需要教练信息

## 验收标准

系统SHALL满足以下验收条件：

1. **功能完整性**
   - 6个API接口全部实现并通过测试
   - 所有CRUD操作正常工作

2. **数据完整性**
   - 数据库表包含正确的字段、索引、外键约束
   - 教练列表查询返回完整的关联信息（性别、学科）

3. **业务规则**
   - 创建教练支持关联学生（多对多）
   - 删除教练时正确级联删除student_coach关联
   - 教练状态切换功能正常

4. **前端兼容性**
   - 响应格式与原Python项目完全一致
   - 前端页面可以正常使用所有功能
   - 返回消息与前端判断条件匹配

5. **架构规范**
   - 代码遵循DDD四层架构
   - 各层职责清晰，依赖方向正确

6. **代码质量**
   - 代码通过lint检查
   - 关键操作有日志记录
   - 使用事务保证数据一致性

7. **集成测试**
   - 与学生管理模块集成正常
   - 学生列表正确显示教练信息
