package bootstrap

import (
	"fmt"
	"log"

	"mini-social/internal/config"
	"mini-social/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db failed: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database failed: %w", err)
	}

	log.Println("database connected")

	return db, nil
}

// 自动迁移
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Moment{},
		&model.Comment{},
	); err != nil {
		return fmt.Errorf("auto migrate failed:%w", err)
	}

	log.Println("database migrated")
	return nil
}
