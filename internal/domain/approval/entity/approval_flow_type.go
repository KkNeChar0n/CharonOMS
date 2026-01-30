package entity

import "time"

// ApprovalFlowType 审批流类型实体
type ApprovalFlowType struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Status     int8      `gorm:"column:status;not null;default:0" json:"status"` // 0=启用，1=禁用
	CreateTime time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 指定表名
func (ApprovalFlowType) TableName() string {
	return "approval_flow_type"
}
