package expense

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type ExpenseRepository interface {
	FindAll(ctx context.Context, userID string, filter ExpenseFilter) ([]Expense, error)
	FindByID(ctx context.Context, expenseID string, userID string) (*Expense, error)
	Insert(ctx context.Context, e *Expense) error
	UpdateByID(ctx context.Context, e *Expense) error
	DeleteByID(ctx context.Context, expenseID string, userID string) error
}

func NewExpenseRepository(db *sqlx.DB) ExpenseRepository {
	return &repo{db: db}
}

func (r *repo) FindAll(ctx context.Context, userID string, filter ExpenseFilter) ([]Expense, error) {
	query := `
		SELECT id, user_id, category_id, description, amount, occurred_at, 
			   created_at, updated_at, deleted_at, revision 
		FROM expenses 
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{userID}
	argIndex := 2 // Next placeholder will be $2

	if filter.CategoryID != nil && *filter.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		argIndex++
	}
	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND occurred_at >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}
	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND occurred_at <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	query += " ORDER BY occurred_at DESC"

	var expenses []Expense
	err := r.db.SelectContext(ctx, &expenses, query, args...)
	return expenses, err
}

func (r *repo) FindByID(ctx context.Context, expenseID string, userID string) (*Expense, error) {
	query := `
		SELECT id, user_id, category_id, description, amount, occurred_at, 
			   created_at, updated_at, deleted_at, revision 
		FROM expenses 
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	var expense Expense
	err := r.db.GetContext(ctx, &expense, query, expenseID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &expense, nil
}

func (r *repo) Insert(ctx context.Context, e *Expense) error {
	e.ID = uuid.New().String()

	query := `
		INSERT INTO expenses (id, user_id, category_id, description, amount, occurred_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at, revision
	`
	return r.db.QueryRowxContext(ctx, query, e.ID, e.UserID, e.CategoryID, e.Description, e.Amount, e.OccurredAt).Scan(
		&e.CreatedAt, &e.UpdatedAt, &e.Revision,
	)
}

func (r *repo) UpdateByID(ctx context.Context, e *Expense) error {
	newRevision := e.Revision + 1

	query := `
		UPDATE expenses 
		SET category_id = $1, description = $2, amount = $3, occurred_at = $4, 
			updated_at = NOW(), revision = $5
		WHERE id = $6 AND user_id = $7 AND deleted_at IS NULL AND revision = $8
		RETURNING updated_at, revision
	`
	return r.db.QueryRowxContext(ctx, query, e.CategoryID, e.Description, e.Amount, e.OccurredAt, newRevision, e.ID, e.UserID, e.Revision).Scan(
		&e.UpdatedAt, &e.Revision,
	)
}

func (r *repo) DeleteByID(ctx context.Context, expenseID string, userID string) error {
	query := `
		UPDATE expenses 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, query, expenseID, userID)
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
