package order

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"charonoms/internal/domain/order/entity"
	"charonoms/internal/domain/order/repository"
)

// GormOrderRepository GORM实现的订单仓储
type GormOrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储实例
func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &GormOrderRepository{db: db}
}

// GetOrders 获取订单列表（含学生信息）
func (r *GormOrderRepository) GetOrders(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.WithContext(ctx).
		Table("orders o").
		Select(`
			o.id,
			o.student_id AS uid,
			s.name AS student_name,
			o.expected_payment_time,
			o.amount_receivable,
			o.discount_amount,
			o.amount_received,
			o.create_time,
			o.status
		`).
		Joins("JOIN student s ON o.student_id = s.id").
		Order("o.create_time DESC").
		Find(&results).Error

	return results, err
}

// GetOrderByID 根据ID查询订单
func (r *GormOrderRepository) GetOrderByID(ctx context.Context, id int) (*entity.Order, error) {
	var order OrderDO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &entity.Order{
		ID:                  order.ID,
		StudentID:           order.StudentID,
		ExpectedPaymentTime: order.ExpectedPaymentTime,
		AmountReceivable:    order.AmountReceivable,
		AmountReceived:      order.AmountReceived,
		DiscountAmount:      order.DiscountAmount,
		Status:              order.Status,
		CreateTime:          order.CreateTime,
	}, nil
}

// CreateOrder 创建订单（含事务：订单、子订单、活动关联）
func (r *GormOrderRepository) CreateOrder(ctx context.Context, order *entity.Order, childOrders []*entity.ChildOrder, activityIDs []int) (int, error) {
	var orderID int

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 创建订单
		orderDO := &OrderDO{
			StudentID:           order.StudentID,
			ExpectedPaymentTime: order.ExpectedPaymentTime,
			AmountReceivable:    order.AmountReceivable,
			AmountReceived:      order.AmountReceived,
			DiscountAmount:      order.DiscountAmount,
			Status:              order.Status,
		}
		if err := tx.Create(orderDO).Error; err != nil {
			return fmt.Errorf("创建订单失败: %w", err)
		}
		orderID = orderDO.ID

		// 2. 创建子订单
		if len(childOrders) > 0 {
			childOrderDOs := make([]ChildOrderDO, 0, len(childOrders))
			for _, co := range childOrders {
				childOrderDOs = append(childOrderDOs, ChildOrderDO{
					ParentsID:        orderID,
					GoodsID:          co.GoodsID,
					AmountReceivable: co.AmountReceivable,
					AmountReceived:   co.AmountReceived,
					DiscountAmount:   co.DiscountAmount,
					Status:           co.Status,
				})
			}
			if err := tx.Create(&childOrderDOs).Error; err != nil {
				return fmt.Errorf("创建子订单失败: %w", err)
			}
		}

		// 3. 创建订单活动关联
		if len(activityIDs) > 0 {
			activityDOs := make([]OrdersActivityDO, 0, len(activityIDs))
			for _, activityID := range activityIDs {
				activityDOs = append(activityDOs, OrdersActivityDO{
					OrdersID:   orderID,
					ActivityID: activityID,
				})
			}
			if err := tx.Create(&activityDOs).Error; err != nil {
				return fmt.Errorf("创建订单活动关联失败: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// UpdateOrder 更新订单（含事务：订单、子订单、活动关联）
func (r *GormOrderRepository) UpdateOrder(ctx context.Context, order *entity.Order, childOrders []*entity.ChildOrder, activityIDs []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 更新订单
		updates := map[string]interface{}{
			"amount_receivable":      order.AmountReceivable,
			"amount_received":        order.AmountReceived,
			"discount_amount":        order.DiscountAmount,
			"expected_payment_time":  order.ExpectedPaymentTime,
		}
		if err := tx.Model(&OrderDO{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新订单失败: %w", err)
		}

		// 2. 删除旧的子订单
		if err := tx.Where("parentsid = ?", order.ID).Delete(&ChildOrderDO{}).Error; err != nil {
			return fmt.Errorf("删除旧子订单失败: %w", err)
		}

		// 3. 创建新的子订单
		if len(childOrders) > 0 {
			childOrderDOs := make([]ChildOrderDO, 0, len(childOrders))
			for _, co := range childOrders {
				childOrderDOs = append(childOrderDOs, ChildOrderDO{
					ParentsID:        order.ID,
					GoodsID:          co.GoodsID,
					AmountReceivable: co.AmountReceivable,
					AmountReceived:   co.AmountReceived,
					DiscountAmount:   co.DiscountAmount,
					Status:           co.Status,
				})
			}
			if err := tx.Create(&childOrderDOs).Error; err != nil {
				return fmt.Errorf("创建新子订单失败: %w", err)
			}
		}

		// 4. 删除旧的活动关联
		if err := tx.Where("orders_id = ?", order.ID).Delete(&OrdersActivityDO{}).Error; err != nil {
			return fmt.Errorf("删除旧活动关联失败: %w", err)
		}

		// 5. 创建新的活动关联
		if len(activityIDs) > 0 {
			activityDOs := make([]OrdersActivityDO, 0, len(activityIDs))
			for _, activityID := range activityIDs {
				activityDOs = append(activityDOs, OrdersActivityDO{
					OrdersID:   order.ID,
					ActivityID: activityID,
				})
			}
			if err := tx.Create(&activityDOs).Error; err != nil {
				return fmt.Errorf("创建新活动关联失败: %w", err)
			}
		}

		return nil
	})
}

// UpdateOrderStatus 更新订单状态（含子订单状态级联更新）
func (r *GormOrderRepository) UpdateOrderStatus(ctx context.Context, orderID int, orderStatus int, childOrderStatus int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 更新订单状态
		if err := tx.Model(&OrderDO{}).Where("id = ?", orderID).Update("status", orderStatus).Error; err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}

		// 2. 更新子订单状态
		if err := tx.Model(&ChildOrderDO{}).Where("parentsid = ?", orderID).Update("status", childOrderStatus).Error; err != nil {
			return fmt.Errorf("更新子订单状态失败: %w", err)
		}

		return nil
	})
}

// GetOrderGoods 获取订单商品列表（关联商品、品牌、分类、属性）
func (r *GormOrderRepository) GetOrderGoods(ctx context.Context, orderID int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 1. 查询子订单和商品基本信息
	err := r.db.WithContext(ctx).
		Table("childorders c").
		Select(`
			c.id,
			c.goodsid,
			g.name AS goods_name,
			g.price,
			c.amount_receivable AS total_price,
			c.discount_amount,
			(g.price - c.discount_amount) AS discounted_price,
			g.isgroup,
			b.name AS brand_name,
			cl.name AS classify_name,
			c.amount_receivable,
			c.amount_received
		`).
		Joins("JOIN goods g ON c.goodsid = g.id").
		Joins("LEFT JOIN brand b ON g.brandid = b.id").
		Joins("LEFT JOIN classify cl ON g.classifyid = cl.id").
		Where("c.parentsid = ?", orderID).
		Order("c.id ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// 2. 为每个商品查询属性信息
	for i := range results {
		// 安全地转换 goodsid
		var goodsID int
		switch v := results[i]["goodsid"].(type) {
		case int:
			goodsID = v
		case int64:
			goodsID = int(v)
		case int32:
			goodsID = int(v)
		case uint:
			goodsID = int(v)
		case uint32:
			goodsID = int(v)
		case uint64:
			goodsID = int(v)
		case float64:
			goodsID = int(v)
		default:
			results[i]["attributes"] = []map[string]interface{}{}
			continue
		}

		var attributes []map[string]interface{}
		err := r.db.WithContext(ctx).
			Table("goods_attributevalue gav").
			Select("a.name AS attr_name, av.name AS value_name").
			Joins("JOIN attribute_value av ON gav.attributevalueid = av.id").
			Joins("JOIN attribute a ON av.attributeid = a.id").
			Where("gav.goodsid = ?", goodsID).
			Order("a.id, av.id").
			Find(&attributes).Error

		if err == nil {
			results[i]["attributes"] = attributes
		} else {
			results[i]["attributes"] = []map[string]interface{}{}
		}
	}

	return results, nil
}

// DeleteOrderChildOrders 删除订单的所有子订单
func (r *GormOrderRepository) DeleteOrderChildOrders(ctx context.Context, orderID int) error {
	return r.db.WithContext(ctx).Where("parentsid = ?", orderID).Delete(&ChildOrderDO{}).Error
}

// DeleteOrderActivities 删除订单的所有活动关联
func (r *GormOrderRepository) DeleteOrderActivities(ctx context.Context, orderID int) error {
	return r.db.WithContext(ctx).Where("orders_id = ?", orderID).Delete(&OrdersActivityDO{}).Error
}
