package payment

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

var ErrOverpayment = fmt.Errorf("payment exceeds remaining amount: %w", store.ErrInvalidInput)

const cols = `id, user_id, debt_id, amount, note, paid_at, created_at, updated_at, revision`

const maxListRows = 1000

type repo struct{ db sqlx.ExtContext }

type PaymentRepository interface {
	FindByDebtID(ctx context.Context, debtID, userID string) ([]Payment, error)
	RecordPayment(ctx context.Context, p *Payment) error
}

func NewPaymentRepository(db sqlx.ExtContext) PaymentRepository { return &repo{db} }

func (r *repo) FindByDebtID(ctx context.Context, debtID, userID string) ([]Payment, error) {
	payments := make([]Payment, 0, 8)
	q := `SELECT ` + cols + ` FROM payments
	      WHERE debt_id = $1 AND user_id = $2 AND deleted_at IS NULL
	      ORDER BY paid_at DESC, id DESC
	      LIMIT $3`
	if err := sqlx.SelectContext(ctx, r.db, &payments, q, debtID, userID, maxListRows); err != nil {
		return nil, fmt.Errorf("payment: find by debt: %w", err)
	}
	return payments, nil
}

func (r *repo) RecordPayment(ctx context.Context, p *Payment) error {
	p.Note = strings.TrimSpace(p.Note)
	if err := validatePayment(p.DebtID, p.Amount, p.Note); err != nil {
		return err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("payment: new id: %w", err)
	}

	const q = `
	WITH applied AS (
	    UPDATE debts
	    SET remaining_amount = remaining_amount - $4,
	        status = CASE
	                   WHEN remaining_amount - $4 = 0 THEN 'paid'
	                   ELSE 'partial'
	                 END,
	        settled_at = CASE
	                         WHEN remaining_amount - $4 = 0 THEN COALESCE(settled_at, now())
	                         ELSE NULL
	                     END,
	        updated_at = now(),
	        revision = revision + 1
	    WHERE id = $3 AND user_id = $2 AND deleted_at IS NULL
	      AND remaining_amount >= $4
	    RETURNING id
	)
	INSERT INTO payments (id, user_id, debt_id, amount, note, paid_at)
	SELECT $1, $2, $3, $4, $5, now()
	WHERE EXISTS (SELECT 1 FROM applied)
	RETURNING paid_at, created_at, updated_at, revision`

	var paid, created, updated time.Time
	var rev int64
	err = r.db.QueryRowxContext(ctx, q,
		id.String(), p.UserID, p.DebtID, p.Amount, p.Note,
	).Scan(&paid, &created, &updated, &rev)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return r.classifyMiss(ctx, p.DebtID, p.UserID, p.Amount)
	case err != nil:
		return fmt.Errorf("payment: record: %w", store.Translate(err))
	}
	p.ID, p.PaidAt, p.CreatedAt, p.UpdatedAt, p.Revision = id.String(), paid, created, updated, rev
	return nil
}

func (r *repo) classifyMiss(ctx context.Context, debtID, userID string, amount decimal.Decimal) error {
	var remaining decimal.Decimal
	err := sqlx.GetContext(ctx, r.db, &remaining,
		`SELECT remaining_amount FROM debts
		 WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		debtID, userID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return store.ErrNotFound
	case err != nil:
		return fmt.Errorf("payment: classify miss: %w", err)
	case remaining.LessThan(amount):
		return ErrOverpayment
	default:
		return ErrOverpayment
	}
}

func validatePayment(debtID string, amount decimal.Decimal, note string) error {
	switch {
	case strings.TrimSpace(debtID) == "":
		return fmt.Errorf("%w: debt_id is required", store.ErrInvalidInput)
	case !amount.IsPositive():
		return fmt.Errorf("%w: amount must be greater than zero", store.ErrInvalidInput)
	case len(note) > maxNoteLen:
		return fmt.Errorf("%w: note too long", store.ErrInvalidInput)
	}
	return nil
}
