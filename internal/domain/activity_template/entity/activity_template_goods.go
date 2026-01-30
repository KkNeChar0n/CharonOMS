package entity

// ActivityTemplateGoods 活动模板商品/分类关联实体
type ActivityTemplateGoods struct {
	ID         int  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateID int  `gorm:"column:template_id;not null" json:"template_id"`
	GoodsID    *int `gorm:"column:goods_id" json:"goods_id,omitempty"`       // 当 select_type=2 时使用
	ClassifyID *int `gorm:"column:classify_id" json:"classify_id,omitempty"` // 当 select_type=1 时使用
}

// TableName 指定表名
func (ActivityTemplateGoods) TableName() string {
	return "activity_template_goods"
}
