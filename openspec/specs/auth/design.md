# Design: User Authentication

## Context

CharonOMS 需要一个安全的用户认证系统，从原 Python Flask 项目迁移到 Go 版本时，保持API兼容性的同时提升性能和安全性。

**约束条件：**
- 必须与原Python版本的API响应格式保持一致
- 前端Vue.js应用未修改，需要完全兼容
- 使用JWT实现无状态认证
- 采用DDD分层架构

**利益相关者：**
- 前端开发者：依赖稳定的API契约
- 系统管理员：需要安全的认证机制
- 终端用户：需要简单可靠的登录体验

## Goals / Non-Goals

**Goals:**
- 实现JWT-based无状态认证
- 使用bcrypt安全存储密码
- 保持与原Python版本API格式完全兼容
- 实现角色信息同步机制
- 提供清晰的认证中间件

**Non-Goals:**
- OAuth/SAML等第三方认证（未来功能）
- 多因素认证（MFA）
- 密码重置功能（当前版本不包含）
- 会话管理（使用无状态JWT）

## Architecture

### Layer Structure (DDD)

```
┌─────────────────────────────────────────┐
│  Interfaces Layer                       │
│  - handler/auth/auth_handler.go         │
│  - middleware/jwt.go                    │
│  - router/router.go                     │
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│  Application Layer                      │
│  - service/auth/auth_service.go         │
│  - DTO transformations                  │
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│  Domain Layer                           │
│  - entity/user_account.go               │
│  - repository/auth_repository.go (IF)   │
└─────────────────────────────────────────┘
                   ↑
┌─────────────────────────────────────────┐
│  Infrastructure Layer                   │
│  - persistence/mysql/auth/              │
│    auth_repository_impl.go              │
│  - config/jwt_config.go                 │
└─────────────────────────────────────────┘
```

### Component Interaction

```
Client Request
      ↓
[Router] → [JWT Middleware] → [Auth Handler]
                                      ↓
                              [Auth Service]
                                      ↓
                              [Auth Repository]
                                      ↓
                                  [MySQL DB]
```

## Decisions

### Decision 1: JWT for Authentication
**选择**: 使用JWT (JSON Web Token) 实现无状态认证

**理由**:
- 无状态：不需要服务器端会话存储
- 可扩展：支持水平扩展和负载均衡
- 跨域友好：支持前后端分离架构
- 标准化：使用成熟的 golang-jwt/jwt 库

**备选方案**:
- Session-based认证：需要Redis/数据库存储，增加复杂度
- API Key：缺乏过期机制，安全性较低

**实现细节**:
```go
// Token payload结构
type Claims struct {
    UserID       uint   `json:"user_id"`
    Username     string `json:"username"`
    RoleID       uint   `json:"role_id"`
    IsSuperAdmin bool   `json:"is_super_admin"`
    jwt.RegisteredClaims
}
```

### Decision 2: bcrypt for Password Hashing
**选择**: 使用bcrypt算法加密密码，成本因子设为10

**理由**:
- 自适应：成本因子可调，对抗硬件性能提升
- 自动加盐：内置随机盐，防止彩虹表攻击
- 行业标准：广泛使用，经过时间验证
- Go标准库支持：`golang.org/x/crypto/bcrypt`

**备选方案**:
- SHA256 + Salt：计算太快，容易被暴力破解
- Argon2：更现代，但bcrypt已足够安全且更成熟

**实现细节**:
```go
// 加密密码
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

// 验证密码
err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
```

### Decision 3: Repository Pattern
**选择**: 使用Repository模式抽象数据访问层

**理由**:
- 依赖倒置：Domain层定义接口，Infrastructure层实现
- 可测试性：易于Mock数据访问进行单元测试
- 数据库无关：理论上可切换MySQL到PostgreSQL
- DDD最佳实践：符合领域驱动设计原则

**接口定义**:
```go
type AuthRepository interface {
    FindByUsername(ctx context.Context, username string) (*entity.UserAccount, error)
    GetRoleByUserID(ctx context.Context, userID uint) (*entity.Role, error)
}
```

### Decision 4: Middleware-based Authorization
**选择**: 使用Gin中间件实现JWT验证

**理由**:
- 关注点分离：认证逻辑与业务逻辑解耦
- 可复用：一次定义，多个路由使用
- 清晰的路由保护：显式标记哪些路由需要认证
- 框架标准做法：符合Gin的最佳实践

**实现位置**: `internal/interfaces/http/middleware/jwt.go`

### Decision 5: Response Format Compatibility
**选择**: 保持与原Python版本完全一致的响应格式

**理由**:
- 前端不变：Vue.js前端代码无需修改
- 平滑迁移：可逐步切换到Go后端
- 降低风险：减少因格式变化导致的bug

**标准响应格式**:
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

## Data Flow

### Login Flow
```
1. Client → POST /api/login {username, password}
2. Handler → Service.Login(username, password)
3. Service → Repository.FindByUsername(username)
4. Service → bcrypt.CompareHashAndPassword(stored, input)
5. Service → GenerateJWT(user)
6. Service → Repository.GetRoleByUserID(user.ID)
7. Handler ← Service: LoginResponse{token, username, is_super_admin}
8. Client ← Handler: JSON Response
```

### Protected Resource Access Flow
```
1. Client → GET /api/profile (Authorization: Bearer <token>)
2. Middleware → ValidateJWT(token)
3. Middleware → Extract Claims
4. Middleware → Set user info in Gin Context
5. Handler → Get user from Context
6. Handler → Return user profile
7. Client ← Handler: JSON Response
```

## Security Measures

### 1. Password Security
- **存储**: bcrypt哈希，成本因子10
- **传输**: 仅通过HTTPS（生产环境）
- **验证**: 使用constant-time比较防止时序攻击

### 2. JWT Security
- **Secret管理**: 从配置文件读取，不硬编码
- **Token过期**: 强制24小时过期
- **签名算法**: HS256（HMAC with SHA-256）
- **Payload**: 不包含敏感信息（如密码）

### 3. Error Handling
- **登录失败**: 不泄露用户名是否存在
- **统一错误消息**: "用户名或密码错误"
- **日志记录**: 记录失败尝试但不记录密码

### 4. Rate Limiting
- **当前状态**: 未实现
- **建议**: 后续添加登录尝试次数限制

## Configuration

### config/config.yaml
```yaml
jwt:
  secret: "your-secret-key-here"  # 生产环境使用强密钥
  expire_hours: 24                # Token有效期（小时）

database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  database: "charonoms"
```

### Environment Variables (Optional)
```bash
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24
```

## Error Handling

### Error Types
1. **400 Bad Request**: 缺少必填字段
2. **401 Unauthorized**: 认证失败、Token无效/过期
3. **500 Internal Server Error**: 数据库错误、系统错误

### Error Response Format
```json
{
  "code": <http_status_code>,
  "message": "<error_message>"
}
```

## Testing Strategy

### Unit Tests
- Service层：测试登录逻辑、密码验证、JWT生成
- Repository层（使用Mock）：测试数据访问逻辑

### Integration Tests
- Repository层：使用测试数据库测试实际查询
- API层：测试完整的登录流程（端到端）

### Test Cases
```go
func TestAuthService_Login_Success(t *testing.T) { ... }
func TestAuthService_Login_WrongPassword(t *testing.T) { ... }
func TestAuthService_Login_UserNotFound(t *testing.T) { ... }
func TestJWTMiddleware_ValidToken(t *testing.T) { ... }
func TestJWTMiddleware_ExpiredToken(t *testing.T) { ... }
```

## Migration from Python

### API Compatibility Checklist
- [x] POST /api/login 响应格式一致
- [x] POST /api/logout 响应格式一致
- [x] GET /api/profile 响应格式一致
- [x] GET /api/sync-role 响应格式一致
- [x] 401错误消息格式一致
- [x] JWT token格式兼容（前端可解析）

### Database Schema Compatibility
- [x] useraccount表结构保持不变
- [x] 字段名称完全一致（username, password, role_id等）
- [x] 密码哈希格式兼容（bcrypt）

## Performance Considerations

### Expected Load
- 登录请求：~100 req/s
- 认证检查：~1000 req/s（所有受保护API）

### Optimizations
1. **JWT验证**: 纯计算，无数据库查询，性能高
2. **连接池**: GORM自动管理数据库连接池
3. **缓存**: 暂未实现，可考虑缓存用户角色信息

### Bottlenecks
- 数据库查询：`FindByUsername` 和 `GetRoleByUserID`
- bcrypt验证：CPU密集型（成本因子10约50ms）

## Risks / Trade-offs

### Risk 1: JWT无法主动撤销
**风险**: Token签发后在过期前一直有效，无法主动撤销

**缓解措施**:
- 设置较短的过期时间（24小时）
- 实现角色同步机制（/api/sync-role）检测权限变更
- 后续可添加Token黑名单（Redis）

### Risk 2: 密码重置功能缺失
**风险**: 用户忘记密码无法自助重置

**缓解措施**:
- 当前由管理员手动重置
- 未来版本添加邮件重置功能

### Risk 3: 无Rate Limiting
**风险**: 可能被暴力破解攻击

**缓解措施**:
- bcrypt自带延时（成本因子10约50ms/次）
- 建议后续添加IP-based rate limiting

## Open Questions

1. **Multi-tenancy**: 是否需要支持多租户？
   - **当前答案**: 否，当前版本单租户

2. **Remember Me**: 是否需要"记住我"功能（更长Token有效期）？
   - **当前答案**: 否，统一24小时过期

3. **Refresh Token**: 是否需要实现Token刷新机制？
   - **当前答案**: 否，过期后重新登录

## Future Enhancements

1. **密码策略**: 强制密码复杂度要求
2. **密码重置**: 邮件验证码重置密码
3. **登录日志**: 记录登录历史和IP地址
4. **多因素认证**: TOTP或短信验证码
5. **OAuth集成**: 支持第三方登录（微信、钉钉）
6. **Refresh Token**: 自动刷新Token机制
7. **Rate Limiting**: 防暴力破解
