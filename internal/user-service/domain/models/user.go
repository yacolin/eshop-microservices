package models

import (
	"eshop-microservices/pkg/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 角色常量定义
const (
	RoleAdmin    = "admin"    // 超级管理员
	RoleCustomer = "customer" // 普通用户
	RoleSystem   = "system"   // 系统用户
	RoleMerchant = "merchant" // 商家
	RoleOperator = "operator" // 运营人员
)

// User 用户主表 - 只保留业务核心信息
// 遵循 User.md 设计：User 表保存业务相关的用户信息
type User struct {
	ID     string `gorm:"type:varchar(36);primaryKey" json:"id"`
	Roles  string `gorm:"type:varchar(255);default:'customer'" json:"roles"` // 逗号分隔的角色列表，如: "admin,customer"
	Status int    `gorm:"type:tinyint;default:1" json:"status"`            // 1:正常 2:禁用

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
	// 设置默认值
	if u.Roles == "" {
		u.Roles = RoleCustomer
	}
	return nil
}

// GetRoles 获取角色列表
func (u *User) GetRoles() []string {
	if u.Roles == "" {
		return []string{RoleCustomer}
	}
	return strings.Split(u.Roles, ",")
}

// HasRole 检查是否有指定角色
func (u *User) HasRole(role string) bool {
	roles := u.GetRoles()
	for _, r := range roles {
		if strings.TrimSpace(r) == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查是否有任意一个指定角色
func (u *User) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

// AddRole 添加角色
func (u *User) AddRole(role string) {
	if u.HasRole(role) {
		return
	}
	if u.Roles == "" {
		u.Roles = role
	} else {
		u.Roles = u.Roles + "," + role
	}
}

// RemoveRole 移除角色
func (u *User) RemoveRole(role string) {
	roles := u.GetRoles()
	var newRoles []string
	for _, r := range roles {
		if strings.TrimSpace(r) != role {
			newRoles = append(newRoles, r)
		}
	}
	u.Roles = strings.Join(newRoles, ",")
}

// IsAdmin 检查是否为管理员（支持多角色）
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// IsActive 检查用户是否活跃
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
