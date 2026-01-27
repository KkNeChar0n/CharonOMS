package account

import (
	"charonoms/internal/application/service/auth"
	"charonoms/internal/domain/account/repository"
	"charonoms/internal/domain/auth/entity"
	"charonoms/pkg/errors"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// AccountService 账号应用服务
type AccountService struct {
	accountRepo repository.AccountRepository
}

// NewAccountService 创建账号服务实例
func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// GetAccountListResponse 账号列表响应
type GetAccountListResponse struct {
	Accounts []*AccountDTO `json:"accounts"`
}

// AccountDTO 账号数据传输对象
type AccountDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
	Status   int8   `json:"status"`
}

// GetAccountList 获取账号列表
func (s *AccountService) GetAccountList(ctx context.Context) (*GetAccountListResponse, error) {
	accounts, err := s.accountRepo.GetAccountList(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	accountDTOs := make([]*AccountDTO, 0, len(accounts))
	for _, account := range accounts {
		dto := &AccountDTO{
			ID:       account.ID,
			Username: account.Username,
			Name:     account.Name,
			Phone:    account.Phone,
			RoleID:   account.RoleID,
			Status:   account.Status,
		}
		if account.Role != nil {
			dto.RoleName = account.Role.Name
		}
		accountDTOs = append(accountDTOs, dto)
	}

	return &GetAccountListResponse{
		Accounts: accountDTOs,
	}, nil
}

// CreateAccountRequest 创建账号请求
type CreateAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Phone    string `json:"phone" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
	Status   int8   `json:"status"`
}

// CreateAccount 创建账号
func (s *AccountService) CreateAccount(ctx context.Context, req *CreateAccountRequest) error {
	// 检查用户名是否已存在
	exists, err := s.accountRepo.CheckUsernameExists(ctx, req.Username, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.BadRequest("用户名已存在")
	}

	// 检查手机号是否已存在
	exists, err = s.accountRepo.CheckPhoneExists(ctx, req.Phone, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.BadRequest("手机号已存在")
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建账号
	account := &entity.UserAccount{
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
		Status:   req.Status,
	}

	return s.accountRepo.CreateAccount(ctx, account)
}

// UpdateAccountRequest 更新账号请求
type UpdateAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name"`
	Phone    string `json:"phone" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
	Status   int8   `json:"status"`
}

// UpdateAccount 更新账号
func (s *AccountService) UpdateAccount(ctx context.Context, id uint, req *UpdateAccountRequest) error {
	// 检查账号是否存在
	account, err := s.accountRepo.GetAccountByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFound("账号不存在")
		}
		return err
	}

	// 检查用户名是否已被其他账号使用
	if req.Username != account.Username {
		exists, err := s.accountRepo.CheckUsernameExists(ctx, req.Username, id)
		if err != nil {
			return err
		}
		if exists {
			return errors.BadRequest("用户名已存在")
		}
	}

	// 检查手机号是否已被其他账号使用
	if req.Phone != account.Phone {
		exists, err := s.accountRepo.CheckPhoneExists(ctx, req.Phone, id)
		if err != nil {
			return err
		}
		if exists {
			return errors.BadRequest("手机号已存在")
		}
	}

	// 更新账号信息
	account.Username = req.Username
	account.Name = req.Name
	account.Phone = req.Phone
	account.RoleID = req.RoleID
	account.Status = req.Status

	return s.accountRepo.UpdateAccount(ctx, account)
}

// UpdateAccountStatus 更新账号状态
func (s *AccountService) UpdateAccountStatus(ctx context.Context, id uint, status int8) error {
	// 检查账号是否存在
	_, err := s.accountRepo.GetAccountByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFound("账号不存在")
		}
		return err
	}

	return s.accountRepo.UpdateAccountStatus(ctx, id, status)
}
