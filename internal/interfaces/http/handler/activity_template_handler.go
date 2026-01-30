package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"charonoms/internal/application/activity_template"
	"charonoms/internal/interfaces/http/dto"
)

// ActivityTemplateHandler 活动模板处理器
type ActivityTemplateHandler struct {
	service *activity_template.Service
}

// NewActivityTemplateHandler 创建活动模板处理器实例
func NewActivityTemplateHandler(service *activity_template.Service) *ActivityTemplateHandler {
	return &ActivityTemplateHandler{
		service: service,
	}
}

// CreateTemplate 创建活动模板
func (h *ActivityTemplateHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateActivityTemplateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "活动模板创建成功",
		"id":      id,
	})
}

// UpdateTemplate 更新活动模板
func (h *ActivityTemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateActivityTemplateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateTemplate(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "活动模板更新成功",
	})
}

// DeleteTemplate 删除活动模板
func (h *ActivityTemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "活动模板删除成功",
	})
}

// GetTemplate 获取活动模板详情
func (h *ActivityTemplateHandler) GetTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	template, err := h.service.GetTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"template": template,
	})
}

// ListTemplates 查询活动模板列表
func (h *ActivityTemplateHandler) ListTemplates(c *gin.Context) {
	templates, err := h.service.ListTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
	})
}

// ListActiveTemplates 查询启用的活动模板
func (h *ActivityTemplateHandler) ListActiveTemplates(c *gin.Context) {
	templates, err := h.service.ListActiveTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
	})
}

// UpdateTemplateStatus 更新活动模板状态
func (h *ActivityTemplateHandler) UpdateTemplateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateTemplateStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateTemplateStatus(c.Request.Context(), id, *req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "状态更新成功",
	})
}
