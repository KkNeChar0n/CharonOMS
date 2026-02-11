package refund

import (
	"charonoms/internal/domain/approval/repository"
	"charonoms/internal/domain/financial/refund"
	orderRepo "charonoms/internal/domain/order/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type RefundService struct {
	refundRepo       refund.RefundRepository
	orderRepo        orderRepo.OrderRepository
	approvalTplRepo  repository.ApprovalFlowTemplateRepository
	approvalFlowRepo repository.ApprovalFlowManagementRepository
	approvalNodeRepo repository.ApprovalNodeCaseRepository
	db               *gorm.DB
}

// NewRefundService 创建退费服务实例
func NewRefundService(
	refundRepo refund.RefundRepository,
	orderRepo orderRepo.OrderRepository,
	approvalTplRepo repository.ApprovalFlowTemplateRepository,
	approvalFlowRepo repository.ApprovalFlowManagementRepository,
	approvalNodeRepo repository.ApprovalNodeCaseRepository,
	db *gorm.DB,
) *RefundService {
	return &RefundService{
		refundRepo:       refundRepo,
		orderRepo:        orderRepo,
		approvalTplRepo:  approvalTplRepo,
		approvalFlowRepo: approvalFlowRepo,
		approvalNodeRepo: approvalNodeRepo,
		db:               db,
	}
}

// GetRefundOrders 获取退费订单列表
func (s *RefundService) GetRefundOrders(filters map[string]interface{}) ([]*refund.RefundOrder, error) {
	return s.refundRepo.ListRefundOrders(filters)
}

// GetRefundOrderDetail 获取退费订单详情
func (s *RefundService) GetRefundOrderDetail(refundOrderID int) (map[string]interface{}, error) {
	// 获取退费订单基本信息
	refundOrder, err := s.refundRepo.GetRefundOrderByID(refundOrderID)
	if err != nil {
		return nil, err
	}

	// 获取退费子订单明细
	refundItems, err := s.refundRepo.GetRefundOrderItemsByRefundID(refundOrderID)
	if err != nil {
		return nil, err
	}

	// 获取退费收款分配
	refundPayments, err := s.refundRepo.GetRefundPaymentsByRefundID(refundOrderID)
	if err != nil {
		return nil, err
	}

	// 获取淘宝退费补充信息
	taobaoSupplement, err := s.refundRepo.GetTaobaoSupplementByRefundID(refundOrderID)
	if err != nil {
		return nil, err
	}

	// 获取常规退费补充信息
	regularSupplements, err := s.refundRepo.GetRegularSupplementsByRefundID(refundOrderID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"refund_order":         refundOrder,
		"refund_items":         refundItems,
		"refund_payments":      refundPayments,
		"taobao_supplement":    taobaoSupplement,
		"regular_supplements":  regularSupplements,
	}

	return result, nil
}

// GetRefundChildOrders 获取退费子订单列表
func (s *RefundService) GetRefundChildOrders(filters map[string]interface{}) ([]*refund.RefundOrderItem, error) {
	return s.refundRepo.ListRefundOrderItems(filters)
}

// GetRefundRegularSupplements 获取常规退费补充信息列表
func (s *RefundService) GetRefundRegularSupplements(filters map[string]interface{}) ([]*refund.RefundRegularSupplement, error) {
	return s.refundRepo.ListRegularSupplements(filters)
}

// GetRefundTaobaoSupplements 获取淘宝退费补充信息列表
func (s *RefundService) GetRefundTaobaoSupplements(filters map[string]interface{}) ([]*refund.RefundTaobaoSupplement, error) {
	return s.refundRepo.ListTaobaoSupplements(filters)
}

// GetRefundPaymentDetails 获取退费支付明细列表
func (s *RefundService) GetRefundPaymentDetails(filters map[string]interface{}) ([]*refund.RefundPayment, error) {
	return s.refundRepo.ListRefundPayments(filters)
}

// CreateRefundOrderRequest 创建退费订单请求
type CreateRefundOrderRequest struct {
	OrderID            int                                `json:"order_id"`
	RefundItems        []RefundItemRequest                `json:"refund_items"`
	RefundPayments     []RefundPaymentRequest             `json:"refund_payments"`
	TaobaoSupplement   *TaobaoSupplementRequest           `json:"taobao_supplement"`
	RegularSupplements []RegularSupplementRequest         `json:"regular_supplements"`
}

type RefundItemRequest struct {
	ChildOrderID int     `json:"childorder_id"`
	GoodsID      int     `json:"goods_id"`
	GoodsName    string  `json:"goods_name"`
	RefundAmount float64 `json:"refund_amount,string"`
}

type RefundPaymentRequest struct {
	PaymentID    int     `json:"payment_id"`
	PaymentType  int     `json:"payment_type"`
	RefundAmount float64 `json:"refund_amount,string"`
}

type TaobaoSupplementRequest struct {
	AlipayAccount string  `json:"alipay_account"`
	AlipayName    string  `json:"alipay_name"`
	RefundAmount  float64 `json:"refund_amount"`
}

type RegularSupplementRequest struct {
	PayeeEntity         *int    `json:"payee_entity"`
	IsCorporateTransfer *int    `json:"is_corporate_transfer"`
	Payer               string  `json:"payer"`
	BankAccount         string  `json:"bank_account"`
	PayerReadonly       *bool   `json:"payer_readonly"`
	RefundAmount        float64 `json:"refund_amount"`
}

// CreateRefundOrder 创建退费订单
func (s *RefundService) CreateRefundOrder(ctx context.Context, req *CreateRefundOrderRequest, username string, userID int) (int, error) {
	// 1. 参数校验
	if req.OrderID == 0 || len(req.RefundItems) == 0 || len(req.RefundPayments) == 0 {
		return 0, errors.New("参数不完整")
	}

	// 计算退费总额
	var refundTotal float64
	for _, item := range req.RefundItems {
		refundTotal += item.RefundAmount
	}

	// 计算收款分配总额
	var paymentTotal float64
	for _, p := range req.RefundPayments {
		paymentTotal += p.RefundAmount
	}

	// 验证金额一致性
	if refundTotal-paymentTotal > 0.01 || refundTotal-paymentTotal < -0.01 {
		return 0, fmt.Errorf("退费金额(%.2f)与收款分配金额(%.2f)不一致", refundTotal, paymentTotal)
	}

	// 2. 查询订单信息
	order, err := s.orderRepo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return 0, fmt.Errorf("订单不存在: %w", err)
	}

	// 3. 查询退费类型的审批流模板（status=0表示启用）
	var template struct {
		ID                   int
		ApprovalFlowTypeID   int
	}
	err = s.db.Raw(`
		SELECT aft.id, aft.approval_flow_type_id
		FROM approval_flow_template aft
		INNER JOIN approval_flow_type aftype ON aftype.id = aft.approval_flow_type_id
		WHERE aftype.name = '退费' AND aft.status = 0
		LIMIT 1
	`).Scan(&template).Error
	if err != nil || template.ID == 0 {
		return 0, errors.New("未找到启用的退费审批流模板")
	}

	templateID := template.ID
	flowTypeID := template.ApprovalFlowTypeID

	var refundOrderID int

	// 4. 在事务中执行
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 4.1 创建退费订单主记录
		refundOrder := &refund.RefundOrder{
			OrderID:      req.OrderID,
			StudentID:    order.StudentID,
			RefundAmount: refundTotal,
			Submitter:    username,
			SubmitTime:   time.Now(),
			Status:       refund.RefundStatusPending,
		}
		if err := tx.Create(refundOrder).Error; err != nil {
			return fmt.Errorf("创建退费订单失败: %w", err)
		}
		refundOrderID = refundOrder.ID

		// 4.2 创建退费子订单明细
		for _, item := range req.RefundItems {
			refundItem := &refund.RefundOrderItem{
				RefundOrderID: refundOrderID,
				ChildOrderID:  item.ChildOrderID,
				GoodsID:       item.GoodsID,
				GoodsName:     item.GoodsName,
				RefundAmount:  item.RefundAmount,
				Status:        refund.RefundStatusPending,
			}
			if err := tx.Create(refundItem).Error; err != nil {
				return fmt.Errorf("创建退费子订单明细失败: %w", err)
			}
		}

		// 4.3 创建退费收款分配
		for _, p := range req.RefundPayments {
			refundPayment := &refund.RefundPayment{
				RefundOrderID: refundOrderID,
				PaymentID:     p.PaymentID,
				PaymentType:   p.PaymentType,
				RefundAmount:  p.RefundAmount,
			}
			if err := tx.Create(refundPayment).Error; err != nil {
				return fmt.Errorf("创建退费收款分配失败: %w", err)
			}
		}

		// 4.4 保存淘宝退费补充信息
		if req.TaobaoSupplement != nil {
			taobaoSup := &refund.RefundTaobaoSupplement{
				RefundOrderID: refundOrderID,
				StudentID:     order.StudentID,
				AlipayAccount: req.TaobaoSupplement.AlipayAccount,
				AlipayName:    req.TaobaoSupplement.AlipayName,
				RefundAmount:  req.TaobaoSupplement.RefundAmount,
				Status:        refund.RefundStatusPending,
			}
			if err := tx.Create(taobaoSup).Error; err != nil {
				return fmt.Errorf("保存淘宝退费补充信息失败: %w", err)
			}
		}

		// 4.5 保存常规退费补充信息
		for _, rs := range req.RegularSupplements {
			var isCorporate *bool
			if rs.IsCorporateTransfer != nil {
				val := *rs.IsCorporateTransfer == 1
				isCorporate = &val
			}

			regularSup := &refund.RefundRegularSupplement{
				RefundOrderID:       refundOrderID,
				StudentID:           order.StudentID,
				PayeeEntity:         rs.PayeeEntity,
				IsCorporateTransfer: isCorporate,
				Payer:               rs.Payer,
				BankAccount:         rs.BankAccount,
				PayerReadonly:       rs.PayerReadonly,
				RefundAmount:        rs.RefundAmount,
				Status:              refund.RefundStatusPending,
			}
			if err := tx.Create(regularSup).Error; err != nil {
				return fmt.Errorf("保存常规退费补充信息失败: %w", err)
			}
		}

		// 4.6 更新订单状态为退费中(50)
		if err := tx.Exec("UPDATE orders SET status = 50 WHERE id = ?", req.OrderID).Error; err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}

		// 4.7 创建审批流实例
		if err := s.createApprovalFlowInstance(tx, templateID, flowTypeID, userID); err != nil {
			return fmt.Errorf("创建审批流实例失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return refundOrderID, nil
}

// createApprovalFlowInstance 创建审批流实例
func (s *RefundService) createApprovalFlowInstance(tx *gorm.DB, templateID int, flowTypeID int, createUserID int) error {
	fmt.Printf("[DEBUG] createApprovalFlowInstance: templateID=%d, flowTypeID=%d, createUserID=%d\n", templateID, flowTypeID, createUserID)

	// 1. 创建审批流管理记录
	result := tx.Exec(`
		INSERT INTO approval_flow_management
		(approval_flow_template_id, approval_flow_type_id, step, create_user, status, create_time)
		VALUES (?, ?, 0, ?, 0, NOW())
	`, templateID, flowTypeID, createUserID)
	if result.Error != nil {
		fmt.Printf("[ERROR] 插入审批流管理记录失败: %v\n", result.Error)
		return result.Error
	}
	fmt.Printf("[DEBUG] 审批流管理记录插入成功，影响行数: %d\n", result.RowsAffected)

	// 获取刚插入的ID
	var flowID int64
	if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&flowID).Error; err != nil {
		fmt.Printf("[ERROR] 获取审批流ID失败: %v\n", err)
		return err
	}
	fmt.Printf("[DEBUG] 审批流ID: %d\n", flowID)

	// 2. 查询模板的第一个审批节点（按sort排序）
	var firstNode struct {
		ID   int
		Type int
		Sort int
	}
	if err := tx.Raw(`
		SELECT id, type, sort
		FROM approval_flow_template_node
		WHERE template_id = ?
		ORDER BY sort ASC
		LIMIT 1
	`, templateID).Scan(&firstNode).Error; err != nil {
		fmt.Printf("[ERROR] 查询第一个审批节点失败: %v\n", err)
		return fmt.Errorf("查询第一个审批节点失败: %w", err)
	}
	fmt.Printf("[DEBUG] 第一个审批节点: ID=%d, Type=%d, Sort=%d\n", firstNode.ID, firstNode.Type, firstNode.Sort)

	// 3. 创建第一个审批节点实例
	result = tx.Exec(`
		INSERT INTO approval_node_case
		(node_id, approval_flow_management_id, type, sort, result, create_time)
		VALUES (?, ?, ?, ?, NULL, NOW())
	`, firstNode.ID, flowID, firstNode.Type, firstNode.Sort)
	if result.Error != nil {
		fmt.Printf("[ERROR] 创建审批节点实例失败: %v\n", result.Error)
		return result.Error
	}
	fmt.Printf("[DEBUG] 审批节点实例创建成功，影响行数: %d\n", result.RowsAffected)

	// 获取刚插入的ID
	var nodeCaseID int64
	if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&nodeCaseID).Error; err != nil {
		fmt.Printf("[ERROR] 获取节点实例ID失败: %v\n", err)
		return err
	}
	fmt.Printf("[DEBUG] 节点实例ID: %d\n", nodeCaseID)

	// 4. 查询该节点的所有审批人
	var approverIDs []int
	if err := tx.Raw(`
		SELECT useraccount_id
		FROM approval_node_useraccount
		WHERE node_id = ?
	`, firstNode.ID).Scan(&approverIDs).Error; err != nil {
		fmt.Printf("[ERROR] 查询审批人失败: %v\n", err)
		return fmt.Errorf("查询审批人失败: %w", err)
	}
	fmt.Printf("[DEBUG] 查询到审批人数量: %d, IDs: %v\n", len(approverIDs), approverIDs)

	// 5. 为每个审批人创建审批记录
	for _, approverID := range approverIDs {
		result := tx.Exec(`
			INSERT INTO approval_node_case_user
			(approval_node_case_id, useraccount_id, result, create_time)
			VALUES (?, ?, NULL, NOW())
		`, nodeCaseID, approverID)
		if result.Error != nil {
			fmt.Printf("[ERROR] 创建审批人记录失败 (approverID=%d): %v\n", approverID, result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 审批人记录创建成功: approverID=%d, 影响行数: %d\n", approverID, result.RowsAffected)
	}

	fmt.Println("[DEBUG] createApprovalFlowInstance 完成")
	return nil
}
