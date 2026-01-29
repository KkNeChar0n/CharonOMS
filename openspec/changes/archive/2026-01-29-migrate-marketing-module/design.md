# 设计文档：营销管理模块迁移

## 架构决策

### 1. 模块划分策略

**决策**：营销管理作为独立的领域模块，包含活动模板和活动管理两个子模块

**理由**：
- 活动模板和活动实例是紧密关联的业务领域
- 未来可能扩展活动效果分析、自动化执行等功能
- 符合 DDD 的限界上下文概念

**结构**：
```
internal/domain/marketing/
├── entity/
│   ├── activity_template.go      # 活动模板实体
│   ├── activity.go                # 活动实体
│   └── activity_detail.go         # 活动细节实体
└── repository/
    ├── template_repository.go     # 活动模板仓储接口
    └── activity_repository.go     # 活动仓储接口
```

### 2. 活动模板关联设计

**问题**：活动模板可以按分类选择商品，也可以直接选择商品

**分析**：
- `select_type = 1`：按分类选择，关联表为 `activity_template_classify`
- `select_type = 2`：按商品选择，关联表为 `activity_template_goods`

**设计**：
```go
type ActivityTemplate struct {
    ID         uint
    Name       string
    Type       int       // 1-满减 2-满折 3-满赠
    SelectType int       // 1-按分类 2-按商品
    Status     int

    // 关联（二选一，根据 SelectType 动态加载）
    Classifies []*ActivityTemplateClassify `gorm:"foreignKey:TemplateID"`
    Goods      []*ActivityTemplateGoods    `gorm:"foreignKey:TemplateID"`
}

// 按分类关联
type ActivityTemplateClassify struct {
    ID         uint
    TemplateID uint
    ClassifyID uint
    Classify   *Classify `gorm:"foreignKey:ClassifyID"` // 预加载分类名称
}

// 按商品关联
type ActivityTemplateGoods struct {
    ID         uint
    TemplateID uint
    GoodsID    uint
    Goods      *Goods `gorm:"foreignKey:GoodsID"` // 预加载商品名称
}
```

**加载策略**：
- 根据 `select_type` 动态决定预加载哪个关联
- 创建/更新时根据 `select_type` 操作对应的关联表
- 使用事务确保模板和关联数据的一致性

### 3. 活动与活动细节设计

**业务逻辑**：
- 一个活动基于一个活动模板创建
- 活动有开始时间和结束时间
- 活动细节记录每个参与商品的具体优惠配置

**数据模型**：
```go
type Activity struct {
    ID           uint
    TemplateID   uint      // 关联的活动模板
    Name         string
    StartTime    time.Time
    EndTime      time.Time
    Status       int       // 0-启用 1-禁用
    CreateTime   time.Time
    UpdateTime   time.Time

    // 关联
    Template *ActivityTemplate  `gorm:"foreignKey:TemplateID"`
    Details  []*ActivityDetail  `gorm:"foreignKey:ActivityID"`
}

type ActivityDetail struct {
    ID         uint
    ActivityID uint
    GoodsID    uint
    Discount   float64   // 折扣值：满减活动=优惠金额，满折活动=折扣率

    // 关联
    Goods *Goods `gorm:"foreignKey:GoodsID"`
}
```

**创建流程**：
1. 验证活动模板存在且启用
2. 验证开始时间 < 结束时间
3. 使用事务创建活动和活动细节
4. 根据模板类型和 select_type，验证参与商品的合法性

### 4. 数据库表设计

**需要创建的表**：

```sql
-- 活动表
CREATE TABLE `activity` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `template_id` int NOT NULL COMMENT '关联的活动模板ID',
  `name` varchar(100) NOT NULL COMMENT '活动名称',
  `start_time` timestamp NOT NULL COMMENT '开始时间',
  `end_time` timestamp NOT NULL COMMENT '结束时间',
  `status` int DEFAULT '0' COMMENT '状态：0-启用 1-禁用',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_status` (`status`),
  KEY `idx_time` (`start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动表';

-- 活动细节表
CREATE TABLE `activity_detail` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '细节ID',
  `activity_id` int NOT NULL COMMENT '关联的活动ID',
  `goods_id` int NOT NULL COMMENT '商品ID',
  `discount` decimal(10,2) NOT NULL COMMENT '折扣值',
  PRIMARY KEY (`id`),
  KEY `idx_activity_id` (`activity_id`),
  KEY `idx_goods_id` (`goods_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活动细节表';
```

### 5. API 响应格式

**要求**：保持与 Python 版本的响应格式一致

**标准格式**：
```json
{
  "activity_templates": [...],
  "activities": [...]
}
```

**注意事项**：
- 不使用统一的 `data` 包装层（与其他模块不同）
- 列表接口直接返回数组字段
- 详情接口返回单个对象（不包装）
- 错误响应使用标准格式：`{"error": "错误信息"}`

**活动列表响应示例**：
```json
{
  "activities": [
    {
      "id": 1,
      "template_id": 1,
      "template_name": "春季促销",
      "template_type": 1,
      "name": "春季大促销活动",
      "start_time": "2026-03-01T00:00:00Z",
      "end_time": "2026-03-31T23:59:59Z",
      "status": 0
    }
  ]
}
```

**活动详情响应示例**：
```json
{
  "id": 1,
  "template_id": 1,
  "template_name": "春季促销",
  "template_type": 1,
  "name": "春季大促销活动",
  "start_time": "2026-03-01T00:00:00Z",
  "end_time": "2026-03-31T23:59:59Z",
  "status": 0,
  "details": [
    {
      "id": 1,
      "activity_id": 1,
      "goods_id": 101,
      "goods_name": "篮球",
      "discount": 50.00
    }
  ]
}
```

### 6. 状态管理模式

**决策**：采用统一的状态字段和操作接口

**模式**：
- 所有实体使用 `status` 字段：0-启用，1-禁用
- 提供独立的状态更新接口（`PUT /:id/status`）
- 活动模板：禁用状态下才能编辑
- 活动：禁用状态下才能编辑

**业务规则**：
1. 新创建的活动模板默认为禁用状态
2. 新创建的活动默认为禁用状态
3. 只有启用的活动模板才能用于创建活动
4. 活动时间范围验证：start_time < end_time

### 7. 筛选参数处理

**支持的筛选参数**：

活动模板：
- `id` - 精确匹配
- `name` - 模糊匹配（LIKE）
- `type` - 1/2/3（满减/满折/满赠）
- `status` - 0/1

活动：
- `id` - 精确匹配
- `name` - 模糊匹配（LIKE）
- `template_id` - 精确匹配
- `status` - 0/1

**实现方式**：
```go
// 动态构建查询条件
query := r.db.WithContext(ctx)

if id != 0 {
    query = query.Where("id = ?", id)
}
if name != "" {
    query = query.Where("name LIKE ?", "%"+name+"%")
}
if templateID != 0 {
    query = query.Where("template_id = ?", templateID)
}
if status != nil {
    query = query.Where("status = ?", *status)
}
```

## 数据库表关系

```
activity_template (活动模板)
    ↓ 1:N (when select_type=1)
activity_template_classify (模板-分类关联)
    ↓ N:1
classify (商品分类)

activity_template (活动模板)
    ↓ 1:N (when select_type=2)
activity_template_goods (模板-商品关联)
    ↓ N:1
goods (商品)

activity_template (活动模板)
    ↓ 1:N
activity (活动)
    ↓ 1:N
activity_detail (活动细节)
    ↓ N:1
goods (商品)
```

## 技术约束

1. **GORM 关联加载**：
   - 使用 `Preload()` 预加载关联数据
   - 避免 N+1 查询问题
   - 活动列表需要 Preload Template（获取 template_name 和 template_type）

2. **事务处理**：
   - 创建活动模板时需要事务（模板 + 关联）
   - 更新活动模板时需要事务（先删除旧关联，再创建新关联）
   - 创建活动时需要事务（活动 + 细节）
   - 更新活动时需要事务（活动 + 细节）

3. **字段映射**：
   - 统一使用 `create_time`、`update_time`
   - GORM 标签：`gorm:"column:create_time"`

4. **数据验证**：
   - 活动模板：名称不能为空，类型必须为 1/2/3
   - 活动：开始时间必须早于结束时间
   - 活动细节：折扣值必须大于0

## 测试策略

### 单元测试
- Domain 层：实体创建和验证逻辑
- Application 层：业务用例和筛选逻辑
- Repository 层：数据库查询和关联加载

### 集成测试
- 测试活动模板创建（按分类和按商品）
- 测试活动创建（包含活动细节）
- 测试状态更新逻辑
- 测试时间范围验证

### 兼容性测试
- 对比 Python 版本的 API 响应格式
- 验证前端页面调用是否正常
- 测试权限控制是否生效
