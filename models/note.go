package models

import "time"

// Note write something done
type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"-"`
	User      User      `json:"-"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
