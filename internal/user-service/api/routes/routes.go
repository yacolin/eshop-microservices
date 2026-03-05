package routes

import (
	"eshop-microservices/internal/user-service/api/handlers"
	"eshop-microservices/internal/user-service/api/middleware"
	pkgmiddleware "eshop-microservices/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, permissionHandler *handlers.PermissionHandler) {
	r.Use(middleware.Recovery(), middleware.Logger(), pkgmiddleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "user ok"})
	})

	api := r.Group("/api")
	registerV1(api, userHandler, authHandler, permissionHandler)
}

func registerV1(api *gin.RouterGroup, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, permissionHandler *handlers.PermissionHandler) {
	// 认证相关路由（统一使用 auth handler）
	auth := api.Group("/v1/auth")
	{
		// 登录路由
		auth.POST("/login/password", authHandler.LoginByPassword)
		auth.POST("/login/wechat", authHandler.LoginByWechat)
		auth.POST("/login/phone", authHandler.LoginByPhone)

		// 注册路由
		auth.POST("/register", authHandler.Register)

		// Token管理
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)

		// 需要认证的路由
		auth.Use(pkgmiddleware.JWTAuth())
		auth.GET("/me", authHandler.GetCurrentUser)
	}

	// 用户相关路由（只保留用户资料管理，登录/注册/登出统一使用 /auth）
	users := api.Group("/v1/users")
	{
		// 需要认证的路由
		protected := users.Group("").Use(pkgmiddleware.JWTAuth())
		{
			// 获取用户资料（包含 User 和 UserInfo）
			protected.GET("/profile", userHandler.GetProfile)

			// 用户详细信息管理（Avatar、Nickname 等）
			protected.GET("/info", userHandler.GetUserInfo)
			protected.PUT("/info", userHandler.UpdateUserInfo)

			// 根据ID获取用户（管理员接口）
			protected.GET("/:id", userHandler.GetByID)

			// 更新用户角色（管理员接口）
			protected.PUT("/:user_id/roles", permissionHandler.UpdateUserRole)
		}
	}

	// 权限相关路由
	permissions := api.Group("/v1/permissions")
	{
		// 公开路由（已认证用户可以检查自己的权限）
		permissions.Use(pkgmiddleware.JWTAuth())
		{
			// 查询权限
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.GET("/:id", permissionHandler.GetPermission)

			// 检查权限
			permissions.POST("/check", permissionHandler.CheckPermissions)
		}

		// 管理员路由（需要 admin 角色）
		admin := permissions.Group("").Use(pkgmiddleware.JWTAuth(), middleware.RequireAdmin())
		{
			// 权限 CRUD
			admin.POST("", permissionHandler.CreatePermission)
			admin.PUT("/:id", permissionHandler.UpdatePermission)
			admin.DELETE("/:id", permissionHandler.DeletePermission)
		}
	}

	// 角色权限管理路由
	roles := api.Group("/v1/roles")
	{
		// 需要认证的路由
		roles.Use(pkgmiddleware.JWTAuth())
		{
			// 查看角色权限
			roles.GET("/:role/permissions", permissionHandler.GetRolePermissions)
		}

		// 管理员路由
		admin := roles.Group("").Use(pkgmiddleware.JWTAuth(), middleware.RequireAdmin())
		{
			// 分配/移除权限给角色
			admin.POST("/:role/permissions", permissionHandler.AssignPermissionToRole)
			admin.DELETE("/:role/permissions/:permission_id", permissionHandler.RemovePermissionFromRole)
			admin.POST("/:role/permissions/batch", permissionHandler.BatchAssignPermissionsToRole)
		}
	}
}
