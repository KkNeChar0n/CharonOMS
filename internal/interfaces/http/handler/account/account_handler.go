package account

import (
	"charonoms/internal/application/service/account"
	"charonoms/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AccountHandler 账号处理器
type AccountHandler struct {
	accountService *account.AccountService
}

// NewAccountHandler 创建账号处理器实例
func NewAccountHandler(accountService *account.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// GetAccounts 获取账号列表
// @Summary Get account list
// @Tags Account
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/accounts [get]
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	// 解析查询参数
	filters := make(map[string]interface{})
	if id := c.Query("id"); id != "" {
		filters["id"] = id
	}
	if phone := c.Query("phone"); phone != "" {
		filters["phone"] = phone
	}
	if roleID := c.Query("role_id"); roleID != "" {
		filters["role_id"] = roleID
	}

	resp, err := h.accountService.GetAccountList(c.Request.Context(), filters)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 直接返回accounts数组，与原Python项目格式保持一致
	c.JSON(200, gin.H{
		"accounts": resp.Accounts,
	})
}

// CreateAccount 创建账号
// @Summary Create account
// @Tags Account
// @Accept json
// @Produce json
// @Param body body account.CreateAccountRequest true "Create account request"
// @Success 200 {object} response.Response
// @Router /api/accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req account.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.accountService.CreateAccount(c.Request.Context(), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateAccount 更新账号
// @Summary Update account
// @Tags Account
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param body body account.UpdateAccountRequest true "Update account request"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的账号ID")
		return
	}

	var req account.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.accountService.UpdateAccount(c.Request.Context(), uint(id), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// UpdateAccountStatus 更新账号状态
// @Summary Update account status
// @Tags Account
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param body body object true "Status"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id}/status [put]
func (h *AccountHandler) UpdateAccountStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的账号ID")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.accountService.UpdateAccountStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}
