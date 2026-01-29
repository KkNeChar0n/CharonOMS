package marketing

import (
	"context"

	"gorm.io/gorm"

	"charonoms/internal/domain/marketing/entity"
	"charonoms/internal/domain/marketing/repository"
)

// activityRepositoryImpl 活动仓储实现
type activityRepositoryImpl struct {
	db *gorm.DB
}

// NewActivityRepository 创建活动仓储实例
func NewActivityRepository(db *gorm.DB) repository.ActivityRepository {
	return &activityRepositoryImpl{db: db}
}

// List 查询活动列表，支持筛选
func (r *activityRepositoryImpl) List(ctx context.Context, id uint, name string, templateID uint, status *int) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	query := r.db.WithContext(ctx)

	if id != 0 {
		query = query.Where("activity.id = ?", id)
	}
	if name != "" {
		query = query.Where("activity.name LIKE ?", "%"+name+"%")
	}
	if templateID != 0 {
		query = query.Where("activity.template_id = ?", templateID)
	}
	if status != nil {
		query = query.Where("activity.status = ?", *status)
	}

	// 预加载模板信息
	err := query.Preload("Template").Find(&activities).Error
	if err != nil {
		return nil, err
	}

	// 填充模板名称和类型
	for _, activity := range activities {
		if activity.Template != nil {
			activity.TemplateName = activity.Template.Name
			activity.TemplateType = activity.Template.Type
		}
	}

	return activities, nil
}

// GetByID 根据ID获取活动详情（包含细节）
func (r *activityRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Activity, error) {
	var activity entity.Activity

	// 预加载模板和细节
	err := r.db.WithContext(ctx).Preload("Template").First(&activity, id).Error
	if err != nil {
		return nil, err
	}

	// 填充模板名称和类型
	if activity.Template != nil {
		activity.TemplateName = activity.Template.Name
		activity.TemplateType = activity.Template.Type
	}

	// 加载活动细节，并获取商品名称
	var details []*entity.ActivityDetail
	err = r.db.WithContext(ctx).
		Table("activity_detail").
		Select("activity_detail.id, activity_detail.activity_id, activity_detail.goods_id, activity_detail.discount, goods.name as goods_name").
		Joins("LEFT JOIN goods ON goods.id = activity_detail.goods_id").
		Where("activity_detail.activity_id = ?", id).
		Scan(&details).Error
	if err != nil {
		return nil, err
	}
	activity.Details = details

	return &activity, nil
}

// Create 创建活动（包含细节）
func (r *activityRepositoryImpl) Create(ctx context.Context, activity *entity.Activity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建活动
		if err := tx.Create(activity).Error; err != nil {
			return err
		}

		// 创建活动细节
		if len(activity.Details) > 0 {
			for _, detail := range activity.Details {
				detail.ActivityID = activity.ID
			}
			if err := tx.Create(&activity.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Update 更新活动（包含细节）
func (r *activityRepositoryImpl) Update(ctx context.Context, activity *entity.Activity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新活动基本信息
		if err := tx.Model(&entity.Activity{}).Where("id = ?", activity.ID).Updates(map[string]interface{}{
			"template_id": activity.TemplateID,
			"name":        activity.Name,
			"start_time":  activity.StartTime,
			"end_time":    activity.EndTime,
		}).Error; err != nil {
			return err
		}

		// 删除旧的活动细节
		if err := tx.Where("activity_id = ?", activity.ID).Delete(&entity.ActivityDetail{}).Error; err != nil {
			return err
		}

		// 创建新的活动细节
		if len(activity.Details) > 0 {
			for _, detail := range activity.Details {
				detail.ActivityID = activity.ID
			}
			if err := tx.Create(&activity.Details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateStatus 更新活动状态
func (r *activityRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.db.WithContext(ctx).Model(&entity.Activity{}).Where("id = ?", id).Update("status", status).Error
}
