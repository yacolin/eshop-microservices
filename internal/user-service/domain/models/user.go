package models

import (
	"eshop-microservices/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID     string `gorm:"type:varchar(36);primaryKey" json:"id"`
	Status int    `gorm:"type:tinyint;default:1" json:"status"`

	CreatedAt utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	UserInfo *UserInfo `gorm:"foreignKey:UserID" json:"user_info,omitempty"`
	Roles    []Role    `gorm:"many2many:user_roles;" json:"roles,omitempty"`
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

func (u *User) IsActive() bool {
	return u.Status == 1
}

// GetPrimaryIdentity 获取用户的主身份凭证（用于显示用户名等）
func (u *User) GetPrimaryIdentity(identities []UserIdentity) *UserIdentity {
	if len(identities) == 0 {
		return nil
	}
	// 优先返回 password 类型的身份
	for _, identity := range identities {
		if identity.Provider == ProviderPassword.String() {
			return &identity
		}
	}
	// 否则返回第一个
	return &identities[0]
}

// UserInfo 用户详细信息模型（对应 User.md 中的 user_profile）
// 保存可变个人信息：nickname, avatar, gender 等
type UserInfo struct {
	ID       string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID   string     `gorm:"type:varchar(36);not null;uniqueIndex" json:"user_id"`
	Avatar   string     `gorm:"type:varchar(255)" json:"avatar"`
	Nickname string     `gorm:"type:varchar(50)" json:"nickname"`
	Gender   int        `gorm:"type:tinyint;default:0" json:"gender"` // 0:未知 1:男 2:女
	Birthday *time.Time `json:"birthday"`
	Address  string     `gorm:"type:varchar(255)" json:"address"`
	Bio      string     `gorm:"type:text" json:"bio"`
	Country  string     `gorm:"type:varchar(50)" json:"country"`
	Province string     `gorm:"type:varchar(50)" json:"province"`
	City     string     `gorm:"type:varchar(50)" json:"city"`
	ZipCode  string     `gorm:"type:varchar(20)" json:"zip_code"`
	Language string     `gorm:"type:varchar(20);default:zh-CN" json:"language"`
	Timezone string     `gorm:"type:varchar(50);default:Asia/Shanghai" json:"timezone"`

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
