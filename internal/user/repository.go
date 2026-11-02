package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (*User, error)
	FindByEmail(ctx context.Context, userEmail string) (*User, error)
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) FindByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT * FROM users WHERE user_id = $1
	`

	var user User

	err := r.db.GetContext(ctx, user, query, userID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) FindByEmail(ctx context.Context, userEmail string) (*User, error) {
	query := `
		SELECT * FROM users WHERE email = $1
	`

	var user User

	err := r.db.GetContext(ctx, user, query, userEmail)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
