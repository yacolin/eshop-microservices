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