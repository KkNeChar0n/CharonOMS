package service

import (
	"math"
)

// OrderService 订单领域服务
type OrderService struct{}

// NewOrderService 创建订单领域服务
func NewOrderService() *OrderService {
	return &OrderService{}
}

// GoodsItem 商品项
type GoodsItem struct {
	GoodsID    int
	TotalPrice float64 // 商品总价（组合商品为子商品价格之和）
	Price      float64 // 标准售价
}

// CalculateOrderAmounts 计算订单应收和实收金额
func (s *OrderService) CalculateOrderAmounts(goodsList []GoodsItem, discountAmount float64) (amountReceivable float64, amountReceived float64) {
	// 应收金额 = 所有商品的总价之和（total_price）
	for _, goods := range goodsList {
		amountReceivable += goods.TotalPrice
	}

	// 实收金额 = 应收金额 - 优惠金额
	amountReceived = amountReceivable - discountAmount

	// 保留两位小数
	amountReceivable = roundToTwoDecimal(amountReceivable)
	amountReceived = roundToTwoDecimal(amountReceived)

	return
}

// AllocateChildDiscounts 将订单优惠分摊到子订单
// childDiscounts: 前端已计算好的每个商品的优惠金额
func (s *OrderService) AllocateChildDiscounts(goodsList []GoodsItem, childDiscounts map[int]float64) map[int]float64 {
	result := make(map[int]float64)

	for _, goods := range goodsList {
		// 使用前端传入的优惠分摊，如果没有则为0
		discount := childDiscounts[goods.GoodsID]
		result[goods.GoodsID] = roundToTwoDecimal(discount)
	}

	return result
}

// CalculateChildAmounts 计算子订单的应收和实收金额
// totalPrice: 商品总价（单商品=price，组合商品=子商品价格之和）
// price: 商品标准售价（goods表的price字段，用于显示但不参与计算）
// discountAmount: 优惠金额
func (s *OrderService) CalculateChildAmounts(totalPrice float64, price float64, discountAmount float64) (amountReceivable float64, amountReceived float64) {
	// 应收金额 = 商品总价（total_price）
	amountReceivable = roundToTwoDecimal(totalPrice)
	// 实收金额 = 总价 - 优惠金额
	amountReceived = roundToTwoDecimal(totalPrice - discountAmount)
	return
}

// roundToTwoDecimal 四舍五入到两位小数
func roundToTwoDecimal(value float64) float64 {
	return math.Round(value*100) / 100
}
