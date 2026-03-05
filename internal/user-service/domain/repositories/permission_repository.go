package repositories

import (
	"eshop-microservices/internal/user-service/domain/models"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	// 基础 CRUD
	Create(permission *models.Permission) error
	GetByID(id string) (*models.Permission, error)
	GetByName(name string) (*models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id string) error

	// 查询
	List(limit, offset int) ([]*models.Permission, int64, error)
	ByCategory(category string, limit, offset int) ([]*models.Permission, int64, error)
	ByResource(resource string, limit, offset int) ([]*models.Permission, int64, error)
	ByRole(roleName string, limit, offset int) ([]*models.Permission, int64, error)
	ByStatus(status int, limit, offset int) ([]*models.Permission, int64, error)

	// 检查
	ExistsByName(name string) (bool, error)
	GetPermissionsByRoles(roleNames []string) ([]*models.Permission, error)
	HasPermission(roleNames []string, permissionName string) (bool, error)

	// 角色权限关联
	AssignPermissionToRole(roleName, permissionID string) error
	RemovePermissionFromRole(roleName, permissionID string) error
	GetRolePermissions(roleName string, limit, offset int) ([]*models.RolePermission, int64, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) GetByID(id string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id string) error {
	return r.db.Delete(&models.Permission{}, "id = ?", id).Error
}

func (r *permissionRepository) List(limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ByCategory(category string, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).Where("category = ?", category)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ByResource(resource string, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).Where("resource = ?", resource)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ByRole(roleName string, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_name = ?", roleName).
		Where("role_permissions.deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("permissions.sort ASC, permissions.created_at DESC").
		Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ByStatus(status int, limit, offset int) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) GetPermissionsByRoles(roleNames []string) ([]*models.Permission, error) {
	var permissions []*models.Permission

	err := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_name IN ?", roleNames).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.status = ?", 1).
		Distinct("permissions.*").
		Find(&permissions).Error

	return permissions, err
}

func (r *permissionRepository) HasPermission(roleNames []string, permissionName string) (bool, error) {
	var count int64

	err := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_name IN ?", roleNames).
		Where("role_permissions.deleted_at IS NULL").
		Where("permissions.name = ?", permissionName).
		Where("permissions.status = ?", 1).
		Count(&count).Error

	return count > 0, err
}

func (r *permissionRepository) AssignPermissionToRole(roleName, permissionID string) error {
	rolePermission := &models.RolePermission{
		RoleName:     roleName,
		PermissionID: permissionID,
	}

	return r.db.Create(rolePermission).Error
}

func (r *permissionRepository) RemovePermissionFromRole(roleName, permissionID string) error {
	return r.db.Where("role_name = ? AND permission_id = ?", roleName, permissionID).
		Delete(&models.RolePermission{}).Error
}

func (r *permissionRepository) GetRolePermissions(roleName string, limit, offset int) ([]*models.RolePermission, int64, error) {
	var rolePermissions []*models.RolePermission
	var total int64

	query := r.db.Model(&models.RolePermission{}).Where("role_name = ?", roleName)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Permission").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&rolePermissions).Error

	return rolePermissions, total, err
}
