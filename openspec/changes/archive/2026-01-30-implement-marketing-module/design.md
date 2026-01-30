# 设计文档：营销管理模块实现

## 架构概览

营销管理模块采用 DDD（领域驱动设计）分层架构，与项目现有模块保持一致。模块分为两个主要子领域：

1. **活动模板领域** (Activity Template Domain)
2. **活动领域** (Activity Domain)

## 领域模型

### 1. Activity Template 领域

#### 实体关系

```
ActivityTemplate (活动模板)
├── id: int
├── name: string
├── type: int (1=满减, 2=满折, 3=满赠)
├── select_type: int (1=按分类, 2=按商品)
├── status: int (0=启用, 1=禁用)
├── create_time: timestamp
├── update_time: timestamp
└── relations: []ActivityTemplateGoods

ActivityTemplateGoods (活动模板商品/分类关联)
├── id: int
├── template_id: int
├── goods_id: *int (当 select_type=2 时使用)
└── classify_id: *int (当 select_type=1 时使用)
```

#### 领域规则

1. **创建规则**
   - 模板名称不能为空
   - type 必须为 1/2/3
   - select_type 必须为 1/2
   - 当 select_type=1 时，必须至少关联一个分类
   - 当 select_type=2 时，必须至少关联一个商品
   - 新创建的模板默认状态为禁用（status=1）

2. **更新规则**
   - 只有禁用状态（status=1）的模板才能编辑
   - 更新时先删除所有旧的关联记录，再创建新的关联记录
   - 模板名称不能为空

3. **删除规则**
   - 删除前必须检查是否有关联的活动
   - 如果存在关联的活动，禁止删除
   - 删除模板时，关联的 ActivityTemplateGoods 记录会级联删除（数据库外键）

4. **状态更新规则**
   - 可以在启用和禁用之间切换
   - 启用状态的模板才能被用于创建活动

### 2. Activity 领域

#### 实体关系

```
Activity (活动)
├── id: int
├── name: string
├── template_id: int (关联 ActivityTemplate)
├── start_time: datetime
├── end_time: datetime
├── status: int (0=启用, 1=禁用)
├── create_time: timestamp
└── details: []ActivityDetail

ActivityDetail (活动详情 - 满折规则)
├── id: int
├── activity_id: int
├── threshold_amount: decimal(10,2) (门槛金额)
└── discount_value: decimal(10,2) (折扣值)
```

#### 领域规则

1. **创建规则**
   - 活动名称不能为空
   - 关联的模板必须存在且状态为启用
   - 开始时间必须早于结束时间
   - 满折活动（template.type=2）必须至少有一条 ActivityDetail
   - 新创建的活动默认状态为禁用（status=1）

2. **更新规则**
   - 只有禁用状态（status=1）的活动才能编辑
   - 更新时先删除所有旧的详情记录，再创建新的详情记录
   - 验证规则与创建相同

3. **删除规则**
   - 直接删除活动
   - ActivityDetail 记录会级联删除（数据库外键）
   - orders_activity 关联记录会级联删除（数据库外键）

4. **状态更新规则**
   - 可以在启用和禁用之间切换
   - 启用状态的活动才会被应用到订单

5. **按日期范围查询规则**
   - 查询在指定时间范围内的启用活动
   - WHERE start_time <= payment_time AND end_time >= payment_time AND status = 0
   - 检测同类型活动冲突：同一时间范围内不能有两个相同类型的活动同时启用
   - 仅满折活动（type=2）需要加载 details

## 数据访问层设计

### Repository 接口

#### ActivityTemplateRepository

```go
type ActivityTemplateRepository interface {
    // 基础 CRUD
    Create(ctx context.Context, template *entity.ActivityTemplate) error
    Update(ctx context.Context, template *entity.ActivityTemplate) error
    Delete(ctx context.Context, id int) error
    FindByID(ctx context.Context, id int) (*entity.ActivityTemplate, error)
    List(ctx context.Context) ([]*entity.ActivityTemplate, error)

    // 业务查询
    FindActiveTemplates(ctx context.Context) ([]*entity.ActivityTemplate, error)
    CountRelatedActivities(ctx context.Context, templateID int) (int, error)

    // 关联操作
    UpdateTemplateGoods(ctx context.Context, templateID int, goods []*entity.ActivityTemplateGoods) error
}
```

#### ActivityRepository

```go
type ActivityRepository interface {
    // 基础 CRUD
    Create(ctx context.Context, activity *entity.Activity) error
    Update(ctx context.Context, activity *entity.Activity) error
    Delete(ctx context.Context, id int) error
    FindByID(ctx context.Context, id int) (*entity.Activity, error)
    List(ctx context.Context) ([]*entity.Activity, error)

    // 业务查询
    FindByDateRange(ctx context.Context, paymentTime time.Time) ([]*entity.Activity, error)

    // 关联操作
    UpdateActivityDetails(ctx context.Context, activityID int, details []*entity.ActivityDetail) error
}
```

### 事务处理

#### 需要事务的操作

1. **创建活动模板**
   - 插入 activity_template
   - 批量插入 activity_template_goods

2. **更新活动模板**
   - 更新 activity_template
   - 删除旧的 activity_template_goods
   - 批量插入新的 activity_template_goods

3. **创建活动**
   - 插入 activity
   - 批量插入 activity_detail（如果有）

4. **更新活动**
   - 更新 activity
   - 删除旧的 activity_detail
   - 批量插入新的 activity_detail（如果有）

#### 事务实现模式

使用 GORM 的事务 API：

```go
func (s *ActivityService) Create(ctx context.Context, dto *dto.CreateActivityDTO) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建活动主记录
        if err := tx.Create(&activity).Error; err != nil {
            return err
        }

        // 2. 创建活动详情
        if len(details) > 0 {
            if err := tx.Create(&details).Error; err != nil {
                return err
            }
        }

        return nil
    })
}
```

## 应用服务层设计

### ActivityTemplateService

职责：
- 编排活动模板的业务流程
- DTO 与 Entity 的转换
- 调用 Repository 完成数据持久化
- 执行业务验证

主要方法：
```go
type ActivityTemplateService interface {
    CreateTemplate(ctx context.Context, dto *dto.CreateTemplateDTO) (*dto.TemplateDTO, error)
    UpdateTemplate(ctx context.Context, id int, dto *dto.UpdateTemplateDTO) error
    DeleteTemplate(ctx context.Context, id int) error
    GetTemplate(ctx context.Context, id int) (*dto.TemplateDetailDTO, error)
    ListTemplates(ctx context.Context) ([]*dto.TemplateDTO, error)
    ListActiveTemplates(ctx context.Context) ([]*dto.TemplateDTO, error)
    UpdateTemplateStatus(ctx context.Context, id int, status int) error
}
```

### ActivityService

职责：
- 编排活动的业务流程
- DTO 与 Entity 的转换
- 调用 Repository 完成数据持久化
- 执行业务验证
- 处理活动时间范围查询和冲突检测

主要方法：
```go
type ActivityService interface {
    CreateActivity(ctx context.Context, dto *dto.CreateActivityDTO) (*dto.ActivityDTO, error)
    UpdateActivity(ctx context.Context, id int, dto *dto.UpdateActivityDTO) error
    DeleteActivity(ctx context.Context, id int) error
    GetActivity(ctx context.Context, id int) (*dto.ActivityDetailDTO, error)
    ListActivities(ctx context.Context) ([]*dto.ActivityDTO, error)
    GetActivitiesByDateRange(ctx context.Context, paymentTime time.Time) (*dto.ActivitiesByDateRangeDTO, error)
    UpdateActivityStatus(ctx context.Context, id int, status int) error
}
```

## HTTP 接口层设计

### 响应格式

为了与 Python 版本保持完全一致，所有接口返回格式如下：

**成功响应**:
```json
{
  "field_name": "value"
}
```

**错误响应**:
```json
{
  "error": "错误消息"
}
```

**注意**: Go 项目的标准响应格式为 `{"code": 0, "message": "", "data": {}}`，但为了与 Python 版本兼容，营销管理模块的接口将使用 Python 的响应格式。

### 错误处理

HTTP 状态码与 Python 版本保持一致：

- 200: 操作成功
- 201: 创建成功
- 400: 请求参数错误
- 404: 资源不存在
- 500: 服务器内部错误

错误消息使用中文，与 Python 版本的错误消息完全一致。

### Handler 实现模式

```go
func (h *ActivityTemplateHandler) CreateTemplate(c *gin.Context) {
    var req dto.CreateTemplateDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "请求参数错误"})
        return
    }

    result, err := h.service.CreateTemplate(c.Request.Context(), &req)
    if err != nil {
        // 根据错误类型返回不同的状态码
        c.JSON(determineStatusCode(err), gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, result)
}
```

## 数据传输对象 (DTO)

### ActivityTemplate DTOs

```go
// 创建活动模板请求
type CreateTemplateDTO struct {
    Name        string `json:"name" binding:"required"`
    Type        int    `json:"type" binding:"required,oneof=1 2 3"`
    SelectType  int    `json:"select_type" binding:"required,oneof=1 2"`
    ClassifyIDs []int  `json:"classify_ids"`
    GoodsIDs    []int  `json:"goods_ids"`
    Status      int    `json:"status"`
}

// 活动模板列表响应
type TemplateDTO struct {
    ID         int       `json:"id"`
    Name       string    `json:"name"`
    Type       int       `json:"type"`
    SelectType int       `json:"select_type"`
    Status     int       `json:"status"`
    CreateTime time.Time `json:"create_time"`
    UpdateTime time.Time `json:"update_time"`
}

// 活动模板详情响应
type TemplateDetailDTO struct {
    TemplateDTO
    ClassifyList []ClassifyRelationDTO `json:"classify_list,omitempty"`
    GoodsList    []GoodsRelationDTO    `json:"goods_list,omitempty"`
}
```

### Activity DTOs

```go
// 创建活动请求
type CreateActivityDTO struct {
    TemplateID int                    `json:"template_id" binding:"required"`
    Name       string                 `json:"name" binding:"required"`
    StartTime  time.Time              `json:"start_time" binding:"required"`
    EndTime    time.Time              `json:"end_time" binding:"required"`
    Details    []ActivityDetailDTO    `json:"details"`
    Status     int                    `json:"status"`
}

// 活动详情
type ActivityDetailDTO struct {
    ID              int     `json:"id,omitempty"`
    ActivityID      int     `json:"activity_id,omitempty"`
    ThresholdAmount float64 `json:"threshold_amount" binding:"required"`
    DiscountValue   float64 `json:"discount_value" binding:"required"`
}

// 活动列表响应
type ActivityDTO struct {
    ID           int       `json:"id"`
    TemplateID   int       `json:"template_id"`
    TemplateName string    `json:"template_name"`
    TemplateType int       `json:"template_type"`
    Name         string    `json:"name"`
    StartTime    time.Time `json:"start_time"`
    EndTime      time.Time `json:"end_time"`
    Status       int       `json:"status"`
    CreateTime   time.Time `json:"create_time"`
    UpdateTime   time.Time `json:"update_time"`
}

// 按日期范围查询活动响应
type ActivitiesByDateRangeDTO struct {
    HasDuplicate  bool          `json:"has_duplicate"`
    DuplicateType *int          `json:"duplicate_type"`
    TypeName      string        `json:"type_name,omitempty"`
    Activities    []ActivityDTO `json:"activities"`
}
```

## 关键业务流程

### 1. 创建活动模板流程

```
1. Handler 接收请求，验证 DTO
2. 检查 select_type 与 classify_ids/goods_ids 的一致性
3. 调用 Service.CreateTemplate
4. Service 验证业务规则
5. 创建 ActivityTemplate 实体
6. 创建 ActivityTemplateGoods 实体列表
7. 调用 Repository.Create 在事务中：
   - 插入 activity_template
   - 批量插入 activity_template_goods
8. 返回创建的模板 ID
```

### 2. 按日期范围查询活动流程

```
1. Handler 接收 payment_time 参数
2. 调用 Service.GetActivitiesByDateRange
3. Service 调用 Repository.FindByDateRange 查询活动
4. 检测同类型活动冲突：
   - 遍历活动列表
   - 统计每种类型的活动数量
   - 如果某种类型数量 > 1，返回冲突信息
5. 如果无冲突，为满折活动加载 details
6. 返回活动列表和冲突信息
```

### 3. 更新活动模板流程

```
1. Handler 接收请求
2. Service 查询模板是否存在
3. 检查模板状态是否为禁用
4. 如果已启用，返回错误
5. 在事务中：
   - 更新 activity_template
   - 删除旧的 activity_template_goods
   - 插入新的 activity_template_goods
6. 返回成功
```

## 测试策略

### 单元测试

#### 领域层测试
- Entity 的验证逻辑测试
- 领域服务的业务规则测试

#### 应用层测试
- Service 方法的业务流程测试
- DTO 转换逻辑测试
- 使用 Mock Repository

#### 接口层测试
- Handler 的请求参数验证测试
- 响应格式测试
- 使用 Mock Service

### 集成测试

- 使用测试数据库
- 测试完整的 API 调用流程
- 验证数据库状态变化
- 测试事务回滚

### 关键测试场景

1. **活动模板**
   - 创建按分类选择的模板
   - 创建按商品选择的模板
   - 更新启用状态的模板（应失败）
   - 删除有关联活动的模板（应失败）
   - 状态更新

2. **活动**
   - 创建满折活动（with details）
   - 创建满减活动（without details）
   - 使用未启用模板创建活动（应失败）
   - 结束时间早于开始时间（应失败）
   - 按日期范围查询（有冲突）
   - 按日期范围查询（无冲突）

## 性能考虑

### 查询优化

1. **活动模板详情查询**
   - 使用 JOIN 查询关联的商品/分类信息
   - 避免 N+1 查询问题

2. **活动列表查询**
   - 使用 JOIN 查询关联的模板信息
   - 按 ID 降序排序

3. **按日期范围查询活动**
   - 在 activity 表的 start_time, end_time, status 上创建复合索引
   - 仅在需要时加载 details（满折活动）

### 缓存策略

暂不实现缓存，后续可以考虑：
- 缓存启用的活动模板列表
- 缓存当前生效的活动列表

## 日志和监控

### 关键日志点

1. 业务操作日志
   - 创建/更新/删除活动模板
   - 创建/更新/删除活动
   - 状态变更

2. 错误日志
   - 数据库操作错误
   - 业务验证失败
   - 请求参数错误

3. 性能日志
   - 慢查询（> 100ms）
   - 大批量操作

### 日志格式

使用项目统一的 Zap 日志格式：
```go
logger.Info("创建活动模板",
    zap.String("name", template.Name),
    zap.Int("type", template.Type),
    zap.Int("id", template.ID),
)
```

## 兼容性检查清单

- [ ] API 路由与 Python 版本完全一致
- [ ] 请求 DTO 字段名与 Python 版本一致（snake_case）
- [ ] 响应 DTO 字段名与 Python 版本一致（snake_case）
- [ ] HTTP 状态码与 Python 版本一致
- [ ] 错误消息文本与 Python 版本一致
- [ ] 业务验证规则与 Python 版本一致
- [ ] 数据库字段映射正确
- [ ] 时间字段格式一致
- [ ] 列表排序方式一致（ID 降序）
