package service

import (
	"testing"
)

func TestOrderService_CalculateOrderAmounts(t *testing.T) {
	s := NewOrderService()

	tests := []struct {
		name                    string
		goodsList               []GoodsItem
		discountAmount          float64
		wantAmountReceivable    float64
		wantAmountReceived      float64
	}{
		{
			name: "单个商品无优惠",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			},
			discountAmount:       0,
			wantAmountReceivable: 100.00,
			wantAmountReceived:   100.00,
		},
		{
			name: "单个商品有优惠",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			},
			discountAmount:       10.00,
			wantAmountReceivable: 100.00,
			wantAmountReceived:   90.00,
		},
		{
			name: "多个商品无优惠",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
				{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
				{GoodsID: 3, TotalPrice: 300.00, Price: 300.00},
			},
			discountAmount:       0,
			wantAmountReceivable: 600.00,
			wantAmountReceived:   600.00,
		},
		{
			name: "多个商品有优惠",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
				{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
				{GoodsID: 3, TotalPrice: 300.00, Price: 300.00},
			},
			discountAmount:       50.00,
			wantAmountReceivable: 600.00,
			wantAmountReceived:   550.00,
		},
		{
			name: "组合商品：TotalPrice与Price不同",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 150.00, Price: 100.00}, // 组合商品
				{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
			},
			discountAmount:       20.00,
			wantAmountReceivable: 350.00,
			wantAmountReceived:   280.00, // (100 + 200) - 20
		},
		{
			name:                 "空商品列表",
			goodsList:            []GoodsItem{},
			discountAmount:       0,
			wantAmountReceivable: 0,
			wantAmountReceived:   0,
		},
		{
			name: "精度测试：需要四舍五入",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 33.333, Price: 33.333},
				{GoodsID: 2, TotalPrice: 66.666, Price: 66.666},
			},
			discountAmount:       9.999,
			wantAmountReceivable: 100.00, // 33.333 + 66.666 = 99.999 -> 100.00
			wantAmountReceived:   90.00,  // 99.999 - 9.999 = 90.00
		},
		{
			name: "全额优惠",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			},
			discountAmount:       100.00,
			wantAmountReceivable: 100.00,
			wantAmountReceived:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAmountReceivable, gotAmountReceived := s.CalculateOrderAmounts(tt.goodsList, tt.discountAmount)
			if gotAmountReceivable != tt.wantAmountReceivable {
				t.Errorf("CalculateOrderAmounts() amountReceivable = %v, want %v", gotAmountReceivable, tt.wantAmountReceivable)
			}
			if gotAmountReceived != tt.wantAmountReceived {
				t.Errorf("CalculateOrderAmounts() amountReceived = %v, want %v", gotAmountReceived, tt.wantAmountReceived)
			}
		})
	}
}

func TestOrderService_AllocateChildDiscounts(t *testing.T) {
	s := NewOrderService()

	tests := []struct {
		name           string
		goodsList      []GoodsItem
		childDiscounts map[int]float64
		want           map[int]float64
	}{
		{
			name: "单个商品分摊",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			},
			childDiscounts: map[int]float64{
				1: 10.00,
			},
			want: map[int]float64{
				1: 10.00,
			},
		},
		{
			name: "多个商品分摊",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
				{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
				{GoodsID: 3, TotalPrice: 300.00, Price: 300.00},
			},
			childDiscounts: map[int]float64{
				1: 5.00,
				2: 10.00,
				3: 15.00,
			},
			want: map[int]float64{
				1: 5.00,
				2: 10.00,
				3: 15.00,
			},
		},
		{
			name: "部分商品无优惠",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
				{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
			},
			childDiscounts: map[int]float64{
				1: 10.00,
				// GoodsID 2 没有优惠
			},
			want: map[int]float64{
				1: 10.00,
				2: 0, // 默认为0
			},
		},
		{
			name: "精度四舍五入",
			goodsList: []GoodsItem{
				{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			},
			childDiscounts: map[int]float64{
				1: 10.555, // 应该四舍五入到10.56
			},
			want: map[int]float64{
				1: 10.56,
			},
		},
		{
			name:           "空商品列表",
			goodsList:      []GoodsItem{},
			childDiscounts: map[int]float64{},
			want:           map[int]float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.AllocateChildDiscounts(tt.goodsList, tt.childDiscounts)
			if len(got) != len(tt.want) {
				t.Errorf("AllocateChildDiscounts() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for goodsID, wantDiscount := range tt.want {
				gotDiscount, exists := got[goodsID]
				if !exists {
					t.Errorf("AllocateChildDiscounts() missing goodsID %d", goodsID)
					continue
				}
				if gotDiscount != wantDiscount {
					t.Errorf("AllocateChildDiscounts() goodsID %d = %v, want %v", goodsID, gotDiscount, wantDiscount)
				}
			}
		})
	}
}

func TestOrderService_CalculateChildAmounts(t *testing.T) {
	s := NewOrderService()

	tests := []struct {
		name                 string
		totalPrice           float64
		price                float64
		discountAmount       float64
		wantAmountReceivable float64
		wantAmountReceived   float64
	}{
		{
			name:                 "无优惠",
			totalPrice:           100.00,
			price:                100.00,
			discountAmount:       0,
			wantAmountReceivable: 100.00,
			wantAmountReceived:   100.00,
		},
		{
			name:                 "有优惠",
			totalPrice:           100.00,
			price:                100.00,
			discountAmount:       10.00,
			wantAmountReceivable: 100.00,
			wantAmountReceived:   90.00,
		},
		{
			name:                 "组合商品：TotalPrice与Price不同",
			totalPrice:           150.00,
			price:                100.00,
			discountAmount:       10.00,
			wantAmountReceivable: 150.00,
			wantAmountReceived:   90.00,
		},
		{
			name:                 "全额优惠",
			totalPrice:           100.00,
			price:                100.00,
			discountAmount:       100.00,
			wantAmountReceivable: 100.00,
			wantAmountReceived:   0,
		},
		{
			name:                 "精度测试",
			totalPrice:           99.999,
			price:                99.999,
			discountAmount:       9.999,
			wantAmountReceivable: 100.00,
			wantAmountReceived:   90.00,
		},
		{
			name:                 "零金额",
			totalPrice:           0,
			price:                0,
			discountAmount:       0,
			wantAmountReceivable: 0,
			wantAmountReceived:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAmountReceivable, gotAmountReceived := s.CalculateChildAmounts(tt.totalPrice, tt.price, tt.discountAmount)
			if gotAmountReceivable != tt.wantAmountReceivable {
				t.Errorf("CalculateChildAmounts() amountReceivable = %v, want %v", gotAmountReceivable, tt.wantAmountReceivable)
			}
			if gotAmountReceived != tt.wantAmountReceived {
				t.Errorf("CalculateChildAmounts() amountReceived = %v, want %v", gotAmountReceived, tt.wantAmountReceived)
			}
		})
	}
}

func TestRoundToTwoDecimal(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{"整数", 100.00, 100.00},
		{"一位小数", 100.5, 100.5},
		{"两位小数", 100.55, 100.55},
		{"三位小数向上舍入", 100.555, 100.56},
		{"三位小数向下舍入", 100.554, 100.55},
		{"多位小数", 100.556789, 100.56},
		{"负数向上舍入", -100.555, -100.56},
		{"负数向下舍入", -100.554, -100.55},
		{"零", 0, 0},
		{"极小值", 0.001, 0.00},
		{"极小值舍入", 0.005, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundToTwoDecimal(tt.value); got != tt.want {
				t.Errorf("roundToTwoDecimal(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
