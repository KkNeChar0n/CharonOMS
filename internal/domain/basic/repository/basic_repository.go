package repository

import (
	"charonoms/internal/domain/basic/entity"
	"context"
)

// BasicRepository 基础数据仓储接口
type BasicRepository interface {
	// GetAllSexes 获取所有性别
	GetAllSexes(ctx context.Context) ([]*entity.Sex, error)

	// GetActiveGrades 获取启用的年级
	GetActiveGrades(ctx context.Context) ([]*entity.Grade, error)

	// GetActiveSubjects 获取启用的学科
	GetActiveSubjects(ctx context.Context) ([]*entity.Subject, error)
}
