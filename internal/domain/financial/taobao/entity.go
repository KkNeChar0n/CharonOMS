package taobao

import "time"

// TaobaoPayment 淘宝收款实体
type TaobaoPayment struct {
	ID              int       `json:"id"`
	OrderID         *int      `json:"order_id"`         // 关联订单ID
	StudentID       *int      `json:"student_id"`       // 学生ID
	Payer           *string   `json:"payer"`            // 付款方
	ZhifubaoAccount *string   `json:"zhifubao_account"` // 支付宝账号
	PaymentAmount   float64   `json:"payment_amount"`   // 金额
	OrderTime       *time.Time `json:"order_time"`      // 下单时间
	ArrivalTime     *time.Time `json:"arrival_time"`    // 到账时间
	MerchantOrder   *string   `json:"merchant_order"`   // 商户订单号
	Status          int       `json:"status"`           // 状态：0-已下单，10-待认领，20-已认领，30-已到账，40-已退单
	Claimer         *int      `json:"claimer"`          // 认领人ID
	CreateTime      time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime      time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (TaobaoPayment) TableName() string {
	return "taobao_payment"
}

// 淘宝收款状态常量
const (
	TaobaoPaymentStatusOrdered    = 0  // 已下单
	TaobaoPaymentStatusUnclaimed  = 10 // 待认领
	TaobaoPaymentStatusClaimed    = 20 // 已认领
	TaobaoPaymentStatusArrived    = 30 // 已到账
	TaobaoPaymentStatusRefunded   = 40 // 已退单
)
