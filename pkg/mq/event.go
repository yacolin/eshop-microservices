package mq

import "encoding/json"

// Event 基础事件结构
type Event struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp string          `json:"timestamp"`
	Source    string          `json:"source"`
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	TotalAmount int64  `json:"total_amount"` // 单位：分
	Status      string `json:"status"`
}

// OrderUpdatedEvent 订单更新事件
type OrderUpdatedEvent struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	ID string `json:"id"`
}

// OrderCompletedEvent 订单完成事件
type OrderCompletedEvent struct {
	ID string `json:"id"`
}

// ProductCreatedEvent 产品创建事件
type ProductCreatedEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	SKU         string `json:"sku"`
}

// ProductUpdatedEvent 产品更新事件
type ProductUpdatedEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

// ProductDeletedEvent 产品删除事件
type ProductDeletedEvent struct {
	ID string `json:"id"`
}

// InventoryCreatedEvent 库存创建事件
type InventoryCreatedEvent struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

// InventoryUpdatedEvent 库存更新事件
type InventoryUpdatedEvent struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

// InventoryDeletedEvent 库存删除事件
type InventoryDeletedEvent struct {
	ID string `json:"id"`
}

// CategoryCreatedEvent 分类创建事件
type CategoryCreatedEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CategoryUpdatedEvent 分类更新事件
type CategoryUpdatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CategoryDeletedEvent 分类删除事件
type CategoryDeletedEvent struct {
	ID string `json:"id"`
}
