package data

import "time"

type Authenticator struct {
	CredentialID         string    `json:"credential_id"`
	UserID               int       `json:"user_id"`
	ProviderAccountID    string    `json:"provider_account_id"`
	CredentialPublicKey  string    `json:"credential_public_key"`  // base64 encoded
	Counter              int       `json:"counter"`                // number of times the credential has been used
	CredentialDeviceType string    `json:"credential_device_type"` // "internal" if the credential is backed up
	CredentialBackedUp   bool      `json:"credential_backed_up"`   // true if the credential is backed up
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
