package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "eshop-events"
	ExchangeType = "topic"
)

// Client RabbitMQ 客户端
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	exchange string
}

// NewClient 创建 MQ 客户端并声明 topic exchange
func NewClient(url, exchange string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}
	if err := ch.ExchangeDeclare(exchange, ExchangeType, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}
	return &Client{conn: conn, channel: ch, exchange: exchange}, nil
}

// Publish 发布消息，routingKey 如 order.created
func (c *Client) Publish(ctx context.Context, routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	return c.channel.PublishWithContext(ctx, c.exchange, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	})
}

// Subscribe 订阅队列，bindingKey 如 order.*
func (c *Client) Subscribe(queueName, bindingKey string) (<-chan amqp.Delivery, error) {
	q, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}
	if err := c.channel.QueueBind(q.Name, bindingKey, c.exchange, false, nil); err != nil {
		return nil, fmt.Errorf("bind queue: %w", err)
	}
	deliveries, err := c.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return deliveries, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("close channel: %v", err)
		}
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
