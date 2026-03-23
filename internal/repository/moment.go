package repository

import (
	"mini-social/internal/model"

	"gorm.io/gorm"
)

type MomentRepository struct {
	db *gorm.DB
}

func NewMomentRepository(db *gorm.DB) *MomentRepository {
	return &MomentRepository{
		db: db,
	}
}

func (r *MomentRepository) Create(moment *model.Moment) error {
	return r.db.Create(moment).Error
}

// limit限制返回的动态数目，offset表示跳过前多少条记录
func (r *MomentRepository) List(offset, limit int) ([]model.Moment, error) {
	var moments []model.Moment
	err := r.db.
		Preload("User"). // 自动关联查询出用户信息
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&moments).Error
	if err != nil {
		return nil, err
	}
	return moments, nil
}
