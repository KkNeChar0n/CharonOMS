package entity

import "time"

// Classify 分类实体
type Classify struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Level     int       `gorm:"column:level;not null" json:"level"`     // 0=一级分类，1=二级分类
	ParentID  *int      `gorm:"column:parentid" json:"parentid"`        // 父级分类ID，一级分类为NULL
	Status    int       `gorm:"column:status;default:0" json:"status"`  // 0=启用，1=禁用
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Classify) TableName() string {
	return "classify"
}
