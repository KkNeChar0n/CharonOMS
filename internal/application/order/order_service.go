package order

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"charonoms/internal/domain/financial/payment"
	"charonoms/internal/domain/financial/taobao"
	"charonoms/internal/domain/goods/repository"
	"charonoms/internal/domain/order/entity"
	orderRepo "charonoms/internal/domain/order/repository"
	"charonoms/internal/domain/order/service"
)

// Service 订单应用服务
type Service struct {
	orderRepo       orderRepo.OrderRepository
	childOrderRepo  orderRepo.ChildOrderRepository
	goodsRepo       repository.GoodsRepository
	paymentRepo     payment.PaymentRepository
	taobaoRepo      taobao.TaobaoPaymentRepository
	orderService    *service.OrderService
	discountService *service.DiscountService
	db              *gorm.DB
}

// NewService 创建订单服务实例
func NewService(
	orderRepo orderRepo.OrderRepository,
	childOrderRepo orderRepo.ChildOrderRepository,
	goodsRepo repository.GoodsRepository,
	paymentRepo payment.PaymentRepository,
	taobaoRepo taobao.TaobaoPaymentRepository,
	db *gorm.DB,
) *Service {
	return &Service{
		orderRepo:       orderRepo,
		childOrderRepo:  childOrderRepo,
		goodsRepo:       goodsRepo,
		paymentRepo:     paymentRepo,
		taobaoRepo:      taobaoRepo,
		orderService:    service.NewOrderService(),
		discountService: service.NewDiscountService(db),
		db:              db,
	}
}

// CreateOrder 创建订单
func (s *Service) CreateOrder(ctx context.Context, req *CreateOrderRequest) (int, error) {
	// 1. 验证请求
	if req.StudentID == 0 {
		return 0, errors.New("学生ID不能为空")
	}
	if len(req.GoodsList) == 0 {
		return 0, errors.New("必须至少选择一个商品")
	}

	// 2. 构建商品列表用于金额计算
	goodsItems := make([]service.GoodsItem, 0, len(req.GoodsList))
	for _, g := range req.GoodsList {
		goodsItems = append(goodsItems, service.GoodsItem{
			GoodsID:    g.GoodsID,
			TotalPrice: g.TotalPrice,
			Price:      g.Price,
		})
	}

	// 3. 计算订单金额
	amountReceivable, amountReceived := s.orderService.CalculateOrderAmounts(goodsItems, req.DiscountAmount)

	// 4. 创建订单实体
	order := &entity.Order{
		StudentID:           req.StudentID,
		ExpectedPaymentTime: req.ExpectedPaymentTime,
		AmountReceivable:    amountReceivable,
		AmountReceived:      amountReceived,
		DiscountAmount:      req.DiscountAmount,
		Status:              entity.OrderStatusDraft,
	}

	// 5. 验证订单金额
	if !order.ValidateAmounts() {
		return 0, errors.New("订单金额验证失败")
	}

	// 6. 创建子订单
	childOrders := make([]*entity.ChildOrder, 0, len(req.GoodsList))
	for _, g := range req.GoodsList {
		childDiscount := req.ChildDiscounts[g.GoodsID]
		childAmountReceivable, childAmountReceived := s.orderService.CalculateChildAmounts(
			g.TotalPrice,
			g.Price,
			childDiscount,
		)

		childOrder := &entity.ChildOrder{
			GoodsID:          g.GoodsID,
			AmountReceivable: childAmountReceivable,
			AmountReceived:   childAmountReceived,
			DiscountAmount:   childDiscount,
			Status:           entity.ChildOrderStatusInit,
		}

		if !childOrder.ValidateAmounts() {
			return 0, fmt.Errorf("商品 %d 的金额验证失败", g.GoodsID)
		}

		childOrders = append(childOrders, childOrder)
	}

	// 7. 保存订单
	orderID, err := s.orderRepo.CreateOrder(ctx, order, childOrders, req.ActivityIDs)
	if err != nil {
		return 0, fmt.Errorf("创建订单失败: %w", err)
	}

	return orderID, nil
}

// GetOrders 获取订单列表
func (s *Service) GetOrders(ctx context.Context) ([]map[string]interface{}, error) {
	return s.orderRepo.GetOrders(ctx)
}

// GetOrderGoods 获取订单商品列表
func (s *Service) GetOrderGoods(ctx context.Context, orderID int) ([]map[string]interface{}, error) {
	return s.orderRepo.GetOrderGoods(ctx, orderID)
}

// UpdateOrder 更新订单
func (s *Service) UpdateOrder(ctx context.Context, orderID int, req *UpdateOrderRequest) error {
	// 1. 查询订单
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return fmt.Errorf("查询订单失败: %w", err)
	}

	// 2. 验证订单状态
	if !order.CanEdit() {
		return errors.New("只能编辑草稿状态的订单")
	}

	// 3. 验证请求
	if len(req.GoodsList) == 0 {
		return errors.New("必须至少选择一个商品")
	}

	// 4. 构建商品列表用于金额计算
	goodsItems := make([]service.GoodsItem, 0, len(req.GoodsList))
	for _, g := range req.GoodsList {
		goodsItems = append(goodsItems, service.GoodsItem{
			GoodsID:    g.GoodsID,
			TotalPrice: g.TotalPrice,
			Price:      g.Price,
		})
	}

	// 5. 重新计算订单金额
	amountReceivable, amountReceived := s.orderService.CalculateOrderAmounts(goodsItems, req.DiscountAmount)

	// 6. 更新订单实体
	order.ExpectedPaymentTime = req.ExpectedPaymentTime
	order.AmountReceivable = amountReceivable
	order.AmountReceived = amountReceived
	order.DiscountAmount = req.DiscountAmount

	// 7. 验证订单金额
	if !order.ValidateAmounts() {
		return errors.New("订单金额验证失败")
	}

	// 8. 创建新的子订单列表
	childOrders := make([]*entity.ChildOrder, 0, len(req.GoodsList))
	for _, g := range req.GoodsList {
		childDiscount := req.ChildDiscounts[g.GoodsID]
		childAmountReceivable, childAmountReceived := s.orderService.CalculateChildAmounts(
			g.TotalPrice,
			g.Price,
			childDiscount,
		)

		childOrder := &entity.ChildOrder{
			GoodsID:          g.GoodsID,
			AmountReceivable: childAmountReceivable,
			AmountReceived:   childAmountReceived,
			DiscountAmount:   childDiscount,
			Status:           entity.ChildOrderStatusInit,
		}

		if !childOrder.ValidateAmounts() {
			return fmt.Errorf("商品 %d 的金额验证失败", g.GoodsID)
		}

		childOrders = append(childOrders, childOrder)
	}

	// 9. 保存更新
	err = s.orderRepo.UpdateOrder(ctx, order, childOrders, req.ActivityIDs)
	if err != nil {
		return fmt.Errorf("更新订单失败: %w", err)
	}

	return nil
}

// SubmitOrder 提交订单
func (s *Service) SubmitOrder(ctx context.Context, orderID int) error {
	// 1. 查询订单
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return fmt.Errorf("查询订单失败: %w", err)
	}

	// 2. 验证订单状态
	if !order.CanSubmit() {
		return errors.New("只能提交草稿状态的订单")
	}

	// 3. 更新订单和子订单状态
	err = s.orderRepo.UpdateOrderStatus(ctx, orderID, entity.OrderStatusUnpaid, entity.ChildOrderStatusUnpaid)
	if err != nil {
		return fmt.Errorf("提交订单失败: %w", err)
	}

	return nil
}

// CancelOrder 作废订单
func (s *Service) CancelOrder(ctx context.Context, orderID int) error {
	// 1. 查询订单
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return fmt.Errorf("查询订单失败: %w", err)
	}

	// 2. 验证订单状态
	if !order.CanCancel() {
		return errors.New("只能作废草稿状态的订单")
	}

	// 3. 更新订单和子订单状态
	err = s.orderRepo.UpdateOrderStatus(ctx, orderID, entity.OrderStatusCancelled, entity.ChildOrderStatusCancelled)
	if err != nil {
		return fmt.Errorf("作废订单失败: %w", err)
	}

	return nil
}

// GetChildOrders 获取子订单列表
func (s *Service) GetChildOrders(ctx context.Context) ([]map[string]interface{}, error) {
	return s.childOrderRepo.GetChildOrders(ctx)
}

// GetActiveGoodsForOrder 获取启用商品列表（用于订单）
func (s *Service) GetActiveGoodsForOrder(ctx context.Context) ([]map[string]interface{}, error) {
	return s.goodsRepo.GetActiveGoodsForOrder(ctx)
}

// GetGoodsTotalPrice 获取商品总价
func (s *Service) GetGoodsTotalPrice(ctx context.Context, goodsID int) (map[string]interface{}, error) {
	return s.goodsRepo.GetGoodsTotalPrice(ctx, goodsID)
}

// CalculateOrderDiscount 计算订单优惠
func (s *Service) CalculateOrderDiscount(ctx context.Context, goodsList []GoodsItemRequest, activityIDs []int) (float64, map[int]float64, error) {
	// 转换为领域服务需要的格式
	goodsForDiscount := make([]service.GoodsForDiscount, 0, len(goodsList))
	for _, g := range goodsList {
		goodsForDiscount = append(goodsForDiscount, service.GoodsForDiscount{
			GoodsID: g.GoodsID,
			Price:   g.Price,
		})
	}

	// 调用领域服务计算优惠
	return s.discountService.CalculateDiscount(ctx, goodsForDiscount, activityIDs)
}

// GetUnpaidOrdersByStudentID 获取学生的未付款订单列表
func (s *Service) GetUnpaidOrdersByStudentID(ctx context.Context, studentID int) ([]*entity.Order, error) {
	return s.orderRepo.GetUnpaidOrdersByStudentID(ctx, studentID)
}

// GetOrderPendingAmount 获取订单待付金额
func (s *Service) GetOrderPendingAmount(ctx context.Context, orderID int) (float64, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return 0, err
	}

	// 获取常规收款已付金额
	regularPaidAmount, err := s.paymentRepo.GetTotalPaidAmount(orderID)
	if err != nil {
		return 0, err
	}

	// 获取淘宝收款已付金额
	taobaoPaidAmount, err := s.taobaoRepo.GetTotalPaid(orderID)
	if err != nil {
		return 0, err
	}

	// 总已付金额 = 常规收款 + 淘宝收款
	totalPaidAmount := regularPaidAmount + taobaoPaidAmount

	// 计算待付金额 = 实收金额 - 已付金额
	// 注意：使用 AmountReceived（实收金额），因为这是用户实际需要支付的金额（已扣除优惠）
	pendingAmount := order.AmountReceived - totalPaidAmount
	if pendingAmount < 0 {
		pendingAmount = 0
	}

	return pendingAmount, nil
}

// GetOrderRefundInfo 获取订单退费信息
func (s *Service) GetOrderRefundInfo(ctx context.Context, orderID int) (map[string]interface{}, error) {
	// 1. 获取订单信息
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	// 2. 检查订单状态（必须是部分支付或已支付）
	if order.Status != 30 && order.Status != 40 {
		return nil, errors.New("订单状态不允许申请退费")
	}

	// 3. 获取子订单列表
	childOrders, err := s.childOrderRepo.ListByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("查询子订单失败: %w", err)
	}

	// 4. 计算每个子订单的可退金额
	childOrderData := make([]map[string]interface{}, 0, len(childOrders))
	for _, child := range childOrders {
		// 获取商品信息
		goods, err := s.goodsRepo.GetByID(child.GoodsID)
		if err != nil {
			continue
		}
		goodsName, _ := goods["name"].(string)

		// 计算已退金额：查询 refund_order_item 表中所有已通过（status=10）的退费
		var refundedAmount float64
		err = s.db.Raw(`
			SELECT COALESCE(SUM(refund_amount), 0)
			FROM refund_order_item
			WHERE childorder_id = ? AND status = 10
		`, child.ID).Scan(&refundedAmount).Error
		if err != nil {
			return nil, fmt.Errorf("查询已退金额失败: %w", err)
		}

		// 可退金额 = 实收金额 - 已退金额
		availableRefund := child.AmountReceived - refundedAmount
		if availableRefund < 0 {
			availableRefund = 0
		}

		childOrderData = append(childOrderData, map[string]interface{}{
			"childorder_id":    child.ID,
			"goods_id":         child.GoodsID,
			"goods_name":       goodsName,
			"amount_received":  child.AmountReceived,
			"refunded_amount":  refundedAmount,
			"available_refund": availableRefund,
		})
	}

	// 5. 返回订单信息和子订单列表
	result := map[string]interface{}{
		"order": map[string]interface{}{
			"id":         order.ID,
			"student_id": order.StudentID,
			"status":     order.Status,
		},
		"child_orders": childOrderData,
	}

	return result, nil
}

// toInt 安全地将interface{}转换为int，支持int32、int64、int
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint32:
		return int(val)
	case uint64:
		return int(val)
	default:
		return 0
	}
}

// toFloat64 安全地将interface{}转换为float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case string:
		// 尝试解析字符串为float64
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0
	}
}

// GetRefundPayments 获取退费收款列表（严格按照Python版本实现）
func (s *Service) GetRefundPayments(ctx context.Context, orderID int, refundItems []map[string]interface{}) (map[string]interface{}, error) {
	// 1. 提取子订单ID列表
	childOrderIDs := make([]int, 0, len(refundItems))
	for _, item := range refundItems {
		childOrderID := int(item["childorder_id"].(float64))
		childOrderIDs = append(childOrderIDs, childOrderID)
	}

	// 2. 查询该订单的所有常规收款（status=20已支付），并计算历史累计已退金额
	type PaymentWithRefund struct {
		PaymentID           int
		PaymentType         int
		PaymentAmount       float64
		PayeeEntity         *int
		Payer               *string
		IsCorporateTransfer int
		RefundedAmount      float64
	}

	regularPayments := []PaymentWithRefund{}
	err := s.db.Raw(`
		SELECT
			pc.id as payment_id,
			0 as payment_type,
			pc.payment_amount,
			pc.payee_entity,
			pc.payer,
			CASE WHEN pc.payment_method = 9 THEN 1 ELSE 0 END as is_corporate_transfer,
			COALESCE(SUM(rp.refund_amount), 0) as refunded_amount
		FROM payment_collection pc
		LEFT JOIN refund_payment rp ON pc.id = rp.payment_id AND rp.payment_type = 0
		LEFT JOIN refund_order ro ON ro.id = rp.refund_order_id AND ro.status = 10
		WHERE pc.order_id = ? AND pc.status = 20
		GROUP BY pc.id
		ORDER BY pc.id ASC
	`, orderID).Scan(&regularPayments).Error
	if err != nil {
		return nil, fmt.Errorf("查询常规收款失败: %w", err)
	}

	// 3. 查询该订单的所有淘宝收款（status=30已到账），并计算历史累计已退金额
	taobaoPayments := []PaymentWithRefund{}
	err = s.db.Raw(`
		SELECT
			tp.id as payment_id,
			1 as payment_type,
			tp.payment_amount,
			NULL as payee_entity,
			tp.payer,
			NULL as is_corporate_transfer,
			COALESCE(SUM(rp.refund_amount), 0) as refunded_amount
		FROM taobao_payment tp
		LEFT JOIN refund_payment rp ON tp.id = rp.payment_id AND rp.payment_type = 1
		LEFT JOIN refund_order ro ON ro.id = rp.refund_order_id AND ro.status = 10
		WHERE tp.order_id = ? AND tp.status IN (30, 40)
		GROUP BY tp.id
		ORDER BY tp.id ASC
	`, orderID).Scan(&taobaoPayments).Error
	if err != nil {
		return nil, fmt.Errorf("查询淘宝收款失败: %w", err)
	}

	// 4. 组装所有收款并计算可退金额和分账金额
	payments := make([]map[string]interface{}, 0)

	// 处理常规收款
	for _, p := range regularPayments {
		// 计算可退金额 = 收款金额 - 历史累计已退金额
		availableRefund := p.PaymentAmount - p.RefundedAmount
		if availableRefund < 0 {
			availableRefund = 0
		}

		// 计算该收款与待退费子订单关联的分账金额
		var separateAmount float64
		if len(childOrderIDs) > 0 {
			err := s.db.Raw(`
				SELECT COALESCE(SUM(separate_amount), 0)
				FROM separate_account
				WHERE payment_id = ? AND payment_type = 0 AND childorders_id IN ? AND type = 0
			`, p.PaymentID, childOrderIDs).Scan(&separateAmount).Error
			if err != nil {
				separateAmount = 0
			}
		}

		payment := map[string]interface{}{
			"payment_id":            p.PaymentID,
			"payment_type":          0,
			"payment_amount":        p.PaymentAmount,
			"refunded_amount":       p.RefundedAmount,
			"available_refund":      availableRefund,
			"separate_amount":       separateAmount,
			"is_corporate_transfer": p.IsCorporateTransfer,
		}

		if p.PayeeEntity != nil {
			payment["payee_entity"] = *p.PayeeEntity
		}
		if p.Payer != nil {
			payment["payer"] = *p.Payer
		}

		payments = append(payments, payment)
	}

	// 处理淘宝收款
	for _, p := range taobaoPayments {
		// 计算可退金额 = 收款金额 - 历史累计已退金额
		availableRefund := p.PaymentAmount - p.RefundedAmount
		if availableRefund < 0 {
			availableRefund = 0
		}

		// 计算该收款与待退费子订单关联的分账金额
		var separateAmount float64
		if len(childOrderIDs) > 0 {
			err := s.db.Raw(`
				SELECT COALESCE(SUM(separate_amount), 0)
				FROM separate_account
				WHERE payment_id = ? AND payment_type = 1 AND childorders_id IN ? AND type = 0
			`, p.PaymentID, childOrderIDs).Scan(&separateAmount).Error
			if err != nil {
				separateAmount = 0
			}
		}

		payment := map[string]interface{}{
			"payment_id":       p.PaymentID,
			"payment_type":     1,
			"payment_amount":   p.PaymentAmount,
			"refunded_amount":  p.RefundedAmount,
			"available_refund": availableRefund,
			"separate_amount":  separateAmount,
		}

		if p.Payer != nil {
			payment["payer"] = *p.Payer
		}

		payments = append(payments, payment)
	}

	// 5. 计算每个子订单的总分账金额
	childorderSeparateAmounts := make(map[int]float64)
	for _, childOrderID := range childOrderIDs {
		var totalSeparate float64
		err := s.db.Raw(`
			SELECT COALESCE(SUM(separate_amount), 0)
			FROM separate_account
			WHERE childorders_id = ? AND type = 0
		`, childOrderID).Scan(&totalSeparate).Error
		if err != nil {
			continue
		}
		childorderSeparateAmounts[childOrderID] = totalSeparate
	}

	result := map[string]interface{}{
		"payments":                    payments,
		"childorder_separate_amounts": childorderSeparateAmounts,
	}

	return result, nil
}
