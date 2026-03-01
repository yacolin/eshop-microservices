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

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 创建一个新的分类
// @Tags 分类
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryDTO true "分类信息"
// @Success 200 {object} models.Category "成功"
// @Router /inventory/api/v1/categories [post]
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

// GetCategoryByID 获取分类详情
// @Summary 获取分类详情
// @Description 根据ID获取分类详细信息
// @Tags 分类
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} models.Category "成功"
// @Router /inventory/api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.inventorySvc.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, category)
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Description 根据ID更新分类信息
// @Tags 分类
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Param category body dto.UpdateCategoryDTO true "分类信息"
// @Success 200 {object} models.Category "成功"
// @Router /inventory/api/v1/categories/{id} [put]
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

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 根据ID删除分类
// @Tags 分类
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/categories/{id} [delete]
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

// ListCategories 获取分类列表
// @Summary 获取分类列表
// @Description 获取分类列表，支持分页
// @Tags 分类
// @Produce json
// @Param page query int false "页码，默认1"
// @Param size query int false "每页大小，默认10"
// @Success 200 {object} dto.CategoryListResult "成功"
// @Router /inventory/api/v1/categories [get]
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
