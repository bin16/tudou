package models

import "time"

// Todo item, one todo for one day
type Todo struct {
	ID        int       `json:"id"`
	Event     Event     `json:"event"`
	EventID   int       `json:"-"`
	Status    string    `json:"status"`
	User      User      `json:"-"`
	UserID    int       `json:"-"`
	Date      string    `json:"date"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	TimeStart string    `json:"time"`
}
