package routes

import (
	"eshop-microservices/internal/user-service/api/handlers"
	"eshop-microservices/internal/user-service/api/middleware"
	pkgmiddleware "eshop-microservices/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, userHandler *handlers.UserHandler) {
	r.Use(middleware.Recovery(), middleware.Logger(), pkgmiddleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	registerV1(api, userHandler)
}

func registerV1(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := api.Group("/v1/users")

	// 公开路由（不需要认证）
	public := users.Group("")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	protected := users.Group("").Use(pkgmiddleware.JWTAuth())
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.PUT("/profile", userHandler.UpdateProfile)
		protected.POST("/logout", userHandler.Logout)
		protected.GET("/info", userHandler.GetUserInfo)
		protected.PUT("/info", userHandler.UpdateUserInfo)
	}
}
