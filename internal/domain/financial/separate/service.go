package separate

import (
	"errors"
	"fmt"
	"math"

	orderEntity "charonoms/internal/domain/order/entity"
	orderRepo "charonoms/internal/domain/order/repository"
	"charonoms/internal/domain/financial/payment"
)

// SeparateAccountDomainService 分账明细领域服务
type SeparateAccountDomainService struct {
	separateRepo   SeparateAccountRepository
	paymentRepo    payment.PaymentRepository
	childOrderRepo orderRepo.ChildOrderRepository
}

// NewSeparateAccountDomainService 创建分账明细领域服务
func NewSeparateAccountDomainService(
	separateRepo SeparateAccountRepository,
	paymentRepo payment.PaymentRepository,
	childOrderRepo orderRepo.ChildOrderRepository,
) *SeparateAccountDomainService {
	return &SeparateAccountDomainService{
		separateRepo:   separateRepo,
		paymentRepo:    paymentRepo,
		childOrderRepo: childOrderRepo,
	}
}

// GenerateSeparateAccounts 生成分账明细
func (s *SeparateAccountDomainService) GenerateSeparateAccounts(paymentID int, orderID int) error {
	// 1. 查询收款信息
	paymentCollection, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return fmt.Errorf("查询收款信息失败: %w", err)
	}
	if paymentCollection == nil {
		return errors.New("收款记录不存在")
	}

	studentID := paymentCollection.StudentID
	paymentAmount := paymentCollection.PaymentAmount

	// 2. 查询子订单列表（按ID升序）
	childOrders, err := s.childOrderRepo.ListByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("查询子订单列表失败: %w", err)
	}

	if len(childOrders) == 0 {
		// 没有子订单，记录警告但不报错
		return nil
	}

	// 3. 防重复检查
	exists, err := s.separateRepo.ExistsByPaymentAndOrder(paymentID, orderID, PaymentTypeRegular)
	if err != nil {
		return fmt.Errorf("检查分账是否存在失败: %w", err)
	}
	if exists {
		// 已生成分账，幂等性保证
		return nil
	}

	// 4. 遍历子订单，按序分配收款金额
	remainingAmount := paymentAmount
	var accounts []*SeparateAccount

	for _, child := range childOrders {
		// 4.1 计算子订单已分配金额
		allocatedAmount, err := s.separateRepo.GetChildOrderAllocatedAmount(child.ID)
		if err != nil {
			return fmt.Errorf("计算子订单已分配金额失败: %w", err)
		}

		// 4.2 计算子订单还需金额
		neededAmount := child.AmountReceived - allocatedAmount

		if neededAmount <= 0 {
			// 子订单已满额，跳过
			continue
		}

		// 4.3 计算本次分账金额（取剩余收款和还需金额的较小值）
		separateAmount := math.Min(remainingAmount, neededAmount)
		separateAmount = roundToTwoDecimal(separateAmount)

		// 4.4 创建分账明细记录
		account := &SeparateAccount{
			UID:            studentID,
			OrdersID:       orderID,
			ChildOrdersID:  child.ID,
			PaymentID:      paymentID,
			PaymentType:    PaymentTypeRegular,
			GoodsID:        child.GoodsID,
			GoodsName:      "", // 商品名称需要从商品表查询，暂时留空
			SeparateAmount: separateAmount,
			Type:           SeparateTypeSale,
		}
		accounts = append(accounts, account)

		// 4.5 更新剩余收款金额
		remainingAmount -= separateAmount
		remainingAmount = roundToTwoDecimal(remainingAmount)

		if remainingAmount <= 0 {
			break
		}
	}

	// 5. 批量插入分账明细
	if len(accounts) > 0 {
		err = s.separateRepo.BatchCreate(accounts)
		if err != nil {
			return fmt.Errorf("插入分账明细失败: %w", err)
		}

		// 6. 更新所有相关子订单的状态
		for _, account := range accounts {
			err = s.UpdateChildOrderStatus(account.ChildOrdersID)
			if err != nil {
				return fmt.Errorf("更新子订单状态失败: %w", err)
			}
		}
	}

	return nil
}

// UpdateChildOrderStatus 更新子订单状态
func (s *SeparateAccountDomainService) UpdateChildOrderStatus(childOrderID int) error {
	// 1. 查询子订单
	childOrder, err := s.childOrderRepo.GetByID(childOrderID)
	if err != nil {
		return fmt.Errorf("查询子订单失败: %w", err)
	}
	if childOrder == nil {
		return errors.New("子订单不存在")
	}

	actualAmount := childOrder.AmountReceived

	// 2. 计算子订单总分账金额（只统计售卖类型）
	totalSeparate, err := s.separateRepo.GetChildOrderTotalSeparate(childOrderID)
	if err != nil {
		return fmt.Errorf("计算子订单总分账金额失败: %w", err)
	}

	// 3. 根据分账情况确定子订单状态
	var newStatus int
	if totalSeparate == 0 {
		newStatus = orderEntity.ChildOrderStatusUnpaid // 10-未支付
	} else if totalSeparate < actualAmount {
		newStatus = orderEntity.ChildOrderStatusPartialPaid // 20-部分支付
	} else {
		newStatus = orderEntity.ChildOrderStatusPaid // 30-已支付
	}

	// 4. 更新子订单状态
	if childOrder.Status != newStatus {
		childOrder.Status = newStatus
		err = s.childOrderRepo.Update(childOrder)
		if err != nil {
			return fmt.Errorf("更新子订单状态失败: %w", err)
		}
	}

	return nil
}

// roundToTwoDecimal 四舍五入到两位小数
func roundToTwoDecimal(value float64) float64 {
	return math.Round(value*100) / 100
}
