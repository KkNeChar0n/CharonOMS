# 任务清单：实现学生管理模块

## Phase 1: 数据库Schema创建

- [x] 创建数据库迁移脚本
  - [ ] 编写 `scripts/migrations/001_create_student_tables.sql`
  - [ ] 创建student表（包含字段、索引、外键约束）
  - [ ] 创建student_coach关联表
  - [ ] 添加必要的索引（sex_id, grade_id, status）

- [x] 执行数据库迁移
  - [ ] 在开发数据库运行迁移脚本
  - [ ] 验证表结构正确性
  - [ ] 检查外键约束是否生效

- [x] 准备测试数据（可选）
  - [ ] 插入5-10条测试学生数据
  - [ ] 测试数据包含不同性别、年级、状态

## Phase 2: Domain层实现

- [x] 创建学生实体模型
  - [ ] 编写 `internal/domain/student/entity/student.go`
  - [ ] 定义Student结构体（包含所有字段）
  - [ ] 定义StudentCoach关联结构体
  - [ ] 添加GORM标签和JSON标签

- [x] 创建仓储接口
  - [ ] 编写 `internal/domain/student/repository/student_repository.go`
  - [ ] 定义StudentRepository接口
  - [ ] 方法清单：
    - [ ] GetStudentList() - 获取学生列表（含关联）
    - [ ] GetActiveStudents() - 获取启用学生
    - [ ] GetStudentByID() - 根据ID查询学生
    - [ ] CreateStudent() - 创建学生
    - [ ] UpdateStudent() - 更新学生信息
    - [ ] UpdateStudentStatus() - 更新学生状态
    - [ ] DeleteStudent() - 删除学生
    - [ ] CheckStudentHasOrders() - 检查学生是否有订单
    - [ ] AddStudentCoaches() - 添加学生教练关联
    - [ ] RemoveStudentCoaches() - 删除学生教练关联

## Phase 3: Infrastructure层实现

- [x] 创建仓储实现
  - [ ] 编写 `internal/infrastructure/persistence/mysql/student/student_repository_impl.go`
  - [ ] 实现StudentRepositoryImpl结构体
  - [ ] 实现NewStudentRepository构造函数

- [x] 实现查询方法
  - [ ] 实现GetStudentList（使用LEFT JOIN查询）
    - [ ] JOIN sex表获取性别名称
    - [ ] JOIN grade表获取年级名称
    - [ ] LEFT JOIN student_coach和coach表
    - [ ] 使用GROUP_CONCAT聚合教练名称
  - [ ] 实现GetActiveStudents（status=0过滤）
  - [ ] 实现GetStudentByID（包含关联信息）

- [x] 实现写入方法
  - [ ] 实现CreateStudent（事务处理）
    - [ ] 插入student记录
    - [ ] 批量插入student_coach关联（如果有）
  - [ ] 实现UpdateStudent
  - [ ] 实现UpdateStudentStatus
  - [ ] 实现DeleteStudent（事务处理）
    - [ ] 检查订单关联
    - [ ] 删除student_coach关联
    - [ ] 删除student记录

- [x] 实现辅助方法
  - [ ] 实现CheckStudentHasOrders
  - [ ] 实现AddStudentCoaches
  - [ ] 实现RemoveStudentCoaches

## Phase 4: Application层实现

- [x] 创建DTO定义
  - [ ] 编写 `internal/application/service/student/dto.go`
  - [ ] StudentListResponse - 学生列表响应
  - [ ] StudentDetailDTO - 学生详情
  - [ ] CreateStudentRequest - 创建学生请求
  - [ ] UpdateStudentRequest - 更新学生请求
  - [ ] UpdateStudentStatusRequest - 更新状态请求
  - [ ] ActiveStudentDTO - 简化学生信息（ID+姓名）

- [x] 创建业务服务
  - [ ] 编写 `internal/application/service/student/student_service.go`
  - [ ] 定义StudentService结构体
  - [ ] 实现NewStudentService构造函数

- [x] 实现业务方法
  - [ ] GetStudentList - 获取学生列表
  - [ ] GetActiveStudents - 获取启用学生
  - [ ] GetStudentByID - 获取学生详情
  - [ ] CreateStudent - 创建学生
    - [ ] 验证必填字段（name, sex_id, grade_id, phone）
    - [ ] 调用仓储创建学生
    - [ ] 如果有coach_ids，添加关联
  - [ ] UpdateStudent - 更新学生
    - [ ] 验证必填字段
    - [ ] 检查学生是否存在
    - [ ] 调用仓储更新
  - [ ] UpdateStudentStatus - 更新状态
    - [ ] 验证status值（0或1）
    - [ ] 调用仓储更新
  - [ ] DeleteStudent - 删除学生
    - [ ] 检查是否有关联订单
    - [ ] 如果有订单，返回错误
    - [ ] 调用仓储删除（级联删除关联）

## Phase 5: Interface层实现

- [x] 创建HTTP处理器
  - [ ] 编写 `internal/interfaces/http/handler/student/student_handler.go`
  - [ ] 定义StudentHandler结构体
  - [ ] 实现NewStudentHandler构造函数

- [x] 实现HTTP处理方法
  - [ ] GetStudents - GET /api/students
    - [ ] 调用业务服务
    - [ ] 返回格式：`{"students": [...]}`
  - [ ] GetActiveStudents - GET /api/students/active
    - [ ] 调用业务服务
    - [ ] 返回格式：`{"students": [{"id": 1, "student_name": "..."}]}`
  - [ ] CreateStudent - POST /api/students
    - [ ] 绑定JSON请求体
    - [ ] 调用业务服务
    - [ ] 返回201和学生ID
  - [ ] UpdateStudent - PUT /api/students/:id
    - [ ] 解析路径参数
    - [ ] 绑定JSON请求体
    - [ ] 调用业务服务
    - [ ] 返回成功消息
  - [ ] UpdateStudentStatus - PUT /api/students/:id/status
    - [ ] 解析路径参数
    - [ ] 绑定JSON请求体
    - [ ] 调用业务服务
    - [ ] 返回成功消息
  - [ ] DeleteStudent - DELETE /api/students/:id
    - [ ] 解析路径参数
    - [ ] 调用业务服务
    - [ ] 处理订单关联错误
    - [ ] 返回成功消息

- [x] 注册路由
  - [ ] 修改 `internal/interfaces/http/router/router.go`
  - [ ] 导入student相关包
  - [ ] 初始化StudentRepository
  - [ ] 初始化StudentService
  - [ ] 初始化StudentHandler
  - [ ] 注册students路由组
    - [ ] GET /students/active（必须在/:id之前）
    - [ ] GET /students
    - [ ] POST /students
    - [ ] PUT /students/:id
    - [ ] PUT /students/:id/status
    - [ ] DELETE /students/:id
  - [ ] 删除占位符路由

## Phase 6: 测试验证

- [x] 单元测试
  - [ ] 测试StudentService业务逻辑
  - [ ] 测试字段验证
  - [ ] 测试状态值验证
  - [ ] 测试订单关联检查

- [x] 集成测试
  - [ ] 测试创建学生（基本字段）
  - [ ] 测试创建学生并关联教练
  - [ ] 测试查询学生列表
  - [ ] 测试查询启用学生
  - [ ] 测试更新学生信息
  - [ ] 测试更新学生状态
  - [ ] 测试删除学生（无订单）
  - [ ] 测试删除学生（有订单，应失败）

- [x] API测试（使用curl或Postman）
  - [ ] POST /api/students - 创建学生
  - [ ] GET /api/students - 获取列表
  - [ ] GET /api/students/active - 获取启用学生
  - [ ] PUT /api/students/1 - 更新学生
  - [ ] PUT /api/students/1/status - 切换状态
  - [ ] DELETE /api/students/1 - 删除学生

- [x] 前端联调测试
  - [ ] 启动服务器
  - [ ] 浏览器访问 http://localhost:5001/#/students
  - [ ] 验证学生列表正常显示
  - [ ] 测试新增学生功能
  - [ ] 测试编辑学生功能
  - [ ] 测试删除学生功能
  - [ ] 测试状态切换功能
  - [ ] 测试筛选功能（ID、姓名、年级、状态）
  - [ ] 测试分页功能

## Phase 7: 文档更新

- [x] 更新项目文档
  - [ ] 更新 `REFACTORING_STATUS.md`
    - [ ] 将学生管理模块标记为100%完成
    - [ ] 更新API接口统计（+6个接口）
    - [ ] 更新数据库表统计（+2个表）
    - [ ] 更新完成度百分比

- [x] 更新README.md
  - [ ] 添加学生管理API文档
  - [ ] 列出6个接口的路径、方法、参数、响应
  - [ ] 添加使用示例

- [x] 创建API文档（可选）
  - [ ] 编写 `docs/api/student.md`
  - [ ] 详细说明每个接口的用法
  - [ ] 提供请求示例和响应示例
  - [ ] 说明错误码和异常情况

## 验收标准

- [x] 所有6个API接口实现并通过测试
- [x] 数据库表创建成功，包含正确的索引和约束
- [x] 学生列表查询返回完整的关联信息（性别、年级、教练）
- [x] 创建学生支持关联教练（多对多）
- [x] 删除学生时正确检查订单关联
- [x] 学生状态切换功能正常
- [x] 前端页面可以正常使用所有功能
- [x] 代码遵循DDD四层架构
- [x] 代码通过lint检查
- [x] 所有测试用例通过

## 注意事项

1. **路由顺序**：`/students/active` 必须在 `/students/:id` 之前注册，避免路由冲突
2. **事务处理**：创建和删除学生涉及多表操作，需要使用数据库事务
3. **响应格式**：严格按照原项目格式返回，字段名使用下划线（student_name而非studentName）
4. **教练依赖**：教练管理模块尚未实现，教练关联功能可以先实现接口，实际关联待教练模块完成后测试
5. **订单检查**：订单表未创建时，CheckStudentHasOrders可以暂时返回false，待订单模块实现后补充
6. **字符集**：确保数据库连接使用utf8mb4字符集
7. **错误处理**：使用统一的错误处理机制（response.HandleError）
8. **日志记录**：关键操作添加日志，便于调试和追踪
