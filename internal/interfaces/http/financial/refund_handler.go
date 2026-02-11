package financial

import (
	"bytes"
	"charonoms/internal/application/financial/refund"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RefundHandler struct {
	service *refund.RefundService
}

// NewRefundHandler 创建退费处理器
func NewRefundHandler(service *refund.RefundService) *RefundHandler {
	return &RefundHandler{service: service}
}

// GetRefundOrders 获取退费订单列表
// GET /api/refund_orders
func (h *RefundHandler) GetRefundOrders(c *gin.Context) {
	filters := make(map[string]interface{})

	// 解析查询参数
	if id := c.Query("id"); id != "" {
		filters["id"] = id
	}
	if uid := c.Query("uid"); uid != "" {
		filters["uid"] = uid
	}
	if orderID := c.Query("order_id"); orderID != "" {
		filters["order_id"] = orderID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	list, err := h.service.GetRefundOrders(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"refund_orders": list})
}

// GetRefundOrderDetail 获取退费订单详情
// GET /api/refund_orders/:id
func (h *RefundHandler) GetRefundOrderDetail(c *gin.Context) {
	// 获取退费订单ID
	refundOrderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的退费订单ID",
		})
		return
	}

	detail, err := h.service.GetRefundOrderDetail(refundOrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, detail)
}

// GetRefundChildOrders 获取退费子订单列表
// GET /api/refund_childorders
func (h *RefundHandler) GetRefundChildOrders(c *gin.Context) {
	filters := make(map[string]interface{})

	// 解析查询参数
	if id := c.Query("id"); id != "" {
		filters["id"] = id
	}
	if studentID := c.Query("student_id"); studentID != "" {
		filters["student_id"] = studentID
	}
	if orderID := c.Query("order_id"); orderID != "" {
		filters["order_id"] = orderID
	}
	if refundOrderID := c.Query("refund_order_id"); refundOrderID != "" {
		filters["refund_order_id"] = refundOrderID
	}
	if childOrderID := c.Query("childorder_id"); childOrderID != "" {
		filters["childorder_id"] = childOrderID
	}
	if goodsID := c.Query("goods_id"); goodsID != "" {
		filters["goods_id"] = goodsID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	list, err := h.service.GetRefundChildOrders(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"refund_childorders": list})
}

// GetRefundRegularSupplements 获取常规退费补充信息列表
// GET /api/refund-regular-supplements
func (h *RefundHandler) GetRefundRegularSupplements(c *gin.Context) {
	filters := make(map[string]interface{})

	// 解析查询参数
	if id := c.Query("id"); id != "" {
		filters["id"] = id
	}
	if refundOrderID := c.Query("refund_order_id"); refundOrderID != "" {
		filters["refund_order_id"] = refundOrderID
	}
	// 支持 uid 或 student_id 参数
	if uid := c.Query("uid"); uid != "" {
		filters["uid"] = uid
	} else if studentID := c.Query("student_id"); studentID != "" {
		filters["uid"] = studentID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	list, err := h.service.GetRefundRegularSupplements(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"regular_supplements": list})
}

// GetRefundTaobaoSupplements 获取淘宝退费补充信息列表
// GET /api/refund-taobao-supplements
func (h *RefundHandler) GetRefundTaobaoSupplements(c *gin.Context) {
	filters := make(map[string]interface{})

	// 解析查询参数
	if id := c.Query("id"); id != "" {
		filters["id"] = id
	}
	if refundOrderID := c.Query("refund_order_id"); refundOrderID != "" {
		filters["refund_order_id"] = refundOrderID
	}
	// 支持 uid 或 student_id 参数
	if uid := c.Query("uid"); uid != "" {
		filters["uid"] = uid
	} else if studentID := c.Query("student_id"); studentID != "" {
		filters["uid"] = studentID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	list, err := h.service.GetRefundTaobaoSupplements(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"taobao_supplements": list})
}

// GetRefundPaymentDetails 获取退费支付明细列表
// GET /api/refund-payment-details
func (h *RefundHandler) GetRefundPaymentDetails(c *gin.Context) {
	filters := make(map[string]interface{})

	// 解析查询参数
	if id := c.Query("id"); id != "" {
		filters["id"] = id
	}
	if studentID := c.Query("student_id"); studentID != "" {
		filters["student_id"] = studentID
	}
	if orderID := c.Query("order_id"); orderID != "" {
		filters["order_id"] = orderID
	}
	if refundOrderID := c.Query("refund_order_id"); refundOrderID != "" {
		filters["refund_order_id"] = refundOrderID
	}
	if paymentID := c.Query("payment_id"); paymentID != "" {
		filters["payment_id"] = paymentID
	}
	if paymentType := c.Query("payment_type"); paymentType != "" {
		filters["payment_type"] = paymentType
	}

	list, err := h.service.GetRefundPaymentDetails(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"refund_payment_details": list})
}

// CreateRefundOrder 创建退费订单
// POST /api/refund-orders
func (h *RefundHandler) CreateRefundOrder(c *gin.Context) {
	// 读取并记录原始请求体
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	println("收到退费申请请求，原始Body:", string(bodyBytes))

	// 重新设置请求体以便后续绑定
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req refund.CreateRefundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		println("JSON绑定失败:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	println("JSON绑定成功，订单ID:", req.OrderID)

	// 从JWT中获取用户信息
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户ID"})
		return
	}

	// 创建退费订单
	refundOrderID, err := h.service.CreateRefundOrder(
		c.Request.Context(),
		&req,
		username.(string),
		int(userID.(uint)),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"refund_order_id": refundOrderID,
		"message":         "退费申请提交成功",
	})
}
