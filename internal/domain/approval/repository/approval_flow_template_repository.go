package repository

import "charonoms/internal/domain/approval/entity"

// ApprovalFlowTemplateRepository 审批流模板仓储接口
type ApprovalFlowTemplateRepository interface {
	// GetList 获取审批流模板列表（含关联类型名称）
	GetList(filters map[string]interface{}) ([]map[string]interface{}, error)

	// GetByID 获取模板基本信息
	GetByID(id int) (*entity.ApprovalFlowTemplate, error)

	// GetDetailByID 获取模板完整详情（含节点、审批人员、抄送人员）
	GetDetailByID(id int) (map[string]interface{}, error)

	// Create 创建审批流模板（事务：模板、节点、人员、抄送）
	Create(template *entity.ApprovalFlowTemplate, nodes []entity.ApprovalFlowTemplateNode,
		nodeApprovers map[int][]int, copyUsers []int) error

	// UpdateStatus 更新模板状态
	UpdateStatus(id int, status int8) error

	// DisableSameTypeTemplates 禁用同类型的其他模板
	DisableSameTypeTemplates(typeID int, excludeID int) error

	// GetNodesByTemplateID 获取模板的所有节点（按sort排序）
	GetNodesByTemplateID(templateID int) ([]entity.ApprovalFlowTemplateNode, error)

	// GetNodeApprovers 获取节点的审批人员
	GetNodeApprovers(nodeID int) ([]int, error)

	// GetCopyUsers 获取模板的抄送人员
	GetCopyUsers(templateID int) ([]int, error)
}
