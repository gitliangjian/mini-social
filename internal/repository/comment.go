package repository

import (
	"mini-social/internal/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// 列出某条动态下的所有评论（带分页+Preload 用户信息）
func (r *CommentRepository) ListByMomentID(momentID uint, offset, limit int) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.
		Preload("User").
		Where("moment_id=?", momentID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) DeleteWithAuth(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
