package repository

import (
	"context"
)

// ChildOrderRepository 子订单仓储接口
type ChildOrderRepository interface {
	// GetChildOrders 获取子订单列表（含商品信息）
	GetChildOrders(ctx context.Context) ([]map[string]interface{}, error)

	// GetChildOrdersByParentID 根据父订单ID获取子订单列表
	GetChildOrdersByParentID(ctx context.Context, parentID int) ([]map[string]interface{}, error)

	// UpdateChildOrderStatus 批量更新子订单状态
	UpdateChildOrderStatus(ctx context.Context, parentID int, status int) error
}
