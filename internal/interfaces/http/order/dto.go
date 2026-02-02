package order

import (
	"encoding/json"
	"strings"
	"time"
)

// CustomTime 自定义时间类型，支持多种格式解析
type CustomTime struct {
	time.Time
}

// UnmarshalJSON 自定义JSON解析，支持多种时间格式
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		ct.Time = time.Time{}
		return nil
	}

	// 尝试多种时间格式
	formats := []string{
		"2006-01-02T15:04:05Z07:00",     // RFC3339
		"2006-01-02T15:04:05.999Z07:00", // RFC3339 with milliseconds
		"2006-01-02T15:04:05Z",          // ISO8601 with timezone
		"2006-01-02T15:04:05",           // ISO8601 without timezone
		"2006-01-02T15:04",              // ISO8601 without seconds (前端使用)
		"2006-01-02 15:04:05",           // MySQL datetime
		"2006-01-02 15:04",              // datetime without seconds
		"2006-01-02",                    // Date only
		time.RFC3339,
		time.RFC3339Nano,
	}

	var err error
	for _, format := range formats {
		ct.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}

// MarshalJSON 自定义JSON序列化
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ct.Time.Format("2006-01-02 15:04:05"))
}

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
	ExpectedPaymentTime *CustomTime          `json:"expected_payment_time"`
	ActivityIDs         []int                `json:"activity_ids"`
	DiscountAmount      float64              `json:"discount_amount"`
	ChildDiscounts      map[int]float64      `json:"child_discounts"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	GoodsList           []GoodsItemRequest   `json:"goods_list"`
	ExpectedPaymentTime *CustomTime          `json:"expected_payment_time"`
	ActivityIDs         []int                `json:"activity_ids"`
	DiscountAmount      float64              `json:"discount_amount"`
	ChildDiscounts      map[int]float64      `json:"child_discounts"`
}

// CalculateDiscountRequest 优惠计算请求
type CalculateDiscountRequest struct {
	GoodsList   []GoodsItemRequest `json:"goods_list"`
	ActivityIDs []int              `json:"activity_ids"`
}

// CalculateDiscountResponse 优惠计算响应
type CalculateDiscountResponse struct {
	TotalDiscount  float64         `json:"total_discount"`
	ChildDiscounts map[int]float64 `json:"child_discounts"`
}
