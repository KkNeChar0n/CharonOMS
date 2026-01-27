package auth

import (
	"charonoms/internal/application/service/auth"
	"charonoms/internal/interfaces/http/middleware"
	"charonoms/pkg/response"

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

	response.SuccessWithMessage(c, "鐧诲綍鎴愬姛", resp)
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

	response.Success(c, gin.H{
		"username": user.Username,
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

	response.Success(c, resp)
}

// Logout user logout
// @Summary User logout
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT is stateless, logout only needs frontend to delete token
	// If token blacklist is needed, add logic here
	response.SuccessWithMessage(c, "Logout successful", nil)
}