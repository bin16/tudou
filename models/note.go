package models

import "time"

// Note write something done
type Note struct {
	ID          int       `json:"id"`
	UserID      int       `json:"-"`
	User        User      `json:"user"`
	Content     string    `json:"content"`
	Title       string    `json:"title"`
	Archived    bool      `json:"archived"`
	Attachments string    `json:"attachments"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"-"`
}
