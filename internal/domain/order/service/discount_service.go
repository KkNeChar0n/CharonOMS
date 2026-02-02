package service

import (
	"context"

	"gorm.io/gorm"
)

// DiscountService 优惠计算服务
type DiscountService struct {
	db *gorm.DB
}

// NewDiscountService 创建优惠计算服务
func NewDiscountService(db *gorm.DB) *DiscountService {
	return &DiscountService{db: db}
}

// GoodsForDiscount 用于优惠计算的商品信息
type GoodsForDiscount struct {
	GoodsID int
	Price   float64
}

// CalculateDiscount 计算订单优惠
func (s *DiscountService) CalculateDiscount(ctx context.Context, goodsList []GoodsForDiscount, activityIDs []int) (totalDiscount float64, childDiscounts map[int]float64, err error) {
	childDiscounts = make(map[int]float64)

	// 初始化每个商品的优惠为0
	for _, goods := range goodsList {
		childDiscounts[goods.GoodsID] = 0
	}

	// 如果没有商品或活动，返回零优惠
	if len(goodsList) == 0 || len(activityIDs) == 0 {
		return 0, childDiscounts, nil
	}

	// 遍历每个活动计算优惠
	for _, activityID := range activityIDs {
		discount, childDisc, err := s.calculateActivityDiscount(ctx, goodsList, activityID)
		if err != nil {
			continue // 忽略单个活动的错误，继续处理其他活动
		}

		totalDiscount += discount

		// 累加子订单优惠
		for goodsID, amount := range childDisc {
			childDiscounts[goodsID] += amount
		}
	}

	// 四舍五入到两位小数
	totalDiscount = roundToTwoDecimal(totalDiscount)
	for goodsID := range childDiscounts {
		childDiscounts[goodsID] = roundToTwoDecimal(childDiscounts[goodsID])
	}

	return totalDiscount, childDiscounts, nil
}

// calculateActivityDiscount 计算单个活动的优惠
func (s *DiscountService) calculateActivityDiscount(ctx context.Context, goodsList []GoodsForDiscount, activityID int) (float64, map[int]float64, error) {
	childDiscounts := make(map[int]float64)

	// 1. 查询活动信息
	var activity struct {
		ID         int
		TemplateID int
		Type       int
		SelectType int
	}

	err := s.db.WithContext(ctx).
		Table("activity a").
		Select("a.id, a.template_id, t.type, t.select_type").
		Joins("JOIN activity_template t ON a.template_id = t.id").
		Where("a.id = ?", activityID).
		First(&activity).Error

	if err != nil {
		return 0, childDiscounts, err
	}

	// 2. 仅处理满折类型（type=2）
	if activity.Type != 2 {
		return 0, childDiscounts, nil
	}

	// 3. 查询活动折扣规则（按threshold_amount降序）
	var details []struct {
		ThresholdAmount float64
		DiscountValue   float64
	}

	err = s.db.WithContext(ctx).
		Table("activity_detail").
		Select("threshold_amount, discount_value").
		Where("activity_id = ?", activityID).
		Order("threshold_amount DESC").
		Find(&details).Error

	if err != nil || len(details) == 0 {
		return 0, childDiscounts, err
	}

	// 4. 查询活动模板关联的商品/类型
	var templateGoods []struct {
		GoodsID    *int
		ClassifyID *int
	}

	err = s.db.WithContext(ctx).
		Table("activity_template_goods").
		Select("goods_id, classify_id").
		Where("template_id = ?", activity.TemplateID).
		Find(&templateGoods).Error

	if err != nil || len(templateGoods) == 0 {
		return 0, childDiscounts, err
	}

	// 5. 筛选参与活动的商品
	eligibleGoods, err := s.filterEligibleGoods(ctx, goodsList, templateGoods, activity.SelectType)
	if err != nil || len(eligibleGoods) == 0 {
		return 0, childDiscounts, err
	}

	// 6. 匹配最大满足的折扣档位（按商品数量匹配）
	eligibleCount := float64(len(eligibleGoods))
	var matchedDiscount float64 = 0

	for _, detail := range details {
		if eligibleCount >= detail.ThresholdAmount {
			matchedDiscount = detail.DiscountValue
			break
		}
	}

	if matchedDiscount == 0 {
		return 0, childDiscounts, nil
	}

	// 7. 计算总优惠：(1 - 折扣/100) × 参与商品标准售价之和
	// discount_value=90表示9折(付90%)，discount_value=80表示8折(付80%)
	var eligiblePriceSum float64
	for _, goods := range eligibleGoods {
		eligiblePriceSum += goods.Price
	}

	discountRate := 1 - matchedDiscount/100 // 80 -> 1 - 0.8 = 0.2 (优惠20%)
	activityDiscount := discountRate * eligiblePriceSum

	// 8. 按比例分摊优惠到各商品
	if activityDiscount > 0 && eligiblePriceSum > 0 {
		allocated := 0.0
		for i, goods := range eligibleGoods {
			var childDiscount float64
			if i == len(eligibleGoods)-1 {
				// 最后一个商品用减法避免精度误差
				childDiscount = activityDiscount - allocated
			} else {
				ratio := goods.Price / eligiblePriceSum
				childDiscount = roundToTwoDecimal(activityDiscount * ratio)
				allocated += childDiscount
			}
			childDiscounts[goods.GoodsID] = childDiscount
		}
	}

	return activityDiscount, childDiscounts, nil
}

// filterEligibleGoods 筛选参与活动的商品
func (s *DiscountService) filterEligibleGoods(ctx context.Context, goodsList []GoodsForDiscount, templateGoods []struct {
	GoodsID    *int
	ClassifyID *int
}, selectType int) ([]GoodsForDiscount, error) {
	var eligibleGoods []GoodsForDiscount

	for _, goods := range goodsList {
		isEligible := false

		if selectType == 2 {
			// 按商品：检查商品ID是否在模板关联的商品中
			for _, tg := range templateGoods {
				if tg.GoodsID != nil && *tg.GoodsID == goods.GoodsID {
					isEligible = true
					break
				}
			}
		} else {
			// 按分类：查询商品的classifyid并匹配
			var classifyID int
			err := s.db.WithContext(ctx).
				Table("goods").
				Select("classifyid").
				Where("id = ?", goods.GoodsID).
				Scan(&classifyID).Error

			if err == nil {
				for _, tg := range templateGoods {
					if tg.ClassifyID != nil && *tg.ClassifyID == classifyID {
						isEligible = true
						break
					}
				}
			}
		}

		if isEligible {
			eligibleGoods = append(eligibleGoods, goods)
		}
	}

	return eligibleGoods, nil
}
