package persistence

import (
	"context"

	"gorm.io/gorm"

	"charonoms/internal/domain/activity_template/entity"
	"charonoms/internal/domain/activity_template/repository"
)

// GormActivityTemplateRepository GORM实现的活动模板仓储
type GormActivityTemplateRepository struct {
	db *gorm.DB
}

// NewActivityTemplateRepository 创建活动模板仓储实例
func NewActivityTemplateRepository(db *gorm.DB) repository.ActivityTemplateRepository {
	return &GormActivityTemplateRepository{db: db}
}

// Create 创建活动模板
func (r *GormActivityTemplateRepository) Create(ctx context.Context, template *entity.ActivityTemplate, goods []*entity.ActivityTemplateGoods) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建模板
		if err := tx.Create(template).Error; err != nil {
			return err
		}

		// 创建关联数据
		if len(goods) > 0 {
			for _, g := range goods {
				g.TemplateID = template.ID
			}
			if err := tx.Create(&goods).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Update 更新活动模板
func (r *GormActivityTemplateRepository) Update(ctx context.Context, template *entity.ActivityTemplate) error {
	return r.db.WithContext(ctx).Model(template).Updates(template).Error
}

// UpdateWithGoods 更新活动模板及关联商品/分类
func (r *GormActivityTemplateRepository) UpdateWithGoods(ctx context.Context, template *entity.ActivityTemplate, goods []*entity.ActivityTemplateGoods) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新模板
		if err := tx.Model(template).Updates(template).Error; err != nil {
			return err
		}

		// 删除旧的关联数据
		if err := tx.Where("template_id = ?", template.ID).Delete(&entity.ActivityTemplateGoods{}).Error; err != nil {
			return err
		}

		// 创建新的关联数据
		if len(goods) > 0 {
			for _, g := range goods {
				g.TemplateID = template.ID
			}
			if err := tx.Create(&goods).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete 删除活动模板
func (r *GormActivityTemplateRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entity.ActivityTemplate{}, id).Error
}

// FindByID 根据ID查询活动模板
func (r *GormActivityTemplateRepository) FindByID(ctx context.Context, id int) (*entity.ActivityTemplate, error) {
	var template entity.ActivityTemplate
	err := r.db.WithContext(ctx).First(&template, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrTemplateNotFound
		}
		return nil, err
	}
	return &template, nil
}

// FindByIDWithGoods 根据ID查询活动模板及关联商品/分类
func (r *GormActivityTemplateRepository) FindByIDWithGoods(ctx context.Context, id int) (*entity.ActivityTemplate, []*entity.ActivityTemplateGoods, error) {
	template, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	var goods []*entity.ActivityTemplateGoods
	if err := r.db.WithContext(ctx).Where("template_id = ?", id).Find(&goods).Error; err != nil {
		return nil, nil, err
	}

	return template, goods, nil
}

// List 查询所有活动模板
func (r *GormActivityTemplateRepository) List(ctx context.Context) ([]*entity.ActivityTemplate, error) {
	var templates []*entity.ActivityTemplate
	err := r.db.WithContext(ctx).Order("id DESC").Find(&templates).Error
	return templates, err
}

// FindActiveTemplates 查询启用的活动模板
func (r *GormActivityTemplateRepository) FindActiveTemplates(ctx context.Context) ([]*entity.ActivityTemplate, error) {
	var templates []*entity.ActivityTemplate
	err := r.db.WithContext(ctx).Where("status = ?", 0).Order("id DESC").Find(&templates).Error
	return templates, err
}

// CountRelatedActivities 统计关联的活动数量
func (r *GormActivityTemplateRepository) CountRelatedActivities(ctx context.Context, templateID int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("activity").Where("template_id = ?", templateID).Count(&count).Error
	return count, err
}

// UpdateStatus 更新活动模板状态
func (r *GormActivityTemplateRepository) UpdateStatus(ctx context.Context, id int, status int) error {
	return r.db.WithContext(ctx).Model(&entity.ActivityTemplate{}).Where("id = ?", id).Update("status", status).Error
}
