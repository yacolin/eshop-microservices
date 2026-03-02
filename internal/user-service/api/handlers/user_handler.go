package handlers

import (
	"strconv"

	"eshop-microservices/internal/user-service/api/dto"
	"eshop-microservices/internal/user-service/mq"
	"eshop-microservices/internal/user-service/service"
	"eshop-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc   *service.UserService
	publisher *mq.Publisher
}

func NewUserHandler(userSvc *service.UserService, publisher *mq.Publisher) *UserHandler {
	return &UserHandler{userSvc: userSvc, publisher: publisher}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := h.userSvc.Register(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	if h.publisher != nil {
		h.publisher.PublishUserCreated(user.ID, user.Username, user.Email)
	}

	response.Success(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	tokens, err := h.userSvc.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, tokens)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	user, err := h.userSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := h.userSvc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	if h.publisher != nil {
		h.publisher.PublishUserUpdated(user.ID, user.Username, user.Email)
	}

	response.Success(c, user)
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	userInfo, err := h.userSvc.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, userInfo)
}

func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	var req dto.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	userInfo, err := h.userSvc.UpdateUserInfo(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, userInfo)
}

func (h *UserHandler) Logout(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	if err := h.userSvc.Logout(c.Request.Context(), userID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "logged out"})
}

func (h *UserHandler) getUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Abort()
		return "", nil
	}

	switch v := userID.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case string:
		return v, nil
	default:
		c.Abort()
		return "", nil
	}
}
