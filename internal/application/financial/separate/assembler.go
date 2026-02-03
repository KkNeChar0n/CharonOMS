package separate

import (
	"charonoms/internal/application/financial"
	domainSeparate "charonoms/internal/domain/financial/separate"
)

// ToSeparateAccountDTO 实体转DTO
func ToSeparateAccountDTO(s *domainSeparate.SeparateAccount) *financial.SeparateAccountDTO {
	if s == nil {
		return nil
	}
	return &financial.SeparateAccountDTO{
		ID:             s.ID,
		UID:            s.UID,
		OrdersID:       s.OrdersID,
		ChildOrdersID:  s.ChildOrdersID,
		PaymentID:      s.PaymentID,
		PaymentType:    s.PaymentType,
		GoodsID:        s.GoodsID,
		GoodsName:      s.GoodsName,
		SeparateAmount: s.SeparateAmount,
		Type:           s.Type,
		CreateTime:     s.CreateTime,
	}
}

// ToSeparateAccountDTOList 实体列表转DTO列表
func ToSeparateAccountDTOList(accounts []*domainSeparate.SeparateAccount) []*financial.SeparateAccountDTO {
	result := make([]*financial.SeparateAccountDTO, 0, len(accounts))
	for _, s := range accounts {
		result = append(result, ToSeparateAccountDTO(s))
	}
	return result
}
