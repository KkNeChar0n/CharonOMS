package financial

import (
	"net/http"
	"strconv"

	separateApp "charonoms/internal/application/financial/separate"
	"github.com/gin-gonic/gin"
)

// SeparateAccountHandler 分账明细接口处理器
type SeparateAccountHandler struct {
	separateService *separateApp.SeparateAccountApplicationService
}

// NewSeparateAccountHandler 创建分账明细接口处理器
func NewSeparateAccountHandler(separateService *separateApp.SeparateAccountApplicationService) *SeparateAccountHandler {
	return &SeparateAccountHandler{
		separateService: separateService,
	}
}

// GetSeparateAccounts 获取分账明细列表
// GET /api/separate-accounts
func (h *SeparateAccountHandler) GetSeparateAccounts(c *gin.Context) {
	// 解析查询参数
	var id, uid, ordersID, childOrdersID, goodsID, paymentID, paymentType, separateType *int
	var page, pageSize int = 1, 20

	if idStr := c.Query("id"); idStr != "" {
		if idVal, err := strconv.Atoi(idStr); err == nil {
			id = &idVal
		}
	}
	if uidStr := c.Query("uid"); uidStr != "" {
		if uidVal, err := strconv.Atoi(uidStr); err == nil {
			uid = &uidVal
		}
	}
	if ordersIDStr := c.Query("orders_id"); ordersIDStr != "" {
		if ordersIDVal, err := strconv.Atoi(ordersIDStr); err == nil {
			ordersID = &ordersIDVal
		}
	}
	if childOrdersIDStr := c.Query("childorders_id"); childOrdersIDStr != "" {
		if childOrdersIDVal, err := strconv.Atoi(childOrdersIDStr); err == nil {
			childOrdersID = &childOrdersIDVal
		}
	}
	if goodsIDStr := c.Query("goods_id"); goodsIDStr != "" {
		if goodsIDVal, err := strconv.Atoi(goodsIDStr); err == nil {
			goodsID = &goodsIDVal
		}
	}
	if paymentIDStr := c.Query("payment_id"); paymentIDStr != "" {
		if paymentIDVal, err := strconv.Atoi(paymentIDStr); err == nil {
			paymentID = &paymentIDVal
		}
	}
	if paymentTypeStr := c.Query("payment_type"); paymentTypeStr != "" {
		if paymentTypeVal, err := strconv.Atoi(paymentTypeStr); err == nil {
			paymentType = &paymentTypeVal
		}
	}
	if typeStr := c.Query("type"); typeStr != "" {
		if typeVal, err := strconv.Atoi(typeStr); err == nil {
			separateType = &typeVal
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
	response, err := h.separateService.GetSeparateAccounts(
		id, uid, ordersID, childOrdersID, goodsID, paymentID, paymentType, separateType, page, pageSize,
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
