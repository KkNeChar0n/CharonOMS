package contract

// ContractListResponse 合同列表响应
type ContractListResponse struct {
	Contracts []map[string]interface{} `json:"contracts"`
}

// CreateContractRequest 创建合同请求
// 前端通过Vue select/input发送的数字字段为字符串，使用 interface{} 接收后在service层转换
type CreateContractRequest struct {
	Name           string      `json:"name"`
	StudentID      interface{} `json:"student_id"`
	Type           interface{} `json:"type"`
	SignatureForm  interface{} `json:"signature_form"`
	ContractAmount interface{} `json:"contract_amount"`
	Signatory      string      `json:"signatory"`
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
