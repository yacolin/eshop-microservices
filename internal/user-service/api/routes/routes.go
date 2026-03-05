package routes

import (
	"eshop-microservices/internal/user-service/api/handlers"
	"eshop-microservices/internal/user-service/api/middleware"
	"eshop-microservices/internal/user-service/domain/repositories"
	pkgmiddleware "eshop-microservices/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, permissionHandler *handlers.PermissionHandler, roleHandler *handlers.RoleHandler, roleRepo repositories.RoleRepository) {
	r.Use(middleware.Recovery(), middleware.Logger(), pkgmiddleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "user ok"})
	})

	api := r.Group("/api")
	registerV1(api, userHandler, authHandler, permissionHandler, roleHandler, roleRepo)
}

func registerV1(api *gin.RouterGroup, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, permissionHandler *handlers.PermissionHandler, roleHandler *handlers.RoleHandler, roleRepo repositories.RoleRepository) {
	auth := api.Group("/v1/auth")
	{
		auth.POST("/login/password", authHandler.LoginByPassword)
		auth.POST("/login/wechat", authHandler.LoginByWechat)
		auth.POST("/login/phone", authHandler.LoginByPhone)

		auth.POST("/register", authHandler.Register)

		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)

		auth.Use(pkgmiddleware.JWTAuth())
		auth.GET("/me", authHandler.GetCurrentUser)
	}

	users := api.Group("/v1/users")
	{
		protected := users.Group("").Use(pkgmiddleware.JWTAuth())
		{
			protected.GET("/profile", userHandler.GetProfile)
			protected.GET("/info", userHandler.GetUserInfo)
			protected.PUT("/info", userHandler.UpdateUserInfo)
			protected.GET("/:id", userHandler.GetByID)
		}
	}

	permissions := api.Group("/v1/permissions")
	{
		permissions.Use(pkgmiddleware.JWTAuth())
		{
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.GET("/:id", permissionHandler.GetPermission)
			permissions.POST("/check", permissionHandler.CheckPermissions)
		}

		roleConfig := middleware.NewRequireRoleConfig(roleRepo)
		admin := permissions.Group("").Use(pkgmiddleware.JWTAuth(), middleware.RequireAdmin(roleConfig))
		{
			admin.POST("", permissionHandler.CreatePermission)
			admin.PUT("/:id", permissionHandler.UpdatePermission)
			admin.DELETE("/:id", permissionHandler.DeletePermission)
		}
	}

	roles := api.Group("/v1/roles")
	{
		roles.Use(pkgmiddleware.JWTAuth())
		{
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.GET("/name/:name", roleHandler.GetRoleByName)
		}

		roleConfig := middleware.NewRequireRoleConfig(roleRepo)
		admin := roles.Group("").Use(pkgmiddleware.JWTAuth(), middleware.RequireAdmin(roleConfig))
		{
			admin.POST("", roleHandler.CreateRole)
			admin.PUT("/:id", roleHandler.UpdateRole)
			admin.DELETE("/:id", roleHandler.DeleteRole)

			admin.POST("/:id/permissions", roleHandler.AssignPermissionsToRole)
			admin.DELETE("/:id/permissions", roleHandler.RemovePermissionsFromRole)
		}
	}

	userRoles := api.Group("/v1/users/:user_id/roles")
	{
		userRoles.Use(pkgmiddleware.JWTAuth())
		{
			userRoles.GET("", roleHandler.GetUserRoles)
		}

		roleConfig := middleware.NewRequireRoleConfig(roleRepo)
		admin := userRoles.Group("").Use(pkgmiddleware.JWTAuth(), middleware.RequireMerchant(roleConfig))
		{
			admin.POST("", roleHandler.AssignRoleToUser)
			admin.DELETE("/:role_id", roleHandler.RemoveRoleFromUser)
		}
	}
}
