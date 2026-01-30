# 提案：实现营销管理模块

## 为什么

教育培训机构需要营销管理功能来促进课程销售和提升客户满意度。原 Python 项目已经实现了完整的营销管理模块，包括活动模板和活动管理功能。为了实现技术栈统一和性能提升，需要将这些功能从 Python 迁移到 Go 项目。

本次迁移将保持数据库表结构不变，确保 API 接口完全兼容，使前端无需任何修改即可切换到新后端。

## 变更内容

本次变更完整实现了营销管理模块，包括：

1. **活动模板管理** (`activity-template` 规范)
   - 实现活动模板的 CRUD 操作
   - 支持按分类或按商品配置参与范围
   - 实现模板状态管理和编辑权限控制
   - 添加模板删除前的关联活动检查

2. **活动管理** (`activity` 规范)
   - 实现活动的 CRUD 操作
   - 实现活动详情（满折规则）配置
   - 实现活动状态管理和编辑权限控制
   - 实现按日期范围查询活动及冲突检测
   - 支持订单关联活动

3. **商品筛选增强** (`goods` 规范)
   - 为商品查询接口添加 classifyid 和 status 参数支持
   - 支持按分类和状态筛选商品

4. **Bug 修复**
   - 修复活动模板和活动列表返回格式不一致问题
   - 修复状态更新时 required 验证失败问题（使用指针类型）
   - 修复时间解析错误（支持多种时间格式）
   - 修复折扣值验证逻辑（支持百分比格式 0-100）
   - 移除 ActivityTemplateGoods 实体中不存在的 CreateTime 字段

详细的规范变更请参见 `specs/activity-template/spec.md` 和 `specs/activity/spec.md`。

## 概述

从 Python 项目 (ZhixinStudentSaaS) 迁移营销管理模块到 Go 项目 (CharonOMS)，使用 Golang 重写所有业务逻辑，保持数据库表结构不变，完全兼容原有 API 接口。

## 动机

1. **功能完整性**：原 Python 项目已实现完整的营销管理功能，包括活动模板管理和活动管理
2. **业务连续性**：营销管理是教育培训机构的核心功能，与订单系统紧密结合
3. **技术统一**：将 Python 代码迁移到 Go，统一技术栈，提升性能和可维护性
4. **数据库已就绪**：数据库表结构已经创建完成，无需额外的数据库迁移工作

## 范围

### 包含的功能

#### 1. 活动模板管理 (Activity Template)
- 活动模板 CRUD 操作
- 活动模板状态管理（启用/禁用）
- 支持按分类或按商品配置参与范围
- 支持三种活动类型：满减（type=1）、满折（type=2）、满赠（type=3）

#### 2. 活动管理 (Activity)
- 活动 CRUD 操作
- 活动状态管理（启用/禁用）
- 活动时间范围管理
- 活动详情配置（满折规则）
- 按日期范围查询启用的活动
- 同类型活动冲突检测

#### 3. 与订单系统集成
- 订单关联活动（orders_activity 表）
- 活动参与商品筛选
- 活动折扣计算逻辑

### 不包含的功能

- 前端页面修改（前端保持不变）
- 数据库表结构修改（使用现有表结构）
- 其他模块的修改（仅限营销管理模块）

## 目标

1. **API 完全兼容**：所有 API 路由、请求格式、响应格式与 Python 版本完全一致
2. **业务逻辑一致**：所有验证规则、错误处理、业务流程与 Python 版本一致
3. **代码质量**：遵循 Go 项目的 DDD 架构，代码清晰可维护
4. **测试覆盖**：核心业务逻辑单元测试覆盖率 > 80%
5. **无缝迁移**：前端无需任何修改即可切换到新后端

## 技术方案

### 架构设计

采用项目现有的 DDD 分层架构：

```
internal/
├── domain/
│   ├── activity/
│   │   ├── entity/           # Activity, ActivityDetail 实体
│   │   ├── repository/       # Repository 接口定义
│   │   └── service/          # 领域服务
│   └── activity_template/
│       ├── entity/           # ActivityTemplate, ActivityTemplateGoods 实体
│       ├── repository/       # Repository 接口定义
│       └── service/          # 领域服务
├── application/
│   ├── activity/
│   │   ├── service.go       # Activity 应用服务
│   │   └── assembler.go     # DTO <-> Entity 转换
│   └── activity_template/
│       ├── service.go       # ActivityTemplate 应用服务
│       └── assembler.go     # DTO <-> Entity 转换
├── infrastructure/
│   └── persistence/
│       ├── activity_repository.go
│       └── activity_template_repository.go
└── interfaces/
    └── http/
        ├── handler/
        │   ├── activity/
        │   │   └── activity_handler.go
        │   └── activity_template/
        │       └── activity_template_handler.go
        ├── dto/
        │   ├── activity_dto.go
        │   └── activity_template_dto.go
        └── router/
            └── activity_routes.go
```

### API 路由映射

从 Python 项目迁移以下 API 接口：

**活动模板 API**:
- `GET /api/activity-templates` - 查询活动模板列表
- `GET /api/activity-templates/:id` - 获取活动模板详情
- `POST /api/activity-templates` - 创建活动模板
- `PUT /api/activity-templates/:id` - 更新活动模板
- `DELETE /api/activity-templates/:id` - 删除活动模板
- `PUT /api/activity-templates/:id/status` - 更新活动模板状态
- `GET /api/activity-templates/active` - 获取启用的活动模板

**活动 API**:
- `GET /api/activities` - 查询活动列表
- `GET /api/activities/:id` - 获取活动详情
- `POST /api/activities` - 创建活动
- `PUT /api/activities/:id` - 更新活动
- `DELETE /api/activities/:id` - 删除活动
- `PUT /api/activities/:id/status` - 更新活动状态
- `GET /api/activities/by-date-range` - 按日期范围查询活动

### 数据库表

使用已创建的数据库表（不做任何修改）：

1. **activity_template** - 活动模板表
   - id, name, type, select_type, status, create_time, update_time

2. **activity_template_goods** - 活动模板商品/分类关联表
   - id, template_id, goods_id, classify_id

3. **activity** - 活动表
   - id, name, template_id, start_time, end_time, status, create_time

4. **activity_detail** - 活动详情表（满折规则）
   - id, activity_id, threshold_amount, discount_value

5. **orders_activity** - 订单活动关联表
   - id, orders_id, activity_id, create_time

## 成功标准

1. **功能完整性**
   - 所有 Python 版本的 API 接口都已实现
   - 所有业务逻辑与 Python 版本一致

2. **API 兼容性**
   - 请求格式与 Python 版本完全一致
   - 响应格式与 Python 版本完全一致
   - HTTP 状态码与 Python 版本完全一致
   - 错误消息与 Python 版本完全一致

3. **数据完整性**
   - 所有数据库操作正确执行
   - 外键约束正确处理
   - 事务正确使用

4. **测试覆盖**
   - 核心业务逻辑单元测试覆盖率 > 80%
   - 所有 API 接口集成测试通过
   - 边界条件和错误情况测试通过

5. **前端兼容**
   - 前端无需任何修改
   - 所有前端功能正常工作

## 风险与缓解

### 风险1：API 响应格式不一致
**影响**: 高
**概率**: 中
**缓解**:
- 详细对比 Python 版本的 API 响应格式
- 编写集成测试验证响应格式
- 使用 Python 代码中的错误消息文本

### 风险2：业务逻辑理解偏差
**影响**: 高
**概率**: 中
**缓解**:
- 详细分析 Python 代码的业务逻辑
- 编写详细的场景测试
- 与业务需求文档对齐

### 风险3：数据库事务处理不当
**影响**: 高
**概率**: 低
**缓解**:
- 使用 GORM 事务机制
- 参考现有模块的事务处理模式
- 编写事务回滚测试

### 风险4：时间字段格式不一致
**影响**: 中
**概率**: 中
**缓解**:
- 统一使用 RFC3339 时间格式
- 编写时间格式转换测试
- 参考 Python 代码的时间处理逻辑

## 依赖关系

### 前置依赖
- 数据库表已创建 ✅（已完成）
- 商品管理模块已实现 ✅（已完成）
- 分类管理模块已实现 ✅（已完成）
- 订单管理模块基础框架已存在 ⚠️（需确认）

### 被依赖关系
- 订单管理模块需要使用活动查询接口
- 前端营销管理页面依赖这些 API

## 替代方案

### 方案A：渐进式迁移
逐步迁移功能，先实现活动模板，再实现活动管理。

**优点**:
- 降低单次变更风险
- 可以更早获得反馈

**缺点**:
- 两个模块紧密相关，分开实现可能导致重复工作
- 增加整体迁移时间

### 方案B：一次性完整迁移（推荐）
一次性实现所有营销管理功能。

**优点**:
- 模块完整性好
- 避免重复工作
- 缩短整体迁移时间

**缺点**:
- 单次变更较大
- 需要更多测试工作

**选择理由**: 营销管理模块功能相对独立且完整，一次性迁移更合理。

## 后续工作

完成本次迁移后，可能需要：

1. 订单创建时关联活动的逻辑实现
2. 活动优惠计算逻辑集成到订单系统
3. 营销报表和统计功能（如果原系统有）
4. 性能优化（如缓存启用的活动列表）

## 相关规范

- `openspec/specs/activity-template/spec.md` - 活动模板功能规范
- `openspec/specs/activity/spec.md` - 活动管理功能规范
- `openspec/project.md` - 项目架构和约定

## 参考资料

- Python 源代码: `D:\claude space\ZhixinStudentSaaS\backend\app.py` (行 2619-3186)
- 数据库恢复脚本: `D:\claude space\CharonOMS\scripts\restore_missing_tables_charonoms.sql`
- 数据库恢复报告: `D:\claude space\CharonOMS\scripts\database_restore_report.md`
