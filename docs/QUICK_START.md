# CharonOMS 快速启动指南

## 前置要求

- Go 1.21+
- MySQL 5.7+
- Git

## 步骤 1：创建数据库

打开 MySQL 客户端，执行：

```sql
CREATE DATABASE charonoms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 步骤 2：导入测试数据

执行初始化脚本：

```bash
mysql -u root -p charonoms < scripts/init_test_data.sql
```

或者在 MySQL 客户端中：

```sql
USE charonoms;
SOURCE d:/claude space/CharonOMS/scripts/init_test_data.sql;
```

## 步骤 3：配置数据库连接

修改 `config/config.yaml`：

```yaml
database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"    # 修改为你的密码
  database: "charonoms"
```

## 步骤 4：安装依赖

```bash
cd "d:/claude space/CharonOMS"
go mod download
```

## 步骤 5：运行项目

```bash
go run cmd/server/main.go
```

或使用 Makefile：

```bash
make run
```

## 步骤 6：测试登录

### 方法 1：使用 curl

```bash
curl -X POST http://localhost:5001/api/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"password\"}"
```

### 方法 2：使用浏览器

打开浏览器访问：http://localhost:5001

## 测试账号

系统预置了 3 个测试账号：

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | password | 超级管理员 | 拥有所有权限 |
| manager | password | 普通管理员 | 拥有学生、教练、商品等模块的管理权限 |
| operator | password | 操作员 | 仅拥有查看权限 |

## API 测试示例

### 1. 登录

```bash
curl -X POST http://localhost:5001/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

响应示例：
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

### 2. 获取用户菜单

```bash
curl -X GET http://localhost:5001/api/menu \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 获取角色列表

```bash
curl -X GET http://localhost:5001/api/roles \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. 获取权限树

```bash
curl -X GET http://localhost:5001/api/permissions/tree \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 常见问题

### 1. 无法连接数据库

检查：
- MySQL 服务是否启动
- 数据库用户名和密码是否正确
- 数据库 charonoms 是否已创建
- config/config.yaml 中的配置是否正确

### 2. Token 无效

检查：
- Authorization header 格式是否为 `Bearer <token>`
- Token 是否已过期（默认 24 小时）
- JWT secret 配置是否正确

### 3. 端口被占用

修改 `config/config.yaml` 中的 `server.port` 配置

## 下一步

- 查看 [API 文档](../README.md#api-文档)
- 查看 [开发指南](../README.md#开发指南)
- 开始实现其他业务模块

## 项目结构说明

```
CharonOMS/
├── cmd/server/          # 主程序入口
├── internal/
│   ├── interfaces/      # API 层（Handler、中间件、路由）
│   ├── application/     # 应用层（Service）
│   ├── domain/          # 领域层（Entity、Repository 接口）
│   └── infrastructure/  # 基础设施层（数据库、配置、日志）
├── pkg/                 # 公共库（JWT、响应、错误）
├── config/              # 配置文件
├── scripts/             # SQL 脚本
├── docs/                # 文档
└── frontend/            # 前端代码
```

## 已实现的功能

✅ 用户认证（登录、登出）
✅ JWT Token 认证
✅ RBAC 权限系统（角色、权限、菜单）
✅ 用户菜单树（基于权限）
✅ 权限中间件
✅ 日志系统
✅ 配置管理
✅ 统一响应格式
✅ 错误处理

## 待实现的功能

⏳ 学生管理
⏳ 教练管理
⏳ 商品管理
⏳ 订单管理
⏳ 活动管理
⏳ 合同管理
⏳ 收款管理
⏳ 退款管理
⏳ 审批流程
