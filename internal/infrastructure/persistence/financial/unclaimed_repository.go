package financial

import (
	"charonoms/internal/domain/financial/unclaimed"
	"time"

	"gorm.io/gorm"
)

type unclaimedRepository struct {
	db *gorm.DB
}

// NewUnclaimedRepository 创建常规待认领仓储实例
func NewUnclaimedRepository(db *gorm.DB) unclaimed.UnclaimedRepository {
	return &unclaimedRepository{db: db}
}

// Create 创建待认领记录
func (r *unclaimedRepository) Create(u *unclaimed.Unclaimed) error {
	return r.db.Create(u).Error
}

// GetByID 根据ID获取待认领记录
func (r *unclaimedRepository) GetByID(id int) (*unclaimed.Unclaimed, error) {
	var u unclaimed.Unclaimed
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// List 获取待认领列表
func (r *unclaimedRepository) List(filters map[string]interface{}) ([]*unclaimed.Unclaimed, error) {
	var list []*unclaimed.Unclaimed
	query := r.db.Model(&unclaimed.Unclaimed{})

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("id = ?", id)
	}
	if payer, ok := filters["payer"]; ok && payer != "" {
		query = query.Where("payer LIKE ?", "%"+payer.(string)+"%")
	}
	if paymentMethod, ok := filters["payment_method"]; ok && paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}
	if arrivalDate, ok := filters["arrival_date"]; ok && arrivalDate != "" {
		query = query.Where("DATE(arrival_time) = ?", arrivalDate)
	}
	if claimer, ok := filters["claimer"]; ok && claimer != "" {
		query = query.Where("claimer = ?", claimer)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// 按创建时间倒序
	query = query.Order("create_time DESC")

	err := query.Find(&list).Error
	return list, err
}

// Update 更新待认领记录
func (r *unclaimedRepository) Update(u *unclaimed.Unclaimed) error {
	return r.db.Save(u).Error
}

// Delete 删除待认领记录
func (r *unclaimedRepository) Delete(id int) error {
	return r.db.Delete(&unclaimed.Unclaimed{}, id).Error
}

// FindMatchingPayment 查找匹配的未核验收款记录
func (r *unclaimedRepository) FindMatchingPayment(merchantOrder string, paymentMethod int, paymentAmount float64, payeeEntity int) (int, error) {
	var paymentID int
	err := r.db.Table("payment_collection").
		Select("id").
		Where("payment_scenario = ?", 1). // 线下场景
		Where("status = ?", 10).           // 未核验状态
		Where("payment_method = ?", paymentMethod).
		Where("payment_amount = ?", paymentAmount).
		Where("merchant_order = ?", merchantOrder).
		Where("payee_entity = ?", payeeEntity).
		Limit(1).
		Scan(&paymentID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return paymentID, nil
}

// UpdatePaymentStatus 更新收款记录状态为已支付
func (r *unclaimedRepository) UpdatePaymentStatus(paymentID int, arrivalTime *time.Time) error {
	return r.db.Table("payment_collection").
		Where("id = ?", paymentID).
		Updates(map[string]interface{}{
			"status":       20, // 已支付
			"arrival_time": arrivalTime,
		}).Error
}
