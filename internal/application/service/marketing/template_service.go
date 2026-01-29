package marketing

import (
	"context"
	"errors"
	"time"

	"charonoms/internal/domain/marketing/entity"
	"charonoms/internal/domain/marketing/repository"
)

// TemplateService 活动模板服务
type TemplateService struct {
	templateRepo repository.TemplateRepository
}

// NewTemplateService 创建活动模板服务实例
func NewTemplateService(templateRepo repository.TemplateRepository) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
	}
}

// TemplateListRequest 活动模板列表请求
type TemplateListRequest struct {
	ID     uint   `form:"id"`
	Name   string `form:"name"`
	Type   *int   `form:"type"`
	Status *int   `form:"status"`
}

// CreateTemplateRequest 创建活动模板请求
type CreateTemplateRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	SelectType  int    `json:"select_type" binding:"required"`
	ClassifyIDs []uint `json:"classify_ids"`
	GoodsIDs    []uint `json:"goods_ids"`
}

// UpdateTemplateRequest 更新活动模板请求
type UpdateTemplateRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	SelectType  int    `json:"select_type" binding:"required"`
	ClassifyIDs []uint `json:"classify_ids"`
	GoodsIDs    []uint `json:"goods_ids"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status int `json:"status" binding:"required"`
}

// GetTemplateList 获取活动模板列表
func (s *TemplateService) GetTemplateList(ctx context.Context, req *TemplateListRequest) ([]*entity.ActivityTemplate, error) {
	return s.templateRepo.List(ctx, req.ID, req.Name, req.Type, req.Status)
}

// GetActiveTemplateList 获取启用的活动模板列表
func (s *TemplateService) GetActiveTemplateList(ctx context.Context) ([]*entity.ActivityTemplate, error) {
	return s.templateRepo.GetActiveList(ctx)
}

// GetTemplateByID 获取活动模板详情
func (s *TemplateService) GetTemplateByID(ctx context.Context, id uint) (*entity.ActivityTemplate, error) {
	return s.templateRepo.GetByID(ctx, id)
}

// CreateTemplate 创建活动模板
func (s *TemplateService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) error {
	// 验证活动类型
	if req.Type < 1 || req.Type > 3 {
		return errors.New("活动类型必须为1-满减、2-满折、3-满赠")
	}

	// 验证选择方式
	if req.SelectType < 1 || req.SelectType > 2 {
		return errors.New("选择方式必须为1-按分类、2-按商品")
	}

	// 验证关联配置
	if req.SelectType == 1 && len(req.ClassifyIDs) == 0 {
		return errors.New("按分类选择时MUST提供分类ID列表")
	}
	if req.SelectType == 2 && len(req.GoodsIDs) == 0 {
		return errors.New("按商品选择时MUST提供商品ID列表")
	}

	// 构建实体
	template := &entity.ActivityTemplate{
		Name:       req.Name,
		Type:       req.Type,
		SelectType: req.SelectType,
		Status:     1, // 默认禁用
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	return s.templateRepo.Create(ctx, template, req.ClassifyIDs, req.GoodsIDs)
}

// UpdateTemplate 更新活动模板
func (s *TemplateService) UpdateTemplate(ctx context.Context, id uint, req *UpdateTemplateRequest) error {
	// 获取模板
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证状态
	if template.Status == 0 {
		return errors.New("活动模板启用中，无法编辑")
	}

	// 验证活动类型
	if req.Type < 1 || req.Type > 3 {
		return errors.New("活动类型必须为1-满减、2-满折、3-满赠")
	}

	// 验证选择方式
	if req.SelectType < 1 || req.SelectType > 2 {
		return errors.New("选择方式必须为1-按分类、2-按商品")
	}

	// 验证关联配置
	if req.SelectType == 1 && len(req.ClassifyIDs) == 0 {
		return errors.New("按分类选择时MUST提供分类ID列表")
	}
	if req.SelectType == 2 && len(req.GoodsIDs) == 0 {
		return errors.New("按商品选择时MUST提供商品ID列表")
	}

	// 更新实体
	template.Name = req.Name
	template.Type = req.Type
	template.SelectType = req.SelectType
	template.UpdateTime = time.Now()

	return s.templateRepo.Update(ctx, template, req.ClassifyIDs, req.GoodsIDs)
}

// UpdateTemplateStatus 更新活动模板状态
func (s *TemplateService) UpdateTemplateStatus(ctx context.Context, id uint, status int) error {
	// 验证状态值
	if status != 0 && status != 1 {
		return errors.New("状态必须为0-启用或1-禁用")
	}

	return s.templateRepo.UpdateStatus(ctx, id, status)
}
