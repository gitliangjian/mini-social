package model

import (
	"time"

	"gorm.io/gorm"
)

type Moment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"notnull;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	//实现软删除
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
