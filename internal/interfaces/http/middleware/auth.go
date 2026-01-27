package middleware

import (
	"charonoms/internal/infrastructure/config"
	"charonoms/pkg/jwt"
	"charonoms/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// Token 格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "Token 格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 Token
		claims, err := jwt.ParseToken(tokenString, config.GlobalConfig.JWT)
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)
		c.Set("is_super_admin", claims.IsSuperAdmin)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetRoleID 从上下文获取角色ID
func GetRoleID(c *gin.Context) uint {
	if roleID, exists := c.Get("role_id"); exists {
		return roleID.(uint)
	}
	return 0
}

// IsSuperAdmin 判断是否为超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	if isSuperAdmin, exists := c.Get("is_super_admin"); exists {
		return isSuperAdmin.(bool)
	}
	return false
}
