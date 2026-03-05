package models

import (
	"eshop-microservices/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          string `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string `gorm:"type:text" json:"description"`
	Status      int    `gorm:"type:tinyint;default:1" json:"status"`
	Sort        int    `gorm:"type:int;default:0" json:"sort"`
	IsSystem    bool   `gorm:"type:tinyint(1);default:0" json:"is_system"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	Permissions []Permission `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

type UserRole struct {
	ID     string `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID string `gorm:"type:varchar(36);not null;index:idx_user_role" json:"user_id"`
	RoleID string `gorm:"type:varchar(36);not null;index:idx_user_role" json:"role_id"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (ur *UserRole) BeforeCreate(tx *gorm.DB) error {
	if ur.ID == "" {
		ur.ID = uuid.New().String()
	}
	return nil
}
