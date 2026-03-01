package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复 panic 并记录错误日志
func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err any) {
		log.Printf("panic: %v", err)
		c.AbortWithStatus(500)
	})
}
