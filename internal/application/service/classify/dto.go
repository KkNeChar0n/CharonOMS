package classify

// CreateClassifyRequest 创建分类请求
type CreateClassifyRequest struct {
	Name     string      `json:"name" binding:"required"`
	Level    interface{} `json:"level" binding:"required"` // 可以接收字符串或整数
	ParentID interface{} `json:"parent_id"`                // 可以接收字符串、整数或null
}

// UpdateClassifyRequest 更新分类请求
type UpdateClassifyRequest struct {
	Name     string      `json:"name" binding:"required"`
	Level    interface{} `json:"level" binding:"required"` // 可以接收字符串或整数
	ParentID interface{} `json:"parent_id"`                // 可以接收字符串、整数或null
}

// UpdateClassifyStatusRequest 更新分类状态请求
type UpdateClassifyStatusRequest struct {
	Status *int `json:"status" binding:"required,oneof=0 1"` // 使用指针类型，可以区分0和未提供的值
}

// ClassifyListResponse 分类列表响应
type ClassifyListResponse struct {
	Classifies []map[string]interface{} `json:"classifies"`
}

// ParentsListResponse 一级分类列表响应
type ParentsListResponse struct {
	Parents []map[string]interface{} `json:"parents"`
}

// ActiveClassifyResponse 启用分类列表响应
type ActiveClassifyResponse struct {
	Classifies []map[string]interface{} `json:"classifies"`
}
