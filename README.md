# CharonOMS - 订单管理系统

基于 DDD 分层架构的订单管理系统，使用 Go + Gin + GORM + JWT 重写。

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **日志**: Zap
- **配置**: Viper
- **前端**: Vue.js 3

## 项目结构

```
CharonOMS/
├── cmd/                          # 应用入口
│   └── server/
│       └── main.go              # 主程序
│
├── internal/                     # 内部代码
│   ├── interfaces/              # 接口层（API）
│   │   └── http/
│   │       ├── handler/        # 处理器
│   │       ├── middleware/     # 中间件
│   │       ├── router/         # 路由
│   │       └── dto/            # DTO
│   │
│   ├── application/             # 应用层
│   │   ├── service/           # 应用服务
│   │   └── assembler/         # 对象转换
│   │
│   ├── domain/                  # 领域层
│   │   ├── entity/            # 实体
│   │   ├── vo/                # 值对象
│   │   ├── repository/        # 仓储接口
│   │   └── service/           # 领域服务
│   │
│   └── infrastructure/          # 基础设施层
│       ├── persistence/       # 持久化
│       ├── config/            # 配置
│       ├── logger/            # 日志
│       └── utils/             # 工具
│
├── pkg/                          # 公共库
│   ├── jwt/                    # JWT 工具
│   ├── response/               # 统一响应
│   └── errors/                 # 错误定义
│
├── config/                       # 配置文件
├── docs/                         # 文档
├── frontend/                     # 前端代码
├── go.mod                        # Go 模块
└── README.md                     # 项目说明
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+

### 配置数据库

1. 创建数据库

```sql
CREATE DATABASE charonoms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改配置文件 `config/config.yaml`

```yaml
database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"
  database: "charonoms"
```

### 安装依赖

```bash
go mod download
```

### 运行项目

```bash
go run cmd/server/main.go
```

服务器将启动在 `http://localhost:5001`

### 访问前端

打开浏览器访问：`http://localhost:5001`

## API 文档

### 认证接口

#### 登录
```
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}

Response:
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

#### 获取用户信息
```
GET /api/profile
Authorization: Bearer <token>

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "username": "admin"
  }
}
```

#### 同步角色
```
GET /api/sync-role
Authorization: Bearer <token>

Response:
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

#### 登出
```
POST /api/logout
Authorization: Bearer <token>

Response:
{
  "code": 0,
  "message": "登出成功"
}
```

## 业务模块

### 已实现模块
- ✅ 认证模块（登录、登出、角色同步）

### 待实现模块
- ⏳ RBAC 权限系统（角色、权限、菜单）
- ⏳ 学生管理
- ⏳ 教练管理
- ⏳ 商品管理（品牌、分类、属性、商品）
- ⏳ 订单管理
- ⏳ 活动管理
- ⏳ 合同管理
- ⏳ 收款管理
- ⏳ 退款管理
- ⏳ 审批流程

## 开发指南

### 添加新模块

1. 在 `internal/domain/{module}/entity` 创建实体
2. 在 `internal/domain/{module}/repository` 创建仓储接口
3. 在 `internal/infrastructure/persistence/mysql/{module}` 实现仓储
4. 在 `internal/application/service/{module}` 创建应用服务
5. 在 `internal/interfaces/http/handler/{module}` 创建处理器
6. 在 `router.go` 中注册路由

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 导出的函数和结构体必须添加注释

### Git 提交规范

```
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试
chore: 构建/工具变动
```

## 部署

### 编译

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o charonoms cmd/server/main.go
```

### Docker 部署

```bash
docker build -t charonoms .
docker run -d -p 5001:5001 --name charonoms charonoms
```

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue。
