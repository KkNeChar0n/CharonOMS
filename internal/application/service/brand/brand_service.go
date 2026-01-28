package brand

import (
	"fmt"
	"strings"

	"charonoms/internal/domain/brand/entity"
	"charonoms/internal/domain/brand/repository"
)

// BrandService 品牌业务服务
type BrandService struct {
	brandRepo repository.BrandRepository
}

// NewBrandService 创建品牌业务服务实例
func NewBrandService(brandRepo repository.BrandRepository) *BrandService {
	return &BrandService{
		brandRepo: brandRepo,
	}
}

// GetBrandList 获取所有品牌列表
func (s *BrandService) GetBrandList() ([]entity.Brand, error) {
	return s.brandRepo.GetAll()
}

// GetActiveBrands 获取启用状态的品牌列表
func (s *BrandService) GetActiveBrands() ([]entity.Brand, error) {
	return s.brandRepo.GetActive()
}

// CreateBrand 创建品牌
func (s *BrandService) CreateBrand(req *CreateBrandRequest) error {
	// 验证品牌名称不能为空
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return fmt.Errorf("品牌名称不能为空")
	}

	// 检查品牌名称是否已存在
	existingBrand, err := s.brandRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check brand name: %w", err)
	}
	if existingBrand != nil {
		return fmt.Errorf("该品牌名称已存在")
	}

	// 创建品牌
	_, err = s.brandRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}

	return nil
}

// UpdateBrand 更新品牌信息
func (s *BrandService) UpdateBrand(id int, req *UpdateBrandRequest) error {
	// 验证品牌名称不能为空
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return fmt.Errorf("品牌名称不能为空")
	}

	// 检查品牌是否存在
	brand, err := s.brandRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get brand: %w", err)
	}
	if brand == nil {
		return fmt.Errorf("品牌不存在")
	}

	// 如果名称发生变化，检查新名称是否已被其他品牌使用
	if brand.Name != name {
		existingBrand, err := s.brandRepo.GetByName(name)
		if err != nil {
			return fmt.Errorf("failed to check brand name: %w", err)
		}
		if existingBrand != nil && existingBrand.ID != id {
			return fmt.Errorf("该品牌名称已存在")
		}
	}

	// 更新品牌信息
	return s.brandRepo.Update(id, name)
}

// UpdateBrandStatus 更新品牌状态
func (s *BrandService) UpdateBrandStatus(id int, status int) error {
	// 验证状态值
	if status != 0 && status != 1 {
		return fmt.Errorf("状态不能为空")
	}

	// 检查品牌是否存在
	brand, err := s.brandRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get brand: %w", err)
	}
	if brand == nil {
		return fmt.Errorf("品牌不存在")
	}

	// 更新状态
	return s.brandRepo.UpdateStatus(id, status)
}
