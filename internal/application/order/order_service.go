package order

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

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
	orderService    *service.OrderService
	discountService *service.DiscountService
	db              *gorm.DB
}

// NewService 创建订单服务实例
func NewService(
	orderRepo orderRepo.OrderRepository,
	childOrderRepo orderRepo.ChildOrderRepository,
	goodsRepo repository.GoodsRepository,
	db *gorm.DB,
) *Service {
	return &Service{
		orderRepo:       orderRepo,
		childOrderRepo:  childOrderRepo,
		goodsRepo:       goodsRepo,
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
