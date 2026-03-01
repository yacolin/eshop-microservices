package service

import (
	"context"
	"fmt"

	"eshop-microservices/internal/order-service/api/dto"
	"eshop-microservices/internal/order-service/domain/models"
	"eshop-microservices/internal/order-service/domain/repositories"
	"eshop-microservices/pkg/errcode"
	"eshop-microservices/pkg/query"
)

// OrderService 订单业务
type OrderService struct {
	repo repositories.OrderRepository
}

// NewOrderService 创建订单服务
func NewOrderService(repo repositories.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CustomerID string               `json:"customer_id" binding:"required"`
	Currency   string               `json:"currency"` // 可选，默认 CNY
	Items      []CreateOrderItemReq `json:"items" binding:"required,min=1,dive"`
}

// CreateOrderItemReq 订单项
type CreateOrderItemReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	UnitPrice int64  `json:"unit_price" binding:"required,min=0"` // 单价，单位：分
}

// Create 创建订单
func (s *OrderService) Create(ctx context.Context, req dto.CreateOrderDTO) (*models.Order, error) {
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}
	order := &models.Order{
		CustomerID: req.CustomerID,
		Currency:   currency,
		Status:     models.OrderStatusPending,
	}
	var total int64
	for _, it := range req.Items {
		amount := it.UnitPrice * int64(it.Quantity)
		order.Items = append(order.Items, models.OrderItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			Amount:    amount,
		})
		total += amount
	}
	order.TotalAmount = total
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

// GetByID 获取订单详情
func (s *OrderService) GetByID(ctx context.Context, id string) (*models.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return o, nil
}

// List 订单列表
func (s *OrderService) List(ctx context.Context, q dto.OrderListQuery) (*query.ListResult[models.Order], error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 10
	}
	offset := (q.Page - 1) * q.Size

	list, err := s.repo.ListByQuery(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountByQuery(ctx, q)
	if err != nil {
		return nil, err
	}

	return &query.ListResult[models.Order]{
		List:  list,
		Total: total,
	}, nil
}

// UpdateStatus 更新订单状态
func (s *OrderService) UpdateStatus(ctx context.Context, id, status string) error {
	valid := map[string]bool{
		models.OrderStatusPending: true, models.OrderStatusConfirmed: true,
		models.OrderStatusShipped: true, models.OrderStatusCompleted: true, models.OrderStatusCancelled: true,
	}
	if !valid[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

// Cancel 取消订单（软删或改状态，这里用更新状态）
func (s *OrderService) Cancel(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, models.OrderStatusCancelled)
}
