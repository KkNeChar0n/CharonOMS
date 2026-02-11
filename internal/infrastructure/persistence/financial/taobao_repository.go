package financial

import (
	"charonoms/internal/domain/financial/taobao"
	"fmt"

	"gorm.io/gorm"
)

type taobaoPaymentRepository struct {
	db *gorm.DB
}

// NewTaobaoPaymentRepository 创建淘宝收款仓储实例
func NewTaobaoPaymentRepository(db *gorm.DB) taobao.TaobaoPaymentRepository {
	return &taobaoPaymentRepository{db: db}
}

func (r *taobaoPaymentRepository) Create(payment *taobao.TaobaoPayment) error {
	return r.db.Table("taobao_payment").Create(payment).Error
}

func (r *taobaoPaymentRepository) GetByID(id int) (*taobao.TaobaoPayment, error) {
	var payment taobao.TaobaoPayment
	err := r.db.Table("taobao_payment").Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *taobaoPaymentRepository) List(filters map[string]interface{}) ([]*taobao.TaobaoPayment, error) {
	var payments []*taobao.TaobaoPayment
	query := r.db.Table("taobao_payment")

	// 只返回状态为 0,30,40 的记录（已下单、已到账、已退单）
	query = query.Where("status IN ?", []int{
		taobao.TaobaoPaymentStatusOrdered,
		taobao.TaobaoPaymentStatusArrived,
		taobao.TaobaoPaymentStatusRefunded,
	})

	if id, ok := filters["id"].(int); ok && id > 0 {
		query = query.Where("id = ?", id)
	}
	if studentID, ok := filters["student_id"].(int); ok && studentID > 0 {
		query = query.Where("student_id = ?", studentID)
	}
	if orderID, ok := filters["order_id"].(int); ok && orderID > 0 {
		query = query.Where("order_id = ?", orderID)
	}
	if orderDate, ok := filters["order_date"].(string); ok && orderDate != "" {
		query = query.Where("DATE(order_time) = ?", orderDate)
	}
	if status, ok := filters["status"].(int); ok && status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Order("order_time DESC").Find(&payments).Error
	return payments, err
}

func (r *taobaoPaymentRepository) Update(payment *taobao.TaobaoPayment) error {
	return r.db.Table("taobao_payment").Where("id = ?", payment.ID).Updates(payment).Error
}

func (r *taobaoPaymentRepository) Delete(id int) error {
	return r.db.Table("taobao_payment").Where("id = ?", id).Delete(&taobao.TaobaoPayment{}).Error
}

func (r *taobaoPaymentRepository) GetTotalPaid(orderID int) (float64, error) {
	var total float64
	err := r.db.Table("taobao_payment").
		Select("COALESCE(SUM(payment_amount), 0) as total_paid").
		Where("order_id = ? AND status IN ?", orderID, []int{
			taobao.TaobaoPaymentStatusClaimed,  // 20-已认领
			taobao.TaobaoPaymentStatusArrived,  // 30-已到账
		}).
		Scan(&total).Error
	return total, err
}

func (r *taobaoPaymentRepository) ListUnclaimed(filters map[string]interface{}) ([]*taobao.TaobaoPayment, error) {
	var payments []*taobao.TaobaoPayment
	query := r.db.Table("taobao_payment")

	// 只返回状态为 10,20 的记录（待认领、已认领）
	query = query.Where("status IN ?", []int{
		taobao.TaobaoPaymentStatusUnclaimed,
		taobao.TaobaoPaymentStatusClaimed,
	})

	if id, ok := filters["id"].(int); ok && id > 0 {
		query = query.Where("id = ?", id)
	}
	if arrivalDate, ok := filters["arrival_date"].(string); ok && arrivalDate != "" {
		query = query.Where("DATE(arrival_time) = ?", arrivalDate)
	}
	if status, ok := filters["status"].(int); ok && status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Order("arrival_time DESC").Find(&payments).Error
	return payments, err
}

func (r *taobaoPaymentRepository) FindByMerchantOrderAndAmount(merchantOrder string, amount float64) (*taobao.TaobaoPayment, error) {
	var payment taobao.TaobaoPayment
	err := r.db.Table("taobao_payment").
		Where("merchant_order = ? AND payment_amount = ? AND status = ?",
			merchantOrder, amount, taobao.TaobaoPaymentStatusOrdered).
		First(&payment).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil // 返回 nil 表示未找到匹配记录
	}
	if err != nil {
		return nil, fmt.Errorf("查询淘宝收款记录失败: %v", err)
	}
	return &payment, nil
}
