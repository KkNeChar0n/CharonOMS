package order

import "time"

// OrderDO orders表数据对象
type OrderDO struct {
	ID                  int        `gorm:"column:id;primaryKey;autoIncrement"`
	StudentID           int        `gorm:"column:student_id;not null"`
	ExpectedPaymentTime *time.Time `gorm:"column:expected_payment_time"`
	AmountReceivable    float64    `gorm:"column:amount_receivable;type:decimal(10,2);default:0"`
	AmountReceived      float64    `gorm:"column:amount_received;type:decimal(10,2);default:0"`
	DiscountAmount      float64    `gorm:"column:discount_amount;type:decimal(10,2);default:0"`
	Status              int        `gorm:"column:status;default:10"`
	CreateTime          time.Time  `gorm:"column:create_time;autoCreateTime"`
}

// TableName 指定表名
func (OrderDO) TableName() string {
	return "orders"
}

// ChildOrderDO childorders表数据对象
type ChildOrderDO struct {
	ID               int       `gorm:"column:id;primaryKey;autoIncrement"`
	ParentsID        int       `gorm:"column:parentsid;not null"`
	GoodsID          int       `gorm:"column:goodsid;not null"`
	AmountReceivable float64   `gorm:"column:amount_receivable;type:decimal(10,2);default:0"`
	AmountReceived   float64   `gorm:"column:amount_received;type:decimal(10,2);default:0"`
	DiscountAmount   float64   `gorm:"column:discount_amount;type:decimal(10,2);default:0"`
	Status           int       `gorm:"column:status;default:0"`
	CreateTime       time.Time `gorm:"column:create_time;autoCreateTime"`
}

// TableName 指定表名
func (ChildOrderDO) TableName() string {
	return "childorders"
}

// OrdersActivityDO orders_activity表数据对象
type OrdersActivityDO struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement"`
	OrdersID   int       `gorm:"column:orders_id;not null"`
	ActivityID int       `gorm:"column:activity_id;not null"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
}

// TableName 指定表名
func (OrdersActivityDO) TableName() string {
	return "orders_activity"
}
