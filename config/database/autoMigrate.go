package database

import (
	"gorm.io/gorm"
	"live-chat/app/models"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		models.Config{},
		models.User{},
	)
}
