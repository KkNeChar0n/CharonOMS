package attribute

// CreateAttributeRequest 创建属性请求
type CreateAttributeRequest struct {
	Name     string      `json:"name" binding:"required"`
	Classify interface{} `json:"classify" binding:"required"` // 可以接收字符串或整数
}

// UpdateAttributeRequest 更新属性请求
type UpdateAttributeRequest struct {
	Name     string      `json:"name" binding:"required"`
	Classify interface{} `json:"classify" binding:"required"` // 可以接收字符串或整数
}

// UpdateAttributeStatusRequest 更新属性状态请求
type UpdateAttributeStatusRequest struct {
	Status *int `json:"status" binding:"required,oneof=0 1"` // 使用指针类型，可以区分0和未提供的值
}

// SaveValuesRequest 保存属性值请求
type SaveValuesRequest struct {
	Values []string `json:"values" binding:"required"` // 直接接收字符串数组
}

// AttributeListResponse 属性列表响应
type AttributeListResponse struct {
	Attributes []map[string]interface{} `json:"attributes"`
}

// ActiveAttributeResponse 启用属性列表响应
type ActiveAttributeResponse struct {
	Attributes []map[string]interface{} `json:"attributes"`
}

// AttributeValuesResponse 属性值列表响应
type AttributeValuesResponse struct {
	Values []map[string]interface{} `json:"values"`
}
