package account

import (
	"charonoms/internal/domain/account/repository"
	"charonoms/internal/domain/auth/entity"
	"context"

	"gorm.io/gorm"
)

// accountRepositoryImpl 账号仓储实现
type accountRepositoryImpl struct {
	db *gorm.DB
}

// NewAccountRepository 创建账号仓储实例
func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &accountRepositoryImpl{
		db: db,
	}
}

// GetAccountList 获取账号列表
func (r *accountRepositoryImpl) GetAccountList(ctx context.Context, filters map[string]interface{}) ([]*entity.UserAccount, error) {
	var accounts []*entity.UserAccount
	query := r.db.WithContext(ctx).Preload("Role")

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("id = ?", id)
	}
	if phone, ok := filters["phone"]; ok && phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone.(string)+"%")
	}
	if roleID, ok := filters["role_id"]; ok && roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	err := query.Order("id DESC").Find(&accounts).Error
	return accounts, err
}

// GetAccountByID 根据ID获取账号
func (r *accountRepositoryImpl) GetAccountByID(ctx context.Context, id uint) (*entity.UserAccount, error) {
	var account entity.UserAccount
	err := r.db.WithContext(ctx).
		Preload("Role").
		First(&account, id).Error
	return &account, err
}

// CreateAccount 创建账号
func (r *accountRepositoryImpl) CreateAccount(ctx context.Context, account *entity.UserAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// UpdateAccount 更新账号
func (r *accountRepositoryImpl) UpdateAccount(ctx context.Context, account *entity.UserAccount) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserAccount{}).
		Where("id = ?", account.ID).
		Updates(map[string]interface{}{
			"username": account.Username,
			"name":     account.Name,
			"phone":    account.Phone,
			"role_id":  account.RoleID,
			"status":   account.Status,
		}).Error
}

// UpdateAccountStatus 更新账号状态
func (r *accountRepositoryImpl) UpdateAccountStatus(ctx context.Context, id uint, status int8) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserAccount{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// CheckUsernameExists 检查用户名是否存在
func (r *accountRepositoryImpl) CheckUsernameExists(ctx context.Context, username string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.UserAccount{}).Where("username = ?", username)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// CheckPhoneExists 检查手机号是否存在
func (r *accountRepositoryImpl) CheckPhoneExists(ctx context.Context, phone string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.UserAccount{}).Where("phone = ?", phone)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
