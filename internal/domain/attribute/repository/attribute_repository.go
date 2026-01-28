package repository

import (
	"charonoms/internal/domain/attribute/entity"
)

// AttributeRepository 属性仓储接口
type AttributeRepository interface {
	// GetAll 获取所有属性列表（含value_count统计）
	GetAll() ([]map[string]interface{}, error)

	// GetActive 获取启用状态的属性列表（含嵌套values数组）
	GetActive() ([]map[string]interface{}, error)

	// GetByID 根据ID查询属性
	GetByID(id int) (*entity.Attribute, error)

	// Create 创建属性
	Create(attribute *entity.Attribute) error

	// Update 更新属性信息
	Update(attribute *entity.Attribute) error

	// UpdateStatus 更新属性状态
	UpdateStatus(id int, status int) error

	// GetValues 获取属性的所有属性值
	GetValues(attributeID int) ([]entity.AttributeValue, error)

	// SaveValues 保存属性值（DELETE+INSERT全量替换）
	SaveValues(attributeID int, values []entity.AttributeValue) error
}
