package persistence

import (
	"charonoms/internal/domain/basic/entity"
	"charonoms/internal/domain/basic/repository"
	"context"

	"gorm.io/gorm"
)

// BasicRepositoryImpl 基础数据仓储实现
type BasicRepositoryImpl struct {
	db *gorm.DB
}

// NewBasicRepository 创建基础数据仓储实例
func NewBasicRepository(db *gorm.DB) repository.BasicRepository {
	return &BasicRepositoryImpl{
		db: db,
	}
}

// GetAllSexes 获取所有性别
func (r *BasicRepositoryImpl) GetAllSexes(ctx context.Context) ([]*entity.Sex, error) {
	var sexes []*entity.Sex
	err := r.db.WithContext(ctx).Order("id").Find(&sexes).Error
	return sexes, err
}

// GetActiveGrades 获取启用的年级
func (r *BasicRepositoryImpl) GetActiveGrades(ctx context.Context) ([]*entity.Grade, error) {
	var grades []*entity.Grade
	err := r.db.WithContext(ctx).
		Where("status = ?", 0).
		Order("id").
		Find(&grades).Error
	return grades, err
}

// GetActiveSubjects 获取启用的学科
func (r *BasicRepositoryImpl) GetActiveSubjects(ctx context.Context) ([]*entity.Subject, error) {
	var subjects []*entity.Subject
	err := r.db.WithContext(ctx).
		Where("status = ?", 0).
		Order("id").
		Find(&subjects).Error
	return subjects, err
}
