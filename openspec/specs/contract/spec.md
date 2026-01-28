# 规范：合同管理模块

## 新增需求

### 需求：合同基本信息管理

系统SHALL提供合同基本信息的完整管理功能，包括创建、查询操作。

#### 场景：创建合同

**GIVEN** 用户已登录系统
**WHEN** 用户提交创建合同请求，包含必填字段（合同名称、学生ID、合同类型、签署形式、合同金额）
**THEN** 系统SHALL：
- 验证所有必填字段是否存在
- 验证合同类型必须为0（首报）或1（续报）
- 验证签署形式必须为0（线上签署）或1（线下签署）
- 验证合同金额为有效的数值
- 创建新的合同记录
- 设置status为0（待审核）
- 设置payment_status为0（未付款）
- 从JWT token获取当前用户名作为initiator（发起人）
- 返回201状态码和合同ID
- 返回消息："合同新增成功"

**WHEN** 用户提交创建合同请求，缺少必填字段
**THEN** 系统SHALL：
- 返回400错误
- 提供清晰的错误信息说明缺少哪个字段

**WHEN** 用户提交创建合同请求，合同类型或签署形式值非法
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息说明合法值范围

#### 场景：查询合同列表

**GIVEN** 系统中存在合同数据
**WHEN** 用户请求合同列表
**THEN** 系统SHALL：
- 返回所有合同的完整信息
- 包含关联的学生姓名（通过student_id关联student表）
- 按创建时间倒序排列
- 响应格式为：`{"contracts": [...]}`
- 每个合同对象包含字段：id, name, student_id, student_name, type, signature_form, contract_amount, signatory, initiating_party, initiator, status, payment_status, termination_agreement, create_time

**GIVEN** 系统中不存在任何合同
**WHEN** 用户请求合同列表
**THEN** 系统SHALL返回空数组：`{"contracts": []}`

#### 场景：查询合同详情

**GIVEN** 系统中存在指定ID的合同
**WHEN** 用户请求合同详情
**THEN** 系统SHALL：
- 返回该合同的完整信息
- 包含关联的学生姓名
- 所有字段都应返回

**WHEN** 用户请求不存在的合同ID
**THEN** 系统SHALL：
- 返回404错误
- 提供错误信息说明合同不存在

### 需求：合同状态管理

系统SHALL提供合同状态流转的管理功能，支持撤销和中止合作操作。

#### 场景：撤销合同

**GIVEN** 系统中存在status为0（待审核）的合同
**WHEN** 用户提交撤销合同请求
**THEN** 系统SHALL：
- 验证合同是否存在
- 验证当前status必须为0
- 更新status为98（已作废）
- 返回成功消息

**WHEN** 用户尝试撤销status不为0的合同
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息："只有待审核状态的合同可以撤销"

**WHEN** 用户尝试撤销不存在的合同
**THEN** 系统SHALL：
- 返回404错误
- 提供错误信息说明合同不存在

#### 场景：中止合作

**GIVEN** 系统中存在status为50（已通过）的合同
**WHEN** 用户提交中止合作请求，包含termination_agreement（中止协议文件）
**THEN** 系统SHALL：
- 验证合同是否存在
- 验证当前status必须为50
- 验证termination_agreement不为空
- 更新status为99（协议中止）
- 保存termination_agreement文件路径
- 返回成功消息

**WHEN** 用户尝试中止status不为50的合同
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息："只有已通过状态的合同可以中止"

**WHEN** 用户提交中止合作请求，未提供termination_agreement
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息说明必须上传中止协议文件

**WHEN** 用户尝试中止不存在的合同
**THEN** 系统SHALL：
- 返回404错误
- 提供错误信息说明合同不存在

### 需求：合同与学生关联管理

系统SHALL支持合同与学生的关联关系管理。

#### 场景：创建合同时关联学生

**GIVEN** 系统中存在学生数据
**WHEN** 用户创建合同时提供有效的student_id
**THEN** 系统SHALL：
- 验证student_id在student表中存在
- 创建合同记录并关联该学生

**WHEN** 用户创建合同时提供不存在的student_id
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息说明学生不存在

#### 场景：查询合同时显示学生信息

**GIVEN** 合同已关联学生
**WHEN** 用户查询合同列表或详情
**THEN** 系统SHALL：
- 通过LEFT JOIN student表查询学生姓名
- 在响应的student_name字段中返回学生姓名

**GIVEN** 合同关联的学生已被删除
**WHEN** 用户查询合同列表或详情
**THEN** student_name字段SHALL为空或null

### 需求：数据验证

系统SHALL对所有输入数据执行严格验证。

#### 场景：必填字段验证

**WHEN** 创建合同时
**THEN** 系统SHALL验证以下字段必须存在且非空：
- name（合同名称）
- student_id（学生ID）
- type（合同类型）
- signature_form（签署形式）
- contract_amount（合同金额）

#### 场景：外键约束验证

**WHEN** 创建合同时
**THEN** 系统SHALL：
- 验证student_id在student表中存在
- 如果外键约束失败，返回400错误和明确的错误信息

#### 场景：数据类型验证

**WHEN** 接收到合同数据
**THEN** 系统SHALL验证：
- id为正整数
- student_id为正整数
- type为0或1
- signature_form为0或1
- contract_amount为有效的DECIMAL(10,2)格式
- status为0、50、98或99之一
- payment_status为0、10或30之一
- name为字符串，长度不超过200
- signatory为字符串，长度不超过100
- initiating_party为字符串，长度不超过100（可为空）
- initiator为字符串，长度不超过50
- termination_agreement为字符串，长度不超过255（可为空）

### 需求：响应格式规范

系统SHALL返回符合原项目规范的响应格式。

#### 场景：列表查询响应

**WHEN** 查询合同列表（GET /api/contracts）
**THEN** 响应格式SHALL为：
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

#### 场景：详情查询响应

**WHEN** 查询合同详情（GET /api/contracts/:id）
**THEN** 响应格式SHALL为：
```json
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
```

#### 场景：创建成功响应

**WHEN** 成功创建合同
**THEN** 系统SHALL：
- 返回201状态码
- 响应体包含合同ID和消息：`{"id": 1, "message": "合同新增成功"}`

#### 场景：撤销/中止成功响应

**WHEN** 成功撤销或中止合同
**THEN** 系统SHALL：
- 返回200状态码
- 响应体包含成功消息：`{"message": "操作成功"}`

#### 场景：错误响应

**WHEN** 操作失败
**THEN** 系统SHALL：
- 返回适当的HTTP状态码（400/404/500）
- 响应体包含错误信息：`{"code": 400, "message": "错误描述"}`

### 需求：合同状态流转规则

系统SHALL严格按照业务规则管理合同状态流转。

#### 场景：状态流转图

系统SHALL支持以下状态流转：
- 新建合同 → 待审核(0)
- 待审核(0) → 已作废(98)：通过撤销操作
- 已通过(50) → 协议中止(99)：通过中止合作操作

**非法状态流转**：
- 已作废(98) 不能撤销
- 协议中止(99) 不能撤销
- 待审核(0) 不能中止
- 已作废(98) 不能中止
- 协议中止(99) 不能中止

#### 场景：状态验证

**WHEN** 执行撤销操作
**THEN** 系统SHALL验证当前status必须为0

**WHEN** 执行中止操作
**THEN** 系统SHALL验证当前status必须为50

**WHEN** 状态验证失败
**THEN** 系统SHALL返回400错误并提供明确的错误信息

### 需求：付款状态管理

系统SHALL保留付款状态字段，但当前阶段不实现更新逻辑。

#### 场景：创建合同时初始化付款状态

**WHEN** 创建新合同
**THEN** 系统SHALL：
- 设置payment_status为0（未付款）
- 该字段在合同创建后由订单收款模块自动更新

#### 场景：查询时返回付款状态

**WHEN** 查询合同列表或详情
**THEN** 系统SHALL：
- 在响应中包含payment_status字段
- 返回当前数据库中的值

#### 场景：付款状态字段说明

系统SHALL使用以下付款状态值：
- 0 = 未付款
- 10 = 部分付款
- 30 = 已付款

**注意**：当前模块不提供更新payment_status的接口，该字段由未来的订单收款模块自动更新。

### 需求：文件上传管理

系统SHALL支持中止协议文件的上传和存储。

#### 场景：上传中止协议文件

**WHEN** 用户中止合作时上传termination_agreement文件
**THEN** 系统SHALL：
- 接收上传的文件
- 验证文件类型（允许的文件类型：pdf, doc, docx, jpg, png）
- 验证文件大小（最大10MB）
- 保存文件到服务器指定目录
- 生成唯一的文件名（避免冲突）
- 返回文件路径
- 将路径存储到contract表的termination_agreement字段

**WHEN** 上传的文件类型不符合要求
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息说明允许的文件类型

**WHEN** 上传的文件大小超过限制
**THEN** 系统SHALL：
- 返回400错误
- 提供错误信息说明文件大小限制

#### 场景：查询时返回文件路径

**WHEN** 查询合同详情且该合同已中止
**THEN** 系统SHALL：
- 在termination_agreement字段中返回文件路径
- 前端可使用该路径下载或预览文件

### 需求：API路由规范

系统SHALL提供以下REST API端点。

#### 场景：路由定义

系统SHALL注册以下路由：
- GET /api/contracts - 获取合同列表
- GET /api/contracts/:id - 获取合同详情
- POST /api/contracts - 创建合同
- PUT /api/contracts/:id/revoke - 撤销合同
- PUT /api/contracts/:id/terminate - 中止合作

#### 场景：路由顺序

**WHEN** 注册路由
**THEN** 系统SHALL：
- 将具体路径（如 /contracts/:id/revoke）注册在通用路径（如 /contracts/:id）之前
- 避免路由匹配冲突

### 需求：DDD架构实现

系统SHALL遵循领域驱动设计（DDD）的四层架构。

#### 场景：架构分层

系统SHALL按以下层次组织代码：

**Domain层**（领域层）：
- `internal/domain/contract/entity/contract.go` - 合同实体和值对象
- `internal/domain/contract/repository/contract_repository.go` - 仓储接口定义

**Infrastructure层**（基础设施层）：
- `internal/infrastructure/persistence/mysql/contract/contract_repository_impl.go` - 仓储实现
- 使用GORM进行数据库操作

**Application层**（应用层）：
- `internal/application/service/contract/contract_service.go` - 业务逻辑服务
- `internal/application/service/contract/dto.go` - 数据传输对象

**Interface层**（接口层）：
- `internal/interfaces/http/handler/contract/contract_handler.go` - HTTP处理器
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

#### 场景：contract表结构

系统SHALL验证或创建contract表，包含以下字段：
- `id` INT PRIMARY KEY AUTO_INCREMENT
- `name` VARCHAR(200) NOT NULL - 合同名称
- `student_id` INT NOT NULL - 学生ID
- `type` TINYINT NOT NULL - 合同类型（0=首报，1=续报）
- `signature_form` TINYINT NOT NULL - 签署形式（0=线上，1=线下）
- `contract_amount` DECIMAL(10,2) NOT NULL - 合同金额
- `signatory` VARCHAR(100) - 签署方
- `initiating_party` VARCHAR(100) - 发起方（可为空）
- `initiator` VARCHAR(50) - 发起人
- `status` TINYINT DEFAULT 0 - 合同状态（0=待审核，50=已通过，98=已作废，99=协议中止）
- `payment_status` TINYINT DEFAULT 0 - 付款状态（0=未付款，10=部分付款，30=已付款）
- `termination_agreement` VARCHAR(255) - 中止协议文件路径
- `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP

索引SHALL包括：
- PRIMARY KEY (id)
- INDEX idx_student_id (student_id)
- INDEX idx_status (status)
- INDEX idx_payment_status (payment_status)
- INDEX idx_create_time (create_time)

外键约束SHALL包括：
- FOREIGN KEY (student_id) REFERENCES student(id)

### 需求：性能优化

系统SHALL优化查询性能。

#### 场景：列表查询优化

**WHEN** 查询合同列表
**THEN** 系统SHALL：
- 使用LEFT JOIN而非N+1查询获取学生信息
- 在常用过滤字段（student_id, status, payment_status）上创建索引
- 按create_time倒序排列，利用索引提升排序速度

#### 场景：查询条件过滤

**WHEN** 前端提供过滤条件
**THEN** 系统SHALL支持以下过滤维度：
- ID 精确匹配
- student_id 精确匹配
- student_name 模糊匹配
- type 精确匹配
- status 精确匹配
- payment_status 精确匹配

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

**WHEN** 执行合同CRUD操作
**THEN** 系统SHALL记录：
- 操作类型（创建/撤销/中止）
- 操作用户ID（从JWT token获取）
- 合同ID
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
- 确保中文合同名称和签署方名称正确存储和查询

#### 场景：API响应编码

**WHEN** 返回API响应
**THEN** 系统SHALL：
- 设置Content-Type为application/json; charset=utf-8
- 确保中文字符正确编码

### 需求：用户身份识别

系统SHALL从JWT token中识别当前操作用户。

#### 场景：获取当前用户

**WHEN** 创建合同
**THEN** 系统SHALL：
- 从请求的JWT token中提取用户名
- 将用户名作为initiator（发起人）字段值
- 如果无法获取用户名，使用空字符串或默认值

#### 场景：JWT中间件保护

**WHEN** 访问合同管理API
**THEN** 系统SHALL：
- 要求请求包含有效的JWT token
- 验证token的有效性
- 拒绝未认证的请求

## 依赖关系

### 前置依赖（已实现）
- ✅ 学生管理（student表和API）
- ✅ 认证系统（JWT token生成和验证）

### 后置依赖（待实现）
- ❌ 订单管理 - 需要合同ID字段（contract_id）
- ❌ 收款管理 - 需要合同信息进行收款关联
- ❌ 审批流程 - 需要审批合同状态从0变为50

## 验收标准

系统SHALL满足以下验收条件：

1. **功能完整性**
   - 5个API接口全部实现并通过测试
   - 所有CRUD和状态流转操作正常工作

2. **数据完整性**
   - 数据库表包含正确的字段、索引、外键约束
   - 合同列表查询返回完整的学生信息

3. **业务规则**
   - 创建合同正确初始化状态和发起人
   - 撤销合同仅在待审核状态时允许
   - 中止合作仅在已通过状态时允许
   - 状态流转验证严格执行

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
   - 错误处理统一规范

7. **文件上传**
   - 中止协议文件上传功能正常
   - 文件类型和大小验证有效
   - 文件路径正确存储和返回
