package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 记录请求日志
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return "[" + param.TimeStamp.Format(time.RFC3339) + "] " +
			param.Method + " " + param.Path + " " +
			param.ErrorMessage + "\n"
	})
}
