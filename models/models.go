package models

import (
	"github.com/jinzhu/gorm"
	// Sqlite3
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&Todo{})
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Setting{})
	db.AutoMigrate(&Note{})
}