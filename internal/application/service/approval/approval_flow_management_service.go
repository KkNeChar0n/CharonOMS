package approval

import (
	"charonoms/internal/domain/approval/repository"
	"charonoms/internal/domain/approval/service"
	"errors"
)

// ApprovalFlowManagementService 审批流管理应用服务
type ApprovalFlowManagementService struct {
	flowRepo          repository.ApprovalFlowManagementRepository
	nodeCaseRepo      repository.ApprovalNodeCaseRepository
	flowDomainService *service.ApprovalFlowService
}

// NewApprovalFlowManagementService 创建审批流管理应用服务实例
func NewApprovalFlowManagementService(
	flowRepo repository.ApprovalFlowManagementRepository,
	nodeCaseRepo repository.ApprovalNodeCaseRepository,
	flowDomainService *service.ApprovalFlowService,
) *ApprovalFlowManagementService {
	return &ApprovalFlowManagementService{
		flowRepo:          flowRepo,
		nodeCaseRepo:      nodeCaseRepo,
		flowDomainService: flowDomainService,
	}
}

// GetInitiatedFlows 获取用户发起的审批流
func (s *ApprovalFlowManagementService) GetInitiatedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	return s.flowRepo.GetInitiatedFlows(userID, filters)
}

// GetPendingFlows 获取待用户审批的任务
func (s *ApprovalFlowManagementService) GetPendingFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	return s.flowRepo.GetPendingFlows(userID, filters)
}

// GetCompletedFlows 获取用户已处理的审批任务
func (s *ApprovalFlowManagementService) GetCompletedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	return s.flowRepo.GetCompletedFlows(userID, filters)
}

// GetCopiedFlows 获取抄送给用户的通知
func (s *ApprovalFlowManagementService) GetCopiedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	return s.flowRepo.GetCopiedFlows(userID, filters)
}

// GetDetail 获取审批流详情
func (s *ApprovalFlowManagementService) GetDetail(flowID int, userID int) (map[string]interface{}, error) {
	return s.flowRepo.GetDetailByID(flowID, userID)
}

// CreateFromTemplate 从模板创建审批流实例
func (s *ApprovalFlowManagementService) CreateFromTemplate(templateID int, userID int) (int, error) {
	// 调用仓储层从模板创建审批流实例（事务）
	return s.flowRepo.CreateFromTemplate(templateID, userID)
}

// Cancel 撤销审批流
func (s *ApprovalFlowManagementService) Cancel(flowID int, userID int) error {
	// 获取审批流信息
	flow, err := s.flowRepo.GetByID(flowID)
	if err != nil {
		return err
	}

	// 验证是否为发起人
	if flow.CreateUser != userID {
		return errors.New("只有发起人可以撤销审批流")
	}

	// 验证状态是否为待审批
	if flow.Status != 0 {
		return errors.New("只能撤销待审批的流程")
	}

	// 调用仓储层撤销
	return s.flowRepo.Cancel(flowID, userID)
}

// Approve 审批通过
func (s *ApprovalFlowManagementService) Approve(nodeCaseUserID int, userID int) error {
	// 获取审批人员记录
	nodeCaseUser, err := s.nodeCaseRepo.GetNodeUserByID(nodeCaseUserID)
	if err != nil {
		return err
	}

	// 验证审批人员是否为当前用户
	if nodeCaseUser.UserAccountID != userID {
		return errors.New("无权处理此审批")
	}

	// 检查是否已处理
	if nodeCaseUser.Result != nil {
		return errors.New("该审批已处理")
	}

	// 调用领域服务处理审批通过逻辑
	return s.flowDomainService.ProcessApprove(nodeCaseUserID)
}

// Reject 审批驳回
func (s *ApprovalFlowManagementService) Reject(nodeCaseUserID int, userID int) error {
	// 获取审批人员记录
	nodeCaseUser, err := s.nodeCaseRepo.GetNodeUserByID(nodeCaseUserID)
	if err != nil {
		return err
	}

	// 验证审批人员是否为当前用户
	if nodeCaseUser.UserAccountID != userID {
		return errors.New("无权处理此审批")
	}

	// 检查是否已处理
	if nodeCaseUser.Result != nil {
		return errors.New("该审批已处理")
	}

	// 调用领域服务处理审批驳回逻辑
	return s.flowDomainService.ProcessReject(nodeCaseUserID)
}
