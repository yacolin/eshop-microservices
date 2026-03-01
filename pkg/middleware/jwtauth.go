package middleware

import (
	"eshop-microservices/pkg/errcode"
	"eshop-microservices/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth validates the Authorization header and token, but delegates error
// responses to the centralized ErrorHandler middleware by calling `c.Error(err)`.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		// extract user_id and store as typed uint in context
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		// convert to uint and set
		c.Set("user_id", uint(userIDFloat))
		c.Next()
	}
}

// GetUserID returns the typed user id set by `JWTAuth` and a boolean indicating presence.
// Use this helper in handlers to avoid casting from float64 everywhere.
// func GetUserID(c *gin.Context) (uint, bool) {
// 	v, ok := c.Get("user_id")
// 	if !ok {
// 		return 0, false
// 	}
// 	switch id := v.(type) {
// 	case uint:
// 		return id, true
// 	case int:
// 		return uint(id), true
// 	case int64:
// 		return uint(id), true
// 	case float64:
// 		return uint(id), true
// 	default:
// 		return 0, false
// 	}
// }
