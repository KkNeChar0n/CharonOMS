package approval

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormApprovalNodeCaseRepository GORM实现的审批节点实例仓储
type GormApprovalNodeCaseRepository struct {
	db *gorm.DB
}

// NewApprovalNodeCaseRepository 创建审批节点实例仓储实例
func NewApprovalNodeCaseRepository(db *gorm.DB) repository.ApprovalNodeCaseRepository {
	return &GormApprovalNodeCaseRepository{db: db}
}

// GetByID 根据ID获取节点实例
func (r *GormApprovalNodeCaseRepository) GetByID(id int) (*entity.ApprovalNodeCase, error) {
	var nodeCase entity.ApprovalNodeCase
	err := r.db.First(&nodeCase, id).Error
	if err != nil {
		return nil, err
	}
	return &nodeCase, nil
}

// GetByFlowIDAndStep 获取当前节点实例
func (r *GormApprovalNodeCaseRepository) GetByFlowIDAndStep(flowID int, step int) (*entity.ApprovalNodeCase, error) {
	var nodeCase entity.ApprovalNodeCase
	err := r.db.Where("approval_flow_management_id = ? AND sort = ?", flowID, step).
		First(&nodeCase).Error
	if err != nil {
		return nil, err
	}
	return &nodeCase, nil
}

// GetNodeUsers 获取节点的所有审批人员记录
func (r *GormApprovalNodeCaseRepository) GetNodeUsers(nodeCaseID int) ([]entity.ApprovalNodeCaseUser, error) {
	var users []entity.ApprovalNodeCaseUser
	err := r.db.Where("approval_node_case_id = ?", nodeCaseID).Find(&users).Error
	return users, err
}

// GetNodeUserByID 根据ID获取审批人员记录
func (r *GormApprovalNodeCaseRepository) GetNodeUserByID(id int) (*entity.ApprovalNodeCaseUser, error) {
	var user entity.ApprovalNodeCaseUser
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateNextNode 创建下一节点实例（事务：节点实例、审批人员记录）
func (r *GormApprovalNodeCaseRepository) CreateNextNode(flowID int, templateNodeID int, nodeType int8, sort int, approvers []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建节点实例
		nodeCase := entity.ApprovalNodeCase{
			NodeID:                   templateNodeID,
			ApprovalFlowManagementID: flowID,
			Type:                     nodeType,
			Sort:                     sort,
		}
		if err := tx.Create(&nodeCase).Error; err != nil {
			return fmt.Errorf("创建节点实例失败: %w", err)
		}

		// 2. 创建审批人员记录
		if len(approvers) == 0 {
			return fmt.Errorf("审批人员列表为空")
		}

		for _, approverID := range approvers {
			userCase := entity.ApprovalNodeCaseUser{
				ApprovalNodeCaseID: nodeCase.ID,
				UserAccountID:      approverID,
			}
			if err := tx.Create(&userCase).Error; err != nil {
				return fmt.Errorf("创建审批人员记录失败: %w", err)
			}
		}

		return nil
	})
}

// UpdateNodeResult 更新节点结果
func (r *GormApprovalNodeCaseRepository) UpdateNodeResult(nodeCaseID int, result int8) error {
	now := time.Now()
	return r.db.Model(&entity.ApprovalNodeCase{}).
		Where("id = ?", nodeCaseID).
		Updates(map[string]interface{}{
			"result":        &result,
			"complete_time": &now,
		}).Error
}

// UpdateUserResult 更新人员审批结果（含锁定检查）
func (r *GormApprovalNodeCaseRepository) UpdateUserResult(userID int, result int8) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 使用FOR UPDATE锁定记录
		var user entity.ApprovalNodeCaseUser
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", userID).
			First(&user).Error
		if err != nil {
			return fmt.Errorf("获取审批人员记录失败: %w", err)
		}

		// 2. 检查是否已经处理过
		if user.Result != nil {
			return fmt.Errorf("该审批任务已经处理过")
		}

		// 3. 更新审批结果
		now := time.Now()
		err = tx.Model(&entity.ApprovalNodeCaseUser{}).
			Where("id = ?", userID).
			Updates(map[string]interface{}{
				"result":      &result,
				"handle_time": &now,
			}).Error
		if err != nil {
			return fmt.Errorf("更新审批结果失败: %w", err)
		}

		return nil
	})
}

// DeletePendingUsers 删除同节点其他待审批人员
func (r *GormApprovalNodeCaseRepository) DeletePendingUsers(nodeCaseID int, excludeUserID int) error {
	return r.db.Where("approval_node_case_id = ?", nodeCaseID).
		Where("id != ?", excludeUserID).
		Where("result IS NULL").
		Delete(&entity.ApprovalNodeCaseUser{}).Error
}

// CreateCopyRecords 创建抄送记录
func (r *GormApprovalNodeCaseRepository) CreateCopyRecords(flowID int, copyUsers []int, copyInfo string) error {
	if len(copyUsers) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, userID := range copyUsers {
			copyCase := entity.ApprovalCopyUserAccountCase{
				ApprovalFlowManagementID: flowID,
				UserAccountID:            userID,
				CopyInfo:                 copyInfo,
			}
			if err := tx.Create(&copyCase).Error; err != nil {
				return fmt.Errorf("创建抄送记录失败: %w", err)
			}
		}
		return nil
	})
}

// GetTemplateNodeByID 获取模板节点信息
func (r *GormApprovalNodeCaseRepository) GetTemplateNodeByID(nodeID int) (*entity.ApprovalFlowTemplateNode, error) {
	var node entity.ApprovalFlowTemplateNode
	err := r.db.First(&node, nodeID).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// GetNextTemplateNode 获取下一个模板节点
func (r *GormApprovalNodeCaseRepository) GetNextTemplateNode(templateID int, currentSort int) (*entity.ApprovalFlowTemplateNode, error) {
	var node entity.ApprovalFlowTemplateNode
	err := r.db.Where("template_id = ? AND sort > ?", templateID, currentSort).
		Order("sort ASC").
		First(&node).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 没有下一个节点，返回nil而不是错误
			return nil, nil
		}
		return nil, err
	}
	return &node, nil
}

// GetNodeApprovers 获取模板节点的审批人员
func (r *GormApprovalNodeCaseRepository) GetNodeApprovers(nodeID int) ([]int, error) {
	var approvers []entity.ApprovalNodeUserAccount
	err := r.db.Where("node_id = ?", nodeID).Find(&approvers).Error
	if err != nil {
		return nil, err
	}

	if len(approvers) == 0 {
		return nil, fmt.Errorf("节点未配置审批人员")
	}

	approverIDs := make([]int, len(approvers))
	for i, a := range approvers {
		approverIDs[i] = a.UserAccountID
	}
	return approverIDs, nil
}
