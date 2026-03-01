package handlers

import (
	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/mq"
	"eshop-microservices/internal/inventory-service/service"
	"eshop-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类 HTTP 处理
type CategoryHandler struct {
	inventorySvc *service.InventoryService
	publisher    *mq.Publisher
}

// NewCategoryHandler 创建分类 Handler
func NewCategoryHandler(inventorySvc *service.InventoryService, publisher *mq.Publisher) *CategoryHandler {
	return &CategoryHandler{inventorySvc: inventorySvc, publisher: publisher}
}

// CreateCategory 创建分类 POST /api/v1/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
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
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.inventorySvc.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, category)
}

// UpdateCategory 更新分类 PUT /api/v1/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
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
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
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
func (h *CategoryHandler) ListCategories(c *gin.Context) {
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