package payment

import "time"

// 付款场景常量
const (
	PaymentScenarioOnline  = 0 // 线上
	PaymentScenarioOffline = 1 // 线下
)

// 付款方式常量
const (
	PaymentMethodWechat    = 0 // 微信
	PaymentMethodAlipay    = 1 // 支付宝
	PaymentMethodYouli     = 2 // 优利支付
	PaymentMethodLingling  = 3 // 零零购支付
	PaymentMethodPublic    = 9 // 对公转账
)

// 收款主体常量
const (
	PayeeEntityBeijing = 0 // 北京
	PayeeEntityXian    = 1 // 西安
)

// 收款状态常量
const (
	PaymentStatusWaitPay    = 0  // 待支付
	PaymentStatusUnverified = 10 // 未核验
	PaymentStatusPaid       = 20 // 已支付
)

// PaymentCollection 收款记录实体
type PaymentCollection struct {
	ID              int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID         int        `gorm:"column:order_id;not null" json:"order_id"`
	StudentID       int        `gorm:"column:student_id;not null" json:"student_id"`
	PaymentScenario int        `gorm:"column:payment_scenario;not null" json:"payment_scenario"`
	PaymentMethod   int        `gorm:"column:payment_method;not null" json:"payment_method"`
	PaymentAmount   float64    `gorm:"column:payment_amount;type:decimal(10,2);not null" json:"payment_amount"`
	Payer           string     `gorm:"column:payer;type:varchar(100)" json:"payer"`
	PayeeEntity     int        `gorm:"column:payee_entity;not null" json:"payee_entity"`
	TradingHours    *time.Time `gorm:"column:trading_hours" json:"trading_hours"`
	ArrivalTime     *time.Time `gorm:"column:arrival_time" json:"arrival_time"`
	MerchantOrder   string     `gorm:"column:merchant_order;type:varchar(100)" json:"merchant_order"`
	Status          int        `gorm:"column:status;default:10" json:"status"`
	CreateTime      time.Time  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (PaymentCollection) TableName() string {
	return "payment_collection"
}

// CanConfirm 判断是否可以确认到账
func (p *PaymentCollection) CanConfirm() bool {
	return p.Status == PaymentStatusUnverified
}

// CanDelete 判断是否可以删除
func (p *PaymentCollection) CanDelete() bool {
	return p.Status == PaymentStatusUnverified
}

// Confirm 确认到账
func (p *PaymentCollection) Confirm() {
	p.Status = PaymentStatusPaid
	now := time.Now()
	p.ArrivalTime = &now
}
