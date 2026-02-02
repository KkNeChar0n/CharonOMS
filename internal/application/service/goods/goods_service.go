package goods

import (
	"fmt"
	"strconv"
	"strings"

	"charonoms/internal/domain/goods/repository"
)

// GoodsService 商品业务服务
type GoodsService struct {
	goodsRepo repository.GoodsRepository
}

// NewGoodsService 创建商品业务服务实例
func NewGoodsService(goodsRepo repository.GoodsRepository) *GoodsService {
	return &GoodsService{
		goodsRepo: goodsRepo,
	}
}

// toInt 将 interface{} 转换为 int（支持字符串和数字）
func toInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case float64:
		return int(val), true
	case float32:
		return int(val), true
	case string:
		i, err := strconv.Atoi(strings.TrimSpace(val))
		return i, err == nil
	case int:
		return val, true
	case int8:
		return int(val), true
	case int16:
		return int(val), true
	case int32:
		return int(val), true
	case int64:
		return int(val), true
	case uint:
		return int(val), true
	case uint8:
		return int(val), true
	case uint16:
		return int(val), true
	case uint32:
		return int(val), true
	case uint64:
		return int(val), true
	default:
		return 0, false
	}
}

// toFloat64 将 interface{} 转换为 float64（支持字符串和数字）
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		return f, err == nil
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	default:
		return 0, false
	}
}

// GetGoodsList 获取商品列表（支持按分类和状态过滤）
func (s *GoodsService) GetGoodsList(classifyID *int, status *int) ([]map[string]interface{}, error) {
	return s.goodsRepo.GetList(classifyID, status)
}

// GetGoodsByID 获取商品详情
func (s *GoodsService) GetGoodsByID(id int) (map[string]interface{}, error) {
	result, err := s.goodsRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		return nil, fmt.Errorf("商品不存在")
	}

	return result, nil
}

// GetActiveForOrder 获取可用于下单的商品列表
func (s *GoodsService) GetActiveForOrder() ([]map[string]interface{}, error) {
	return s.goodsRepo.GetActiveForOrder()
}

// GetAvailableForCombo 获取可用于组合的单品商品列表
func (s *GoodsService) GetAvailableForCombo(excludeID int) ([]map[string]interface{}, error) {
	return s.goodsRepo.GetAvailableForCombo(excludeID)
}

// GetIncludedGoods 根据父商品ID获取包含的子商品列表
func (s *GoodsService) GetIncludedGoods(parentsID int) ([]map[string]interface{}, error) {
	return s.goodsRepo.GetIncludedGoods(parentsID)
}

// GetTotalPrice 计算商品总价
func (s *GoodsService) GetTotalPrice(id int) (float64, error) {
	return s.goodsRepo.GetTotalPrice(id)
}

// CreateGoods 创建商品
func (s *GoodsService) CreateGoods(req *CreateGoodsRequest) (int, error) {
	// 转换并验证必填字段
	name := req.Name
	if name == "" {
		return 0, fmt.Errorf("商品信息不完整")
	}

	brandID, ok := toInt(req.BrandID)
	if !ok || brandID == 0 {
		return 0, fmt.Errorf("商品信息不完整")
	}

	classifyID, ok := toInt(req.ClassifyID)
	if !ok || classifyID == 0 {
		return 0, fmt.Errorf("商品信息不完整")
	}

	isGroup, ok := toInt(req.IsGroup)
	if !ok {
		return 0, fmt.Errorf("商品信息不完整")
	}

	price, ok := toFloat64(req.Price)
	if !ok || price == 0 {
		return 0, fmt.Errorf("商品信息不完整")
	}

	// 验证组合商品必须至少包含一个子商品
	if isGroup == 0 && len(req.IncludedGoodsIDs) == 0 {
		return 0, fmt.Errorf("组合商品必须至少包含一个子商品")
	}

	// 确保属性值ID数组不为nil
	attributeValueIDs := req.AttributeValueIDs
	if attributeValueIDs == nil {
		attributeValueIDs = []int{}
	}

	// 确保包含的商品ID数组不为nil
	includedGoodsIDs := req.IncludedGoodsIDs
	if includedGoodsIDs == nil {
		includedGoodsIDs = []int{}
	}

	goodsID, err := s.goodsRepo.Create(
		name,
		brandID,
		classifyID,
		isGroup,
		price,
		attributeValueIDs,
		includedGoodsIDs,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create goods: %w", err)
	}

	return goodsID, nil
}

// UpdateGoods 更新商品
func (s *GoodsService) UpdateGoods(id int, req *UpdateGoodsRequest) error {
	// 从数据库获取当前商品信息，读取isgroup（编辑接口不接受isgroup参数）
	current, err := s.goodsRepo.GetByID(id)
	if err != nil || current == nil {
		return fmt.Errorf("商品不存在")
	}

	// 读取当前isgroup值
	isGroup, ok := toInt(current["isgroup"])
	if !ok {
		return fmt.Errorf("商品数据格式错误")
	}

	// 转换并验证必填字段
	name := req.Name
	if name == "" {
		return fmt.Errorf("商品信息不完整")
	}

	brandID, ok := toInt(req.BrandID)
	if !ok || brandID == 0 {
		return fmt.Errorf("商品信息不完整")
	}

	classifyID, ok := toInt(req.ClassifyID)
	if !ok || classifyID == 0 {
		return fmt.Errorf("商品信息不完整")
	}

	price, ok := toFloat64(req.Price)
	if !ok || price == 0 {
		return fmt.Errorf("商品信息不完整")
	}

	// 验证组合商品必须至少包含一个子商品
	if isGroup == 0 && len(req.IncludedGoodsIDs) == 0 {
		return fmt.Errorf("组合商品必须至少包含一个子商品")
	}

	// 确保属性值ID数组不为nil
	attributeValueIDs := req.AttributeValueIDs
	if attributeValueIDs == nil {
		attributeValueIDs = []int{}
	}

	// 确保包含的商品ID数组不为nil
	includedGoodsIDs := req.IncludedGoodsIDs
	if includedGoodsIDs == nil {
		includedGoodsIDs = []int{}
	}

	err = s.goodsRepo.Update(
		id,
		name,
		brandID,
		classifyID,
		isGroup,
		price,
		attributeValueIDs,
		includedGoodsIDs,
	)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("商品不存在")
		}
		return fmt.Errorf("failed to update goods: %w", err)
	}

	return nil
}

// UpdateStatus 更新商品状态
func (s *GoodsService) UpdateStatus(id int, req *UpdateStatusRequest) error {
	// 转换并验证状态
	status, ok := toInt(req.Status)
	if !ok && req.Status != nil {
		return fmt.Errorf("状态不能为空")
	}

	if !ok {
		return fmt.Errorf("状态不能为空")
	}

	err := s.goodsRepo.UpdateStatus(id, status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("商品不存在")
		}
		return fmt.Errorf("failed to update goods status: %w", err)
	}

	return nil
}
