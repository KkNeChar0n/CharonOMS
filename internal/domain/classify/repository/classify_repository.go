package repository

import (
	"charonoms/internal/domain/classify/entity"
)

// ClassifyRepository 分类仓储接口
type ClassifyRepository interface {
	// GetAll 获取所有分类列表（含parent_name）
	GetAll() ([]map[string]interface{}, error)

	// GetParents 获取所有一级分类
	GetParents() ([]entity.Classify, error)

	// GetActive 获取启用状态的分类列表
	GetActive() ([]map[string]interface{}, error)

	// GetByID 根据ID查询分类
	GetByID(id int) (*entity.Classify, error)

	// CheckNameUnique 检查名称在指定级别和父级下是否唯一
	// level: 0=一级分类，1=二级分类
	// parentID: 二级分类时需要提供父级ID
	// excludeID: 更新时排除自身ID
	CheckNameUnique(name string, level int, parentID *int, excludeID int) (bool, error)

	// Create 创建分类
	Create(classify *entity.Classify) error

	// Update 更新分类信息
	Update(classify *entity.Classify) error

	// UpdateStatus 更新分类状态
	UpdateStatus(id int, status int) error
}
