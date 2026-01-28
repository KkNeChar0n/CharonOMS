package entity

import "time"

// Brand 品牌实体
type Brand struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Status     int       `gorm:"column:status;default:0" json:"status"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Brand) TableName() string {
	return "brand"
}
