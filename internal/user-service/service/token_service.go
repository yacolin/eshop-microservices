package service

import (
	"context"
	"encoding/json"
	"time"

	"eshop-microservices/internal/user-service/domain/models"
	"eshop-microservices/internal/user-service/domain/repositories"
	"eshop-microservices/pkg/errcode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenPair Token对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// TokenClaims JWT Claims
type TokenClaims struct {
	UserID     string `json:"user_id"`
	IdentityID string `json:"identity_id"`
	Provider   string `json:"provider"`
	JTI        string `json:"jti"`
	jwt.RegisteredClaims
}

// TokenService Token服务
type TokenService struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	tokenRepo     repositories.AuthTokenRepository
}

// TokenServiceOption Token服务配置选项
type TokenServiceOption func(*TokenService)

// WithAccessExpiry 设置访问令牌过期时间
func WithAccessExpiry(expiry time.Duration) TokenServiceOption {
	return func(s *TokenService) {
		s.accessExpiry = expiry
	}
}

// WithRefreshExpiry 设置刷新令牌过期时间
func WithRefreshExpiry(expiry time.Duration) TokenServiceOption {
	return func(s *TokenService) {
		s.refreshExpiry = expiry
	}
}

// NewTokenService 创建Token服务实例
func NewTokenService(secret string, tokenRepo repositories.AuthTokenRepository, opts ...TokenServiceOption) *TokenService {
	svc := &TokenService{
		secret:        []byte(secret),
		accessExpiry:  2 * time.Hour,      // 默认2小时
		refreshExpiry: 7 * 24 * time.Hour, // 默认7天
		tokenRepo:     tokenRepo,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// GenerateTokenPair 生成Token对
func (s *TokenService) GenerateTokenPair(ctx context.Context, userID, identityID, provider string, meta map[string]interface{}) (*TokenPair, error) {
	now := time.Now()

	// 生成JTI (JWT ID)
	accessJTI := uuid.New().String()
	refreshJTI := uuid.New().String()

	// 生成Access Token
	accessExpiresAt := now.Add(s.accessExpiry)
	accessClaims := TokenClaims{
		UserID:     userID,
		IdentityID: identityID,
		Provider:   provider,
		JTI:        accessJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        accessJTI,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, errcode.ErrGenerateAccessToken
	}

	// 生成Refresh Token
	refreshExpiresAt := now.Add(s.refreshExpiry)
	refreshClaims := TokenClaims{
		UserID: userID,
		JTI:    refreshJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        refreshJTI,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return nil, errcode.ErrGenerateRefreshToken
	}

	// 将Refresh Token存入数据库（用于撤销）
	metaJSON, _ := json.Marshal(meta)
	dbToken := &models.AuthToken{
		UserID:    userID,
		JTI:       refreshJTI,
		TokenType: models.TokenTypeRefreshToken,
		ExpiresAt: refreshExpiresAt,
		Revoked:   false,
		Meta:      string(metaJSON),
	}

	if err := s.tokenRepo.Create(ctx, dbToken); err != nil {
		return nil, errcode.ErrSaveRefreshToken
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiresAt,
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken 生成访问令牌
func (s *TokenService) GenerateAccessToken(userID, identityID, provider string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessExpiry)
	jti := uuid.New().String()

	claims := TokenClaims{
		UserID:     userID,
		IdentityID: identityID,
		Provider:   provider,
		JTI:        jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, errcode.ErrGenerateAccessToken
	}

	return tokenString, expiresAt, nil
}

// ParseToken 解析Token
func (s *TokenService) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errcode.ErrUnexpectedSigningMethod
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, errcode.ErrParseToken
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errcode.ErrInvalidToken
}

// ValidateToken 验证Token
func (s *TokenService) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查token是否已撤销（仅检查refresh token）
	if claims.ExpiresAt.After(time.Now().Add(s.accessExpiry)) {
		// 可能是refresh token，检查是否已撤销
		revoked, err := s.tokenRepo.IsRevoked(ctx, claims.JTI)
		if err != nil {
			return nil, err
		}
		if revoked {
			return nil, errcode.ErrTokenRevoked
		}
	}

	return claims, nil
}

// RefreshToken 刷新Token
func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// 1. 解析refresh token
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 2. 检查refresh token是否已撤销
	revoked, err := s.tokenRepo.IsRevoked(ctx, claims.JTI)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, errcode.ErrTokenRevoked
	}

	// 3. 撤销旧的refresh token
	if err := s.tokenRepo.Revoke(ctx, claims.JTI); err != nil {
		return nil, err
	}

	// 4. 生成新的token对
	return s.GenerateTokenPair(ctx, claims.UserID, "", "", nil)
}

// RevokeToken 撤销Token
func (s *TokenService) RevokeToken(ctx context.Context, jti string) error {
	return s.tokenRepo.Revoke(ctx, jti)
}

// RevokeAllUserTokens 撤销用户的所有Token
func (s *TokenService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.tokenRepo.RevokeAllByUserID(ctx, userID)
}

// IsTokenRevoked 检查Token是否已撤销
func (s *TokenService) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	return s.tokenRepo.IsRevoked(ctx, jti)
}

// GetTokenExpiry 获取Token过期时间
func (s *TokenService) GetTokenExpiry() time.Duration {
	return s.accessExpiry
}

// GetRefreshTokenExpiry 获取Refresh Token过期时间
func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return s.refreshExpiry
}
