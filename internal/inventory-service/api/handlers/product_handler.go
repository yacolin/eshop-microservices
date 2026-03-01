package handlers

import (
	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/mq"
	"eshop-microservices/internal/inventory-service/service"
	"eshop-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

// ProductHandler 产品 HTTP 处理
type ProductHandler struct {
	inventorySvc *service.InventoryService
	publisher    *mq.Publisher
}

// NewProductHandler 创建产品 Handler
func NewProductHandler(inventorySvc *service.InventoryService, publisher *mq.Publisher) *ProductHandler {
	return &ProductHandler{inventorySvc: inventorySvc, publisher: publisher}
}

// CreateProduct 创建产品 POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
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
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.inventorySvc.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, product)
}

// UpdateProduct 更新产品 PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
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
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
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
func (h *ProductHandler) ListProducts(c *gin.Context) {
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