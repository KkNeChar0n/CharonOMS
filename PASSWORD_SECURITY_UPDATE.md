# 密码安全更新说明

## ✅ 已完成的修复

### 1. 修改登录验证逻辑
- **文件**: `internal/application/service/auth/auth_service.go`
- **修改内容**: 将密码验证从明文比对改为使用 bcrypt 加密验证
- **影响**: 所有用户登录都需要使用 bcrypt 加密的密码

### 2. 添加密码工具函数
- **HashPassword()**: 使用 bcrypt 加密密码
- **VerifyPassword()**: 验证密码是否匹配

### 3. 创建密码迁移工具
- **文件**: `scripts/migrate_passwords.go`
- **功能**: 自动将数据库中的明文密码转换为 bcrypt 加密
- **使用方法**: 运行 `scripts/migrate_passwords.bat` (Windows) 或 `go run scripts/migrate_passwords.go`

### 4. 创建bcrypt加密的初始化数据
- **文件**: `scripts/init_data_bcrypt.sql`
- **内容**: 包含bcrypt加密密码的测试用户数据
- **测试用户密码**: 所有用户的密码都是 `password` (已加密)

### 5. 更新文档
- **README.md**: 添加密码迁移章节
- **PASSWORD_MIGRATION_README.md**: 详细的密码迁移指南

## 🔒 安全改进

### Before (不安全)
```go
// 明文密码比对
if user.Password != req.Password {
    return nil, errors.ErrInvalidCredentials
}
```

### After (安全)
```go
// bcrypt加密验证
err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
if err != nil {
    return nil, errors.ErrInvalidCredentials
}
```

### 密码存储对比

| 项目 | 明文密码 | bcrypt加密 |
|------|---------|-----------|
| 示例 | `password` | `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy` |
| 长度 | 8字符 | 60字符 |
| 可逆 | ✅ 是（严重安全隐患） | ❌ 否（单向加密） |
| 防彩虹表 | ❌ 不防护 | ✅ 自动加盐 |
| 安全性 | 🔴 极低 | 🟢 行业标准 |

## 📋 使用指南

### 新项目初始化

如果是全新项目（没有现有用户数据）:

```bash
# 1. 创建数据库
mysql -u root -p -e "CREATE DATABASE charonoms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 2. 运行数据库迁移（创建表结构）
# 这里需要你自己的表结构创建脚本

# 3. 初始化测试数据（使用bcrypt加密的密码）
mysql -u root -p charonoms < scripts/init_data_bcrypt.sql
```

### 从旧系统迁移

如果已有明文密码的用户数据:

```bash
# 1. 备份数据库
mysqldump -u root -p charonoms > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. 运行密码迁移脚本
# Windows:
scripts\migrate_passwords.bat

# Linux/Mac:
go run scripts/migrate_passwords.go

# 3. 验证迁移成功
# 检查密码格式和长度
mysql -u root -p charonoms -e "SELECT id, username, LEFT(password, 10) as pass_prefix, LENGTH(password) as pass_len FROM useraccount;"
```

## 🧪 测试验证

### 1. 测试密码加密功能

```bash
cd "D:\claude space\CharonOMS"
go run scripts/test_password.go
```

预期输出:
```
✅ 加密成功
   明文密码: password
   加密后: $2a$10$7MWoDuoB1Szyvcv2sgGopeybYSrfERSaTDMFpCn8wDHJ4pP6drNRi
   长度: 60 字符
✅ 正确密码验证成功
✅ 错误密码验证失败 (符合预期)
```

### 2. 测试登录API

```bash
# 测试登录
curl -X POST http://localhost:5001/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

预期响应:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin",
    "is_super_admin": true
  }
}
```

### 3. 测试错误密码

```bash
curl -X POST http://localhost:5001/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrongpassword"}'
```

预期响应:
```json
{
  "code": 401,
  "message": "用户名或密码错误"
}
```

## 📝 开发注意事项

### 创建新用户时

在代码中创建新用户时，必须使用 `HashPassword()` 函数加密密码:

```go
import "charonoms/internal/application/service/auth"

// 加密密码
hashedPassword, err := auth.HashPassword("user_input_password")
if err != nil {
    // 处理错误
    return err
}

// 保存到数据库
user := &entity.UserAccount{
    Username: "newuser",
    Password: hashedPassword,  // 保存加密后的密码
    // ... 其他字段
}
```

### 密码重置功能

如果需要实现密码重置功能:

```go
func ResetPassword(userID uint, newPassword string) error {
    // 1. 加密新密码
    hashedPassword, err := auth.HashPassword(newPassword)
    if err != nil {
        return err
    }

    // 2. 更新数据库
    err = db.Model(&entity.UserAccount{}).
        Where("id = ?", userID).
        Update("password", hashedPassword).Error

    return err
}
```

### 密码验证

如果需要验证用户输入的密码是否正确:

```go
import "charonoms/internal/application/service/auth"

// 验证密码
err := auth.VerifyPassword(user.Password, userInputPassword)
if err != nil {
    // 密码错误
    return errors.ErrInvalidCredentials
}
// 密码正确
```

## 🔐 安全最佳实践

### 1. 密码复杂度要求（建议实现）

```go
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("密码长度至少8位")
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

    complexity := 0
    if hasUpper { complexity++ }
    if hasLower { complexity++ }
    if hasDigit { complexity++ }
    if hasSpecial { complexity++ }

    if complexity < 3 {
        return errors.New("密码必须包含大写字母、小写字母、数字、特殊字符中的至少3种")
    }

    return nil
}
```

### 2. 密码哈希验证限流（防暴力破解）

考虑在登录接口添加限流机制:
- 同一IP多次失败登录后锁定
- 验证码机制
- 账号锁定策略

### 3. 密码更新策略

建议实施:
- 定期提醒用户更新密码（如90天）
- 禁止重复使用最近N次密码
- 密码泄露检测（与已知泄露数据库对比）

### 4. 审计日志

记录所有密码相关操作:
- 密码修改
- 密码重置
- 登录失败
- 账号锁定

## 📊 性能影响

### bcrypt性能特征

- **单次加密**: ~100ms (cost=10)
- **单次验证**: ~100ms (cost=10)
- **并发处理**: Go会自动利用多核CPU

这个耗时是有意设计的，用于防止暴力破解。在正常登录场景下，100ms的延迟是可以接受的。

### 优化建议

如果遇到性能瓶颈:

1. **调整cost因子** (谨慎)
   ```go
   // 降低到8（更快，但安全性降低）
   bcrypt.GenerateFromPassword([]byte(password), 8)
   ```

2. **实施登录缓存**
   - Redis缓存有效token
   - 避免频繁查询数据库

3. **使用连接池**
   - 已在项目中配置（max_open_conns: 100）

## ✅ 完成清单

- [x] 修改登录验证逻辑使用bcrypt
- [x] 添加密码加密和验证工具函数
- [x] 创建密码迁移脚本
- [x] 创建bcrypt加密的初始化数据SQL
- [x] 更新README文档
- [x] 创建详细的迁移指南
- [x] 创建测试脚本验证功能
- [ ] 实施密码复杂度验证（可选）
- [ ] 添加登录限流机制（可选）
- [ ] 实施密码更新策略（可选）

## 🎯 下一步建议

1. **立即执行**: 运行密码迁移脚本（如有现有数据）
2. **测试验证**: 确保登录功能正常工作
3. **用户通知**: 通知用户密码已加密保护
4. **实施额外安全措施**: 考虑添加上述安全最佳实践

---

**安全提示**:
- 绝不在日志中记录明文密码
- 定期审计密码安全策略
- 保持bcrypt库更新到最新版本
- 考虑实施双因素认证(2FA)

---

最后更新: 2026-01-27
