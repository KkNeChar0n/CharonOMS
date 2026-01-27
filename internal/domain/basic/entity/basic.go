package entity

import "time"

// Sex 性别实体
type Sex struct {
	ID   uint   `gorm:"column:id;primaryKey" json:"id"`
	Name string `gorm:"column:name;size:10;not null" json:"name"`
}

// TableName 指定表名
func (Sex) TableName() string {
	return "sex"
}

// Grade 年级实体
type Grade struct {
	ID         uint      `gorm:"column:id;primaryKey" json:"id"`
	Name       string    `gorm:"column:name;size:50;not null" json:"grade"`
	Status     int       `gorm:"column:status;default:0" json:"status"` // 0-启用 1-禁用
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Grade) TableName() string {
	return "grade"
}

// Subject 学科实体
type Subject struct {
	ID         uint      `gorm:"column:id;primaryKey" json:"id"`
	Subject    string    `gorm:"column:subject;size:50;not null" json:"subject"`
	Status     int       `gorm:"column:status;default:0" json:"status"` // 0-启用 1-禁用
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Subject) TableName() string {
	return "subject"
}
