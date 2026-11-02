package account

import (
	"api/internal/store"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const cols = `id, user_id, name, type, created_at, updated_at, revision`

const maxListRows = 1000

type repo struct{ db sqlx.ExtContext }

type AccountRepository interface {
	FindAll(ctx context.Context, userID string) ([]Account, error)
	FindByID(ctx context.Context, id, userID string) (*Account, error)
	Insert(ctx context.Context, a *Account) error
	UpdateByID(ctx context.Context, a *Account) error
	DeleteByID(ctx context.Context, id, userID string) error
}

func NewAccountRepository(db sqlx.ExtContext) AccountRepository { return &repo{db} }

func (r *repo) FindAll(ctx context.Context, userID string) ([]Account, error) {
	accounts := make([]Account, 0, 16)
	q := `SELECT ` + cols + ` FROM accounts
	      WHERE user_id = $1 AND deleted_at IS NULL
	      ORDER BY name ASC, id ASC
	      LIMIT $2`
	if err := sqlx.SelectContext(ctx, r.db, &accounts, q, userID, maxListRows); err != nil {
		return nil, fmt.Errorf("account: find all: %w", err)
	}
	return accounts, nil
}

func (r *repo) FindByID(ctx context.Context, id, userID string) (*Account, error) {
	var a Account
	q := `SELECT ` + cols + ` FROM accounts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	if err := sqlx.GetContext(ctx, r.db, &a, q, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("account: find by id: %w", err)
	}
	return &a, nil
}

func (r *repo) Insert(ctx context.Context, a *Account) error {
	a.Name = strings.TrimSpace(a.Name)
	if err := validateFields(a.Name, a.Type); err != nil {
		return err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("account: new id: %w", err)
	}
	var created, updated time.Time
	var rev int64
	q := `INSERT INTO accounts (id, user_id, name, type) VALUES ($1, $2, $3, $4)
	      RETURNING created_at, updated_at, revision`
	if err := r.db.QueryRowxContext(ctx, q, id.String(), a.UserID, a.Name, a.Type).
		Scan(&created, &updated, &rev); err != nil {
		return fmt.Errorf("account: insert: %w", store.Translate(err))
	}

	a.ID, a.CreatedAt, a.UpdatedAt, a.Revision = id.String(), created, updated, rev
	return nil
}

func (r *repo) UpdateByID(ctx context.Context, a *Account) error {
	a.Name = strings.TrimSpace(a.Name)
	if err := validateFields(a.Name, a.Type); err != nil {
		return err
	}
	q := `UPDATE accounts
	      SET name = $1, type = $2, updated_at = now(), revision = revision + 1
	      WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL AND revision = $5
	      RETURNING updated_at, revision`
	var updated time.Time
	var rev int64
	err := r.db.QueryRowxContext(ctx, q, a.Name, a.Type, a.ID, a.UserID, a.Revision).
		Scan(&updated, &rev)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		exists, perr := r.existsActive(ctx, a.ID, a.UserID)
		if perr != nil {
			return perr
		}
		if exists {
			return store.ErrConflict
		}
		return store.ErrNotFound
	case err != nil:
		return fmt.Errorf("account: update: %w", store.Translate(err))
	}
	a.UpdatedAt, a.Revision = updated, rev
	return nil
}

func (r *repo) DeleteByID(ctx context.Context, id, userID string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE accounts
		 SET deleted_at = now(), updated_at = now(), revision = revision + 1
		 WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`, id, userID)
	if err != nil {
		return fmt.Errorf("account: delete: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("account: delete: rows affected: %w", err)
	}
	if n == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (r *repo) existsActive(ctx context.Context, id, userID string) (bool, error) {
	var ok bool
	q := `SELECT EXISTS(
	        SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`
	if err := sqlx.GetContext(ctx, r.db, &ok, q, id, userID); err != nil {
		return false, fmt.Errorf("account: existence probe: %w", err)
	}
	return ok, nil
}

func validateFields(name, typ string) error {
	switch {
	case name == "":
		return fmt.Errorf("%w: name is required", store.ErrInvalidInput)
	case len(name) > maxNameLen:
		return fmt.Errorf("%w: name exceeds %d characters", store.ErrInvalidInput, maxNameLen)
	case !validAccountType(typ):
		return fmt.Errorf("%w: unknown account type %q", store.ErrInvalidInput, typ)
	}
	return nil
}
