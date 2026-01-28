# 设计文档：商品管理模块

## 架构设计

### 模块划分

按 DDD 四层架构，将商品体系分为四个独立子模块：

```
internal/
├── domain/
│   ├── brand/
│   │   ├── entity/brand.go
│   │   └── repository/brand_repository.go
│   ├── classify/
│   │   ├── entity/classify.go
│   │   └── repository/classify_repository.go
│   ├── attribute/
│   │   ├── entity/attribute.go
│   │   ├── entity/attribute_value.go
│   │   └── repository/attribute_repository.go
│   └── goods/
│       ├── entity/goods.go
│       └── repository/goods_repository.go
├── infrastructure/persistence/mysql/
│   ├── brand/brand_repository_impl.go
│   ├── classify/classify_repository_impl.go
│   ├── attribute/attribute_repository_impl.go
│   └── goods/goods_repository_impl.go
├── application/service/
│   ├── brand/brand_service.go + dto.go
│   ├── classify/classify_service.go + dto.go
│   ├── attribute/attribute_service.go + dto.go
│   └── goods/goods_service.go + dto.go
└── interfaces/http/handler/
    ├── brand/brand_handler.go
    ├── classify/classify_handler.go
    ├── attribute/attribute_handler.go
    └── goods/goods_handler.go
```

### 实体设计

#### Brand
```go
type Brand struct {
    ID         int       `gorm:"primaryKey;column:id"`
    Name       string    `gorm:"column:name;uniqueIndex"`
    Status     int       `gorm:"column:status;default:0"`
    CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
    UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}
```

#### Classify
```go
type Classify struct {
    ID         int        `gorm:"primaryKey;column:id"`
    Name       string     `gorm:"column:name"`
    Level      int        `gorm:"column:level;default:0"`      // 0=一级, 1=二级
    ParentID   *int       `gorm:"column:parentid"`             // 二级时指向父级ID
    Status     int        `gorm:"column:status;default:0"`
    CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt  time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}
```

#### Attribute
```go
type Attribute struct {
    ID         int       `gorm:"primaryKey;column:id"`
    Name       string    `gorm:"column:name"`
    Classify   int       `gorm:"column:classify;default:0"`   // 0=属性, 1=规格
    Status     int       `gorm:"column:status;default:0"`
    CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
    UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

type AttributeValue struct {
    ID          int    `gorm:"primaryKey;column:id"`
    Name        string `gorm:"column:name"`
    AttributeID int    `gorm:"column:attributeid"`
}
```

#### Goods
```go
type Goods struct {
    ID          int            `gorm:"primaryKey;column:id"`
    Name        string         `gorm:"column:name"`
    BrandID     int            `gorm:"column:brandid"`
    ClassifyID  int            `gorm:"column:classifyid"`
    IsGroup     int            `gorm:"column:isgroup;default:1"` // 0=套餐, 1=单品
    Price       float64        `gorm:"column:price;type:decimal(10,2)"`
    Status      int            `gorm:"column:status;default:0"`
    CreateTime  time.Time      `gorm:"column:create_time;autoCreateTime"`
    UpdateTime  time.Time      `gorm:"column:update_time;autoUpdateTime"`
}
```

### 仓储接口设计

#### BrandRepository
```go
type BrandRepository interface {
    GetAll() ([]Brand, error)
    GetActive() ([]Brand, error)
    GetByID(id int) (*Brand, error)
    GetByName(name string) (*Brand, error)
    Create(name string) (int, error)
    Update(id int, name string) error
    UpdateStatus(id int, status int) error
}
```

#### ClassifyRepository
```go
type ClassifyRepository interface {
    GetAll() ([]map[string]interface{}, error)           // 含parent_name
    GetParents() ([]Classify, error)                     // 启用的一级分类
    GetActive() ([]Classify, error)                      // 启用的二级分类
    GetByID(id int) (*Classify, error)
    CheckNameUnique(name string, level int, parentID *int, excludeID int) (bool, error)
    Create(name string, level int, parentID *int) (int, error)
    Update(id int, name string, level int, parentID *int) error
    UpdateStatus(id int, status int) error
}
```

#### AttributeRepository
```go
type AttributeRepository interface {
    GetAll() ([]map[string]interface{}, error)           // 含value_count
    GetActive() ([]map[string]interface{}, error)        // 启用属性含值列表
    GetByID(id int) (*Attribute, error)
    Create(name string, classify int) (int, error)
    Update(id int, name string, classify int) error
    UpdateStatus(id int, status int) error
    GetValues(attributeID int) ([]AttributeValue, error)
    SaveValues(attributeID int, values []string) error   // 全量替换
}
```

#### GoodsRepository
```go
type GoodsRepository interface {
    GetList(classifyID *int, status *int) ([]map[string]interface{}, error)  // 含brand_name, classify_name, attributes
    GetByID(id int) (map[string]interface{}, error)                           // 含attributevalue_ids, included_goods_ids
    GetActiveForOrder() ([]map[string]interface{}, error)                     // 含total_price
    GetAvailableForCombo(excludeID *int) ([]map[string]interface{}, error)
    GetIncludedGoods(goodsID int) ([]map[string]interface{}, error)
    GetTotalPrice(goodsID int) (map[string]interface{}, error)
    Create(name string, brandID, classifyID, isGroup int, price float64, attributeValueIDs []int, includedGoodsIDs []int) (int, error)
    Update(id int, name string, brandID, classifyID int, price float64, attributeValueIDs []int, includedGoodsIDs []int) error
    UpdateStatus(id int, status int) error
}
```

### 关键设计决策

1. **属性字符串格式**：商品列表中 `attributes` 字段为 `"属性名:值1/值2,属性名2:值3"` 格式，在仓储层构建。

2. **全量替换策略**：
   - 属性值保存：DELETE + INSERT（替换属性下所有值）
   - 商品属性绑定：DELETE goods_attributevalue + INSERT
   - 组合商品子商品：DELETE goods_goods + INSERT

3. **组合商品总价**：在仓储层通过子查询计算 `SUM(sub_goods.price)`。

4. **分类两级约束**：
   - level=0 名称在全部一级分类中唯一
   - level=1 名称在同 parentid 下唯一

5. **路由顺序**：静态路径（如 `/goods/active-for-order`）必须在参数路径（`/goods/:id`）之前注册。

6. **Vue 前端字段类型**：与合同模块同样，数字字段可能以字符串形式传入，DTO 使用 `interface{}` 接收关键数字字段。

### 路由注册方案

```go
// Brand
brands := authorized.Group("/brands")
{
    brands.GET("", brandHdl.GetBrands)
    brands.GET("/active", brandHdl.GetActiveBrands)
    brands.POST("", brandHdl.CreateBrand)
    brands.PUT("/:id", brandHdl.UpdateBrand)
    brands.PUT("/:id/status", brandHdl.UpdateBrandStatus)
}

// Classify
classifies := authorized.Group("/classifies")
{
    classifies.GET("", classifyHdl.GetClassifies)
    classifies.GET("/parents", classifyHdl.GetParents)
    classifies.GET("/active", classifyHdl.GetActiveClassifies)
    classifies.POST("", classifyHdl.CreateClassify)
    classifies.PUT("/:id", classifyHdl.UpdateClassify)
    classifies.PUT("/:id/status", classifyHdl.UpdateClassifyStatus)
}

// Attribute
attributes := authorized.Group("/attributes")
{
    attributes.GET("", attributeHdl.GetAttributes)
    attributes.GET("/active", attributeHdl.GetActiveAttributes)
    attributes.POST("", attributeHdl.CreateAttribute)
    attributes.PUT("/:id", attributeHdl.UpdateAttribute)
    attributes.PUT("/:id/status", attributeHdl.UpdateAttributeStatus)
    attributes.GET("/:id/values", attributeHdl.GetAttributeValues)
    attributes.POST("/:id/values", attributeHdl.SaveAttributeValues)
}

// Goods (静态路径在参数路径之前)
goods := authorized.Group("/goods")
{
    goods.GET("", goodsHdl.GetGoods)
    goods.GET("/active-for-order", goodsHdl.GetActiveForOrder)
    goods.GET("/available-for-combo", goodsHdl.GetAvailableForCombo)
    goods.POST("", goodsHdl.CreateGoods)
    goods.GET("/:id", goodsHdl.GetGoodsByID)
    goods.GET("/:id/included-goods", goodsHdl.GetIncludedGoods)
    goods.GET("/:id/total-price", goodsHdl.GetTotalPrice)
    goods.PUT("/:id", goodsHdl.UpdateGoods)
    goods.PUT("/:id/status", goodsHdl.UpdateGoodsStatus)
}
```

### 错误响应映射

| Python 错误消息 | HTTP 状态码 | 场景 |
|----------------|-------------|------|
| "品牌名称不能为空" | 400 | 新增/编辑品牌时名称为空 |
| "该品牌名称已存在" | 400 | 品牌名称重复 |
| "品牌不存在" | 404 | 编辑/状态更新时ID不存在 |
| "品牌状态更新成功" | 200 | 状态更新成功 |
| "名称和级别不能为空" | 400 | 新增/编辑分类时必填项为空 |
| "级别值必须为0或1" | 400 | 分类level非法值 |
| "二级类型必须选择父级类型" | 400 | level=1 无 parent_id |
| "该一级类型名称已存在" | 400 | 一级分类名称重复 |
| "该父级类型下已存在同名的二级类型" | 400 | 同父级二级分类名重复 |
| "类型不存在" | 404 | 编辑/状态更新时ID不存在 |
| "名称和分类不能为空" | 400 | 新增/编辑属性时必填项为空 |
| "分类值必须为0或1" | 400 | 属性classify非法值 |
| "属性不存在" | 404 | 编辑/状态更新时ID不存在 |
| "至少需要填入一条属性值" | 400 | 保存属性值时值数组为空 |
| "属性值不能为空" | 400 | 值数组中有空字符串 |
| "商品信息不完整" | 400 | 新增/编辑商品时必填项缺失 |
| "组合商品必须至少包含一个子商品" | 400 | isgroup=0 无子商品 |
| "商品不存在" | 404 | 获取/编辑/状态更新时ID不存在 |
| "状态不能为空" | 400 | 状态更新时status为空 |
| "状态值必须为0或1" | 400 | 状态非法值（分类和属性） |
