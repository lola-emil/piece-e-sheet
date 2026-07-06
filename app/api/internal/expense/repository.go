package expense

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type ExpenseRepository interface {
	// TODO: pagination
	FindAll(ctx context.Context) ([]Expense, error)
	FindByID(ctx context.Context, expenseID string) (*Expense, error)
	Insert(ctx context.Context, e *Expense) error
	UpdateByID(ctx context.Context, e *Expense) error
	DeleteByID(ctx context.Context, expenseID string) error
}

func (r *repo) FindAll(ctx context.Context) ([]Expense, error) {
	query := `SELECT * FROM expenses`

	var expenses []Expense

	err := r.db.SelectContext(ctx, expenses, query)

	return expenses, err
}

func (r *repo) FindByID(ctx context.Context, expenseID string) (*Expense, error) {
	query := `SELECT * FROM expenses WHERE id = $1`

	var expense Expense

	err := r.db.GetContext(ctx, expense, query, expenseID)

	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *repo) Insert(ctx context.Context, e *Expense) error {
	query := `
	INSERT INTO expenses (
		user_id,
		category_id,
		amount,
		occurred_at,
	) VALUES ($1, $2, $3, $4)
	 RETURNING id, created_at, updated_at, revision
	`
	return r.db.QueryRowxContext(ctx, query,
		e.UserID, e.CategoryID, e.Amount, e.OccuredAt,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.Revision)

}

func (r *repo) UpdateByID(ctx context.Context, e *Expense) error {
	query := `
	UPDATE expenses
	SET user_id = $1,
		category_id = $2,
		amount = $3,
		occured_at = $4
	WHERE id = $5	
	`

	res, err := r.db.ExecContext(ctx, query,
		e.UserID,
		e.CategoryID,
		e.Amount,
		e.OccuredAt,
		e.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repo) DeleteByID(ctx context.Context, expenseID string) error {
	query := `
	UPDATE expenses
	SET deleted_at = $1
	WHERE id = $2
	`

	res, err := r.db.ExecContext(ctx, query, time.Now(), expenseID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
