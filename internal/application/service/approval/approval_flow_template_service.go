package approval

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
	"errors"
)

// ApprovalFlowTemplateService 审批流模板应用服务
type ApprovalFlowTemplateService struct {
	templateRepo repository.ApprovalFlowTemplateRepository
	typeRepo     repository.ApprovalFlowTypeRepository
}

// NewApprovalFlowTemplateService 创建审批流模板应用服务
func NewApprovalFlowTemplateService(
	templateRepo repository.ApprovalFlowTemplateRepository,
	typeRepo repository.ApprovalFlowTypeRepository,
) *ApprovalFlowTemplateService {
	return &ApprovalFlowTemplateService{
		templateRepo: templateRepo,
		typeRepo:     typeRepo,
	}
}

// GetList 获取审批流模板列表
func (s *ApprovalFlowTemplateService) GetList(filters map[string]interface{}) ([]map[string]interface{}, error) {
	return s.templateRepo.GetList(filters)
}

// GetDetail 获取审批流模板详情
func (s *ApprovalFlowTemplateService) GetDetail(id int) (map[string]interface{}, error) {
	return s.templateRepo.GetDetailByID(id)
}

// GetByID 获取模板基本信息
func (s *ApprovalFlowTemplateService) GetByID(id int) (*entity.ApprovalFlowTemplate, error) {
	return s.templateRepo.GetByID(id)
}

// DisableSameTypeTemplates 禁用同类型的其他模板
func (s *ApprovalFlowTemplateService) DisableSameTypeTemplates(typeID int, excludeID int) error {
	return s.templateRepo.DisableSameTypeTemplates(typeID, excludeID)
}

// Create 创建审批流模板
func (s *ApprovalFlowTemplateService) Create(
	template *entity.ApprovalFlowTemplate,
	nodes []entity.ApprovalFlowTemplateNode,
	nodeApprovers map[int][]int,
	copyUsers []int,
) error {
	// 验证：至少一个节点
	if len(nodes) == 0 {
		return errors.New("模板至少需要一个节点")
	}

	// 验证：每个节点至少一个审批人
	for i := range nodes {
		if len(nodeApprovers[i]) == 0 {
			return errors.New("每个节点至少需要一个审批人")
		}
	}

	// 创建模板
	return s.templateRepo.Create(template, nodes, nodeApprovers, copyUsers)
}

// UpdateStatus 更新模板状态
func (s *ApprovalFlowTemplateService) UpdateStatus(id int, status int8) error {
	// 获取模板信息
	template, err := s.templateRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 如果启用，先禁用同类型的其他模板
	if status == 0 {
		if err := s.templateRepo.DisableSameTypeTemplates(template.ApprovalFlowTypeID, id); err != nil {
			return err
		}
	}

	// 更新状态
	return s.templateRepo.UpdateStatus(id, status)
}
