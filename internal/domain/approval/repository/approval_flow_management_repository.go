package repository

import "charonoms/internal/domain/approval/entity"

// ApprovalFlowManagementRepository 审批流实例仓储接口
type ApprovalFlowManagementRepository interface {
	// GetInitiatedFlows 获取用户发起的审批流
	GetInitiatedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error)

	// GetPendingFlows 获取待用户审批的任务
	GetPendingFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error)

	// GetCompletedFlows 获取用户已处理的审批任务
	GetCompletedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error)

	// GetCopiedFlows 获取抄送给用户的通知
	GetCopiedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error)

	// GetDetailByID 获取审批流详情
	GetDetailByID(flowID int, userID int) (map[string]interface{}, error)

	// GetByID 根据ID查询审批流
	GetByID(id int) (*entity.ApprovalFlowManagement, error)

	// CreateFromTemplate 从模板创建审批流实例（事务：审批流、第一个节点、审批人员）
	CreateFromTemplate(templateID int, userID int) (int, error)

	// UpdateStatus 更新审批流状态
	UpdateStatus(flowID int, status int8) error

	// Cancel 撤销审批流
	Cancel(flowID int, userID int) error

	// IncrementStep 增加step步骤
	IncrementStep(flowID int) error
}
