package entity

import "time"

// ActivityTemplate 活动模板实体
type ActivityTemplate struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Type       int       `gorm:"not null" json:"type"`             // 1-满减 2-满折 3-满赠
	SelectType int       `gorm:"not null" json:"select_type"`      // 1-按分类 2-按商品
	Status     int       `gorm:"default:1" json:"status"`          // 0-启用 1-禁用
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`

	// 关联（根据 SelectType 动态加载其中之一）
	Classifies []*ActivityTemplateClassify `gorm:"foreignKey:TemplateID" json:"classifies,omitempty"`
	Goods      []*ActivityTemplateGoods    `gorm:"foreignKey:TemplateID" json:"goods,omitempty"`
}

// TableName 指定表名
func (ActivityTemplate) TableName() string {
	return "activity_template"
}

// ActivityTemplateClassify 活动模板-分类关联实体
type ActivityTemplateClassify struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	TemplateID   uint   `gorm:"not null" json:"template_id"`
	ClassifyID   uint   `gorm:"not null" json:"classify_id"`
	ClassifyName string `gorm:"-" json:"classify_name,omitempty"` // 关联查询时填充
}

// TableName 指定表名
func (ActivityTemplateClassify) TableName() string {
	return "activity_template_classify"
}

// ActivityTemplateGoods 活动模板-商品关联实体
type ActivityTemplateGoods struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TemplateID uint   `gorm:"not null" json:"template_id"`
	GoodsID    uint   `gorm:"not null" json:"goods_id"`
	GoodsName  string `gorm:"-" json:"goods_name,omitempty"` // 关联查询时填充
}

// TableName 指定表名
func (ActivityTemplateGoods) TableName() string {
	return "activity_template_goods"
}
