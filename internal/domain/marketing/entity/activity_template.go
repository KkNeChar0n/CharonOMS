package entity

import "time"

// ActivityTemplate 活动模板实体
type ActivityTemplate struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Type       int       `gorm:"not null" json:"type"`             // 1-满减 2-满折 3-满赠
	SelectType int       `gorm:"not null" json:"select_type"`      // 1-按分类 2-按商品
	Status     int       `gorm:"default:0" json:"status"`          // 0-启用 1-禁用
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`

	// 关联的分类列表（当 select_type=1 时）
	ClassifyList []map[string]interface{} `gorm:"-" json:"classify_list,omitempty"`
	// 关联的商品列表（当 select_type=2 时）
	GoodsList []map[string]interface{} `gorm:"-" json:"goods_list,omitempty"`
}

// TableName 指定表名
func (ActivityTemplate) TableName() string {
	return "activity_template"
}

// ActivityTemplateGoods 活动模板关联表（同时存储分类和商品关联）
type ActivityTemplateGoods struct {
	ID         uint  `gorm:"primaryKey" json:"id"`
	TemplateID uint  `gorm:"not null" json:"template_id"`
	ClassifyID *uint `gorm:"column:classify_id" json:"classify_id,omitempty"` // 分类ID（select_type=1时使用）
	GoodsID    *uint `gorm:"column:goods_id" json:"goods_id,omitempty"`       // 商品ID（select_type=2时使用）
}

// TableName 指定表名
func (ActivityTemplateGoods) TableName() string {
	return "activity_template_goods"
}
