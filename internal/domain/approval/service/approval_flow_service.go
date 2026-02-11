package service

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
	orderEntity "charonoms/internal/domain/order/entity"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ApprovalFlowService 审批流领域服务
type ApprovalFlowService struct {
	flowRepo     repository.ApprovalFlowManagementRepository
	nodeCaseRepo repository.ApprovalNodeCaseRepository
	templateRepo repository.ApprovalFlowTemplateRepository
	db           *gorm.DB
}

// NewApprovalFlowService 创建审批流领域服务
func NewApprovalFlowService(
	flowRepo repository.ApprovalFlowManagementRepository,
	nodeCaseRepo repository.ApprovalNodeCaseRepository,
	templateRepo repository.ApprovalFlowTemplateRepository,
	db *gorm.DB,
) *ApprovalFlowService {
	return &ApprovalFlowService{
		flowRepo:     flowRepo,
		nodeCaseRepo: nodeCaseRepo,
		templateRepo: templateRepo,
		db:           db,
	}
}

// ProcessApprove 处理审批通过逻辑
func (s *ApprovalFlowService) ProcessApprove(nodeCaseUserID int) error {
	// 1. 获取审批人员记录
	nodeCaseUser, err := s.nodeCaseRepo.GetNodeUserByID(nodeCaseUserID)
	if err != nil {
		return err
	}

	// 检查是否已处理
	if nodeCaseUser.Result != nil {
		return errors.New("该审批已处理")
	}

	// 2. 更新当前用户审批结果为通过
	result := int8(0)
	if err := s.nodeCaseRepo.UpdateUserResult(nodeCaseUserID, result); err != nil {
		return err
	}

	// 3. 获取节点实例
	nodeCase, err := s.nodeCaseRepo.GetByID(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 4. 获取该节点的所有审批人员
	nodeUsers, err := s.nodeCaseRepo.GetNodeUsers(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 5. 判断节点是否通过
	nodePassed := false
	if nodeCase.Type == 0 {
		// 会签节点：所有人都通过才能通过
		nodePassed = s.isCountersignNodePassed(nodeUsers)
	} else {
		// 或签节点：任意一人通过即可通过
		nodePassed = s.isOrSignNodePassed(nodeUsers)
		if nodePassed {
			// 删除同节点其他待审批人员
			if err := s.nodeCaseRepo.DeletePendingUsers(nodeCaseUser.ApprovalNodeCaseID, nodeCaseUserID); err != nil {
				return err
			}
		}
	}

	// 6. 如果节点通过，更新节点结果并流转
	if nodePassed {
		if err := s.nodeCaseRepo.UpdateNodeResult(nodeCaseUser.ApprovalNodeCaseID, 0); err != nil {
			return err
		}

		// 7. 流转到下一节点或完成审批流
		return s.proceedToNextNodeOrComplete(nodeCase, nodeCaseUser.ApprovalNodeCaseID)
	}

	return nil
}

// ProcessReject 处理审批驳回逻辑
func (s *ApprovalFlowService) ProcessReject(nodeCaseUserID int) error {
	// 1. 获取审批人员记录
	nodeCaseUser, err := s.nodeCaseRepo.GetNodeUserByID(nodeCaseUserID)
	if err != nil {
		return err
	}

	// 检查是否已处理
	if nodeCaseUser.Result != nil {
		return errors.New("该审批已处理")
	}

	// 2. 更新当前用户审批结果为驳回
	result := int8(1)
	if err := s.nodeCaseRepo.UpdateUserResult(nodeCaseUserID, result); err != nil {
		return err
	}

	// 3. 获取节点实例
	nodeCase, err := s.nodeCaseRepo.GetByID(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 4. 获取该节点的所有审批人员
	nodeUsers, err := s.nodeCaseRepo.GetNodeUsers(nodeCaseUser.ApprovalNodeCaseID)
	if err != nil {
		return err
	}

	// 5. 判断节点是否驳回
	nodeRejected := false
	if nodeCase.Type == 0 {
		// 会签节点：任意一人驳回即驳回
		nodeRejected = s.isCountersignNodeRejected(nodeUsers)
		if nodeRejected {
			// 删除同节点其他待审批人员
			if err := s.nodeCaseRepo.DeletePendingUsers(nodeCaseUser.ApprovalNodeCaseID, nodeCaseUserID); err != nil {
				return err
			}
		}
	} else {
		// 或签节点：所有人都驳回才驳回
		nodeRejected = s.isOrSignNodeRejected(nodeUsers)
	}

	// 6. 如果节点驳回，更新节点结果和审批流状态
	if nodeRejected {
		if err := s.nodeCaseRepo.UpdateNodeResult(nodeCaseUser.ApprovalNodeCaseID, 1); err != nil {
			return err
		}

		// 更新审批流状态为已驳回(20)
		// 需要从nodeCase获取flowID
		if err := s.flowRepo.UpdateStatus(nodeCase.ApprovalFlowManagementID, 20); err != nil {
			return err
		}

		// 审批流驳回后的回调处理
		if err := s.handleApprovalComplete(nodeCase.ApprovalFlowManagementID, false); err != nil {
			return err
		}
	}

	return nil
}

// isCountersignNodePassed 判断会签节点是否通过
func (s *ApprovalFlowService) isCountersignNodePassed(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 所有人都审批 && 所有人都通过 → 节点通过
	for _, user := range nodeUsers {
		if user.Result == nil {
			return false // 有人未审批
		}
		if *user.Result != 0 {
			return false // 有人驳回
		}
	}
	return true
}

// isCountersignNodeRejected 判断会签节点是否驳回
func (s *ApprovalFlowService) isCountersignNodeRejected(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 任意一人驳回 → 节点驳回
	for _, user := range nodeUsers {
		if user.Result != nil && *user.Result == 1 {
			return true
		}
	}
	return false
}

// isOrSignNodePassed 判断或签节点是否通过
func (s *ApprovalFlowService) isOrSignNodePassed(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 任意一人通过 → 节点通过
	for _, user := range nodeUsers {
		if user.Result != nil && *user.Result == 0 {
			return true
		}
	}
	return false
}

// isOrSignNodeRejected 判断或签节点是否驳回
func (s *ApprovalFlowService) isOrSignNodeRejected(nodeUsers []entity.ApprovalNodeCaseUser) bool {
	// 所有人都审批 && 所有人都驳回 → 节点驳回
	for _, user := range nodeUsers {
		if user.Result == nil {
			return false // 有人未审批
		}
		if *user.Result != 1 {
			return false // 有人通过
		}
	}
	return true
}

// handleApprovalComplete 处理审批流完成后的回调
func (s *ApprovalFlowService) handleApprovalComplete(flowID int, approved bool) error {
	fmt.Printf("[DEBUG] handleApprovalComplete调用: flowID=%d, approved=%v\n", flowID, approved)

	// 1. 获取审批流信息和类型名称
	var flowInfo struct {
		ID                 int
		CreateTime         string
		ApprovalFlowTypeID int
	}
	err := s.db.Table("approval_flow_management").
		Select("id, create_time, approval_flow_type_id").
		Where("id = ?", flowID).
		Scan(&flowInfo).Error
	if err != nil {
		return fmt.Errorf("查询审批流信息失败: %w", err)
	}

	fmt.Printf("[DEBUG] 审批流信息: ID=%d, CreateTime=%s, TypeID=%d\n", flowInfo.ID, flowInfo.CreateTime, flowInfo.ApprovalFlowTypeID)

	// 2. 获取审批流类型名称
	var flowTypeName string
	err = s.db.Table("approval_flow_type").
		Select("name").
		Where("id = ?", flowInfo.ApprovalFlowTypeID).
		Scan(&flowTypeName).Error
	if err != nil {
		return fmt.Errorf("查询审批流类型失败: %w", err)
	}

	fmt.Printf("[DEBUG] 审批流类型名称: %s\n", flowTypeName)

	// 3. 如果是"退费"类型，处理退费逻辑
	if flowTypeName == "退费" {
		fmt.Println("[DEBUG] 匹配到退费类型，开始处理退费逻辑")
		return s.handleRefundApproval(flowInfo.CreateTime, approved)
	}

	fmt.Println("[DEBUG] 非退费类型，跳过处理")
	return nil
}

// handleRefundApproval 处理退费审批完成
func (s *ApprovalFlowService) handleRefundApproval(createTime string, approved bool) error {
	fmt.Printf("[DEBUG] handleRefundApproval调用: createTime=%s, approved=%v\n", createTime, approved)

	// 1. 根据时间范围查询退费订单（前后5秒范围）
	var refundOrder struct {
		ID      int `gorm:"column:id"`
		OrderID int `gorm:"column:order_id"`
	}

	err := s.db.Table("refund_order").
		Select("id, order_id").
		Where("submit_time >= DATE_SUB(?, INTERVAL 5 SECOND)", createTime).
		Where("submit_time <= DATE_ADD(?, INTERVAL 5 SECOND)", createTime).
		Order("submit_time DESC").
		Limit(1).
		Scan(&refundOrder).Error
	if err != nil {
		return fmt.Errorf("查询退费订单失败: %w", err)
	}

	fmt.Printf("[DEBUG] 查询到退费订单: ID=%d, OrderID=%d\n", refundOrder.ID, refundOrder.OrderID)

	if refundOrder.ID == 0 {
		// 未找到关联的退费订单，可能不是退费审批流
		fmt.Println("[DEBUG] 未找到关联的退费订单，跳过处理")
		return nil
	}

	if approved {
		// 审批通过：调用退费处理逻辑
		fmt.Println("[DEBUG] 审批通过，开始处理退费逻辑")
		return s.processRefundApproval(refundOrder.ID, refundOrder.OrderID)
	} else {
		// 审批驳回：更新退费订单状态为已驳回(20)，订单状态恢复为部分支付(30)
		fmt.Println("[DEBUG] 审批驳回，开始处理驳回逻辑")
		return s.processRefundRejection(refundOrder.ID, refundOrder.OrderID)
	}
}

// processRefundApproval 处理退费审批通过
func (s *ApprovalFlowService) processRefundApproval(refundOrderID int, orderID int) error {
	fmt.Printf("[DEBUG] processRefundApproval调用: refundOrderID=%d, orderID=%d\n", refundOrderID, orderID)

	// TODO: 使用依赖注入将RefundDomainService注入到ApprovalFlowService中
	// 当前为避免循环依赖，暂时直接在这里实现退费处理逻辑
	// 完整实现需要调用 refundDomainService.ProcessRefundApproval(refundOrderID)

	return s.db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("[DEBUG] 开始事务处理")

		// 1. 查询退费订单获取学生ID
		var studentID int
		if err := tx.Table("refund_order").
			Select("student_id").
			Where("id = ?", refundOrderID).
			Scan(&studentID).Error; err != nil {
			fmt.Printf("[ERROR] 查询学生ID失败: %v\n", err)
			return err
		}
		fmt.Printf("[DEBUG] 查询到学生ID: %d\n", studentID)

		// 2. 更新退费订单状态为已通过(10)
		result := tx.Table("refund_order").
			Where("id = ?", refundOrderID).
			Update("status", 10)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新退费订单状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新退费订单状态成功，影响行数: %d\n", result.RowsAffected)

		// 2.1 更新退费子订单状态为已通过(10)
		result = tx.Table("refund_order_item").
			Where("refund_order_id = ?", refundOrderID).
			Update("status", 10)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新退费子订单状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新退费子订单状态成功，影响行数: %d\n", result.RowsAffected)

		// 2.2 更新淘宝退费补充信息状态为已通过(10)
		result = tx.Table("refund_taobao_supplement").
			Where("refund_order_id = ?", refundOrderID).
			Update("status", 10)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新淘宝退费补充信息状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新淘宝退费补充信息状态成功，影响行数: %d\n", result.RowsAffected)

		// 2.3 更新常规退费补充信息状态为已通过(10)
		result = tx.Table("refund_regular_supplement").
			Where("refund_order_id = ?", refundOrderID).
			Update("status", 10)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新常规退费补充信息状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新常规退费补充信息状态成功，影响行数: %d\n", result.RowsAffected)

		// 3. 执行退费冲回和重分账逻辑（按照Python版本实现）
		// 3.1 获取退费相关信息
		// 获取该订单的所有常规收款（已到账）
		var regularPayments []struct {
			PaymentID     int     `gorm:"column:payment_id"`
			PaymentAmount float64 `gorm:"column:payment_amount"`
			PaymentType   int     `gorm:"column:payment_type"`
		}
		if err := tx.Raw(`
			SELECT id as payment_id, payment_amount, 0 as payment_type
			FROM payment_collection
			WHERE order_id = ? AND status IN (10, 20)
			ORDER BY id ASC
		`, orderID).Scan(&regularPayments).Error; err != nil {
			return err
		}

		// 获取该订单的所有淘宝收款（已到账）
		var taobaoPayments []struct {
			PaymentID     int     `gorm:"column:payment_id"`
			PaymentAmount float64 `gorm:"column:payment_amount"`
			PaymentType   int     `gorm:"column:payment_type"`
		}
		if err := tx.Raw(`
			SELECT id as payment_id, payment_amount, 1 as payment_type
			FROM taobao_payment
			WHERE order_id = ? AND status = 30
			ORDER BY id ASC
		`, orderID).Scan(&taobaoPayments).Error; err != nil {
			return err
		}

		// 检查是否有收款信息
		hasPayment := len(regularPayments) > 0 || len(taobaoPayments) > 0
		if !hasPayment {
			// 无收款信息，跳过分账逻辑
			return nil
		}

		// 3.2 检查是否需要执行冲回
		needChargeback := false

		// 获取本次退费涉及的子订单ID列表
		var refundChildOrderIDs []int
		if err := tx.Table("refund_order_item").
			Select("DISTINCT childorder_id").
			Where("refund_order_id = ?", refundOrderID).
			Pluck("childorder_id", &refundChildOrderIDs).Error; err != nil {
			return err
		}

		// 获取收款列表区填写的退费金额（按收款维度）
		var refundPaymentsList []struct {
			PaymentID    int     `gorm:"column:payment_id"`
			PaymentType  int     `gorm:"column:payment_type"`
			RefundAmount float64 `gorm:"column:refund_amount"`
		}
		if err := tx.Table("refund_payment").
			Select("payment_id, payment_type, refund_amount").
			Where("refund_order_id = ?", refundOrderID).
			Scan(&refundPaymentsList).Error; err != nil {
			return err
		}

		// 检查每个收款的分账金额是否充足
		for _, rp := range refundPaymentsList {
			if len(refundChildOrderIDs) == 0 {
				break
			}

			// 计算该收款在被退费子订单上的分账总额（仅计算未冲回的售卖类）
			var totalSeparate float64
			if err := tx.Raw(`
				SELECT COALESCE(SUM(separate_amount), 0) as total_separate
				FROM separate_account
				WHERE payment_id = ? AND payment_type = ?
					AND childorders_id IN (?)
					AND type = 0
					AND id NOT IN (
						SELECT parent_id FROM separate_account
						WHERE payment_id = ? AND payment_type = ? AND parent_id IS NOT NULL
					)
			`, rp.PaymentID, rp.PaymentType, refundChildOrderIDs, rp.PaymentID, rp.PaymentType).
				Scan(&totalSeparate).Error; err != nil {
				return err
			}

			// 如果分账金额 < 退费金额，需要冲回
			if totalSeparate < rp.RefundAmount {
				needChargeback = true
				break
			}
		}

		// 3.3 执行冲回和重新分账
		if needChargeback {
			// 3.3a. 冲回所有未冲回的售卖类分账明细
			// 获取该订单所有未冲回的售卖类分账明细
			var originalSeparates []struct {
				ID             int     `gorm:"column:id"`
				UID            int     `gorm:"column:uid"`
				OrdersID       int     `gorm:"column:orders_id"`
				ChildOrdersID  int     `gorm:"column:childorders_id"`
				PaymentID      int     `gorm:"column:payment_id"`
				PaymentType    int     `gorm:"column:payment_type"`
				GoodsID        int     `gorm:"column:goods_id"`
				GoodsName      string  `gorm:"column:goods_name"`
				SeparateAmount float64 `gorm:"column:separate_amount"`
			}
			if err := tx.Raw(`
				SELECT id, uid, orders_id, childorders_id, payment_id, payment_type,
					   goods_id, goods_name, separate_amount
				FROM separate_account
				WHERE orders_id = ? AND type = 0
					AND id NOT IN (
						SELECT parent_id FROM separate_account
						WHERE orders_id = ? AND parent_id IS NOT NULL
					)
			`, orderID, orderID).Scan(&originalSeparates).Error; err != nil {
				return err
			}

			// 生成冲回记录（type=1，负金额，记录parent_id）
			for _, sep := range originalSeparates {
				if err := tx.Table("separate_account").Create(map[string]interface{}{
					"uid":             sep.UID,
					"orders_id":       sep.OrdersID,
					"childorders_id":  sep.ChildOrdersID,
					"payment_id":      sep.PaymentID,
					"payment_type":    sep.PaymentType,
					"goods_id":        sep.GoodsID,
					"goods_name":      sep.GoodsName,
					"separate_amount": -sep.SeparateAmount,
					"type":            1,
					"parent_id":       sep.ID,
				}).Error; err != nil {
					return err
				}
			}

			// 3.3b. 重新生成售卖类分账明细
			// 获取退费子订单列表
			var refundItemsList []struct {
				ChildOrderID int     `gorm:"column:childorder_id"`
				RefundAmount float64 `gorm:"column:refund_amount"`
				GoodsID      int     `gorm:"column:goodsid"`
				GoodsName    string  `gorm:"column:goods_name"`
			}
			if err := tx.Raw(`
				SELECT roi.childorder_id, roi.refund_amount, co.goodsid, g.name as goods_name
				FROM refund_order_item roi
				INNER JOIN childorders co ON roi.childorder_id = co.id
				LEFT JOIN goods g ON co.goodsid = g.id
				WHERE roi.refund_order_id = ?
				ORDER BY roi.childorder_id ASC
			`, refundOrderID).Scan(&refundItemsList).Error; err != nil {
				return err
			}

			// 第一批售卖分账：用退费收款分配给退费子订单
			// 为退费子订单记录剩余需求
			refundChildRemaining := make(map[int]float64)
			for _, item := range refundItemsList {
				refundChildRemaining[item.ChildOrderID] = item.RefundAmount
			}

			// 按退费收款顺序分配
			for _, rp := range refundPaymentsList {
				remainingToAllocate := rp.RefundAmount

				// 按退费子订单顺序分配
				for _, ri := range refundItemsList {
					if remainingToAllocate <= 0 {
						break
					}

					if refundChildRemaining[ri.ChildOrderID] <= 0 {
						continue
					}

					// 本次分配金额 = min(退费收款剩余, 退费子订单剩余需求)
					allocateAmount := remainingToAllocate
					if refundChildRemaining[ri.ChildOrderID] < allocateAmount {
						allocateAmount = refundChildRemaining[ri.ChildOrderID]
					}

					// 插入第一批售卖分账明细
					if err := tx.Table("separate_account").Create(map[string]interface{}{
						"uid":             studentID,
						"orders_id":       orderID,
						"childorders_id":  ri.ChildOrderID,
						"payment_id":      rp.PaymentID,
						"payment_type":    rp.PaymentType,
						"goods_id":        ri.GoodsID,
						"goods_name":      ri.GoodsName,
						"separate_amount": allocateAmount,
						"type":            0,
					}).Error; err != nil {
						return err
					}

					remainingToAllocate -= allocateAmount
					refundChildRemaining[ri.ChildOrderID] -= allocateAmount
				}
			}

			// 第二批售卖分账：用剩余收款分配给剩余子订单需求
			// 合并所有收款
			allPayments := make([]struct {
				PaymentID     int
				PaymentType   int
				PaymentAmount float64
			}, 0, len(regularPayments)+len(taobaoPayments))
			for _, p := range regularPayments {
				allPayments = append(allPayments, struct {
					PaymentID     int
					PaymentType   int
					PaymentAmount float64
				}{p.PaymentID, p.PaymentType, p.PaymentAmount})
			}
			for _, p := range taobaoPayments {
				allPayments = append(allPayments, struct {
					PaymentID     int
					PaymentType   int
					PaymentAmount float64
				}{p.PaymentID, p.PaymentType, p.PaymentAmount})
			}

			// 构建退费收款map
			refundPaymentMap := make(map[string]float64)
			for _, rp := range refundPaymentsList {
				key := fmt.Sprintf("%d_%d", rp.PaymentID, rp.PaymentType)
				refundPaymentMap[key] = rp.RefundAmount
			}

			// 计算每个收款的剩余金额 = 总金额 - 退费金额
			paymentRemaining := make(map[string]float64)
			for _, payment := range allPayments {
				key := fmt.Sprintf("%d_%d", payment.PaymentID, payment.PaymentType)
				refundAmount := refundPaymentMap[key]
				remaining := payment.PaymentAmount - refundAmount
				paymentRemaining[key] = remaining
			}

			// 获取所有子订单
			var allChildOrders []struct {
				ID             int     `gorm:"column:id"`
				GoodsID        int     `gorm:"column:goodsid"`
				AmountReceived float64 `gorm:"column:amount_received"`
				GoodsName      string  `gorm:"column:goods_name"`
			}
			if err := tx.Raw(`
				SELECT co.id, co.goodsid, co.amount_received, g.name AS goods_name
				FROM childorders co
				LEFT JOIN goods g ON co.goodsid = g.id
				WHERE co.parentsid = ?
				ORDER BY co.id ASC
			`, orderID).Scan(&allChildOrders).Error; err != nil {
				return err
			}

			// 构建退费金额map
			refundItemMap := make(map[int]float64)
			for _, item := range refundItemsList {
				refundItemMap[item.ChildOrderID] = item.RefundAmount
			}

			// 计算每个子订单的剩余需求 = 实收金额 - 退费金额
			childRemainingNeed := make(map[int]float64)
			for _, child := range allChildOrders {
				refundAmount := refundItemMap[child.ID]
				remainingNeed := child.AmountReceived - refundAmount
				childRemainingNeed[child.ID] = remainingNeed
			}

			// 按收款顺序分配剩余金额
			for _, payment := range allPayments {
				key := fmt.Sprintf("%d_%d", payment.PaymentID, payment.PaymentType)

				if paymentRemaining[key] <= 0 {
					continue
				}

				// 按子订单顺序分配
				for _, child := range allChildOrders {
					if paymentRemaining[key] <= 0 {
						break
					}

					if childRemainingNeed[child.ID] <= 0 {
						continue
					}

					// 本次分配金额 = min(收款剩余, 子订单剩余需求)
					allocateAmount := paymentRemaining[key]
					if childRemainingNeed[child.ID] < allocateAmount {
						allocateAmount = childRemainingNeed[child.ID]
					}

					// 插入第二批售卖分账明细
					if err := tx.Table("separate_account").Create(map[string]interface{}{
						"uid":             studentID,
						"orders_id":       orderID,
						"childorders_id":  child.ID,
						"payment_id":      payment.PaymentID,
						"payment_type":    payment.PaymentType,
						"goods_id":        child.GoodsID,
						"goods_name":      child.GoodsName,
						"separate_amount": allocateAmount,
						"type":            0,
					}).Error; err != nil {
						return err
					}

					paymentRemaining[key] -= allocateAmount
					childRemainingNeed[child.ID] -= allocateAmount
				}
			}
		}

		// 3.4 生成退费类分账明细（无论是否冲回都执行）
		// 重新查询所有退费子订单和退费金额（确保数据最新且避免变量冲突）
		var refundItemsForRefundSeparate []struct {
			ChildOrderID int     `gorm:"column:childorder_id"`
			RefundAmount float64 `gorm:"column:refund_amount"`
			GoodsID      int     `gorm:"column:goodsid"`
			GoodsName    string  `gorm:"column:goods_name"`
		}
		if err := tx.Raw(`
			SELECT roi.childorder_id, roi.refund_amount, co.goodsid, g.name as goods_name
			FROM refund_order_item roi
			INNER JOIN childorders co ON roi.childorder_id = co.id
			LEFT JOIN goods g ON co.goodsid = g.id
			WHERE roi.refund_order_id = ?
			ORDER BY roi.childorder_id ASC
		`, refundOrderID).Scan(&refundItemsForRefundSeparate).Error; err != nil {
			return err
		}

		// 为每个退费子订单按照其售卖分账的分布生成退费类分账明细
		for _, item := range refundItemsForRefundSeparate {
			// 查询该子订单当前的售卖分账分布（按收款ID升序）
			var childSeparates []struct {
				PaymentID      int     `gorm:"column:payment_id"`
				PaymentType    int     `gorm:"column:payment_type"`
				SeparateAmount float64 `gorm:"column:separate_amount"`
			}
			if err := tx.Raw(`
				SELECT payment_id, payment_type, separate_amount
				FROM separate_account
				WHERE childorders_id = ? AND type = 0
					AND id NOT IN (
						SELECT parent_id FROM separate_account
						WHERE childorders_id = ? AND parent_id IS NOT NULL
					)
				ORDER BY payment_id ASC
			`, item.ChildOrderID, item.ChildOrderID).Scan(&childSeparates).Error; err != nil {
				return err
			}

			// 按照售卖分账的分布生成退费类分账
			remainingRefund := item.RefundAmount
			for _, separate := range childSeparates {
				if remainingRefund <= 0 {
					break
				}

				// 本次退费金额 = min(剩余退费金额, 该收款的售卖分账金额)
				refundAmount := remainingRefund
				if separate.SeparateAmount < refundAmount {
					refundAmount = separate.SeparateAmount
				}

				// 插入退费类分账明细（负金额）
				if err := tx.Table("separate_account").Create(map[string]interface{}{
					"uid":             studentID,
					"orders_id":       orderID,
					"childorders_id":  item.ChildOrderID,
					"payment_id":      separate.PaymentID,
					"payment_type":    separate.PaymentType,
					"goods_id":        item.GoodsID,
					"goods_name":      item.GoodsName,
					"separate_amount": -refundAmount,
					"type":            2,
				}).Error; err != nil {
					return err
				}

				remainingRefund -= refundAmount
			}
		}

		// 7. 更新子订单状态（根据净分账金额计算）
		// 查询订单的所有子订单
		var allChildOrders []struct {
			ID             int     `gorm:"column:id"`
			AmountReceived float64 `gorm:"column:amount_received"`
		}
		if err := tx.Table("childorders").
			Select("id, amount_received").
			Where("parentsid = ?", orderID).
			Scan(&allChildOrders).Error; err != nil {
			return err
		}

		// 更新每个子订单的状态
		for _, childOrder := range allChildOrders {
			// 计算子订单的净分账金额（售卖类未冲回的 + 退费类）
			var netAllocated float64
			if err := tx.Raw(`
				SELECT COALESCE(SUM(separate_amount), 0) as net_allocated
				FROM separate_account
				WHERE childorders_id = ?
					AND (
						(type = 0 AND id NOT IN (
							SELECT parent_id FROM separate_account
							WHERE childorders_id = ? AND parent_id IS NOT NULL
						))
						OR type = 2
					)
			`, childOrder.ID, childOrder.ID).Scan(&netAllocated).Error; err != nil {
				return err
			}

			// 根据净分账金额确定子订单状态
			var newStatus int
			if netAllocated <= 0 {
				newStatus = orderEntity.ChildOrderStatusUnpaid // 10 未支付
			} else if netAllocated < childOrder.AmountReceived {
				newStatus = orderEntity.ChildOrderStatusPartialPaid // 20 部分支付
			} else {
				newStatus = orderEntity.ChildOrderStatusPaid // 30 已支付
			}

			if err := tx.Table("childorders").
				Where("id = ?", childOrder.ID).
				Update("status", newStatus).Error; err != nil {
				return err
			}
		}

		// 8. 更新订单状态（根据净收款计算）
		// 计算总收款金额（常规+淘宝）
		var regularPaid float64
		if err := tx.Table("payment_collection").
			Select("COALESCE(SUM(payment_amount), 0) as total").
			Where("order_id = ? AND status IN (10, 20)", orderID).
			Scan(&regularPaid).Error; err != nil {
			return err
		}

		var taobaoPaid float64
		if err := tx.Table("taobao_payment").
			Select("COALESCE(SUM(payment_amount), 0) as total").
			Where("order_id = ? AND status = 30", orderID).
			Scan(&taobaoPaid).Error; err != nil {
			return err
		}

		totalPaid := regularPaid + taobaoPaid

		// 计算总退费金额（包含当前退费订单和其他已通过的退费订单）
		var regularRefund float64
		if err := tx.Raw(`
			SELECT COALESCE(SUM(refund_amount), 0) as total
			FROM refund_regular_supplement
			WHERE refund_order_id IN (
				SELECT id FROM refund_order
				WHERE order_id = ? AND (status = 10 OR id = ?)
			)
		`, orderID, refundOrderID).Scan(&regularRefund).Error; err != nil {
			return err
		}

		var taobaoRefund float64
		if err := tx.Raw(`
			SELECT COALESCE(SUM(refund_amount), 0) as total
			FROM refund_taobao_supplement
			WHERE refund_order_id IN (
				SELECT id FROM refund_order
				WHERE order_id = ? AND (status = 10 OR id = ?)
			)
		`, orderID, refundOrderID).Scan(&taobaoRefund).Error; err != nil {
			return err
		}

		totalRefund := regularRefund + taobaoRefund

		// 净收款 = 总收款 - 总退费
		netPaid := totalPaid - totalRefund

		// 获取订单应收金额
		var amountReceived float64
		if err := tx.Table("orders").
			Select("amount_received").
			Where("id = ?", orderID).
			Scan(&amountReceived).Error; err != nil {
			return err
		}

		// 根据净收款确定订单状态
		var newOrderStatus int
		if netPaid <= 0 {
			newOrderStatus = orderEntity.OrderStatusUnpaid // 20 未支付
		} else if netPaid >= amountReceived {
			newOrderStatus = orderEntity.OrderStatusPaid // 40 已支付
		} else {
			newOrderStatus = orderEntity.OrderStatusPartialPaid // 30 部分支付
		}

		if err := tx.Table("orders").
			Where("id = ?", orderID).
			Update("status", newOrderStatus).Error; err != nil {
			return err
		}

		return nil
	})
}

// processRefundRejection 处理退费审批驳回
func (s *ApprovalFlowService) processRefundRejection(refundOrderID int, orderID int) error {
	fmt.Printf("[DEBUG] processRefundRejection调用: refundOrderID=%d, orderID=%d\n", refundOrderID, orderID)

	return s.db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("[DEBUG] 开始驳回事务处理")

		// 1. 更新退费订单状态为已驳回(20)
		result := tx.Table("refund_order").
			Where("id = ?", refundOrderID).
			Update("status", 20)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新退费订单状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新退费订单状态为已驳回，影响行数: %d\n", result.RowsAffected)

		// 1.1 更新退费子订单状态为已驳回(20)
		result = tx.Table("refund_order_item").
			Where("refund_order_id = ?", refundOrderID).
			Update("status", 20)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新退费子订单状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新退费子订单状态为已驳回，影响行数: %d\n", result.RowsAffected)

		// 1.2 更新淘宝退费补充信息状态为已驳回(20)
		result = tx.Table("refund_taobao_supplement").
			Where("refund_order_id = ?", refundOrderID).
			Update("status", 20)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新淘宝退费补充信息状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新淘宝退费补充信息状态为已驳回，影响行数: %d\n", result.RowsAffected)

		// 1.3 更新常规退费补充信息状态为已驳回(20)
		result = tx.Table("refund_regular_supplement").
			Where("refund_order_id = ?", refundOrderID).
			Update("status", 20)
		if result.Error != nil {
			fmt.Printf("[ERROR] 更新常规退费补充信息状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 更新常规退费补充信息状态为已驳回，影响行数: %d\n", result.RowsAffected)

		// 2. 恢复订单状态为部分支付(30)
		result = tx.Table("orders").
			Where("id = ?", orderID).
			Update("status", 30)
		if result.Error != nil {
			fmt.Printf("[ERROR] 恢复订单状态失败: %v\n", result.Error)
			return result.Error
		}
		fmt.Printf("[DEBUG] 恢复订单状态为部分支付，影响行数: %d\n", result.RowsAffected)

		fmt.Println("[DEBUG] processRefundRejection 完成")
		return nil
	})
}

// proceedToNextNodeOrComplete 流转到下一节点或完成审批流
func (s *ApprovalFlowService) proceedToNextNodeOrComplete(nodeCase *entity.ApprovalNodeCase, nodeCaseID int) error {
	// 1. 获取审批流信息
	flow, err := s.flowRepo.GetByID(nodeCase.ApprovalFlowManagementID)
	if err != nil {
		return err
	}

	// 2. 获取模板节点信息
	templateNode, err := s.nodeCaseRepo.GetTemplateNodeByID(nodeCase.NodeID)
	if err != nil {
		return err
	}

	// 3. 查找下一个节点
	nextNode, err := s.nodeCaseRepo.GetNextTemplateNode(flow.ApprovalFlowTemplateID, templateNode.Sort)
	if err != nil && err.Error() != "record not found" {
		return err
	}

	if nextNode != nil {
		// 有下一节点：创建下一节点实例
		approvers, err := s.nodeCaseRepo.GetNodeApprovers(nextNode.ID)
		if err != nil {
			return err
		}

		if err := s.nodeCaseRepo.CreateNextNode(
			flow.ID,
			nextNode.ID,
			nextNode.Type,
			nextNode.Sort,
			approvers,
		); err != nil {
			return err
		}

		// 增加step
		return s.flowRepo.IncrementStep(flow.ID)
	} else {
		// 没有下一节点：审批流完成
		if err := s.flowRepo.UpdateStatus(flow.ID, 10); err != nil {
			return err
		}

		// 创建抄送记录
		copyUsers, err := s.templateRepo.GetCopyUsers(flow.ApprovalFlowTemplateID)
		if err != nil {
			return err
		}

		if len(copyUsers) > 0 {
			copyInfo := "审批流已完成"
			if err := s.nodeCaseRepo.CreateCopyRecords(flow.ID, copyUsers, copyInfo); err != nil {
				return err
			}
		}

		// 审批流完成后的回调处理
		if err := s.handleApprovalComplete(flow.ID, true); err != nil {
			return err
		}

		return nil
	}
}
