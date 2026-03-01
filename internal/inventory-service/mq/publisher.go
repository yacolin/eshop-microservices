package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"eshop-microservices/internal/inventory-service/domain/models"
	"eshop-microservices/pkg/mq"
)

const sourceService = "inventory-service"

// Publisher 库存事件发布
type Publisher struct {
	client *mq.Client
}

// NewPublisher 创建发布者
func NewPublisher(client *mq.Client) *Publisher {
	return &Publisher{client: client}
}

// PublishProductCreated 发布产品创建
func (p *Publisher) PublishProductCreated(product *models.Product) {
	evt := mq.ProductCreatedEvent{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
	}
	body := mq.Event{
		Type:      "product.created",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "product.created", body); err != nil {
		log.Printf("publish product.created: %v", err)
	}
}

// PublishProductUpdated 发布产品更新
func (p *Publisher) PublishProductUpdated(product *models.Product) {
	evt := mq.ProductUpdatedEvent{ID: product.ID, Name: product.Name, Price: product.Price}
	body := mq.Event{
		Type:      "product.updated",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "product.updated", body); err != nil {
		log.Printf("publish product.updated: %v", err)
	}
}

// PublishProductDeleted 发布产品删除
func (p *Publisher) PublishProductDeleted(id string) {
	evt := mq.ProductDeletedEvent{ID: id}
	body := mq.Event{
		Type:      "product.deleted",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "product.deleted", body); err != nil {
		log.Printf("publish product.deleted: %v", err)
	}
}

// PublishInventoryCreated 发布库存创建
func (p *Publisher) PublishInventoryCreated(inventory *models.Inventory) {
	evt := mq.InventoryCreatedEvent{
		ID:        inventory.ID,
		ProductID: inventory.ProductID,
		Quantity:  inventory.Quantity,
		Status:    inventory.Status,
	}
	body := mq.Event{
		Type:      "inventory.created",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "inventory.created", body); err != nil {
		log.Printf("publish inventory.created: %v", err)
	}
}

// PublishInventoryUpdated 发布库存更新
func (p *Publisher) PublishInventoryUpdated(inventory *models.Inventory) {
	evt := mq.InventoryUpdatedEvent{ID: inventory.ID, ProductID: inventory.ProductID, Quantity: inventory.Quantity, Status: inventory.Status}
	body := mq.Event{
		Type:      "inventory.updated",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "inventory.updated", body); err != nil {
		log.Printf("publish inventory.updated: %v", err)
	}
}

// PublishInventoryDeleted 发布库存删除
func (p *Publisher) PublishInventoryDeleted(id string) {
	evt := mq.InventoryDeletedEvent{ID: id}
	body := mq.Event{
		Type:      "inventory.deleted",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "inventory.deleted", body); err != nil {
		log.Printf("publish inventory.deleted: %v", err)
	}
}

// PublishCategoryCreated 发布分类创建
func (p *Publisher) PublishCategoryCreated(category *models.Category) {
	evt := mq.CategoryCreatedEvent{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
	body := mq.Event{
		Type:      "category.created",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "category.created", body); err != nil {
		log.Printf("publish category.created: %v", err)
	}
}

// PublishCategoryUpdated 发布分类更新
func (p *Publisher) PublishCategoryUpdated(category *models.Category) {
	evt := mq.CategoryUpdatedEvent{ID: category.ID, Name: category.Name}
	body := mq.Event{
		Type:      "category.updated",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "category.updated", body); err != nil {
		log.Printf("publish category.updated: %v", err)
	}
}

// PublishCategoryDeleted 发布分类删除
func (p *Publisher) PublishCategoryDeleted(id string) {
	evt := mq.CategoryDeletedEvent{ID: id}
	body := mq.Event{
		Type:      "category.deleted",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "category.deleted", body); err != nil {
		log.Printf("publish category.deleted: %v", err)
	}
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
