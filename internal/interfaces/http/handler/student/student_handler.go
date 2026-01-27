package student

import (
	"charonoms/internal/application/service/student"
	"charonoms/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StudentHandler 学生HTTP处理器
type StudentHandler struct {
	service *student.StudentService
}

// NewStudentHandler 创建学生处理器实例
func NewStudentHandler(service *student.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

// GetStudents 获取学生列表
// @route GET /api/students
func (h *StudentHandler) GetStudents(c *gin.Context) {
	result, err := h.service.GetStudentList()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// GetActiveStudents 获取启用学生列表
// @route GET /api/students/active
func (h *StudentHandler) GetActiveStudents(c *gin.Context) {
	result, err := h.service.GetActiveStudents()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}

// CreateStudent 创建学生
// @route POST /api/students
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req student.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	id, err := h.service.CreateStudent(&req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{
		"id":      id,
		"message": "学生添加成功",
	})
}

// UpdateStudent 更新学生信息
// @route PUT /api/students/:id
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "学生ID无效")
		return
	}

	// 绑定请求体
	var req student.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层更新
	if err := h.service.UpdateStudent(id, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"message": "学生信息更新成功",
	})
}

// UpdateStudentStatus 更新学生状态
// @route PUT /api/students/:id/status
func (h *StudentHandler) UpdateStudentStatus(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "学生ID无效")
		return
	}

	// 绑定请求体
	var req student.UpdateStudentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录详细错误信息用于调试
		c.JSON(400, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层更新状态
	if err := h.service.UpdateStudentStatus(id, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "操作成功", nil)
}

// DeleteStudent 删除学生
// @route DELETE /api/students/:id
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "学生ID无效")
		return
	}

	// 调用服务层删除
	if err := h.service.DeleteStudent(id); err != nil {
		// 如果是订单关联错误，返回400
		if err.Error() == "无法删除，该学生存在关联订单" {
			response.BadRequest(c, err.Error())
			return
		}
		response.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"message": "学生删除成功",
	})
}
