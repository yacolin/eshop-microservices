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

// Create 创建订单
// @Summary 创建订单
// @Description 创建一个新的订单
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderDTO true "订单信息"
// @Success 200 {object} models.Order "成功"
// @Router /order/api/v1/orders [post]
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

// List 订单列表 GET
// @Summary 获取订单列表
// @Description 获取订单列表，支持分页和筛选
// @Tags orders
// @Produce json
// @Param customer_id query int64 false "客户ID"
// @Param status query string false "订单状态"
// @Param min_price query float64 false "价格区间下限"
// @Param max_price query float64 false "价格区间上限"
// @Param sort_by query string false "排序字段，例如 total_amount, created_at"
// @Param order query string false "排序方向，asc 或 desc"
// @Param page query int false "页码，默认1"
// @Param size query int false "每页大小，默认10"
// @Success 200 {object} response.APIResponse "成功"
// @Router /order/api/v1/orders [get]
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

// GetByID 订单详情 GET
// @Summary 获取订单详情
// @Description 根据ID获取订单详细信息
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} models.Order "成功"
// @Router /order/api/v1/orders/{id} [get]
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, order)
}

// UpdateStatus 更新订单状态
// @Summary 更新订单状态
// @Description 根据ID更新订单状态
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param status body object{status string} true "订单状态"
// @Success 200 {object} map[string]string "成功"
// @Router /order/api/v1/orders/{id} [put]
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

// Cancel 取消订单
// @Summary 取消订单
// @Description 根据ID取消订单
// @Tags orders
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} map[string]string "成功"
// @Router /order/api/v1/orders/{id} [delete]
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

// GetSagaStatus 获取 Saga 执行状态
// @Summary 获取 Saga 执行状态
// @Description 根据 Saga ID 获取分布式事务执行状态
// @Tags orders
// @Produce json
// @Param saga_id path string true "Saga ID"
// @Success 200 {object} saga.Saga "成功"
// @Router /order/api/v1/orders/saga/{saga_id} [get]
func (h *OrderHandler) GetSagaStatus(c *gin.Context) {
	sagaID := c.Param("saga_id")
	saga, err := h.orderSvc.GetSagaStatus(c.Request.Context(), sagaID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, saga)
}
