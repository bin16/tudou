package models

import (
	"time"
)

// Event contains text data of todo
type Event struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	UserID    int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
