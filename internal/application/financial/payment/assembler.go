package payment

import (
	"charonoms/internal/application/financial"
	domainPayment "charonoms/internal/domain/financial/payment"
)

// ToPaymentCollectionDTO 实体转DTO
func ToPaymentCollectionDTO(p *domainPayment.PaymentCollection, studentName string) *financial.PaymentCollectionDTO {
	if p == nil {
		return nil
	}
	return &financial.PaymentCollectionDTO{
		ID:              p.ID,
		OrderID:         p.OrderID,
		StudentID:       p.StudentID,
		StudentName:     studentName,
		PaymentScenario: p.PaymentScenario,
		PaymentMethod:   p.PaymentMethod,
		PaymentAmount:   p.PaymentAmount,
		Payer:           p.Payer,
		PayeeEntity:     p.PayeeEntity,
		TradingHours:    p.TradingHours,
		ArrivalTime:     p.ArrivalTime,
		MerchantOrder:   p.MerchantOrder,
		Status:          p.Status,
		CreateTime:      p.CreateTime,
	}
}

// ToPaymentCollectionDTOList 实体列表转DTO列表
func ToPaymentCollectionDTOList(payments []*domainPayment.PaymentCollection, studentNames map[int]string) []*financial.PaymentCollectionDTO {
	result := make([]*financial.PaymentCollectionDTO, 0, len(payments))
	for _, p := range payments {
		studentName := studentNames[p.StudentID]
		result = append(result, ToPaymentCollectionDTO(p, studentName))
	}
	return result
}

// ToPaymentCollectionEntity 请求转实体
func ToPaymentCollectionEntity(req *financial.CreatePaymentCollectionRequest) *domainPayment.PaymentCollection {
	return &domainPayment.PaymentCollection{
		OrderID:         req.OrderID,
		StudentID:       req.StudentID,
		PaymentScenario: req.PaymentScenario,
		PaymentMethod:   req.PaymentMethod,
		PaymentAmount:   req.PaymentAmount,
		Payer:           req.Payer,
		PayeeEntity:     req.PayeeEntity,
		MerchantOrder:   req.MerchantOrder,
		TradingHours:    req.TradingHours,
		Status:          domainPayment.PaymentStatusUnverified, // 默认未核验
	}
}
