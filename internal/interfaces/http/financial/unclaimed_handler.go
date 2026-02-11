package financial

import (
	"charonoms/internal/application/financial/unclaimed"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type UnclaimedHandler struct {
	service *unclaimed.UnclaimedService
}

// NewUnclaimedHandler 创建常规待认领处理器
func NewUnclaimedHandler(service *unclaimed.UnclaimedService) *UnclaimedHandler {
	return &UnclaimedHandler{service: service}
}

// GetList 获取待认领列表
// GET /api/unclaimed
func (h *UnclaimedHandler) GetList(c *gin.Context) {
	filters := make(map[string]interface{})

	// 解析查询参数
	if id := c.Query("id"); id != "" {
		filters["id"] = id
	}
	if payer := c.Query("payer"); payer != "" {
		filters["payer"] = payer
	}
	if paymentMethod := c.Query("payment_method"); paymentMethod != "" {
		filters["payment_method"] = paymentMethod
	}
	if arrivalDate := c.Query("arrival_date"); arrivalDate != "" {
		filters["arrival_date"] = arrivalDate
	}
	if claimer := c.Query("claimer"); claimer != "" {
		filters["claimer"] = claimer
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	list, err := h.service.GetList(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unclaimed": list,
	})
}

// Claim 认领待认领款项
// PUT /api/unclaimed/:id/claim
func (h *UnclaimedHandler) Claim(c *gin.Context) {
	// 获取待认领ID
	unclaimedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的待认领ID",
		})
		return
	}

	// 解析请求体
	var req struct {
		OrderID int `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID（从会话或token中获取）
	userID := 1 // TODO: 从认证中间件获取实际用户ID

	// 执行认领操作
	if err := h.service.Claim(unclaimedID, req.OrderID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "认领成功",
	})
}

// Delete 删除待认领记录
// DELETE /api/unclaimed/:id
func (h *UnclaimedHandler) Delete(c *gin.Context) {
	// 获取待认领ID
	unclaimedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的待认领ID",
		})
		return
	}

	// 执行删除操作
	if err := h.service.Delete(unclaimedID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// DownloadTemplate 下载Excel导入模板
// GET /api/unclaimed/template
func (h *UnclaimedHandler) DownloadTemplate(c *gin.Context) {
	file, err := h.service.GenerateTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成模板失败: " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\"待认领收款模板.xlsx\"")
	c.Header("Content-Transfer-Encoding", "binary")

	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "下载模板失败: " + err.Error()})
		return
	}
}

// ImportExcel 导入Excel数据
// POST /api/unclaimed/import
func (h *UnclaimedHandler) ImportExcel(c *gin.Context) {
	// 获取上传的文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未上传文件"})
		return
	}

	if fileHeader.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名为空"})
		return
	}

	// 检查文件格式
	if !(len(fileHeader.Filename) > 5 && (fileHeader.Filename[len(fileHeader.Filename)-5:] == ".xlsx" || fileHeader.Filename[len(fileHeader.Filename)-4:] == ".xls")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持.xls和.xlsx格式"})
		return
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败"})
		return
	}
	defer file.Close()

	// 读取Excel文件
	f, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取Excel文件失败: " + err.Error()})
		return
	}
	defer f.Close()

	// 导入数据
	successCount, matchedCount, errorRows, err := h.service.ImportExcelFile(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导入失败: " + err.Error()})
		return
	}

	// 构造响应消息（严格按照Python版本格式）
	message := fmt.Sprintf("导入完成，成功%d条", successCount)
	if matchedCount > 0 {
		message += fmt.Sprintf("，其中%d条自动匹配到已收款", matchedCount)
	}

	response := gin.H{
		"message":       message,
		"success_count": successCount,
		"matched_count": matchedCount,
	}

	if len(errorRows) > 0 {
		response["errors"] = errorRows
	}

	// 如果有成功的或有部分成功，返回200；全部失败返回400
	if successCount > 0 {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}
