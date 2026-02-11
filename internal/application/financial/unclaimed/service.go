package unclaimed

import (
	"charonoms/internal/domain/financial/payment"
	"charonoms/internal/domain/financial/separate"
	"charonoms/internal/domain/financial/unclaimed"
	orderRepo "charonoms/internal/domain/order/repository"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type UnclaimedService struct {
	db             *gorm.DB
	unclaimedRepo  unclaimed.UnclaimedRepository
	paymentRepo    payment.PaymentRepository
	orderRepo      orderRepo.OrderRepository
	separateService *separate.SeparateAccountDomainService
}

// NewUnclaimedService 创建常规待认领服务实例
func NewUnclaimedService(
	db *gorm.DB,
	unclaimedRepo unclaimed.UnclaimedRepository,
	paymentRepo payment.PaymentRepository,
	orderRepository orderRepo.OrderRepository,
	separateService *separate.SeparateAccountDomainService,
) *UnclaimedService {
	return &UnclaimedService{
		db:              db,
		unclaimedRepo:   unclaimedRepo,
		paymentRepo:     paymentRepo,
		orderRepo:       orderRepository,
		separateService: separateService,
	}
}

// GetList 获取待认领列表
func (s *UnclaimedService) GetList(filters map[string]interface{}) ([]*unclaimed.Unclaimed, error) {
	return s.unclaimedRepo.List(filters)
}

// Claim 认领待认领款项
func (s *UnclaimedService) Claim(unclaimedID int, orderID int, userID int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		ctx := context.Background()

		// 获取待认领记录
		unclaimedRecord, err := s.unclaimedRepo.GetByID(unclaimedID)
		if err != nil {
			return fmt.Errorf("待认领记录不存在")
		}

		// 验证状态为待认领
		if unclaimedRecord.Status != unclaimed.UnclaimedStatusPending {
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

		// 计算已有收款总额
		totalPaid, err := s.paymentRepo.GetTotalPaidAmount(orderID)
		if err != nil {
			return err
		}

		// 验证认领金额不超过订单实收金额
		if totalPaid+unclaimedRecord.PaymentAmount > order.AmountReceived {
			return fmt.Errorf("收款金额超出订单实收金额")
		}

		// 创建 payment_collection 记录（线下场景，已支付状态）
		paymentCollection := &payment.PaymentCollection{
			OrderID:         orderID,
			StudentID:       order.StudentID,
			PaymentScenario: payment.PaymentScenarioOffline, // 线下场景
			PaymentMethod:   unclaimedRecord.PaymentMethod,
			PaymentAmount:   unclaimedRecord.PaymentAmount,
			Payer:           *unclaimedRecord.Payer,
			PayeeEntity:     unclaimedRecord.PayeeEntity,
			TradingHours:    unclaimedRecord.ArrivalTime,
			ArrivalTime:     unclaimedRecord.ArrivalTime,
			MerchantOrder:   *unclaimedRecord.MerchantOrder,
			Status:          payment.PaymentStatusPaid, // 已支付
		}

		if err := s.paymentRepo.Create(paymentCollection); err != nil {
			return err
		}

		// 生成分账明细
		if err := s.separateService.GenerateSeparateAccounts(paymentCollection.ID, orderID); err != nil {
			return err
		}

		// 更新待认领记录为已认领状态
		unclaimedRecord.Status = unclaimed.UnclaimedStatusClaimed
		unclaimedRecord.Claimer = &userID
		unclaimedRecord.PaymentID = &paymentCollection.ID
		if err := s.unclaimedRepo.Update(unclaimedRecord); err != nil {
			return err
		}

		// 更新订单支付状态
		return s.updateOrderPaymentStatus(orderID)
	})
}

// Delete 删除待认领记录
func (s *UnclaimedService) Delete(unclaimedID int) error {
	// 获取记录
	unclaimedRecord, err := s.unclaimedRepo.GetByID(unclaimedID)
	if err != nil {
		return fmt.Errorf("待认领记录不存在")
	}

	// 只能删除待认领状态的记录
	if unclaimedRecord.Status != unclaimed.UnclaimedStatusPending {
		return fmt.Errorf("只能删除待认领状态的记录")
	}

	return s.unclaimedRepo.Delete(unclaimedID)
}

// ImportExcelRecords 导入Excel数据
func (s *UnclaimedService) ImportExcelRecords(records []map[string]interface{}) (int, int, []string, error) {
	successCount := 0
	matchedCount := 0
	var errors []string

	err := s.db.Transaction(func(tx *gorm.DB) error {
		for i, record := range records {
			rowNum := i + 2 // Excel从第2行开始（第1行是标题）

			// 提取字段
			paymentMethod, ok := record["payment_method"].(int)
			if !ok {
				errors = append(errors, fmt.Sprintf("第%d行：付款方式格式错误", rowNum))
				continue
			}

			paymentAmount, ok := record["payment_amount"].(float64)
			if !ok {
				errors = append(errors, fmt.Sprintf("第%d行：付款金额格式错误", rowNum))
				continue
			}

			payer, _ := record["payer"].(string)
			payeeEntity, ok := record["payee_entity"].(int)
			if !ok {
				errors = append(errors, fmt.Sprintf("第%d行：收款主体格式错误", rowNum))
				continue
			}

			merchantOrder, _ := record["merchant_order"].(string)
			arrivalTime, ok := record["arrival_time"].(time.Time)
			if !ok {
				errors = append(errors, fmt.Sprintf("第%d行：到账时间格式错误", rowNum))
				continue
			}

			// 尝试自动匹配已有的未核验收款记录
			if merchantOrder != "" {
				paymentID, err := s.unclaimedRepo.FindMatchingPayment(merchantOrder, paymentMethod, paymentAmount, payeeEntity)
				if err != nil {
					return err
				}

				if paymentID > 0 {
					// 找到匹配，更新为已支付状态
					if err := s.unclaimedRepo.UpdatePaymentStatus(paymentID, &arrivalTime); err != nil {
						return err
					}
					matchedCount++
					successCount++
					continue
				}
			}

			// 未匹配到，插入为待认领记录
			newUnclaimed := &unclaimed.Unclaimed{
				PaymentMethod: paymentMethod,
				PaymentAmount: paymentAmount,
				Payer:         &payer,
				PayeeEntity:   payeeEntity,
				MerchantOrder: &merchantOrder,
				ArrivalTime:   &arrivalTime,
				Status:        unclaimed.UnclaimedStatusPending,
			}

			if err := s.unclaimedRepo.Create(newUnclaimed); err != nil {
				return err
			}
			successCount++
		}
		return nil
	})

	return successCount, matchedCount, errors, err
}

// updateOrderPaymentStatus 更新订单支付状态
func (s *UnclaimedService) updateOrderPaymentStatus(orderID int) error {
	ctx := context.Background()

	// 获取订单
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 计算总收款
	totalPaid, err := s.paymentRepo.GetTotalPaidAmount(orderID)
	if err != nil {
		return err
	}

	// 更新订单状态
	var newStatus int
	if totalPaid == 0 {
		newStatus = 20 // 未支付
	} else if totalPaid >= order.AmountReceived {
		newStatus = 40 // 已支付
	} else {
		newStatus = 30 // 部分支付
	}

	// 使用 UpdateOrderStatus 更新状态
	return s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatus, 0)
}

// GenerateTemplate 生成常规待认领Excel模板
func (s *UnclaimedService) GenerateTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "待认领收款模板"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	// 设置表头（严格按照Python版本）
	headers := []string{"付款方式", "付款金额", "付款方", "收款主体", "商户订单号", "到账时间"}
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

	// 添加示例数据（严格按照Python版本）
	f.SetCellValue(sheetName, "A2", "微信")
	f.SetCellValue(sheetName, "B2", "1000.00")
	f.SetCellValue(sheetName, "C2", "张三")
	f.SetCellValue(sheetName, "D2", "北京")
	f.SetCellValue(sheetName, "E2", "M202601120001")
	f.SetCellValue(sheetName, "F2", "2026-01-12")

	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	return f, nil
}

// ImportExcelFile 导入常规待认领Excel文件
// 严格按照Python版本逻辑实现
func (s *UnclaimedService) ImportExcelFile(f *excelize.File) (int, int, []string, error) {
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("读取Excel失败: %v", err)
	}

	if len(rows) < 2 {
		return 0, 0, nil, fmt.Errorf("Excel文件为空或只有表头")
	}

	// 付款方式映射（严格按照Python版本）
	paymentMethodMap := map[string]int{
		"微信":     0,
		"支付宝":    1,
		"优利支付":   2,
		"零零购支付":  3,
		"对公转账":   9,
	}

	// 收款主体映射（严格按照Python版本）
	payeeEntityMap := map[string]int{
		"北京": 0,
		"西安": 1,
	}

	var successCount, matchedCount int
	var errorRows []string

	// 正则表达式用于验证日期格式
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 跳过表头，从第二行开始
		for idx := 1; idx < len(rows); idx++ {
			row := rows[idx]
			rowNum := idx + 1

			// 跳过空行
			if len(row) == 0 || (len(row) > 0 && row[0] == "") {
				continue
			}

			// 确保至少有6列数据
			if len(row) < 6 {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：列数不足", rowNum))
				continue
			}

			// 验证付款方式
			paymentMethodStr := ""
			if len(row) > 0 {
				paymentMethodStr = row[0]
			}
			paymentMethod, ok := paymentMethodMap[paymentMethodStr]
			if !ok {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：付款方式不正确，仅支持：微信、支付宝、对公转账、零零购支付、优利支付", rowNum))
				continue
			}

			// 验证付款金额
			paymentAmountStr := ""
			if len(row) > 1 {
				paymentAmountStr = row[1]
			}
			paymentAmount, err := strconv.ParseFloat(paymentAmountStr, 64)
			if err != nil || paymentAmount <= 0 {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：付款金额格式不正确，必须为正数", rowNum))
				continue
			}
			// 保留两位小数
			paymentAmount = float64(int(paymentAmount*100+0.5)) / 100

			// 付款方
			payer := ""
			if len(row) > 2 {
				payer = row[2]
			}

			// 验证收款主体
			payeeEntityStr := ""
			if len(row) > 3 {
				payeeEntityStr = row[3]
			}
			payeeEntity, ok := payeeEntityMap[payeeEntityStr]
			if !ok {
				errorRows = append(errorRows, fmt.Sprintf("第%d行：收款主体不正确，仅支持：北京、西安", rowNum))
				continue
			}

			// 商户订单号（可选）
			merchantOrder := ""
			if len(row) > 4 && row[4] != "" && row[4] != "None" {
				merchantOrder = row[4]
			}

			// 验证到账时间
			arrivalTimeStr := ""
			if len(row) > 5 {
				arrivalTimeStr = row[5]
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

			// 自动匹配逻辑：仅当商户订单号不为空时进行匹配
			matched := false
			if merchantOrder != "" {
				// 查找payment_collection表中符合条件的记录
				paymentID, err := s.unclaimedRepo.FindMatchingPayment(merchantOrder, paymentMethod, paymentAmount, payeeEntity)
				if err != nil {
					return err
				}

				if paymentID > 0 {
					// 匹配成功，更新payment_collection状态为已支付(status=20)，更新到账时间
					if err := s.unclaimedRepo.UpdatePaymentStatus(paymentID, &arrivalTime); err != nil {
						return err
					}

					// 获取payment信息以获取orderID
					paymentInfo, err := s.paymentRepo.GetByID(paymentID)
					if err != nil {
						return err
					}

					// 生成分账明细
					if err := s.separateService.GenerateSeparateAccounts(paymentID, paymentInfo.OrderID); err != nil {
						return err
					}

					matched = true
					matchedCount++
				}
			}

			// 如果没有匹配到，则插入到unclaimed表(status=0待认领)
			if !matched {
				newUnclaimed := &unclaimed.Unclaimed{
					PaymentMethod: paymentMethod,
					PaymentAmount: paymentAmount,
					Payer:         &payer,
					PayeeEntity:   payeeEntity,
					MerchantOrder: &merchantOrder,
					ArrivalTime:   &arrivalTime,
					Status:        unclaimed.UnclaimedStatusPending, // 0-待认领
				}
				if err := s.unclaimedRepo.Create(newUnclaimed); err != nil {
					return err
				}
			}

			successCount++
		}
		return nil
	})

	return successCount, matchedCount, errorRows, err
}
