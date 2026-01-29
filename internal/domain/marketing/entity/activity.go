package entity

import "time"

// Activity 活动实体
type Activity struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TemplateID uint      `gorm:"not null" json:"template_id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	StartTime  time.Time `gorm:"not null;column:start_time" json:"start_time"`
	EndTime    time.Time `gorm:"not null;column:end_time" json:"end_time"`
	Status     int       `gorm:"default:1" json:"status"` // 0-启用 1-禁用
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`

	// 关联
	Template     *ActivityTemplate  `gorm:"foreignKey:TemplateID" json:"-"`
	TemplateName string             `gorm:"-" json:"template_name,omitempty"` // 从 Template 填充
	TemplateType int                `gorm:"-" json:"template_type,omitempty"` // 从 Template 填充
	Details      []*ActivityDetail  `gorm:"foreignKey:ActivityID" json:"details,omitempty"`
}

// TableName 指定表名
func (Activity) TableName() string {
	return "activity"
}

// ActivityDetail 活动细节实体
type ActivityDetail struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	ActivityID uint    `gorm:"not null;column:activity_id" json:"activity_id"`
	GoodsID    uint    `gorm:"not null;column:goods_id" json:"goods_id"`
	Discount   float64 `gorm:"type:decimal(10,2);not null" json:"discount"` // 折扣值或优惠金额
	GoodsName  string  `gorm:"-" json:"goods_name,omitempty"`                // 关联查询时填充
}

// TableName 指定表名
func (ActivityDetail) TableName() string {
	return "activity_detail"
}
