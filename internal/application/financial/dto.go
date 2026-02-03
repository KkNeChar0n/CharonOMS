package financial

import "time"

// 此文件定义财务管理模块的DTO（数据传输对象）

// PaymentCollectionDTO 收款记录响应DTO
type PaymentCollectionDTO struct {
	ID              int        `json:"id"`
	OrderID         int        `json:"order_id"`
	StudentID       int        `json:"student_id"`
	StudentName     string     `json:"student_name,omitempty"`
	PaymentScenario int        `json:"payment_scenario"`
	PaymentMethod   int        `json:"payment_method"`
	PaymentAmount   float64    `json:"payment_amount"`
	Payer           string     `json:"payer"`
	PayeeEntity     int        `json:"payee_entity"`
	TradingHours    *time.Time `json:"trading_hours"`
	ArrivalTime     *time.Time `json:"arrival_time"`
	MerchantOrder   string     `json:"merchant_order"`
	Status          int        `json:"status"`
	CreateTime      time.Time  `json:"create_time"`
}

// CreatePaymentCollectionRequest 新增收款请求DTO
type CreatePaymentCollectionRequest struct {
	OrderID         int        `json:"order_id" binding:"required"`
	StudentID       int        `json:"student_id" binding:"required"`
	PaymentScenario int        `json:"payment_scenario" binding:"required"`
	PaymentMethod   int        `json:"payment_method" binding:"required"`
	PaymentAmount   float64    `json:"payment_amount" binding:"required,gt=0"`
	Payer           string     `json:"payer"`
	PayeeEntity     int        `json:"payee_entity" binding:"required"`
	MerchantOrder   string     `json:"merchant_order"`
	TradingHours    *time.Time `json:"trading_hours"`
}

// PaymentCollectionListResponse 收款列表响应DTO
type PaymentCollectionListResponse struct {
	Collections []*PaymentCollectionDTO `json:"collections"`
	Total       int64                   `json:"total"`
	Page        int                     `json:"page"`
	PageSize    int                     `json:"page_size"`
}

// SeparateAccountDTO 分账明细响应DTO
type SeparateAccountDTO struct {
	ID             int       `json:"id"`
	UID            int       `json:"uid"`
	OrdersID       int       `json:"orders_id"`
	ChildOrdersID  int       `json:"childorders_id"`
	PaymentID      int       `json:"payment_id"`
	PaymentType    int       `json:"payment_type"`
	GoodsID        int       `json:"goods_id"`
	GoodsName      string    `json:"goods_name"`
	SeparateAmount float64   `json:"separate_amount"`
	Type           int       `json:"type"`
	CreateTime     time.Time `json:"create_time"`
}

// SeparateAccountListResponse 分账明细列表响应DTO
type SeparateAccountListResponse struct {
	SeparateAccounts []*SeparateAccountDTO `json:"separate_accounts"`
	Total            int64                 `json:"total"`
	Page             int                   `json:"page"`
	PageSize         int                   `json:"page_size"`
}
