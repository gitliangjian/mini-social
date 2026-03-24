package handler

import (
	"errors"
	"fmt"
	"mini-social/internal/middleware"
	"mini-social/internal/service"
	"mini-social/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

// 创建评论
func (h *CommentHandler) Create(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
	}
	userID := userIDValue.(uint)

	momentID, err := strconv.ParseUint(c.Param("moment_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid moment_id")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	comment, err := h.commentService.Create(service.CreateCommentInput{
		MomentID: uint(momentID),
		UserID:   userID,
		Content:  req.Content,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCommentContent) {
			response.BadRequest(c, "content cannot be empty or too long")
			return
		}
		response.InternalError(c, "create comment failed")
		return
	}

	response.Success(c, comment)
}

// 获取某条动态的评论列表
func (h *CommentHandler) List(c *gin.Context) {
	momentID, err := strconv.ParseUint(c.Param("moment_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid moment_id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.commentService.ListByMoment(service.ListCommentsInput{
		MomentID: uint(momentID),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		fmt.Println("list comments error:", err)
		response.InternalError(c, "list comments failed")
		return
	}

	response.Success(c, gin.H{
		"list":      result.List,
		"page":      result.Page,
		"page_size": result.PageSize,
	})
}

// 删除自己的评论
func (h *CommentHandler) Delete(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	userID := userIDValue.(uint)

	//区别于List和Create方法，这里需要解析路由中的comment_id参数
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid comment id")
		return
	}

	err = h.commentService.Delete(userID, uint(commentID))
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) || errors.Is(err, service.ErrCommentForbidden) {
			response.NotFound(c, "comment not found or no permission")
			return
		}
		response.InternalError(c, "delete comment failed")
		return
	}

	response.Success(c, gin.H{
		"id":      uint(commentID),
		"deleted": true,
	})
}
