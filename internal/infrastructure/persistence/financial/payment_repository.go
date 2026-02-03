package financial

import (
	"charonoms/internal/domain/financial/payment"
	"gorm.io/gorm"
)

// PaymentRepositoryImpl 收款仓储实现
type PaymentRepositoryImpl struct {
	db *gorm.DB
}

// NewPaymentRepository 创建收款仓储实例
func NewPaymentRepository(db *gorm.DB) payment.PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

// Create 创建收款记录
func (r *PaymentRepositoryImpl) Create(p *payment.PaymentCollection) error {
	return r.db.Create(p).Error
}

// GetByID 根据ID查询收款记录
func (r *PaymentRepositoryImpl) GetByID(id int) (*payment.PaymentCollection, error) {
	var p payment.PaymentCollection
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// List 查询收款列表
func (r *PaymentRepositoryImpl) List(filter payment.PaymentListFilter) ([]*payment.PaymentCollection, int64, error) {
	query := r.db.Model(&payment.PaymentCollection{})

	// 构建查询条件
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}
	if filter.StudentID != nil {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.OrderID != nil {
		query = query.Where("order_id = ?", *filter.OrderID)
	}
	if filter.Payer != nil && *filter.Payer != "" {
		query = query.Where("payer LIKE ?", "%"+*filter.Payer+"%")
	}
	if filter.PaymentMethod != nil {
		query = query.Where("payment_method = ?", *filter.PaymentMethod)
	}
	if filter.TradingDate != nil && *filter.TradingDate != "" {
		query = query.Where("DATE(trading_hours) = ?", *filter.TradingDate)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	// 查询总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	var payments []*payment.PaymentCollection
	offset := (filter.Page - 1) * filter.PageSize
	err = query.Order("create_time DESC").Offset(offset).Limit(filter.PageSize).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// Update 更新收款记录
func (r *PaymentRepositoryImpl) Update(p *payment.PaymentCollection) error {
	return r.db.Save(p).Error
}

// Delete 删除收款记录
func (r *PaymentRepositoryImpl) Delete(id int) error {
	return r.db.Delete(&payment.PaymentCollection{}, id).Error
}

// GetTotalPaidAmount 查询订单已收款总额
func (r *PaymentRepositoryImpl) GetTotalPaidAmount(orderID int) (float64, error) {
	var totalRegular float64
	// 常规收款：status=10或20
	err := r.db.Model(&payment.PaymentCollection{}).
		Where("order_id = ? AND status IN (?, ?)", orderID, payment.PaymentStatusUnverified, payment.PaymentStatusPaid).
		Select("COALESCE(SUM(payment_amount), 0)").
		Scan(&totalRegular).Error
	if err != nil {
		return 0, err
	}

	// 淘宝收款暂不实现，返回0
	totalTaobao := 0.0

	return totalRegular + totalTaobao, nil
}

// CountByOrderAndStatus 统计订单指定状态的收款数量
func (r *PaymentRepositoryImpl) CountByOrderAndStatus(orderID int, status int) (int64, error) {
	var count int64
	err := r.db.Model(&payment.PaymentCollection{}).
		Where("order_id = ? AND status = ?", orderID, status).
		Count(&count).Error
	return count, err
}
