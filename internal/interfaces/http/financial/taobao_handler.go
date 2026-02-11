package financial

import (
	"charonoms/internal/application/financial/taobao"
	taobaoEntity "charonoms/internal/domain/financial/taobao"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type TaobaoHandler struct {
	service *taobao.TaobaoPaymentService
}

func NewTaobaoHandler(service *taobao.TaobaoPaymentService) *TaobaoHandler {
	return &TaobaoHandler{service: service}
}

// GetPaidList 获取淘宝已付款列表
func (h *TaobaoHandler) GetPaidList(c *gin.Context) {
	filters := make(map[string]interface{})

	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filters["id"] = id
		}
	}
	if studentIDStr := c.Query("student_id"); studentIDStr != "" {
		if studentID, err := strconv.Atoi(studentIDStr); err == nil {
			filters["student_id"] = studentID
		}
	}
	if orderIDStr := c.Query("order_id"); orderIDStr != "" {
		if orderID, err := strconv.Atoi(orderIDStr); err == nil {
			filters["order_id"] = orderID
		}
	}
	if orderDate := c.Query("order_date"); orderDate != "" {
		filters["order_date"] = orderDate
	}
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filters["status"] = status
		}
	}

	payments, err := h.service.GetList(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

// CreateTaobaoPayment 创建淘宝收款
func (h *TaobaoHandler) CreateTaobaoPayment(c *gin.Context) {
	var req struct {
		StudentID       int     `json:"student_id" binding:"required,min=1"`
		OrderID         int     `json:"order_id" binding:"required,min=1"`
		ZhifubaoAccount string  `json:"zhifubao_account"`
		Payer           string  `json:"payer"`
		PaymentAmount   float64 `json:"payment_amount" binding:"required,gt=0"`
		OrderTime       string  `json:"order_time" binding:"required"`
		MerchantOrder   string  `json:"merchant_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	// 解析时间，支持多种格式
	var orderTime time.Time
	var err error
	// 尝试 datetime-local 格式 (YYYY-MM-DDTHH:mm)
	orderTime, err = time.Parse("2006-01-02T15:04", req.OrderTime)
	if err != nil {
		// 尝试完整的 datetime-local 格式 (YYYY-MM-DDTHH:mm:ss)
		orderTime, err = time.Parse("2006-01-02T15:04:05", req.OrderTime)
		if err != nil {
			// 尝试标准格式 (YYYY-MM-DD HH:mm:ss)
			orderTime, err = time.Parse("2006-01-02 15:04:05", req.OrderTime)
			if err != nil {
				// 尝试简短格式 (YYYY-MM-DD HH:mm)
				orderTime, err = time.Parse("2006-01-02 15:04", req.OrderTime)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "时间格式错误，支持格式: YYYY-MM-DDTHH:mm, YYYY-MM-DD HH:mm:ss"})
					return
				}
			}
		}
	}

	payment := &taobaoEntity.TaobaoPayment{
		StudentID:       &req.StudentID,
		OrderID:         &req.OrderID,
		ZhifubaoAccount: &req.ZhifubaoAccount,
		Payer:           &req.Payer,
		PaymentAmount:   req.PaymentAmount,
		OrderTime:       &orderTime,
		MerchantOrder:   &req.MerchantOrder,
	}

	if err := h.service.Create(payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": payment})
}

// ConfirmArrival 确认淘宝收款到账
func (h *TaobaoHandler) ConfirmArrival(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := h.service.ConfirmArrival(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "确认到账成功"})
}

// DeleteTaobaoPayment 删除淘宝收款
func (h *TaobaoHandler) DeleteTaobaoPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// GetUnclaimedList 获取淘宝待认领列表
func (h *TaobaoHandler) GetUnclaimedList(c *gin.Context) {
	filters := make(map[string]interface{})

	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filters["id"] = id
		}
	}
	if arrivalDate := c.Query("arrival_date"); arrivalDate != "" {
		filters["arrival_date"] = arrivalDate
	}
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filters["status"] = status
		}
	}

	payments, err := h.service.GetUnclaimedList(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取淘宝待认领列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unclaimed": payments})
}

// ClaimUnclaimed 认领淘宝待认领
func (h *TaobaoHandler) ClaimTaobaoPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var req struct {
		OrderID int `json:"order_id" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	// 获取当前用户ID (从JWT中间件获取)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	if err := h.service.ClaimUnclaimed(id, req.OrderID, int(userID.(uint))); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "认领成功"})
}

// DeleteUnclaimed 删除淘宝待认领
func (h *TaobaoHandler) DeleteUnclaimed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := h.service.DeleteUnclaimed(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// DownloadUnclaimedTemplate 下载淘宝待认领导入模板
func (h *TaobaoHandler) DownloadUnclaimedTemplate(c *gin.Context) {
	file, err := h.service.GenerateUnclaimedTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成模板失败: " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\"淘宝待认领模板.xlsx\"")
	c.Header("Content-Transfer-Encoding", "binary")

	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "下载模板失败: " + err.Error()})
		return
	}
}

// ImportUnclaimedExcel 导入淘宝待认领Excel
func (h *TaobaoHandler) ImportUnclaimedExcel(c *gin.Context) {
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
	successCount, matchedCount, errorRows, err := h.service.ImportUnclaimedExcelFile(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导入失败: " + err.Error()})
		return
	}

	// 构造响应消息
	message := fmt.Sprintf("导入完成，成功%d条", successCount)
	if matchedCount > 0 {
		message += fmt.Sprintf("，其中%d条自动匹配到已付款", matchedCount)
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
