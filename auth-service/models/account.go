package data

import "time"

type Account struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"user_id"`
	Type                  string    `json:"type"`
	Provider              string    `json:"provider"`
	ProviderAccountID     string    `json:"provider_account_id"`
	RefreshToken          string    `json:"refresh_token"`
	AccessToken           string    `json:"access_token"`
	ExpiresAt             int       `json:"expires_at"`
	TokenType             string    `json:"token_type"`
	Scope                 string    `json:"scope"`
	IDToken               string    `json:"id_token"`
	SessionState          string    `json:"session_state"`
	RefreshTokenExpiresAt int       `json:"refresh_token_expires_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	CreatedAt             time.Time `json:"created_at"`
}
