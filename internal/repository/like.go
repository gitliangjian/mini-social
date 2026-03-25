package repository

import (
	"mini-social/internal/model"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

// 点赞，创建或恢复软删除
func (r *LikeRepository) Like(userID uint, targetType model.LikeTargetType, targetID uint) error {
	like := &model.Like{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}

	//尝试恢复软删除记录，Unscoped用于关闭软删除过滤器，从而可以查到被软删除的记录
	result := r.db.Unscoped().Where("user_id=? AND target_type=? AND target_id = ?", userID, targetType, targetID).First(&model.Like{})
	if result.Error == nil {
		//点赞记录存在，First查询后结果会存放在result.Statement.Dest中，*model.Like做类型断言，然后取出该点赞记录对应的id
		return r.db.Unscoped().Model(&model.Like{}).Where("id = ?", result.Statement.Dest.(*model.Like).ID).Update("deleted_at", nil).Error
	}

	//没有点赞记录，创建点赞
	return r.db.Create(like).Error
}

// 取消点赞，软删除
func (r *LikeRepository) UnLike(userID uint, targetType model.LikeTargetType, targetID uint) error {
	return r.db.Where("user_id=? AND target_type=? AND target_id = ?", userID, targetType, targetID).Delete(&model.Like{}).Error
}

// 根据targetType获取动态/评论的点赞数
func (r *LikeRepository) CountByTarget(targetType model.LikeTargetType, targetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("target_type = ? AND target_id = ? AND deleted_at IS NULL", targetType, targetID).
		Count(&count).Error
	return count, err
}

// 检查用户是否已点赞
func (r *LikeRepository) IsLiked(userID uint, targetType model.LikeTargetType, targetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		Count(&count).Error
	return count > 0, err
}
