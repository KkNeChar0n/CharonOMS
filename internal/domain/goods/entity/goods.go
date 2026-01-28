package entity

import "time"

// Goods 商品实体
type Goods struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(200);not null" json:"name"`
	BrandID    int       `gorm:"column:brandid;not null" json:"brandid"`
	ClassifyID int       `gorm:"column:classifyid;not null" json:"classifyid"`
	IsGroup    *int      `gorm:"column:isgroup;not null" json:"isgroup"` // 0=套餐，1=单品
	Price      float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Status     int       `gorm:"column:status;default:0" json:"status"` // 0=启用，1=禁用
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Goods) TableName() string {
	return "goods"
}

// GoodsAttributeValue 商品与属性值的关联实体
type GoodsAttributeValue struct {
	ID               int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID          int       `gorm:"column:goodsid;not null" json:"goodsid"`
	AttributeValueID int       `gorm:"column:attributevalueid;not null" json:"attributevalueid"`
	CreateTime       time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (GoodsAttributeValue) TableName() string {
	return "goods_attributevalue"
}

// GoodsGoods 商品组合关系实体（子商品与父商品的关联）
type GoodsGoods struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID    int       `gorm:"column:goodsid;not null" json:"goodsid"`     // 子商品ID
	ParentsID  int       `gorm:"column:parentsid;not null" json:"parentsid"` // 父商品ID（组合商品）
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (GoodsGoods) TableName() string {
	return "goods_goods"
}
