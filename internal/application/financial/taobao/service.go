package taobao

import (
	"charonoms/internal/domain/financial/payment"
	"charonoms/internal/domain/financial/separate"
	"charonoms/internal/domain/financial/taobao"
	orderRepo "charonoms/internal/domain/order/repository"
	"context"
	"fmt"
	"time"
	"regexp"
	"strconv"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type TaobaoPaymentService struct {
	db                 *gorm.DB
	taobaoRepo         taobao.TaobaoPaymentRepository
	paymentRepo        payment.PaymentRepository
	orderRepo          orderRepo.OrderRepository
	childOrderRepo     orderRepo.ChildOrderRepository
	separateRepo       separate.SeparateAccountRepository
}

// NewTaobaoPaymentService 创建淘宝收款服务实例
func NewTaobaoPaymentService(
	db *gorm.DB,
	taobaoRepo taobao.TaobaoPaymentRepository,
	paymentRepo payment.PaymentRepository,
	orderRepository orderRepo.OrderRepository,
	childOrderRepository orderRepo.ChildOrderRepository,
	separateRepo separate.SeparateAccountRepository,
) *TaobaoPaymentService {
	return &TaobaoPaymentService{
		db:             db,
		taobaoRepo:     taobaoRepo,
		paymentRepo:    paymentRepo,
		orderRepo:      orderRepository,
		childOrderRepo: childOrderRepository,
		separateRepo:   separateRepo,
	}
}

// GetList 获取淘宝收款列表
func (s *TaobaoPaymentService) GetList(filters map[string]interface{}) ([]*taobao.TaobaoPayment, error) {
	return s.taobaoRepo.List(filters)
}

// Create 创建淘宝收款记录
func (s *TaobaoPaymentService) Create(payment *taobao.TaobaoPayment) error {
	ctx := context.Background()
	// 验证订单存在且状态为20(未支付)或30(部分支付)
	order, err := s.orderRepo.GetOrderByID(ctx, *payment.OrderID)
	if err != nil {
		return fmt.Errorf("订单不存在")
	}
	if order.Status != 20 && order.Status != 30 {
		return fmt.Errorf("订单状态不允许添加收款")
	}

	// 计算已有收款总额
	regularPaid, err := s.paymentRepo.GetTotalPaidAmount(*payment.OrderID)
	if err != nil {
		return err
	}
	taobaoPaid, err := s.taobaoRepo.GetTotalPaid(*payment.OrderID)
	if err != nil {
		return err
	}
	totalPaid := regularPaid + taobaoPaid

	// 计算待支付金额
	pendingAmount := order.AmountReceived - totalPaid

	// 验证新增金额不超过待支付金额
	if payment.PaymentAmount > pendingAmount {
		return fmt.Errorf("收款金额超过待支付金额")
	}

	// 设置初始状态为0(已下单)
	payment.Status = taobao.TaobaoPaymentStatusOrdered

	// 创建记录
	if err := s.taobaoRepo.Create(payment); err != nil {
		return err
	}

	// 更新订单支付状态
	return s.updateOrderPaymentStatus(*payment.OrderID)
}

// ConfirmArrival 确认淘宝收款到账
func (s *TaobaoPaymentService) ConfirmArrival(paymentID int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 获取记录
		payment, err := s.taobaoRepo.GetByID(paymentID)
		if err != nil {
			return fmt.Errorf("淘宝收款记录不存在")
		}

		// 验证状态为0(已下单)
		if payment.Status != taobao.TaobaoPaymentStatusOrdered {
			return fmt.Errorf("只能确认状态为已下单的记录")
		}

		// 更新状态和到账时间
		now := time.Now()
		payment.Status = taobao.TaobaoPaymentStatusArrived
		payment.ArrivalTime = &now
		if err := s.taobaoRepo.Update(payment); err != nil {
			return err
		}

		// 生成分账明细
		if payment.OrderID != nil {
			if err := s.generateSeparateAccounts(paymentID, *payment.OrderID); err != nil {
				return err
			}
		}

		// 更新订单状态
		if payment.OrderID != nil {
			return s.updateOrderPaymentStatus(*payment.OrderID)
		}

		return nil
	})
}

// Delete 删除淘宝收款记录
func (s *TaobaoPaymentService) Delete(paymentID int) error {
	// 获取记录
	payment, err := s.taobaoRepo.GetByID(paymentID)
	if err != nil {
		return fmt.Errorf("淘宝收款记录不存在")
	}

	// 只能删除状态为0(已下单)或10(待认领)的记录
	if payment.Status != taobao.TaobaoPaymentStatusOrdered && payment.Status != taobao.TaobaoPaymentStatusUnclaimed {
		return fmt.Errorf("只能删除已下单或待认领状态的记录")
	}

	// 删除记录
	if err := s.taobaoRepo.Delete(paymentID); err != nil {
		return err
	}

	// 更新订单状态
	if payment.OrderID != nil {
		return s.updateOrderPaymentStatus(*payment.OrderID)
	}

	return nil
}

// generateSeparateAccounts 为淘宝收款生成分账明细
func (s *TaobaoPaymentService) generateSeparateAccounts(paymentID int, orderID int) error {
	// 获取淘宝收款信息
	payment, err := s.taobaoRepo.GetByID(paymentID)
	if err != nil {
		return err
	}

	// 获取该订单的所有子订单（按ID排序）
	childOrders, err := s.childOrderRepo.ListByOrderID(orderID)
	if err != nil {
		return err
	}

	// 检查是否已生成过分账
	exists, err := s.separateRepo.ExistsByPaymentAndOrder(paymentID, orderID, 1) // 1表示淘宝支付
	if err != nil {
		return err
	}
	if exists {
		return nil // 已生成过，跳过
	}

	// 按子订单顺序分配收款金额
	remainingAmount := payment.PaymentAmount
	for _, childOrder := range childOrders {
		if remainingAmount <= 0 {
			break
		}

		// 获取子订单已分配金额
		allocated, err := s.separateRepo.GetChildOrderTotalSeparate(childOrder.ID)
		if err != nil {
			return err
		}

		// 计算需要金额
		needAmount := childOrder.AmountReceived - allocated
		if needAmount <= 0 {
			continue
		}

		// 计算本次分配金额
		allocAmount := remainingAmount
		if allocAmount > needAmount {
			allocAmount = needAmount
		}

		// 创建分账记录
		separateAccount := &separate.SeparateAccount{
			UID:            *payment.StudentID,
			OrdersID:       orderID,
			ChildOrdersID:  childOrder.ID,
			PaymentID:      paymentID,
			PaymentType:    1, // 淘宝支付
			GoodsID:        childOrder.GoodsID,
			GoodsName:      "", // 商品名称（从商品表关联查询）
			SeparateAmount: allocAmount,
			Type:           0, // 售卖类
			CreateTime:     time.Now(),
		}

		if err := s.separateRepo.Create(separateAccount); err != nil {
			return err
		}

		remainingAmount -= allocAmount
	}

	// 更新所有子订单状态
	for _, childOrder := range childOrders {
		if err := s.updateChildOrderStatus(childOrder.ID); err != nil {
			return err
		}
	}

	return nil
}

// updateChildOrderStatus 更新子订单状态
func (s *TaobaoPaymentService) updateChildOrderStatus(childOrderID int) error {
	// 获取子订单
	childOrder, err := s.childOrderRepo.GetByID(childOrderID)
	if err != nil {
		return err
	}

	// 获取总分账金额
	totalSeparate, err := s.separateRepo.GetChildOrderTotalSeparate(childOrderID)
	if err != nil {
		return err
	}

	// 更新状态
	if totalSeparate <= 0 {
		childOrder.Status = 10 // 未支付
	} else if totalSeparate < childOrder.AmountReceived {
		childOrder.Status = 20 // 部分支付
	} else {
		childOrder.Status = 30 // 已支付
	}

	return s.childOrderRepo.Update(childOrder)
}

// updateOrderPaymentStatus 更新订单支付状态
func (s *TaobaoPaymentService) updateOrderPaymentStatus(orderID int) error {
	ctx := context.Background()
	// 获取订单
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 计算总收款
	regularPaid, err := s.paymentRepo.GetTotalPaidAmount(orderID)
	if err != nil {
		return err
	}
	taobaoPaid, err := s.taobaoRepo.GetTotalPaid(orderID)
	if err != nil {
		return err
	}
	totalPaid := regularPaid + taobaoPaid

	// 更新订单状态
	var newStatus int
	if totalPaid == 0 {
		newStatus = 20 // 未支付
	} else if totalPaid >= order.AmountReceived {
		newStatus = 40 // 已支付
	} else {
		newStatus = 30 // 部分支付
	}

	// 使用 UpdateOrderStatus 更新状态（不涉及子订单级联更新，传入0）
	return s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatus, 0)
}

// GetUnclaimedList 获取淘宝待认领列表
func (s *TaobaoPaymentService) GetUnclaimedList(filters map[string]interface{}) ([]*taobao.TaobaoPayment, error) {
	return s.taobaoRepo.ListUnclaimed(filters)
}

// ClaimUnclaimed 认领淘宝待认领
func (s *TaobaoPaymentService) ClaimUnclaimed(unclaimedID int, orderID int, userID int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		ctx := context.Background()
		// 获取待认领记录
		unclaimed, err := s.taobaoRepo.GetByID(unclaimedID)
		if err != nil {
			return fmt.Errorf("待认领记录不存在")
		}

		// 验证状态为10(待认领)
		if unclaimed.Status != taobao.TaobaoPaymentStatusUnclaimed {
			return fmt.Errorf("只能认领待认领状态的记录")
		}

		// 验证订单存在且状态为20(未支付)或30(部分支付)
		order, err := s.orderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("订单不存在")
		}
		if order.Status != 20 && order.Status != 30 {
			return fmt.Errorf("订单状态不允许认领")
		}

		// 计算已有收款总额 + 待认领金额
		regularPaid, err := s.paymentRepo.GetTotalPaidAmount(orderID)
		if err != nil {
			return err
		}
		taobaoPaid, err := s.taobaoRepo.GetTotalPaid(orderID)
		if err != nil {
			return err
		}
		totalPaid := regularPaid + taobaoPaid + unclaimed.PaymentAmount

		// 验证不超过订单实收金额
		if totalPaid > order.AmountReceived {
			return fmt.Errorf("认领金额超过订单实收金额")
		}

		// 更新待认领记录
		unclaimed.Status = taobao.TaobaoPaymentStatusClaimed
		unclaimed.OrderID = &orderID
		studentID := order.StudentID
		unclaimed.StudentID = &studentID
		unclaimed.Claimer = &userID
		if err := s.taobaoRepo.Update(unclaimed); err != nil {
			return err
		}

		// 生成分账明细
		if err := s.generateSeparateAccounts(unclaimedID, orderID); err != nil {
			return err
		}

		// 更新订单状态
		return s.updateOrderPaymentStatus(orderID)
	})
}

// DeleteUnclaimed 删除淘宝待认领
func (s *TaobaoPaymentService) DeleteUnclaimed(unclaimedID int) error {
	// 获取记录
	unclaimed, err := s.taobaoRepo.GetByID(unclaimedID)
	if err != nil {
		return fmt.Errorf("待认领记录不存在")
	}

	// 只能删除状态为10(待认领)的记录
	if unclaimed.Status != taobao.TaobaoPaymentStatusUnclaimed {
		return fmt.Errorf("只能删除待认领状态的记录")
	}

	return s.taobaoRepo.Delete(unclaimedID)
}

// ImportUnclaimedExcel 导入淘宝待认领Excel
func (s *TaobaoPaymentService) ImportUnclaimedExcel(records []map[string]interface{}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			payer := record["payer"].(string)
			zhifubaoAccount := record["zhifubao_account"].(string)
			amount := record["amount"].(float64)
			merchantOrder := record["merchant_order"].(string)
			arrivalTime := record["arrival_time"].(time.Time)

			// 查找匹配的记录
			matched, err := s.taobaoRepo.FindByMerchantOrderAndAmount(merchantOrder, amount)
			if err != nil {
				return err
			}

			if matched != nil {
				// 匹配成功，更新为已到账
				matched.Status = taobao.TaobaoPaymentStatusArrived
				matched.ArrivalTime = &arrivalTime
				matched.Payer = &payer
				matched.ZhifubaoAccount = &zhifubaoAccount
				if err := s.taobaoRepo.Update(matched); err != nil {
					return err
				}
			} else {
				// 未匹配，插入为待认领
				newPayment := &taobao.TaobaoPayment{
					Payer:           &payer,
					ZhifubaoAccount: &zhifubaoAccount,
					PaymentAmount:   amount,
					MerchantOrder:   &merchantOrder,
					ArrivalTime:     &arrivalTime,
					Status:          taobao.TaobaoPaymentStatusUnclaimed,
				}
				if err := s.taobaoRepo.Create(newPayment); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// GenerateUnclaimedTemplate 生成淘宝待认领Excel模板
func (s *TaobaoPaymentService) GenerateUnclaimedTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "淘宝待认领模板"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	// 设置表头
	headers := []string{"付款方", "支付宝账号", "付款金额", "商户订单号", "到账时间"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
		// 设置表头样式：加粗、居中
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// 添加示例数据
	f.SetCellValue(sheetName, "A2", "张三")
	f.SetCellValue(sheetName, "B2", "zhangsan@example.com")
	f.SetCellValue(sheetName, "C2", "1000.00")
	f.SetCellValue(sheetName, "D2", "TB202601120001")
	f.SetCellValue(sheetName, "E2", "2026-01-12")

	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	return f, nil
}

// ImportUnclaimedExcelFile 导入淘宝待认领Excel文件
// 严格按照Python版本逻辑实现
func (s *TaobaoPaymentService) ImportUnclaimedExcelFile(f *excelize.File) (int, int, []string, error) {
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("读取Excel失败: %v", err)
	}

	if len(rows) < 2 {
		return 0, 0, nil, fmt.Errorf("Excel文件为空或只有表头")
	}

	var successCount, matchedCount int
	var errorRows []string

	// 正则表达式用于验证日期格式
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 跳过表头，从第二行开始
		for idx := 1; idx < len(rows); idx++ {
			row := rows[idx]
			rowNum := idx + 1 // Excel行号从1开始，但我们跳过了表头

			// 跳过空行
			if len(row) == 0 || (len(row) > 0 && row[0] == "") {
				continue
			}

			// 确保至少有5列数据
			if len(row) < 5 {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：列数不足", rowNum))
				continue
			}

			// 提取字段
			payer := ""
			if len(row) > 0 && row[0] != "" && row[0] != "None" {
				payer = row[0]
			}

			zhifubaoAccount := ""
			if len(row) > 1 && row[1] != "" && row[1] != "None" {
				zhifubaoAccount = row[1]
			}

			// 验证付款金额
			paymentAmountStr := ""
			if len(row) > 2 {
				paymentAmountStr = row[2]
			}
			paymentAmount, err := strconv.ParseFloat(paymentAmountStr, 64)
			if err != nil || paymentAmount <= 0 {
				fmt.Printf("DEBUG 第%d行：付款金额字符串='%s'(长度=%d), err=%v, amount=%v\n", rowNum, paymentAmountStr, len(paymentAmountStr), err, paymentAmount)
			errorRows = append(errorRows, fmt.Sprintf("第%d行：付款金额格式不正确，必须为正数", rowNum))
				continue
			}
			// 保留两位小数
			paymentAmount = float64(int(paymentAmount*100+0.5)) / 100

			// 验证商户订单号（必填）
			merchantOrder := ""
			if len(row) > 3 {
				merchantOrder = row[3]
			}
			if merchantOrder == "" {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：商户订单号不能为空", rowNum))
				continue
			}

			// 验证到账时间
			arrivalTimeStr := ""
			if len(row) > 4 {
				arrivalTimeStr = row[4]
			}
			if !datePattern.MatchString(arrivalTimeStr) {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：到账时间格式不正确，必须为YYYY-MM-DD格式", rowNum))
				continue
			}
			arrivalTime, err := time.Parse("2006-01-02", arrivalTimeStr)
			if err != nil {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：到账时间格式不正确，必须为YYYY-MM-DD格式", rowNum))
				continue
			}

			// 自动匹配逻辑：查找状态为0(已下单)且金额和商户订单号匹配的记录
			matched, err := s.taobaoRepo.FindByMerchantOrderAndAmount(merchantOrder, paymentAmount)
			if err != nil {
				return err
			}

			if matched != nil {
				// 匹配成功，更新为已到账(status=30)，同时更新到账时间、付款方和支付宝账号
				matched.Status = taobao.TaobaoPaymentStatusArrived
				matched.ArrivalTime = &arrivalTime
				if payer != "" {
					matched.Payer = &payer
				}
				if zhifubaoAccount != "" {
					matched.ZhifubaoAccount = &zhifubaoAccount
				}
				if err := s.taobaoRepo.Update(matched); err != nil {
					return err
				}

				// 生成分账明细
				if matched.OrderID != nil {
					if err := s.generateSeparateAccounts(matched.ID, *matched.OrderID); err != nil {
						return err
					}
				}

				matchedCount++
			} else {
				// 未匹配，插入为待认领(status=10)
				newPayment := &taobao.TaobaoPayment{
					Status:          taobao.TaobaoPaymentStatusUnclaimed,
					PaymentAmount:   paymentAmount,
					MerchantOrder:   &merchantOrder,
					ArrivalTime:     &arrivalTime,
				}
				if payer != "" {
					newPayment.Payer = &payer
				}
				if zhifubaoAccount != "" {
					newPayment.ZhifubaoAccount = &zhifubaoAccount
				}
				if err := s.taobaoRepo.Create(newPayment); err != nil {
					return err
				}
			}

			successCount++
		}
		return nil
	})

	return successCount, matchedCount, errorRows, err
}

