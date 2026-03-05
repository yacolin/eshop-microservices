package handlers

import (
	"strconv"

	"eshop-microservices/internal/user-service/service"
	"eshop-microservices/pkg/errcode"
	"eshop-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	permissionSvc service.PermissionService
}

func NewRoleHandler(permissionSvc service.PermissionService) *RoleHandler {
	return &RoleHandler{
		permissionSvc: permissionSvc,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色（需要管理员权限）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param request body service.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	role, err := h.permissionSvc.CreateRole(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, role)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	role, err := h.permissionSvc.GetRole(id)
	if err != nil {
		c.Error(errcode.ErrNotFound)
		return
	}

	response.Success(c, role)
}

// GetRoleByName 根据名称获取角色
// @Summary 根据名称获取角色
// @Description 根据角色名称获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param name path string true "角色名称"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles/name/{name} [get]
func (h *RoleHandler) GetRoleByName(c *gin.Context) {
	name := c.Param("name")
	role, err := h.permissionSvc.GetRoleByName(name)
	if err != nil {
		c.Error(errcode.ErrNotFound)
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息（需要管理员权限）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param request body service.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	role, err := h.permissionSvc.UpdateRole(id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色（需要管理员权限）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.permissionSvc.DeleteRole(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Role deleted successfully"})
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 分页获取角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=service.ListRolesResponse}
// @Router /api/v1/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.permissionSvc.ListRoles(page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// AssignRoleToUser 分配角色给用户
// @Summary 分配角色给用户
// @Description 为指定用户分配角色（需要管理员权限）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param request body object true "角色信息"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{user_id}/roles [post]
func (h *RoleHandler) AssignRoleToUser(c *gin.Context) {
	userID := c.Param("user_id")
	var req struct {
		RoleID string `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionSvc.AssignRoleToUser(userID, req.RoleID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Role assigned to user successfully"})
}

// RemoveRoleFromUser 移除用户的角色
// @Summary 移除用户的角色
// @Description 移除指定用户的角色（需要管理员权限）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param role_id path string true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{user_id}/roles/{role_id} [delete]
func (h *RoleHandler) RemoveRoleFromUser(c *gin.Context) {
	userID := c.Param("user_id")
	roleID := c.Param("role_id")

	if err := h.permissionSvc.RemoveRoleFromUser(userID, roleID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Role removed from user successfully"})
}

// GetUserRoles 获取用户的角色列表
// @Summary 获取用户的角色列表
// @Description 获取指定用户的角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} response.Response{data=[]models.Role}
// @Router /api/v1/users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("user_id")

	roles, err := h.permissionSvc.GetUserRoles(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, roles)
}

// AssignPermissionsToRole 为角色分配权限
// @Summary 为角色分配权限
// @Description 为指定角色分配权限（需要管理员权限）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param request body object true "权限ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID := c.Param("id")
	var req struct {
		PermissionIDs []string `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionSvc.AssignPermissionsToRole(roleID, req.PermissionIDs); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permissions assigned to role successfully"})
}

// RemovePermissionsFromRole 移除角色的权限
// @Summary 移除角色的权限
// @Description 移除指定角色的权限（需要管理员权限）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param request body object true "权限ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{id}/permissions [delete]
func (h *RoleHandler) RemovePermissionsFromRole(c *gin.Context) {
	roleID := c.Param("id")
	var req struct {
		PermissionIDs []string `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionSvc.RemovePermissionsFromRole(roleID, req.PermissionIDs); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permissions removed from role successfully"})
}
