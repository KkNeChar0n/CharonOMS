package marketing

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"charonoms/internal/application/service/marketing"
)

// TemplateHandler 活动模板处理器
type TemplateHandler struct {
	templateService *marketing.TemplateService
}

// NewTemplateHandler 创建活动模板处理器实例
func NewTemplateHandler(templateService *marketing.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}

// GetTemplates 获取活动模板列表
func (h *TemplateHandler) GetTemplates(c *gin.Context) {
	var req marketing.TemplateListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templates, err := h.templateService.GetTemplateList(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetActiveTemplates 获取启用的活动模板列表
func (h *TemplateHandler) GetActiveTemplates(c *gin.Context) {
	templates, err := h.templateService.GetActiveTemplateList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetTemplateDetail 获取活动模板详情
func (h *TemplateHandler) GetTemplateDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	template, err := h.templateService.GetTemplateByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

// CreateTemplate 创建活动模板
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req marketing.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.templateService.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活动模板创建成功"})
}

// UpdateTemplate 更新活动模板
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	var req marketing.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.templateService.UpdateTemplate(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活动模板更新成功"})
}

// UpdateTemplateStatus 更新活动模板状态
func (h *TemplateHandler) UpdateTemplateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	var req marketing.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.templateService.UpdateTemplateStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活动模板状态更新成功"})
}
