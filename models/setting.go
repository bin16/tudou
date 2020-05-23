package models

import "time"

// Setting ...
type Setting struct {
	ID         int       `json:"id"`
	UserID     int       `json:"-"`
	Timezone   string    `json:"tz"`
	Language   string    `json:"lang"`
	GameMode   bool      `json:"gameMode"`
	ThemeColor string    `json:"themeColor"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
