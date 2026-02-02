package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"charonoms/internal/application/order"
	orderDTO "charonoms/internal/interfaces/http/order"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	service *order.Service
}

// NewOrderHandler 创建订单处理器实例
func NewOrderHandler(service *order.Service) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// GetOrders 获取订单列表
func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.service.GetOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保返回空数组而不是null
	if orders == nil {
		orders = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req orderDTO.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为应用层DTO
	var expectedPaymentTime *time.Time
	if req.ExpectedPaymentTime != nil && !req.ExpectedPaymentTime.IsZero() {
		t := req.ExpectedPaymentTime.Time
		expectedPaymentTime = &t
	}

	appReq := &order.CreateOrderRequest{
		StudentID:           req.StudentID,
		GoodsList:           make([]order.GoodsItemRequest, len(req.GoodsList)),
		ExpectedPaymentTime: expectedPaymentTime,
		ActivityIDs:         req.ActivityIDs,
		DiscountAmount:      req.DiscountAmount,
		ChildDiscounts:      req.ChildDiscounts,
	}

	for i, g := range req.GoodsList {
		appReq.GoodsList[i] = order.GoodsItemRequest{
			GoodsID:    g.GoodsID,
			TotalPrice: g.TotalPrice,
			Price:      g.Price,
		}
	}

	// 初始化空切片和map以避免nil
	if appReq.ActivityIDs == nil {
		appReq.ActivityIDs = []int{}
	}
	if appReq.ChildDiscounts == nil {
		appReq.ChildDiscounts = make(map[int]float64)
	}

	orderID, err := h.service.CreateOrder(c.Request.Context(), appReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "订单创建成功",
		"order_id": orderID,
	})
}

// GetOrderGoods 获取订单商品列表
func (h *OrderHandler) GetOrderGoods(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	goods, err := h.service.GetOrderGoods(c.Request.Context(), orderID)
	if err != nil {
		// 输出详细错误日志
		fmt.Printf("GetOrderGoods error for orderID %d: %v\n", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保返回空数组而不是null
	if goods == nil {
		goods = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"goods": goods,
	})
}

// UpdateOrder 更新订单
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req orderDTO.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为应用层DTO
	var expectedPaymentTime *time.Time
	if req.ExpectedPaymentTime != nil && !req.ExpectedPaymentTime.IsZero() {
		t := req.ExpectedPaymentTime.Time
		expectedPaymentTime = &t
	}

	appReq := &order.UpdateOrderRequest{
		GoodsList:           make([]order.GoodsItemRequest, len(req.GoodsList)),
		ExpectedPaymentTime: expectedPaymentTime,
		ActivityIDs:         req.ActivityIDs,
		DiscountAmount:      req.DiscountAmount,
		ChildDiscounts:      req.ChildDiscounts,
	}

	for i, g := range req.GoodsList {
		appReq.GoodsList[i] = order.GoodsItemRequest{
			GoodsID:    g.GoodsID,
			TotalPrice: g.TotalPrice,
			Price:      g.Price,
		}
	}

	// 初始化空切片和map以避免nil
	if appReq.ActivityIDs == nil {
		appReq.ActivityIDs = []int{}
	}
	if appReq.ChildDiscounts == nil {
		appReq.ChildDiscounts = make(map[int]float64)
	}

	err = h.service.UpdateOrder(c.Request.Context(), orderID, appReq)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "只能编辑草稿状态的订单" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("UpdateOrder error for orderID %d: %v\n", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Order %d updated successfully\n", orderID)
	c.JSON(http.StatusOK, gin.H{
		"message":  "订单更新成功",
		"order_id": orderID,
	})
}

// SubmitOrder 提交订单
func (h *OrderHandler) SubmitOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.SubmitOrder(c.Request.Context(), orderID)
	if err != nil {
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "只能提交草稿状态的订单" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "订单提交成功",
	})
}

// CancelOrder 作废订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.CancelOrder(c.Request.Context(), orderID)
	if err != nil {
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "只能作废草稿状态的订单" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "订单已作废",
	})
}

// GetChildOrders 获取子订单列表
func (h *OrderHandler) GetChildOrders(c *gin.Context) {
	childOrders, err := h.service.GetChildOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保返回空数组而不是null
	if childOrders == nil {
		childOrders = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"childorders": childOrders,
	})
}

// GetActiveGoodsForOrder 获取启用商品列表（用于订单）
func (h *OrderHandler) GetActiveGoodsForOrder(c *gin.Context) {
	goods, err := h.service.GetActiveGoodsForOrder(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保返回空数组而不是null
	if goods == nil {
		goods = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"goods": goods,
	})
}

// GetGoodsTotalPrice 获取商品总价
func (h *OrderHandler) GetGoodsTotalPrice(c *gin.Context) {
	goodsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.service.GetGoodsTotalPrice(c.Request.Context(), goodsID)
	if err != nil {
		if err.Error() == "goods not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CalculateOrderDiscount 计算订单优惠
func (h *OrderHandler) CalculateOrderDiscount(c *gin.Context) {
	var req orderDTO.CalculateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换为应用层DTO
	goodsList := make([]order.GoodsItemRequest, len(req.GoodsList))
	for i, g := range req.GoodsList {
		goodsList[i] = order.GoodsItemRequest{
			GoodsID:    g.GoodsID,
			TotalPrice: g.TotalPrice,
			Price:      g.Price,
		}
	}

	activityIDs := req.ActivityIDs
	if activityIDs == nil {
		activityIDs = []int{}
	}

	// 调用服务计算优惠
	totalDiscount, childDiscounts, err := h.service.CalculateOrderDiscount(c.Request.Context(), goodsList, activityIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_discount":  totalDiscount,
		"child_discounts": childDiscounts,
	})
}
