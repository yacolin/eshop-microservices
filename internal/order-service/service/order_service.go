package service

import (
	"context"
	"fmt"

	"eshop-microservices/api/proto/inventorypb"
	"eshop-microservices/internal/order-service/api/dto"
	"eshop-microservices/internal/order-service/clients"
	"eshop-microservices/internal/order-service/domain/models"
	"eshop-microservices/internal/order-service/domain/repositories"
	ordersaga "eshop-microservices/internal/order-service/saga"
	"eshop-microservices/pkg/errcode"
	"eshop-microservices/pkg/logger"
	"eshop-microservices/pkg/query"
	"eshop-microservices/pkg/saga"

	"go.uber.org/zap"
)

// OrderService 订单业务
type OrderService struct {
	repo            repositories.OrderRepository
	inventoryClient *clients.InventoryClient
	createOrderSaga *ordersaga.CreateOrderSaga
}

// NewOrderService 创建订单服务
func NewOrderService(repo repositories.OrderRepository, inventoryClient *clients.InventoryClient, sagaLog saga.SagaLog) *OrderService {
	return &OrderService{
		repo:            repo,
		inventoryClient: inventoryClient,
		createOrderSaga: ordersaga.NewCreateOrderSaga(repo, inventoryClient, sagaLog),
	}
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

// Create 创建订单 - 使用 Saga 模式保证分布式事务
func (s *OrderService) Create(ctx context.Context, req dto.CreateOrderDTO) (*models.Order, error) {
	logger.Info("creating order with saga pattern",
		zap.String("customer_id", req.CustomerID),
		zap.Int("item_count", len(req.Items)))

	// 使用 Saga 模式执行创建订单流程
	result, err := s.createOrderSaga.Execute(ctx, req)
	if err != nil {
		logger.Error("create order saga failed",
			zap.Error(err),
			zap.String("saga_id", result.SagaID))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	logger.Info("order created successfully with saga",
		zap.String("order_id", result.Order.ID),
		zap.String("saga_id", result.SagaID),
		zap.String("reservation_id", result.ReservationID))

	return result.Order, nil
}

// CreateWithSaga 使用 Saga 模式创建订单（返回 Saga ID 用于追踪）
func (s *OrderService) CreateWithSaga(ctx context.Context, req dto.CreateOrderDTO) (*ordersaga.CreateOrderResult, error) {
	return s.createOrderSaga.Execute(ctx, req)
}

// GetSagaStatus 获取 Saga 执行状态
func (s *OrderService) GetSagaStatus(ctx context.Context, sagaID string) (*saga.Saga, error) {
	return s.createOrderSaga.GetSagaStatus(ctx, sagaID)
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

// Cancel 取消订单 - 释放预占库存
func (s *OrderService) Cancel(ctx context.Context, id string) error {
	// 获取订单信息
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}

	// 如果订单已确认且有库存客户端，释放预占库存
	if s.inventoryClient != nil && order.Status == models.OrderStatusConfirmed {
		stockItems := make([]*inventorypb.StockItem, 0, len(order.Items))
		for _, item := range order.Items {
			stockItems = append(stockItems, &inventorypb.StockItem{
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
			})
		}

		// 释放库存
		_, releaseErr := s.inventoryClient.ReleaseStock(ctx, order.ID, order.ID, stockItems)
		if releaseErr != nil {
			logger.Error("failed to release stock when cancelling order",
				zap.String("order_id", order.ID),
				zap.Error(releaseErr))
			// 继续取消订单，但记录错误
		}
	}

	return s.repo.UpdateStatus(ctx, id, models.OrderStatusCancelled)
}

// ConfirmOrder 确认订单 - 用于 MQ 消费者最终确认
func (s *OrderService) ConfirmOrder(ctx context.Context, orderID string) error {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return errcode.ErrNotFound
	}

	// 只有已确认的订单才能最终扣减
	if order.Status != models.OrderStatusConfirmed {
		return fmt.Errorf("order status is not confirmed: %s", order.Status)
	}

	// 调用 gRPC 确认扣减库存
	if s.inventoryClient != nil {
		stockItems := make([]*inventorypb.StockItem, 0, len(order.Items))
		for _, item := range order.Items {
			stockItems = append(stockItems, &inventorypb.StockItem{
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
			})
		}

		resp, err := s.inventoryClient.ConfirmDeduct(ctx, order.ID, order.ID, stockItems)
		if err != nil {
			return fmt.Errorf("failed to confirm deduct: %w", err)
		}

		if !resp.Success {
			return fmt.Errorf("confirm deduct failed: %s", resp.Message)
		}

		logger.Info("stock deducted successfully for order", zap.String("order_id", orderID))
	}

	// 更新订单状态为已发货
	return s.repo.UpdateStatus(ctx, orderID, models.OrderStatusShipped)
}
