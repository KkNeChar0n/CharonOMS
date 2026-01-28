package classify

import (
	"charonoms/internal/domain/classify/entity"
	"charonoms/internal/domain/classify/repository"

	"gorm.io/gorm"
)

// ClassifyRepositoryImpl 分类仓储实现
type ClassifyRepositoryImpl struct {
	db *gorm.DB
}

// NewClassifyRepository 创建分类仓储实例
func NewClassifyRepository(db *gorm.DB) repository.ClassifyRepository {
	return &ClassifyRepositoryImpl{db: db}
}

// GetAll 获取所有分类列表（含parent_name）
func (r *ClassifyRepositoryImpl) GetAll() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Table("classify c").
		Select(`c.id,
				c.name,
				c.level,
				c.parentid,
				parent.name as parent_name,
				c.status`).
		Joins("LEFT JOIN classify parent ON c.parentid = parent.id").
		Order("c.level ASC, c.id ASC").
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

// GetParents 获取所有一级分类
func (r *ClassifyRepositoryImpl) GetParents() ([]entity.Classify, error) {
	var classifies []entity.Classify

	err := r.db.Where("level = ?", 0).
		Order("id ASC").
		Find(&classifies).Error

	if err != nil {
		return nil, err
	}

	// 确保返回空数组而不是nil
	if classifies == nil {
		classifies = []entity.Classify{}
	}

	return classifies, nil
}

// GetActive 获取启用状态的分类列表
func (r *ClassifyRepositoryImpl) GetActive() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Table("classify c").
		Select(`c.id,
				c.name,
				c.level,
				c.parentid,
				parent.name as parent_name,
				c.status`).
		Joins("LEFT JOIN classify parent ON c.parentid = parent.id").
		Where("c.status = ?", 0).
		Order("c.level ASC, c.id ASC").
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

// GetByID 根据ID查询分类
func (r *ClassifyRepositoryImpl) GetByID(id int) (*entity.Classify, error) {
	var classify entity.Classify

	err := r.db.Where("id = ?", id).First(&classify).Error
	if err != nil {
		return nil, err
	}

	return &classify, nil
}

// CheckNameUnique 检查名称在指定级别和父级下是否唯一
func (r *ClassifyRepositoryImpl) CheckNameUnique(name string, level int, parentID *int, excludeID int) (bool, error) {
	var count int64

	query := r.db.Model(&entity.Classify{}).
		Where("name = ? AND level = ?", name, level)

	// 对于二级分类，需要在同一父级下检查
	if level == 1 && parentID != nil {
		query = query.Where("parentid = ?", *parentID)
	}

	// 更新时排除自身
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// Create 创建分类
func (r *ClassifyRepositoryImpl) Create(classify *entity.Classify) error {
	return r.db.Create(classify).Error
}

// Update 更新分类信息
func (r *ClassifyRepositoryImpl) Update(classify *entity.Classify) error {
	return r.db.Model(classify).
		Where("id = ?", classify.ID).
		Updates(map[string]interface{}{
			"name":     classify.Name,
			"level":    classify.Level,
			"parentid": classify.ParentID,
		}).Error
}

// UpdateStatus 更新分类状态
func (r *ClassifyRepositoryImpl) UpdateStatus(id int, status int) error {
	return r.db.Model(&entity.Classify{}).
		Where("id = ?", id).
		Update("status", status).Error
}
