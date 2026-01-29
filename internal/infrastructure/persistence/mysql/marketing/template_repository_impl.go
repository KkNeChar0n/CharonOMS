package marketing

import (
	"context"

	"gorm.io/gorm"

	"charonoms/internal/domain/marketing/entity"
	"charonoms/internal/domain/marketing/repository"
)

// templateRepositoryImpl 活动模板仓储实现
type templateRepositoryImpl struct {
	db *gorm.DB
}

// NewTemplateRepository 创建活动模板仓储实例
func NewTemplateRepository(db *gorm.DB) repository.TemplateRepository {
	return &templateRepositoryImpl{db: db}
}

// List 查询活动模板列表，支持筛选
func (r *templateRepositoryImpl) List(ctx context.Context, id uint, name string, typ *int, status *int) ([]*entity.ActivityTemplate, error) {
	var templates []*entity.ActivityTemplate
	query := r.db.WithContext(ctx)

	if id != 0 {
		query = query.Where("id = ?", id)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if typ != nil {
		query = query.Where("type = ?", *typ)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Find(&templates).Error
	return templates, err
}

// GetActiveList 获取启用的活动模板列表
func (r *templateRepositoryImpl) GetActiveList(ctx context.Context) ([]*entity.ActivityTemplate, error) {
	var templates []*entity.ActivityTemplate
	err := r.db.WithContext(ctx).Where("status = ?", 0).Find(&templates).Error
	return templates, err
}

// GetByID 根据ID获取活动模板详情（包含关联）
func (r *templateRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.ActivityTemplate, error) {
	var template entity.ActivityTemplate

	// 先查询模板基本信息
	err := r.db.WithContext(ctx).First(&template, id).Error
	if err != nil {
		return nil, err
	}

	// 根据 select_type 动态加载关联
	if template.SelectType == 1 {
		// 按分类：预加载分类关联，并获取分类名称
		var classifies []*entity.ActivityTemplateClassify
		err = r.db.WithContext(ctx).
			Table("activity_template_classify").
			Select("activity_template_classify.id, activity_template_classify.template_id, activity_template_classify.classify_id, classify.name as classify_name").
			Joins("LEFT JOIN classify ON classify.id = activity_template_classify.classify_id").
			Where("activity_template_classify.template_id = ?", id).
			Scan(&classifies).Error
		if err != nil {
			return nil, err
		}
		template.Classifies = classifies
	} else if template.SelectType == 2 {
		// 按商品：预加载商品关联，并获取商品名称
		var goods []*entity.ActivityTemplateGoods
		err = r.db.WithContext(ctx).
			Table("activity_template_goods").
			Select("activity_template_goods.id, activity_template_goods.template_id, activity_template_goods.goods_id, goods.name as goods_name").
			Joins("LEFT JOIN goods ON goods.id = activity_template_goods.goods_id").
			Where("activity_template_goods.template_id = ?", id).
			Scan(&goods).Error
		if err != nil {
			return nil, err
		}
		template.Goods = goods
	}

	return &template, nil
}

// Create 创建活动模板（包含关联）
func (r *templateRepositoryImpl) Create(ctx context.Context, template *entity.ActivityTemplate) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建模板
		if err := tx.Create(template).Error; err != nil {
			return err
		}

		// 根据 select_type 创建关联
		if template.SelectType == 1 && len(template.Classifies) > 0 {
			// 按分类：创建分类关联
			for _, classify := range template.Classifies {
				classify.TemplateID = template.ID
			}
			if err := tx.Create(&template.Classifies).Error; err != nil {
				return err
			}
		} else if template.SelectType == 2 && len(template.Goods) > 0 {
			// 按商品：创建商品关联
			for _, goods := range template.Goods {
				goods.TemplateID = template.ID
			}
			if err := tx.Create(&template.Goods).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Update 更新活动模板（包含关联）
func (r *templateRepositoryImpl) Update(ctx context.Context, template *entity.ActivityTemplate) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新模板基本信息
		if err := tx.Model(&entity.ActivityTemplate{}).Where("id = ?", template.ID).Updates(map[string]interface{}{
			"name":        template.Name,
			"type":        template.Type,
			"select_type": template.SelectType,
		}).Error; err != nil {
			return err
		}

		// 删除旧关联
		if template.SelectType == 1 {
			if err := tx.Where("template_id = ?", template.ID).Delete(&entity.ActivityTemplateClassify{}).Error; err != nil {
				return err
			}
			// 创建新的分类关联
			if len(template.Classifies) > 0 {
				for _, classify := range template.Classifies {
					classify.TemplateID = template.ID
				}
				if err := tx.Create(&template.Classifies).Error; err != nil {
					return err
				}
			}
		} else if template.SelectType == 2 {
			if err := tx.Where("template_id = ?", template.ID).Delete(&entity.ActivityTemplateGoods{}).Error; err != nil {
				return err
			}
			// 创建新的商品关联
			if len(template.Goods) > 0 {
				for _, goods := range template.Goods {
					goods.TemplateID = template.ID
				}
				if err := tx.Create(&template.Goods).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// UpdateStatus 更新活动模板状态
func (r *templateRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.db.WithContext(ctx).Model(&entity.ActivityTemplate{}).Where("id = ?", id).Update("status", status).Error
}
