package models

import (
	"eshop-microservices/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InventoryStatus 枚举类型
type InventoryStatus string

// 定义枚举值
const (
	InventoryStatusInStock    = "instock"
	InventoryStatusOutOfStock = "outofstock"
	InventoryStatusLowStock   = "lowstock"
)

// Inventory 库存
type Inventory struct {
	ID        string  `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProductID string  `gorm:"type:varchar(36);not null;index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null;default:0" json:"quantity"`
	Status    string  `gorm:"type:varchar(20);not null;default:'instock'" json:"status"`
	Reserved  int     `gorm:"not null;default:0" json:"reserved"`   // 已预订数量
	Threshold int     `gorm:"not null;default:10" json:"threshold"` // 低库存阈值

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 库存表名
func (Inventory) TableName() string {
	return "inventories"
}

// BeforeCreate GORM 钩子：生成 UUID
func (i *Inventory) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	// 设置初始状态
	if i.Status == "" {
		if i.Quantity <= i.Threshold {
			if i.Quantity <= 0 {
				i.Status = InventoryStatusOutOfStock
			} else {
				i.Status = InventoryStatusLowStock
			}
		} else {
			i.Status = InventoryStatusInStock
		}
	}
	return nil
}

// UpdateStatus 更新库存状态
func (i *Inventory) UpdateStatus() {
	if i.Quantity <= 0 {
		i.Status = InventoryStatusOutOfStock
	} else if i.Quantity <= i.Threshold {
		i.Status = InventoryStatusLowStock
	} else {
		i.Status = InventoryStatusInStock
	}
}