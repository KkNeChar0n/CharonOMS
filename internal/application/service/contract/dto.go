package contract

// ContractListResponse 合同列表响应
type ContractListResponse struct {
	Contracts []map[string]interface{} `json:"contracts"`
}

// CreateContractRequest 创建合同请求
type CreateContractRequest struct {
	Name           string  `json:"name"`
	StudentID      int     `json:"student_id"`
	Type           *int    `json:"type"`
	SignatureForm  *int    `json:"signature_form"`
	ContractAmount float64 `json:"contract_amount"`
	Signatory      string  `json:"signatory"`
}

// TerminateContractRequest 中止合作请求
type TerminateContractRequest struct {
	TerminationAgreement string `json:"termination_agreement"`
}

// CreateContractResponse 创建合同响应
type CreateContractResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// MessageResponse 通用消息响应
type MessageResponse struct {
	Message string `json:"message"`
}
