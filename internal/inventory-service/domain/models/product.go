package models

import (
	"eshop-microservices/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product 产品
type Product struct {
	ID          string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Price       int64      `gorm:"type:bigint;not null" json:"price"` // 价格，单位：分
	SKU         string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	CategoryID  *string    `gorm:"type:varchar(36);index" json:"category_id"`       // 分类ID
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category"`           // 所属分类
	Categories  []Category `gorm:"many2many:product_categories;" json:"categories"` // 多个分类

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
