package contract

import (
	"fmt"

	"charonoms/internal/domain/contract/repository"
)

// ContractService 合同业务服务
type ContractService struct {
	contractRepo repository.ContractRepository
}

// NewContractService 创建合同业务服务实例
func NewContractService(contractRepo repository.ContractRepository) *ContractService {
	return &ContractService{
		contractRepo: contractRepo,
	}
}

// GetContractList 获取合同列表
func (s *ContractService) GetContractList() ([]map[string]interface{}, error) {
	return s.contractRepo.GetContractList()
}

// GetContractByID 获取合同详情
func (s *ContractService) GetContractByID(id int) (map[string]interface{}, error) {
	result, err := s.contractRepo.GetContractByID(id)
	if err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		return nil, fmt.Errorf("contract not found")
	}

	return result, nil
}

// CreateContract 创建合同
func (s *ContractService) CreateContract(req *CreateContractRequest, initiator string) (int, error) {
	// 验证必填字段
	if req.Name == "" {
		return 0, fmt.Errorf("请填写所有必填项")
	}
	if req.StudentID == 0 {
		return 0, fmt.Errorf("请填写所有必填项")
	}
	if req.Type == nil {
		return 0, fmt.Errorf("请填写所有必填项")
	}
	if req.SignatureForm == nil {
		return 0, fmt.Errorf("请填写所有必填项")
	}
	if req.ContractAmount == 0 {
		return 0, fmt.Errorf("请填写所有必填项")
	}

	// 验证合同类型
	if *req.Type != 0 && *req.Type != 1 {
		return 0, fmt.Errorf("合同类型必须为0（首报）或1（续报）")
	}

	// 验证签署形式
	if *req.SignatureForm != 0 && *req.SignatureForm != 1 {
		return 0, fmt.Errorf("签署形式必须为0（线上签署）或1（线下签署）")
	}

	contractID, err := s.contractRepo.CreateContract(
		req.Name,
		req.StudentID,
		*req.Type,
		*req.SignatureForm,
		req.ContractAmount,
		req.Signatory,
		initiator,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create contract: %w", err)
	}

	return contractID, nil
}

// RevokeContract 撤销合同
func (s *ContractService) RevokeContract(id int) error {
	// 获取合同当前状态
	status, exists, err := s.contractRepo.GetContractStatus(id)
	if err != nil {
		return fmt.Errorf("failed to get contract status: %w", err)
	}
	if !exists {
		return fmt.Errorf("contract not found")
	}

	// 验证状态必须为0（待审核）
	if status != 0 {
		return fmt.Errorf("只有待审核状态的合同可以撤销")
	}

	return s.contractRepo.RevokeContract(id)
}

// TerminateContract 中止合作
func (s *ContractService) TerminateContract(id int, req *TerminateContractRequest) error {
	// 验证必填字段
	if req.TerminationAgreement == "" {
		return fmt.Errorf("请上传中止协议文件")
	}

	// 获取合同当前状态
	status, exists, err := s.contractRepo.GetContractStatus(id)
	if err != nil {
		return fmt.Errorf("failed to get contract status: %w", err)
	}
	if !exists {
		return fmt.Errorf("contract not found")
	}

	// 验证状态必须为50（已通过）
	if status != 50 {
		return fmt.Errorf("只有已通过状态的合同可以中止")
	}

	return s.contractRepo.TerminateContract(id, req.TerminationAgreement)
}
