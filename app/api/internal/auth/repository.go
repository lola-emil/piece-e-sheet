package auth

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, user *User) error {
	user.ID = uuid.New().String()

	query := `
		INSERT INTO users (id, email, display_name, password_hash) 
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query, user.ID, user.Email, user.DisplayName, user.Password).Scan(
		&user.CreatedAt, &user.UpdatedAt,
	)
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, display_name, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
