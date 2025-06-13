package data

import (
	"context"
	"database/sql"
	"time"
)

type VerificationToken struct {
	Identifier string    `json:"identifier"`
	Token      string    `json:"token"`
	Expires    time.Time `json:"expires"`
	CreatedAt  time.Time `json:"created_at"`
}

type VerificationTokenRepository struct {
	DB *sql.DB
}

func NewVerificationTokenRepository(db *sql.DB) *VerificationTokenRepository {
	return &VerificationTokenRepository{
		DB: db,
	}
}

func (r *VerificationTokenRepository) NewVerificationToken(ctx context.Context, verificationToken VerificationToken) (VerificationToken, error) {
	query := `INSERT INTO verification_token (
		identifier,
		token,
		expires,
		created_at,
	) VALUES ($1, $2, $3, $4) RETURNING identifier`

	var identifier string
	err := r.DB.QueryRowContext(ctx, query,
		verificationToken.Identifier,
		verificationToken.Token,
		verificationToken.Expires,
		time.Now(),
	).Scan(&identifier)

	if err != nil {
		return VerificationToken{}, err
	}

	verificationToken.Identifier = identifier
	return verificationToken, nil
}

func (r *VerificationTokenRepository) DeleteVerificationToken(ctx context.Context, identifier string) error {
	query := `DELETE FROM verification_token WHERE identifier = $1`
	_, err := r.DB.ExecContext(ctx, query, identifier)
	if err != nil {
		return err
	}

	return nil
}
