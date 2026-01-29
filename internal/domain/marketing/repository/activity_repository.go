package repository

import (
	"context"

	"charonoms/internal/domain/marketing/entity"
)

// ActivityRepository 活动仓储接口
type ActivityRepository interface {
	// List 查询活动列表，支持筛选
	List(ctx context.Context, id uint, name string, templateID uint, status *int) ([]*entity.Activity, error)

	// GetByID 根据ID获取活动详情（包含细节）
	GetByID(ctx context.Context, id uint) (*entity.Activity, error)

	// Create 创建活动（包含细节）
	Create(ctx context.Context, activity *entity.Activity) error

	// Update 更新活动（包含细节）
	Update(ctx context.Context, activity *entity.Activity) error

	// UpdateStatus 更新活动状态
	UpdateStatus(ctx context.Context, id uint, status int) error
}
