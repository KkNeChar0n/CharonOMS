package repository

import (
	"context"

	"charonoms/internal/domain/activity_template/entity"
)

// ActivityTemplateRepository 活动模板仓储接口
type ActivityTemplateRepository interface {
	// Create 创建活动模板
	Create(ctx context.Context, template *entity.ActivityTemplate, goods []*entity.ActivityTemplateGoods) error

	// Update 更新活动模板
	Update(ctx context.Context, template *entity.ActivityTemplate) error

	// UpdateWithGoods 更新活动模板及关联商品/分类
	UpdateWithGoods(ctx context.Context, template *entity.ActivityTemplate, goods []*entity.ActivityTemplateGoods) error

	// Delete 删除活动模板
	Delete(ctx context.Context, id int) error

	// FindByID 根据ID查询活动模板
	FindByID(ctx context.Context, id int) (*entity.ActivityTemplate, error)

	// FindByIDWithGoods 根据ID查询活动模板及关联商品/分类
	FindByIDWithGoods(ctx context.Context, id int) (*entity.ActivityTemplate, []*entity.ActivityTemplateGoods, error)

	// List 查询所有活动模板
	List(ctx context.Context) ([]*entity.ActivityTemplate, error)

	// FindActiveTemplates 查询启用的活动模板
	FindActiveTemplates(ctx context.Context) ([]*entity.ActivityTemplate, error)

	// CountRelatedActivities 统计关联的活动数量
	CountRelatedActivities(ctx context.Context, templateID int) (int64, error)

	// UpdateStatus 更新活动模板状态
	UpdateStatus(ctx context.Context, id int, status int) error
}
