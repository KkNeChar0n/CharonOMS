package service

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
	"errors"
)

// ApprovalFlowService 审批流领域服务
type ApprovalFlowService struct {
	flowRepo     repository.ApprovalFlowManagementRepository
	nodeCaseRepo repository.ApprovalNodeCaseRepository
	templateRepo repository.ApprovalFlowTemplateRepository
}

// NewApprovalFlowService 创建审批流领域服务
func NewApprovalFlowService(
	flowRepo repository.ApprovalFlowManagementRepository,
	nodeCaseRepo repository.ApprovalNodeCaseRepository,
	templateRepo repository.ApprovalFlowTemplateRepository,
) *ApprovalFlowService {
	return &ApprovalFlowService{
		flowRepo:     flowRepo,
		nodeCaseRepo: nodeCaseRepo,
		templateRepo: templateRepo,
	}
}

// ProcessApprove 处理审批通过逻辑
func (s *ApprovalFlowService) ProcessApprove(nodeCaseUserID int) error {
	// 1. 获取审批人员记录
	nodeCaseUser, err := s.nodeCaseRepo.GetNodeUserByID(nodeCaseUserID)
	if err != nil {
		return err
	}

	// 检查是否已处理
	if nodeCaseUser.Result != nil {
		return errors.New("该审批已处理")
	}

	// 2. 更新当前用户审批结果为通过
	result := int8(0)
	if err := s.nodeCaseRepo.UpdateUserResult(nodeCaseUserID, result); err != nil {
		return err
	}

	// 3. 获取节点实例
	nodeCase, err := s.nodeCaseRepo.GetByID(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 4. 获取该节点的所有审批人员
	nodeUsers, err := s.nodeCaseRepo.GetNodeUsers(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 5. 判断节点是否通过
	nodePassed := false
	if nodeCase.Type == 0 {
		// 会签节点：所有人都通过才能通过
		nodePassed = s.isCountersignNodePassed(nodeUsers)
	} else {
		// 或签节点：任意一人通过即可通过
		nodePassed = s.isOrSignNodePassed(nodeUsers)
		if nodePassed {
			// 删除同节点其他待审批人员
			if err := s.nodeCaseRepo.DeletePendingUsers(nodeCaseUser.ApprovalNodeCaseID, nodeCaseUserID); err != nil {
				return err
			}
		}
	}

	// 6. 如果节点通过，更新节点结果并流转
	if nodePassed {
		if err := s.nodeCaseRepo.UpdateNodeResult(nodeCaseUser.ApprovalNodeCaseID, 0); err != nil {
			return err
		}

		// 7. 流转到下一节点或完成审批流
		return s.proceedToNextNodeOrComplete(nodeCase, nodeCaseUser.ApprovalNodeCaseID)
	}

	return nil
}

// ProcessReject 处理审批驳回逻辑
func (s *ApprovalFlowService) ProcessReject(nodeCaseUserID int) error {
	// 1. 获取审批人员记录
	nodeCaseUser, err := s.nodeCaseRepo.GetNodeUserByID(nodeCaseUserID)
	if err != nil {
		return err
	}

	// 检查是否已处理
	if nodeCaseUser.Result != nil {
		return errors.New("该审批已处理")
	}

	// 2. 更新当前用户审批结果为驳回
	result := int8(1)
	if err := s.nodeCaseRepo.UpdateUserResult(nodeCaseUserID, result); err != nil {
		return err
	}

	// 3. 获取节点实例
	nodeCase, err := s.nodeCaseRepo.GetByID(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 4. 获取该节点的所有审批人员
	nodeUsers, err := s.nodeCaseRepo.GetNodeUsers(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 5. 判断节点是否驳回
	nodeRejected := false
	if nodeCase.Type == 0 {
		// 会签节点：任意一人驳回即驳回
		nodeRejected = s.isCountersignNodeRejected(nodeUsers)
		if nodeRejected {
			// 删除同节点其他待审批人员
			if err := s.nodeCaseRepo.DeletePendingUsers(nodeCaseUser.ApprovalNodeCaseID, nodeCaseUserID); err != nil {
				return err
			}
		}
	} else {
		// 或签节点：所有人都驳回才驳回
		nodeRejected = s.isOrSignNodeRejected(nodeUsers)
	}

	// 6. 如果节点驳回，更新节点结果和审批流状态
	if nodeRejected {
		if err := s.nodeCaseRepo.UpdateNodeResult(nodeCaseUser.ApprovalNodeCaseID, 1); err != nil {
			return err
		}

		// 更新审批流状态为已驳回(20)
		// 需要从nodeCase获取flowID
		return s.flowRepo.UpdateStatus(nodeCase.ApprovalFlowManagementID, 20)
	}

	return nil
}

// isCountersignNodePassed 判断会签节点是否通过
func (s *ApprovalFlowService) isCountersignNodePassed(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 所有人都审批 && 所有人都通过 → 节点通过
	for _, user := range nodeUsers {
		if user.Result == nil {
			return false // 有人未审批
		}
		if *user.Result != 0 {
			return false // 有人驳回
		}
	}
	return true
}

// isCountersignNodeRejected 判断会签节点是否驳回
func (s *ApprovalFlowService) isCountersignNodeRejected(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 任意一人驳回 → 节点驳回
	for _, user := range nodeUsers {
		if user.Result != nil && *user.Result == 1 {
			return true
		}
	}
	return false
}

// isOrSignNodePassed 判断或签节点是否通过
func (s *ApprovalFlowService) isOrSignNodePassed(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 任意一人通过 → 节点通过
	for _, user := range nodeUsers {
		if user.Result != nil && *user.Result == 0 {
			return true
		}
	}
	return false
}

// isOrSignNodeRejected 判断或签节点是否驳回
func (s *ApprovalFlowService) isOrSignNodeRejected(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 所有人都审批 && 所有人都驳回 → 节点驳回
	for _, user := range nodeUsers {
		if user.Result == nil {
			return false // 有人未审批
		}
		if *user.Result != 1 {
			return false // 有人通过
		}
	}
	return true
}

// proceedToNextNodeOrComplete 流转到下一节点或完成审批流
func (s *ApprovalFlowService) proceedToNextNodeOrComplete(nodeCase *entity.ApprovalNodeCase, nodeCaseID int) error {
	// 1. 获取审批流信息
	flow, err := s.flowRepo.GetByID(nodeCase.ApprovalFlowManagementID)
	if err != nil {
		return err
	}

	// 2. 获取模板节点信息
	templateNode, err := s.nodeCaseRepo.GetTemplateNodeByID(nodeCase.NodeID)
	if err != nil {
		return err
	}

	// 3. 查找下一个节点
	nextNode, err := s.nodeCaseRepo.GetNextTemplateNode(flow.ApprovalFlowTemplateID, templateNode.Sort)
	if err != nil && err.Error() != "record not found" {
		return err
	}

	if nextNode != nil {
		// 有下一节点：创建下一节点实例
		approvers, err := s.nodeCaseRepo.GetNodeApprovers(nextNode.ID)
		if err != nil {
			return err
		}

		if err := s.nodeCaseRepo.CreateNextNode(
			flow.ID,
			nextNode.ID,
			nextNode.Type,
			nextNode.Sort,
			approvers,
		); err != nil {
			return err
		}

		// 增加step
		return s.flowRepo.IncrementStep(flow.ID)
	} else {
		// 没有下一节点：审批流完成
		if err := s.flowRepo.UpdateStatus(flow.ID, 10); err != nil {
			return err
		}

		// 创建抄送记录
		copyUsers, err := s.templateRepo.GetCopyUsers(flow.ApprovalFlowTemplateID)
		if err != nil {
			return err
		}

		if len(copyUsers) > 0 {
			copyInfo := "审批流已完成"
			return s.nodeCaseRepo.CreateCopyRecords(flow.ID, copyUsers, copyInfo)
		}

		return nil
	}
}
