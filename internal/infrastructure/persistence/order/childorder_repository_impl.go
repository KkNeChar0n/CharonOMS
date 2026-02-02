package order

import (
	"context"

	"gorm.io/gorm"

	"charonoms/internal/domain/order/repository"
)

// GormChildOrderRepository GORM实现的子订单仓储
type GormChildOrderRepository struct {
	db *gorm.DB
}

// NewChildOrderRepository 创建子订单仓储实例
func NewChildOrderRepository(db *gorm.DB) repository.ChildOrderRepository {
	return &GormChildOrderRepository{db: db}
}

// GetChildOrders 获取子订单列表（含商品信息）
func (r *GormChildOrderRepository) GetChildOrders(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.WithContext(ctx).
		Table("childorders c").
		Select(`
			c.id,
			c.parentsid,
			c.goodsid,
			g.name AS goods_name,
			c.amount_receivable,
			c.discount_amount,
			c.amount_received,
			c.status,
			c.create_time
		`).
		Joins("JOIN goods g ON c.goodsid = g.id").
		Order("c.id DESC").
		Find(&results).Error

	return results, err
}

// GetChildOrdersByParentID 根据父订单ID获取子订单列表
func (r *GormChildOrderRepository) GetChildOrdersByParentID(ctx context.Context, parentID int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.WithContext(ctx).
		Table("childorders c").
		Select(`
			c.id,
			c.parentsid,
			c.goodsid,
			g.name AS goods_name,
			c.amount_receivable,
			c.discount_amount,
			c.amount_received,
			c.status,
			c.create_time
		`).
		Joins("JOIN goods g ON c.goodsid = g.id").
		Where("c.parentsid = ?", parentID).
		Order("c.id ASC").
		Find(&results).Error

	return results, err
}

// UpdateChildOrderStatus 批量更新子订单状态
func (r *GormChildOrderRepository) UpdateChildOrderStatus(ctx context.Context, parentID int, status int) error {
	return r.db.WithContext(ctx).
		Model(&ChildOrderDO{}).
		Where("parentsid = ?", parentID).
		Update("status", status).Error
}
