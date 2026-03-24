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

// 通过Preload拿到User信息
func (r *MomentRepository) GetByID(id uint) (*model.Moment, error) {
	var moment model.Moment
	err := r.db.Preload("User").First(&moment, id).Error
	if err != nil {
		return nil, err
	}
	return &moment, nil
}

// // 只查询该动态所属的用户ID,减少查询开销
// func (r *MomentRepository) GetUserIDByID(id uint) (uint, error) {
// 	//用结构体进行First查询和接收
// 	var result struct {
// 		UserID uint
// 	}
// 	// Select("user_id") 不Preload，不取Content
// 	err := r.db.Model(&model.Moment{}).Select("user_id").Where("id = ?", id).First(&result).Error
// 	return result.UserID, err
// }

// // 删除动态
// func (r *MomentRepository) Delete(id uint) error {
// 	return r.db.Delete(&model.Moment{}, id).Error
// }

// 只有当 ID 匹配且所有者匹配时才删除
func (r *MomentRepository) DeleteWithAuth(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Moment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound //没找到
	}
	return nil
}
