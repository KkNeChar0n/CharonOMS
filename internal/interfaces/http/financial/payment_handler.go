package financial

import (
	"net/http"
	"strconv"

	"charonoms/internal/application/financial"
	paymentApp "charonoms/internal/application/financial/payment"
	"github.com/gin-gonic/gin"
)

// PaymentHandler 收款接口处理器
type PaymentHandler struct {
	paymentService *paymentApp.PaymentApplicationService
}

// NewPaymentHandler 创建收款接口处理器
func NewPaymentHandler(paymentService *paymentApp.PaymentApplicationService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// GetPaymentCollections 获取收款列表
// GET /api/payment-collections
func (h *PaymentHandler) GetPaymentCollections(c *gin.Context) {
	// 解析查询参数
	var id, studentID, orderID, paymentMethod, status *int
	var payer, tradingDate *string
	var page, pageSize int = 1, 20

	if idStr := c.Query("id"); idStr != "" {
		if idVal, err := strconv.Atoi(idStr); err == nil {
			id = &idVal
		}
	}
	if studentIDStr := c.Query("student_id"); studentIDStr != "" {
		if studentIDVal, err := strconv.Atoi(studentIDStr); err == nil {
			studentID = &studentIDVal
		}
	}
	if orderIDStr := c.Query("order_id"); orderIDStr != "" {
		if orderIDVal, err := strconv.Atoi(orderIDStr); err == nil {
			orderID = &orderIDVal
		}
	}
	if payerStr := c.Query("payer"); payerStr != "" {
		payer = &payerStr
	}
	if paymentMethodStr := c.Query("payment_method"); paymentMethodStr != "" {
		if paymentMethodVal, err := strconv.Atoi(paymentMethodStr); err == nil {
			paymentMethod = &paymentMethodVal
		}
	}
	if tradingDateStr := c.Query("trading_date"); tradingDateStr != "" {
		tradingDate = &tradingDateStr
	}
	if statusStr := c.Query("status"); statusStr != "" {
		if statusVal, err := strconv.Atoi(statusStr); err == nil {
			status = &statusVal
		}
	}
	if pageStr := c.Query("page"); pageStr != "" {
		if pageVal, err := strconv.Atoi(pageStr); err == nil && pageVal > 0 {
			page = pageVal
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeVal > 0 {
			pageSize = pageSizeVal
		}
	}

	// 调用应用服务
	response, err := h.paymentService.GetPaymentCollections(
		id, studentID, orderID, payer, paymentMethod, tradingDate, status, page, pageSize,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    response,
	})
}

// CreatePaymentCollection 新增收款
// POST /api/payment-collections
func (h *PaymentHandler) CreatePaymentCollection(c *gin.Context) {
	var req financial.CreatePaymentCollectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用应用服务
	paymentID, err := h.paymentService.CreatePaymentCollection(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "收款记录创建成功",
		"data": gin.H{
			"id": paymentID,
		},
	})
}

// ConfirmPaymentCollection 确认收款到账
// PUT /api/payment-collections/:id/confirm
func (h *PaymentHandler) ConfirmPaymentCollection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的收款ID",
			"data":    nil,
		})
		return
	}

	// 调用应用服务
	err = h.paymentService.ConfirmPaymentCollection(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "确认到账成功",
		"data":    nil,
	})
}

// DeletePaymentCollection 删除收款
// DELETE /api/payment-collections/:id
func (h *PaymentHandler) DeletePaymentCollection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的收款ID",
			"data":    nil,
		})
		return
	}

	// 调用应用服务
	err = h.paymentService.DeletePaymentCollection(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}
