package marketing

import (
	"context"
	"errors"
	"time"

	"charonoms/internal/domain/marketing/entity"
	"charonoms/internal/domain/marketing/repository"
)

// ActivityService 活动服务
type ActivityService struct {
	activityRepo repository.ActivityRepository
	templateRepo repository.TemplateRepository
}

// NewActivityService 创建活动服务实例
func NewActivityService(activityRepo repository.ActivityRepository, templateRepo repository.TemplateRepository) *ActivityService {
	return &ActivityService{
		activityRepo: activityRepo,
		templateRepo: templateRepo,
	}
}

// ActivityListRequest 活动列表请求
type ActivityListRequest struct {
	ID         uint   `form:"id"`
	Name       string `form:"name"`
	TemplateID uint   `form:"template_id"`
	Status     *int   `form:"status"`
}

// ActivityDetailRequest 活动细节请求
type ActivityDetailRequest struct {
	GoodsID  uint    `json:"goods_id" binding:"required"`
	Discount float64 `json:"discount" binding:"required"`
}

// CreateActivityRequest 创建活动请求
type CreateActivityRequest struct {
	TemplateID uint                        `json:"template_id" binding:"required"`
	Name       string                      `json:"name" binding:"required"`
	StartTime  time.Time                   `json:"start_time" binding:"required"`
	EndTime    time.Time                   `json:"end_time" binding:"required"`
	Details    []ActivityDetailRequest     `json:"details" binding:"required,min=1"`
}

// UpdateActivityRequest 更新活动请求
type UpdateActivityRequest struct {
	TemplateID uint                        `json:"template_id" binding:"required"`
	Name       string                      `json:"name" binding:"required"`
	StartTime  time.Time                   `json:"start_time" binding:"required"`
	EndTime    time.Time                   `json:"end_time" binding:"required"`
	Details    []ActivityDetailRequest     `json:"details" binding:"required,min=1"`
}

// GetActivityList 获取活动列表
func (s *ActivityService) GetActivityList(ctx context.Context, req *ActivityListRequest) ([]*entity.Activity, error) {
	return s.activityRepo.List(ctx, req.ID, req.Name, req.TemplateID, req.Status)
}

// GetActivityByID 获取活动详情
func (s *ActivityService) GetActivityByID(ctx context.Context, id uint) (*entity.Activity, error) {
	return s.activityRepo.GetByID(ctx, id)
}

// CreateActivity 创建活动
func (s *ActivityService) CreateActivity(ctx context.Context, req *CreateActivityRequest) error {
	// 验证时间范围
	if !req.StartTime.Before(req.EndTime) {
		return errors.New("开始时间MUST早于结束时间")
	}

	// 验证活动细节
	if len(req.Details) == 0 {
		return errors.New("MUST至少包含一个活动细节")
	}

	// 验证折扣值
	for _, detail := range req.Details {
		if detail.Discount <= 0 {
			return errors.New("活动细节中的折扣值MUST大于0")
		}
	}

	// 验证活动模板存在且启用
	template, err := s.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return errors.New("活动模板不存在")
	}
	if template.Status != 0 {
		return errors.New("活动模板未启用")
	}

	// 构建实体
	activity := &entity.Activity{
		TemplateID: req.TemplateID,
		Name:       req.Name,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Status:     1, // 默认禁用
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	// 构建活动细节
	for _, detail := range req.Details {
		activity.Details = append(activity.Details, &entity.ActivityDetail{
			GoodsID:  detail.GoodsID,
			Discount: detail.Discount,
		})
	}

	return s.activityRepo.Create(ctx, activity)
}

// UpdateActivity 更新活动
func (s *ActivityService) UpdateActivity(ctx context.Context, id uint, req *UpdateActivityRequest) error {
	// 获取活动
	activity, err := s.activityRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证状态
	if activity.Status == 0 {
		return errors.New("活动启用中，无法编辑")
	}

	// 验证时间范围
	if !req.StartTime.Before(req.EndTime) {
		return errors.New("开始时间MUST早于结束时间")
	}

	// 验证活动细节
	if len(req.Details) == 0 {
		return errors.New("MUST至少包含一个活动细节")
	}

	// 验证折扣值
	for _, detail := range req.Details {
		if detail.Discount <= 0 {
			return errors.New("活动细节中的折扣值MUST大于0")
		}
	}

	// 验证活动模板存在且启用
	template, err := s.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return errors.New("活动模板不存在")
	}
	if template.Status != 0 {
		return errors.New("活动模板未启用")
	}

	// 更新实体
	activity.TemplateID = req.TemplateID
	activity.Name = req.Name
	activity.StartTime = req.StartTime
	activity.EndTime = req.EndTime
	activity.UpdateTime = time.Now()

	// 清空旧细节
	activity.Details = nil

	// 构建新细节
	for _, detail := range req.Details {
		activity.Details = append(activity.Details, &entity.ActivityDetail{
			GoodsID:  detail.GoodsID,
			Discount: detail.Discount,
		})
	}

	return s.activityRepo.Update(ctx, activity)
}

// UpdateActivityStatus 更新活动状态
func (s *ActivityService) UpdateActivityStatus(ctx context.Context, id uint, status int) error {
	// 验证状态值
	if status != 0 && status != 1 {
		return errors.New("状态必须为0-启用或1-禁用")
	}

	return s.activityRepo.UpdateStatus(ctx, id, status)
}
