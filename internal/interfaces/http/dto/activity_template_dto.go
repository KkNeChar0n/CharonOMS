package dto

import "time"

// CreateActivityTemplateDTO 创建活动模板请求
type CreateActivityTemplateDTO struct {
	Name        string `json:"name" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	SelectType  int    `json:"select_type" binding:"required"`
	ClassifyIDs []int  `json:"classify_ids"`
	GoodsIDs    []int  `json:"goods_ids"`
	Status      int    `json:"status"`
}

// UpdateActivityTemplateDTO 更新活动模板请求
type UpdateActivityTemplateDTO struct {
	Name        string `json:"name" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	SelectType  int    `json:"select_type" binding:"required"`
	ClassifyIDs []int  `json:"classify_ids"`
	GoodsIDs    []int  `json:"goods_ids"`
	Status      int    `json:"status"`
}

// UpdateTemplateStatusDTO 更新模板状态请求
type UpdateTemplateStatusDTO struct {
	Status *int `json:"status" binding:"required,oneof=0 1"`
}

// ActivityTemplateDTO 活动模板响应
type ActivityTemplateDTO struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Type       int       `json:"type"`
	SelectType int       `json:"select_type"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// ClassifyRelationDTO 分类关联响应
type ClassifyRelationDTO struct {
	ClassifyID   int    `json:"classify_id"`
	ClassifyName string `json:"classify_name"`
}

// GoodsRelationDTO 商品关联响应
type GoodsRelationDTO struct {
	GoodsID      int     `json:"goods_id"`
	GoodsName    string  `json:"goods_name"`
	Price        float64 `json:"price"`
	BrandName    string  `json:"brand_name"`
	ClassifyName string  `json:"classify_name"`
}

// ActivityTemplateDetailDTO 活动模板详情响应
type ActivityTemplateDetailDTO struct {
	ID           int                   `json:"id"`
	Name         string                `json:"name"`
	Type         int                   `json:"type"`
	SelectType   int                   `json:"select_type"`
	Status       int                   `json:"status"`
	CreateTime   time.Time             `json:"create_time"`
	UpdateTime   time.Time             `json:"update_time"`
	ClassifyList []ClassifyRelationDTO `json:"classify_list,omitempty"`
	GoodsList    []GoodsRelationDTO    `json:"goods_list,omitempty"`
}
