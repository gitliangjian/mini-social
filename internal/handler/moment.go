package handler

import (
	"errors"
	"mini-social/internal/middleware"
	"mini-social/internal/service"
	"mini-social/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MomentHandler struct {
	momentService *service.MomentService
}

func NewMomentHandler(momentService *service.MomentService) *MomentHandler {
	return &MomentHandler{
		momentService: momentService,
	}
}

type CreateMomentRequest struct {
	Content string `json:"content"`
}

// 发布动态
func (h *MomentHandler) Create(c *gin.Context) {
	//获取userID
	userIDValue, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.InternalError(c, "invalid user id type")
		return
	}

	var req CreateMomentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	//调用service
	moment, err := h.momentService.Create(service.CreateMomentInput{
		UserID:  uint(userID),
		Content: req.Content,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidMomentContent) {
			response.BadRequest(c, "content cannot be empty")
			return
		}
		response.InternalError(c, "create moment failed")
		return
	}

	//响应
	response.Success(c, gin.H{
		"id":         moment.ID,
		"user_id":    moment.UserID,
		"content":    moment.Content,
		"created_at": moment.CreatedAt,
		"updated_at": moment.UpdatedAt,
	})
}

func (h *MomentHandler) List(c *gin.Context) {
	//获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.momentService.List(service.ListMomentsInput{
		Page:     page,
		PageSize: pageSize,
	})

	if err != nil {
		response.InternalError(c, "list moments failed")
		return
	}

	response.Success(c, gin.H{
		"list":      result.List,
		"page":      result.Page,
		"page_size": result.PageSize,
		"total":     len(result.List), // 增加一个当前页数量，方便前端
	})
}

// 动态详情
func (h *MomentHandler) Detail(c *gin.Context) {
	//使用无符号整数解析ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid moment id")
		return
	}

	moment, err := h.momentService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrMomentNotFound) {
			response.NotFound(c, "moment not found")
			return
		}
		response.InternalError(c, "get moment detail failed")
		return
	}

	// data := gin.H{
	// 	"id":         moment.ID,
	// 	"user_id":    moment.UserID,
	// 	"content":    moment.Content,
	// 	"created_at": moment.CreatedAt,
	// 	"updated_at": moment.UpdatedAt,
	// }

	// if moment.User != nil {
	// 	data["user"] = gin.H{
	// 		"id":       moment.User.ID,
	// 		"username": moment.User.Username,
	// 	}
	// }

	// response.Success(c, data)

	//直接返回结构体，利用json Tag自动格式化
	response.Success(c, moment)
}

// 删除动态
func (h *MomentHandler) Delete(c *gin.Context) {
	//获取UserID
	userIDValue, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.InternalError(c, "invalid user id type")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid moment id")
		return
	}

	err = h.momentService.Delete(userID, uint(id))
	if err != nil {
		if errors.Is(err, service.ErrMomentNotFound) {
			response.NotFound(c, "moment not found")
			return
		}
		if errors.Is(err, service.ErrMomentForbidden) {
			response.Forbidden(c, "no permission to delete this moment")
			return
		}
		response.InternalError(c, "delete moment failed")
		return
	}

	response.Success(c, gin.H{
		"id":      uint(id),
		"deleted": true,
	})
}
