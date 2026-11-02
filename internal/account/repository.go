package account

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

type AccountRepository interface {
	FindAll(ctx context.Context, userID string) ([]Account, error)
	FindByID(ctx context.Context, id, userID string) (*Account, error)
	Insert(ctx context.Context, a *Account) error
	UpdateByID(ctx context.Context, a *Account) error
	DeleteByID(ctx context.Context, id, userID string) error
}

func NewAccountRepository(db *sqlx.DB) AccountRepository { return &repo{db} }

func (r *repo) FindAll(ctx context.Context, userID string) ([]Account, error) {
	var accounts []Account
	err := r.db.SelectContext(ctx, &accounts, `SELECT * FROM accounts WHERE user_id = $1 AND deleted_at IS NULL ORDER BY name ASC`, userID)
	return accounts, err
}

func (r *repo) FindByID(ctx context.Context, id, userID string) (*Account, error) {
	var a Account
	err := r.db.GetContext(ctx, &a, `SELECT * FROM accounts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`, id, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) Insert(ctx context.Context, a *Account) error {
	a.ID = uuid.New().String()
	return r.db.QueryRowxContext(ctx, `INSERT INTO accounts (id, user_id, name, type) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at, revision`, a.ID, a.UserID, a.Name, a.Type).Scan(&a.CreatedAt, &a.UpdatedAt, &a.Revision)
}

func (r *repo) UpdateByID(ctx context.Context, a *Account) error {
	a.UpdatedAt = time.Now()
	rev := a.Revision + 1
	return r.db.QueryRowxContext(ctx, `UPDATE accounts SET name=$1, type=$2, updated_at=NOW(), revision=$3 WHERE id=$4 AND user_id=$5 AND deleted_at IS NULL AND revision=$6 RETURNING updated_at, revision`, a.Name, a.Type, rev, a.ID, a.UserID, a.Revision).Scan(&a.UpdatedAt, &a.Revision)
}

func (r *repo) DeleteByID(ctx context.Context, id, userID string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE accounts SET deleted_at=NOW() WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL`, id, userID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
