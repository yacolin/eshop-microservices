package handlers

import (
	"eshop-microservices/internal/user-service/api/dto"
	"eshop-microservices/internal/user-service/domain/models"
	"eshop-microservices/internal/user-service/service"
	"eshop-microservices/pkg/errcode"
	"eshop-microservices/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService service.PermissionService
	userService       *service.UserService
}

func NewPermissionHandler(
	permissionService service.PermissionService,
	userService *service.UserService,
) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
		userService:       userService,
	}
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Description 创建新权限（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body dto.CreatePermissionRequest true "权限信息"
// @Success 200 {object} response.Response{data=models.Permission}
// @Router /api/v1/permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	permission, err := h.permissionService.CreatePermission(&service.CreatePermissionRequest{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
		Sort:        req.Sort,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permission)
}

// GetPermission 获取权限详情
// @Summary 获取权限详情
// @Description 根据ID获取权限详情
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Success 200 {object} response.Response{data=models.Permission}
// @Router /api/v1/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	permission, err := h.permissionService.GetPermission(id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Description 更新权限信息（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Param request body dto.UpdatePermissionRequest true "权限信息"
// @Success 200 {object} response.Response{data=models.Permission}
// @Router /api/v1/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	permission, err := h.permissionService.UpdatePermission(id, &service.UpdatePermissionRequest{
		DisplayName: req.DisplayName,
		Description:  req.Description,
		Category:     req.Category,
		Sort:         req.Sort,
		Status:       req.Status,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permission)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Description 删除权限（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Success 200 {object} response.Response
// @Router /api/v1/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	if err := h.permissionService.DeletePermission(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permission deleted successfully"})
}

// ListPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 获取权限列表，支持分页和筛选
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Param category query string false "分类"
// @Param resource query string false "资源"
// @Param role query string false "角色"
// @Success 200 {object} response.Response{data=service.ListPermissionsResponse}
// @Router /api/v1/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	var query dto.PermissionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}
	query.Normalize()

	var result *service.ListPermissionsResponse
	var err error

	switch {
	case query.Category != "":
		result, err = h.permissionService.GetPermissionsByCategory(query.Category, query.Page, query.Size)
	case query.Resource != "":
		result, err = h.permissionService.GetPermissionsByResource(query.Resource, query.Page, query.Size)
	case query.Role != "":
		result, err = h.permissionService.GetPermissionsByRole(query.Role, query.Page, query.Size)
	default:
		result, err = h.permissionService.ListPermissions(query.Page, query.Size)
	}

	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// AssignPermissionToRole 分配权限给角色
// @Summary 分配权限给角色
// @Description 为指定角色分配权限（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param role path string true "角色名称"
// @Param request body dto.AssignPermissionToRoleRequest true "权限ID"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{role}/permissions [post]
func (h *PermissionHandler) AssignPermissionToRole(c *gin.Context) {
	roleName := c.Param("role")
	if roleName == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.AssignPermissionToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionService.AssignPermissionToRole(roleName, req.PermissionID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permission assigned to role successfully"})
}

// RemovePermissionFromRole 移除角色的权限
// @Summary 移除角色的权限
// @Description 移除指定角色的权限（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param role path string true "角色名称"
// @Param permission_id path string true "权限ID"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{role}/permissions/{permission_id} [delete]
func (h *PermissionHandler) RemovePermissionFromRole(c *gin.Context) {
	roleName := c.Param("role")
	permissionID := c.Param("permission_id")

	if roleName == "" || permissionID == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	if err := h.permissionService.RemovePermissionFromRole(roleName, permissionID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permission removed from role successfully"})
}

// GetRolePermissions 获取角色的权限列表
// @Summary 获取角色的权限列表
// @Description 获取指定角色的所有权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param role path string true "角色名称"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=service.ListRolePermissionsResponse}
// @Router /api/v1/roles/{role}/permissions [get]
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	roleName := c.Param("role")
	if roleName == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.permissionService.GetRolePermissions(roleName, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// BatchAssignPermissionsToRole 批量分配权限给角色
// @Summary 批量分配权限给角色
// @Description 为指定角色批量分配权限（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param role path string true "角色名称"
// @Param request body dto.BatchAssignPermissionsToRoleRequest true "权限ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{role}/permissions/batch [post]
func (h *PermissionHandler) BatchAssignPermissionsToRole(c *gin.Context) {
	roleName := c.Param("role")
	if roleName == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.BatchAssignPermissionsToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionService.BatchAssignPermissionsToRole(roleName, req.PermissionIDs); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permissions assigned to role successfully"})
}

// CheckPermissions 检查当前用户的权限
// @Summary 检查当前用户的权限
// @Description 检查当前登录用户是否有指定权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body dto.CheckPermissionsRequest true "权限名称列表"
// @Success 200 {object} response.Response{data=dto.CheckPermissionsResponse}
// @Router /api/v1/permissions/check [post]
func (h *PermissionHandler) CheckPermissions(c *gin.Context) {
	var req dto.CheckPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	result, err := h.permissionService.CheckUserPermissions(userIDStr, req.PermissionNames)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dto.CheckPermissionsResponse{
		Permissions: result,
	})
}

// UpdateUserRole 更新用户角色
// @Summary 更新用户角色
// @Description 更新指定用户的角色（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param request body dto.UpdateUserRoleRequest true "角色列表"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{user_id}/roles [put]
func (h *PermissionHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	// 验证角色是否合法
	validRoles := map[string]bool{
		models.RoleAdmin:    true,
		models.RoleCustomer: true,
		models.RoleSystem:   true,
		models.RoleMerchant: true,
		models.RoleOperator: true,
	}

	for _, role := range req.Roles {
		if !validRoles[role] {
			c.Error(errcode.ErrInvalidRoleName)
			return
		}
	}

	if err := h.userService.UpdateUserRoles(userID, req.Roles); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "User roles updated successfully"})
}
