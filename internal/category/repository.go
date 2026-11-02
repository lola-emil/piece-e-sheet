package category

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

const cols = `id, user_id, name, created_at, updated_at, revision`

const maxListRows = 1000

type repo struct{ db sqlx.ExtContext }

type CategoryRepository interface {
	FindAll(ctx context.Context, userID string) ([]Category, error)
	FindByID(ctx context.Context, id, userID string) (*Category, error)
	Insert(ctx context.Context, c *Category) error
	UpdateByID(ctx context.Context, c *Category) error
	DeleteByID(ctx context.Context, id, userID string) error
}

func NewCategoryRepository(db sqlx.ExtContext) CategoryRepository { return &repo{db} }

func (r *repo) FindAll(ctx context.Context, userID string) ([]Category, error) {
	categories := make([]Category, 0, 16)
	q := `SELECT ` + cols + ` FROM categories
	      WHERE user_id = $1 AND deleted_at IS NULL
	      ORDER BY name ASC, id ASC
	      LIMIT $2`
	if err := sqlx.SelectContext(ctx, r.db, &categories, q, userID, maxListRows); err != nil {
		return nil, fmt.Errorf("category: find all: %w", err)
	}
	return categories, nil
}

func (r *repo) FindByID(ctx context.Context, id, userID string) (*Category, error) {
	var c Category
	q := `SELECT ` + cols + ` FROM categories
	      WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	if err := sqlx.GetContext(ctx, r.db, &c, q, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("category: find by id: %w", err)
	}
	return &c, nil
}

func (r *repo) Insert(ctx context.Context, c *Category) error {
	c.Name = strings.TrimSpace(c.Name)
	if err := validateName(c.Name); err != nil {
		return err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("category: new id: %w", err)
	}
	var created, updated time.Time
	var rev int64
	q := `INSERT INTO categories (id, user_id, name) VALUES ($1, $2, $3)
	      RETURNING created_at, updated_at, revision`
	if err := r.db.QueryRowxContext(ctx, q, id.String(), c.UserID, c.Name).
		Scan(&created, &updated, &rev); err != nil {
		return fmt.Errorf("category: insert: %w", store.Translate(err))
	}
	c.ID, c.CreatedAt, c.UpdatedAt, c.Revision = id.String(), created, updated, rev
	return nil
}

func (r *repo) UpdateByID(ctx context.Context, c *Category) error {
	c.Name = strings.TrimSpace(c.Name)
	if err := validateName(c.Name); err != nil {
		return err
	}
	q := `UPDATE categories
	      SET name = $1, updated_at = now(), revision = revision + 1
	      WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL AND revision = $4
	      RETURNING updated_at, revision`
	var updated time.Time
	var rev int64
	err := r.db.QueryRowxContext(ctx, q, c.Name, c.ID, c.UserID, c.Revision).Scan(&updated, &rev)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		exists, perr := r.existsActive(ctx, c.ID, c.UserID)
		if perr != nil {
			return perr
		}
		if exists {
			return store.ErrConflict
		}
		return store.ErrNotFound
	case err != nil:
		return fmt.Errorf("category: update: %w", store.Translate(err))
	}
	c.UpdatedAt, c.Revision = updated, rev
	return nil
}

func (r *repo) DeleteByID(ctx context.Context, id, userID string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE categories
		 SET deleted_at = now(), updated_at = now(), revision = revision + 1
		 WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`, id, userID)
	if err != nil {
		return fmt.Errorf("category: delete: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("category: delete: rows affected: %w", err)
	}
	if n == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (r *repo) existsActive(ctx context.Context, id, userID string) (bool, error) {
	var ok bool
	q := `SELECT EXISTS(
	        SELECT 1 FROM categories WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`
	if err := sqlx.GetContext(ctx, r.db, &ok, q, id, userID); err != nil {
		return false, fmt.Errorf("category: existence probe: %w", err)
	}
	return ok, nil
}

func validateName(name string) error {
	switch {
	case name == "":
		return fmt.Errorf("%w: name is required", store.ErrInvalidInput)
	case len(name) > maxNameLen:
		return fmt.Errorf("%w: name exceeds %d characters", store.ErrInvalidInput, maxNameLen)
	}
	return nil
}
