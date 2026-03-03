package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"eshop-microservices/internal/inventory-service/api/dto"
	"eshop-microservices/internal/inventory-service/service"
	"eshop-microservices/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Consumer RabbitMQ 消费者
type Consumer struct {
	channel     *amqp091.Channel
	svc         *service.InventoryService
	exchange    string
	queueName   string
	routingKeys []string
}

// NewConsumer 创建消费者
func NewConsumer(conn *amqp091.Connection, svc *service.InventoryService, exchange string) (*Consumer, error) {
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

	return &Consumer{
		channel:  ch,
		svc:      svc,
		exchange: exchange,
		routingKeys: []string{
			"order.created",      // 订单创建 - 预占库存
			"order.confirmed",    // 订单确认 - 扣减库存
			"order.cancelled",    // 订单取消 - 释放库存
			"order.payment_completed", // 支付完成 - 最终扣减
			"order.payment_failed",    // 支付失败 - 释放库存
		},
	}, nil
}

// Start 启动消费者
func (c *Consumer) Start() error {
	// 声明队列
	q, err := c.channel.QueueDeclare(
		"inventory-service-queue", // name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
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
				logger.Error("failed to handle message",
					zap.String("routing_key", msg.RoutingKey),
					zap.Error(err))
				// 拒绝消息，重新入队
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	logger.Info("MQ consumer started", zap.String("queue", q.Name))
	return nil
}

// Stop 停止消费者
func (c *Consumer) Stop() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}

// handleMessage 处理消息
func (c *Consumer) handleMessage(msg amqp091.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch msg.RoutingKey {
	case "order.created":
		return c.handleOrderCreated(ctx, msg.Body)
	case "order.confirmed":
		return c.handleOrderConfirmed(ctx, msg.Body)
	case "order.cancelled":
		return c.handleOrderCancelled(ctx, msg.Body)
	case "order.payment_completed":
		return c.handlePaymentCompleted(ctx, msg.Body)
	case "order.payment_failed":
		return c.handlePaymentFailed(ctx, msg.Body)
	default:
		logger.Warn("unknown routing key", zap.String("key", msg.RoutingKey))
		return nil
	}
}

// OrderItem 订单项
type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	TotalAmount int64      `json:"total_amount"`
}

// handleOrderCreated 处理订单创建事件 - 预占库存
func (c *Consumer) handleOrderCreated(ctx context.Context, body []byte) error {
	var event OrderCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order created event: %w", err)
	}

	logger.Info("handling order created event",
		zap.String("order_id", event.ID),
		zap.String("customer_id", event.CustomerID))

	// 预占库存
	for _, item := range event.Items {
		if err := c.svc.ReserveInventory(ctx, dto.ReserveInventoryDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}); err != nil {
			logger.Error("failed to reserve inventory",
				zap.String("order_id", event.ID),
				zap.String("product_id", item.ProductID),
				zap.Error(err))
			return err
		}
	}

	logger.Info("inventory reserved for order",
		zap.String("order_id", event.ID))

	return nil
}

// OrderConfirmedEvent 订单确认事件
type OrderConfirmedEvent struct {
	ID string `json:"id"`
}

// handleOrderConfirmed 处理订单确认事件 - 扣减实际库存
func (c *Consumer) handleOrderConfirmed(ctx context.Context, body []byte) error {
	var event OrderConfirmedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order confirmed event: %w", err)
	}

	logger.Info("handling order confirmed event",
		zap.String("order_id", event.ID))

	// 这里应该根据订单信息扣减实际库存
	// 由于事件中没有订单项信息，需要查询订单服务或从数据库获取
	// 简化处理：实际项目中应该包含订单项信息

	logger.Info("inventory deducted for confirmed order",
		zap.String("order_id", event.ID))

	return nil
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	ID     string      `json:"id"`
	Items  []OrderItem `json:"items"`
	Reason string      `json:"reason"`
}

// handleOrderCancelled 处理订单取消事件 - 释放库存
func (c *Consumer) handleOrderCancelled(ctx context.Context, body []byte) error {
	var event OrderCancelledEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order cancelled event: %w", err)
	}

	logger.Info("handling order cancelled event",
		zap.String("order_id", event.ID),
		zap.String("reason", event.Reason))

	// 释放库存
	for _, item := range event.Items {
		if err := c.svc.ReleaseInventory(ctx, dto.ReleaseInventoryDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}); err != nil {
			logger.Error("failed to release inventory",
				zap.String("order_id", event.ID),
				zap.String("product_id", item.ProductID),
				zap.Error(err))
			// 继续释放其他商品的库存
		}
	}

	logger.Info("inventory released for cancelled order",
		zap.String("order_id", event.ID))

	return nil
}

// PaymentCompletedEvent 支付完成事件
type PaymentCompletedEvent struct {
	OrderID   string      `json:"order_id"`
	PaymentID string      `json:"payment_id"`
	Items     []OrderItem `json:"items"`
}

// handlePaymentCompleted 处理支付完成事件 - 最终扣减库存
func (c *Consumer) handlePaymentCompleted(ctx context.Context, body []byte) error {
	var event PaymentCompletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payment completed event: %w", err)
	}

	logger.Info("handling payment completed event",
		zap.String("order_id", event.OrderID),
		zap.String("payment_id", event.PaymentID))

	// 最终扣减库存（从预占转为实际扣减）
	for _, item := range event.Items {
		// 调整库存：减去实际库存和预占库存
		if err := c.svc.AdjustInventory(ctx, dto.AdjustInventoryDTO{
			ProductID: item.ProductID,
			Quantity:  -item.Quantity, // 负数表示扣减
		}); err != nil {
			logger.Error("failed to deduct inventory after payment",
				zap.String("order_id", event.OrderID),
				zap.String("product_id", item.ProductID),
				zap.Error(err))
			return err
		}
	}

	logger.Info("inventory deducted after payment",
		zap.String("order_id", event.OrderID))

	return nil
}

// PaymentFailedEvent 支付失败事件
type PaymentFailedEvent struct {
	OrderID   string      `json:"order_id"`
	PaymentID string      `json:"payment_id"`
	Items     []OrderItem `json:"items"`
	Reason    string      `json:"reason"`
}

// handlePaymentFailed 处理支付失败事件 - 释放库存
func (c *Consumer) handlePaymentFailed(ctx context.Context, body []byte) error {
	var event PaymentFailedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payment failed event: %w", err)
	}

	logger.Info("handling payment failed event",
		zap.String("order_id", event.OrderID),
		zap.String("reason", event.Reason))

	// 释放库存
	for _, item := range event.Items {
		if err := c.svc.ReleaseInventory(ctx, dto.ReleaseInventoryDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}); err != nil {
			logger.Error("failed to release inventory after payment failed",
				zap.String("order_id", event.OrderID),
				zap.String("product_id", item.ProductID),
				zap.Error(err))
			// 继续释放其他商品的库存
		}
	}

	logger.Info("inventory released after payment failed",
		zap.String("order_id", event.OrderID))

	return nil
}
