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
	CreatePermission(req *CreatePermissionRequest) (*models.Permission, error)
	GetPermission(id string) (*models.Permission, error)
	GetPermissionByName(name string) (*models.Permission, error)
	UpdatePermission(id string, req *UpdatePermissionRequest) (*models.Permission, error)
	DeletePermission(id string) error

	ListPermissions(page, pageSize int) (*ListPermissionsResponse, error)
	GetPermissionsByCategory(category string, page, pageSize int) (*ListPermissionsResponse, error)
	GetPermissionsByResource(resource string, page, pageSize int) (*ListPermissionsResponse, error)
	GetPermissionsByRoleID(roleID string, page, pageSize int) (*ListPermissionsResponse, error)

	CheckPermissionsByRoleIDs(roleIDs []string, permissionNames []string) (map[string]bool, error)
	CheckUserPermissions(userID string, permissionNames []string) (map[string]bool, error)

	CreateRole(req *CreateRoleRequest) (*models.Role, error)
	GetRole(id string) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	UpdateRole(id string, req *UpdateRoleRequest) (*models.Role, error)
	DeleteRole(id string) error
	ListRoles(page, pageSize int) (*ListRolesResponse, error)
	AssignRoleToUser(userID, roleID string) error
	RemoveRoleFromUser(userID, roleID string) error
	GetUserRoles(userID string) ([]models.Role, error)
	AssignPermissionsToRole(roleID string, permissionIDs []string) error
	RemovePermissionsFromRole(roleID string, permissionIDs []string) error
}

type permissionService struct {
	permissionRepo repositories.PermissionRepository
	userRepo       repositories.UserRepository
	roleRepo       repositories.RoleRepository
}

func NewPermissionService(
	permissionRepo repositories.PermissionRepository,
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
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

func (s *permissionService) GetPermissionsByRoleID(roleID string, page, pageSize int) (*ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.ByRoleID(roleID, pageSize, offset)
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

func (s *permissionService) CheckPermissionsByRoleIDs(roleIDs []string, permissionNames []string) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, permissionName := range permissionNames {
		has, err := s.permissionRepo.HasPermissionByRoleIDs(roleIDs, permissionName)
		if err != nil {
			return nil, err
		}
		result[permissionName] = has
	}

	return result, nil
}

func (s *permissionService) CheckUserPermissions(userID string, permissionNames []string) (map[string]bool, error) {
	roles, err := s.roleRepo.GetUserRoles(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	return s.CheckPermissionsByRoleIDs(roleIDs, permissionNames)
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
	IsSystem    bool   `json:"is_system"`
}

type UpdateRoleRequest struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
	Sort        *int    `json:"sort"`
}

type ListRolesResponse struct {
	Roles  []models.Role `json:"roles"`
	Total  int64         `json:"total"`
	Page   int           `json:"page"`
	PageSize int         `json:"page_size"`
}

func (s *permissionService) CreateRole(req *CreateRoleRequest) (*models.Role, error) {
	role := &models.Role{
		ID:          uuid.New().String(),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Status:      req.Status,
		Sort:        req.Sort,
		IsSystem:    req.IsSystem,
	}

	if role.Status == 0 {
		role.Status = 1
	}

	if err := s.roleRepo.Create(context.Background(), role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *permissionService) GetRole(id string) (*models.Role, error) {
	return s.roleRepo.GetByID(context.Background(), id)
}

func (s *permissionService) GetRoleByName(name string) (*models.Role, error) {
	return s.roleRepo.GetByName(context.Background(), name)
}

func (s *permissionService) UpdateRole(id string, req *UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if req.DisplayName != nil {
		role.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if req.Sort != nil {
		role.Sort = *req.Sort
	}

	if err := s.roleRepo.Update(context.Background(), role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *permissionService) DeleteRole(id string) error {
	return s.roleRepo.Delete(context.Background(), id)
}

func (s *permissionService) ListRoles(page, pageSize int) (*ListRolesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	roles, total, err := s.roleRepo.List(context.Background(), pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListRolesResponse{
		Roles:    roles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *permissionService) AssignRoleToUser(userID, roleID string) error {
	return s.roleRepo.AssignRoleToUser(context.Background(), userID, roleID)
}

func (s *permissionService) RemoveRoleFromUser(userID, roleID string) error {
	return s.roleRepo.RemoveRoleFromUser(context.Background(), userID, roleID)
}

func (s *permissionService) GetUserRoles(userID string) ([]models.Role, error) {
	return s.roleRepo.GetUserRoles(context.Background(), userID)
}

func (s *permissionService) AssignPermissionsToRole(roleID string, permissionIDs []string) error {
	for _, permissionID := range permissionIDs {
		if err := s.roleRepo.AssignPermissionToRole(context.Background(), roleID, permissionID); err != nil {
			return err
		}
	}
	return nil
}

func (s *permissionService) RemovePermissionsFromRole(roleID string, permissionIDs []string) error {
	for _, permissionID := range permissionIDs {
		if err := s.roleRepo.RemovePermissionFromRole(context.Background(), roleID, permissionID); err != nil {
			return err
		}
	}
	return nil
}
