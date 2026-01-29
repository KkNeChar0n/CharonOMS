package goods

import (
	"net/http"
	"strconv"
	"strings"

	goodsService "charonoms/internal/application/service/goods"

	"github.com/gin-gonic/gin"
)

// GoodsHandler 商品HTTP处理器
type GoodsHandler struct {
	goodsService *goodsService.GoodsService
}

// NewGoodsHandler 创建商品HTTP处理器实例
func NewGoodsHandler(goodsService *goodsService.GoodsService) *GoodsHandler {
	return &GoodsHandler{
		goodsService: goodsService,
	}
}

// GetGoods 获取商品列表（支持按分类和状态筛选）
// GET /api/goods?classifyid=1&status=0
func (h *GoodsHandler) GetGoods(c *gin.Context) {
	// 获取查询参数
	classifyIDStr := c.Query("classifyid")
	statusStr := c.Query("status")

	goods, err := h.goodsService.GetGoodsList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 按分类ID筛选
	if classifyIDStr != "" {
		classifyID, err := strconv.Atoi(classifyIDStr)
		if err == nil {
			filtered := make([]map[string]interface{}, 0)
			for _, g := range goods {
				// 处理不同的整数类型
				var cid int
				switch v := g["classifyid"].(type) {
				case int:
					cid = v
				case int64:
					cid = int(v)
				case float64:
					cid = int(v)
				default:
					continue
				}
				if cid == classifyID {
					filtered = append(filtered, g)
				}
			}
			goods = filtered
		}
	}

	// 按状态筛选
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			filtered := make([]map[string]interface{}, 0)
			for _, g := range goods {
				// 处理不同的整数类型
				var s int
				switch v := g["status"].(type) {
				case int:
					s = v
				case int64:
					s = int(v)
				case float64:
					s = int(v)
				default:
					continue
				}
				if s == status {
					filtered = append(filtered, g)
				}
			}
			goods = filtered
		}
	}

	c.JSON(http.StatusOK, gin.H{"goods": goods})
}

// GetGoodsByID 获取商品详情
// GET /api/goods/:id
func (h *GoodsHandler) GetGoodsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goods id"})
		return
	}

	goods, err := h.goodsService.GetGoodsByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "商品不存在") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goods": goods})
}

// GetActiveForOrder 获取可用于下单的商品列表
// GET /api/goods/active-for-order
func (h *GoodsHandler) GetActiveForOrder(c *gin.Context) {
	goods, err := h.goodsService.GetActiveForOrder()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goods": goods})
}

// GetAvailableForCombo 获取可用于组合的单品商品列表
// GET /api/goods/available-for-combo
func (h *GoodsHandler) GetAvailableForCombo(c *gin.Context) {
	excludeID := 0
	if excludeIDStr := c.Query("exclude_id"); excludeIDStr != "" {
		if id, err := strconv.Atoi(excludeIDStr); err == nil {
			excludeID = id
		}
	}

	goods, err := h.goodsService.GetAvailableForCombo(excludeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goods": goods})
}

// GetIncludedGoods 根据父商品ID获取包含的子商品列表
// GET /api/goods/:id/included-goods
func (h *GoodsHandler) GetIncludedGoods(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goods id"})
		return
	}

	goods, err := h.goodsService.GetIncludedGoods(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"included_goods": goods})
}

// GetTotalPrice 计算商品总价
// GET /api/goods/:id/total-price
func (h *GoodsHandler) GetTotalPrice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goods id"})
		return
	}

	totalPrice, err := h.goodsService.GetTotalPrice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_price": totalPrice})
}

// CreateGoods 创建商品
// POST /api/goods
func (h *GoodsHandler) CreateGoods(c *gin.Context) {
	var req goodsService.CreateGoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goodsID, err := h.goodsService.CreateGoods(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      goodsID,
		"message": "商品添加成功",
	})
}

// UpdateGoods 更新商品
// PUT /api/goods/:id
func (h *GoodsHandler) UpdateGoods(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goods id"})
		return
	}

	var req goodsService.UpdateGoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.goodsService.UpdateGoods(id, &req); err != nil {
		if strings.Contains(err.Error(), "商品不存在") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品信息更新成功"})
}

// UpdateStatus 更新商品状态
// PUT /api/goods/:id/status
func (h *GoodsHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goods id"})
		return
	}

	var req goodsService.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.goodsService.UpdateStatus(id, &req); err != nil {
		if strings.Contains(err.Error(), "商品不存在") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品状态更新成功"})
}
