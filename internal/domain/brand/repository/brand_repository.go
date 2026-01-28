package repository

import "charonoms/internal/domain/brand/entity"

// BrandRepository 品牌仓储接口
type BrandRepository interface {
	// GetAll 获取所有品牌列表
	GetAll() ([]entity.Brand, error)

	// GetActive 获取启用状态的品牌列表
	GetActive() ([]entity.Brand, error)

	// GetByID 根据ID获取品牌
	GetByID(id int) (*entity.Brand, error)

	// GetByName 根据名称获取品牌
	GetByName(name string) (*entity.Brand, error)

	// Create 创建品牌
	Create(name string) (int, error)

	// Update 更新品牌信息
	Update(id int, name string) error

	// UpdateStatus 更新品牌状态
	UpdateStatus(id int, status int) error
}
