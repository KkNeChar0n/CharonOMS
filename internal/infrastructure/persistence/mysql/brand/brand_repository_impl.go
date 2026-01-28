package brand

import (
	"fmt"

	"charonoms/internal/domain/brand/entity"
	"charonoms/internal/domain/brand/repository"

	"gorm.io/gorm"
)

// BrandRepositoryImpl 品牌仓储实现
type BrandRepositoryImpl struct {
	db *gorm.DB
}

// NewBrandRepository 创建品牌仓储实例
func NewBrandRepository(db *gorm.DB) repository.BrandRepository {
	return &BrandRepositoryImpl{db: db}
}

// GetAll 获取所有品牌列表
func (r *BrandRepositoryImpl) GetAll() ([]entity.Brand, error) {
	var brands []entity.Brand

	if err := r.db.Order("id DESC").Find(&brands).Error; err != nil {
		return nil, fmt.Errorf("failed to get brand list: %w", err)
	}

	// 确保返回空数组而不是nil
	if brands == nil {
		brands = []entity.Brand{}
	}

	return brands, nil
}

// GetActive 获取启用状态的品牌列表
func (r *BrandRepositoryImpl) GetActive() ([]entity.Brand, error) {
	var brands []entity.Brand

	if err := r.db.Where("status = ?", 0).Order("id DESC").Find(&brands).Error; err != nil {
		return nil, fmt.Errorf("failed to get active brands: %w", err)
	}

	// 确保返回空数组而不是nil
	if brands == nil {
		brands = []entity.Brand{}
	}

	return brands, nil
}

// GetByID 根据ID获取品牌
func (r *BrandRepositoryImpl) GetByID(id int) (*entity.Brand, error) {
	var brand entity.Brand

	if err := r.db.Where("id = ?", id).First(&brand).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get brand by id: %w", err)
	}

	return &brand, nil
}

// GetByName 根据名称获取品牌
func (r *BrandRepositoryImpl) GetByName(name string) (*entity.Brand, error) {
	var brand entity.Brand

	if err := r.db.Where("name = ?", name).First(&brand).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get brand by name: %w", err)
	}

	return &brand, nil
}

// Create 创建品牌
func (r *BrandRepositoryImpl) Create(name string) (int, error) {
	brand := &entity.Brand{
		Name:   name,
		Status: 0, // 默认启用状态
	}

	if err := r.db.Create(brand).Error; err != nil {
		return 0, fmt.Errorf("failed to create brand: %w", err)
	}

	return brand.ID, nil
}

// Update 更新品牌信息
func (r *BrandRepositoryImpl) Update(id int, name string) error {
	result := r.db.Model(&entity.Brand{}).Where("id = ?", id).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("failed to update brand: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("brand not found")
	}

	return nil
}

// UpdateStatus 更新品牌状态
func (r *BrandRepositoryImpl) UpdateStatus(id int, status int) error {
	result := r.db.Model(&entity.Brand{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update brand status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("brand not found")
	}

	return nil
}
