package handler

import (
	"errors"
	"mini-social/internal/middleware"
	"mini-social/internal/service"
	"mini-social/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type RegisterRequest struct {
	//自动校验非空、最小长度和最大长度
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	req.Password = strings.TrimSpace(req.Password)
	req.Username = strings.TrimSpace(req.Username)

	user, err := h.userService.Register(service.RegisterInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUsernameAlreadyExists) {
			response.BadRequest(c, "username already exists")
			return
		}
		response.InternalError(c, "register failed")
		return
	}
	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	result, err := h.userService.Login(service.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.BadRequest(c, "用户名或密码错误")
			return
		}
		c.Error(err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": result.Token,
		"user": gin.H{
			"id":       result.User.ID,
			"username": result.User.Username,
		},
	})
}

func (h *UserHandler) Me(c *gin.Context) {
	//从JWT鉴权的中间件获取userID，即当前登录的用户
	userIDValue, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.BadRequest(c, "user id not found in context")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.BadRequest(c, "invalid user id type")
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.BadRequest(c, "user not found")
			return
		}
		response.InternalError(c, "get current user failed")
		// c.Error(err)
		// response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
