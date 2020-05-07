package models

import "time"

// Todo item, one todo for one day
type Setting struct {
	ID        int       `json:"id"`
	UserID    int       `json:"-"`
	Timezone  string    `json:"tz"`
	Language  string    `json:"lang"`
	GameMode  bool      `json:"gameMode"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
