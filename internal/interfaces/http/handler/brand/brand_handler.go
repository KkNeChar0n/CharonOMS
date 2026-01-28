package brand

import (
	"net/http"
	"strconv"

	brandService "charonoms/internal/application/service/brand"

	"github.com/gin-gonic/gin"
)

// BrandHandler 品牌HTTP处理器
type BrandHandler struct {
	brandService *brandService.BrandService
}

// NewBrandHandler 创建品牌HTTP处理器实例
func NewBrandHandler(brandService *brandService.BrandService) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
	}
}

// GetBrands 获取所有品牌列表
// GET /api/brands
func (h *BrandHandler) GetBrands(c *gin.Context) {
	brands, err := h.brandService.GetBrandList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"brands": brands})
}

// GetActiveBrands 获取启用状态的品牌列表
// GET /api/brands/active
func (h *BrandHandler) GetActiveBrands(c *gin.Context) {
	brands, err := h.brandService.GetActiveBrands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"brands": brands})
}

// CreateBrand 创建品牌
// POST /api/brands
func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var req brandService.CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.brandService.CreateBrand(&req); err != nil {
		// 处理业务错误
		errMsg := err.Error()
		if errMsg == "品牌名称不能为空" || errMsg == "该品牌名称已存在" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "品牌添加成功"})
}

// UpdateBrand 更新品牌信息
// PUT /api/brands/:id
func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	var req brandService.UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.brandService.UpdateBrand(id, &req); err != nil {
		errMsg := err.Error()
		// 处理业务错误
		if errMsg == "品牌名称不能为空" || errMsg == "该品牌名称已存在" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		if errMsg == "品牌不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "品牌信息更新成功"})
}

// UpdateBrandStatus 更新品牌状态
// PUT /api/brands/:id/status
func (h *BrandHandler) UpdateBrandStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	var req brandService.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态不能为空"})
		return
	}

	status := *req.Status
	if err := h.brandService.UpdateBrandStatus(id, status); err != nil {
		errMsg := err.Error()
		if errMsg == "状态不能为空" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		if errMsg == "品牌不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "品牌状态更新成功"})
}
