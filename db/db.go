package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB *gorm.DB
)

func Open() error {
	var err error
	DB, err = gorm.Open("sqlite3", "data/db.sqlite3")

	return err
}
