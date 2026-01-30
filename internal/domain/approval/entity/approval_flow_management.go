package entity

import "time"

// ApprovalFlowManagement 审批流实例实体
type ApprovalFlowManagement struct {
	ID                     int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApprovalFlowTemplateID int        `gorm:"column:approval_flow_template_id;not null" json:"approval_flow_template_id"`
	ApprovalFlowTypeID     int        `gorm:"column:approval_flow_type_id;not null" json:"approval_flow_type_id"`
	Step                   int        `gorm:"column:step;not null;default:0" json:"step"`         // 当前执行到第几步节点
	CreateUser             int        `gorm:"column:create_user;not null" json:"create_user"`     // 发起人ID
	CreateTime             time.Time  `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	Status                 int8       `gorm:"column:status;not null;default:0" json:"status"` // 0=待审批，10=已通过，20=已驳回，99=已撤销
	CompleteTime           *time.Time `gorm:"column:complete_time" json:"complete_time"`      // 完成时间
}

// TableName 指定表名
func (ApprovalFlowManagement) TableName() string {
	return "approval_flow_management"
}

// ApprovalNodeCase 审批节点实例实体
type ApprovalNodeCase struct {
	ID                       int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NodeID                   int        `gorm:"column:node_id;not null" json:"node_id"` // 关联的模板节点ID
	ApprovalFlowManagementID int        `gorm:"column:approval_flow_management_id;not null" json:"approval_flow_management_id"`
	Type                     int8       `gorm:"column:type;not null" json:"type"` // 0=会签，1=或签
	Sort                     int        `gorm:"column:sort;not null" json:"sort"`
	Result                   *int8      `gorm:"column:result" json:"result"`               // NULL=审批中，0=通过，1=驳回
	CreateTime               time.Time  `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	CompleteTime             *time.Time `gorm:"column:complete_time" json:"complete_time"` // 完成时间
}

// TableName 指定表名
func (ApprovalNodeCase) TableName() string {
	return "approval_node_case"
}

// ApprovalNodeCaseUser 审批人员记录实体
type ApprovalNodeCaseUser struct {
	ID                  int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApprovalNodeCaseID  int        `gorm:"column:approval_node_case_id;not null" json:"approval_node_case_id"`
	UserAccountID       int        `gorm:"column:useraccount_id;not null" json:"useraccount_id"`
	Result              *int8      `gorm:"column:result" json:"result"`               // NULL=待审批，0=通过，1=驳回
	CreateTime          time.Time  `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	HandleTime          *time.Time `gorm:"column:handle_time" json:"handle_time"` // 处理时间
}

// TableName 指定表名
func (ApprovalNodeCaseUser) TableName() string {
	return "approval_node_case_user"
}

// ApprovalCopyUserAccountCase 抄送记录实体
type ApprovalCopyUserAccountCase struct {
	ID                       int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApprovalFlowManagementID int       `gorm:"column:approval_flow_management_id;not null" json:"approval_flow_management_id"`
	UserAccountID            int       `gorm:"column:useraccount_id;not null" json:"useraccount_id"`
	CopyInfo                 string    `gorm:"column:copy_info;type:varchar(500);not null" json:"copy_info"`
	CreateTime               time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
}

// TableName 指定表名
func (ApprovalCopyUserAccountCase) TableName() string {
	return "approval_copy_useraccount_case"
}
