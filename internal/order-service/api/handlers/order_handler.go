package handlers

import (
	"eshop-microservices/internal/order-service/api/dto"
	"eshop-microservices/internal/order-service/domain/models"
	"eshop-microservices/internal/order-service/mq"
	"eshop-microservices/internal/order-service/service"
	"eshop-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

// OrderHandler 订单 HTTP 处理
type OrderHandler struct {
	orderSvc  *service.OrderService
	publisher *mq.Publisher
}

// NewOrderHandler 创建订单 Handler
func NewOrderHandler(orderSvc *service.OrderService, publisher *mq.Publisher) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc, publisher: publisher}
}

// Create 创建订单 POST /api/orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	order, err := h.orderSvc.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishOrderCreated(order)
	}
	response.Success(c, order)
}

// List 订单列表 GET /api/orders
func (h *OrderHandler) List(c *gin.Context) {
	var q dto.OrderListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.orderSvc.List(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)

}

// GetByID 订单详情 GET /api/orders/:id
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, order)
}

// UpdateStatus 更新订单状态 PUT /api/orders/:id
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}
	if err := h.orderSvc.UpdateStatus(c.Request.Context(), id, body.Status); err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishOrderUpdated(id, body.Status)
		if body.Status == models.OrderStatusCompleted {
			h.publisher.PublishOrderCompleted(id)
		}
	}
	response.Success(c, gin.H{"message": "updated"})
}

// Cancel 取消订单 DELETE /api/orders/:id
func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	if err := h.orderSvc.Cancel(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	if h.publisher != nil {
		h.publisher.PublishOrderCancelled(id)
	}
	response.Success(c, gin.H{"message": "cancelled"})
}
