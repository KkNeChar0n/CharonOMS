package attribute

import (
	"charonoms/internal/domain/attribute/entity"
	"charonoms/internal/domain/attribute/repository"

	"gorm.io/gorm"
)

// AttributeRepositoryImpl 属性仓储实现
type AttributeRepositoryImpl struct {
	db *gorm.DB
}

// NewAttributeRepository 创建属性仓储实例
func NewAttributeRepository(db *gorm.DB) repository.AttributeRepository {
	return &AttributeRepositoryImpl{db: db}
}

// GetAll 获取所有属性列表（含value_count统计）
func (r *AttributeRepositoryImpl) GetAll() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Table("attribute a").
		Select(`a.id,
				a.name,
				a.classify,
				a.status,
				COALESCE(COUNT(av.id), 0) as value_count`).
		Joins("LEFT JOIN attribute_value av ON a.id = av.attributeid").
		Group("a.id").
		Order("a.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 确保返回空数组而不是nil
	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetActive 获取启用状态的属性列表（含嵌套values数组）
func (r *AttributeRepositoryImpl) GetActive() ([]map[string]interface{}, error) {
	// 查询所有启用的属性
	var attributes []entity.Attribute
	err := r.db.Where("status = ?", 0).Order("id DESC").Find(&attributes).Error
	if err != nil {
		return nil, err
	}

	// 构造返回结果
	results := make([]map[string]interface{}, 0, len(attributes))

	for _, attr := range attributes {
		// 查询该属性的所有属性值
		var values []entity.AttributeValue
		err := r.db.Where("attributeid = ?", attr.ID).Find(&values).Error
		if err != nil {
			return nil, err
		}

		// 确保values是空数组而不是nil
		if values == nil {
			values = []entity.AttributeValue{}
		}

		// 构造结果
		result := map[string]interface{}{
			"id":       attr.ID,
			"name":     attr.Name,
			"classify": attr.Classify,
			"values":   values,
		}
		results = append(results, result)
	}

	// 确保返回空数组而不是nil
	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetByID 根据ID查询属性
func (r *AttributeRepositoryImpl) GetByID(id int) (*entity.Attribute, error) {
	var attribute entity.Attribute
	err := r.db.First(&attribute, id).Error
	if err != nil {
		return nil, err
	}
	return &attribute, nil
}

// Create 创建属性
func (r *AttributeRepositoryImpl) Create(attribute *entity.Attribute) error {
	return r.db.Create(attribute).Error
}

// Update 更新属性信息
func (r *AttributeRepositoryImpl) Update(attribute *entity.Attribute) error {
	return r.db.Model(&entity.Attribute{}).
		Where("id = ?", attribute.ID).
		Updates(map[string]interface{}{
			"name":     attribute.Name,
			"classify": attribute.Classify,
		}).Error
}

// UpdateStatus 更新属性状态
func (r *AttributeRepositoryImpl) UpdateStatus(id int, status int) error {
	return r.db.Model(&entity.Attribute{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetValues 获取属性的所有属性值
func (r *AttributeRepositoryImpl) GetValues(attributeID int) ([]entity.AttributeValue, error) {
	var values []entity.AttributeValue
	err := r.db.Where("attributeid = ?", attributeID).Find(&values).Error
	if err != nil {
		return nil, err
	}

	// 确保返回空数组而不是nil
	if values == nil {
		values = []entity.AttributeValue{}
	}

	return values, nil
}

// SaveValues 保存属性值（DELETE+INSERT全量替换）
func (r *AttributeRepositoryImpl) SaveValues(attributeID int, values []entity.AttributeValue) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除该属性的所有旧属性值
		if err := tx.Where("attributeid = ?", attributeID).Delete(&entity.AttributeValue{}).Error; err != nil {
			return err
		}

		// 批量插入新属性值
		if len(values) > 0 {
			if err := tx.Create(&values).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
