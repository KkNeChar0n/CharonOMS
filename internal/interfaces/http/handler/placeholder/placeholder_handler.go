package placeholder

import (
	"github.com/gin-gonic/gin"
)

// PlaceholderHandler 占位符处理器，用于尚未实现的功能
type PlaceholderHandler struct{}

// NewPlaceholderHandler 创建占位符处理器实例
func NewPlaceholderHandler() *PlaceholderHandler {
	return &PlaceholderHandler{}
}

// HandlePlaceholder 返回占位符响应
func (h *PlaceholderHandler) HandlePlaceholder(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "功能开发中",
		"status":  "pending",
	})
}
