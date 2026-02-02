package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

func TestDiscountService_CalculateDiscount_NoActivities(t *testing.T) {
	db, _ := setupMockDB(t)
	s := NewDiscountService(db)
	ctx := context.Background()

	goodsList := []GoodsForDiscount{
		{GoodsID: 1, Price: 100.00},
		{GoodsID: 2, Price: 200.00},
	}

	// 没有活动时应该返回零优惠
	totalDiscount, childDiscounts, err := s.CalculateDiscount(ctx, goodsList, []int{})
	if err != nil {
		t.Errorf("CalculateDiscount() error = %v, want nil", err)
	}
	if totalDiscount != 0 {
		t.Errorf("CalculateDiscount() totalDiscount = %v, want 0", totalDiscount)
	}
	if len(childDiscounts) != 2 {
		t.Errorf("CalculateDiscount() len(childDiscounts) = %v, want 2", len(childDiscounts))
	}
	for goodsID, discount := range childDiscounts {
		if discount != 0 {
			t.Errorf("CalculateDiscount() childDiscounts[%d] = %v, want 0", goodsID, discount)
		}
	}
}

func TestDiscountService_CalculateDiscount_NoGoods(t *testing.T) {
	db, _ := setupMockDB(t)
	s := NewDiscountService(db)
	ctx := context.Background()

	// 没有商品时应该返回零优惠
	totalDiscount, childDiscounts, err := s.CalculateDiscount(ctx, []GoodsForDiscount{}, []int{1})
	if err != nil {
		t.Errorf("CalculateDiscount() error = %v, want nil", err)
	}
	if totalDiscount != 0 {
		t.Errorf("CalculateDiscount() totalDiscount = %v, want 0", totalDiscount)
	}
	if len(childDiscounts) != 0 {
		t.Errorf("CalculateDiscount() len(childDiscounts) = %v, want 0", len(childDiscounts))
	}
}

func TestDiscountService_CalculateDiscount_ByGoods(t *testing.T) {
	db, mock := setupMockDB(t)
	s := NewDiscountService(db)
	ctx := context.Background()

	// 模拟查询活动信息
	activityRows := sqlmock.NewRows([]string{"id", "template_id", "type", "select_type"}).
		AddRow(100, 1, 2, 2) // type=2(满折), select_type=2(按商品)
	mock.ExpectQuery("SELECT (.+) FROM activity (.+)").
		WithArgs(100, 1).
		WillReturnRows(activityRows)

	// 模拟查询活动折扣规则
	detailRows := sqlmock.NewRows([]string{"threshold_amount", "discount_value"}).
		AddRow(2.0, 80.0) // 满2件享8折
	mock.ExpectQuery("SELECT (.+) FROM `activity_detail` (.+)").
		WithArgs(100).
		WillReturnRows(detailRows)

	// 模拟查询活动模板关联的商品
	templateGoodsRows := sqlmock.NewRows([]string{"goods_id", "classify_id"}).
		AddRow(1, nil).
		AddRow(2, nil)
	mock.ExpectQuery("SELECT (.+) FROM `activity_template_goods` (.+)").
		WithArgs(1).
		WillReturnRows(templateGoodsRows)

	goodsList := []GoodsForDiscount{
		{GoodsID: 1, Price: 100.00},
		{GoodsID: 2, Price: 200.00},
	}

	totalDiscount, childDiscounts, err := s.CalculateDiscount(ctx, goodsList, []int{100})
	if err != nil {
		t.Fatalf("CalculateDiscount() error = %v", err)
	}

	// 总价300，8折优惠20% = 60
	if totalDiscount != 60.00 {
		t.Errorf("CalculateDiscount() totalDiscount = %v, want 60.00", totalDiscount)
	}

	// 验证子订单优惠按比例分摊
	expectedChildDiscounts := map[int]float64{
		1: 20.00, // 100/300 * 60 = 20
		2: 40.00, // 200/300 * 60 = 40
	}

	for goodsID, expected := range expectedChildDiscounts {
		if childDiscounts[goodsID] != expected {
			t.Errorf("CalculateDiscount() childDiscounts[%d] = %v, want %v", goodsID, childDiscounts[goodsID], expected)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDiscountService_CalculateDiscount_ByClassify(t *testing.T) {
	db, mock := setupMockDB(t)
	s := NewDiscountService(db)
	ctx := context.Background()

	// 模拟查询活动信息
	activityRows := sqlmock.NewRows([]string{"id", "template_id", "type", "select_type"}).
		AddRow(200, 2, 2, 1) // type=2(满折), select_type=1(按分类)
	mock.ExpectQuery("SELECT (.+) FROM activity (.+)").
		WithArgs(200, 1).
		WillReturnRows(activityRows)

	// 模拟查询活动折扣规则
	detailRows := sqlmock.NewRows([]string{"threshold_amount", "discount_value"}).
		AddRow(2.0, 90.0) // 满2件享9折
	mock.ExpectQuery("SELECT (.+) FROM `activity_detail` (.+)").
		WithArgs(200).
		WillReturnRows(detailRows)

	// 模拟查询活动模板关联的分类
	templateGoodsRows := sqlmock.NewRows([]string{"goods_id", "classify_id"}).
		AddRow(nil, 10)
	mock.ExpectQuery("SELECT (.+) FROM `activity_template_goods` (.+)").
		WithArgs(2).
		WillReturnRows(templateGoodsRows)

	// 模拟查询商品分类（按分类选择时需要）
	goodsClassifyRows1 := sqlmock.NewRows([]string{"classifyid"}).AddRow(10)
	mock.ExpectQuery("SELECT (.+) FROM `goods` (.+)").
		WithArgs(1).
		WillReturnRows(goodsClassifyRows1)

	goodsClassifyRows2 := sqlmock.NewRows([]string{"classifyid"}).AddRow(10)
	mock.ExpectQuery("SELECT (.+) FROM `goods` (.+)").
		WithArgs(2).
		WillReturnRows(goodsClassifyRows2)

	goodsList := []GoodsForDiscount{
		{GoodsID: 1, Price: 100.00},
		{GoodsID: 2, Price: 200.00},
	}

	totalDiscount, childDiscounts, err := s.CalculateDiscount(ctx, goodsList, []int{200})
	if err != nil {
		t.Fatalf("CalculateDiscount() error = %v", err)
	}

	// 总价300，9折优惠10% = 30
	if totalDiscount != 30.00 {
		t.Errorf("CalculateDiscount() totalDiscount = %v, want 30.00", totalDiscount)
	}

	// 验证子订单优惠按比例分摊
	expectedChildDiscounts := map[int]float64{
		1: 10.00, // 100/300 * 30 = 10
		2: 20.00, // 200/300 * 30 = 20
	}

	for goodsID, expected := range expectedChildDiscounts {
		if childDiscounts[goodsID] != expected {
			t.Errorf("CalculateDiscount() childDiscounts[%d] = %v, want %v", goodsID, childDiscounts[goodsID], expected)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDiscountService_CalculateDiscount_NotMeetThreshold(t *testing.T) {
	db, mock := setupMockDB(t)
	s := NewDiscountService(db)
	ctx := context.Background()

	// 模拟查询活动信息
	activityRows := sqlmock.NewRows([]string{"id", "template_id", "type", "select_type"}).
		AddRow(300, 3, 2, 2)
	mock.ExpectQuery("SELECT (.+) FROM activity (.+)").
		WithArgs(300, 1).
		WillReturnRows(activityRows)

	// 模拟查询活动折扣规则：满3件享8折
	detailRows := sqlmock.NewRows([]string{"threshold_amount", "discount_value"}).
		AddRow(3.0, 80.0)
	mock.ExpectQuery("SELECT (.+) FROM `activity_detail` (.+)").
		WithArgs(300).
		WillReturnRows(detailRows)

	// 模拟查询活动模板关联的商品
	templateGoodsRows := sqlmock.NewRows([]string{"goods_id", "classify_id"}).
		AddRow(1, nil).
		AddRow(2, nil)
	mock.ExpectQuery("SELECT (.+) FROM `activity_template_goods` (.+)").
		WithArgs(3).
		WillReturnRows(templateGoodsRows)

	// 只有2件商品，不满足3件门槛
	goodsList := []GoodsForDiscount{
		{GoodsID: 1, Price: 100.00},
		{GoodsID: 2, Price: 200.00},
	}

	totalDiscount, childDiscounts, err := s.CalculateDiscount(ctx, goodsList, []int{300})
	if err != nil {
		t.Fatalf("CalculateDiscount() error = %v", err)
	}

	// 不满足门槛，应该没有优惠
	if totalDiscount != 0 {
		t.Errorf("CalculateDiscount() totalDiscount = %v, want 0", totalDiscount)
	}

	for goodsID, discount := range childDiscounts {
		if discount != 0 {
			t.Errorf("CalculateDiscount() childDiscounts[%d] = %v, want 0", goodsID, discount)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDiscountService_CalculateDiscount_NonDiscountActivity(t *testing.T) {
	db, mock := setupMockDB(t)
	s := NewDiscountService(db)
	ctx := context.Background()

	// 模拟查询活动信息：非满折活动（type != 2）
	activityRows := sqlmock.NewRows([]string{"id", "template_id", "type", "select_type"}).
		AddRow(600, 6, 1, 2) // type=1，不是满折类型
	mock.ExpectQuery("SELECT (.+) FROM activity (.+)").
		WithArgs(600, 1).
		WillReturnRows(activityRows)

	goodsList := []GoodsForDiscount{
		{GoodsID: 1, Price: 100.00},
		{GoodsID: 2, Price: 200.00},
	}

	totalDiscount, childDiscounts, err := s.CalculateDiscount(ctx, goodsList, []int{600})
	if err != nil {
		t.Fatalf("CalculateDiscount() error = %v", err)
	}

	// 非满折活动，应该没有优惠
	if totalDiscount != 0 {
		t.Errorf("CalculateDiscount() totalDiscount = %v, want 0", totalDiscount)
	}

	for goodsID, discount := range childDiscounts {
		if discount != 0 {
			t.Errorf("CalculateDiscount() childDiscounts[%d] = %v, want 0", goodsID, discount)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDiscountService_CalculateDiscount_ActivityNotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	s := NewDiscountService(db)
	ctx := context.Background()

	// 模拟查询活动信息：活动不存在
	mock.ExpectQuery("SELECT (.+) FROM activity (.+)").
		WithArgs(999, 1).
		WillReturnError(sql.ErrNoRows)

	goodsList := []GoodsForDiscount{
		{GoodsID: 1, Price: 100.00},
	}

	totalDiscount, _, err := s.CalculateDiscount(ctx, goodsList, []int{999})
	if err != nil {
		t.Fatalf("CalculateDiscount() error = %v", err)
	}

	// 活动不存在时，忽略错误继续处理，返回零优惠
	if totalDiscount != 0 {
		t.Errorf("CalculateDiscount() totalDiscount = %v, want 0", totalDiscount)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
