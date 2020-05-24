package db

import (
	"github.com/bin16/tudou/models"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

// Use one connected *gorm.DB
func Use(db *gorm.DB) {
	migrate(db)
	DB = db
}

// AutoMigrate all tables
func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Todo{})
	db.AutoMigrate(&models.Log{})
	db.AutoMigrate(&models.Setting{})
	db.AutoMigrate(&models.Note{})
	db.AutoMigrate(&models.TokenRecord{})
}
