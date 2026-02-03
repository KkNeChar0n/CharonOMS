package payment

import (
	"errors"

	"charonoms/internal/application/financial"
	domainPayment "charonoms/internal/domain/financial/payment"
	domainSeparate "charonoms/internal/domain/financial/separate"
	studentRepo "charonoms/internal/domain/student/repository"
	"gorm.io/gorm"
)

// PaymentApplicationService 收款应用服务
type PaymentApplicationService struct {
	db                   *gorm.DB
	paymentRepo          domainPayment.PaymentRepository
	studentRepo          studentRepo.StudentRepository
	paymentDomainService *domainPayment.PaymentDomainService
	separateDomainService *domainSeparate.SeparateAccountDomainService
}

// NewPaymentApplicationService 创建收款应用服务
func NewPaymentApplicationService(
	db *gorm.DB,
	paymentRepo domainPayment.PaymentRepository,
	studentRepo studentRepo.StudentRepository,
	paymentDomainService *domainPayment.PaymentDomainService,
	separateDomainService *domainSeparate.SeparateAccountDomainService,
) *PaymentApplicationService {
	return &PaymentApplicationService{
		db:                    db,
		paymentRepo:           paymentRepo,
		studentRepo:           studentRepo,
		paymentDomainService:  paymentDomainService,
		separateDomainService: separateDomainService,
	}
}

// GetPaymentCollections 获取收款列表
func (s *PaymentApplicationService) GetPaymentCollections(
	id, studentID, orderID *int,
	payer *string,
	paymentMethod *int,
	tradingDate *string,
	status *int,
	page, pageSize int,
) (*financial.PaymentCollectionListResponse, error) {
	// 默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 查询收款列表
	filter := domainPayment.PaymentListFilter{
		ID:            id,
		StudentID:     studentID,
		OrderID:       orderID,
		Payer:         payer,
		PaymentMethod: paymentMethod,
		TradingDate:   tradingDate,
		Status:        status,
		Page:          page,
		PageSize:      pageSize,
	}

	payments, total, err := s.paymentRepo.List(filter)
	if err != nil {
		return nil, err
	}

	// 查询学生姓名
	studentNames := make(map[int]string)
	for _, p := range payments {
		student, err := s.studentRepo.GetStudentByID(p.StudentID)
		if err == nil && student != nil {
			studentNames[p.StudentID] = student.Name
		}
	}

	// 转换为DTO
	dtos := ToPaymentCollectionDTOList(payments, studentNames)

	return &financial.PaymentCollectionListResponse{
		Collections: dtos,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

// CreatePaymentCollection 新增收款
func (s *PaymentApplicationService) CreatePaymentCollection(req *financial.CreatePaymentCollectionRequest) (int, error) {
	// 验证付款金额
	err := s.paymentDomainService.ValidatePaymentAmount(req.OrderID, req.PaymentAmount)
	if err != nil {
		return 0, err
	}

	// 转换为实体
	paymentEntity := ToPaymentCollectionEntity(req)

	// 在事务中执行
	var paymentID int
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 插入收款记录
		err := s.paymentRepo.Create(paymentEntity)
		if err != nil {
			return err
		}
		paymentID = paymentEntity.ID

		// 更新订单状态
		err = s.paymentDomainService.UpdateOrderPaymentStatus(req.OrderID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return paymentID, nil
}

// ConfirmPaymentCollection 确认收款到账
func (s *PaymentApplicationService) ConfirmPaymentCollection(id int) error {
	// 查询收款记录
	paymentEntity, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return err
	}
	if paymentEntity == nil {
		return errors.New("收款记录不存在")
	}

	// 检查是否可以确认
	if !paymentEntity.CanConfirm() {
		return errors.New("只能确认未核验的收款")
	}

	// 在事务中执行
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新收款状态
		paymentEntity.Confirm()
		err := s.paymentRepo.Update(paymentEntity)
		if err != nil {
			return err
		}

		// 更新订单状态
		err = s.paymentDomainService.UpdateOrderPaymentStatus(paymentEntity.OrderID)
		if err != nil {
			return err
		}

		// 生成分账明细
		err = s.separateDomainService.GenerateSeparateAccounts(id, paymentEntity.OrderID)
		if err != nil {
			return err
		}

		return nil
	})
}

// DeletePaymentCollection 删除收款
func (s *PaymentApplicationService) DeletePaymentCollection(id int) error {
	// 查询收款记录
	paymentEntity, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return err
	}
	if paymentEntity == nil {
		return errors.New("收款记录不存在")
	}

	// 检查是否可以删除
	if !paymentEntity.CanDelete() {
		return errors.New("只能删除未核验的收款")
	}

	// 在事务中执行
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除收款记录
		err := s.paymentRepo.Delete(id)
		if err != nil {
			return err
		}

		// 更新订单状态
		err = s.paymentDomainService.UpdateOrderPaymentStatus(paymentEntity.OrderID)
		if err != nil {
			return err
		}

		return nil
	})
}
