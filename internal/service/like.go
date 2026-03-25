package service

import (
	"errors"
	"mini-social/internal/model"
	"mini-social/internal/repository"
)

var (
	ErrLikeTargetNotFound = errors.New("like target not found")
	ErrAlreadyLiked       = errors.New("already liked")
)

type LikeService struct {
	likeRepo *repository.LikeRepository
}

func NewLikeService(likeRepo *repository.LikeRepository) *LikeService {
	return &LikeService{
		likeRepo: likeRepo,
	}
}

type LikeInput struct {
	UserID     uint
	TargetType model.LikeTargetType
	TargetID   uint
}

// 点赞
func (s *LikeService) Like(input LikeInput) (int64, error) {
	//点赞失败返回0
	if err := s.likeRepo.Like(input.UserID, input.TargetType, input.TargetID); err != nil {
		return 0, nil
	}

	//点赞成功返回点赞数
	count, _ := s.likeRepo.CountByTarget(input.TargetType, input.TargetID)
	return count, nil
}

// 取消点赞
func (s *LikeService) UnLike(input LikeInput) (int64, error) {
	if err := s.likeRepo.UnLike(input.UserID, input.TargetType, input.TargetID); err != nil {
		return 0, nil
	}

	count, _ := s.likeRepo.CountByTarget(input.TargetType, input.TargetID)
	return count, nil
}

// 检查是否已经点赞
func (s *LikeService) IsLiked(userID uint, targetType model.LikeTargetType, targetID uint) (bool, error) {
	return s.likeRepo.IsLiked(userID, targetType, targetID)
}
