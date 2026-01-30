package activity_template

import (
	"context"

	"gorm.io/gorm"

	"charonoms/internal/domain/activity_template/entity"
	"charonoms/internal/domain/activity_template/repository"
	"charonoms/internal/interfaces/http/dto"
)

// Service 活动模板应用服务
type Service struct {
	repo repository.ActivityTemplateRepository
	db   *gorm.DB
}

// NewService 创建活动模板服务实例
func NewService(repo repository.ActivityTemplateRepository, db *gorm.DB) *Service {
	return &Service{
		repo: repo,
		db:   db,
	}
}

// CreateTemplate 创建活动模板
func (s *Service) CreateTemplate(ctx context.Context, req *dto.CreateActivityTemplateDTO) (int, error) {
	// 验证关联配置
	if req.SelectType == 1 && len(req.ClassifyIDs) == 0 {
		return 0, entity.ErrMissingRelationConfig
	}
	if req.SelectType == 2 && len(req.GoodsIDs) == 0 {
		return 0, entity.ErrMissingRelationConfig
	}

	// 创建模板实体
	template := &entity.ActivityTemplate{
		Name:       req.Name,
		Type:       req.Type,
		SelectType: req.SelectType,
		Status:     req.Status,
	}

	// 验证模板
	if err := template.Validate(); err != nil {
		return 0, err
	}

	// 创建关联实体
	var goods []*entity.ActivityTemplateGoods
	if req.SelectType == 1 {
		for _, classifyID := range req.ClassifyIDs {
			goods = append(goods, &entity.ActivityTemplateGoods{
				ClassifyID: &classifyID,
			})
		}
	} else {
		for _, goodsID := range req.GoodsIDs {
			goods = append(goods, &entity.ActivityTemplateGoods{
				GoodsID: &goodsID,
			})
		}
	}

	// 保存
	if err := s.repo.Create(ctx, template, goods); err != nil {
		return 0, err
	}

	return template.ID, nil
}

// UpdateTemplate 更新活动模板
func (s *Service) UpdateTemplate(ctx context.Context, id int, req *dto.UpdateActivityTemplateDTO) error {
	// 查询模板
	template, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查状态
	if template.IsEnabled() {
		return entity.ErrTemplateIsEnabled
	}

	// 验证关联配置
	if req.SelectType == 1 && len(req.ClassifyIDs) == 0 {
		return entity.ErrMissingRelationConfig
	}
	if req.SelectType == 2 && len(req.GoodsIDs) == 0 {
		return entity.ErrMissingRelationConfig
	}

	// 更新模板
	template.Name = req.Name
	template.Type = req.Type
	template.SelectType = req.SelectType

	// 验证模板
	if err := template.Validate(); err != nil {
		return err
	}

	// 创建关联实体
	var goods []*entity.ActivityTemplateGoods
	if req.SelectType == 1 {
		for _, classifyID := range req.ClassifyIDs {
			goods = append(goods, &entity.ActivityTemplateGoods{
				ClassifyID: &classifyID,
			})
		}
	} else {
		for _, goodsID := range req.GoodsIDs {
			goods = append(goods, &entity.ActivityTemplateGoods{
				GoodsID: &goodsID,
			})
		}
	}

	// 保存
	return s.repo.UpdateWithGoods(ctx, template, goods)
}

// DeleteTemplate 删除活动模板
func (s *Service) DeleteTemplate(ctx context.Context, id int) error {
	// 检查是否有关联活动
	count, err := s.repo.CountRelatedActivities(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return entity.ErrTemplateHasActivities
	}

	return s.repo.Delete(ctx, id)
}

// GetTemplate 获取活动模板详情
func (s *Service) GetTemplate(ctx context.Context, id int) (*dto.ActivityTemplateDetailDTO, error) {
	template, goods, err := s.repo.FindByIDWithGoods(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &dto.ActivityTemplateDetailDTO{
		ID:         template.ID,
		Name:       template.Name,
		Type:       template.Type,
		SelectType: template.SelectType,
		Status:     template.Status,
		CreateTime: template.CreateTime,
		UpdateTime: template.UpdateTime,
	}

	// 查询关联数据
	if template.SelectType == 1 {
		// 按分类选择
		result.ClassifyList = s.getClassifyList(ctx, goods)
	} else {
		// 按商品选择
		result.GoodsList = s.getGoodsList(ctx, goods)
	}

	return result, nil
}

// ListTemplates 查询活动模板列表
func (s *Service) ListTemplates(ctx context.Context) ([]*dto.ActivityTemplateDTO, error) {
	templates, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dto.ActivityTemplateDTO
	for _, t := range templates {
		result = append(result, &dto.ActivityTemplateDTO{
			ID:         t.ID,
			Name:       t.Name,
			Type:       t.Type,
			SelectType: t.SelectType,
			Status:     t.Status,
			CreateTime: t.CreateTime,
			UpdateTime: t.UpdateTime,
		})
	}

	return result, nil
}

// ListActiveTemplates 查询启用的活动模板
func (s *Service) ListActiveTemplates(ctx context.Context) ([]*dto.ActivityTemplateDTO, error) {
	templates, err := s.repo.FindActiveTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dto.ActivityTemplateDTO
	for _, t := range templates {
		result = append(result, &dto.ActivityTemplateDTO{
			ID:         t.ID,
			Name:       t.Name,
			Type:       t.Type,
			SelectType: t.SelectType,
			Status:     t.Status,
			CreateTime: t.CreateTime,
			UpdateTime: t.UpdateTime,
		})
	}

	return result, nil
}

// UpdateTemplateStatus 更新活动模板状态
func (s *Service) UpdateTemplateStatus(ctx context.Context, id int, status int) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

// getClassifyList 获取分类列表
func (s *Service) getClassifyList(ctx context.Context, goods []*entity.ActivityTemplateGoods) []dto.ClassifyRelationDTO {
	var result []dto.ClassifyRelationDTO
	for _, g := range goods {
		if g.ClassifyID != nil {
			var classify struct {
				ID   int
				Name string
			}
			// 只查询启用的分类
			s.db.Table("classify").Select("id, name").Where("id = ? AND status = ?", *g.ClassifyID, 0).First(&classify)
			// 只有查询到启用的分类才添加到结果中
			if classify.ID > 0 {
				result = append(result, dto.ClassifyRelationDTO{
					ClassifyID:   *g.ClassifyID,
					ClassifyName: classify.Name,
				})
			}
		}
	}
	return result
}

// getGoodsList 获取商品列表
func (s *Service) getGoodsList(ctx context.Context, goods []*entity.ActivityTemplateGoods) []dto.GoodsRelationDTO {
	var result []dto.GoodsRelationDTO
	for _, g := range goods {
		if g.GoodsID != nil {
			var item struct {
				GoodsID      int
				GoodsName    string
				Price        float64
				BrandName    string
				ClassifyName string
			}
			// 只查询启用的商品
			s.db.Table("goods g").
				Select("g.id as goods_id, g.name as goods_name, g.price, b.name as brand_name, c.name as classify_name").
				Joins("LEFT JOIN brand b ON g.brandid = b.id").
				Joins("LEFT JOIN classify c ON g.classifyid = c.id").
				Where("g.id = ? AND g.status = ?", *g.GoodsID, 0).
				First(&item)
			// 只有查询到启用的商品才添加到结果中
			if item.GoodsID > 0 {
				result = append(result, dto.GoodsRelationDTO{
					GoodsID:      *g.GoodsID,
					GoodsName:    item.GoodsName,
					Price:        item.Price,
					BrandName:    item.BrandName,
					ClassifyName: item.ClassifyName,
				})
			}
		}
	}
	return result
}
