package middleware

import (
	"eshop-microservices/pkg/errcode"
	"strings"

	"github.com/gin-gonic/gin"
)

// PermissionService 权限服务接口（用于中间件）
type PermissionService interface {
	CheckUserPermissions(userID string, permissionNames []string) (map[string]bool, error)
}

// RequirePermission 权限校验中间件
// 使用方式: RequirePermission("order:create", "product:read")
func RequirePermission(service PermissionService, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 JWT 中获取用户信息
		userID, exists := c.Get("user_id")
		if !exists {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查用户权限
		result, err := service.CheckUserPermissions(userIDStr, permissions)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// 检查是否有任一权限
		for _, permission := range permissions {
			if has, ok := result[permission]; ok && has {
				c.Next()
				return
			}
		}

		// 无权限
		c.Error(errcode.ErrInsufficientPermissions)
		c.Abort()
	}
}

// RequireRole 角色校验中间件
// 使用方式: RequireRole("admin", "operator")
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 JWT 中获取用户角色
		rolesClaim, exists := c.Get("roles")
		if !exists {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		userRoles, ok := rolesClaim.(string)
		if !ok {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查用户是否有任一指定角色
		userRoleList := strings.Split(userRoles, ",")
		for _, role := range roles {
			for _, userRole := range userRoleList {
				if strings.TrimSpace(userRole) == role {
					c.Next()
					return
				}
			}
		}

		// 无角色权限
		c.Error(errcode.ErrInsufficientPermissions)
		c.Abort()
	}
}

// RequireAdmin 需要管理员角色的中间件
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireMerchant 需要商家或管理员角色的中间件
func RequireMerchant() gin.HandlerFunc {
	return RequireRole("merchant", "admin")
}
