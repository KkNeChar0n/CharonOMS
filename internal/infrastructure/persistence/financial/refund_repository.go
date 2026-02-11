package financial

import (
	"charonoms/internal/domain/financial/refund"

	"gorm.io/gorm"
)

type refundRepository struct {
	db *gorm.DB
}

// NewRefundRepository 创建退费订单仓储实例
func NewRefundRepository(db *gorm.DB) refund.RefundRepository {
	return &refundRepository{db: db}
}

// CreateRefundOrder 创建退费订单
func (r *refundRepository) CreateRefundOrder(refundOrder *refund.RefundOrder) error {
	return r.db.Create(refundOrder).Error
}

// GetRefundOrderByID 根据ID获取退费订单
func (r *refundRepository) GetRefundOrderByID(id int) (*refund.RefundOrder, error) {
	var order refund.RefundOrder
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListRefundOrders 获取退费订单列表
func (r *refundRepository) ListRefundOrders(filters map[string]interface{}) ([]*refund.RefundOrder, error) {
	var list []*refund.RefundOrder
	query := r.db.Model(&refund.RefundOrder{})

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("id = ?", id)
	}
	if uid, ok := filters["uid"]; ok && uid != "" {
		query = query.Where("uid = ?", uid)
	}
	if orderID, ok := filters["order_id"]; ok && orderID != "" {
		query = query.Where("order_id = ?", orderID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// 按创建时间倒序
	query = query.Order("create_time DESC")

	err := query.Find(&list).Error
	return list, err
}

// UpdateRefundOrder 更新退费订单
func (r *refundRepository) UpdateRefundOrder(refundOrder *refund.RefundOrder) error {
	return r.db.Save(refundOrder).Error
}

// CreateRefundOrderItem 创建退费子订单明细
func (r *refundRepository) CreateRefundOrderItem(item *refund.RefundOrderItem) error {
	return r.db.Create(item).Error
}

// ListRefundOrderItems 获取退费子订单明细列表
func (r *refundRepository) ListRefundOrderItems(filters map[string]interface{}) ([]*refund.RefundOrderItem, error) {
	var list []*refund.RefundOrderItem

	// 使用JOIN查询获取关联的uid和order_id
	query := r.db.Table("refund_order_item").
		Select("refund_order_item.*, refund_order.student_id as uid, refund_order.order_id").
		Joins("LEFT JOIN refund_order ON refund_order_item.refund_order_id = refund_order.id")

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("refund_order_item.id = ?", id)
	}
	if studentID, ok := filters["student_id"]; ok && studentID != "" {
		query = query.Where("refund_order.student_id = ?", studentID)
	}
	if orderID, ok := filters["order_id"]; ok && orderID != "" {
		query = query.Where("refund_order.order_id = ?", orderID)
	}
	if refundOrderID, ok := filters["refund_order_id"]; ok && refundOrderID != "" {
		query = query.Where("refund_order_item.refund_order_id = ?", refundOrderID)
	}
	if childOrderID, ok := filters["childorder_id"]; ok && childOrderID != "" {
		query = query.Where("refund_order_item.childorder_id = ?", childOrderID)
	}
	if goodsID, ok := filters["goods_id"]; ok && goodsID != "" {
		query = query.Where("refund_order_item.goods_id = ?", goodsID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("refund_order_item.status = ?", status)
	}

	// 按创建时间倒序
	query = query.Order("refund_order_item.create_time DESC")

	err := query.Scan(&list).Error
	return list, err
}

// GetRefundOrderItemsByRefundID 根据退费订单ID获取所有明细
func (r *refundRepository) GetRefundOrderItemsByRefundID(refundOrderID int) ([]*refund.RefundOrderItem, error) {
	var items []*refund.RefundOrderItem
	err := r.db.Where("refund_order_id = ?", refundOrderID).Find(&items).Error
	return items, err
}

// CreateRefundPayment 创建退费收款分配
func (r *refundRepository) CreateRefundPayment(payment *refund.RefundPayment) error {
	return r.db.Create(payment).Error
}

// ListRefundPayments 获取退费收款分配列表
func (r *refundRepository) ListRefundPayments(filters map[string]interface{}) ([]*refund.RefundPayment, error) {
	var list []*refund.RefundPayment

	// 使用JOIN查询获取关联的uid、order_id和status
	query := r.db.Table("refund_payment").
		Select("refund_payment.*, refund_order.student_id as uid, refund_order.order_id, refund_order.status").
		Joins("LEFT JOIN refund_order ON refund_payment.refund_order_id = refund_order.id")

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("refund_payment.id = ?", id)
	}
	if studentID, ok := filters["student_id"]; ok && studentID != "" {
		query = query.Where("refund_order.student_id = ?", studentID)
	}
	if orderID, ok := filters["order_id"]; ok && orderID != "" {
		query = query.Where("refund_order.order_id = ?", orderID)
	}
	if refundOrderID, ok := filters["refund_order_id"]; ok && refundOrderID != "" {
		query = query.Where("refund_payment.refund_order_id = ?", refundOrderID)
	}
	if paymentID, ok := filters["payment_id"]; ok && paymentID != "" {
		query = query.Where("refund_payment.payment_id = ?", paymentID)
	}
	if paymentType, ok := filters["payment_type"]; ok && paymentType != "" {
		query = query.Where("refund_payment.payment_type = ?", paymentType)
	}

	// 按创建时间倒序
	query = query.Order("refund_payment.create_time DESC")

	err := query.Scan(&list).Error
	return list, err
}

// GetRefundPaymentsByRefundID 根据退费订单ID获取所有收款分配
func (r *refundRepository) GetRefundPaymentsByRefundID(refundOrderID int) ([]*refund.RefundPayment, error) {
	var payments []*refund.RefundPayment
	err := r.db.Where("refund_order_id = ?", refundOrderID).Find(&payments).Error
	return payments, err
}

// CreateTaobaoSupplement 创建淘宝退费补充信息
func (r *refundRepository) CreateTaobaoSupplement(supplement *refund.RefundTaobaoSupplement) error {
	return r.db.Create(supplement).Error
}

// ListTaobaoSupplements 获取淘宝退费补充信息列表
func (r *refundRepository) ListTaobaoSupplements(filters map[string]interface{}) ([]*refund.RefundTaobaoSupplement, error) {
	var list []*refund.RefundTaobaoSupplement
	query := r.db.Model(&refund.RefundTaobaoSupplement{})

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("id = ?", id)
	}
	if refundOrderID, ok := filters["refund_order_id"]; ok && refundOrderID != "" {
		query = query.Where("refund_order_id = ?", refundOrderID)
	}
	if uid, ok := filters["uid"]; ok && uid != "" {
		query = query.Where("student_id = ?", uid)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// 按创建时间倒序
	query = query.Order("create_time DESC")

	err := query.Find(&list).Error
	return list, err
}

// GetTaobaoSupplementByRefundID 根据退费订单ID获取淘宝补充信息
func (r *refundRepository) GetTaobaoSupplementByRefundID(refundOrderID int) (*refund.RefundTaobaoSupplement, error) {
	var supplement refund.RefundTaobaoSupplement
	err := r.db.Where("refund_order_id = ?", refundOrderID).First(&supplement).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &supplement, nil
}

// CreateRegularSupplement 创建常规退费补充信息
func (r *refundRepository) CreateRegularSupplement(supplement *refund.RefundRegularSupplement) error {
	return r.db.Create(supplement).Error
}

// ListRegularSupplements 获取常规退费补充信息列表
func (r *refundRepository) ListRegularSupplements(filters map[string]interface{}) ([]*refund.RefundRegularSupplement, error) {
	var list []*refund.RefundRegularSupplement
	query := r.db.Model(&refund.RefundRegularSupplement{})

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("id = ?", id)
	}
	if refundOrderID, ok := filters["refund_order_id"]; ok && refundOrderID != "" {
		query = query.Where("refund_order_id = ?", refundOrderID)
	}
	if uid, ok := filters["uid"]; ok && uid != "" {
		query = query.Where("student_id = ?", uid)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// 按创建时间倒序
	query = query.Order("create_time DESC")

	err := query.Find(&list).Error
	return list, err
}

// GetRegularSupplementsByRefundID 根据退费订单ID获取所有常规补充信息
func (r *refundRepository) GetRegularSupplementsByRefundID(refundOrderID int) ([]*refund.RefundRegularSupplement, error) {
	var supplements []*refund.RefundRegularSupplement
	err := r.db.Where("refund_order_id = ?", refundOrderID).Find(&supplements).Error
	return supplements, err
}
