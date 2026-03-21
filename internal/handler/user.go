package handler

import (
	"errors"
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

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误: "+err.Error())
		return
	}

	req.Password = strings.TrimSpace(req.Password)

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
