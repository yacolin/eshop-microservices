package routes

import (
	"eshop-microservices/internal/order-service/api/handlers"
	"eshop-microservices/internal/order-service/api/middleware"

	"github.com/gin-gonic/gin"
)

// Setup 注册路由
// 版本策略：v1 / v2 并行运行，待 v2 稳定后再移除旧版本注册即可
func Setup(r *gin.Engine, orderHandler *handlers.OrderHandler) {
	r.Use(middleware.Recovery(), middleware.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "order ok"})
	})

	api := r.Group("/api")
	registerV1(api, orderHandler)
	// 后续增加 v2：复制 registerV1 为 registerV2，改 handler/逻辑，在此处加一行 registerV2(api, orderHandler)
	// 下线 v1：删掉本行 registerV1 及下方 registerV1 函数即可
}

func registerV1(api *gin.RouterGroup, orderHandler *handlers.OrderHandler) {
	// 统一前缀：/api/v1/orders（支持带或不带末尾 /）
	orders := api.Group("/v1/orders")
	{
		orders.POST("", orderHandler.Create)
		orders.GET("", orderHandler.List)
		orders.GET("/:id", orderHandler.GetByID)
		orders.PUT("/:id", orderHandler.UpdateStatus)
		orders.DELETE("/:id", orderHandler.Cancel)
	}
}
