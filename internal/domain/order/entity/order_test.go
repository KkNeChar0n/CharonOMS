package entity

import (
	"testing"
	"time"
)

func TestOrder_CanEdit(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"草稿状态可以编辑", OrderStatusDraft, true},
		{"未支付状态不能编辑", OrderStatusUnpaid, false},
		{"部分支付状态不能编辑", OrderStatusPartialPaid, false},
		{"已支付状态不能编辑", OrderStatusPaid, false},
		{"退费中状态不能编辑", OrderStatusRefunding, false},
		{"已作废状态不能编辑", OrderStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{Status: tt.status}
			if got := o.CanEdit(); got != tt.want {
				t.Errorf("Order.CanEdit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_CanSubmit(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"草稿状态可以提交", OrderStatusDraft, true},
		{"未支付状态不能提交", OrderStatusUnpaid, false},
		{"部分支付状态不能提交", OrderStatusPartialPaid, false},
		{"已支付状态不能提交", OrderStatusPaid, false},
		{"退费中状态不能提交", OrderStatusRefunding, false},
		{"已作废状态不能提交", OrderStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{Status: tt.status}
			if got := o.CanSubmit(); got != tt.want {
				t.Errorf("Order.CanSubmit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_CanCancel(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"草稿状态可以作废", OrderStatusDraft, true},
		{"未支付状态不能作废", OrderStatusUnpaid, false},
		{"部分支付状态不能作废", OrderStatusPartialPaid, false},
		{"已支付状态不能作废", OrderStatusPaid, false},
		{"退费中状态不能作废", OrderStatusRefunding, false},
		{"已作废状态不能作废", OrderStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{Status: tt.status}
			if got := o.CanCancel(); got != tt.want {
				t.Errorf("Order.CanCancel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_ValidateAmounts(t *testing.T) {
	tests := []struct {
		name              string
		amountReceivable  float64
		amountReceived    float64
		discountAmount    float64
		want              bool
	}{
		{"正常金额", 100.00, 90.00, 10.00, true},
		{"应收金额为零", 0, 0, 0, true},
		{"应收金额为负", -100.00, 90.00, 10.00, false},
		{"实收金额为负", 100.00, -90.00, 10.00, false},
		{"优惠金额为负", 100.00, 90.00, -10.00, false},
		{"实收大于应收", 100.00, 110.00, 0, false},
		{"边界情况：实收等于应收", 100.00, 100.00, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				AmountReceivable: tt.amountReceivable,
				AmountReceived:   tt.amountReceived,
				DiscountAmount:   tt.discountAmount,
			}
			if got := o.ValidateAmounts(); got != tt.want {
				t.Errorf("Order.ValidateAmounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_StatusConstants(t *testing.T) {
	// 测试状态常量值
	if OrderStatusDraft != 10 {
		t.Errorf("OrderStatusDraft = %d, want 10", OrderStatusDraft)
	}
	if OrderStatusUnpaid != 20 {
		t.Errorf("OrderStatusUnpaid = %d, want 20", OrderStatusUnpaid)
	}
	if OrderStatusPartialPaid != 30 {
		t.Errorf("OrderStatusPartialPaid = %d, want 30", OrderStatusPartialPaid)
	}
	if OrderStatusPaid != 40 {
		t.Errorf("OrderStatusPaid = %d, want 40", OrderStatusPaid)
	}
	if OrderStatusRefunding != 50 {
		t.Errorf("OrderStatusRefunding = %d, want 50", OrderStatusRefunding)
	}
	if OrderStatusCancelled != 99 {
		t.Errorf("OrderStatusCancelled = %d, want 99", OrderStatusCancelled)
	}
}

func TestOrder_TableName(t *testing.T) {
	o := Order{}
	if got := o.TableName(); got != "orders" {
		t.Errorf("Order.TableName() = %v, want orders", got)
	}
}

func TestOrder_Creation(t *testing.T) {
	now := time.Now()
	expectedPaymentTime := now.Add(24 * time.Hour)

	order := &Order{
		ID:                  1,
		StudentID:           100,
		ExpectedPaymentTime: &expectedPaymentTime,
		AmountReceivable:    1000.00,
		AmountReceived:      900.00,
		DiscountAmount:      100.00,
		Status:              OrderStatusDraft,
		CreateTime:          now,
	}

	if order.ID != 1 {
		t.Errorf("Order.ID = %d, want 1", order.ID)
	}
	if order.StudentID != 100 {
		t.Errorf("Order.StudentID = %d, want 100", order.StudentID)
	}
	if order.AmountReceivable != 1000.00 {
		t.Errorf("Order.AmountReceivable = %f, want 1000.00", order.AmountReceivable)
	}
	if order.AmountReceived != 900.00 {
		t.Errorf("Order.AmountReceived = %f, want 900.00", order.AmountReceived)
	}
	if order.DiscountAmount != 100.00 {
		t.Errorf("Order.DiscountAmount = %f, want 100.00", order.DiscountAmount)
	}
	if !order.ValidateAmounts() {
		t.Error("Order.ValidateAmounts() should return true for valid amounts")
	}
}
