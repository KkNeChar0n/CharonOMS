# Specification: User Authentication

用户认证模块，提供基于JWT的登录、登出和用户信息管理功能。

## Requirements

### Requirement: 用户登录
系统 SHALL 允许用户通过用户名和密码进行登录认证。

#### Scenario: 登录成功
- **WHEN** 用户提供有效的用户名和密码
- **THEN** 系统返回JWT token和用户信息（用户名、是否超级管理员）
- **AND** 响应格式为 `{"code": 0, "message": "success", "data": {"token": "...", "username": "...", "is_super_admin": true/false}}`

#### Scenario: 用户名不存在
- **WHEN** 用户提供不存在的用户名
- **THEN** 系统返回401错误
- **AND** 错误消息为 "用户名或密码错误"

#### Scenario: 密码错误
- **WHEN** 用户名存在但密码不正确
- **THEN** 系统返回401错误
- **AND** 错误消息为 "用户名或密码错误"

#### Scenario: 缺少必填字段
- **WHEN** 请求缺少用户名或密码字段
- **THEN** 系统返回400错误
- **AND** 错误消息提示缺失的字段

#### Scenario: 用户未关联角色
- **WHEN** 用户账号存在但未关联任何角色
- **THEN** 系统允许登录
- **AND** 返回的角色信息为空

### Requirement: 用户登出
系统 SHALL 允许已认证用户登出系统。

#### Scenario: 登出成功
- **WHEN** 已认证用户请求登出
- **THEN** 系统清除会话状态
- **AND** 返回成功响应 `{"code": 0, "message": "登出成功"}`

#### Scenario: 未认证用户登出
- **WHEN** 未提供有效JWT token的用户请求登出
- **THEN** 系统返回401错误
- **AND** 错误消息为 "未授权"

### Requirement: 获取用户信息
系统 SHALL 允许已认证用户获取当前登录的用户信息。

#### Scenario: 获取用户信息成功
- **WHEN** 已认证用户请求获取个人信息
- **THEN** 系统返回用户名
- **AND** 响应格式为 `{"code": 0, "message": "success", "data": {"username": "..."}}`

#### Scenario: Token无效
- **WHEN** 用户提供无效或过期的JWT token
- **THEN** 系统返回401错误
- **AND** 错误消息为 "token无效或已过期"

#### Scenario: Token缺失
- **WHEN** 请求未包含Authorization header
- **THEN** 系统返回401错误
- **AND** 错误消息为 "未提供token"

### Requirement: 同步角色信息
系统 SHALL 允许已认证用户同步其角色信息，检测角色变更。

#### Scenario: 角色未变更
- **WHEN** 用户的角色信息与token中的一致
- **THEN** 系统返回 `{"role_changed": false, "role_id": <id>, "is_super_admin": <bool>}`

#### Scenario: 角色已变更
- **WHEN** 用户的角色信息与token中的不一致
- **THEN** 系统返回 `{"role_changed": true, "role_id": <new_id>, "is_super_admin": <new_bool>}`
- **AND** 前端应提示用户重新登录

#### Scenario: 用户角色被删除
- **WHEN** 用户关联的角色在数据库中不存在
- **THEN** 系统返回空角色信息
- **AND** `role_changed` 标记为 true

### Requirement: JWT Token认证
系统 SHALL 使用JWT token保护需要认证的API端点。

#### Scenario: 有效Token访问受保护资源
- **WHEN** 用户使用有效JWT token访问受保护的API
- **THEN** 系统验证token并允许访问

#### Scenario: Token过期
- **WHEN** 用户使用过期的JWT token访问受保护的API
- **THEN** 系统返回401错误
- **AND** 错误消息为 "token无效或已过期"

#### Scenario: Token格式错误
- **WHEN** 用户提供格式错误的JWT token
- **THEN** 系统返回401错误
- **AND** 错误消息为 "token格式错误"

#### Scenario: 未提供Token
- **WHEN** 用户未提供JWT token访问受保护的API
- **THEN** 系统返回401错误
- **AND** 错误消息为 "未提供token"

### Requirement: 密码安全
系统 SHALL 使用bcrypt算法加密存储用户密码。

#### Scenario: 密码加密存储
- **WHEN** 创建或更新用户密码
- **THEN** 系统使用bcrypt（成本因子10）加密密码
- **AND** 数据库中不存储明文密码

#### Scenario: 密码验证
- **WHEN** 用户登录时
- **THEN** 系统使用bcrypt比对输入密码与存储的哈希值

### Requirement: Token生成
系统 SHALL 在用户登录成功后生成JWT token。

#### Scenario: Token包含用户信息
- **WHEN** 生成JWT token
- **THEN** Token payload包含用户ID、用户名、角色ID和是否超级管理员标识

#### Scenario: Token过期时间
- **WHEN** 生成JWT token
- **THEN** Token有效期为配置文件中指定的时长（默认24小时）

## API Endpoints

### POST /api/login
登录接口

**Request:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin",
    "is_super_admin": true
  }
}
```

**Response (401):**
```json
{
  "code": 401,
  "message": "用户名或密码错误"
}
```

### POST /api/logout
登出接口（需要认证）

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "code": 0,
  "message": "登出成功"
}
```

### GET /api/profile
获取用户信息（需要认证）

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "username": "admin"
  }
}
```

### GET /api/sync-role
同步角色信息（需要认证）

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "role_changed": false,
    "role_id": 1,
    "is_super_admin": true
  }
}
```

## Data Models

### UserAccount
```go
type UserAccount struct {
    ID         uint      `gorm:"column:id;primaryKey"`
    Username   string    `gorm:"column:username;size:50;uniqueIndex;not null"`
    Name       string    `gorm:"column:name;size:100"`
    Phone      string    `gorm:"column:phone;size:20"`
    RoleID     *uint     `gorm:"column:role_id"`
    Password   string    `gorm:"column:password;size:255;not null"`
    Status     int       `gorm:"column:status;default:0"`
    CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
    UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}
```

## Configuration

### JWT设置
- **Secret**: 从配置文件读取（`config.yaml` 的 `jwt.secret`）
- **过期时间**: 从配置文件读取（默认24小时）
- **签名算法**: HS256

### 密码加密
- **算法**: bcrypt
- **成本因子**: 10

## Security Considerations

1. **密码存储**: 使用bcrypt加密，成本因子10
2. **Token传输**: 通过Authorization header的Bearer token方式传输
3. **敏感信息**: 登录失败时不泄露用户名是否存在
4. **Token过期**: 强制Token过期时间，防止永久有效token
5. **HTTPS**: 生产环境必须使用HTTPS传输
