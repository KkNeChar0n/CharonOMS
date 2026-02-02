package entity

import "testing"

func TestChildOrder_ValidateAmounts(t *testing.T) {
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
			c := &ChildOrder{
				AmountReceivable: tt.amountReceivable,
				AmountReceived:   tt.amountReceived,
				DiscountAmount:   tt.discountAmount,
			}
			if got := c.ValidateAmounts(); got != tt.want {
				t.Errorf("ChildOrder.ValidateAmounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChildOrder_StatusConstants(t *testing.T) {
	// 测试状态常量值
	if ChildOrderStatusInit != 0 {
		t.Errorf("ChildOrderStatusInit = %d, want 0", ChildOrderStatusInit)
	}
	if ChildOrderStatusUnpaid != 10 {
		t.Errorf("ChildOrderStatusUnpaid = %d, want 10", ChildOrderStatusUnpaid)
	}
	if ChildOrderStatusPartialPaid != 20 {
		t.Errorf("ChildOrderStatusPartialPaid = %d, want 20", ChildOrderStatusPartialPaid)
	}
	if ChildOrderStatusPaid != 30 {
		t.Errorf("ChildOrderStatusPaid = %d, want 30", ChildOrderStatusPaid)
	}
	if ChildOrderStatusCancelled != 99 {
		t.Errorf("ChildOrderStatusCancelled = %d, want 99", ChildOrderStatusCancelled)
	}
}

func TestChildOrder_TableName(t *testing.T) {
	c := ChildOrder{}
	if got := c.TableName(); got != "childorders" {
		t.Errorf("ChildOrder.TableName() = %v, want childorders", got)
	}
}

func TestChildOrder_Creation(t *testing.T) {
	childOrder := &ChildOrder{
		ID:               1,
		ParentsID:        100,
		GoodsID:          200,
		AmountReceivable: 500.00,
		AmountReceived:   450.00,
		DiscountAmount:   50.00,
		Status:           ChildOrderStatusInit,
	}

	if childOrder.ID != 1 {
		t.Errorf("ChildOrder.ID = %d, want 1", childOrder.ID)
	}
	if childOrder.ParentsID != 100 {
		t.Errorf("ChildOrder.ParentsID = %d, want 100", childOrder.ParentsID)
	}
	if childOrder.GoodsID != 200 {
		t.Errorf("ChildOrder.GoodsID = %d, want 200", childOrder.GoodsID)
	}
	if childOrder.AmountReceivable != 500.00 {
		t.Errorf("ChildOrder.AmountReceivable = %f, want 500.00", childOrder.AmountReceivable)
	}
	if childOrder.AmountReceived != 450.00 {
		t.Errorf("ChildOrder.AmountReceived = %f, want 450.00", childOrder.AmountReceived)
	}
	if childOrder.DiscountAmount != 50.00 {
		t.Errorf("ChildOrder.DiscountAmount = %f, want 50.00", childOrder.DiscountAmount)
	}
	if !childOrder.ValidateAmounts() {
		t.Error("ChildOrder.ValidateAmounts() should return true for valid amounts")
	}
}
