package coach

import (
	"net/http"
	"strconv"

	coachService "charonoms/internal/application/service/coach"

	"github.com/gin-gonic/gin"
)

// CoachHandler 教练HTTP处理器
type CoachHandler struct {
	coachService *coachService.CoachService
}

// NewCoachHandler 创建教练HTTP处理器实例
func NewCoachHandler(coachService *coachService.CoachService) *CoachHandler {
	return &CoachHandler{
		coachService: coachService,
	}
}

// GetCoaches 获取教练列表
// GET /api/coaches
func (h *CoachHandler) GetCoaches(c *gin.Context) {
	coaches, err := h.coachService.GetCoachList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coaches": coaches})
}

// GetActiveCoaches 获取启用教练列表
// GET /api/coaches/active
func (h *CoachHandler) GetActiveCoaches(c *gin.Context) {
	coaches, err := h.coachService.GetActiveCoaches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coaches": coaches})
}

// CreateCoach 创建教练
// POST /api/coaches
func (h *CoachHandler) CreateCoach(c *gin.Context) {
	var req coachService.CreateCoachRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coachID, err := h.coachService.CreateCoach(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      coachID,
		"message": "教练添加成功",
	})
}

// UpdateCoach 更新教练信息
// PUT /api/coaches/:id
func (h *CoachHandler) UpdateCoach(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coach id"})
		return
	}

	var req coachService.UpdateCoachRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.coachService.UpdateCoach(id, &req); err != nil {
		if err.Error() == "coach not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "教练信息更新成功"})
}

// UpdateCoachStatus 更新教练状态
// PUT /api/coaches/:id/status
func (h *CoachHandler) UpdateCoachStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coach id"})
		return
	}

	var req coachService.UpdateCoachStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := *req.Status
	if err := h.coachService.UpdateCoachStatus(id, status); err != nil {
		if err.Error() == "coach not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "操作成功"})
}

// DeleteCoach 删除教练
// DELETE /api/coaches/:id
func (h *CoachHandler) DeleteCoach(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coach id"})
		return
	}

	if err := h.coachService.DeleteCoach(id); err != nil {
		if err.Error() == "coach not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "教练删除成功"})
}
