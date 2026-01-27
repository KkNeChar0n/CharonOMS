# OpenSpec for CharonOMS

本项目已采用 OpenSpec 进行规范驱动开发。OpenSpec 帮助我们在实施代码之前明确需求和行为，确保所有变更都经过充分的规划和审查。

## 目录结构

```
openspec/
├── README.md              # 本文件
├── AGENTS.md              # AI助手工作流程指南
├── project.md             # 项目上下文和约定
│
├── specs/                 # 已实现功能的规范（事实）
│   ├── auth/             # 用户认证模块
│   │   ├── spec.md       # 认证需求和场景
│   │   └── design.md     # 认证技术设计
│   ├── rbac/             # 角色权限控制模块
│   │   └── spec.md       # RBAC需求和场景
│   └── basic/            # 基础数据模块
│       └── spec.md       # 基础数据需求和场景
│
└── changes/              # 待实施的变更提案
    └── archive/          # 已归档的变更
```

## 已实现模块

### ✅ Auth - 用户认证
- JWT-based 无状态认证
- 登录、登出、用户信息获取
- 角色信息同步
- bcrypt 密码加密

**规范位置**: `specs/auth/spec.md`
**设计文档**: `specs/auth/design.md`

### ✅ RBAC - 角色权限控制
- 角色管理（创建、更新、启用/禁用）
- 权限管理（查询、更新状态）
- 菜单管理（查询、更新）
- 角色-权限绑定
- 超级管理员机制

**规范位置**: `specs/rbac/spec.md`

### ✅ Basic - 基础数据
- 性别字典（男、女）
- 年级列表（一年级~高三）
- 学科列表（语文、数学等12个学科）

**规范位置**: `specs/basic/spec.md`

## 工作流程

### 阶段 1：创建变更提案

当您计划添加新功能或进行重大变更时：

1. **创建变更目录**:
   ```bash
   mkdir -p openspec/changes/add-student-management
   ```

2. **编写提案文件** (`proposal.md`):
   ```markdown
   # Change: 添加学生管理功能

   ## Why
   需要管理学生信息和教练关联

   ## What Changes
   - 添加学生 CRUD API
   - 添加学生-教练多对多关联

   ## Impact
   - 影响的规范：student, coach
   ```

3. **创建任务清单** (`tasks.md`):
   ```markdown
   ## 1. 数据层
   - [ ] 创建 Student 实体模型
   - [ ] 创建 StudentRepository

   ## 2. 业务层
   - [ ] 创建 StudentService
   ```

4. **编写规范增量** (`specs/student/spec.md`):
   ```markdown
   ## ADDED Requirements

   ### Requirement: 学生信息管理
   系统 SHALL 允许创建、查看、更新和删除学生信息。

   #### Scenario: 创建学生成功
   - **WHEN** 用户提供有效的学生信息
   - **THEN** 系统创建学生记录并返回ID
   ```

5. **验证提案** (如果安装了 openspec CLI):
   ```bash
   openspec validate add-student-management --strict
   ```

6. **等待批准** - 与团队成员或AI助手讨论并获得批准

### 阶段 2：实施变更

1. 阅读 `proposal.md` 了解目标
2. 阅读 `design.md`（如果有）了解技术决策
3. 按照 `tasks.md` 逐项实施
4. 完成一项标记 `[x]` 一项
5. 确保所有测试通过

### 阶段 3：归档变更

部署后归档变更：

```bash
openspec archive add-student-management --yes
```

这会：
- 将 `changes/add-student-management/` 移至 `changes/archive/YYYY-MM-DD-add-student-management/`
- 将规范增量合并到 `specs/student/spec.md`

## 快速命令参考

```bash
# 查看当前活跃的变更
openspec list

# 查看已实现的功能规范
openspec list --specs

# 查看特定变更详情
openspec show add-student-management

# 验证变更格式
openspec validate add-student-management --strict

# 归档已完成的变更
openspec archive add-student-management --yes
```

## 规范编写规则

### 需求格式
```markdown
### Requirement: 需求标题
系统 SHALL 描述行为...

#### Scenario: 场景名称
- **WHEN** 触发条件
- **THEN** 预期结果
- **AND** 额外条件（可选）
```

### 重要提示
1. ✅ 使用 `####` （4个井号）标记 Scenario
2. ✅ 每个需求至少有一个 Scenario
3. ✅ 使用 SHALL/MUST 表示规范性需求
4. ⛔ 不要使用 `- **Scenario:**` 或 `### Scenario:`

### 增量操作
- `## ADDED Requirements` - 新增需求
- `## MODIFIED Requirements` - 修改已有需求（需完整粘贴）
- `## REMOVED Requirements` - 删除需求
- `## RENAMED Requirements` - 重命名需求

## 何时创建提案

### ✅ 需要提案的情况
- 添加新功能或模块
- 破坏性 API 变更
- 数据库模式变更
- 架构或设计模式变更
- 性能优化（改变行为）

### ⛔ 跳过提案的情况
- Bug 修复（恢复预期行为）
- 拼写、格式、注释修改
- 非破坏性依赖更新
- 配置变更
- 现有功能的单元测试

## 与 AI 助手协作

AI 助手（Claude）会遵循以下流程：

1. **接收需求** - 您描述想要的功能
2. **检查现状** - AI 查看现有规范避免重复
3. **创建提案** - AI 起草完整的变更提案
4. **等待批准** - AI 不会在批准前开始编码
5. **实施变更** - 获得批准后 AI 按任务清单实施
6. **归档变更** - 部署后 AI 帮助归档

详细指南见：`openspec/AGENTS.md`

## 查看现有规范

### Auth 模块
```bash
# 查看认证规范
cat openspec/specs/auth/spec.md

# 查看认证设计
cat openspec/specs/auth/design.md
```

### RBAC 模块
```bash
# 查看权限规范
cat openspec/specs/rbac/spec.md
```

### Basic 模块
```bash
# 查看基础数据规范
cat openspec/specs/basic/spec.md
```

## 项目约定

详细的项目约定、技术栈、架构模式见：`openspec/project.md`

关键约定：
- **架构**: DDD 分层架构（接口层、应用层、领域层、基础设施层）
- **API格式**: `{"code": 0, "message": "success", "data": {...}}`
- **数据库**: MySQL 8.0+, utf8mb4 字符集
- **认证**: JWT Bearer Token
- **编码**: UTF-8 无 BOM

## 最佳实践

1. **在编码前规划** - 先写规范，后写代码
2. **保持规范更新** - 代码变更必须同步更新规范
3. **使用清晰场景** - 每个需求都有具体可测试的场景
4. **简单优先** - 默认实现最简单的方案
5. **明确影响范围** - 在提案中列出所有受影响的模块

## 问题和反馈

如有疑问，请：
1. 查阅 `AGENTS.md` 了解工作流程
2. 查阅 `project.md` 了解项目约定
3. 查看现有规范作为参考
4. 向团队成员或 AI 助手提问

---

**规范是事实，变更是提案。保持它们同步。**
