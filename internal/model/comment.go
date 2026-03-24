package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MomentID  uint      `gorm:"not null;index" json:"moment_id"`         //属于哪条动态
	UserID    uint      `gorm:"not null;index" json:"user_id"`           //属于哪个用户
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"` //自动Preload评论者相关信息
	Content   string    `gorm:"text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` //软删除
}
