package category

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, categoryID string) (*Category, error)
	Insert(ctx context.Context, e *Category) error

	// UpdateByID(ctx context.Context, categoryID string) error
	// DeleteByID(ctx context.Context, categoryID string) error
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) FindAll(ctx context.Context) ([]Category, error) {
	query := `
		SELECT * FROM categories
	`

	var category []Category

	err := r.db.SelectContext(ctx, category, query)

	return category, err
}

func (r *repo) FindByID(ctx context.Context, categoryID string) (*Category, error) {
	query := `
		SELECT * FROM categories WHERE id = $1
	`

	var category Category

	err := r.db.GetContext(ctx, category, query, categoryID)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *repo) Insert(ctx context.Context, e *Category) error {
	query := `
	INSERT INTO categories (
		user_id,
		name
	) VALUES ($1, $2)
	RETURNING id, created_at, updated_at, revision
	`

	return r.db.QueryRowxContext(ctx, query, e.UserID, e.Name).Scan(
		&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.Revision,
	)
}

// func (r *repo) UpdateByID(ctx context.Context, categoryID string) error {
// 	query := `
// 	`
// }
// func (r *repo) DeleteByID(ctx context.Context, categoryID string) error {

// }
