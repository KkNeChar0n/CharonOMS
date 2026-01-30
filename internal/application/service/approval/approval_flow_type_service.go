package approval

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
)

// ApprovalFlowTypeService 审批流类型应用服务
type ApprovalFlowTypeService struct {
	typeRepo repository.ApprovalFlowTypeRepository
}

// NewApprovalFlowTypeService 创建审批流类型应用服务实例
func NewApprovalFlowTypeService(typeRepo repository.ApprovalFlowTypeRepository) *ApprovalFlowTypeService {
	return &ApprovalFlowTypeService{
		typeRepo: typeRepo,
	}
}

// GetList 获取审批流类型列表
func (s *ApprovalFlowTypeService) GetList(filters map[string]interface{}) ([]entity.ApprovalFlowType, error) {
	return s.typeRepo.GetList(filters)
}

// Create 创建审批流类型
func (s *ApprovalFlowTypeService) Create(flowType *entity.ApprovalFlowType) error {
	return s.typeRepo.Create(flowType)
}

// UpdateStatus 更新审批流类型状态
func (s *ApprovalFlowTypeService) UpdateStatus(id int, status int8) error {
	return s.typeRepo.UpdateStatus(id, status)
}
