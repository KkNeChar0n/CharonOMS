package classify

import (
	"charonoms/internal/application/service/classify"
	"charonoms/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ClassifyHandler 分类HTTP处理器
type ClassifyHandler struct {
	service *classify.ClassifyService
}

// NewClassifyHandler 创建分类处理器实例
func NewClassifyHandler(service *classify.ClassifyService) *ClassifyHandler {
	return &ClassifyHandler{service: service}
}

// GetClassifies 获取分类列表
// @route GET /api/classifies
func (h *ClassifyHandler) GetClassifies(c *gin.Context) {
	result, err := h.service.GetClassifyList()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// GetParents 获取一级分类列表
// @route GET /api/classifies/parents
func (h *ClassifyHandler) GetParents(c *gin.Context) {
	result, err := h.service.GetParents()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// GetActiveClassifies 获取启用分类列表
// @route GET /api/classifies/active
func (h *ClassifyHandler) GetActiveClassifies(c *gin.Context) {
	result, err := h.service.GetActiveClassifies()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// GetClassifyByID 获取分类详情
// @route GET /api/classifies/:id
func (h *ClassifyHandler) GetClassifyByID(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "分类ID无效")
		return
	}

	// 调用服务层获取详情
	result, err := h.service.GetClassifyByID(id)
	if err != nil {
		// 如果是"类型不存在"错误，返回404
		if err.Error() == "类型不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// CreateClassify 创建分类
// @route POST /api/classifies
func (h *ClassifyHandler) CreateClassify(c *gin.Context) {
	var req classify.CreateClassifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	id, err := h.service.CreateClassify(&req)
	if err != nil {
		// 业务逻辑错误返回400
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"message":     "类型创建成功",
		"classify_id": id,
	})
}

// UpdateClassify 更新分类信息
// @route PUT /api/classifies/:id
func (h *ClassifyHandler) UpdateClassify(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "分类ID无效")
		return
	}

	// 绑定请求体
	var req classify.UpdateClassifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层更新
	if err := h.service.UpdateClassify(id, &req); err != nil {
		// 如果是"类型不存在"错误，返回404
		if err.Error() == "类型不存在" {
			response.NotFound(c, err.Error())
			return
		}
		// 其他业务逻辑错误返回400
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": "类型更新成功",
	})
}

// UpdateClassifyStatus 更新分类状态
// @route PUT /api/classifies/:id/status
func (h *ClassifyHandler) UpdateClassifyStatus(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "分类ID无效")
		return
	}

	// 绑定请求体
	var req classify.UpdateClassifyStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录详细错误信息用于调试
		c.JSON(400, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层更新状态
	if err := h.service.UpdateClassifyStatus(id, &req); err != nil {
		// 如果是"类型不存在"错误，返回404
		if err.Error() == "类型不存在" {
			response.NotFound(c, err.Error())
			return
		}
		// 其他业务逻辑错误返回400
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": "状态更新成功",
	})
}
