package rbac

import (
	"charonoms/internal/application/service/rbac"
	"charonoms/internal/interfaces/http/middleware"
	"charonoms/pkg/response"
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

// GetRoles 获取角色列表
func (h *RBACHandler) GetRoles(c *gin.Context) {
	roles, err := h.rbacService.GetRoleList(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// CreateRole 创建角色
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req rbac.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	role, err := h.rbacService.CreateRole(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req rbac.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.rbacService.UpdateRole(c.Request.Context(), uint(id), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// UpdateRoleStatus 更新角色状态
func (h *RBACHandler) UpdateRoleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.rbacService.UpdateRoleStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// ===== 权限管理 =====

// GetPermissions 获取权限列表
func (h *RBACHandler) GetPermissions(c *gin.Context) {
	permissions, err := h.rbacService.GetPermissionList(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}

// UpdatePermissionStatus 更新权限状态
func (h *RBACHandler) UpdatePermissionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.rbacService.UpdatePermissionStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// GetPermissionTree 获取权限树
func (h *RBACHandler) GetPermissionTree(c *gin.Context) {
	tree, err := h.rbacService.GetPermissionTree(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tree)
}

// ===== 角色权限关联 =====

// GetRolePermissions 获取角色的权限列表
func (h *RBACHandler) GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	permissions, err := h.rbacService.GetRolePermissions(c.Request.Context(), uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}

// UpdateRolePermissions 更新角色权限
func (h *RBACHandler) UpdateRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req rbac.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.rbacService.UpdateRolePermissions(c.Request.Context(), uint(id), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// ===== 菜单管理 =====

// GetMenus 获取菜单列表（管理用）
func (h *RBACHandler) GetMenus(c *gin.Context) {
	menus, err := h.rbacService.GetMenuList(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"menus": menus})
}

// GetMenu 获取用户的菜单树（前端导航用）
func (h *RBACHandler) GetMenu(c *gin.Context) {
	roleID := middleware.GetRoleID(c)
	isSuperAdmin := middleware.IsSuperAdmin(c)

	menuTree, err := h.rbacService.GetUserMenuTree(c.Request.Context(), roleID, isSuperAdmin)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"menus": menuTree})
}

// UpdateMenu 更新菜单
func (h *RBACHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	var req rbac.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.rbacService.UpdateMenu(c.Request.Context(), uint(id), &req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// UpdateMenuStatus 更新菜单状态
func (h *RBACHandler) UpdateMenuStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.rbacService.UpdateMenuStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}
