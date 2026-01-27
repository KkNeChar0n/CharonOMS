package student

import (
	"charonoms/internal/domain/student/entity"
	"charonoms/internal/domain/student/repository"
	"errors"
	"fmt"
)

// StudentService 学生业务服务
type StudentService struct {
	repo repository.StudentRepository
}

// NewStudentService 创建学生服务实例
func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// GetStudentList 获取学生列表
func (s *StudentService) GetStudentList() (*StudentListResponse, error) {
	students, err := s.repo.GetStudentList()
	if err != nil {
		return nil, fmt.Errorf("获取学生列表失败: %w", err)
	}

	return &StudentListResponse{Students: students}, nil
}

// GetActiveStudents 获取启用学生列表
func (s *StudentService) GetActiveStudents() (*ActiveStudentResponse, error) {
	students, err := s.repo.GetActiveStudents()
	if err != nil {
		return nil, fmt.Errorf("获取启用学生列表失败: %w", err)
	}

	return &ActiveStudentResponse{Students: students}, nil
}

// GetStudentByID 获取学生详情
func (s *StudentService) GetStudentByID(id int) (*entity.Student, error) {
	student, err := s.repo.GetStudentByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取学生详情失败: %w", err)
	}
	return student, nil
}

// CreateStudent 创建学生
func (s *StudentService) CreateStudent(req *CreateStudentRequest) (int, error) {
	// 验证必填字段（由binding已处理）
	if req.Name == "" || req.SexID == 0 || req.GradeID == 0 || req.Phone == "" {
		return 0, errors.New("必填字段缺失")
	}

	// 创建学生实体
	student := &entity.Student{
		Name:    req.Name,
		SexID:   req.SexID,
		GradeID: req.GradeID,
		Phone:   req.Phone,
		Status:  0, // 默认启用
	}

	// 保存学生
	if err := s.repo.CreateStudent(student); err != nil {
		return 0, fmt.Errorf("创建学生失败: %w", err)
	}

	// 如果有教练关联，添加关联
	if len(req.CoachIDs) > 0 {
		if err := s.repo.AddStudentCoaches(student.ID, req.CoachIDs); err != nil {
			return 0, fmt.Errorf("添加学生教练关联失败: %w", err)
		}
	}

	return student.ID, nil
}

// UpdateStudent 更新学生信息
func (s *StudentService) UpdateStudent(id int, req *UpdateStudentRequest) error {
	// 验证必填字段
	if req.Name == "" || req.SexID == 0 || req.GradeID == 0 || req.Phone == "" {
		return errors.New("必填字段缺失")
	}

	// 检查学生是否存在
	_, err := s.repo.GetStudentByID(id)
	if err != nil {
		return fmt.Errorf("学生不存在: %w", err)
	}

	// 更新学生信息
	student := &entity.Student{
		ID:      id,
		Name:    req.Name,
		SexID:   req.SexID,
		GradeID: req.GradeID,
		Phone:   req.Phone,
	}

	if err := s.repo.UpdateStudent(student); err != nil {
		return fmt.Errorf("更新学生失败: %w", err)
	}

	return nil
}

// UpdateStudentStatus 更新学生状态
func (s *StudentService) UpdateStudentStatus(id int, req *UpdateStudentStatusRequest) error {
	// Status是指针类型，binding已经确保了它不为nil且值在0或1之间
	status := *req.Status

	// 更新状态
	if err := s.repo.UpdateStudentStatus(id, status); err != nil {
		return fmt.Errorf("更新学生状态失败: %w", err)
	}

	return nil
}

// DeleteStudent 删除学生
func (s *StudentService) DeleteStudent(id int) error {
	// 检查是否有关联订单
	hasOrders, err := s.repo.CheckStudentHasOrders(id)
	if err != nil {
		return fmt.Errorf("检查学生订单失败: %w", err)
	}

	if hasOrders {
		return errors.New("无法删除，该学生存在关联订单")
	}

	// 删除学生（级联删除关联）
	if err := s.repo.DeleteStudent(id); err != nil {
		return fmt.Errorf("删除学生失败: %w", err)
	}

	return nil
}
