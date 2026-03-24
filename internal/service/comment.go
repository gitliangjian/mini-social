package service

import (
	"errors"
	"mini-social/internal/model"
	"mini-social/internal/repository"
	"strings"
)

var (
	ErrInvalidCommentContent = errors.New("invalid comment content")
	ErrCommentNotFound       = errors.New("comment not found")
	ErrCommentForbidden      = errors.New("no permission to delete this comment")
)

type CommentService struct {
	commentRepo *repository.CommentRepository
}

func NewCommentService(commentRepo *repository.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

type CreateCommentInput struct {
	MomentID uint
	UserID   uint
	Content  string
}

type ListCommentsInput struct {
	MomentID uint
	Page     int
	PageSize int
}

type ListCommentsResult struct {
	List     []model.Comment
	Page     int
	PageSize int
}

// 创建评论
func (s *CommentService) Create(input CreateCommentInput) (*model.Comment, error) {
	content := strings.TrimSpace(input.Content)
	if content == "" || len(content) > 500 {
		return nil, ErrInvalidCommentContent
	}

	comment := &model.Comment{
		MomentID: input.MomentID,
		UserID:   input.UserID,
		Content:  content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) ListByMoment(input ListCommentsInput) (*ListCommentsResult, error) {
	page := input.Page
	pageSize := input.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	comments, err := s.commentRepo.ListByMomentID(input.MomentID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	return &ListCommentsResult{
		List:     comments,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *CommentService) Delete(userID, commentID uint) error {
	return s.commentRepo.DeleteWithAuth(commentID, userID)
}
