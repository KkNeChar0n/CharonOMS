package payment

// PaymentListFilter 收款列表查询条件
type PaymentListFilter struct {
	ID            *int
	StudentID     *int
	OrderID       *int
	Payer         *string
	PaymentMethod *int
	TradingDate   *string // 格式：YYYY-MM-DD
	Status        *int
	Page          int
	PageSize      int
}

// PaymentRepository 收款仓储接口
type PaymentRepository interface {
	// Create 创建收款记录
	Create(payment *PaymentCollection) error

	// GetByID 根据ID查询收款记录
	GetByID(id int) (*PaymentCollection, error)

	// List 查询收款列表
	List(filter PaymentListFilter) ([]*PaymentCollection, int64, error)

	// Update 更新收款记录
	Update(payment *PaymentCollection) error

	// Delete 删除收款记录
	Delete(id int) error

	// GetTotalPaidAmount 查询订单已收款总额
	// 包括常规收款（status=10或20）和淘宝收款（status=30）
	GetTotalPaidAmount(orderID int) (float64, error)

	// CountByOrderAndStatus 统计订单指定状态的收款数量
	CountByOrderAndStatus(orderID int, status int) (int64, error)
}
