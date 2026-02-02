package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"charonoms/internal/application/activity"
	"charonoms/internal/interfaces/http/dto"
)

// ActivityHandler 活动处理器
type ActivityHandler struct {
	service *activity.Service
}

// NewActivityHandler 创建活动处理器实例
func NewActivityHandler(service *activity.Service) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

// CreateActivity 创建活动
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req dto.CreateActivityDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateActivity(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "活动创建成功",
		"id":      id,
	})
}

// UpdateActivity 更新活动
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateActivityDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateActivity(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "活动更新成功",
	})
}

// DeleteActivity 删除活动
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteActivity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "活动删除成功",
	})
}

// GetActivity 获取活动详情
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	activity, err := h.service.GetActivity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activity": activity,
	})
}

// ListActivities 查询活动列表
func (h *ActivityHandler) ListActivities(c *gin.Context) {
	activities, err := h.service.ListActivities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
	})
}

// GetActivitiesByDateRange 按日期范围查询活动
func (h *ActivityHandler) GetActivitiesByDateRange(c *gin.Context) {
	paymentTimeStr := c.Query("payment_time")
	if paymentTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_time is required"})
		return
	}

	// 解析时间，支持多种格式
	var paymentTime time.Time
	var err error

	// 支持的时间格式列表
	formats := []string{
		time.RFC3339,                   // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05",          // ISO8601 without timezone
		"2006-01-02T15:04",             // ISO8601 without seconds
		"2006-01-02 15:04:05",          // MySQL datetime
		"2006-01-02 15:04",             // datetime without seconds
		"2006-01-02",                   // date only
		time.RFC3339Nano,               // with nanoseconds
	}

	for _, format := range formats {
		paymentTime, err = time.Parse(format, paymentTimeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment_time format"})
		return
	}

	result, err := h.service.GetActivitiesByDateRange(c.Request.Context(), paymentTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateActivityStatus 更新活动状态
func (h *ActivityHandler) UpdateActivityStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateActivityStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateActivityStatus(c.Request.Context(), id, *req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "状态更新成功",
	})
}
