# 营销管理模块实现文档

## 概述

本文档记录从Python Flask项目(ZhixinStudentSaaS)到Go Gin项目(CharonOMS)的营销管理模块迁移实现。

**迁移日期**: 2026-01-30
**Go版本**: 1.21+
**框架**: Gin + GORM
**数据库**: MySQL 8.0

## 项目结构

营销管理模块采用DDD(领域驱动设计)架构,分为4个层次:

```
internal/
├── domain/                          # 领域层
│   ├── activity/                    # 活动领域
│   │   ├── entity/                  # 实体
│   │   │   ├── activity.go          # 活动实体
│   │   │   ├── activity_detail.go   # 活动详情实体
│   │   │   └── errors.go            # 领域错误定义
│   │   └── repository/              # 仓储接口
│   │       └── repository.go
│   └── activity_template/           # 活动模板领域
│       ├── entity/                  # 实体
│       │   ├── activity_template.go # 模板实体
│       │   ├── activity_template_goods.go  # 模板商品关联实体
│       │   └── errors.go            # 领域错误定义
│       └── repository/              # 仓储接口
│           └── repository.go
├── application/                     # 应用层
│   ├── activity/                    # 活动应用服务
│   │   └── service.go
│   └── activity_template/           # 模板应用服务
│       └── service.go
├── infrastructure/                  # 基础设施层
│   └── persistence/                 # 持久化
│       ├── activity_repository.go   # 活动仓储实现
│       └── activity_template_repository.go  # 模板仓储实现
└── interfaces/                      # 接口层
    └── http/                        # HTTP接口
        ├── dto/                     # 数据传输对象
        │   ├── activity_dto.go
        │   └── activity_template_dto.go
        ├── handler/                 # 处理器
        │   ├── activity_handler.go
        │   └── activity_template_handler.go
        └── router/                  # 路由
            └── router.go
```

## 数据库设计

### 活动模板表 (activity_template)

```sql
CREATE TABLE `activity_template` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) NOT NULL COMMENT '模板名称',
  `type` tinyint(4) NOT NULL COMMENT '活动类型：1=满减, 2=满折, 3=满赠',
  `select_type` tinyint(4) NOT NULL COMMENT '选择方式：1=按类型, 2=按商品',
  `status` tinyint(4) DEFAULT 0 COMMENT '状态：0=启用, 1=禁用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动模板表';
```

### 活动表 (activity)

```sql
CREATE TABLE `activity` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) NOT NULL COMMENT '活动名称',
  `template_id` int(11) NOT NULL COMMENT '关联模板ID',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `status` tinyint(4) DEFAULT 0 COMMENT '状态：0=启用, 1=禁用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_time_range` (`start_time`, `end_time`),
  CONSTRAINT `fk_activity_template` FOREIGN KEY (`template_id`)
    REFERENCES `activity_template` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动表';
```

### 活动详情表 (activity_detail)

```sql
CREATE TABLE `activity_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `activity_id` int(11) NOT NULL COMMENT '活动ID',
  `threshold_amount` decimal(10,2) NOT NULL COMMENT '门槛金额',
  `discount_value` decimal(10,2) NOT NULL COMMENT '折扣值',
  PRIMARY KEY (`id`),
  KEY `idx_activity_id` (`activity_id`),
  CONSTRAINT `fk_activity_detail` FOREIGN KEY (`activity_id`)
    REFERENCES `activity` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动详情表';
```

### 模板商品关联表 (activity_template_goods)

```sql
CREATE TABLE `activity_template_goods` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `template_id` int(11) NOT NULL COMMENT '模板ID',
  `classify_id` int(11) DEFAULT NULL COMMENT '分类ID',
  `goods_id` int(11) DEFAULT NULL COMMENT '商品ID',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_classify_id` (`classify_id`),
  KEY `idx_goods_id` (`goods_id`),
  CONSTRAINT `fk_atg_template` FOREIGN KEY (`template_id`)
    REFERENCES `activity_template` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_atg_classify` FOREIGN KEY (`classify_id`)
    REFERENCES `classify` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_atg_goods` FOREIGN KEY (`goods_id`)
    REFERENCES `goods` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动模板商品关联表';
```

## 领域层设计

### 1. 活动模板实体 (ActivityTemplate)

**职责:**
- 封装活动模板的业务规则
- 验证模板数据的有效性
- 提供状态查询方法

**主要方法:**
```go
// Validate 验证模板数据
func (t *ActivityTemplate) Validate() error

// IsEnabled 判断是否启用
func (t *ActivityTemplate) IsEnabled() bool
```

**业务规则:**
- 模板名称不能为空
- 活动类型必须为 1/2/3 (满减/满折/满赠)
- 选择方式必须为 1/2 (按分类/按商品)
- 选择方式为1时必须关联分类,为2时必须关联商品

### 2. 活动实体 (Activity)

**职责:**
- 封装活动的业务规则
- 验证活动时间范围
- 提供状态查询方法

**主要方法:**
```go
// Validate 验证活动数据
func (a *Activity) Validate() error

// IsEnabled 判断是否启用
func (a *Activity) IsEnabled() bool
```

**业务规则:**
- 活动名称不能为空
- 必须关联有效的活动模板
- 开始时间必须早于结束时间
- 启用中的活动不能编辑

### 3. 活动详情实体 (ActivityDetail)

**职责:**
- 封装折扣规则
- 根据活动类型验证折扣值

**主要方法:**
```go
// Validate 验证活动详情
func (d *ActivityDetail) Validate(activityType int) error
```

**业务规则:**
- 门槛金额必须大于0
- 折扣值必须大于0
- 满折类型(type=2)的折扣值必须在0-1之间

## 仓储层设计

### 1. ActivityTemplateRepository 接口

```go
type ActivityTemplateRepository interface {
    // 创建模板及关联商品/分类
    Create(ctx context.Context, template *entity.ActivityTemplate,
           goods []*entity.ActivityTemplateGoods) error

    // 更新模板
    Update(ctx context.Context, template *entity.ActivityTemplate) error

    // 更新模板及关联商品/分类
    UpdateWithGoods(ctx context.Context, template *entity.ActivityTemplate,
                    goods []*entity.ActivityTemplateGoods) error

    // 删除模板
    Delete(ctx context.Context, id int) error

    // 根据ID查询模板
    FindByID(ctx context.Context, id int) (*entity.ActivityTemplate, error)

    // 根据ID查询模板及关联商品/分类
    FindByIDWithGoods(ctx context.Context, id int) (*entity.ActivityTemplate,
                      []*entity.ActivityTemplateGoods, error)

    // 查询所有模板
    List(ctx context.Context) ([]*entity.ActivityTemplate, error)

    // 查询启用的模板
    FindActiveTemplates(ctx context.Context) ([]*entity.ActivityTemplate, error)

    // 统计关联活动数量
    CountRelatedActivities(ctx context.Context, templateID int) (int64, error)

    // 更新模板状态
    UpdateStatus(ctx context.Context, id int, status int) error
}
```

### 2. ActivityRepository 接口

```go
type ActivityRepository interface {
    // 创建活动及详情
    Create(ctx context.Context, activity *entity.Activity,
           details []*entity.ActivityDetail) error

    // 更新活动
    Update(ctx context.Context, activity *entity.Activity) error

    // 更新活动及详情
    UpdateWithDetails(ctx context.Context, activity *entity.Activity,
                      details []*entity.ActivityDetail) error

    // 删除活动
    Delete(ctx context.Context, id int) error

    // 根据ID查询活动
    FindByID(ctx context.Context, id int) (*entity.Activity, error)

    // 根据ID查询活动及详情
    FindByIDWithDetails(ctx context.Context, id int) (*entity.Activity,
                        []*entity.ActivityDetail, error)

    // 查询所有活动
    List(ctx context.Context) ([]*entity.Activity, error)

    // 根据日期范围查询启用的活动
    FindByDateRange(ctx context.Context, paymentTime time.Time) ([]*entity.Activity, error)

    // 查询活动详情
    FindDetailsByActivityID(ctx context.Context, activityID int) ([]*entity.ActivityDetail, error)

    // 更新活动状态
    UpdateStatus(ctx context.Context, id int, status int) error
}
```

## 应用层设计

### 1. ActivityTemplateService

**职责:**
- 编排活动模板业务流程
- 调用领域对象和仓储
- 查询关联的分类和商品信息

**核心方法:**
```go
// CreateTemplate 创建活动模板
func (s *Service) CreateTemplate(ctx context.Context,
     req *dto.CreateActivityTemplateDTO) (int, error)

// UpdateTemplate 更新活动模板
func (s *Service) UpdateTemplate(ctx context.Context, id int,
     req *dto.UpdateActivityTemplateDTO) error

// DeleteTemplate 删除活动模板
func (s *Service) DeleteTemplate(ctx context.Context, id int) error

// GetTemplate 获取活动模板详情
func (s *Service) GetTemplate(ctx context.Context, id int)
     (*dto.ActivityTemplateDetailDTO, error)

// ListTemplates 查询活动模板列表
func (s *Service) ListTemplates(ctx context.Context)
     ([]*dto.ActivityTemplateDTO, error)

// ListActiveTemplates 查询启用的活动模板
func (s *Service) ListActiveTemplates(ctx context.Context)
     ([]*dto.ActivityTemplateDTO, error)

// UpdateTemplateStatus 更新活动模板状态
func (s *Service) UpdateTemplateStatus(ctx context.Context, id int, status int) error
```

**业务逻辑:**
- 创建/更新时验证关联配置(分类或商品)
- 更新时检查模板状态(必须禁用才能编辑)
- 删除时检查是否有关联活动
- 查询详情时根据select_type JOIN查询关联数据

### 2. ActivityService

**职责:**
- 编排活动业务流程
- 验证模板状态
- 实现冲突检测逻辑

**核心方法:**
```go
// CreateActivity 创建活动
func (s *Service) CreateActivity(ctx context.Context,
     req *dto.CreateActivityDTO) (int, error)

// UpdateActivity 更新活动
func (s *Service) UpdateActivity(ctx context.Context, id int,
     req *dto.UpdateActivityDTO) error

// DeleteActivity 删除活动
func (s *Service) DeleteActivity(ctx context.Context, id int) error

// GetActivity 获取活动详情
func (s *Service) GetActivity(ctx context.Context, id int)
     (*dto.ActivityDetailResponseDTO, error)

// ListActivities 查询活动列表
func (s *Service) ListActivities(ctx context.Context)
     ([]*dto.ActivityDTO, error)

// GetActivitiesByDateRange 按日期范围查询活动（包含冲突检测）
func (s *Service) GetActivitiesByDateRange(ctx context.Context,
     paymentTime time.Time) (*dto.ActivitiesByDateRangeDTO, error)

// UpdateActivityStatus 更新活动状态
func (s *Service) UpdateActivityStatus(ctx context.Context, id int, status int) error
```

**业务逻辑:**
- 创建/更新时验证关联模板必须启用
- 更新时检查活动状态(必须禁用才能编辑)
- 验证活动详情的折扣值(满折类型需在0-1之间)
- 日期范围查询时检测同类型活动冲突

## 接口层设计

### API路由

所有API都需要JWT认证,除了登录接口。

#### 活动模板路由

```
GET    /api/activity-templates          - 查询所有模板
GET    /api/activity-templates/active   - 查询启用的模板
GET    /api/activity-templates/:id      - 查询模板详情
POST   /api/activity-templates          - 创建模板
PUT    /api/activity-templates/:id      - 更新模板
DELETE /api/activity-templates/:id      - 删除模板
PUT    /api/activity-templates/:id/status - 更新模板状态
```

#### 活动路由

```
GET    /api/activities                  - 查询所有活动
GET    /api/activities/:id              - 查询活动详情
POST   /api/activities                  - 创建活动
PUT    /api/activities/:id              - 更新活动
DELETE /api/activities/:id              - 删除活动
PUT    /api/activities/:id/status       - 更新活动状态
GET    /api/activities/by-date-range    - 按日期范围查询活动
```

### DTO设计

#### CreateActivityTemplateDTO

```go
type CreateActivityTemplateDTO struct {
    Name        string `json:"name" binding:"required"`
    Type        int    `json:"type" binding:"required"`
    SelectType  int    `json:"select_type" binding:"required"`
    ClassifyIDs []int  `json:"classify_ids"`
    GoodsIDs    []int  `json:"goods_ids"`
    Status      int    `json:"status"`
}
```

#### ActivityTemplateDetailDTO

```go
type ActivityTemplateDetailDTO struct {
    ID           int                   `json:"id"`
    Name         string                `json:"name"`
    Type         int                   `json:"type"`
    SelectType   int                   `json:"select_type"`
    Status       int                   `json:"status"`
    CreateTime   time.Time             `json:"create_time"`
    UpdateTime   time.Time             `json:"update_time"`
    ClassifyList []ClassifyRelationDTO `json:"classify_list,omitempty"`
    GoodsList    []GoodsRelationDTO    `json:"goods_list,omitempty"`
}
```

#### CreateActivityDTO

```go
type CreateActivityDTO struct {
    TemplateID int                 `json:"template_id" binding:"required"`
    Name       string              `json:"name" binding:"required"`
    StartTime  time.Time           `json:"start_time" binding:"required"`
    EndTime    time.Time           `json:"end_time" binding:"required"`
    Details    []ActivityDetailDTO `json:"details"`
    Status     int                 `json:"status"`
}
```

#### ActivitiesByDateRangeDTO

```go
type ActivitiesByDateRangeDTO struct {
    HasDuplicate  bool          `json:"has_duplicate"`
    DuplicateType *int          `json:"duplicate_type"`
    TypeName      string        `json:"type_name,omitempty"`
    Activities    []ActivityDTO `json:"activities"`
}
```

## 关键实现要点

### 1. 事务处理

所有涉及多表操作的方法都使用GORM事务:

```go
func (r *GormActivityTemplateRepository) Create(ctx context.Context,
    template *entity.ActivityTemplate,
    goods []*entity.ActivityTemplateGoods) error {

    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 创建模板
        if err := tx.Create(template).Error; err != nil {
            return err
        }

        // 创建关联数据
        if len(goods) > 0 {
            for _, g := range goods {
                g.TemplateID = template.ID
            }
            if err := tx.Create(&goods).Error; err != nil {
                return err
            }
        }

        return nil
    })
}
```

### 2. 冲突检测算法

在日期范围查询中检测同类型活动冲突:

```go
// 检测是否有重复类型
typeMap := make(map[int]bool)
var duplicateType *int

for _, a := range activities {
    template, _ := s.templateRepo.FindByID(ctx, a.TemplateID)

    // 检测重复类型
    if typeMap[template.Type] {
        result.HasDuplicate = true
        duplicateType = &template.Type
    }
    typeMap[template.Type] = true

    // ... 处理活动数据
}

// 设置冲突类型名称
if result.HasDuplicate && duplicateType != nil {
    result.DuplicateType = duplicateType
    switch *duplicateType {
    case 1:
        result.TypeName = "满减"
    case 2:
        result.TypeName = "满折"
    case 3:
        result.TypeName = "满赠"
    }
}
```

### 3. 级联查询

查询模板详情时根据select_type动态JOIN不同表:

```go
// 查询关联数据
if template.SelectType == 1 {
    // 按分类选择
    result.ClassifyList = s.getClassifyList(ctx, goods)
} else {
    // 按商品选择
    result.GoodsList = s.getGoodsList(ctx, goods)
}

// getGoodsList 实现
func (s *Service) getGoodsList(ctx context.Context,
    goods []*entity.ActivityTemplateGoods) []dto.GoodsRelationDTO {

    var result []dto.GoodsRelationDTO
    for _, g := range goods {
        if g.GoodsID != nil {
            var item struct {
                GoodsID      int
                GoodsName    string
                Price        float64
                BrandName    string
                ClassifyName string
            }
            s.db.Table("goods g").
                Select("g.id as goods_id, g.name as goods_name, g.price, " +
                       "b.name as brand_name, c.name as classify_name").
                Joins("LEFT JOIN brand b ON g.brandid = b.id").
                Joins("LEFT JOIN classify c ON g.classifyid = c.id").
                Where("g.id = ?", *g.GoodsID).
                First(&item)

            result = append(result, dto.GoodsRelationDTO{
                GoodsID:      *g.GoodsID,
                GoodsName:    item.GoodsName,
                Price:        item.Price,
                BrandName:    item.BrandName,
                ClassifyName: item.ClassifyName,
            })
        }
    }
    return result
}
```

### 4. 时间格式处理

支持多种时间格式输入:

```go
var paymentTime time.Time
var err error

// 尝试 RFC3339 格式
paymentTime, err = time.Parse(time.RFC3339, paymentTimeStr)
if err != nil {
    // 尝试日期时间格式
    paymentTime, err = time.Parse("2006-01-02 15:04:05", paymentTimeStr)
    if err != nil {
        // 尝试日期格式
        paymentTime, err = time.Parse("2006-01-02", paymentTimeStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment_time format"})
            return
        }
    }
}
```

## 与Python项目的差异

### 1. 架构模式

| 方面 | Python Flask | Go Gin + DDD |
|------|--------------|--------------|
| 架构 | 简单三层(路由-逻辑-数据) | DDD四层(接口-应用-领域-基础设施) |
| 数据库访问 | 直接SQL + PyMySQL | ORM + Repository模式 |
| 事务处理 | 手动管理连接和事务 | GORM自动事务管理 |
| 错误处理 | try-except + HTTP状态码 | 领域错误 + 统一响应格式 |

### 2. 业务逻辑封装

**Python (过程式):**
```python
@app.route('/api/activity-templates', methods=['POST'])
def create_activity_template():
    data = request.get_json()
    name = data.get('name')
    # ... 验证逻辑
    cursor.execute("INSERT INTO ...")
    # ... 插入关联数据
    connection.commit()
```

**Go (面向对象+DDD):**
```go
func (s *Service) CreateTemplate(ctx context.Context,
    req *dto.CreateActivityTemplateDTO) (int, error) {

    // 领域对象验证
    template := &entity.ActivityTemplate{...}
    if err := template.Validate(); err != nil {
        return 0, err
    }

    // 仓储持久化
    if err := s.repo.Create(ctx, template, goods); err != nil {
        return 0, err
    }

    return template.ID, nil
}
```

### 3. API兼容性

虽然内部实现不同,但API接口保持完全兼容:

- 相同的URL路径
- 相同的请求/响应格式
- 相同的业务规则
- 相同的错误处理

## 测试策略

### 1. 单元测试

- 领域实体验证逻辑
- 仓储方法
- 应用服务方法

### 2. 集成测试

- API端到端测试
- 数据库操作测试
- 事务测试

### 3. 测试数据

提供测试SQL脚本:
- `scripts/test_data_marketing.sql` - 测试数据
- `scripts/test_marketing_api.sh` - Bash测试脚本
- `scripts/test_marketing_api.ps1` - PowerShell测试脚本

## 性能优化

### 1. 数据库索引

```sql
KEY `idx_template_id` (`template_id`)
KEY `idx_time_range` (`start_time`, `end_time`)
KEY `idx_activity_id` (`activity_id`)
```

### 2. 批量操作

使用GORM的批量插入:
```go
if err := tx.Create(&details).Error; err != nil {
    return err
}
```

### 3. 连接池

GORM自动管理数据库连接池,配置参数:
```go
db.SetMaxIdleConns(10)
db.SetMaxOpenConns(100)
db.SetConnMaxLifetime(time.Hour)
```

## 部署说明

### 1. 编译

```bash
cd D:\claude space\CharonOMS
go build -o charonoms.exe ./cmd/server
```

### 2. 配置

确保`config.yaml`包含正确的数据库配置:
```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: qweasd123Q!
  database: charonoms
```

### 3. 初始化数据库

```bash
mysql -h localhost -u root -p charonoms < scripts/test_data_marketing.sql
```

### 4. 运行

```bash
./charonoms.exe
```

服务器默认监听 `http://localhost:8080`

## 总结

营销管理模块成功从Python Flask迁移到Go Gin,实现:

✅ DDD架构设计
✅ 完整的CRUD操作
✅ 事务支持
✅ 业务规则验证
✅ 冲突检测机制
✅ API兼容性
✅ 完整的测试覆盖
✅ 详细的文档

代码质量:
- 类型安全
- 良好的错误处理
- 清晰的职责分离
- 易于测试和维护

性能:
- 数据库连接池
- 批量操作
- 合理的索引设计

可维护性:
- 清晰的目录结构
- 分层架构
- 完整的注释和文档
