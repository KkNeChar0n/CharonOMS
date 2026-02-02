package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	orderApp "charonoms/internal/application/order"
	goodsRepo "charonoms/internal/infrastructure/persistence/mysql/goods"
	orderPersistence "charonoms/internal/infrastructure/persistence/order"
	"charonoms/internal/interfaces/http/handler"
	orderDTO "charonoms/internal/interfaces/http/order"
)

var (
	testDB     *gorm.DB
	testRouter *gin.Engine
)

// setupTestEnv 设置测试环境
func setupTestEnv(t *testing.T) {
	// 连接测试数据库（请根据实际情况修改连接字符串）
	dsn := "root:123456@tcp(localhost:3306)/zhixinstudent_dev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	testDB = db

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化仓储和服务
	orderRepository := orderPersistence.NewOrderRepository(testDB)
	childOrderRepository := orderPersistence.NewChildOrderRepository(testDB)
	goodsRepository := goodsRepo.NewGoodsRepository(testDB)

	orderService := orderApp.NewService(
		orderRepository,
		childOrderRepository,
		goodsRepository,
		testDB,
	)

	// 创建处理器
	orderHandler := handler.NewOrderHandler(orderService)

	// 设置路由
	testRouter = gin.New()
	api := testRouter.Group("/api")
	{
		orders := api.Group("/orders")
		{
			orders.GET("", orderHandler.GetOrders)
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id/goods", orderHandler.GetOrderGoods)
			orders.PUT("/:id", orderHandler.UpdateOrder)
			orders.PUT("/:id/submit", orderHandler.SubmitOrder)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			orders.POST("/calculate-discount", orderHandler.CalculateOrderDiscount)
		}

		api.GET("/childorders", orderHandler.GetChildOrders)
		api.GET("/goods/active-for-order", orderHandler.GetActiveGoodsForOrder)
		api.GET("/goods/:id/total-price", orderHandler.GetGoodsTotalPrice)
	}
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T, orderID int) {
	if orderID > 0 {
		// 删除订单相关数据
		testDB.Exec("DELETE FROM orders_activity WHERE orderid = ?", orderID)
		testDB.Exec("DELETE FROM childorders WHERE parentsid = ?", orderID)
		testDB.Exec("DELETE FROM orders WHERE id = ?", orderID)
	}
}

func TestOrderAPI_CreateOrder(t *testing.T) {
	setupTestEnv(t)

	// 准备测试数据
	expectedPaymentTime := time.Now().Add(24 * time.Hour)
	reqBody := orderDTO.CreateOrderRequest{
		StudentID:           1, // 假设学生ID为1存在
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	assert.NotNil(t, result["order_id"])
	orderID := int(result["order_id"].(float64))

	// 清理测试数据
	defer cleanupTestData(t, orderID)

	// 验证订单是否创建成功
	var orderCount int64
	testDB.Table("orders").Where("id = ?", orderID).Count(&orderCount)
	assert.Equal(t, int64(1), orderCount)

	// 验证子订单是否创建成功
	var childOrderCount int64
	testDB.Table("childorders").Where("parentsid = ?", orderID).Count(&childOrderCount)
	assert.Equal(t, int64(2), childOrderCount)
}

func TestOrderAPI_GetOrders(t *testing.T) {
	setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	assert.NotNil(t, result["orders"])
}

func TestOrderAPI_UpdateOrder(t *testing.T) {
	setupTestEnv(t)

	// 先创建一个订单
	expectedPaymentTime := time.Now().Add(24 * time.Hour)
	createReq := orderDTO.CreateOrderRequest{
		StudentID: 1,
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	var createResult map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &createResult)
	orderID := int(createResult["order_id"].(float64))

	defer cleanupTestData(t, orderID)

	// 更新订单
	newExpectedPaymentTime := time.Now().Add(48 * time.Hour)
	updateReq := orderDTO.UpdateOrderRequest{
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
		},
		ExpectedPaymentTime: &newExpectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      10.00,
		ChildDiscounts: map[int]float64{
			1: 5.00,
			2: 5.00,
		},
	}

	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%d", orderID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// 验证更新后的子订单数量
	var childOrderCount int64
	testDB.Table("childorders").Where("parentsid = ?", orderID).Count(&childOrderCount)
	assert.Equal(t, int64(2), childOrderCount)
}

func TestOrderAPI_SubmitOrder(t *testing.T) {
	setupTestEnv(t)

	// 先创建一个订单
	expectedPaymentTime := time.Now().Add(24 * time.Hour)
	createReq := orderDTO.CreateOrderRequest{
		StudentID: 1,
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	var createResult map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &createResult)
	orderID := int(createResult["order_id"].(float64))

	defer cleanupTestData(t, orderID)

	// 提交订单
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%d/submit", orderID), nil)
	resp = httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// 验证订单状态是否为未支付（20）
	var orderStatus int
	testDB.Table("orders").Select("status").Where("id = ?", orderID).Scan(&orderStatus)
	assert.Equal(t, 20, orderStatus)
}

func TestOrderAPI_CancelOrder(t *testing.T) {
	setupTestEnv(t)

	// 先创建一个订单
	expectedPaymentTime := time.Now().Add(24 * time.Hour)
	createReq := orderDTO.CreateOrderRequest{
		StudentID: 1,
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	var createResult map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &createResult)
	orderID := int(createResult["order_id"].(float64))

	defer cleanupTestData(t, orderID)

	// 作废订单
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%d/cancel", orderID), nil)
	resp = httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// 验证订单状态是否为已作废（99）
	var orderStatus int
	testDB.Table("orders").Select("status").Where("id = ?", orderID).Scan(&orderStatus)
	assert.Equal(t, 99, orderStatus)
}

func TestOrderAPI_GetOrderGoods(t *testing.T) {
	setupTestEnv(t)

	// 先创建一个订单
	expectedPaymentTime := time.Now().Add(24 * time.Hour)
	createReq := orderDTO.CreateOrderRequest{
		StudentID: 1,
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	var createResult map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &createResult)
	orderID := int(createResult["order_id"].(float64))

	defer cleanupTestData(t, orderID)

	// 获取订单商品列表
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/orders/%d/goods", orderID), nil)
	resp = httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	assert.NotNil(t, result["goods"])
	goodsList := result["goods"].([]interface{})
	assert.Equal(t, 2, len(goodsList))
}

func TestOrderAPI_GetChildOrders(t *testing.T) {
	setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/childorders", nil)
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	assert.NotNil(t, result["childorders"])
}

func TestOrderAPI_GetActiveGoodsForOrder(t *testing.T) {
	setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/goods/active-for-order", nil)
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	assert.NotNil(t, result["goods"])
}

func TestOrderAPI_CalculateDiscount(t *testing.T) {
	setupTestEnv(t)

	reqBody := orderDTO.CalculateDiscountRequest{
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
		},
		ActivityIDs: []int{},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/orders/calculate-discount", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	assert.NotNil(t, result["total_discount"])
	assert.NotNil(t, result["child_discounts"])
}

func TestOrderAPI_SubmitNonDraftOrder(t *testing.T) {
	setupTestEnv(t)

	// 先创建并提交一个订单
	expectedPaymentTime := time.Now().Add(24 * time.Hour)
	createReq := orderDTO.CreateOrderRequest{
		StudentID: 1,
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	var createResult map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &createResult)
	orderID := int(createResult["order_id"].(float64))

	defer cleanupTestData(t, orderID)

	// 提交订单
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%d/submit", orderID), nil)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// 再次尝试提交，应该失败
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%d/submit", orderID), nil)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Contains(t, result["error"], "只能提交草稿状态的订单")
}

func TestOrderAPI_UpdateNonDraftOrder(t *testing.T) {
	setupTestEnv(t)

	// 先创建并提交一个订单
	expectedPaymentTime := time.Now().Add(24 * time.Hour)
	createReq := orderDTO.CreateOrderRequest{
		StudentID: 1,
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	var createResult map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &createResult)
	orderID := int(createResult["order_id"].(float64))

	defer cleanupTestData(t, orderID)

	// 提交订单
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%d/submit", orderID), nil)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	// 尝试更新已提交的订单，应该失败
	updateReq := orderDTO.UpdateOrderRequest{
		GoodsList: []orderDTO.GoodsItemRequest{
			{GoodsID: 1, TotalPrice: 100.00, Price: 100.00},
			{GoodsID: 2, TotalPrice: 200.00, Price: 200.00},
		},
		ExpectedPaymentTime: &expectedPaymentTime,
		ActivityIDs:         []int{},
		DiscountAmount:      0,
		ChildDiscounts:      map[int]float64{},
	}

	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%d", orderID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Contains(t, result["error"], "只能编辑草稿状态的订单")
}
