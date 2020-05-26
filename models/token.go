package models

import "time"

// TokenRecord JWT token and key's store
type TokenRecord struct {
	ID        int       `json:"id"`
	JwtID     string    `json:"jti"`     // JWT ID, jti
	Key       []byte    `json:"-"`       // Key
	Revoked   bool      `json:"revoked"` // Manually revoked
	Subject   string    `json:"subject"` // User.UUID
	IssuedAt  int64     `json:"iat"`     // 0
	ExpiresAt int64     `json:"exp"`     // 0
	NotBefore int64     `json:"nbf"`     // 0
	Issuer    string    `json:"iss"`
	Audience  string    `json:"aud"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
