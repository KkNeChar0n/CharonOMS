package rbac

import (
	"charonoms/internal/domain/rbac/entity"
	"charonoms/internal/domain/rbac/repository"
	"charonoms/pkg/errors"
	"context"

	"gorm.io/gorm"
)

// RBACService RBAC 应用服务
type RBACService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	menuRepo       repository.MenuRepository
}

// NewRBACService 创建 RBAC 服务实例
func NewRBACService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	menuRepo repository.MenuRepository,
) *RBACService {
	return &RBACService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		menuRepo:       menuRepo,
	}
}

// ===== 角色管理 =====

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"comment"` // 前端使用comment字段名
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"comment"` // 前端使用comment字段名
}

// CreateRole 创建角色（默认禁用，需显式启用才能使用）
func (s *RBACService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*entity.Role, error) {
	role := &entity.Role{
		Name:        req.Name,
		Description: req.Description,
		Status:      1,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// UpdateRole 更新角色（仅禁用状态可编辑）
func (s *RBACService) UpdateRole(ctx context.Context, id uint, req *UpdateRoleRequest) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	if role.Status == 0 {
		return errors.BadRequest("角色启用中，无法编辑")
	}

	role.Name = req.Name
	role.Description = req.Description

	return s.roleRepo.Update(ctx, role)
}

// UpdateRoleStatus 更新角色状态
func (s *RBACService) UpdateRoleStatus(ctx context.Context, id uint, status int8) error {
	return s.roleRepo.UpdateStatus(ctx, id, status)
}

// GetRoleList 获取角色列表
func (s *RBACService) GetRoleList(ctx context.Context) ([]*entity.Role, error) {
	return s.roleRepo.List(ctx)
}

// GetRoleByID 根据ID获取角色
func (s *RBACService) GetRoleByID(ctx context.Context, id uint) (*entity.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return role, nil
}

// ===== 权限管理 =====

// GetPermissionList 获取权限列表
func (s *RBACService) GetPermissionList(ctx context.Context) ([]*entity.Permission, error) {
	return s.permissionRepo.List(ctx)
}

// UpdatePermissionStatus 更新权限状态
func (s *RBACService) UpdatePermissionStatus(ctx context.Context, id uint, status int8) error {
	return s.permissionRepo.UpdateStatus(ctx, id, status)
}

// GetPermissionTree 获取权限树
func (s *RBACService) GetPermissionTree(ctx context.Context) (interface{}, error) {
	return s.permissionRepo.GetTree(ctx)
}

// ===== 角色权限关联 =====

// GetRolePermissions 获取角色的权限列表
func (s *RBACService) GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error) {
	return s.roleRepo.GetRolePermissions(ctx, roleID)
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

// UpdateRolePermissions 更新角色权限（仅禁用状态可修改）
func (s *RBACService) UpdateRolePermissions(ctx context.Context, roleID uint, req *UpdateRolePermissionsRequest) error {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	if role.Status == 0 {
		return errors.BadRequest("角色启用中，无法编辑")
	}

	return s.roleRepo.UpdateRolePermissions(ctx, roleID, req.PermissionIDs)
}

// ===== 菜单管理 =====

// GetMenuList 获取菜单列表
func (s *RBACService) GetMenuList(ctx context.Context) ([]*entity.Menu, error) {
	return s.menuRepo.List(ctx)
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	Name   string `json:"name" binding:"required"`
	Route  string `json:"route"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
}

// UpdateMenu 更新菜单（仅禁用状态可编辑，同级排序不可重复）
func (s *RBACService) UpdateMenu(ctx context.Context, id uint, req *UpdateMenuRequest) error {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	if menu.Status == 0 {
		return errors.BadRequest("菜单启用中，无法编辑")
	}

	// 校验同级菜单sort_order唯一性
	allMenus, err := s.menuRepo.List(ctx)
	if err != nil {
		return err
	}
	for _, m := range allMenus {
		if m.ID == id {
			continue
		}
		sameParent := (m.ParentID == nil && menu.ParentID == nil) ||
			(m.ParentID != nil && menu.ParentID != nil && *m.ParentID == *menu.ParentID)
		if sameParent && m.SortOrder == req.Sort {
			return errors.BadRequest("同级菜单中排序已存在")
		}
	}

	menu.Name = req.Name
	menu.Route = req.Route
	menu.SortOrder = req.Sort

	return s.menuRepo.Update(ctx, menu)
}

// UpdateMenuStatus 更新菜单状态
func (s *RBACService) UpdateMenuStatus(ctx context.Context, id uint, status int8) error {
	return s.menuRepo.UpdateStatus(ctx, id, status)
}

// GetMenuTree 获取菜单树
func (s *RBACService) GetMenuTree(ctx context.Context) ([]*entity.MenuTree, error) {
	return s.menuRepo.GetMenuTree(ctx)
}

// GetUserMenuTree 获取用户的菜单树（基于权限）
func (s *RBACService) GetUserMenuTree(ctx context.Context, roleID uint, isSuperAdmin bool) ([]*entity.MenuTree, error) {
	return s.menuRepo.GetUserMenuTree(ctx, roleID, isSuperAdmin)
}
