package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"charonoms/internal/domain/activity/entity"
	"charonoms/internal/domain/activity/repository"
)

// GormActivityRepository GORM实现的活动仓储
type GormActivityRepository struct {
	db *gorm.DB
}

// NewActivityRepository 创建活动仓储实例
func NewActivityRepository(db *gorm.DB) repository.ActivityRepository {
	return &GormActivityRepository{db: db}
}

// Create 创建活动
func (r *GormActivityRepository) Create(ctx context.Context, activity *entity.Activity, details []*entity.ActivityDetail) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建活动
		if err := tx.Create(activity).Error; err != nil {
			return err
		}

		// 创建活动详情
		if len(details) > 0 {
			for _, d := range details {
				d.ActivityID = activity.ID
			}
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Update 更新活动
func (r *GormActivityRepository) Update(ctx context.Context, activity *entity.Activity) error {
	return r.db.WithContext(ctx).Model(activity).Updates(activity).Error
}

// UpdateWithDetails 更新活动及详情
func (r *GormActivityRepository) UpdateWithDetails(ctx context.Context, activity *entity.Activity, details []*entity.ActivityDetail) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新活动
		if err := tx.Model(activity).Updates(activity).Error; err != nil {
			return err
		}

		// 删除旧的详情数据
		if err := tx.Where("activity_id = ?", activity.ID).Delete(&entity.ActivityDetail{}).Error; err != nil {
			return err
		}

		// 创建新的详情数据
		if len(details) > 0 {
			for _, d := range details {
				d.ActivityID = activity.ID
			}
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete 删除活动
func (r *GormActivityRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entity.Activity{}, id).Error
}

// FindByID 根据ID查询活动
func (r *GormActivityRepository) FindByID(ctx context.Context, id int) (*entity.Activity, error) {
	var activity entity.Activity
	err := r.db.WithContext(ctx).First(&activity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrActivityNotFound
		}
		return nil, err
	}
	return &activity, nil
}

// FindByIDWithDetails 根据ID查询活动及详情
func (r *GormActivityRepository) FindByIDWithDetails(ctx context.Context, id int) (*entity.Activity, []*entity.ActivityDetail, error) {
	activity, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	var details []*entity.ActivityDetail
	if err := r.db.WithContext(ctx).Where("activity_id = ?", id).Order("threshold_amount ASC").Find(&details).Error; err != nil {
		return nil, nil, err
	}

	return activity, details, nil
}

// List 查询所有活动
func (r *GormActivityRepository) List(ctx context.Context) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	err := r.db.WithContext(ctx).Order("id DESC").Find(&activities).Error
	return activities, err
}

// FindByDateRange 根据日期范围查询启用的活动
func (r *GormActivityRepository) FindByDateRange(ctx context.Context, paymentTime time.Time) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	err := r.db.WithContext(ctx).
		Where("start_time <= ? AND end_time >= ? AND status = ?", paymentTime, paymentTime, 0).
		Order("id ASC").
		Find(&activities).Error
	return activities, err
}

// FindDetailsByActivityID 查询活动详情
func (r *GormActivityRepository) FindDetailsByActivityID(ctx context.Context, activityID int) ([]*entity.ActivityDetail, error) {
	var details []*entity.ActivityDetail
	err := r.db.WithContext(ctx).Where("activity_id = ?", activityID).Order("threshold_amount ASC").Find(&details).Error
	return details, err
}

// UpdateStatus 更新活动状态
func (r *GormActivityRepository) UpdateStatus(ctx context.Context, id int, status int) error {
	return r.db.WithContext(ctx).Model(&entity.Activity{}).Where("id = ?", id).Update("status", status).Error
}
