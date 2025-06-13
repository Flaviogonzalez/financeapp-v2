package data

import "time"

type Session struct {
	ID           int       `json:"id"`
	SessionToken string    `json:"session_token"`
	UserID       int       `json:"user_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}
