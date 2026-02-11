package refund

import "time"

// RefundOrder 退费订单实体
type RefundOrder struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID      int       `json:"order_id" gorm:"not null;comment:关联的主订单ID"`
	StudentID    int       `json:"uid" gorm:"column:student_id;not null;comment:学生ID"`
	RefundAmount float64   `json:"refund_amount" gorm:"type:decimal(10,2);not null;comment:总退费金额"`
	Submitter    string    `json:"submitter" gorm:"type:varchar(100);comment:提交人用户名"`
	SubmitTime   time.Time `json:"submit_time" gorm:"comment:提交时间"`
	Status       int       `json:"status" gorm:"type:tinyint;default:0;comment:状态：0-待审批、10-已通过、20-已驳回"`
	CreateTime   time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

// TableName 表名
func (RefundOrder) TableName() string {
	return "refund_order"
}

// RefundOrderItem 退费子订单明细实体
type RefundOrderItem struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RefundOrderID int       `json:"refund_order_id" gorm:"not null;comment:所属退费订单ID"`
	ChildOrderID  int       `json:"childorder_id" gorm:"column:childorder_id;not null;comment:子订单ID"`
	GoodsID       int       `json:"goods_id" gorm:"not null;comment:商品ID"`
	GoodsName     string    `json:"goods_name" gorm:"type:varchar(100);comment:商品名称"`
	RefundAmount  float64   `json:"refund_amount" gorm:"type:decimal(10,2);not null;comment:退费金额"`
	Status        int       `json:"status" gorm:"type:tinyint;default:0;comment:状态：0-待审批、10-已通过、20-已驳回"`
	CreateTime    time.Time `json:"create_time" gorm:"autoCreateTime"`
	// 关联字段（从refund_order表获取，通过JOIN查询填充）
	UID     int `json:"uid" gorm:"column:uid;->"`
	OrderID int `json:"order_id" gorm:"column:order_id;->"`
}

// TableName 表名
func (RefundOrderItem) TableName() string {
	return "refund_order_item"
}

// RefundPayment 退费收款分配实体
type RefundPayment struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RefundOrderID int       `json:"refund_order_id" gorm:"not null;comment:所属退费订单ID"`
	PaymentID     int       `json:"payment_id" gorm:"not null;comment:收款记录ID"`
	PaymentType   int       `json:"payment_type" gorm:"type:tinyint;not null;comment:收款类型：0-常规收款、1-淘宝收款"`
	RefundAmount  float64   `json:"refund_amount" gorm:"type:decimal(10,2);not null;comment:从该收款退费的金额"`
	CreateTime    time.Time `json:"create_time" gorm:"autoCreateTime"`
	// 关联字段（从refund_order表获取，通过JOIN查询填充）
	UID     int `json:"uid" gorm:"column:uid;->"`
	OrderID int `json:"order_id" gorm:"column:order_id;->"`
	Status  int `json:"status" gorm:"column:status;->"`
}

// TableName 表名
func (RefundPayment) TableName() string {
	return "refund_payment"
}

// RefundTaobaoSupplement 淘宝退费补充信息实体
type RefundTaobaoSupplement struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RefundOrderID  int       `json:"refund_order_id" gorm:"not null;comment:所属退费订单ID"`
	StudentID      int       `json:"uid" gorm:"column:student_id;not null;comment:学生ID"`
	AlipayAccount  string    `json:"alipay_account" gorm:"type:varchar(100);comment:支付宝账号"`
	AlipayName     string    `json:"alipay_name" gorm:"type:varchar(100);comment:支付宝账户名"`
	RefundAmount   float64   `json:"refund_amount" gorm:"type:decimal(10,2);not null;comment:淘宝退费金额"`
	Status         int       `json:"status" gorm:"type:tinyint;default:0;comment:状态：0-待审批、10-已通过、20-已驳回"`
	CreateTime     time.Time `json:"create_time" gorm:"autoCreateTime"`
}

// TableName 表名
func (RefundTaobaoSupplement) TableName() string {
	return "refund_taobao_supplement"
}

// RefundRegularSupplement 常规退费补充信息实体
type RefundRegularSupplement struct {
	ID                   int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RefundOrderID        int       `json:"refund_order_id" gorm:"not null;comment:所属退费订单ID"`
	StudentID            int       `json:"uid" gorm:"column:student_id;not null;comment:学生ID"`
	PayeeEntity          *int      `json:"payee_entity" gorm:"type:tinyint;comment:收款实体"`
	IsCorporateTransfer  *bool     `json:"is_corporate_transfer" gorm:"comment:是否企业转账"`
	Payer                string    `json:"payer" gorm:"type:varchar(100);comment:付款人名称"`
	BankAccount          string    `json:"bank_account" gorm:"type:varchar(100);comment:银行账户"`
	PayerReadonly        *bool     `json:"payer_readonly" gorm:"comment:付款人是否只读"`
	RefundAmount         float64   `json:"refund_amount" gorm:"type:decimal(10,2);not null;comment:常规退费金额"`
	Status               int       `json:"status" gorm:"type:tinyint;default:0;comment:状态：0-待审批、10-已通过、20-已驳回"`
	CreateTime           time.Time `json:"create_time" gorm:"autoCreateTime"`
}

// TableName 表名
func (RefundRegularSupplement) TableName() string {
	return "refund_regular_supplement"
}

// 状态常量
const (
	RefundStatusPending  = 0  // 待审批
	RefundStatusApproved = 10 // 已通过
	RefundStatusRejected = 20 // 已驳回
)

// 收款类型常量
const (
	PaymentTypeRegular = 0 // 常规收款
	PaymentTypeTaobao  = 1 // 淘宝收款
)
