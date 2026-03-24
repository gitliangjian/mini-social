package service

import (
	"errors"
	"mini-social/internal/model"
	"mini-social/internal/repository"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrInvalidMomentContent = errors.New("invalid moment content")
	ErrMomentNotFound       = errors.New("moment not found")
	ErrMomentForbidden      = errors.New("no permission to delete this moment")
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

func (s *MomentService) GetByID(id uint) (*model.Moment, error) {
	moment, err := s.momentRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMomentNotFound
		}
		return nil, err
	}
	return moment, nil
}

func (s *MomentService) Delete(userID, momentID uint) error {
	//直接尝试删除属于该用户的特定动态
	result := s.momentRepo.DeleteWithAuth(momentID, userID)
	//不区分无权限和动态不存在的情况
	return result
}

// //查找动态对应的userid
// id, err := s.momentRepo.GetUserIDByID(momentID)
// if err != nil {
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return ErrMomentNotFound
// 	}
// 	return err
// }
// //比较查到的id和handler层传下来的userID是否相同
// if id != userID {
// 	return ErrMomentForbidden
// }

// //删除
// if err := s.momentRepo.Delete(momentID); err != nil {
// 	return err
// }
// return nil
