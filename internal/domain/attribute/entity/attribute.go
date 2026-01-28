package entity

import "time"

// Attribute 属性实体
type Attribute struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Classify   int       `gorm:"column:classify;not null;default:0" json:"classify"` // 0=属性，1=规格
	Status     int       `gorm:"column:status;not null;default:0" json:"status"`      // 0=启用，1=禁用
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Attribute) TableName() string {
	return "attribute"
}
