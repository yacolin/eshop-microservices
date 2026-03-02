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

// CreateProduct 创建产品
// @Summary 创建产品
// @Description 创建一个新的产品
// @Tags 产品
// @Accept json
// @Produce json
// @Param product body dto.CreateProductDTO true "产品信息"
// @Success 200 {object} models.Product "成功"
// @Router /inventory/api/v1/products [post]
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

// GetProductByID 获取产品详情
// @Summary 获取产品详情
// @Description 根据ID获取产品详细信息
// @Tags 产品
// @Produce json
// @Param id path string true "产品ID"
// @Success 200 {object} models.Product "成功"
// @Router /inventory/api/v1/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.inventorySvc.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, product)
}

// UpdateProduct 更新产品
// @Summary 更新产品
// @Description 根据ID更新产品信息
// @Tags 产品
// @Accept json
// @Produce json
// @Param id path string true "产品ID"
// @Param product body dto.UpdateProductDTO true "产品信息"
// @Success 200 {object} models.Product "成功"
// @Router /inventory/api/v1/products/{id} [put]
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

// DeleteProduct 删除产品
// @Summary 删除产品
// @Description 根据ID删除产品
// @Tags 产品
// @Produce json
// @Param id path string true "产品ID"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/products/{id} [delete]
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

// ListProducts 获取产品列表
// @Summary 获取产品列表
// @Description 获取产品列表，支持分页和筛选
// @Tags 产品
// @Produce json
// @Param page query int false "页码，默认1"
// @Param size query int false "每页大小，默认10"
// @Param category_id query string false "分类ID"
// @Param name query string false "产品名称"
// @Success 200 {object} response.APIResponse "成功"
// @Router /inventory/api/v1/products [get]
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
