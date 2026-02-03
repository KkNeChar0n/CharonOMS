package separate

import "time"

// 收款类型常量
const (
	PaymentTypeRegular = 0 // 常规收款
	PaymentTypeTaobao  = 1 // 淘宝收款
)

// 分账类型常量
const (
	SeparateTypeSale   = 0 // 售卖
	SeparateTypeRevert = 1 // 冲回
	SeparateTypeRefund = 2 // 退费
)

// SeparateAccount 分账明细实体
type SeparateAccount struct {
	ID             int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UID            int       `gorm:"column:uid;not null" json:"uid"`
	OrdersID       int       `gorm:"column:orders_id;not null" json:"orders_id"`
	ChildOrdersID  int       `gorm:"column:childorders_id;not null" json:"childorders_id"`
	PaymentID      int       `gorm:"column:payment_id;not null" json:"payment_id"`
	PaymentType    int       `gorm:"column:payment_type;not null" json:"payment_type"`
	GoodsID        int       `gorm:"column:goods_id;not null" json:"goods_id"`
	GoodsName      string    `gorm:"column:goods_name;type:varchar(100);not null" json:"goods_name"`
	SeparateAmount float64   `gorm:"column:separate_amount;type:decimal(10,2);not null" json:"separate_amount"`
	Type           int       `gorm:"column:type;default:0" json:"type"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (SeparateAccount) TableName() string {
	return "separate_account"
}
