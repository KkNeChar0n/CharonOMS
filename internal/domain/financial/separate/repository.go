package separate

// SeparateListFilter 分账明细列表查询条件
type SeparateListFilter struct {
	ID            *int
	UID           *int
	OrdersID      *int
	ChildOrdersID *int
	GoodsID       *int
	PaymentID     *int
	PaymentType   *int
	Type          *int
	Page          int
	PageSize      int
}

// SeparateAccountRepository 分账明细仓储接口
type SeparateAccountRepository interface {
	// Create 创建分账明细
	Create(account *SeparateAccount) error

	// BatchCreate 批量创建分账明细（在事务中）
	BatchCreate(accounts []*SeparateAccount) error

	// GetByID 根据ID查询分账明细
	GetByID(id int) (*SeparateAccount, error)

	// List 查询分账明细列表
	List(filter SeparateListFilter) ([]*SeparateAccount, int64, error)

	// ExistsByPaymentAndOrder 检查指定收款和订单是否已生成分账
	ExistsByPaymentAndOrder(paymentID int, orderID int, paymentType int) (bool, error)

	// GetChildOrderTotalSeparate 查询子订单的总分账金额
	// 只统计售卖类型（type=0）
	GetChildOrderTotalSeparate(childOrderID int) (float64, error)

	// GetChildOrderAllocatedAmount 查询子订单已分配金额
	GetChildOrderAllocatedAmount(childOrderID int) (float64, error)
}
