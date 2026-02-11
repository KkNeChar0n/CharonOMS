package refund

// RefundRepository 退费订单仓储接口
type RefundRepository interface {
	// CreateRefundOrder 创建退费订单
	CreateRefundOrder(refundOrder *RefundOrder) error

	// GetRefundOrderByID 根据ID获取退费订单
	GetRefundOrderByID(id int) (*RefundOrder, error)

	// ListRefundOrders 获取退费订单列表
	ListRefundOrders(filters map[string]interface{}) ([]*RefundOrder, error)

	// UpdateRefundOrder 更新退费订单
	UpdateRefundOrder(refundOrder *RefundOrder) error

	// CreateRefundOrderItem 创建退费子订单明细
	CreateRefundOrderItem(item *RefundOrderItem) error

	// ListRefundOrderItems 获取退费子订单明细列表
	ListRefundOrderItems(filters map[string]interface{}) ([]*RefundOrderItem, error)

	// GetRefundOrderItemsByRefundID 根据退费订单ID获取所有明细
	GetRefundOrderItemsByRefundID(refundOrderID int) ([]*RefundOrderItem, error)

	// CreateRefundPayment 创建退费收款分配
	CreateRefundPayment(payment *RefundPayment) error

	// ListRefundPayments 获取退费收款分配列表
	ListRefundPayments(filters map[string]interface{}) ([]*RefundPayment, error)

	// GetRefundPaymentsByRefundID 根据退费订单ID获取所有收款分配
	GetRefundPaymentsByRefundID(refundOrderID int) ([]*RefundPayment, error)

	// CreateTaobaoSupplement 创建淘宝退费补充信息
	CreateTaobaoSupplement(supplement *RefundTaobaoSupplement) error

	// ListTaobaoSupplements 获取淘宝退费补充信息列表
	ListTaobaoSupplements(filters map[string]interface{}) ([]*RefundTaobaoSupplement, error)

	// GetTaobaoSupplementByRefundID 根据退费订单ID获取淘宝补充信息
	GetTaobaoSupplementByRefundID(refundOrderID int) (*RefundTaobaoSupplement, error)

	// CreateRegularSupplement 创建常规退费补充信息
	CreateRegularSupplement(supplement *RefundRegularSupplement) error

	// ListRegularSupplements 获取常规退费补充信息列表
	ListRegularSupplements(filters map[string]interface{}) ([]*RefundRegularSupplement, error)

	// GetRegularSupplementsByRefundID 根据退费订单ID获取所有常规补充信息
	GetRegularSupplementsByRefundID(refundOrderID int) ([]*RefundRegularSupplement, error)
}
