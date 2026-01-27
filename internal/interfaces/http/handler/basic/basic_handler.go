package basic

import (
	"charonoms/internal/application/service/basic"
	"charonoms/pkg/response"

	"github.com/gin-gonic/gin"
)

// BasicHandler 基础数据处理器
type BasicHandler struct {
	basicService *basic.BasicService
}

// NewBasicHandler 创建基础数据处理器实例
func NewBasicHandler(basicService *basic.BasicService) *BasicHandler {
	return &BasicHandler{
		basicService: basicService,
	}
}

// GetAllSexes 获取所有性别
// @Summary Get all sexes
// @Tags Basic
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/sexes [get]
func (h *BasicHandler) GetAllSexes(c *gin.Context) {
	sexes, err := h.basicService.GetAllSexes(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"sexes": sexes,
	})
}

// GetActiveGrades 获取启用的年级
// @Summary Get active grades
// @Tags Basic
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/grades/active [get]
func (h *BasicHandler) GetActiveGrades(c *gin.Context) {
	grades, err := h.basicService.GetActiveGrades(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"grades": grades,
	})
}

// GetActiveSubjects 获取启用的学科
// @Summary Get active subjects
// @Tags Basic
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/subjects/active [get]
func (h *BasicHandler) GetActiveSubjects(c *gin.Context) {
	subjects, err := h.basicService.GetActiveSubjects(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"subjects": subjects,
	})
}
