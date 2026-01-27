# 密码迁移指南

## 概述

此文档说明如何将数据库中的明文密码迁移到bcrypt加密密码。

## 背景

原Python Flask项目使用明文密码存储，存在严重的安全隐患。Go重构项目使用bcrypt加密算法来保护用户密码。

### bcrypt优势
- 自动加盐（salt）防止彩虹表攻击
- 可调节的计算成本（cost factor）
- 单向加密，无法反向解密
- 行业标准的密码加密算法

## 迁移前准备

### 1. 备份数据库

**⚠️ 强烈建议在迁移前备份数据库！**

```bash
# 备份整个数据库
mysqldump -u root -p charonoms > backup_before_migration_$(date +%Y%m%d_%H%M%S).sql

# 或仅备份useraccount表
mysqldump -u root -p charonoms useraccount > backup_useraccount_$(date +%Y%m%d_%H%M%S).sql
```

### 2. 配置环境变量

设置数据库连接信息（可选，默认使用localhost）：

**Windows:**
```cmd
set DB_HOST=localhost
set DB_PORT=3306
set DB_USER=root
set DB_PASSWORD=your_password
set DB_NAME=charonoms
```

**Linux/Mac:**
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=charonoms
```

## 迁移步骤

### 方式一：使用批处理脚本（推荐，Windows）

```cmd
cd D:\claude space\CharonOMS
scripts\migrate_passwords.bat
```

### 方式二：直接运行Go脚本

```bash
cd D:\claude space\CharonOMS
go run scripts/migrate_passwords.go
```

### 方式三：手动SQL迁移（高级用户）

如果Go脚本无法运行，可以手动为每个用户更新密码。

**注意**: 你需要先生成bcrypt哈希，然后手动更新数据库。

使用Python生成bcrypt哈希示例：
```python
import bcrypt
password = "your_password"
hashed = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt(rounds=10))
print(hashed.decode('utf-8'))
```

使用在线工具：https://bcrypt-generator.com/ (不推荐用于生产密码)

然后更新数据库：
```sql
UPDATE useraccount SET password = '$2a$10$...' WHERE username = 'admin';
```

## 迁移输出

脚本会显示以下信息：

```
数据库连接成功
[成功] 用户 admin (ID: 1) 密码已加密
[成功] 用户 user1 (ID: 2) 密码已加密
[跳过] 用户 user2 (ID: 3) 密码已经加密

=== 密码迁移完成 ===
总用户数: 3
成功加密: 2
已跳过: 1
失败: 0

✅ 所有密码已成功迁移到bcrypt加密
```

## 验证迁移

### 1. 检查密码格式

bcrypt加密的密码应该：
- 长度为60个字符
- 以 `$2a$`、`$2b$` 或 `$2y$` 开头

```sql
SELECT
    id,
    username,
    LEFT(password, 10) as password_prefix,
    LENGTH(password) as password_length
FROM useraccount;
```

预期输出：
```
+----+----------+------------------+-----------------+
| id | username | password_prefix  | password_length |
+----+----------+------------------+-----------------+
|  1 | admin    | $2a$10$abc      |              60 |
|  2 | user1    | $2a$10$xyz      |              60 |
+----+----------+------------------+-----------------+
```

### 2. 测试登录

使用API测试登录功能：

```bash
curl -X POST http://localhost:5001/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

应该返回：
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

## 故障排查

### 问题1: 数据库连接失败

**错误信息**: `无法连接数据库: ...`

**解决方案**:
1. 检查数据库是否运行
2. 验证环境变量或默认配置是否正确
3. 确认数据库用户有足够的权限

### 问题2: 某些用户迁移失败

**错误信息**: `[失败] 用户 xxx 密码加密失败`

**解决方案**:
1. 查看详细错误日志
2. 检查用户密码是否包含特殊字符
3. 手动为该用户重置密码

### 问题3: 迁移后无法登录

**可能原因**:
1. 迁移脚本未正常运行
2. 密码字段被截断
3. 数据库字段长度不足

**解决方案**:
1. 检查 `useraccount` 表的 `password` 字段长度（应至少为255）
   ```sql
   SHOW COLUMNS FROM useraccount LIKE 'password';
   ```
2. 如果长度不足，修改字段：
   ```sql
   ALTER TABLE useraccount MODIFY COLUMN password VARCHAR(255) NOT NULL;
   ```
3. 恢复备份并重新迁移

## 回滚步骤

如果迁移出现问题，可以从备份恢复：

```bash
# 恢复整个数据库
mysql -u root -p charonoms < backup_before_migration_YYYYMMDD_HHMMSS.sql

# 或仅恢复useraccount表
mysql -u root -p charonoms < backup_useraccount_YYYYMMDD_HHMMSS.sql
```

## 安全建议

### 迁移后

1. ✅ 删除包含明文密码的备份文件（在确认系统正常后）
2. ✅ 通知所有用户密码已加密
3. ✅ 考虑强制用户首次登录时更改密码
4. ✅ 实施密码复杂度策略

### 新用户创建

所有新创建的用户密码必须使用 `auth.HashPassword()` 函数加密：

```go
hashedPassword, err := auth.HashPassword("user_password")
if err != nil {
    // 处理错误
}
// 保存 hashedPassword 到数据库
```

## 技术细节

### bcrypt参数

- **算法**: bcrypt
- **Cost Factor**: 10
- **Salt**: 自动生成（16字节）
- **哈希长度**: 60字符
- **格式**: `$2a$10$<salt><hash>`

### 性能影响

- 单次加密耗时: ~100ms (cost=10)
- 单次验证耗时: ~100ms (cost=10)
- 数据库存储: 60字节 per password

这个耗时是有意设计的，用于防止暴力破解攻击。

## 相关文件

- `scripts/migrate_passwords.go` - 迁移脚本
- `scripts/migrate_passwords.bat` - Windows批处理启动器
- `internal/application/service/auth/auth_service.go` - 认证服务（包含密码验证）

## 联系支持

如有问题，请：
1. 检查本文档的故障排查部分
2. 查看GitHub Issues
3. 联系开发团队

---

最后更新: 2026-01-27
