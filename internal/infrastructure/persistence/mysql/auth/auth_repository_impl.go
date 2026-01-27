package auth

import (
	"charonoms/internal/domain/auth/entity"
	"charonoms/internal/domain/auth/repository"
	"context"

	"gorm.io/gorm"
)

// AuthRepositoryImpl 认证仓储实现
type AuthRepositoryImpl struct {
	db *gorm.DB
}

// NewAuthRepository 创建认证仓储实例
func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &AuthRepositoryImpl{db: db}
}

// GetUserByUsername 根据用户名获取用户
func (r *AuthRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*entity.UserAccount, error) {
	var user entity.UserAccount
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID 根据ID获取用户
func (r *AuthRepositoryImpl) GetUserByID(ctx context.Context, id uint) (*entity.UserAccount, error) {
	var user entity.UserAccount
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserRole 更新用户角色
func (r *AuthRepositoryImpl) UpdateUserRole(ctx context.Context, userID, roleID uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserAccount{}).
		Where("id = ?", userID).
		Update("role_id", roleID).Error
}
