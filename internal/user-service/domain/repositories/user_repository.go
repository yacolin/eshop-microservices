package repositories

import (
	"context"

	"eshop-microservices/internal/user-service/domain/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByIDWithInfo(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]models.User, int64, error)
	UpdateUserRoles(ctx context.Context, userID string, roles []string) error

	CreateUserInfo(ctx context.Context, userInfo *models.UserInfo) error
	GetUserInfoByUserID(ctx context.Context, userID string) (*models.UserInfo, error)
	UpdateUserInfo(ctx context.Context, userInfo *models.UserInfo) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if user.UserInfo != nil {
			user.UserInfo.UserID = user.ID
			if err := tx.Create(user.UserInfo).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByIDWithInfo(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("UserInfo").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&models.UserInfo{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.User{}, "id = ?", id).Error
	})
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	var list []models.User
	query := r.db.WithContext(ctx).Model(&models.User{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *userRepository) CreateUserInfo(ctx context.Context, userInfo *models.UserInfo) error {
	return r.db.WithContext(ctx).Create(userInfo).Error
}

func (r *userRepository) GetUserInfoByUserID(ctx context.Context, userID string) (*models.UserInfo, error) {
	var userInfo models.UserInfo
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (r *userRepository) UpdateUserInfo(ctx context.Context, userInfo *models.UserInfo) error {
	return r.db.WithContext(ctx).Save(userInfo).Error
}

func (r *userRepository) UpdateUserRoles(ctx context.Context, userID string, roles []string) error {
	// 使用 User 模型中的方法构建 roles 字符串
	tempUser := &models.User{ID: userID}
	for _, role := range roles {
		tempUser.AddRole(role)
	}
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("roles", tempUser.Roles).Error
}
