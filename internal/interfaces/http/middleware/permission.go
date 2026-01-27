package middleware

import (
	"charonoms/internal/infrastructure/persistence/mysql"
	"charonoms/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Permission 权限检查中间件
func Permission(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 超级管理员拥有所有权限
		if IsSuperAdmin(c) {
			c.Next()
			return
		}

		// 获取用户角色ID
		roleID := GetRoleID(c)
		if roleID == 0 {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		// 根据请求方法获取操作类型
		action := getActionFromMethod(c.Request.Method)
		actionID := fmt.Sprintf("%s_%s", action, resource)

		// 检查权限
		hasPermission := checkPermission(roleID, actionID)
		if !hasPermission {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// getActionFromMethod 根据HTTP方法获取操作类型
func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "view"
	case "POST":
		return "add"
	case "PUT", "PATCH":
		return "edit"
	case "DELETE":
		return "delete"
	default:
		return "view"
	}
}

// checkPermission 检查角色是否有权限
func checkPermission(roleID uint, actionID string) bool {
	var count int64
	err := mysql.DB.Table("role_permissions rp").
		Joins("JOIN permissions p ON rp.permissions_id = p.id").
		Where("rp.role_id = ? AND p.action_id = ? AND p.status = 0", roleID, actionID).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}
