package handlers

import (
	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/mq"
	"eshop-microservices/internal/inventory-service/service"
	"eshop-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

// InventoryHandler 库存 HTTP 处理
type InventoryHandler struct {
	inventorySvc *service.InventoryService
	publisher    *mq.Publisher
}

// NewInventoryHandler 创建库存 Handler
func NewInventoryHandler(inventorySvc *service.InventoryService, publisher *mq.Publisher) *InventoryHandler {
	return &InventoryHandler{inventorySvc: inventorySvc, publisher: publisher}
}

// CreateProduct 创建产品 POST /api/v1/products
func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	product, err := h.inventorySvc.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishProductCreated(product)
	}
	response.Success(c, product)
}

// GetProductByID 产品详情 GET /api/v1/products/:id
func (h *InventoryHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.inventorySvc.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, product)
}

// UpdateProduct 更新产品 PUT /api/v1/products/:id
func (h *InventoryHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	product, err := h.inventorySvc.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishProductUpdated(product)
	}
	response.Success(c, product)
}

// DeleteProduct 删除产品 DELETE /api/v1/products/:id
func (h *InventoryHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.inventorySvc.DeleteProduct(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishProductDeleted(id)
	}
	response.Success(c, gin.H{"message": "deleted"})
}

// ListProducts 产品列表 GET /api/v1/products
func (h *InventoryHandler) ListProducts(c *gin.Context) {
	var q dto.ProductListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.inventorySvc.ListProducts(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// CreateInventory 创建库存 POST /api/v1/inventories
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.CreateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	inventory, err := h.inventorySvc.CreateInventory(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishInventoryCreated(inventory)
	}
	response.Success(c, inventory)
}

// GetInventoryByID 库存详情 GET /api/v1/inventories/:id
func (h *InventoryHandler) GetInventoryByID(c *gin.Context) {
	id := c.Param("id")
	inventory, err := h.inventorySvc.GetInventoryByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, inventory)
}

// GetInventoryByProductID 根据产品ID获取库存 GET /api/v1/inventories/product/:productId
func (h *InventoryHandler) GetInventoryByProductID(c *gin.Context) {
	productId := c.Param("productId")
	inventory, err := h.inventorySvc.GetInventoryByProductID(c.Request.Context(), productId)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, inventory)
}

// UpdateInventory 更新库存 PUT /api/v1/inventories/:id
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	inventory, err := h.inventorySvc.UpdateInventory(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishInventoryUpdated(inventory)
	}
	response.Success(c, inventory)
}

// DeleteInventory 删除库存 DELETE /api/v1/inventories/:id
func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	id := c.Param("id")
	if err := h.inventorySvc.DeleteInventory(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishInventoryDeleted(id)
	}
	response.Success(c, gin.H{"message": "deleted"})
}

// ListInventories 库存列表 GET /api/v1/inventories
func (h *InventoryHandler) ListInventories(c *gin.Context) {
	var q dto.InventoryListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.inventorySvc.ListInventories(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ReserveInventory 预订库存 POST /api/v1/inventories/reserve
func (h *InventoryHandler) ReserveInventory(c *gin.Context) {
	var req dto.ReserveInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.inventorySvc.ReserveInventory(c.Request.Context(), req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "reserved"})
}

// ReleaseInventory 释放库存 POST /api/v1/inventories/release
func (h *InventoryHandler) ReleaseInventory(c *gin.Context) {
	var req dto.ReleaseInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.inventorySvc.ReleaseInventory(c.Request.Context(), req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "released"})
}

// AdjustInventory 调整库存 POST /api/v1/inventories/adjust
func (h *InventoryHandler) AdjustInventory(c *gin.Context) {
	var req dto.AdjustInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.inventorySvc.AdjustInventory(c.Request.Context(), req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "adjusted"})
}

// CheckAvailability 检查库存可用性 GET /api/v1/inventories/check-availability
func (h *InventoryHandler) CheckAvailability(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id" form:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" form:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		return
	}

	available, err := h.inventorySvc.CheckAvailability(c.Request.Context(), req.ProductID, req.Quantity)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"available":  available,
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
	})
}

// CreateCategory 创建分类 POST /api/v1/categories
func (h *InventoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	category, err := h.inventorySvc.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishCategoryCreated(category)
	}
	response.Success(c, category)
}

// GetCategoryByID 分类详情 GET /api/v1/categories/:id
func (h *InventoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.inventorySvc.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, category)
}

// UpdateCategory 更新分类 PUT /api/v1/categories/:id
func (h *InventoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	category, err := h.inventorySvc.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishCategoryUpdated(category)
	}
	response.Success(c, category)
}

// DeleteCategory 删除分类 DELETE /api/v1/categories/:id
func (h *InventoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.inventorySvc.DeleteCategory(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishCategoryDeleted(id)
	}
	response.Success(c, gin.H{"message": "deleted"})
}

// ListCategories 分类列表 GET /api/v1/categories
func (h *InventoryHandler) ListCategories(c *gin.Context) {
	var q dto.CategoryListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.inventorySvc.ListCategories(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}
