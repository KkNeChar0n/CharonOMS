# 提案：迁移营销管理模块

## 变更标识
- **ID**: `migrate-marketing-module`
- **类型**: 功能迁移
- **优先级**: 中
- **提案日期**: 2026-01-29

## 为什么

### 背景

CharonOMS 从 Python Flask 版本迁移到 Go 版本的过程中，营销管理模块尚未迁移。营销管理模块包含两个子模块：

1. **活动模板管理**：创建营销活动模板（满减、满折、满赠等），定义活动类型和适用商品/分类
2. **活动管理**：基于活动模板创建具体的营销活动实例，设置活动时间、参与商品和优惠详情

这些模块在原 Python 版本中已经实现，前端页面已存在，但后端接口缺失。

### 动机

1. **业务核心功能**：营销活动是培训机构吸引客户、促进销售的重要手段，是业务运营的核心功能
2. **前端已就绪**：前端页面和交互逻辑已经完整实现，只需后端接口支持即可启用
3. **模块完整性**：活动模板和活动管理是配套功能，需要一起迁移才能形成完整的营销管理能力
4. **权限管理完善**：现有 RBAC 权限系统已经成熟，正是集成营销模块权限的最佳时机

## 目标

将营销管理模块从 Python 版本迁移到 Go 版本，实现以下功能：

### 活动模板管理
- 查询活动模板列表（支持ID、名称、类型、状态筛选）
- 获取活动模板详情（包含关联商品或分类）
- 新增活动模板（基本信息 + 关联配置）
- 更新活动模板（修改基本信息和关联配置）
- 更新活动模板状态（启用/禁用）
- 获取启用的活动模板列表（供活动管理使用）

### 活动管理
- 查询活动列表（支持ID、名称、模板、状态筛选）
- 获取活动详情（包含活动细节配置）
- 新增活动（基本信息 + 活动细节）
- 更新活动（修改基本信息和细节配置）
- 更新活动状态（启用/禁用）

### 权限管理集成
- 为所有操作注册相应权限（view, add, edit, enable, disable）
- 权限与菜单关联，支持基于角色的访问控制
- 前端已有权限检查逻辑，需要后端提供权限数据

## 范围

### 包含
- 活动模板管理的 6 个接口
- 活动管理的 5 个接口
- 权限数据播种（SQL 脚本）
- DDD 分层架构实现（Domain、Application、Infrastructure、Interface）
- 与前端现有页面的 API 兼容

### 不包含
- 审批流类型管理
- 审批流模板管理
- 审批流管理
- 活动效果统计和报表
- 活动自动执行引擎

## 受影响的系统

### 后端模块
- `internal/domain/marketing/` - 新增营销领域模块
- `internal/application/service/marketing/` - 新增营销应用服务
- `internal/infrastructure/persistence/mysql/marketing/` - 新增营销仓储实现
- `internal/interfaces/http/handler/marketing/` - 新增营销HTTP处理器
- `internal/interfaces/http/router/router.go` - 注册营销模块路由
- `scripts/seed_marketing_permissions.sql` - 权限数据播种脚本

### 数据库表
- `activity_template` - 活动模板表（已存在）
- `activity_template_classify` - 活动模板-分类关联表（已存在）
- `activity_template_goods` - 活动模板-商品关联表（已存在）
- `activity` - 活动表（需创建）
- `activity_detail` - 活动细节表（需创建）
- `permissions` - 权限表（新增营销模块权限数据）
- `menu` - 菜单表（已存在活动模板、活动管理菜单）

### 前端
- 前端已有完整页面和交互逻辑，无需修改
- API 路径需保持与 Python 版本一致

## 技术方案

### 架构设计
采用 DDD 四层架构：
1. **Domain 层**：定义营销领域实体和仓储接口
2. **Application 层**：实现业务用例和编排逻辑
3. **Infrastructure 层**：实现 MySQL 仓储
4. **Interface 层**：提供 HTTP API 接口

### 数据模型

#### 活动模板 (ActivityTemplate)
```go
type ActivityTemplate struct {
    ID         uint      `gorm:"primaryKey"`
    Name       string    `gorm:"size:100;not null"`
    Type       int       `gorm:"not null"` // 1-满减 2-满折 3-满赠
    SelectType int       `gorm:"not null"` // 1-按分类 2-按商品
    Status     int       `gorm:"default:0"`
    CreateTime time.Time `gorm:"column:create_time"`
    UpdateTime time.Time `gorm:"column:update_time"`

    // 关联（二选一）
    Classifies []*ActivityTemplateClassify `gorm:"foreignKey:TemplateID"`
    Goods      []*ActivityTemplateGoods    `gorm:"foreignKey:TemplateID"`
}
```

#### 活动 (Activity)
```go
type Activity struct {
    ID           uint      `gorm:"primaryKey"`
    TemplateID   uint      `gorm:"not null"`
    Name         string    `gorm:"size:100;not null"`
    StartTime    time.Time `gorm:"not null"`
    EndTime      time.Time `gorm:"not null"`
    Status       int       `gorm:"default:0"`
    CreateTime   time.Time `gorm:"column:create_time"`
    UpdateTime   time.Time `gorm:"column:update_time"`

    // 关联
    Template *ActivityTemplate  `gorm:"foreignKey:TemplateID"`
    Details  []*ActivityDetail  `gorm:"foreignKey:ActivityID"`
}
```

#### 活动细节 (ActivityDetail)
```go
type ActivityDetail struct {
    ID         uint      `gorm:"primaryKey"`
    ActivityID uint      `gorm:"not null"`
    GoodsID    uint      `gorm:"not null"`
    Discount   float64   `gorm:"type:decimal(10,2)"` // 折扣值或优惠金额
}
```

### API 设计

#### 活动模板
- `GET /api/activity-templates` - 获取列表（支持筛选）
- `GET /api/activity-templates/active` - 获取启用的模板列表
- `GET /api/activity-templates/:id` - 获取详情
- `POST /api/activity-templates` - 新增模板
- `PUT /api/activity-templates/:id` - 更新模板
- `PUT /api/activity-templates/:id/status` - 更新状态

#### 活动管理
- `GET /api/activities` - 获取列表（支持筛选）
- `GET /api/activities/:id` - 获取详情
- `POST /api/activities` - 新增活动
- `PUT /api/activities/:id` - 更新活动
- `PUT /api/activities/:id/status` - 更新状态

## 权限设计

### 菜单关联
- 活动模板管理（menu_id: 待确认）
  - view_activity_template
  - add_activity_template
  - edit_activity_template
  - enable_activity_template
  - disable_activity_template

- 活动管理（menu_id: 待确认）
  - view_activity
  - add_activity
  - edit_activity
  - enable_activity
  - disable_activity

## 风险和假设

### 风险
1. activity 和 activity_detail 表可能不存在，需要创建表结构
2. 活动细节的折扣计算逻辑需要与前端对齐
3. 活动时间范围验证规则需要明确

### 假设
1. 前端页面无需修改，API 格式完全兼容
2. activity 和 activity_detail 表结构与 Python 版本一致
3. 活动的优惠计算在前端完成，后端只存储配置

## 后续工作

1. 活动效果统计功能
2. 活动自动开启/关闭定时任务
3. 活动与订单系统的集成

## 参考资料

- Python 原版代码：`../ZhixinStudentSaaS/backend/app.py`
- 前端页面：`frontend/index.html` (activity_template, activity_management 相关部分)
- 数据库表结构：`charonoms_backup_20260127.sql`
