package model

import (
	"time"

	"gorm.io/gorm"
)

type LikeTargetType string

const LikeTargetMoment LikeTargetType = "moment"
const LikeTargetComment LikeTargetType = "comment"

type Like struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//唯一符合索引uniqueIndex，防止点赞重复
	UserID     uint           `gorm:"uniqueIndex:idx_user_target;not null;index" json:"user_id"`
	TargetType LikeTargetType `gorm:"uniqueIndex:idx_user_target;type:varchar(20);not null;index" json:"target_type"` //moment/comment
	TargetID   uint           `gorm:"uniqueIndex:idx_user_target;not null;index" json:"target_id"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` //支持软取消
}

// 自定义表名
func (Like) TableName() string {
	return "likes"
}
