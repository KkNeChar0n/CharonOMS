package order

import "time"

// GoodsItemRequest 商品项请求
type GoodsItemRequest struct {
	GoodsID    int     `json:"goods_id"`
	TotalPrice float64 `json:"total_price"`
	Price      float64 `json:"price"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	StudentID           int                  `json:"student_id"`
	GoodsList           []GoodsItemRequest   `json:"goods_list"`
	ExpectedPaymentTime *time.Time           `json:"expected_payment_time"`
	ActivityIDs         []int                `json:"activity_ids"`
	DiscountAmount      float64              `json:"discount_amount"`
	ChildDiscounts      map[int]float64      `json:"child_discounts"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	GoodsList           []GoodsItemRequest   `json:"goods_list"`
	ExpectedPaymentTime *time.Time           `json:"expected_payment_time"`
	ActivityIDs         []int                `json:"activity_ids"`
	DiscountAmount      float64              `json:"discount_amount"`
	ChildDiscounts      map[int]float64      `json:"child_discounts"`
}
