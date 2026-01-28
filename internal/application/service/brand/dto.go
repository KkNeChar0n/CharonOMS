package brand

import "charonoms/internal/domain/brand/entity"

// BrandListResponse 品牌列表响应
type BrandListResponse struct {
	Brands []entity.Brand `json:"brands"`
}

// CreateBrandRequest 创建品牌请求
type CreateBrandRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateBrandRequest 更新品牌请求
type UpdateBrandRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status *int `json:"status" binding:"required,oneof=0 1"`
}

// MessageResponse 通用消息响应
type MessageResponse struct {
	Message string `json:"message"`
}
