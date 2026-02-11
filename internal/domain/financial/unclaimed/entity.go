package unclaimed

import "time"

// Unclaimed 常规待认领款项实体
type Unclaimed struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentMethod int        `json:"payment_method" gorm:"type:tinyint;comment:付款方式：0-微信、1-支付宝、2-优利支付、3-零零购支付、9-对公转账"`
	PaymentAmount float64    `json:"payment_amount" gorm:"type:decimal(10,2);comment:付款金额"`
	Payer         *string    `json:"payer" gorm:"type:varchar(100);comment:付款方"`
	PayeeEntity   int        `json:"payee_entity" gorm:"type:tinyint;comment:收款主体：0-北京、1-西安"`
	MerchantOrder *string    `json:"merchant_order" gorm:"type:varchar(100);comment:商户订单号"`
	ArrivalTime   *time.Time `json:"arrival_time" gorm:"comment:到账时间"`
	Claimer       *int       `json:"claimer" gorm:"comment:认领人ID"`
	PaymentID     *int       `json:"payment_id" gorm:"comment:关联的payment_collection记录ID"`
	Status        int        `json:"status" gorm:"type:tinyint;default:0;comment:状态：0-待认领、1-已认领"`
	CreateTime    time.Time  `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime    time.Time  `json:"update_time" gorm:"autoUpdateTime"`
}

// 表名
func (Unclaimed) TableName() string {
	return "unclaimed"
}

// 状态常量
const (
	UnclaimedStatusPending = 0 // 待认领
	UnclaimedStatusClaimed = 1 // 已认领
)

// 付款方式常量
const (
	PaymentMethodWechat       = 0 // 微信
	PaymentMethodAlipay       = 1 // 支付宝
	PaymentMethodYouli        = 2 // 优利支付
	PaymentMethodLinglinggou  = 3 // 零零购支付
	PaymentMethodBankTransfer = 9 // 对公转账
)

// 收款主体常量
const (
	PayeeEntityBeijing = 0 // 北京
	PayeeEntityXian    = 1 // 西安
)
