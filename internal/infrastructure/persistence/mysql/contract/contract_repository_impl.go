package contract

import (
	"fmt"

	"charonoms/internal/domain/contract/entity"
	"charonoms/internal/domain/contract/repository"

	"gorm.io/gorm"
)

// ContractRepositoryImpl 合同仓储实现
type ContractRepositoryImpl struct {
	db *gorm.DB
}

// NewContractRepository 创建合同仓储实例
func NewContractRepository(db *gorm.DB) repository.ContractRepository {
	return &ContractRepositoryImpl{db: db}
}

// GetContractList 获取合同列表（含学生信息）
func (r *ContractRepositoryImpl) GetContractList() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			c.id,
			c.name,
			c.student_id,
			s.name as student_name,
			c.type,
			c.signature_form,
			c.contract_amount,
			c.signatory,
			c.initiating_party,
			c.initiator,
			c.status,
			c.payment_status,
			c.termination_agreement,
			c.create_time
		FROM contract c
		LEFT JOIN student s ON c.student_id = s.id
		ORDER BY c.create_time DESC
	`

	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get contract list: %w", err)
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetContractByID 根据ID获取合同详情（含学生信息）
func (r *ContractRepositoryImpl) GetContractByID(id int) (map[string]interface{}, error) {
	var result map[string]interface{}

	query := `
		SELECT
			c.id,
			c.name,
			c.student_id,
			s.name as student_name,
			c.type,
			c.signature_form,
			c.contract_amount,
			c.signatory,
			c.initiating_party,
			c.initiator,
			c.status,
			c.payment_status,
			c.termination_agreement,
			c.create_time
		FROM contract c
		LEFT JOIN student s ON c.student_id = s.id
		WHERE c.id = ?
	`

	if err := r.db.Raw(query, id).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get contract by id: %w", err)
	}

	return result, nil
}

// CreateContract 创建合同
func (r *ContractRepositoryImpl) CreateContract(name string, studentID int, contractType int, signatureForm int, contractAmount float64, signatory string, initiator string) (int, error) {
	contract := &entity.Contract{
		Name:            name,
		StudentID:       studentID,
		Type:            contractType,
		SignatureForm:   signatureForm,
		ContractAmount:  contractAmount,
		Signatory:       signatory,
		Initiator:       initiator,
		Status:          0,
		PaymentStatus:   0,
	}

	if err := r.db.Create(contract).Error; err != nil {
		return 0, fmt.Errorf("failed to create contract: %w", err)
	}

	return contract.ID, nil
}

// RevokeContract 撤销合同（status=0 → status=98）
func (r *ContractRepositoryImpl) RevokeContract(id int) error {
	result := r.db.Model(&entity.Contract{}).Where("id = ?", id).Update("status", 98)
	if result.Error != nil {
		return fmt.Errorf("failed to revoke contract: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}

	return nil
}

// TerminateContract 中止合作（status=50 → status=99）
func (r *ContractRepositoryImpl) TerminateContract(id int, terminationAgreement string) error {
	updates := map[string]interface{}{
		"status":                 99,
		"termination_agreement": terminationAgreement,
	}

	result := r.db.Model(&entity.Contract{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to terminate contract: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}

	return nil
}

// GetContractStatus 获取合同当前状态
func (r *ContractRepositoryImpl) GetContractStatus(id int) (int, bool, error) {
	var status int
	result := r.db.Model(&entity.Contract{}).Where("id = ?", id).Pluck("status", &status)
	if result.Error != nil {
		return 0, false, fmt.Errorf("failed to get contract status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return 0, false, nil
	}

	return status, true, nil
}
