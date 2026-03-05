package service

import (
	"context"
	"eshop-microservices/internal/user-service/domain/models"
	"eshop-microservices/internal/user-service/domain/repositories"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionService interface {
	// 基础 CRUD
	CreatePermission(req *CreatePermissionRequest) (*models.Permission, error)
	GetPermission(id string) (*models.Permission, error)
	GetPermissionByName(name string) (*models.Permission, error)
	UpdatePermission(id string, req *UpdatePermissionRequest) (*models.Permission, error)
	DeletePermission(id string) error

	// 查询
	ListPermissions(page, pageSize int) (*ListPermissionsResponse, error)
	GetPermissionsByCategory(category string, page, pageSize int) (*ListPermissionsResponse, error)
	GetPermissionsByResource(resource string, page, pageSize int) (*ListPermissionsResponse, error)
	GetPermissionsByRole(roleName string, page, pageSize int) (*ListPermissionsResponse, error)

	// 角色权限管理
	AssignPermissionToRole(roleName, permissionID string) error
	RemovePermissionFromRole(roleName, permissionID string) error
	GetRolePermissions(roleName string, page, pageSize int) (*ListRolePermissionsResponse, error)
	BatchAssignPermissionsToRole(roleName string, permissionIDs []string) error
	BatchRemovePermissionsFromRole(roleName string, permissionIDs []string) error

	// 权限检查
	CheckPermissions(roleNames []string, permissionNames []string) (map[string]bool, error)
	CheckUserPermissions(userID string, permissionNames []string) (map[string]bool, error)
}

type permissionService struct {
	permissionRepo repositories.PermissionRepository
	userRepo       repositories.UserRepository
}

func NewPermissionService(
	permissionRepo repositories.PermissionRepository,
	userRepo repositories.UserRepository,
) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Category    string `json:"category"`
	Sort        int    `json:"sort"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Sort        *int    `json:"sort"`
	Status      *int    `json:"status"`
}

// ListPermissionsResponse 权限列表响应
type ListPermissionsResponse struct {
	Permissions []*models.Permission `json:"permissions"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
}

// ListRolePermissionsResponse 角色权限列表响应
type ListRolePermissionsResponse struct {
	RolePermissions []*models.RolePermission `json:"role_permissions"`
	Total           int64                    `json:"total"`
	Page            int                      `json:"page"`
	PageSize        int                      `json:"page_size"`
}

func (s *permissionService) CreatePermission(req *CreatePermissionRequest) (*models.Permission, error) {
	// 检查权限名称是否已存在
	exists, err := s.permissionRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("permission name already exists")
	}

	permission := &models.Permission{
		ID:          uuid.New().String(),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
		Sort:        req.Sort,
		Status:      1,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) GetPermission(id string) (*models.Permission, error) {
	return s.permissionRepo.GetByID(id)
}

func (s *permissionService) GetPermissionByName(name string) (*models.Permission, error) {
	return s.permissionRepo.GetByName(name)
}

func (s *permissionService) UpdatePermission(id string, req *UpdatePermissionRequest) (*models.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	if req.DisplayName != nil {
		permission.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}
	if req.Category != nil {
		permission.Category = *req.Category
	}
	if req.Sort != nil {
		permission.Sort = *req.Sort
	}
	if req.Status != nil {
		permission.Status = *req.Status
	}

	if err := s.permissionRepo.Update(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) DeletePermission(id string) error {
	// TODO: 检查是否有角色使用此权限
	return s.permissionRepo.Delete(id)
}

func (s *permissionService) ListPermissions(page, pageSize int) (*ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.List(pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (s *permissionService) GetPermissionsByCategory(category string, page, pageSize int) (*ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.ByCategory(category, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (s *permissionService) GetPermissionsByResource(resource string, page, pageSize int) (*ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.ByResource(resource, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (s *permissionService) GetPermissionsByRole(roleName string, page, pageSize int) (*ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.ByRole(roleName, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (s *permissionService) AssignPermissionToRole(roleName, permissionID string) error {
	// 验证角色是否存在
	validRoles := []string{models.RoleAdmin, models.RoleCustomer, models.RoleSystem, models.RoleMerchant, models.RoleOperator}
	isValid := false
	for _, r := range validRoles {
		if r == roleName {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid role name")
	}

	// 验证权限是否存在
	_, err := s.permissionRepo.GetByID(permissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return err
	}

	return s.permissionRepo.AssignPermissionToRole(roleName, permissionID)
}

func (s *permissionService) RemovePermissionFromRole(roleName, permissionID string) error {
	return s.permissionRepo.RemovePermissionFromRole(roleName, permissionID)
}

func (s *permissionService) GetRolePermissions(roleName string, page, pageSize int) (*ListRolePermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	rolePermissions, total, err := s.permissionRepo.GetRolePermissions(roleName, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListRolePermissionsResponse{
		RolePermissions: rolePermissions,
		Total:           total,
		Page:            page,
		PageSize:        pageSize,
	}, nil
}

func (s *permissionService) BatchAssignPermissionsToRole(roleName string, permissionIDs []string) error {
	for _, permissionID := range permissionIDs {
		if err := s.AssignPermissionToRole(roleName, permissionID); err != nil {
			return err
		}
	}
	return nil
}

func (s *permissionService) BatchRemovePermissionsFromRole(roleName string, permissionIDs []string) error {
	for _, permissionID := range permissionIDs {
		if err := s.RemovePermissionFromRole(roleName, permissionID); err != nil {
			return err
		}
	}
	return nil
}

func (s *permissionService) CheckPermissions(roleNames []string, permissionNames []string) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, permissionName := range permissionNames {
		has, err := s.permissionRepo.HasPermission(roleNames, permissionName)
		if err != nil {
			return nil, err
		}
		result[permissionName] = has
	}

	return result, nil
}

func (s *permissionService) CheckUserPermissions(userID string, permissionNames []string) (map[string]bool, error) {
	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	roleNames := user.GetRoles()
	return s.CheckPermissions(roleNames, permissionNames)
}
