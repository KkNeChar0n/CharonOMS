package coach

import (
	"fmt"

	"charonoms/internal/domain/coach/repository"
)

// CoachService 教练业务服务
type CoachService struct {
	coachRepo repository.CoachRepository
}

// NewCoachService 创建教练业务服务实例
func NewCoachService(coachRepo repository.CoachRepository) *CoachService {
	return &CoachService{
		coachRepo: coachRepo,
	}
}

// GetCoachList 获取教练列表
func (s *CoachService) GetCoachList() ([]map[string]interface{}, error) {
	return s.coachRepo.GetCoachList()
}

// GetActiveCoaches 获取启用教练列表
func (s *CoachService) GetActiveCoaches() ([]map[string]interface{}, error) {
	return s.coachRepo.GetActiveCoaches()
}

// GetCoachByID 获取教练详情
func (s *CoachService) GetCoachByID(id int) (map[string]interface{}, error) {
	return s.coachRepo.GetCoachByID(id)
}

// CreateCoach 创建教练
func (s *CoachService) CreateCoach(req *CreateCoachRequest) (int, error) {
	// 验证必填字段
	if req.CoachName == "" {
		return 0, fmt.Errorf("coach name is required")
	}
	if req.SexID == 0 {
		return 0, fmt.Errorf("sex_id is required")
	}
	if req.SubjectID == 0 {
		return 0, fmt.Errorf("subject_id is required")
	}
	if req.Phone == "" {
		return 0, fmt.Errorf("phone is required")
	}

	// 创建教练（默认状态为0-启用）
	coachID, err := s.coachRepo.CreateCoach(req.CoachName, req.SexID, req.SubjectID, req.Phone, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to create coach: %w", err)
	}

	// 如果有关联学生，添加关联
	if len(req.StudentIDs) > 0 {
		if err := s.coachRepo.AddCoachStudents(coachID, req.StudentIDs); err != nil {
			return 0, fmt.Errorf("failed to add coach-student associations: %w", err)
		}
	}

	return coachID, nil
}

// UpdateCoach 更新教练信息
func (s *CoachService) UpdateCoach(id int, req *UpdateCoachRequest) error {
	// 验证必填字段
	if req.CoachName == "" {
		return fmt.Errorf("coach name is required")
	}
	if req.SexID == 0 {
		return fmt.Errorf("sex_id is required")
	}
	if req.SubjectID == 0 {
		return fmt.Errorf("subject_id is required")
	}
	if req.Phone == "" {
		return fmt.Errorf("phone is required")
	}

	// 检查教练是否存在
	exists, err := s.coachRepo.CoachExists(id)
	if err != nil {
		return fmt.Errorf("failed to check coach existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("coach not found")
	}

	// 更新教练信息
	return s.coachRepo.UpdateCoach(id, req.CoachName, req.SexID, req.SubjectID, req.Phone)
}

// UpdateCoachStatus 更新教练状态
func (s *CoachService) UpdateCoachStatus(id int, status int) error {
	// 验证状态值
	if status != 0 && status != 1 {
		return fmt.Errorf("invalid status value, must be 0 or 1")
	}

	// 检查教练是否存在
	exists, err := s.coachRepo.CoachExists(id)
	if err != nil {
		return fmt.Errorf("failed to check coach existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("coach not found")
	}

	// 更新状态
	return s.coachRepo.UpdateCoachStatus(id, status)
}

// DeleteCoach 删除教练
func (s *CoachService) DeleteCoach(id int) error {
	// 检查教练是否存在
	exists, err := s.coachRepo.CoachExists(id)
	if err != nil {
		return fmt.Errorf("failed to check coach existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("coach not found")
	}

	// 删除教练（级联删除关联）
	return s.coachRepo.DeleteCoach(id)
}
