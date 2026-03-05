package service

import (
	"context"
	"time"

	"eshop-microservices/internal/user-service/api/dto"
	"eshop-microservices/internal/user-service/domain/models"
	"eshop-microservices/internal/user-service/domain/repositories"
	"eshop-microservices/pkg/errcode"

	"gorm.io/gorm"
)

type UserService struct {
	repo      repositories.UserRepository
	jwtSecret string
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SetJWTSecret(secret string) {
	s.jwtSecret = secret
}

// GetProfile 获取用户资料（包含 UserInfo）
func (s *UserService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.GetByIDWithInfo(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUserInfo 获取用户详细信息
func (s *UserService) GetUserInfo(ctx context.Context, userID string) (*models.UserInfo, error) {
	userInfo, err := s.repo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return userInfo, nil
}

// UpdateUserInfo 更新用户详细信息（包含 Avatar、Nickname 等）
func (s *UserService) UpdateUserInfo(ctx context.Context, userID string, req dto.UpdateUserInfoRequest) (*models.UserInfo, error) {
	userInfo, err := s.repo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userInfo = &models.UserInfo{UserID: userID}
		} else {
			return nil, err
		}
	}

	if req.Nickname != "" {
		userInfo.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		userInfo.Avatar = req.Avatar
	}
	if req.Gender != 0 {
		userInfo.Gender = req.Gender
	}
	if req.Birthday != "" {
		// Parse birthday string to time.Time
		// Assuming format: 2006-01-02
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			userInfo.Birthday = &birthday
		}
	}
	if req.Address != "" {
		userInfo.Address = req.Address
	}
	if req.Bio != "" {
		userInfo.Bio = req.Bio
	}
	if req.Country != "" {
		userInfo.Country = req.Country
	}
	if req.Province != "" {
		userInfo.Province = req.Province
	}
	if req.City != "" {
		userInfo.City = req.City
	}
	if req.ZipCode != "" {
		userInfo.ZipCode = req.ZipCode
	}
	if req.Language != "" {
		userInfo.Language = req.Language
	}
	if req.Timezone != "" {
		userInfo.Timezone = req.Timezone
	}

	if userInfo.ID == "" {
		if err := s.repo.CreateUserInfo(ctx, userInfo); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.UpdateUserInfo(ctx, userInfo); err != nil {
			return nil, err
		}
	}

	return userInfo, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateUserRoles 更新用户角色
func (s *UserService) UpdateUserRoles(userID string, roles []string) error {
	return s.repo.UpdateUserRoles(context.Background(), userID, roles)
}
