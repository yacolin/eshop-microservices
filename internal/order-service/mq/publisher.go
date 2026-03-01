package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"eshop-microservices/internal/order-service/domain/models"
	"eshop-microservices/pkg/mq"
)

const sourceService = "order-service"

// Publisher 订单事件发布
type Publisher struct {
	client *mq.Client
}

// NewPublisher 创建发布者
func NewPublisher(client *mq.Client) *Publisher {
	return &Publisher{client: client}
}

// PublishOrderCreated 发布订单创建
func (p *Publisher) PublishOrderCreated(order *models.Order) {
	evt := mq.OrderCreatedEvent{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}
	body := mq.Event{
		Type:      "order.created",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "order.created", body); err != nil {
		log.Printf("publish order.created: %v", err)
	}
}

// PublishOrderUpdated 发布订单更新
func (p *Publisher) PublishOrderUpdated(id, status string) {
	evt := mq.OrderUpdatedEvent{ID: id, Status: status}
	body := mq.Event{
		Type:      "order.updated",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "order.updated", body); err != nil {
		log.Printf("publish order.updated: %v", err)
	}
}

// PublishOrderCancelled 发布订单取消
func (p *Publisher) PublishOrderCancelled(id string) {
	evt := mq.OrderCancelledEvent{ID: id}
	body := mq.Event{
		Type:      "order.cancelled",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "order.cancelled", body); err != nil {
		log.Printf("publish order.cancelled: %v", err)
	}
}

// PublishOrderCompleted 发布订单完成
func (p *Publisher) PublishOrderCompleted(id string) {
	evt := mq.OrderCompletedEvent{ID: id}
	body := mq.Event{
		Type:      "order.completed",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "order.completed", body); err != nil {
		log.Printf("publish order.completed: %v", err)
	}
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
