package basic

import (
	"charonoms/internal/domain/basic/entity"
	"charonoms/internal/domain/basic/repository"
	"context"
)

// BasicService 基础数据服务
type BasicService struct {
	basicRepo repository.BasicRepository
}

// NewBasicService 创建基础数据服务实例
func NewBasicService(basicRepo repository.BasicRepository) *BasicService {
	return &BasicService{
		basicRepo: basicRepo,
	}
}

// GetAllSexes 获取所有性别
func (s *BasicService) GetAllSexes(ctx context.Context) ([]*entity.Sex, error) {
	return s.basicRepo.GetAllSexes(ctx)
}

// GetActiveGrades 获取启用的年级
func (s *BasicService) GetActiveGrades(ctx context.Context) ([]*entity.Grade, error) {
	return s.basicRepo.GetActiveGrades(ctx)
}

// GetActiveSubjects 获取启用的学科
func (s *BasicService) GetActiveSubjects(ctx context.Context) ([]*entity.Subject, error) {
	return s.basicRepo.GetActiveSubjects(ctx)
}
