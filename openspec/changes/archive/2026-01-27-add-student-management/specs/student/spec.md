# 规范：学生管理模块

## ADDED Requirements

### 需求：学生基本信息管理

系统SHALL提供学生基本信息的完整管理功能，包括创建、查询、更新和删除操作。

#### 场景：创建学生

**GIVEN** 用户已登录系统
**WHEN** 用户提交创建学生请求，包含必填字段（姓名、性别ID、年级ID、联系电话）
**THEN** 系统SHALL：
- 验证所有必填字段是否存在
- 创建新的学生记录
- 返回201状态码和学生ID
- 如果提供了coach_ids，SHALL同时创建student_coach关联记录

**WHEN** 用户提交创建学生请求，缺少必填字段
**THEN** 系统SHALL：
- 返回400错误
- 提供清晰的错误信息说明缺少哪个字段

#### 场景：查询学生列表

**GIVEN** 系统中存在学生数据
**WHEN** 用户请求学生列表
**THEN** 系统SHALL：
- 返回所有学生的完整信息
- 包含关联的性别名称（通过sex_id关联sex表）
- 包含关联的年级名称（通过grade_id关联grade表）
- 包含关联的教练名称（通过student_coach表和coach表，使用逗号分隔多个教练）
- 响应格式为：`{"students": [...]}`
- 每个学生对象包含字段：id, student_name, sex_id, sex, grade_id, grade, phone, status, coach_names

#### 场景：查询启用状态的学生

**GIVEN** 系统中存在启用和禁用的学生
**WHEN** 用户请求启用学生列表（GET /api/students/active）
**THEN** 系统SHALL：
- 仅返回status=0的学生
- 返回简化格式：`{"students": [{"id": 1, "student_name": "..."}]}`
- 不包含禁用（status=1）的学生

#### 场景：更新学生信息

**GIVEN** 系统中存在指定ID的学生
**WHEN** 用户提交更新学生请求
**THEN** 系统SHALL：
- 验证学生ID是否存在
- 验证必填字段是否完整
- 更新学生基本信息
- 返回成功消息

**WHEN** 用户尝试更新不存在的学生
**THEN** 系统SHALL：
- 返回404错误
- 提供错误信息说明学生不存在

#### 场景：删除学生

**GIVEN** 系统中存在指定ID的学生
**WHEN** 用户提交删除学生请求
**THEN** 系统SHALL：
- 检查学生是否有关联订单
- 如果有关联订单，SHALL返回错误并阻止删除
- 如果没有关联订单，SHALL：
  - 删除student_coach表中的所有关联记录
  - 删除student表中的学生记录
  - 使用数据库事务保证操作原子性
  - 返回成功消息

**WHEN** 用户尝试删除有订单的学生
**THEN** 系统SHALL：
- 返回错误信息："无法删除，该学生存在关联订单"
- 不执行删除操作

### 需求：学生状态管理

系统SHALL提供学生启用/禁用状态的管理功能。

#### 场景：更新学生状态

**GIVEN** 系统中存在指定ID的学生
**WHEN** 用户提交更新状态请求（PUT /api/students/:id/status）
**THEN** 系统SHALL：
- 验证status值必须为0（启用）或1（禁用）
- 更新学生的status字段
- 返回成功消息

**WHEN** 用户提交非法状态值（既非0也非1）
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息说明status值必须为0或1

#### 场景：状态过滤查询

**GIVEN** 系统中存在不同状态的学生
**WHEN** 用户请求启用学生列表
**THEN** 系统SHALL仅返回status=0的学生

**WHEN** 用户请求完整学生列表
**THEN** 系统SHALL返回所有状态的学生（包括启用和禁用）

### 需求：学生与教练关联管理

系统SHALL支持学生与教练的多对多关系管理。

#### 场景：创建学生时关联教练

**GIVEN** 系统中存在教练数据
**WHEN** 用户创建学生时提供coach_ids数组
**THEN** 系统SHALL：
- 创建学生记录
- 为每个coach_id在student_coach表中创建关联记录
- 使用数据库事务保证操作原子性

**WHEN** 用户创建学生时不提供coach_ids
**THEN** 系统SHALL：
- 仅创建学生记录
- 不创建student_coach关联记录

#### 场景：查询学生的教练信息

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

**WHEN** 创建或更新学生时
**THEN** 系统SHALL验证以下字段必须存在且非空：
- name（学生姓名）
- sex_id（性别ID）
- grade_id（年级ID）
- phone（联系电话）

#### 场景：外键约束验证

**WHEN** 创建或更新学生时
**THEN** 系统SHALL：
- 验证sex_id在sex表中存在
- 验证grade_id在grade表中存在
- 如果提供coach_ids，验证每个coach_id在coach表中存在
- 如果外键约束失败，返回400错误和明确的错误信息

#### 场景：数据类型验证

**WHEN** 接收到学生数据
**THEN** 系统SHALL验证：
- id为正整数
- sex_id为正整数
- grade_id为正整数
- status为0或1
- phone为字符串，长度不超过20
- name为字符串，长度不超过100

### 需求：响应格式规范

系统SHALL返回符合原项目规范的响应格式。

#### 场景：列表查询响应

**WHEN** 查询学生列表（GET /api/students）
**THEN** 响应格式SHALL为：
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

#### 场景：启用学生查询响应

**WHEN** 查询启用学生（GET /api/students/active）
**THEN** 响应格式SHALL为：
```json
{
  "students": [
    {
      "id": 1,
      "student_name": "张三"
    }
  ]
}
```

#### 场景：创建成功响应

**WHEN** 成功创建学生
**THEN** 系统SHALL：
- 返回201状态码
- 响应体包含学生ID：`{"id": 1, "message": "创建成功"}`

#### 场景：更新/删除成功响应

**WHEN** 成功更新或删除学生
**THEN** 系统SHALL：
- 返回200状态码
- 响应体包含成功消息：`{"message": "操作成功"}`

#### 场景：错误响应

**WHEN** 操作失败
**THEN** 系统SHALL：
- 返回适当的HTTP状态码（400/404/500）
- 响应体包含错误信息：`{"error": "错误描述"}`

### 需求：数据库事务处理

系统SHALL使用数据库事务保证多表操作的原子性。

#### 场景：创建学生的事务处理

**WHEN** 创建学生并关联教练
**THEN** 系统SHALL：
- 在单个事务中执行以下操作：
  1. 插入student记录
  2. 批量插入student_coach记录
- 如果任何步骤失败，SHALL回滚整个事务
- 保证数据一致性

#### 场景：删除学生的事务处理

**WHEN** 删除学生
**THEN** 系统SHALL：
- 在单个事务中执行以下操作：
  1. 检查是否有关联订单
  2. 删除student_coach关联记录
  3. 删除student记录
- 如果任何步骤失败，SHALL回滚整个事务

### 需求：API路由规范

系统SHALL提供以下REST API端点。

#### 场景：路由定义

系统SHALL注册以下路由：
- GET /api/students - 获取学生列表
- GET /api/students/active - 获取启用学生（MUST在/:id之前注册）
- GET /api/students/:id - 获取学生详情（待实现）
- POST /api/students - 创建学生
- PUT /api/students/:id - 更新学生信息
- PUT /api/students/:id/status - 更新学生状态
- DELETE /api/students/:id - 删除学生

#### 场景：路由顺序

**WHEN** 注册路由
**THEN** 系统SHALL：
- 将 /students/active 路由注册在 /students/:id 之前
- 避免active被误识别为学生ID

### 需求：DDD架构实现

系统SHALL遵循领域驱动设计（DDD）的四层架构。

#### 场景：架构分层

系统SHALL按以下层次组织代码：

**Domain层**（领域层）：
- `internal/domain/student/entity/student.go` - 学生实体和值对象
- `internal/domain/student/repository/student_repository.go` - 仓储接口定义

**Infrastructure层**（基础设施层）：
- `internal/infrastructure/persistence/mysql/student/student_repository_impl.go` - 仓储实现
- 使用GORM进行数据库操作

**Application层**（应用层）：
- `internal/application/service/student/student_service.go` - 业务逻辑服务
- `internal/application/service/student/dto.go` - 数据传输对象

**Interface层**（接口层）：
- `internal/interfaces/http/handler/student/student_handler.go` - HTTP处理器
- 路由注册在 `internal/interfaces/http/router/router.go`

#### 场景：依赖方向

**WHEN** 实现各层代码
**THEN** 系统SHALL：
- 保持依赖方向：Interface → Application → Domain ← Infrastructure
- Domain层不依赖其他层
- Infrastructure层实现Domain层定义的接口
- 使用依赖注入传递仓储实现到应用服务

### 需求：数据库Schema

系统SHALL创建符合规范的数据库表结构。

#### 场景：student表结构

系统SHALL创建student表，包含以下字段：
- `id` INT PRIMARY KEY AUTO_INCREMENT
- `name` VARCHAR(100) NOT NULL - 学生姓名
- `sex_id` INT NOT NULL - 性别ID
- `grade_id` INT NOT NULL - 年级ID
- `phone` VARCHAR(20) NOT NULL - 联系电话
- `status` TINYINT DEFAULT 0 - 状态（0=启用，1=禁用）
- `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
- `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

索引SHALL包括：
- PRIMARY KEY (id)
- INDEX idx_sex_id (sex_id)
- INDEX idx_grade_id (grade_id)
- INDEX idx_status (status)

外键约束SHALL包括：
- FOREIGN KEY (sex_id) REFERENCES sex(id)
- FOREIGN KEY (grade_id) REFERENCES grade(id)

#### 场景：student_coach关联表结构

系统SHALL创建student_coach表，包含以下字段：
- `id` INT PRIMARY KEY AUTO_INCREMENT
- `student_id` INT NOT NULL - 学生ID
- `coach_id` INT NOT NULL - 教练ID
- `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP

索引SHALL包括：
- PRIMARY KEY (id)
- UNIQUE KEY uk_student_coach (student_id, coach_id) - 防止重复关联
- INDEX idx_student_id (student_id)
- INDEX idx_coach_id (coach_id)

外键约束SHALL包括：
- FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE
- FOREIGN KEY (coach_id) REFERENCES coach(id) ON DELETE CASCADE

### 需求：性能优化

系统SHALL优化查询性能。

#### 场景：列表查询优化

**WHEN** 查询学生列表
**THEN** 系统SHALL：
- 使用LEFT JOIN而非N+1查询获取关联数据
- 使用GROUP_CONCAT聚合教练名称，避免后处理
- 在常用过滤字段（sex_id, grade_id, status）上创建索引

#### 场景：启用学生查询优化

**WHEN** 查询启用学生
**THEN** 系统SHALL：
- 使用WHERE status=0条件过滤
- 利用status字段的索引提升查询速度
- 仅返回必要字段（id, student_name）

### 需求：错误处理

系统SHALL提供统一的错误处理机制。

#### 场景：统一错误响应

**WHEN** 发生错误
**THEN** 系统SHALL：
- 使用response.HandleError统一处理
- 返回适当的HTTP状态码
- 提供清晰的中文错误消息
- 记录详细的错误日志（包括堆栈跟踪）

#### 场景：业务错误处理

**WHEN** 业务规则验证失败（如删除有订单的学生）
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

**WHEN** 执行学生CRUD操作
**THEN** 系统SHALL记录：
- 操作类型（创建/更新/删除）
- 操作用户ID（从JWT token获取）
- 学生ID
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

## 依赖关系

### 前置依赖（已实现）
- ✅ 性别管理（sex表和API）
- ✅ 年级管理（grade表和API）

### 可选依赖（待实现）
- ⚠️ 教练管理（coach表和API） - 学生可以不关联教练，但关联功能需要教练表存在

### 后置依赖（待实现）
- ❌ 订单管理（order表） - 删除学生前需要检查订单关联

## 验收标准

系统SHALL满足以下验收条件：

1. **功能完整性**
   - 6个API接口全部实现并通过测试
   - 所有CRUD操作正常工作

2. **数据完整性**
   - 数据库表包含正确的字段、索引、外键约束
   - 学生列表查询返回完整的关联信息（性别、年级、教练）

3. **业务规则**
   - 创建学生支持关联教练（多对多）
   - 删除学生时正确检查订单关联
   - 学生状态切换功能正常

4. **前端兼容性**
   - 响应格式与原Python项目完全一致
   - 前端页面可以正常使用所有功能

5. **架构规范**
   - 代码遵循DDD四层架构
   - 各层职责清晰，依赖方向正确

6. **代码质量**
   - 代码通过lint检查
   - 所有测试用例通过
   - 关键操作有日志记录
