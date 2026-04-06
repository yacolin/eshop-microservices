package routes

import (
	"eshop-microservices/internal/inventory-service/api/handlers"
	pkgmiddleware "eshop-microservices/pkg/middleware"

	"github.com/gin-gonic/gin"

	_ "eshop-microservices/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 注册路由
// 版本策略：v1 / v2 并行运行，待 v2 稳定后再移除旧版本注册即可
func Setup(
	r *gin.Engine,
	productHandler *handlers.ProductHandler,
	inventoryHandler *handlers.InventoryHandler,
	categoryHandler *handlers.CategoryHandler,
) {
	r.Use(pkgmiddleware.Recovery(), pkgmiddleware.Logger(), pkgmiddleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "inventory ok"})
	})

	// Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.StaticFile("/swagger-json", "./docs/swagger.json")

	api := r.Group("/api")
	registerV1(api, productHandler, inventoryHandler, categoryHandler)
	// 后续增加 v2：复制 registerV1 为 registerV2，改 handler/逻辑，在此处加一行 registerV2(api, productHandler, inventoryHandler, categoryHandler)
	// 下线 v1：删掉本行 registerV1 及下方 registerV1 函数即可
}

func registerV1(
	api *gin.RouterGroup,
	productHandler *handlers.ProductHandler,
	inventoryHandler *handlers.InventoryHandler,
	categoryHandler *handlers.CategoryHandler,
) {
	// 产品相关路由：/api/v1/products（支持带或不带末尾 /）
	products := api.Group("/v1/products")
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.ListProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
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
		categories.POST("", categoryHandler.CreateCategory)
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// 评论相关路由：/api/v1/comments（支持带或不带末尾 /）
	comments := api.Group("/v1/comments")
	comments.Use(pkgmiddleware.JWTAuth())
	{
		comments.POST("", productHandler.CreateComment)
		comments.GET("", productHandler.ListComments)
		comments.DELETE("/:id", productHandler.DeleteComment)
	}
}
