package handlers

import (
	"strconv"

	"eshop-microservices/internal/user-service/api/dto"
	"eshop-microservices/internal/user-service/mq"
	"eshop-microservices/internal/user-service/service"
	"eshop-microservices/pkg/errcode"
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

// GetProfile 获取用户资料（包含 User 和 UserInfo）
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

// GetUserInfo 获取用户详细信息
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

// UpdateUserInfo 更新用户详细信息（Avatar、Nickname 等）
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

	// 发布用户更新事件（使用 Nickname 代替原来的 Username）
	if h.publisher != nil {
		h.publisher.PublishUserUpdated(userID, userInfo.Nickname, "")
	}

	response.Success(c, userInfo)
}

// GetByID 根据ID获取用户信息（管理员接口）
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, user)
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
