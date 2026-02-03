package payment

import (
	"context"
	"errors"
	"fmt"

	orderEntity "charonoms/internal/domain/order/entity"
	orderRepo "charonoms/internal/domain/order/repository"
)

// PaymentDomainService 收款领域服务
type PaymentDomainService struct {
	paymentRepo     PaymentRepository
	orderRepo       orderRepo.OrderRepository
	childOrderRepo  orderRepo.ChildOrderRepository
}

// NewPaymentDomainService 创建收款领域服务
func NewPaymentDomainService(
	paymentRepo PaymentRepository,
	orderRepo orderRepo.OrderRepository,
	childOrderRepo orderRepo.ChildOrderRepository,
) *PaymentDomainService {
	return &PaymentDomainService{
		paymentRepo:    paymentRepo,
		orderRepo:      orderRepo,
		childOrderRepo: childOrderRepo,
	}
}

// ValidatePaymentAmount 验证付款金额不超过待支付金额
func (s *PaymentDomainService) ValidatePaymentAmount(orderID int, paymentAmount float64) error {
	// 1. 查询订单实收金额
	ctx := context.Background()
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("查询订单失败: %w", err)
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	actualAmount := order.AmountReceived

	// 2. 计算已收款总额
	totalPaid, err := s.paymentRepo.GetTotalPaidAmount(orderID)
	if err != nil {
		return fmt.Errorf("计算已收款总额失败: %w", err)
	}

	// 3. 计算待支付金额
	unpaidAmount := actualAmount - totalPaid

	// 4. 验证付款金额不超过待支付金额
	if paymentAmount > unpaidAmount {
		return fmt.Errorf("付款金额%.2f不能超过待支付金额%.2f", paymentAmount, unpaidAmount)
	}

	return nil
}

// UpdateOrderPaymentStatus 更新订单支付状态
func (s *PaymentDomainService) UpdateOrderPaymentStatus(orderID int) error {
	// 1. 查询订单实收金额
	ctx := context.Background()
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("查询订单失败: %w", err)
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	actualAmount := order.AmountReceived

	// 2. 计算总收款金额（常规收款 + 淘宝收款）
	totalPaid, err := s.paymentRepo.GetTotalPaidAmount(orderID)
	if err != nil {
		return fmt.Errorf("计算总收款金额失败: %w", err)
	}

	// 3. 根据收款情况确定订单状态
	var newStatus int
	if totalPaid == 0 {
		newStatus = orderEntity.OrderStatusUnpaid // 20-未支付
	} else if totalPaid < actualAmount {
		newStatus = orderEntity.OrderStatusPartialPaid // 30-部分支付
	} else {
		newStatus = orderEntity.OrderStatusPaid // 40-已支付
	}

	// 4. 更新订单状态
	if order.Status != newStatus {
		err = s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatus, 0)
		if err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}
	}

	return nil
}
