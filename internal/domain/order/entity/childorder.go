package entity

import "time"

// 子订单状态常量
const (
	ChildOrderStatusInit        = 0  // 初始
	ChildOrderStatusUnpaid      = 10 // 未支付
	ChildOrderStatusPartialPaid = 20 // 部分支付
	ChildOrderStatusPaid        = 30 // 已支付
	ChildOrderStatusCancelled   = 99 // 已作废
)

// ChildOrder 子订单实体
type ChildOrder struct {
	ID               int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentsID        int       `gorm:"column:parentsid;not null" json:"parentsid"`
	GoodsID          int       `gorm:"column:goodsid;not null" json:"goodsid"`
	AmountReceivable float64   `gorm:"column:amount_receivable;type:decimal(10,2);default:0" json:"amount_receivable"`
	AmountReceived   float64   `gorm:"column:amount_received;type:decimal(10,2);default:0" json:"amount_received"`
	DiscountAmount   float64   `gorm:"column:discount_amount;type:decimal(10,2);default:0" json:"discount_amount"`
	Status           int       `gorm:"column:status;default:0" json:"status"`
	CreateTime       time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (ChildOrder) TableName() string {
	return "childorders"
}

// ValidateAmounts 验证子订单金额的合理性
func (c *ChildOrder) ValidateAmounts() bool {
	// 应收金额应该大于等于0
	if c.AmountReceivable < 0 {
		return false
	}
	// 实收金额应该大于等于0
	if c.AmountReceived < 0 {
		return false
	}
	// 优惠金额应该大于等于0
	if c.DiscountAmount < 0 {
		return false
	}
	// 实收金额不应该大于应收金额
	if c.AmountReceived > c.AmountReceivable {
		return false
	}
	return true
}
