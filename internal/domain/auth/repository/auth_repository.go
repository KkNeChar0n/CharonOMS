package repository

import (
	"charonoms/internal/domain/auth/entity"
	"context"
)

// AuthRepository 认证仓储接口
type AuthRepository interface {
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*entity.UserAccount, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, id uint) (*entity.UserAccount, error)

	// UpdateUserRole 更新用户角色
	UpdateUserRole(ctx context.Context, userID, roleID uint) error
}
