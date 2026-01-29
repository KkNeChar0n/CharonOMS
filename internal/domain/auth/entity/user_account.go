package entity

// UserAccount 用户账号实体
type UserAccount struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"` // 不返回给前端
	Name     string `gorm:"size:100" json:"name"`
	Phone    string `gorm:"size:20;uniqueIndex" json:"phone"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	Status   int8   `gorm:"default:0" json:"status"` // 0-正常 1-禁用

	// 关联
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName 指定表名
func (UserAccount) TableName() string {
	return "useraccount"
}

// Role 角色实体（简化版）
type Role struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"size:50;not null" json:"name"`
	IsSuperAdmin int8   `gorm:"default:0" json:"is_super_admin"` // 0-否 1-是
	Status       int8   `gorm:"default:0" json:"status"`          // 0-正常 1-禁用
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}
