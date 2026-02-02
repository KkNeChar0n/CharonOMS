package entity

import "time"

// 订单状态常量
const (
	OrderStatusDraft      = 10 // 草稿
	OrderStatusUnpaid     = 20 // 未支付
	OrderStatusPartialPaid = 30 // 部分支付
	OrderStatusPaid       = 40 // 已支付
	OrderStatusRefunding  = 50 // 退费中
	OrderStatusCancelled  = 99 // 已作废
)

// Order 订单实体
type Order struct {
	ID                  int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentID           int       `gorm:"column:student_id;not null" json:"student_id"`
	ExpectedPaymentTime *time.Time `gorm:"column:expected_payment_time" json:"expected_payment_time"`
	AmountReceivable    float64   `gorm:"column:amount_receivable;type:decimal(10,2);default:0" json:"amount_receivable"`
	AmountReceived      float64   `gorm:"column:amount_received;type:decimal(10,2);default:0" json:"amount_received"`
	DiscountAmount      float64   `gorm:"column:discount_amount;type:decimal(10,2);default:0" json:"discount_amount"`
	Status              int       `gorm:"column:status;default:10" json:"status"`
	CreateTime          time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// CanEdit 判断订单是否可以编辑
func (o *Order) CanEdit() bool {
	return o.Status == OrderStatusDraft
}

// CanSubmit 判断订单是否可以提交
func (o *Order) CanSubmit() bool {
	return o.Status == OrderStatusDraft
}

// CanCancel 判断订单是否可以作废
func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusDraft
}

// ValidateAmounts 验证订单金额的合理性
func (o *Order) ValidateAmounts() bool {
	// 应收金额应该大于等于0
	if o.AmountReceivable < 0 {
		return false
	}
	// 实收金额应该大于等于0
	if o.AmountReceived < 0 {
		return false
	}
	// 优惠金额应该大于等于0
	if o.DiscountAmount < 0 {
		return false
	}
	// 实收金额不应该大于应收金额
	if o.AmountReceived > o.AmountReceivable {
		return false
	}
	return true
}
