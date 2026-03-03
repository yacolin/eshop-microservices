package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"eshop-microservices/internal/order-service/service"
	"eshop-microservices/pkg/logger"
	"eshop-microservices/pkg/mq"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Consumer RabbitMQ 消费者
type Consumer struct {
	channel     *amqp091.Channel
	orderSvc    *service.OrderService
	exchange    string
	queueName   string
	routingKeys []string
}

// NewConsumer 创建消费者
func NewConsumer(conn *amqp091.Connection, orderSvc *service.OrderService, exchange string) (*Consumer, error) {
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
		orderSvc: orderSvc,
		exchange: exchange,
		routingKeys: []string{
			"order.created",
			"order.cancelled",
			"payment.completed",
			"payment.failed",
		},
	}, nil
}

// Start 启动消费者
func (c *Consumer) Start() error {
	// 声明队列
	q, err := c.channel.QueueDeclare(
		"order-service-queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
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
	case "order.cancelled":
		return c.handleOrderCancelled(ctx, msg.Body)
	case "payment.completed":
		return c.handlePaymentCompleted(ctx, msg.Body)
	case "payment.failed":
		return c.handlePaymentFailed(ctx, msg.Body)
	default:
		logger.Warn("unknown routing key", zap.String("key", msg.RoutingKey))
		return nil
	}
}

// handleOrderCreated 处理订单创建事件
func (c *Consumer) handleOrderCreated(ctx context.Context, body []byte) error {
	var event mq.OrderCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order created event: %w", err)
	}

	logger.Info("handling order created event",
		zap.String("order_id", event.ID),
		zap.String("customer_id", event.CustomerID))

	// 订单创建时已经通过 gRPC 预占了库存
	// 这里可以添加其他业务逻辑，如发送通知等

	return nil
}

// handleOrderCancelled 处理订单取消事件
func (c *Consumer) handleOrderCancelled(ctx context.Context, body []byte) error {
	var event mq.OrderCancelledEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order cancelled event: %w", err)
	}

	logger.Info("handling order cancelled event",
		zap.String("order_id", event.ID))

	// 订单取消时释放库存已经在 Cancel 方法中处理
	// 这里可以添加其他业务逻辑

	return nil
}

// handlePaymentCompleted 处理支付完成事件 - 最终扣减库存
func (c *Consumer) handlePaymentCompleted(ctx context.Context, body []byte) error {
	var event struct {
		OrderID       string    `json:"order_id"`
		PaymentID     string    `json:"payment_id"`
		Amount        int64     `json:"amount"`
		Status        string    `json:"status"`
		PaidAt        time.Time `json:"paid_at"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payment completed event: %w", err)
	}

	logger.Info("handling payment completed event",
		zap.String("order_id", event.OrderID),
		zap.String("payment_id", event.PaymentID))

	// 支付完成，确认扣减库存
	if err := c.orderSvc.ConfirmOrder(ctx, event.OrderID); err != nil {
		return fmt.Errorf("failed to confirm order after payment: %w", err)
	}

	logger.Info("order confirmed after payment",
		zap.String("order_id", event.OrderID))

	return nil
}

// handlePaymentFailed 处理支付失败事件 - 释放库存
func (c *Consumer) handlePaymentFailed(ctx context.Context, body []byte) error {
	var event struct {
		OrderID   string `json:"order_id"`
		PaymentID string `json:"payment_id"`
		Reason    string `json:"reason"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payment failed event: %w", err)
	}

	logger.Info("handling payment failed event",
		zap.String("order_id", event.OrderID),
		zap.String("reason", event.Reason))

	// 支付失败，取消订单并释放库存
	if err := c.orderSvc.Cancel(ctx, event.OrderID); err != nil {
		return fmt.Errorf("failed to cancel order after payment failed: %w", err)
	}

	logger.Info("order cancelled after payment failed",
		zap.String("order_id", event.OrderID))

	return nil
}
