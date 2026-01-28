package repository

// ContractRepository 合同仓储接口
type ContractRepository interface {
	// GetContractList 获取合同列表（含学生信息）
	GetContractList() ([]map[string]interface{}, error)

	// GetContractByID 根据ID获取合同详情（含学生信息）
	GetContractByID(id int) (map[string]interface{}, error)

	// CreateContract 创建合同
	CreateContract(name string, studentID int, contractType int, signatureForm int, contractAmount float64, signatory string, initiator string) (int, error)

	// RevokeContract 撤销合同（status=0 → status=98）
	RevokeContract(id int) error

	// TerminateContract 中止合作（status=50 → status=99）
	TerminateContract(id int, terminationAgreement string) error

	// GetContractStatus 获取合同当前状态
	GetContractStatus(id int) (int, bool, error)
}
