package models

import (
	"eshop-microservices/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category 分类
type Category struct {
	ID          string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	ParentID    *string    `gorm:"type:varchar(36);index" json:"parent_id"` // 父分类ID，支持层级结构
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children"`
	Products    []Product  `gorm:"many2many:product_categories;" json:"products"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 分类表名
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate GORM 钩子：生成 UUID
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// InventoryStatus 枚举类型
type InventoryStatus string

// 定义枚举值
const (
	InventoryStatusInStock    = "instock"
	InventoryStatusOutOfStock = "outofstock"
	InventoryStatusLowStock   = "lowstock"
)

// Product 产品
type Product struct {
	ID          string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Price       int64      `gorm:"type:bigint;not null" json:"price"` // 价格，单位：分
	SKU         string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	CategoryID  *string    `gorm:"type:varchar(36);index" json:"category_id"` // 分类ID
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"-"`            // 所属分类
	Categories  []Category `gorm:"many2many:product_categories;" json:"-"`    // 多个分类

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 产品表名
func (Product) TableName() string {
	return "products"
}

// BeforeCreate GORM 钩子：生成 UUID
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

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
