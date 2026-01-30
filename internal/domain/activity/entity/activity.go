package entity

import "time"

// Activity 活动实体
type Activity struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	TemplateID int       `gorm:"column:template_id;not null" json:"template_id"`
	StartTime  time.Time `gorm:"column:start_time;not null" json:"start_time"`
	EndTime    time.Time `gorm:"column:end_time;not null" json:"end_time"`
	Status     int       `gorm:"column:status;default:0" json:"status"` // 0=启用, 1=禁用
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (Activity) TableName() string {
	return "activity"
}

// IsEnabled 判断活动是否启用
func (a *Activity) IsEnabled() bool {
	return a.Status == 0
}

// IsDisabled 判断活动是否禁用
func (a *Activity) IsDisabled() bool {
	return a.Status == 1
}

// Validate 验证活动数据
func (a *Activity) Validate() error {
	if a.Name == "" {
		return ErrActivityNameRequired
	}
	if a.TemplateID == 0 {
		return ErrTemplateIDRequired
	}
	if a.StartTime.IsZero() || a.EndTime.IsZero() {
		return ErrTimeRequired
	}
	if !a.StartTime.Before(a.EndTime) {
		return ErrInvalidTimeRange
	}
	return nil
}
