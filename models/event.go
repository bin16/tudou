package models

import (
	"time"
)

// Event contains text data of todo
type Event struct {
	ID          int       `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"-"`
	UserID      int       `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
