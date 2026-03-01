package routes

import (
	"eshop-microservices/internal/inventory-service/api/handlers"
	"eshop-microservices/internal/inventory-service/api/middleware"

	"github.com/gin-gonic/gin"
)

// Setup 注册路由
// 版本策略：v1 / v2 并行运行，待 v2 稳定后再移除旧版本注册即可
func Setup(r *gin.Engine, inventoryHandler *handlers.InventoryHandler) {
	r.Use(middleware.Recovery(), middleware.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	registerV1(api, inventoryHandler)
	// 后续增加 v2：复制 registerV1 为 registerV2，改 handler/逻辑，在此处加一行 registerV2(api, inventoryHandler)
	// 下线 v1：删掉本行 registerV1 及下方 registerV1 函数即可
}

func registerV1(api *gin.RouterGroup, inventoryHandler *handlers.InventoryHandler) {
	// 产品相关路由：/api/v1/products（支持带或不带末尾 /）
	products := api.Group("/v1/products")
	{
		products.POST("", inventoryHandler.CreateProduct)
		products.GET("", inventoryHandler.ListProducts)
		products.GET("/:id", inventoryHandler.GetProductByID)
		products.PUT("/:id", inventoryHandler.UpdateProduct)
		products.DELETE("/:id", inventoryHandler.DeleteProduct)
	}

	// 库存相关路由：/api/v1/inventories（支持带或不带末尾 /）
	inventories := api.Group("/v1/inventories")
	{
		inventories.POST("", inventoryHandler.CreateInventory)
		inventories.GET("", inventoryHandler.ListInventories)
		inventories.GET("/:id", inventoryHandler.GetInventoryByID)
		inventories.PUT("/:id", inventoryHandler.UpdateInventory)
		inventories.DELETE("/:id", inventoryHandler.DeleteInventory)

		// 库存操作相关
		inventories.POST("/reserve", inventoryHandler.ReserveInventory)
		inventories.POST("/release", inventoryHandler.ReleaseInventory)
		inventories.POST("/adjust", inventoryHandler.AdjustInventory)
		inventories.GET("/check-availability", inventoryHandler.CheckAvailability)
		inventories.GET("/product/:productId", inventoryHandler.GetInventoryByProductID)
	}

	// 分类相关路由：/api/v1/categories（支持带或不带末尾 /）
	categories := api.Group("/v1/categories")
	{
		categories.POST("", inventoryHandler.CreateCategory)
		categories.GET("", inventoryHandler.ListCategories)
		categories.GET("/:id", inventoryHandler.GetCategoryByID)
		categories.PUT("/:id", inventoryHandler.UpdateCategory)
		categories.DELETE("/:id", inventoryHandler.DeleteCategory)
	}
}
