package handler

import (
	"errors"
	"mini-social/internal/middleware"
	"mini-social/internal/model"
	"mini-social/internal/service"
	"mini-social/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// 点赞
func (h *LikeHandler) Like(c *gin.Context) {
	//获取当前登录的用户ID
	userIDValue, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}
	userID := userIDValue.(uint)

	//判断点赞是评论还是动态
	commentIDStr := c.Param("comment_id")
	var targetType model.LikeTargetType
	var targetID uint

	if commentIDStr == "" {
		targetType = model.LikeTargetMoment
		id, _ := strconv.ParseUint(c.Param("moment_id"), 10, 64)
		targetID = uint(id)
	} else {
		targetType = model.LikeTargetComment
		id, _ := strconv.ParseUint(commentIDStr, 10, 64)
		targetID = uint(id)
	}

	count, err := h.likeService.Like(service.LikeInput{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   uint(targetID),
	})

	if err != nil {
		if errors.Is(err, service.ErrAlreadyLiked) {
			response.BadRequest(c, "already liked")
			return
		}
		response.InternalError(c, "like failed")
		return
	}

	response.Success(c, gin.H{"liked": true, "likes_count": count})
}

func (h *LikeHandler) UnLike(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	userID := userIDValue.(uint)

	//判断点赞是评论还是动态
	commentIDStr := c.Param("comment_id")
	var targetType model.LikeTargetType
	var targetID uint

	if commentIDStr == "" {
		targetType = model.LikeTargetMoment
		id, _ := strconv.ParseUint(c.Param("moment_id"), 10, 64)
		targetID = uint(id)
	} else {
		targetType = model.LikeTargetComment
		id, _ := strconv.ParseUint(commentIDStr, 10, 64)
		targetID = uint(id)
	}

	count, err := h.likeService.UnLike(service.LikeInput{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   uint(targetID),
	})
	if err != nil {
		response.InternalError(c, "unlike failed")
		return
	}

	response.Success(c, gin.H{"liked": false, "likes_count": count})
}
