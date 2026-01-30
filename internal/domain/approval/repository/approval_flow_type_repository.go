package repository

import "charonoms/internal/domain/approval/entity"

// ApprovalFlowTypeRepository 审批流类型仓储接口
type ApprovalFlowTypeRepository interface {
	// GetList 获取审批流类型列表
	GetList(filters map[string]interface{}) ([]entity.ApprovalFlowType, error)

	// GetByID 根据ID查询审批流类型
	GetByID(id int) (*entity.ApprovalFlowType, error)

	// Create 创建审批流类型
	Create(flowType *entity.ApprovalFlowType) error

	// UpdateStatus 更新审批流类型状态
	UpdateStatus(id int, status int8) error
}
