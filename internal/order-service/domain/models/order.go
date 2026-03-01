package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderStatus 枚举类型
type OrderStatus string

// 定义枚举值
const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusShipped   = "shipped"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

// 金额统一以「分」为单位存储，避免浮点精度问题（1 元 = 100 分）

// Order 订单
type Order struct {
	ID          string `gorm:"type:varchar(36);primaryKey" json:"id"`
	CustomerID  string `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	TotalAmount int64  `gorm:"type:bigint;not null" json:"total_amount"` // 订单总金额，单位：分
	Currency    string `gorm:"type:varchar(10);default:CNY" json:"currency"`
	Status      string `gorm:"type:varchar(20);not null;index" json:"status"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// TableName 表名
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate GORM 钩子：生成 UUID
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// OrderItem 订单项
type OrderItem struct {
	ID        string `gorm:"type:varchar(36);primaryKey" json:"id"`
	OrderID   string `gorm:"type:varchar(36);not null;index" json:"order_id"`
	ProductID string `gorm:"type:varchar(36);not null" json:"product_id"`
	Quantity  int    `gorm:"not null" json:"quantity"`
	UnitPrice int64  `gorm:"type:bigint;not null" json:"unit_price"` // 单价，单位：分
	Amount    int64  `gorm:"type:bigint;not null" json:"amount"`     // 单项小计，单位：分 = UnitPrice * Quantity
}

// TableName 表名
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate GORM 钩子
func (i *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}
