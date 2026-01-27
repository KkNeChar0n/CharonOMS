package rbac

import (
	"charonoms/internal/domain/rbac/entity"
	"charonoms/internal/domain/rbac/repository"
	"context"

	"gorm.io/gorm"
)

// RoleRepositoryImpl 角色仓储实现
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create 创建角色
func (r *RoleRepositoryImpl) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// Update 更新角色
func (r *RoleRepositoryImpl) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Model(role).Updates(role).Error
}

// UpdateStatus 更新角色状态
func (r *RoleRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status int8) error {
	return r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除角色
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Role{}, id).Error
}

// GetByID 根据ID获取角色
func (r *RoleRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *RoleRepositoryImpl) List(ctx context.Context) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.WithContext(ctx).Order("id ASC").Find(&roles).Error
	return roles, err
}

// GetRolePermissions 获取角色的权限列表
func (r *RoleRepositoryImpl) GetRolePermissions(ctx context.Context, roleID uint) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Table("role_permissions rp").
		Select("p.*").
		Joins("JOIN permissions p ON rp.permissions_id = p.id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// UpdateRolePermissions 更新角色权限
func (r *RoleRepositoryImpl) UpdateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧的权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.RolePermission{}).Error; err != nil {
			return err
		}

		// 添加新的权限关联
		if len(permissionIDs) > 0 {
			rolePermissions := make([]entity.RolePermission, len(permissionIDs))
			for i, permID := range permissionIDs {
				rolePermissions[i] = entity.RolePermission{
					RoleID:        roleID,
					PermissionsID: permID,
				}
			}
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
