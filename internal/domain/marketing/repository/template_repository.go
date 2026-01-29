package repository

import (
	"context"

	"charonoms/internal/domain/marketing/entity"
)

// TemplateRepository 活动模板仓储接口
type TemplateRepository interface {
	// List 查询活动模板列表，支持筛选
	List(ctx context.Context, id uint, name string, typ *int, status *int) ([]*entity.ActivityTemplate, error)

	// GetActiveList 获取启用的活动模板列表
	GetActiveList(ctx context.Context) ([]*entity.ActivityTemplate, error)

	// GetByID 根据ID获取活动模板详情（包含关联）
	GetByID(ctx context.Context, id uint) (*entity.ActivityTemplate, error)

	// Create 创建活动模板（包含关联）
	Create(ctx context.Context, template *entity.ActivityTemplate) error

	// Update 更新活动模板（包含关联）
	Update(ctx context.Context, template *entity.ActivityTemplate) error

	// UpdateStatus 更新活动模板状态
	UpdateStatus(ctx context.Context, id uint, status int) error
}
