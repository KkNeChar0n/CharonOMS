package middleware

import (
	"charonoms/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查来源是否在允许列表中
		allowed := false
		for _, allowedOrigin := range cfg.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else if len(cfg.AllowOrigins) > 0 {
				c.Header("Access-Control-Allow-Origin", cfg.AllowOrigins[0])
			}

			c.Header("Access-Control-Allow-Methods", joinStrings(cfg.AllowMethods, ", "))
			c.Header("Access-Control-Allow-Headers", joinStrings(cfg.AllowHeaders, ", "))

			if cfg.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			c.Header("Access-Control-Max-Age", "86400")
		}

		// 处理 OPTIONS 请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// joinStrings 连接字符串数组
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
