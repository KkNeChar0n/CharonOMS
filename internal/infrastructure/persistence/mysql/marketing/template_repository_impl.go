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
		// 按分类：查询关联的分类
		var classifyList []map[string]interface{}
		err = r.db.WithContext(ctx).
			Table("activity_template_goods atg").
			Select("atg.classify_id, c.name as classify_name").
			Joins("JOIN classify c ON atg.classify_id = c.id").
			Where("atg.template_id = ?", id).
			Scan(&classifyList).Error
		if err != nil {
			return nil, err
		}
		template.ClassifyList = classifyList
	} else if template.SelectType == 2 {
		// 按商品：查询关联的商品
		var goodsList []map[string]interface{}
		err = r.db.WithContext(ctx).
			Table("activity_template_goods atg").
			Select("atg.goods_id, g.name as goods_name, g.price, b.name as brand_name, c.name as classify_name").
			Joins("JOIN goods g ON atg.goods_id = g.id").
			Joins("LEFT JOIN brand b ON g.brandid = b.id").
			Joins("LEFT JOIN classify c ON g.classifyid = c.id").
			Where("atg.template_id = ?", id).
			Scan(&goodsList).Error
		if err != nil {
			return nil, err
		}
		template.GoodsList = goodsList
	}

	return &template, nil
}

// Create 创建活动模板（包含关联）
func (r *templateRepositoryImpl) Create(ctx context.Context, template *entity.ActivityTemplate, classifyIDs []uint, goodsIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建模板
		if err := tx.Create(template).Error; err != nil {
			return err
		}

		// 根据 select_type 创建关联
		if template.SelectType == 1 && len(classifyIDs) > 0 {
			// 按分类：在 activity_template_goods 表中插入 classify_id
			for _, classifyID := range classifyIDs {
				relation := &entity.ActivityTemplateGoods{
					TemplateID: template.ID,
					ClassifyID: &classifyID,
				}
				if err := tx.Create(relation).Error; err != nil {
					return err
				}
			}
		} else if template.SelectType == 2 && len(goodsIDs) > 0 {
			// 按商品：在 activity_template_goods 表中插入 goods_id
			for _, goodsID := range goodsIDs {
				relation := &entity.ActivityTemplateGoods{
					TemplateID: template.ID,
					GoodsID:    &goodsID,
				}
				if err := tx.Create(relation).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Update 更新活动模板（包含关联）
func (r *templateRepositoryImpl) Update(ctx context.Context, template *entity.ActivityTemplate, classifyIDs []uint, goodsIDs []uint) error {
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
		if err := tx.Where("template_id = ?", template.ID).Delete(&entity.ActivityTemplateGoods{}).Error; err != nil {
			return err
		}

		// 根据 select_type 创建新关联
		if template.SelectType == 1 && len(classifyIDs) > 0 {
			// 按分类：在 activity_template_goods 表中插入 classify_id
			for _, classifyID := range classifyIDs {
				relation := &entity.ActivityTemplateGoods{
					TemplateID: template.ID,
					ClassifyID: &classifyID,
				}
				if err := tx.Create(relation).Error; err != nil {
					return err
				}
			}
		} else if template.SelectType == 2 && len(goodsIDs) > 0 {
			// 按商品：在 activity_template_goods 表中插入 goods_id
			for _, goodsID := range goodsIDs {
				relation := &entity.ActivityTemplateGoods{
					TemplateID: template.ID,
					GoodsID:    &goodsID,
				}
				if err := tx.Create(relation).Error; err != nil {
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
