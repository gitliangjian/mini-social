package service

import (
	"errors"
	"mini-social/internal/model"
	"mini-social/internal/repository"
	"strings"
)

var (
	ErrInvalidMomentContent = errors.New("invalid moment content")
)

type MomentService struct {
	momentRepo *repository.MomentRepository
}

func NewMomentService(momentRepo *repository.MomentRepository) *MomentService {
	return &MomentService{
		momentRepo: momentRepo,
	}
}

type CreateMomentInput struct {
	UserID  uint
	Content string
}

type ListMomentsInput struct {
	Page     int
	PageSize int
}

type ListMomentsResult struct {
	List     []model.Moment
	Page     int
	PageSize int
}

func (s *MomentService) Create(input CreateMomentInput) (*model.Moment, error) {
	content := strings.TrimSpace(input.Content)
	//内容不能为空以及长度限制
	if content == "" || len(content) > 2000 {
		return nil, ErrInvalidMomentContent
	}

	moment := &model.Moment{
		UserID:  input.UserID,
		Content: content,
	}

	if err := s.momentRepo.Create(moment); err != nil {
		return nil, err
	}
	return moment, nil
}

func (s *MomentService) List(input ListMomentsInput) (*ListMomentsResult, error) {
	page := input.Page
	pageSize := input.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	moments, err := s.momentRepo.List(offset, pageSize)
	if err != nil {
		return nil, err
	}

	return &ListMomentsResult{
		List:     moments,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
