package saga

import (
	"context"
	"encoding/json"
	"fmt"

	"eshop-microservices/api/proto/inventorypb"
	"eshop-microservices/internal/order-service/api/dto"
	"eshop-microservices/internal/order-service/clients"
	"eshop-microservices/internal/order-service/domain/models"
	"eshop-microservices/internal/order-service/domain/repositories"
	"eshop-microservices/pkg/logger"
	"eshop-microservices/pkg/saga"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateOrderSaga 创建订单 Saga
type CreateOrderSaga struct {
	coordinator     *saga.Coordinator
	repo            repositories.OrderRepository
	inventoryClient *clients.InventoryClient
}

// CreateOrderResult 创建订单结果
type CreateOrderResult struct {
	Order         *models.Order
	SagaID        string
	ReservationID string
}

// NewCreateOrderSaga 创建订单 Saga 协调器
func NewCreateOrderSaga(repo repositories.OrderRepository, inventoryClient *clients.InventoryClient, log saga.SagaLog) *CreateOrderSaga {
	return &CreateOrderSaga{
		coordinator:     saga.NewCoordinator(log),
		repo:            repo,
		inventoryClient: inventoryClient,
	}
}

// Execute 执行创建订单 Saga
func (s *CreateOrderSaga) Execute(ctx context.Context, req dto.CreateOrderDTO) (*CreateOrderResult, error) {
	logger.Info("starting create order saga",
		zap.String("customer_id", req.CustomerID),
		zap.Int("item_count", len(req.Items)))

	// 创建 Saga 实例
	sg := saga.NewSaga("create_order")

	// 准备数据
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}

	// 构建订单对象（但不保存）
	order := &models.Order{
		CustomerID: req.CustomerID,
		Currency:   currency,
		Status:     models.OrderStatusPending,
	}

	var total int64
	for _, it := range req.Items {
		amount := it.UnitPrice * int64(it.Quantity)
		order.Items = append(order.Items, models.OrderItem{
			ID:        uuid.New().String(),
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			Amount:    amount,
		})
		total += amount
	}
	order.TotalAmount = total

	// 构建库存请求
	stockItems := make([]*inventorypb.StockItem, 0, len(req.Items))
	for _, item := range req.Items {
		stockItems = append(stockItems, &inventorypb.StockItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}

	// 存储数据到 Saga
	sg.SetData("order", order)
	sg.SetData("stock_items", stockItems)
	sg.SetData("reservation_id", "")

	// 步骤 1: 创建订单
	sg.AddStep(
		"create_order",
		// 正向操作：创建订单
		func(ctx context.Context) error {
			orderData, _ := sg.GetData("order")
			order := orderData.(*models.Order)

			if err := s.repo.Create(ctx, order); err != nil {
				return fmt.Errorf("failed to create order: %w", err)
			}

			logger.Info("order created in saga",
				zap.String("order_id", order.ID))
			return nil
		},
		// 补偿操作：删除订单
		func(ctx context.Context) error {
			orderData, _ := sg.GetData("order")
			order := orderData.(*models.Order)

			logger.Info("compensating: deleting order",
				zap.String("order_id", order.ID))

			if err := s.repo.Delete(ctx, order.ID); err != nil {
				logger.Error("failed to delete order in compensation",
					zap.String("order_id", order.ID),
					zap.Error(err))
				return err
			}
			return nil
		},
	)

	// 步骤 2: 预占库存（如果有库存客户端）
	if s.inventoryClient != nil {
		sg.AddStep(
			"reserve_stock",
			// 正向操作：预占库存
			func(ctx context.Context) error {
				orderData, _ := sg.GetData("order")
				order := orderData.(*models.Order)

				stockItemsData, _ := sg.GetData("stock_items")
				stockItems := stockItemsData.([]*inventorypb.StockItem)

				resp, err := s.inventoryClient.ReserveStock(ctx, order.ID, stockItems)
				if err != nil {
					return fmt.Errorf("failed to reserve stock: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("stock reserve failed: %s", resp.Message)
				}

				// 保存预留ID
				sg.SetData("reservation_id", resp.ReservationId)

				logger.Info("stock reserved in saga",
					zap.String("order_id", order.ID),
					zap.String("reservation_id", resp.ReservationId))

				return nil
			},
			// 补偿操作：释放库存
			func(ctx context.Context) error {
				orderData, _ := sg.GetData("order")
				order := orderData.(*models.Order)

				stockItemsData, _ := sg.GetData("stock_items")
				stockItems := stockItemsData.([]*inventorypb.StockItem)

				reservationIDData, _ := sg.GetData("reservation_id")
				reservationID := reservationIDData.(string)

				logger.Info("compensating: releasing stock",
					zap.String("order_id", order.ID),
					zap.String("reservation_id", reservationID))

				_, err := s.inventoryClient.ReleaseStock(ctx, order.ID, reservationID, stockItems)
				if err != nil {
					logger.Error("failed to release stock in compensation",
						zap.String("order_id", order.ID),
						zap.Error(err))
					return err
				}
				return nil
			},
		)

		// 步骤 3: 更新订单状态为已确认
		sg.AddStep(
			"confirm_order",
			// 正向操作：更新订单状态
			func(ctx context.Context) error {
				orderData, _ := sg.GetData("order")
				order := orderData.(*models.Order)

				if err := s.repo.UpdateStatus(ctx, order.ID, models.OrderStatusConfirmed); err != nil {
					return fmt.Errorf("failed to confirm order: %w", err)
				}

				order.Status = models.OrderStatusConfirmed

				logger.Info("order confirmed in saga",
					zap.String("order_id", order.ID))
				return nil
			},
			// 补偿操作：恢复订单状态为待处理
			func(ctx context.Context) error {
				orderData, _ := sg.GetData("order")
				order := orderData.(*models.Order)

				logger.Info("compensating: reverting order status",
					zap.String("order_id", order.ID))

				if err := s.repo.UpdateStatus(ctx, order.ID, models.OrderStatusPending); err != nil {
					logger.Error("failed to revert order status in compensation",
						zap.String("order_id", order.ID),
						zap.Error(err))
					return err
				}
				return nil
			},
		)
	}

	// 执行 Saga
	if err := s.coordinator.Execute(ctx, sg); err != nil {
		logger.Error("create order saga failed",
			zap.Error(err),
			zap.String("saga_id", sg.ID))
		return nil, err
	}

	// 获取最终结果
	orderData, _ := sg.GetData("order")
	order = orderData.(*models.Order)

	reservationIDData, _ := sg.GetData("reservation_id")
	reservationID := reservationIDData.(string)

	logger.Info("create order saga completed successfully",
		zap.String("order_id", order.ID),
		zap.String("saga_id", sg.ID))

	return &CreateOrderResult{
		Order:         order,
		SagaID:        sg.ID,
		ReservationID: reservationID,
	}, nil
}

// GetSagaStatus 获取 Saga 状态
func (s *CreateOrderSaga) GetSagaStatus(ctx context.Context, sagaID string) (*saga.Saga, error) {
	return s.coordinator.GetSaga(ctx, sagaID)
}

// CreateOrderSagaData Saga 数据序列化辅助函数
func CreateOrderSagaData(order *models.Order, stockItems []*inventorypb.StockItem) map[string]interface{} {
	return map[string]interface{}{
		"order":        order,
		"stock_items":  stockItems,
		"reservation_id": "",
	}
}

// SerializeOrder 序列化订单
func SerializeOrder(order *models.Order) ([]byte, error) {
	return json.Marshal(order)
}

// DeserializeOrder 反序列化订单
func DeserializeOrder(data []byte) (*models.Order, error) {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return &order, nil
}
