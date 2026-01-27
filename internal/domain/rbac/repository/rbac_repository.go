package repository

import (
	"charonoms/internal/domain/rbac/entity"
	"context"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *entity.Role) error
	// Update 更新角色
	Update(ctx context.Context, role *entity.Role) error
	// UpdateStatus 更新角色状态
	UpdateStatus(ctx context.Context, id uint, status int8) error
	// Delete 删除角色
	Delete(ctx context.Context, id uint) error
	// GetByID 根据ID获取角色
	GetByID(ctx context.Context, id uint) (*entity.Role, error)
	// List 获取角色列表
	List(ctx context.Context) ([]*entity.Role, error)
	// GetRolePermissions 获取角色的权限列表
	GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error)
	// UpdateRolePermissions 更新角色权限
	UpdateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	// List 获取权限列表
	List(ctx context.Context) ([]*entity.Permission, error)
	// ListByStatus 根据状态获取权限列表
	ListByStatus(ctx context.Context, status int8) ([]*entity.Permission, error)
	// GetByID 根据ID获取权限
	GetByID(ctx context.Context, id uint) (*entity.Permission, error)
	// UpdateStatus 更新权限状态
	UpdateStatus(ctx context.Context, id uint, status int8) error
	// GetTree 获取权限树
	GetTree(ctx context.Context) (interface{}, error)
}

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	// List 获取菜单列表
	List(ctx context.Context) ([]*entity.Menu, error)
	// GetByID 根据ID获取菜单
	GetByID(ctx context.Context, id uint) (*entity.Menu, error)
	// Update 更新菜单
	Update(ctx context.Context, menu *entity.Menu) error
	// UpdateStatus 更新菜单状态
	UpdateStatus(ctx context.Context, id uint, status int8) error
	// GetMenuTree 获取菜单树
	GetMenuTree(ctx context.Context) ([]*entity.MenuTree, error)
	// GetUserMenuTree 获取用户的菜单树（基于权限）
	GetUserMenuTree(ctx context.Context, roleID uint, isSuperAdmin bool) ([]*entity.MenuTree, error)
}
