package debt

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
	"github.com/shopspring/decimal"
)

const cols = `id, user_id, person_name, description, amount, remaining_amount,
              type, status, due_date, settled_at, created_at, updated_at, revision`

const maxListRows = 1000

type repo struct{ db sqlx.ExtContext }

type DebtRepository interface {
	FindAll(ctx context.Context, userID string) ([]Debt, error)
	FindByID(ctx context.Context, id, userID string) (*Debt, error)
	Insert(ctx context.Context, d *Debt) error
	UpdateByID(ctx context.Context, d *Debt) error
	DeleteByID(ctx context.Context, id, userID string) error
}

func NewDebtRepository(db sqlx.ExtContext) DebtRepository { return &repo{db} }

func (r *repo) FindAll(ctx context.Context, userID string) ([]Debt, error) {
	debts := make([]Debt, 0, 16)
	q := `SELECT ` + cols + ` FROM debts
	      WHERE user_id = $1 AND deleted_at IS NULL
	      ORDER BY created_at DESC, id DESC
	      LIMIT $2`
	if err := sqlx.SelectContext(ctx, r.db, &debts, q, userID, maxListRows); err != nil {
		return nil, fmt.Errorf("debt: find all: %w", err)
	}
	return debts, nil
}

func (r *repo) FindByID(ctx context.Context, id, userID string) (*Debt, error) {
	var d Debt
	q := `SELECT ` + cols + ` FROM debts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	if err := sqlx.GetContext(ctx, r.db, &d, q, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("debt: find by id: %w", err)
	}
	return &d, nil
}

func (r *repo) Insert(ctx context.Context, d *Debt) error {
	d.PersonName = strings.TrimSpace(d.PersonName)
	d.Description = strings.TrimSpace(d.Description)
	if err := validateDebt(d.PersonName, d.Description, d.Amount, d.Amount, d.Type, "pending"); err != nil {
		return err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("debt: new id: %w", err)
	}

	d.RemainingAmount = d.Amount
	d.Status = "pending"

	var created, updated time.Time
	var rev int64
	q := `INSERT INTO debts
	      (id, user_id, person_name, description, amount, remaining_amount, type, status, due_date)
	      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	      RETURNING created_at, updated_at, revision`
	if err := r.db.QueryRowxContext(ctx, q,
		id.String(), d.UserID, d.PersonName, d.Description,
		d.Amount, d.RemainingAmount, d.Type, d.Status, d.DueDate,
	).Scan(&created, &updated, &rev); err != nil {
		return fmt.Errorf("debt: insert: %w", store.Translate(err))
	}
	d.ID, d.CreatedAt, d.UpdatedAt, d.Revision = id.String(), created, updated, rev
	return nil
}

func (r *repo) UpdateByID(ctx context.Context, d *Debt) error {
	d.PersonName = strings.TrimSpace(d.PersonName)
	d.Description = strings.TrimSpace(d.Description)
	if err := validateDebt(d.PersonName, d.Description, d.Amount, d.RemainingAmount, d.Type, d.Status); err != nil {
		return err
	}
	q := `UPDATE debts
	      SET person_name = $1, description = $2, amount = $3, remaining_amount = $4,
	          type = $5, status = $6, due_date = $7, settled_at = $8,
	          updated_at = now(), revision = revision + 1
	      WHERE id = $9 AND user_id = $10 AND deleted_at IS NULL AND revision = $11
	      RETURNING updated_at, revision`
	var updated time.Time
	var rev int64
	err := r.db.QueryRowxContext(ctx, q,
		d.PersonName, d.Description, d.Amount, d.RemainingAmount,
		d.Type, d.Status, d.DueDate, d.SettledAt,
		d.ID, d.UserID, d.Revision,
	).Scan(&updated, &rev)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		exists, perr := r.existsActive(ctx, d.ID, d.UserID)
		if perr != nil {
			return perr
		}
		if exists {
			return store.ErrConflict
		}
		return store.ErrNotFound
	case err != nil:
		return fmt.Errorf("debt: update: %w", store.Translate(err))
	}
	d.UpdatedAt, d.Revision = updated, rev
	return nil
}

func (r *repo) DeleteByID(ctx context.Context, id, userID string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE debts
		 SET deleted_at = now(), updated_at = now(), revision = revision + 1
		 WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`, id, userID)
	if err != nil {
		return fmt.Errorf("debt: delete: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("debt: delete: rows affected: %w", err)
	}
	if n == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (r *repo) existsActive(ctx context.Context, id, userID string) (bool, error) {
	var ok bool
	q := `SELECT EXISTS(
	        SELECT 1 FROM debts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`
	if err := sqlx.GetContext(ctx, r.db, &ok, q, id, userID); err != nil {
		return false, fmt.Errorf("debt: existence probe: %w", err)
	}
	return ok, nil
}

func validateDebt(person, desc string, amount, remaining decimal.Decimal, typ, status string) error {
	switch {
	case person == "":
		return fmt.Errorf("%w: person_name is required", store.ErrInvalidInput)
	case len(person) > maxPersonNameLen:
		return fmt.Errorf("%w: person_name too long", store.ErrInvalidInput)
	case len(desc) > maxDescriptionLen:
		return fmt.Errorf("%w: description too long", store.ErrInvalidInput)
	case !amount.IsPositive():
		return fmt.Errorf("%w: amount must be greater than zero", store.ErrInvalidInput)
	case remaining.IsNegative():
		return fmt.Errorf("%w: remaining_amount cannot be negative", store.ErrInvalidInput)
	case remaining.GreaterThan(amount):
		return fmt.Errorf("%w: remaining_amount cannot exceed amount", store.ErrInvalidInput)
	case !validDebtType(typ):
		return fmt.Errorf("%w: unknown debt type %q", store.ErrInvalidInput, typ)
	case !validDebtStatus(status):
		return fmt.Errorf("%w: unknown debt status %q", store.ErrInvalidInput, status)
	}
	return nil
}
