package attribute

import (
	"charonoms/internal/application/service/attribute"
	"charonoms/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AttributeHandler 属性HTTP处理器
type AttributeHandler struct {
	service *attribute.AttributeService
}

// NewAttributeHandler 创建属性处理器实例
func NewAttributeHandler(service *attribute.AttributeService) *AttributeHandler {
	return &AttributeHandler{service: service}
}

// GetAttributes 获取属性列表
// @route GET /api/attributes
func (h *AttributeHandler) GetAttributes(c *gin.Context) {
	result, err := h.service.GetAttributeList()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// GetActiveAttributes 获取启用属性列表
// @route GET /api/attributes/active
func (h *AttributeHandler) GetActiveAttributes(c *gin.Context) {
	result, err := h.service.GetActiveAttributes()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// CreateAttribute 创建属性
// @route POST /api/attributes
func (h *AttributeHandler) CreateAttribute(c *gin.Context) {
	var req attribute.CreateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	id, err := h.service.CreateAttribute(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"message":      "属性创建成功",
		"attribute_id": id,
	})
}

// UpdateAttribute 更新属性信息
// @route PUT /api/attributes/:id
func (h *AttributeHandler) UpdateAttribute(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "属性ID无效")
		return
	}

	// 绑定请求体
	var req attribute.UpdateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层更新
	if err := h.service.UpdateAttribute(id, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": "属性更新成功",
	})
}

// UpdateAttributeStatus 更新属性状态
// @route PUT /api/attributes/:id/status
func (h *AttributeHandler) UpdateAttributeStatus(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "属性ID无效")
		return
	}

	// 绑定请求体
	var req attribute.UpdateAttributeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "状态值必须为0或1")
		return
	}

	// 调用服务层更新状态
	if err := h.service.UpdateAttributeStatus(id, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": "状态更新成功",
	})
}

// GetAttributeValues 获取属性值列表
// @route GET /api/attributes/:id/values
func (h *AttributeHandler) GetAttributeValues(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "属性ID无效")
		return
	}

	result, err := h.service.GetAttributeValues(id)
	if err != nil {
		if err.Error() == "属性不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// SaveAttributeValues 保存属性值
// @route POST /api/attributes/:id/values
func (h *AttributeHandler) SaveAttributeValues(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "属性ID无效")
		return
	}

	// 绑定请求体
	var req attribute.SaveValuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层保存
	if err := h.service.SaveAttributeValues(id, &req); err != nil {
		if err.Error() == "属性不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": "属性值保存成功",
	})
}
