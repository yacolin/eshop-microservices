package models

import (
	"eshop-microservices/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"type:varchar(36);primaryKey" json:"id"`
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string `gorm:"type:varchar(100)" json:"full_name"`
	Phone        string `gorm:"type:varchar(20)" json:"phone"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	UserInfo *UserInfo `gorm:"foreignKey:UserID" json:"user_info,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type UserInfo struct {
	ID       string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID   string     `gorm:"type:varchar(36);not null;uniqueIndex" json:"user_id"`
	Avatar   string     `gorm:"type:varchar(255)" json:"avatar"`
	Gender   int        `gorm:"type:tinyint;default:0" json:"gender"` // 0:未知 1:男 2:女
	Birthday *time.Time `json:"birthday"`
	Address  string     `gorm:"type:varchar(255)" json:"address"`
	Bio      string     `gorm:"type:text" json:"bio"`
	Nickname string     `gorm:"type:varchar(50)" json:"nickname"`
	Country  string     `gorm:"type:varchar(50)" json:"country"`
	Province string     `gorm:"type:varchar(50)" json:"province"`
	City     string     `gorm:"type:varchar(50)" json:"city"`
	ZipCode  string     `gorm:"type:varchar(20)" json:"zip_code"`
	Language string     `gorm:"type:varchar(20);default:zh-CN" json:"language"`
	Timezone string     `gorm:"type:varchar(50);default:Asia/Shanghai" json:"timezone"`
	Status   int        `gorm:"type:tinyint;default:1" json:"status"` // 1:正常 2:禁用

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (UserInfo) TableName() string {
	return "user_infos"
}

func (u *UserInfo) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
