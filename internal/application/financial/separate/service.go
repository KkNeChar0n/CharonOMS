package separate

import (
	"charonoms/internal/application/financial"
	domainSeparate "charonoms/internal/domain/financial/separate"
)

// SeparateAccountApplicationService 分账明细应用服务
type SeparateAccountApplicationService struct {
	separateRepo domainSeparate.SeparateAccountRepository
}

// NewSeparateAccountApplicationService 创建分账明细应用服务
func NewSeparateAccountApplicationService(
	separateRepo domainSeparate.SeparateAccountRepository,
) *SeparateAccountApplicationService {
	return &SeparateAccountApplicationService{
		separateRepo: separateRepo,
	}
}

// GetSeparateAccounts 获取分账明细列表
func (s *SeparateAccountApplicationService) GetSeparateAccounts(
	id, uid, ordersID, childOrdersID, goodsID, paymentID, paymentType, separateType *int,
	page, pageSize int,
) (*financial.SeparateAccountListResponse, error) {
	// 默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 查询分账明细列表
	filter := domainSeparate.SeparateListFilter{
		ID:            id,
		UID:           uid,
		OrdersID:      ordersID,
		ChildOrdersID: childOrdersID,
		GoodsID:       goodsID,
		PaymentID:     paymentID,
		PaymentType:   paymentType,
		Type:          separateType,
		Page:          page,
		PageSize:      pageSize,
	}

	accounts, total, err := s.separateRepo.List(filter)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	dtos := ToSeparateAccountDTOList(accounts)

	return &financial.SeparateAccountListResponse{
		SeparateAccounts: dtos,
		Total:            total,
		Page:             page,
		PageSize:         pageSize,
	}, nil
}
