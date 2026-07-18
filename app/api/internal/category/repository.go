package category

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type CategoryRepository interface {
	FindAll(ctx context.Context, userID string) ([]Category, error)
	FindByID(ctx context.Context, categoryID string) (*Category, error)
	Insert(ctx context.Context, e *Category) error
	UpdateByID(ctx context.Context, category *Category) error
	DeleteByID(ctx context.Context, categoryID string) error
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &repo{db: db}
}

func (r *repo) FindAll(ctx context.Context, userID string) ([]Category, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at, deleted_at, revision 
		FROM categories 
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC
	`

	var categories []Category
	err := r.db.SelectContext(ctx, &categories, query, userID)
	return categories, err
}

func (r *repo) FindByID(ctx context.Context, categoryID string) (*Category, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at, deleted_at, revision 
		FROM categories 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var category Category
	err := r.db.GetContext(ctx, &category, query, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *repo) Insert(ctx context.Context, e *Category) error {
	e.ID = uuid.New().String()

	// Postgres allows us to RETURN the generated timestamps!
	query := `
		INSERT INTO categories (id, user_id, name) 
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at, revision
	`

	return r.db.QueryRowxContext(ctx, query, e.ID, e.UserID, e.Name).Scan(
		&e.CreatedAt, &e.UpdatedAt, &e.Revision,
	)
}

func (r *repo) UpdateByID(ctx context.Context, category *Category) error {
	newRevision := category.Revision + 1

	query := `
		UPDATE categories 
		SET name = $1, updated_at = NOW(), revision = $2
		WHERE id = $3 AND deleted_at IS NULL AND revision = $4
		RETURNING updated_at, revision
	`

	return r.db.QueryRowxContext(ctx, query, category.Name, newRevision, category.ID, category.Revision).Scan(
		&category.UpdatedAt, &category.Revision,
	)
}

func (r *repo) DeleteByID(ctx context.Context, categoryID string) error {
	query := `
		UPDATE categories 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, categoryID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
