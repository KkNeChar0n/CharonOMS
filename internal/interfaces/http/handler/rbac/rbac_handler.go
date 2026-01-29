package rbac

import (
	"charonoms/internal/application/service/rbac"
	"charonoms/internal/domain/rbac/entity"
	"charonoms/internal/interfaces/http/middleware"
	"charonoms/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RBACHandler RBAC 处理器
type RBACHandler struct {
	rbacService *rbac.RBACService
}

// NewRBACHandler 创建 RBAC 处理器实例
func NewRBACHandler(rbacService *rbac.RBACService) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
	}
}

// ===== 角色管理 =====

// GetRoles 获取角色列表（支持按状态过滤）
func (h *RBACHandler) GetRoles(c *gin.Context) {
	roles, err := h.rbacService.GetRoleList(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 支持按状态过滤
	statusParam := c.Query("status")
	if statusParam != "" {
		var status int8
		if statusParam == "0" {
			status = 0
		} else if statusParam == "1" {
			status = 1
		} else {
			c.JSON(http.StatusOK, gin.H{"roles": roles})
			return
		}

		// 过滤角色
		filteredRoles := make([]*entity.Role, 0)
		for _, role := range roles {
			if role.Status == status {
				filteredRoles = append(filteredRoles, role)
			}
		}
		c.JSON(http.StatusOK, gin.H{"roles": filteredRoles})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

// CreateRole 创建角色
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req rbac.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	role, err := h.rbacService.CreateRole(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, role)
}

// UpdateRole 更新角色
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req rbac.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.rbacService.UpdateRole(c.Request.Context(), uint(id), &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// UpdateRoleStatus 更新角色状态
func (h *RBACHandler) UpdateRoleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Status == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status参数不能为空"})
		return
	}

	if err := h.rbacService.UpdateRoleStatus(c.Request.Context(), uint(id), *req.Status); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ===== 权限管理 =====

// GetPermissions 获取权限列表
func (h *RBACHandler) GetPermissions(c *gin.Context) {
	permissions, err := h.rbacService.GetPermissionList(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// UpdatePermissionStatus 更新权限状态
func (h *RBACHandler) UpdatePermissionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	var req struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Status == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status参数不能为空"})
		return
	}

	if err := h.rbacService.UpdatePermissionStatus(c.Request.Context(), uint(id), *req.Status); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetPermissionTree 获取权限树
func (h *RBACHandler) GetPermissionTree(c *gin.Context) {
	tree, err := h.rbacService.GetPermissionTree(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 前端期望 { tree: [...] } 格式
	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

// ===== 角色权限关联 =====

// GetRolePermissions 获取角色的权限列表
func (h *RBACHandler) GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	permissions, err := h.rbacService.GetRolePermissions(c.Request.Context(), uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 前端期望 permission_ids 数组（只要ID）
	permissionIDs := make([]uint, 0, len(permissions))
	for _, perm := range permissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	c.JSON(http.StatusOK, gin.H{"permission_ids": permissionIDs})
}

// UpdateRolePermissions 更新角色权限
func (h *RBACHandler) UpdateRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req rbac.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.rbacService.UpdateRolePermissions(c.Request.Context(), uint(id), &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ===== 菜单管理 =====

// GetMenus 获取菜单列表（管理用）
func (h *RBACHandler) GetMenus(c *gin.Context) {
	menus, err := h.rbacService.GetMenuList(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"menus": menus})
}

// GetMenu 获取用户的菜单树（前端导航用）
func (h *RBACHandler) GetMenu(c *gin.Context) {
	roleID := middleware.GetRoleID(c)
	isSuperAdmin := middleware.IsSuperAdmin(c)

	menuTree, err := h.rbacService.GetUserMenuTree(c.Request.Context(), roleID, isSuperAdmin)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 前端期望 response.data.data.menus 格式
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"menus": menuTree,
		},
	})
}

// UpdateMenu 更新菜单
func (h *RBACHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	var req rbac.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.rbacService.UpdateMenu(c.Request.Context(), uint(id), &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// UpdateMenuStatus 更新菜单状态
func (h *RBACHandler) UpdateMenuStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	var req struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Status == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status参数不能为空"})
		return
	}

	if err := h.rbacService.UpdateMenuStatus(c.Request.Context(), uint(id), *req.Status); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
