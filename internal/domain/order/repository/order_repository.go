package repository

import (
	"charonoms/internal/domain/order/entity"
	"context"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// GetOrders 获取订单列表（含学生信息）
	GetOrders(ctx context.Context) ([]map[string]interface{}, error)

	// GetOrderByID 根据ID查询订单
	GetOrderByID(ctx context.Context, id int) (*entity.Order, error)

	// CreateOrder 创建订单（含事务：订单、子订单、活动关联）
	CreateOrder(ctx context.Context, order *entity.Order, childOrders []*entity.ChildOrder, activityIDs []int) (int, error)

	// UpdateOrder 更新订单（含事务：订单、子订单、活动关联）
	UpdateOrder(ctx context.Context, order *entity.Order, childOrders []*entity.ChildOrder, activityIDs []int) error

	// UpdateOrderStatus 更新订单状态（含子订单状态级联更新）
	UpdateOrderStatus(ctx context.Context, orderID int, orderStatus int, childOrderStatus int) error

	// GetOrderGoods 获取订单商品列表（关联商品、品牌、分类、属性）
	GetOrderGoods(ctx context.Context, orderID int) ([]map[string]interface{}, error)

	// DeleteOrderChildOrders 删除订单的所有子订单
	DeleteOrderChildOrders(ctx context.Context, orderID int) error

	// DeleteOrderActivities 删除订单的所有活动关联
	DeleteOrderActivities(ctx context.Context, orderID int) error

	// GetUnpaidOrdersByStudentID 获取学生的未付款订单列表
	GetUnpaidOrdersByStudentID(ctx context.Context, studentID int) ([]*entity.Order, error)
}
