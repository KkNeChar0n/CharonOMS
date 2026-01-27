package entity

import "time"

// Role 角色实体
type Role struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:50;not null" json:"name"`
	Description  string    `gorm:"size:255" json:"description"`
	IsSuperAdmin int8      `gorm:"default:0" json:"is_super_admin"` // 0-否 1-是
	Status       int8      `gorm:"default:0" json:"status"`          // 0-正常 1-禁用
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:role_id;References:ID;joinReferences:permissions_id" json:"permissions,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}

// Permission 权限实体
type Permission struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	ActionID string `gorm:"size:100;uniqueIndex" json:"action_id"` // 如: view_student, add_student
	MenuID   uint   `gorm:"default:0" json:"menu_id"`
	Status   int8   `gorm:"default:0" json:"status"` // 0-启用 1-禁用

	// 关联
	Menu *Menu `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// Menu 菜单实体
type Menu struct {
	ID        uint   `gorm:"column:id;primaryKey" json:"id"`
	Name      string `gorm:"column:name;size:100;not null" json:"name"`
	Route     string `gorm:"column:route;size:100" json:"route"`
	ParentID  *uint  `gorm:"column:parent_id" json:"parent_id"` // 使用指针类型支持NULL
	SortOrder int    `gorm:"column:sort_order;default:0" json:"sort_order"`
	Status    int    `gorm:"column:status;default:0" json:"status"` // 0-启用 1-禁用

	// 关联
	Parent   *Menu   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []*Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menu"
}

// RolePermission 角色权限关联实体
type RolePermission struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	RoleID        uint `gorm:"not null;index:idx_role_permission" json:"role_id"`
	PermissionsID uint `gorm:"not null;index:idx_role_permission" json:"permissions_id"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// MenuTree 菜单树结构（用于前端展示）
type MenuTree struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Route     string      `json:"route"`
	ParentID  *uint       `json:"parent_id"`
	SortOrder int         `json:"sort_order"`
	Status    int         `json:"status"`
	Children  []*MenuTree `json:"children,omitempty"`
}
