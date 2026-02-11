package repository

import (
	"charonoms/internal/domain/auth/entity"
	"context"
)

// AccountRepository 账号仓储接口
type AccountRepository interface {
	// GetAccountList 获取账号列表
	GetAccountList(ctx context.Context, filters map[string]interface{}) ([]*entity.UserAccount, error)

	// GetAccountByID 根据ID获取账号
	GetAccountByID(ctx context.Context, id uint) (*entity.UserAccount, error)

	// CreateAccount 创建账号
	CreateAccount(ctx context.Context, account *entity.UserAccount) error

	// UpdateAccount 更新账号
	UpdateAccount(ctx context.Context, account *entity.UserAccount) error

	// UpdateAccountStatus 更新账号状态
	UpdateAccountStatus(ctx context.Context, id uint, status int8) error

	// CheckUsernameExists 检查用户名是否存在
	CheckUsernameExists(ctx context.Context, username string, excludeID uint) (bool, error)

	// CheckPhoneExists 检查手机号是否存在
	CheckPhoneExists(ctx context.Context, phone string, excludeID uint) (bool, error)
}
