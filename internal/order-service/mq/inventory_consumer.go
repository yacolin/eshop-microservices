package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"eshop-microservices/internal/order-service/service"
	"eshop-microservices/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// InventoryConsumer 库存事件消费者
type InventoryConsumer struct {
	channel     *amqp091.Channel
	orderSvc    *service.OrderService
	exchange    string
	queueName   string
	routingKeys []string
}

// NewInventoryConsumer 创建库存事件消费者
func NewInventoryConsumer(conn *amqp091.Connection, orderSvc *service.OrderService, exchange string) (*InventoryConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明 exchange
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &InventoryConsumer{
		channel:  ch,
		orderSvc: orderSvc,
		exchange: exchange,
		routingKeys: []string{
			"inventory.reserved",      // 库存预占成功
			"inventory.reserved_failed", // 库存预占失败
			"inventory.released",      // 库存释放成功
			"inventory.deducted",      // 库存扣减成功
			"inventory.deduct_failed", // 库存扣减失败
			"inventory.low_stock",     // 库存不足警告
		},
	}, nil
}

// Start 启动消费者
func (c *InventoryConsumer) Start() error {
	// 声明队列
	q, err := c.channel.QueueDeclare(
		"order-service-inventory-queue", // name - 与订单事件队列区分开
		true,                            // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	c.queueName = q.Name

	// 绑定路由键
	for _, key := range c.routingKeys {
		err = c.channel.QueueBind(
			q.Name,   // queue name
			key,      // routing key
			c.exchange, // exchange
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue with key %s: %w", key, err)
		}
	}

	// 消费消息
	msgs, err := c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (手动确认)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := c.handleMessage(msg); err != nil {
				logger.Error("failed to handle inventory message",
					zap.String("routing_key", msg.RoutingKey),
					zap.Error(err))
				// 拒绝消息，重新入队
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	logger.Info("Inventory MQ consumer started", zap.String("queue", q.Name))
	return nil
}

// Stop 停止消费者
func (c *InventoryConsumer) Stop() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}

// handleMessage 处理消息
func (c *InventoryConsumer) handleMessage(msg amqp091.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch msg.RoutingKey {
	case "inventory.reserved":
		return c.handleInventoryReserved(ctx, msg.Body)
	case "inventory.reserved_failed":
		return c.handleInventoryReservedFailed(ctx, msg.Body)
	case "inventory.released":
		return c.handleInventoryReleased(ctx, msg.Body)
	case "inventory.deducted":
		return c.handleInventoryDeducted(ctx, msg.Body)
	case "inventory.deduct_failed":
		return c.handleInventoryDeductFailed(ctx, msg.Body)
	case "inventory.low_stock":
		return c.handleLowStock(ctx, msg.Body)
	default:
		logger.Warn("unknown inventory routing key", zap.String("key", msg.RoutingKey))
		return nil
	}
}

// InventoryReservedEvent 库存预占成功事件
type InventoryReservedEvent struct {
	OrderID       string `json:"order_id"`
	ReservationID string `json:"reservation_id"`
	ProductID     string `json:"product_id"`
	Quantity      int    `json:"quantity"`
	Timestamp     string `json:"timestamp"`
}

// handleInventoryReserved 处理库存预占成功事件
func (c *InventoryConsumer) handleInventoryReserved(ctx context.Context, body []byte) error {
	var event InventoryReservedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal inventory reserved event: %w", err)
	}

	logger.Info("handling inventory reserved event",
		zap.String("order_id", event.OrderID),
		zap.String("reservation_id", event.ReservationID),
		zap.String("product_id", event.ProductID))

	// 可以在这里更新订单的库存预占状态，或发送通知
	// 例如：更新订单项的 reservation_id

	return nil
}

// InventoryReservedFailedEvent 库存预占失败事件
type InventoryReservedFailedEvent struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

// handleInventoryReservedFailed 处理库存预占失败事件
func (c *InventoryConsumer) handleInventoryReservedFailed(ctx context.Context, body []byte) error {
	var event InventoryReservedFailedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal inventory reserved failed event: %w", err)
	}

	logger.Info("handling inventory reserved failed event",
		zap.String("order_id", event.OrderID),
		zap.String("product_id", event.ProductID),
		zap.String("reason", event.Reason))

	// 库存预占失败，需要取消订单或通知用户
	// 注意：如果使用 Saga 模式，这个逻辑已经在 Saga 补偿中处理
	// 这里可以作为额外的保险机制

	return nil
}

// InventoryReleasedEvent 库存释放成功事件
type InventoryReleasedEvent struct {
	OrderID       string `json:"order_id"`
	ReservationID string `json:"reservation_id"`
	ProductID     string `json:"product_id"`
	Quantity      int    `json:"quantity"`
	Timestamp     string `json:"timestamp"`
}

// handleInventoryReleased 处理库存释放成功事件
func (c *InventoryConsumer) handleInventoryReleased(ctx context.Context, body []byte) error {
	var event InventoryReleasedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal inventory released event: %w", err)
	}

	logger.Info("handling inventory released event",
		zap.String("order_id", event.OrderID),
		zap.String("reservation_id", event.ReservationID))

	// 库存已释放，可以更新订单状态或发送通知

	return nil
}

// InventoryDeductedEvent 库存扣减成功事件
type InventoryDeductedEvent struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Timestamp string `json:"timestamp"`
}

// handleInventoryDeducted 处理库存扣减成功事件
func (c *InventoryConsumer) handleInventoryDeducted(ctx context.Context, body []byte) error {
	var event InventoryDeductedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal inventory deducted event: %w", err)
	}

	logger.Info("handling inventory deducted event",
		zap.String("order_id", event.OrderID),
		zap.String("product_id", event.ProductID),
		zap.Int("quantity", event.Quantity))

	// 库存扣减成功，可以更新订单状态为已发货
	// 注意：如果使用 MQ 模式，这里可以触发订单状态更新

	return nil
}

// InventoryDeductFailedEvent 库存扣减失败事件
type InventoryDeductFailedEvent struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

// handleInventoryDeductFail 处理库存扣减失败事件
func (c *InventoryConsumer) handleInventoryDeductFail(ctx context.Context, body []byte) error {
	var event InventoryDeductFailedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal inventory deduct failed event: %w", err)
	}

	logger.Error("handling inventory deduct failed event",
		zap.String("order_id", event.OrderID),
		zap.String("product_id", event.ProductID),
		zap.String("reason", event.Reason))

	// 库存扣减失败，需要人工介入或补偿处理

	return nil
}

// handleInventoryDeductFailed 处理库存扣减失败事件
func (c *InventoryConsumer) handleInventoryDeductFailed(ctx context.Context, body []byte) error {
	var event InventoryDeductFailedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal inventory deduct failed event: %w", err)
	}

	logger.Error("handling inventory deduct failed event",
		zap.String("order_id", event.OrderID),
		zap.String("product_id", event.ProductID),
		zap.String("reason", event.Reason))

	// 库存扣减失败，需要人工介入或补偿处理

	return nil
}

// LowStockEvent 库存不足警告事件
type LowStockEvent struct {
	ProductID         string `json:"product_id"`
	CurrentQuantity   int    `json:"current_quantity"`
	ReservedQuantity  int    `json:"reserved_quantity"`
	Threshold         int    `json:"threshold"`
	Timestamp         string `json:"timestamp"`
}

// handleLowStock 处理库存不足警告事件
func (c *InventoryConsumer) handleLowStock(ctx context.Context, body []byte) error {
	var event LowStockEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal low stock event: %w", err)
	}

	logger.Warn("handling low stock event",
		zap.String("product_id", event.ProductID),
		zap.Int("current_quantity", event.CurrentQuantity),
		zap.Int("reserved_quantity", event.ReservedQuantity),
		zap.Int("threshold", event.Threshold))

	// 库存不足警告，可以发送通知给管理员或触发补货流程

	return nil
}
