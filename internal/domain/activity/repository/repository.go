package repository

import (
	"context"
	"time"

	"charonoms/internal/domain/activity/entity"
)

// ActivityRepository 活动仓储接口
type ActivityRepository interface {
	// Create 创建活动
	Create(ctx context.Context, activity *entity.Activity, details []*entity.ActivityDetail) error

	// Update 更新活动
	Update(ctx context.Context, activity *entity.Activity) error

	// UpdateWithDetails 更新活动及详情
	UpdateWithDetails(ctx context.Context, activity *entity.Activity, details []*entity.ActivityDetail) error

	// Delete 删除活动
	Delete(ctx context.Context, id int) error

	// FindByID 根据ID查询活动
	FindByID(ctx context.Context, id int) (*entity.Activity, error)

	// FindByIDWithDetails 根据ID查询活动及详情
	FindByIDWithDetails(ctx context.Context, id int) (*entity.Activity, []*entity.ActivityDetail, error)

	// List 查询所有活动
	List(ctx context.Context) ([]*entity.Activity, error)

	// FindByDateRange 根据日期范围查询启用的活动
	FindByDateRange(ctx context.Context, paymentTime time.Time) ([]*entity.Activity, error)

	// FindDetailsByActivityID 查询活动详情
	FindDetailsByActivityID(ctx context.Context, activityID int) ([]*entity.ActivityDetail, error)

	// UpdateStatus 更新活动状态
	UpdateStatus(ctx context.Context, id int, status int) error
}
