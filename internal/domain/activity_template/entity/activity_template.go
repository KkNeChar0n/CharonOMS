package entity

import "time"

// ActivityTemplate 活动模板实体
type ActivityTemplate struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Type       int       `gorm:"column:type;not null" json:"type"`             // 1=满减, 2=满折, 3=满赠
	SelectType int       `gorm:"column:select_type;not null" json:"select_type"` // 1=按分类, 2=按商品
	Status     int       `gorm:"column:status;default:0" json:"status"`         // 0=启用, 1=禁用
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (ActivityTemplate) TableName() string {
	return "activity_template"
}

// IsEnabled 判断模板是否启用
func (t *ActivityTemplate) IsEnabled() bool {
	return t.Status == 0
}

// IsDisabled 判断模板是否禁用
func (t *ActivityTemplate) IsDisabled() bool {
	return t.Status == 1
}

// Validate 验证模板数据
func (t *ActivityTemplate) Validate() error {
	if t.Name == "" {
		return ErrTemplateNameRequired
	}
	if t.Type < 1 || t.Type > 3 {
		return ErrInvalidTemplateType
	}
	if t.SelectType < 1 || t.SelectType > 2 {
		return ErrInvalidSelectType
	}
	return nil
}
