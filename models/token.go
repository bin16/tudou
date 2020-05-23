package models

import "time"

// TokenRecord JWT token and key's store
type TokenRecord struct {
	ID        int       `json:"id"`
	Token     string    `json:"-"`
	Key       []byte    `json:"-"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
