package service

import (
	"context"
	"errors"

	"eshop-microservices/internal/user-service/api/dto"
	"eshop-microservices/internal/user-service/domain/models"
	"eshop-microservices/internal/user-service/domain/repositories"
	"eshop-microservices/pkg/errcode"
	"eshop-microservices/pkg/utils"

	"golang.org/x/crypto/bcrypt"
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

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*models.User, error) {
	existingUser, err := s.repo.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errcode.ErrUserAlreadyRegistered
	}

	existingUser, err = s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errcode.ErrEmailAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Phone:        req.Phone,
		UserInfo: &models.UserInfo{
			Nickname: req.Username,
		},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (map[string]string, error) {
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrUserNotFound
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errcode.ErrInvalidCredentials
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return tokens, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.GetByIDWithInfo(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrUserNotFound
		}
		return nil, err
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, userID string) (*models.UserInfo, error) {
	userInfo, err := s.repo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return userInfo, nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, userID string, req dto.UpdateUserInfoRequest) (*models.UserInfo, error) {
	userInfo, err := s.repo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
		// userInfo.Birthday = &birthday
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

func (s *UserService) Logout(ctx context.Context, userID string) error {
	return nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}
