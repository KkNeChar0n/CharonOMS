package entity

import "time"

// ApprovalFlowTemplate 审批流模板实体
type ApprovalFlowTemplate struct {
	ID                  int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	ApprovalFlowTypeID  int       `gorm:"column:approval_flow_type_id;not null" json:"approval_flow_type_id"`
	Creator             string    `gorm:"column:creator;type:varchar(100)" json:"creator"`
	Status              int8      `gorm:"column:status;not null;default:0" json:"status"` // 0=启用，1=禁用
	CreateTime          time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime          time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 指定表名
func (ApprovalFlowTemplate) TableName() string {
	return "approval_flow_template"
}

// ApprovalFlowTemplateNode 审批流模板节点实体
type ApprovalFlowTemplateNode struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateID int       `gorm:"column:template_id;not null" json:"template_id"`
	Name       string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Sort       int       `gorm:"column:sort;not null" json:"sort"`
	Type       int8      `gorm:"column:type;not null" json:"type"` // 0=会签，1=或签
	CreateTime time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
}

// TableName 指定表名
func (ApprovalFlowTemplateNode) TableName() string {
	return "approval_flow_template_node"
}

// ApprovalNodeUserAccount 审批节点人员配置实体
type ApprovalNodeUserAccount struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NodeID        int       `gorm:"column:node_id;not null" json:"node_id"`
	UserAccountID int       `gorm:"column:useraccount_id;not null" json:"useraccount_id"`
	CreateTime    time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
}

// TableName 指定表名
func (ApprovalNodeUserAccount) TableName() string {
	return "approval_node_useraccount"
}

// ApprovalCopyUserAccount 抄送人员配置实体
type ApprovalCopyUserAccount struct {
	ID                     int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApprovalFlowTemplateID int       `gorm:"column:approval_flow_template_id;not null" json:"approval_flow_template_id"`
	UserAccountID          int       `gorm:"column:useraccount_id;not null" json:"useraccount_id"`
	CreateTime             time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
}

// TableName 指定表名
func (ApprovalCopyUserAccount) TableName() string {
	return "approval_copy_useraccount"
}
