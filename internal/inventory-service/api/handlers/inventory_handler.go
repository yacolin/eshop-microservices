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

// CreateInventory 创建库存
// @Summary 创建库存
// @Description 创建一个新的库存记录
// @Tags 库存
// @Accept json
// @Produce json
// @Param inventory body dto.CreateInventoryDTO true "库存信息"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories [post]
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

// GetInventoryByID 获取库存详情
// @Summary 获取库存详情
// @Description 根据ID获取库存详细信息
// @Tags 库存
// @Produce json
// @Param id path string true "库存ID"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories/{id} [get]
func (h *InventoryHandler) GetInventoryByID(c *gin.Context) {
	id := c.Param("id")
	inventory, err := h.inventorySvc.GetInventoryByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, inventory)
}

// GetInventoryByProductID 根据产品ID获取库存
// @Summary 根据产品ID获取库存
// @Description 根据产品ID获取库存信息
// @Tags 库存
// @Produce json
// @Param productId path string true "产品ID"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories/product/{productId} [get]
func (h *InventoryHandler) GetInventoryByProductID(c *gin.Context) {
	productId := c.Param("productId")
	inventory, err := h.inventorySvc.GetInventoryByProductID(c.Request.Context(), productId)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, inventory)
}

// UpdateInventory 更新库存
// @Summary 更新库存
// @Description 根据ID更新库存信息
// @Tags 库存
// @Accept json
// @Produce json
// @Param id path string true "库存ID"
// @Param inventory body dto.UpdateInventoryDTO true "库存信息"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories/{id} [put]
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

// DeleteInventory 删除库存
// @Summary 删除库存
// @Description 根据ID删除库存
// @Tags 库存
// @Produce json
// @Param id path string true "库存ID"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/inventories/{id} [delete]
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

// ListInventories 获取库存列表
// @Summary 获取库存列表
// @Description 获取库存列表，支持分页和筛选
// @Tags 库存
// @Produce json
// @Param page query int false "页码，默认1"
// @Param size query int false "每页大小，默认10"
// @Param product_id query string false "产品ID"
// @Success 200 {object} dto.InventoryListResult "成功"
// @Router /inventory/api/v1/inventories [get]
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

// ReserveInventory 预订库存
// @Summary 预订库存
// @Description 预订指定产品的库存
// @Tags 库存操作
// @Accept json
// @Produce json
// @Param reserve body dto.ReserveInventoryDTO true "预订信息"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/inventories/reserve [post]
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

// ReleaseInventory 释放库存
// @Summary 释放库存
// @Description 释放之前预订的库存
// @Tags 库存操作
// @Accept json
// @Produce json
// @Param release body dto.ReleaseInventoryDTO true "释放信息"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/inventories/release [post]
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

// AdjustInventory 调整库存
// @Summary 调整库存
// @Description 调整指定产品的库存数量
// @Tags 库存操作
// @Accept json
// @Produce json
// @Param adjust body dto.AdjustInventoryDTO true "调整信息"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/inventories/adjust [post]
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

// CheckAvailability 检查库存可用性
// @Summary 检查库存可用性
// @Description 检查指定产品是否有足够的库存
// @Tags 库存操作
// @Produce json
// @Param product_id query string true "产品ID"
// @Param quantity query int true "需要的数量"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /inventory/api/v1/inventories/check-availability [get]
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
