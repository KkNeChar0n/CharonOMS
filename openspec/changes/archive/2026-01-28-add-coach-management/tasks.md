# 任务清单：实现教练管理模�?
## Phase 1: 数据库Schema验证和准�?
- [x] 验证coach表是否存�?  - [x] 检查数据库中是否已有coach�?  - [x] 如不存在，准备创建脚�?  - [x] 验证字段完整性（id, name, sex_id, subject_id, phone, status, created_at, updated_at�?
- [x] 验证student_coach表的外键约束
  - [x] 检查是否存在coach_id外键约束
  - [x] 如不存在，添加外键约�?  - [x] 验证外键的级联删除设置（ON DELETE CASCADE�?
- [x] 准备测试数据（可选）
  - [x] 插入3-5条测试教练数�?  - [x] 测试数据包含不同性别、学科、状�?  - [x] 创建部分教练与学生的关联数据

## Phase 2: Domain层实�?
- [x] 创建教练实体模型
  - [x] 编写 `internal/domain/coach/entity/coach.go`
  - [x] 定义Coach结构体（包含所有字段）
  - [x] 添加GORM标签和JSON标签
  - [x] 定义TableName()方法

- [x] 创建仓储接口
  - [x] 编写 `internal/domain/coach/repository/coach_repository.go`
  - [x] 定义CoachRepository接口
  - [x] 方法清单�?    - [x] GetCoachList() - 获取教练列表（含关联�?    - [x] GetActiveCoaches() - 获取启用教练
    - [x] GetCoachByID() - 根据ID查询教练
    - [x] CreateCoach() - 创建教练
    - [x] UpdateCoach() - 更新教练信息
    - [x] UpdateCoachStatus() - 更新教练状�?    - [x] DeleteCoach() - 删除教练
    - [x] AddCoachStudents() - 添加教练学生关联
    - [x] RemoveCoachStudents() - 删除教练学生关联

## Phase 3: Infrastructure层实�?
- [x] 创建仓储实现
  - [x] 编写 `internal/infrastructure/persistence/mysql/coach/coach_repository_impl.go`
  - [x] 实现CoachRepositoryImpl结构�?  - [x] 实现NewCoachRepository构造函�?
- [x] 实现查询方法
  - [x] 实现GetCoachList（使用LEFT JOIN查询�?    - [x] JOIN sex表获取性别名称
    - [x] JOIN subject表获取学科名�?    - [x] 返回完整的教练信息列�?  - [x] 实现GetActiveCoaches（status=0过滤�?    - [x] 仅返回id和name字段
  - [x] 实现GetCoachByID（包含关联信息）

- [x] 实现写入方法
  - [x] 实现CreateCoach（事务处理）
    - [x] 插入coach记录
    - [x] 批量插入student_coach关联（如果有�?  - [x] 实现UpdateCoach
  - [x] 实现UpdateCoachStatus
  - [x] 实现DeleteCoach（事务处理）
    - [x] 删除student_coach关联
    - [x] 删除coach记录

- [x] 实现辅助方法
  - [x] 实现AddCoachStudents
  - [x] 实现RemoveCoachStudents
  - [x] 实现ValidateCoachExists（验证教练是否存在）

## Phase 4: Application层实�?
- [x] 创建DTO定义
  - [x] 编写 `internal/application/service/coach/dto.go`
  - [x] CoachListResponse - 教练列表响应
  - [x] CoachDetailDTO - 教练详情
  - [x] CreateCoachRequest - 创建教练请求
  - [x] UpdateCoachRequest - 更新教练请求
  - [x] UpdateCoachStatusRequest - 更新状态请�?  - [x] ActiveCoachDTO - 简化教练信息（ID+姓名�?
- [x] 创建业务服务
  - [x] 编写 `internal/application/service/coach/coach_service.go`
  - [x] 定义CoachService结构�?  - [x] 实现NewCoachService构造函�?
- [x] 实现业务方法
  - [x] GetCoachList - 获取教练列表
  - [x] GetActiveCoaches - 获取启用教练
  - [x] GetCoachByID - 获取教练详情
  - [x] CreateCoach - 创建教练
    - [x] 验证必填字段（name, sex_id, subject_id, phone�?    - [x] 调用仓储创建教练
    - [x] 如果有student_ids，添加关�?  - [x] UpdateCoach - 更新教练
    - [x] 验证必填字段
    - [x] 检查教练是否存�?    - [x] 调用仓储更新
  - [x] UpdateCoachStatus - 更新状�?    - [x] 验证status值（0�?�?    - [x] 调用仓储更新
  - [x] DeleteCoach - 删除教练
    - [x] 调用仓储删除（级联删除关联）

## Phase 5: Interface层实�?
- [x] 创建HTTP处理�?  - [x] 编写 `internal/interfaces/http/handler/coach/coach_handler.go`
  - [x] 定义CoachHandler结构�?  - [x] 实现NewCoachHandler构造函�?
- [x] 实现HTTP处理方法
  - [x] GetCoaches - GET /api/coaches
    - [x] 调用业务服务
    - [x] 返回格式：`{"coaches": [...]}`
  - [x] GetActiveCoaches - GET /api/coaches/active
    - [x] 调用业务服务
    - [x] 返回格式：`{"coaches": [{"id": 1, "coach_name": "..."}]}`
  - [x] CreateCoach - POST /api/coaches
    - [x] 绑定JSON请求�?    - [x] 调用业务服务
    - [x] 返回201和教练ID
    - [x] 返回消息�?教练添加成功"
  - [x] UpdateCoach - PUT /api/coaches/:id
    - [x] 解析路径参数
    - [x] 绑定JSON请求�?    - [x] 调用业务服务
    - [x] 返回消息�?教练信息更新成功"
  - [x] UpdateCoachStatus - PUT /api/coaches/:id/status
    - [x] 解析路径参数
    - [x] 绑定JSON请求体（使用指针类型*int处理0值）
    - [x] 调用业务服务
    - [x] 返回成功消息
  - [x] DeleteCoach - DELETE /api/coaches/:id
    - [x] 解析路径参数
    - [x] 调用业务服务
    - [x] 返回消息�?教练删除成功"

- [x] 更新路由注册
  - [x] 修改 `internal/interfaces/http/router/router.go`
  - [x] 导入coach相关�?  - [x] 初始化CoachRepository
  - [x] 初始化CoachService
  - [x] 初始化CoachHandler
  - [x] 替换占位符路由为真实路由
    - [x] GET /coaches/active（必须在/:id之前�?    - [x] GET /coaches
    - [x] POST /coaches
    - [x] PUT /coaches/:id
    - [x] PUT /coaches/:id/status
    - [x] DELETE /coaches/:id

## Phase 6: 测试验证

- [x] 单元测试（可选）
  - [x] 测试CoachService业务逻辑
  - [x] 测试字段验证
  - [x] 测试状态值验�?
- [x] 集成测试
  - [x] 测试创建教练（基本字段）
  - [x] 测试创建教练并关联学�?  - [x] 测试查询教练列表
  - [x] 测试查询启用教练
  - [x] 测试更新教练信息
  - [x] 测试更新教练状态（包括status=0的情况）
  - [x] 测试删除教练

- [x] API测试（使用curl或Postman�?  - [x] POST /api/coaches - 创建教练
  - [x] GET /api/coaches - 获取列表
  - [x] GET /api/coaches/active - 获取启用教练
  - [x] PUT /api/coaches/1 - 更新教练
  - [x] PUT /api/coaches/1/status - 切换状�?  - [x] DELETE /api/coaches/1 - 删除教练

- [x] 前端联调测试
  - [x] 启动服务�?  - [x] 浏览器访�?http://localhost:5001/#/coaches
  - [x] 验证教练列表正常显示
  - [x] 测试新增教练功能
    - [x] 验证可以选择学科
    - [x] 验证可以选择关联学生
    - [x] 验证弹窗自动关闭
  - [x] 测试编辑教练功能
    - [x] 验证页面自动刷新
  - [x] 测试删除教练功能
    - [x] 验证页面自动刷新
  - [x] 测试状态切换功�?    - [x] 验证启用/禁用按钮正常
  - [x] 测试筛选功能（ID、姓名、性别、学科、状态）
  - [x] 测试分页功能

- [x] 与学生管理模块的集成测试
  - [x] 验证学生列表中教练名称正常显�?  - [x] 验证新增学生时可以选择教练
  - [x] 验证删除教练后学生列表中教练信息更新

## Phase 7: 文档更新

- [x] 更新项目文档
  - [x] 更新 `REFACTORING_STATUS.md`
    - [x] 将教练管理模块标记为100%完成
    - [x] 更新API接口统计�?6个接口）
    - [x] 更新完成度百分比�?5% �?30%�?
- [x] 更新README.md（可选）
  - [x] 添加教练管理API文档
  - [x] 列出6个接口的路径、方法、参数、响�?  - [x] 添加使用示例

- [x] 创建API文档（可选）
  - [x] 编写 `docs/api/coach.md`
  - [x] 详细说明每个接口的用�?  - [x] 提供请求示例和响应示�?  - [x] 说明错误码和异常情况

## 验收标准

- [x] 所�?个API接口实现并通过测试
- [x] 数据库表验证完成，包含正确的索引和约�?- [x] 教练列表查询返回完整的关联信息（性别、学科）
- [x] 创建教练支持关联学生（多对多�?- [x] 删除教练时正确级联删除student_coach关联
- [x] 教练状态切换功能正�?- [x] 前端页面可以正常使用所有功�?- [x] 代码遵循DDD四层架构
- [x] 代码通过lint检�?- [x] 与学生管理模块集成正�?
## 注意事项

1. **路由顺序**：`/coaches/active` 必须�?`/coaches/:id` 之前注册，避免路由冲�?2. **事务处理**：创建和删除教练涉及多表操作，需要使用数据库事务
3. **响应格式**：严格按照原项目格式返回，字段名使用下划线（coach_name而非coachName�?4. **字段名匹�?*：前端使用coach_name，DTO的JSON标签需要匹�?5. **状态字段验�?*：使用指针类型`*int`处理status=0的验证问题（参考学生管理模块）
6. **student_coach表复�?*：该表已存在，确保coach_id外键约束完整
7. **字符�?*：确保数据库连接使用utf8mb4字符�?8. **错误处理**：使用统一的错误处理机制（response包）
9. **日志记录**：关键操作添加日志，便于调试和追�?10. **前端兼容**：返回消息需要与前端判断条件完全匹配
