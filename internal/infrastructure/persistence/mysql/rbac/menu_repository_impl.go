package rbac

import (
	"charonoms/internal/domain/rbac/entity"
	"charonoms/internal/domain/rbac/repository"
	"context"

	"gorm.io/gorm"
)

// MenuRepositoryImpl 菜单仓储实现
type MenuRepositoryImpl struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储实例
func NewMenuRepository(db *gorm.DB) repository.MenuRepository {
	return &MenuRepositoryImpl{db: db}
}

// List 获取菜单列表
func (r *MenuRepositoryImpl) List(ctx context.Context, filters map[string]interface{}) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	query := r.db.WithContext(ctx).Preload("Parent")

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("id = ?", id)
	}
	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("sort_order ASC, id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 填充扁平化字段
	for _, menu := range menus {
		if menu.Parent != nil {
			menu.ParentName = menu.Parent.Name
		}
	}

	return menus, nil
}

// GetByID 根据ID获取菜单
func (r *MenuRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.db.WithContext(ctx).
		Preload("Parent").
		First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// Update 更新菜单
func (r *MenuRepositoryImpl) Update(ctx context.Context, menu *entity.Menu) error {
	return r.db.WithContext(ctx).Model(menu).Updates(menu).Error
}

// UpdateStatus 更新菜单状态
func (r *MenuRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status int8) error {
	return r.db.WithContext(ctx).
		Model(&entity.Menu{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetMenuTree 获取菜单树
func (r *MenuRepositoryImpl) GetMenuTree(ctx context.Context) ([]*entity.MenuTree, error) {
	var menus []*entity.Menu
	err := r.db.WithContext(ctx).
		Where("status = 0").
		Order("sort_order ASC, id ASC").
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	return buildMenuTree(menus), nil
}

// GetUserMenuTree 获取用户的菜单树（基于权限）
func (r *MenuRepositoryImpl) GetUserMenuTree(ctx context.Context, roleID uint, isSuperAdmin bool) ([]*entity.MenuTree, error) {
	var menus []*entity.Menu

	if isSuperAdmin {
		// 超级管理员查看所有菜单
		err := r.db.WithContext(ctx).
			Where("status = 0").
			Order("sort_order ASC, id ASC").
			Find(&menus).Error
		if err != nil {
			return nil, err
		}
	} else {
		// 普通用户根据权限查看菜单
		err := r.db.WithContext(ctx).
			Table("menu m").
			Select("DISTINCT m.*").
			Joins("JOIN permissions p ON m.id = p.menu_id").
			Joins("JOIN role_permissions rp ON p.id = rp.permissions_id").
			Where("rp.role_id = ? AND m.status = 0 AND p.status = 0", roleID).
			Order("m.sort_order ASC, m.id ASC").
			Find(&menus).Error
		if err != nil {
			return nil, err
		}

		// 补充父级菜单
		parentIDs := make(map[uint]bool)
		for _, menu := range menus {
			if menu.ParentID != nil && *menu.ParentID > 0 {
				parentIDs[*menu.ParentID] = true
			}
		}

		if len(parentIDs) > 0 {
			var parentMenus []*entity.Menu
			ids := make([]uint, 0, len(parentIDs))
			for id := range parentIDs {
				ids = append(ids, id)
			}

			err := r.db.WithContext(ctx).
				Where("id IN ? AND status = 0", ids).
				Find(&parentMenus).Error
			if err == nil {
				// 合并父级菜单
				menuMap := make(map[uint]bool)
				for _, m := range menus {
					menuMap[m.ID] = true
				}
				for _, pm := range parentMenus {
					if !menuMap[pm.ID] {
						menus = append(menus, pm)
					}
				}
			}
		}
	}

	return buildMenuTree(menus), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []*entity.Menu) []*entity.MenuTree {
	menuMap := make(map[uint]*entity.MenuTree)
	var rootMenus []*entity.MenuTree

	// 第一遍：创建所有菜单节点
	for _, menu := range menus {
		menuTree := &entity.MenuTree{
			ID:        menu.ID,
			Name:      menu.Name,
			Route:     menu.Route,
			ParentID:  menu.ParentID,
			SortOrder: menu.SortOrder,
			Status:    menu.Status,
			Children:  make([]*entity.MenuTree, 0),
		}
		menuMap[menu.ID] = menuTree

		if menu.ParentID == nil {
			rootMenus = append(rootMenus, menuTree)
		}
	}

	// 第二遍：建立父子关系
	for _, menu := range menus {
		if menu.ParentID != nil && *menu.ParentID > 0 {
			if parent, exists := menuMap[*menu.ParentID]; exists {
				parent.Children = append(parent.Children, menuMap[menu.ID])
			}
		}
	}

	return rootMenus
}
