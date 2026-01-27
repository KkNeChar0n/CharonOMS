package auth

import (
	"charonoms/internal/domain/auth/entity"
	"charonoms/internal/domain/auth/repository"
	"charonoms/internal/infrastructure/config"
	"charonoms/pkg/errors"
	"charonoms/pkg/jwt"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证应用服务
type AuthService struct {
	authRepo repository.AuthRepository
	jwtCfg   config.JWTConfig
}

// NewAuthService 创建认证服务实例
func NewAuthService(authRepo repository.AuthRepository, jwtCfg config.JWTConfig) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		jwtCfg:   jwtCfg,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string `json:"token"`
	Username     string `json:"username"`
	IsSuperAdmin bool   `json:"is_super_admin"`
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 查询用户
	user, err := s.authRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, err
	}

	// 验证密码（使用 bcrypt）
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// 密码错误
		return nil, errors.ErrInvalidCredentials
	}

	// 检查账号状态
	if user.Status == 1 {
		return nil, errors.ErrAccountDisabled
	}

	// 生成 JWT Token
	isSuperAdmin := false
	if user.Role != nil && user.Role.IsSuperAdmin == 1 {
		isSuperAdmin = true
	}

	token, err := jwt.GenerateToken(user.ID, user.RoleID, user.Username, isSuperAdmin, s.jwtCfg)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		Username:     user.Username,
		IsSuperAdmin: isSuperAdmin,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, userID uint) (*entity.UserAccount, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// SyncRoleResponse 同步角色响应
type SyncRoleResponse struct {
	RoleChanged  bool `json:"role_changed"`
	RoleID       uint `json:"role_id"`
	IsSuperAdmin bool `json:"is_super_admin"`
}

// SyncRole 同步用户角色信息
func (s *AuthService) SyncRole(ctx context.Context, userID uint, oldRoleID uint, oldIsSuperAdmin bool) (*SyncRoleResponse, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}

	newIsSuperAdmin := false
	if user.Role != nil && user.Role.IsSuperAdmin == 1 {
		newIsSuperAdmin = true
	}

	roleChanged := (oldRoleID != user.RoleID) || (oldIsSuperAdmin != newIsSuperAdmin)

	return &SyncRoleResponse{
		RoleChanged:  roleChanged,
		RoleID:       user.RoleID,
		IsSuperAdmin: newIsSuperAdmin,
	}, nil
}

// HashPassword 密码加密（工具方法）
// 使用 bcrypt 算法加密密码，cost = 10
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// VerifyPassword 验证密码（工具方法）
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
