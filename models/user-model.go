package models

import "time"

// User ah, user
type User struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Coins     int64     `json:"coins"`
	GithubID  int       `json:"-"`
	Setting   Setting   `json:"settings"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
