package auth

import (
	"charonoms/internal/application/service/auth"
	"charonoms/internal/interfaces/http/middleware"
	"charonoms/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handler
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler create auth handler instance
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login user login
// @Summary User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.LoginRequest true "Login request"
// @Success 200 {object} response.Response
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid parameters")
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 设置JWT token到cookie（前端兼容性）
	c.SetCookie("auth_token", resp.Token, 86400, "/", "", false, true)

	// 返回前端期望的格式（包含token和data包装层）
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":          resp.Token,
			"username":       resp.Username,
			"is_super_admin": resp.IsSuperAdmin,
		},
	})
}

// GetProfile get current user info
// @Summary Get current user info
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "Not logged in")
		return
	}

	user, err := h.authService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// Don't return password
	user.Password = ""

	// 获取是否为超级管理员
	isSuperAdmin := false
	if user.Role != nil && user.Role.IsSuperAdmin == 1 {
		isSuperAdmin = true
	}

	// 返回格式同时包含顶层和data层（前端checkLoginStatus需要）
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"data": gin.H{
			"username":       user.Username,
			"is_super_admin": isSuperAdmin,
		},
	})
}

// SyncRole sync user role
// @Summary Sync user role
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/sync-role [get]
func (h *AuthHandler) SyncRole(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "Not logged in")
		return
	}

	oldRoleID := middleware.GetRoleID(c)
	oldIsSuperAdmin := middleware.IsSuperAdmin(c)

	resp, err := h.authService.SyncRole(c.Request.Context(), userID, oldRoleID, oldIsSuperAdmin)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 返回格式保持与其他API一致（带data包装层）
	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// GetUserPermissions get current user permissions
// @Summary Get current user permissions
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/user/permissions [get]
func (h *AuthHandler) GetUserPermissions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "Not logged in")
		return
	}

	roleID := middleware.GetRoleID(c)
	isSuperAdmin := middleware.IsSuperAdmin(c)

	permissions, err := h.authService.GetUserPermissions(c.Request.Context(), roleID, isSuperAdmin)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 返回权限的 action_id 列表
	actionIDs := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		actionIDs = append(actionIDs, perm.ActionID)
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": actionIDs,
	})
}

// Logout user logout
// @Summary User logout
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 清除cookie中的auth_token
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	// 返回平铺JSON格式
	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}