package repository

import "charonoms/internal/domain/approval/entity"

// ApprovalNodeCaseRepository 审批节点实例仓储接口
type ApprovalNodeCaseRepository interface {
	// GetByID 根据ID获取节点实例
	GetByID(id int) (*entity.ApprovalNodeCase, error)

	// GetByFlowIDAndStep 获取当前节点实例
	GetByFlowIDAndStep(flowID int, step int) (*entity.ApprovalNodeCase, error)

	// GetNodeUsers 获取节点的所有审批人员记录
	GetNodeUsers(nodeCaseID int) ([]entity.ApprovalNodeCaseUser, error)

	// GetNodeUserByID 根据ID获取审批人员记录
	GetNodeUserByID(id int) (*entity.ApprovalNodeCaseUser, error)

	// CreateNextNode 创建下一节点实例（事务：节点实例、审批人员记录）
	CreateNextNode(flowID int, templateNodeID int, nodeType int8, sort int, approvers []int) error

	// UpdateNodeResult 更新节点结果
	UpdateNodeResult(nodeCaseID int, result int8) error

	// UpdateUserResult 更新人员审批结果（含锁定检查）
	UpdateUserResult(userID int, result int8) error

	// DeletePendingUsers 删除同节点其他待审批人员
	DeletePendingUsers(nodeCaseID int, excludeUserID int) error

	// CreateCopyRecords 创建抄送记录
	CreateCopyRecords(flowID int, copyUsers []int, copyInfo string) error

	// GetTemplateNodeByID 获取模板节点信息
	GetTemplateNodeByID(nodeID int) (*entity.ApprovalFlowTemplateNode, error)

	// GetNextTemplateNode 获取下一个模板节点
	GetNextTemplateNode(templateID int, currentSort int) (*entity.ApprovalFlowTemplateNode, error)

	// GetNodeApprovers 获取模板节点的审批人员
	GetNodeApprovers(nodeID int) ([]int, error)
}
