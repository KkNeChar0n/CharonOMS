package rbac

import (
	"charonoms/internal/domain/rbac/entity"
	"charonoms/internal/domain/rbac/repository"
	"context"

	"gorm.io/gorm"
)

// PermissionRepositoryImpl 权限仓储实现
type PermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓储实例
func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

// List 获取权限列表
func (r *PermissionRepositoryImpl) List(ctx context.Context) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Preload("Menu").
		Order("id ASC").
		Find(&permissions).Error
	return permissions, err
}

// ListByStatus 根据状态获取权限列表
func (r *PermissionRepositoryImpl) ListByStatus(ctx context.Context, status int8) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Preload("Menu").
		Where("status = ?", status).
		Order("id ASC").
		Find(&permissions).Error
	return permissions, err
}

// GetByID 根据ID获取权限
func (r *PermissionRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.WithContext(ctx).
		Preload("Menu").
		First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// UpdateStatus 更新权限状态
func (r *PermissionRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status int8) error {
	return r.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetTree 获取权限树
func (r *PermissionRepositoryImpl) GetTree(ctx context.Context) (interface{}, error) {
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Preload("Menu").
		Where("status = 0").
		Order("menu_id ASC, id ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	// 按菜单分组
	menuMap := make(map[uint][]map[string]interface{})
	for _, perm := range permissions {
		menuID := perm.MenuID
		permData := map[string]interface{}{
			"id":        perm.ID,
			"name":      perm.Name,
			"action_id": perm.ActionID,
			"status":    perm.Status,
		}
		menuMap[menuID] = append(menuMap[menuID], permData)
	}

	// 获取菜单列表
	var menus []*entity.Menu
	err = r.db.WithContext(ctx).
		Where("status = 0").
		Order("parent_id ASC, sort ASC, id ASC").
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	// 构建树结构
	result := make([]map[string]interface{}, 0)
	for _, menu := range menus {
		if menu.ParentID == nil { // 一级菜单
			menuData := map[string]interface{}{
				"id":          menu.ID,
				"name":        menu.Name,
				"route":       menu.Route,
				"permissions": menuMap[menu.ID],
				"children":    make([]map[string]interface{}, 0),
			}

			// 查找二级菜单
			for _, childMenu := range menus {
				if childMenu.ParentID != nil && *childMenu.ParentID == menu.ID {
					childData := map[string]interface{}{
						"id":          childMenu.ID,
						"name":        childMenu.Name,
						"route":       childMenu.Route,
						"permissions": menuMap[childMenu.ID],
					}
					menuData["children"] = append(menuData["children"].([]map[string]interface{}), childData)
				}
			}

			result = append(result, menuData)
		}
	}

	return result, nil
}
