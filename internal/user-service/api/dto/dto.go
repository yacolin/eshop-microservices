package dto

import (
	pkgQuery "eshop-microservices/pkg/query"
)

// ========== 认证相关 DTO ==========

// LoginRequest 登录请求（通用）
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
	IsNewUser    bool   `json:"is_new_user"` // 是否新用户
}

// PasswordLoginRequest 密码登录请求
type PasswordLoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code   string `json:"code" binding:"required" example:"wx_code_xxx"`
	AppID  string `json:"appid" binding:"required" example:"wx_appid_xxx"`
	Source string `json:"source,omitempty" example:"miniapp"`
}

// PhoneLoginRequest 手机号登录请求
type PhoneLoginRequest struct {
	Phone      string `json:"phone" binding:"required" example:"13800138000"`
	VerifyCode string `json:"verify_code" binding:"required" example:"123456"`
}

// EmailLoginRequest 邮箱登录请求
type EmailLoginRequest struct {
	Email      string `json:"email" binding:"required,email" example:"john@example.com"`
	VerifyCode string `json:"verify_code" binding:"required" example:"123456"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username,omitempty" example:"john_doe"`
	Email    string `json:"email,omitempty" example:"john@example.com"`
	Phone    string `json:"phone,omitempty" example:"13800138000"`
	Password string `json:"password,omitempty" example:"password123"`
	Provider string `json:"provider" binding:"required" example:"password"` // 注册方式：password, phone, email
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// BindIdentityRequest 绑定身份请求
type BindIdentityRequest struct {
	Provider   string `json:"provider" binding:"required"`
	Identifier string `json:"identifier" binding:"required"`
	Credential string `json:"credential,omitempty"`
}

// ========== 用户资料相关 DTO ==========

// UpdateUserInfoRequest 更新用户详细信息请求
// 对应 User.md 中的 user_profile，包含可变个人信息
type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
	Address  string `json:"address" binding:"max=255"`
	Bio      string `json:"bio" binding:"max=500"`
	Country  string `json:"country" binding:"max=50"`
	Province string `json:"province" binding:"max=50"`
	City     string `json:"city" binding:"max=50"`
	ZipCode  string `json:"zip_code" binding:"max=20"`
	Language string `json:"language" binding:"max=20"`
	Timezone string `json:"timezone" binding:"max=50"`
}

// UserListQuery 用户列表查询
type UserListQuery struct {
	pkgQuery.Pagination
}

// ========== 权限相关 DTO ==========

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required" example:"order:create"`
	DisplayName string `json:"display_name" binding:"required" example:"创建订单"`
	Description string `json:"description" example:"创建新订单"`
	Resource    string `json:"resource" binding:"required" example:"order"`
	Action      string `json:"action" binding:"required" example:"create"`
	Category    string `json:"category" example:"business"`
	Sort        int    `json:"sort" example:"0"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	DisplayName *string `json:"display_name" example:"更新显示名称"`
	Description *string `json:"description" example:"更新描述"`
	Category    *string `json:"category" example:"admin"`
	Sort        *int    `json:"sort" example:"10"`
	Status      *int    `json:"status" example:"1"`
}

// PermissionListQuery 权限列表查询
type PermissionListQuery struct {
	pkgQuery.Pagination
	Category string `form:"category"`
	Resource string `form:"resource"`
	Role     string `form:"role"`
}

// AssignPermissionToRoleRequest 分配权限给角色请求
type AssignPermissionToRoleRequest struct {
	PermissionID string `json:"permission_id" binding:"required" example:"permission-id"`
}

// BatchAssignPermissionsToRoleRequest 批量分配权限给角色请求
type BatchAssignPermissionsToRoleRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

// CheckPermissionsRequest 检查权限请求
type CheckPermissionsRequest struct {
	PermissionNames []string `json:"permission_names" binding:"required" example:"order:create,product:read"`
}

// CheckPermissionsResponse 检查权限响应
type CheckPermissionsResponse struct {
	Permissions map[string]bool `json:"permissions"` // key: permission_name, value: 是否有权限
}

// UpdateUserRoleRequest 更新用户角色请求
type UpdateUserRoleRequest struct {
	Roles []string `json:"roles" binding:"required" example:"admin,customer"`
}
