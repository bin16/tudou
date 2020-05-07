package models

import "time"

type Log struct {
	ID        int       `json:"id"`
	EventID   int       `json:"-"`
	Event     Event     `json:"event"`
	Action    string    `json:"action"`
	TodoID    int       `json:"-"`
	UserID    int       `json:"-"`
	Coins     int64     `json:"coins"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"-"`
}
