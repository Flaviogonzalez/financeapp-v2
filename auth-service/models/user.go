package data

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	Password        string    `json:"password"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (User, error) {
	query := `
		SELECT id, name, email, email_verified_at, password, updated_at, created_at
		FROM users
		WHERE id = $1
	`
	var user User
	if err := r.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerifiedAt, &user.Password, &user.UpdatedAt, &user.CreatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT id, name, email, email_verified_at, password, updated_at, created_at
		FROM users
		WHERE email = $1
	`
	var user User
	if err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerifiedAt, &user.Password, &user.UpdatedAt, &user.CreatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user User) (User, error) {
	query := `
		INSERT INTO users (name, email, email_verified_at, password, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int
	err := r.DB.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.EmailVerifiedAt,
		user.Password,
		time.Now(),
		time.Now(),
	).Scan(&id)
	if err != nil {
		return User{}, err
	}
	user.ID = id
	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user User) (User, error) {
	query := `
		UPDATE users
		SET name = $1, email = $2, email_verified_at = $3, password = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, name, email, email_verified_at, password, updated_at, created_at
	`
	var updatedUser User
	err := r.DB.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.EmailVerifiedAt,
		user.Password,
		user.UpdatedAt,
		user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.EmailVerifiedAt,
		&updatedUser.Password,
		time.Now(),
		&updatedUser.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, name, email, email_verified_at, password, updated_at, created_at
		FROM users
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerifiedAt, &user.Password, &user.UpdatedAt, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
