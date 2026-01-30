package activity

import (
	"context"
	"time"

	"gorm.io/gorm"

	"charonoms/internal/domain/activity/entity"
	"charonoms/internal/domain/activity/repository"
	activityTemplateEntity "charonoms/internal/domain/activity_template/entity"
	activityTemplateRepo "charonoms/internal/domain/activity_template/repository"
	"charonoms/internal/interfaces/http/dto"
)

// Service 活动应用服务
type Service struct {
	repo         repository.ActivityRepository
	templateRepo activityTemplateRepo.ActivityTemplateRepository
	db           *gorm.DB
}

// NewService 创建活动服务实例
func NewService(repo repository.ActivityRepository, templateRepo activityTemplateRepo.ActivityTemplateRepository, db *gorm.DB) *Service {
	return &Service{
		repo:         repo,
		templateRepo: templateRepo,
		db:           db,
	}
}

// CreateActivity 创建活动
func (s *Service) CreateActivity(ctx context.Context, req *dto.CreateActivityDTO) (int, error) {
	// 验证模板是否存在且已启用
	template, err := s.templateRepo.FindByID(ctx, req.TemplateID)
	if err != nil {
		return 0, err
	}
	if !template.IsEnabled() {
		return 0, activityTemplateEntity.ErrTemplateIsDisabled
	}

	// 创建活动实体
	activity := &entity.Activity{
		Name:       req.Name,
		TemplateID: req.TemplateID,
		StartTime:  req.StartTime.Time,
		EndTime:    req.EndTime.Time,
		Status:     req.Status,
	}

	// 验证活动
	if err := activity.Validate(); err != nil {
		return 0, err
	}

	// 创建活动详情
	var details []*entity.ActivityDetail
	for _, d := range req.Details {
		detail := &entity.ActivityDetail{
			ThresholdAmount: d.ThresholdAmount,
			DiscountValue:   d.DiscountValue,
		}
		if err := detail.Validate(template.Type); err != nil {
			return 0, err
		}
		details = append(details, detail)
	}

	// 保存
	if err := s.repo.Create(ctx, activity, details); err != nil {
		return 0, err
	}

	return activity.ID, nil
}

// UpdateActivity 更新活动
func (s *Service) UpdateActivity(ctx context.Context, id int, req *dto.UpdateActivityDTO) error {
	// 查询活动
	activity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查状态
	if activity.IsEnabled() {
		return entity.ErrActivityIsEnabled
	}

	// 验证模板是否存在且已启用
	template, err := s.templateRepo.FindByID(ctx, req.TemplateID)
	if err != nil {
		return err
	}
	if !template.IsEnabled() {
		return activityTemplateEntity.ErrTemplateIsDisabled
	}

	// 更新活动
	activity.Name = req.Name
	activity.TemplateID = req.TemplateID
	activity.StartTime = req.StartTime.Time
	activity.EndTime = req.EndTime.Time
	activity.Status = req.Status

	// 验证活动
	if err := activity.Validate(); err != nil {
		return err
	}

	// 创建活动详情
	var details []*entity.ActivityDetail
	for _, d := range req.Details {
		detail := &entity.ActivityDetail{
			ThresholdAmount: d.ThresholdAmount,
			DiscountValue:   d.DiscountValue,
		}
		if err := detail.Validate(template.Type); err != nil {
			return err
		}
		details = append(details, detail)
	}

	// 保存
	return s.repo.UpdateWithDetails(ctx, activity, details)
}

// DeleteActivity 删除活动
func (s *Service) DeleteActivity(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// GetActivity 获取活动详情
func (s *Service) GetActivity(ctx context.Context, id int) (*dto.ActivityDetailResponseDTO, error) {
	activity, details, err := s.repo.FindByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	// 查询模板信息
	template, err := s.templateRepo.FindByID(ctx, activity.TemplateID)
	if err != nil {
		return nil, err
	}

	result := &dto.ActivityDetailResponseDTO{
		ID:           activity.ID,
		Name:         activity.Name,
		TemplateID:   activity.TemplateID,
		TemplateName: template.Name,
		TemplateType: template.Type,
		SelectType:   template.SelectType,
		StartTime:    activity.StartTime,
		EndTime:      activity.EndTime,
		Status:       activity.Status,
	}

	// 转换详情
	for _, d := range details {
		result.Details = append(result.Details, dto.ActivityDetailDTO{
			ID:              d.ID,
			ActivityID:      d.ActivityID,
			ThresholdAmount: d.ThresholdAmount,
			DiscountValue:   d.DiscountValue,
		})
	}

	return result, nil
}

// ListActivities 查询活动列表
func (s *Service) ListActivities(ctx context.Context) ([]*dto.ActivityDTO, error) {
	activities, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dto.ActivityDTO
	for _, a := range activities {
		// 查询模板信息
		template, err := s.templateRepo.FindByID(ctx, a.TemplateID)
		if err != nil {
			continue
		}

		result = append(result, &dto.ActivityDTO{
			ID:           a.ID,
			Name:         a.Name,
			TemplateID:   a.TemplateID,
			TemplateName: template.Name,
			TemplateType: template.Type,
			StartTime:    a.StartTime,
			EndTime:      a.EndTime,
			Status:       a.Status,
			CreateTime:   a.CreateTime,
		})
	}

	return result, nil
}

// GetActivitiesByDateRange 按日期范围查询活动（包含冲突检测）
func (s *Service) GetActivitiesByDateRange(ctx context.Context, paymentTime time.Time) (*dto.ActivitiesByDateRangeDTO, error) {
	// 查询日期范围内的活动
	activities, err := s.repo.FindByDateRange(ctx, paymentTime)
	if err != nil {
		return nil, err
	}

	result := &dto.ActivitiesByDateRangeDTO{
		Activities: make([]dto.ActivityDTO, 0),
	}

	// 检测是否有重复类型
	typeMap := make(map[int]bool)
	var duplicateType *int

	for _, a := range activities {
		// 查询模板信息
		template, err := s.templateRepo.FindByID(ctx, a.TemplateID)
		if err != nil {
			continue
		}

		// 检测重复类型
		if typeMap[template.Type] {
			result.HasDuplicate = true
			duplicateType = &template.Type
		}
		typeMap[template.Type] = true

		// 查询活动详情
		details, err := s.repo.FindDetailsByActivityID(ctx, a.ID)
		if err != nil {
			continue
		}

		activityDTO := dto.ActivityDTO{
			ID:                 a.ID,
			Name:               a.Name,
			TemplateID:         a.TemplateID,
			TemplateName:       template.Name,
			TemplateType:       template.Type,
			TemplateSelectType: template.SelectType,
			StartTime:          a.StartTime,
			EndTime:            a.EndTime,
			Status:             a.Status,
			CreateTime:         a.CreateTime,
		}

		// 转换详情
		for _, d := range details {
			activityDTO.Details = append(activityDTO.Details, dto.ActivityDetailDTO{
				ID:              d.ID,
				ActivityID:      d.ActivityID,
				ThresholdAmount: d.ThresholdAmount,
				DiscountValue:   d.DiscountValue,
			})
		}

		result.Activities = append(result.Activities, activityDTO)
	}

	// 设置冲突类型名称
	if result.HasDuplicate && duplicateType != nil {
		result.DuplicateType = duplicateType
		switch *duplicateType {
		case 1:
			result.TypeName = "满减"
		case 2:
			result.TypeName = "满折"
		case 3:
			result.TypeName = "满赠"
		}
	}

	return result, nil
}

// UpdateActivityStatus 更新活动状态
func (s *Service) UpdateActivityStatus(ctx context.Context, id int, status int) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
