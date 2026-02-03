package financial

import (
	"charonoms/internal/domain/financial/separate"
	"gorm.io/gorm"
)

// SeparateAccountRepositoryImpl 分账明细仓储实现
type SeparateAccountRepositoryImpl struct {
	db *gorm.DB
}

// NewSeparateAccountRepository 创建分账明细仓储实例
func NewSeparateAccountRepository(db *gorm.DB) separate.SeparateAccountRepository {
	return &SeparateAccountRepositoryImpl{db: db}
}

// Create 创建分账明细
func (r *SeparateAccountRepositoryImpl) Create(account *separate.SeparateAccount) error {
	return r.db.Create(account).Error
}

// BatchCreate 批量创建分账明细
func (r *SeparateAccountRepositoryImpl) BatchCreate(accounts []*separate.SeparateAccount) error {
	if len(accounts) == 0 {
		return nil
	}
	return r.db.Create(&accounts).Error
}

// GetByID 根据ID查询分账明细
func (r *SeparateAccountRepositoryImpl) GetByID(id int) (*separate.SeparateAccount, error) {
	var account separate.SeparateAccount
	err := r.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// List 查询分账明细列表
func (r *SeparateAccountRepositoryImpl) List(filter separate.SeparateListFilter) ([]*separate.SeparateAccount, int64, error) {
	query := r.db.Model(&separate.SeparateAccount{})

	// 构建查询条件
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}
	if filter.UID != nil {
		query = query.Where("uid = ?", *filter.UID)
	}
	if filter.OrdersID != nil {
		query = query.Where("orders_id = ?", *filter.OrdersID)
	}
	if filter.ChildOrdersID != nil {
		query = query.Where("childorders_id = ?", *filter.ChildOrdersID)
	}
	if filter.GoodsID != nil {
		query = query.Where("goods_id = ?", *filter.GoodsID)
	}
	if filter.PaymentID != nil {
		query = query.Where("payment_id = ?", *filter.PaymentID)
	}
	if filter.PaymentType != nil {
		query = query.Where("payment_type = ?", *filter.PaymentType)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	// 查询总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	var accounts []*separate.SeparateAccount
	offset := (filter.Page - 1) * filter.PageSize
	err = query.Order("create_time DESC").Offset(offset).Limit(filter.PageSize).Find(&accounts).Error
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// ExistsByPaymentAndOrder 检查指定收款和订单是否已生成分账
func (r *SeparateAccountRepositoryImpl) ExistsByPaymentAndOrder(paymentID int, orderID int, paymentType int) (bool, error) {
	var count int64
	err := r.db.Model(&separate.SeparateAccount{}).
		Where("payment_id = ? AND orders_id = ? AND payment_type = ?", paymentID, orderID, paymentType).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetChildOrderTotalSeparate 查询子订单的总分账金额（只统计售卖类型）
func (r *SeparateAccountRepositoryImpl) GetChildOrderTotalSeparate(childOrderID int) (float64, error) {
	var total float64
	err := r.db.Model(&separate.SeparateAccount{}).
		Where("childorders_id = ? AND type = ?", childOrderID, separate.SeparateTypeSale).
		Select("COALESCE(SUM(separate_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetChildOrderAllocatedAmount 查询子订单已分配金额
func (r *SeparateAccountRepositoryImpl) GetChildOrderAllocatedAmount(childOrderID int) (float64, error) {
	var allocated float64
	err := r.db.Model(&separate.SeparateAccount{}).
		Where("childorders_id = ?", childOrderID).
		Select("COALESCE(SUM(separate_amount), 0)").
		Scan(&allocated).Error
	return allocated, err
}
