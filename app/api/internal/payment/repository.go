package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

type PaymentRepository interface {
	FindByDebtID(ctx context.Context, debtID, userID string) ([]Payment, error)
	RecordPayment(ctx context.Context, p *Payment) error
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository { return &repo{db} }

func (r *repo) FindByDebtID(ctx context.Context, debtID, userID string) ([]Payment, error) {
	var payments []Payment
	err := r.db.SelectContext(ctx, &payments,
		`SELECT * FROM payments WHERE debt_id = $1 AND user_id = $2 AND deleted_at IS NULL ORDER BY paid_at DESC`,
		debtID, userID)
	return payments, err
}

func (r *repo) RecordPayment(ctx context.Context, p *Payment) error {
	p.ID = uuid.New().String()
	p.PaidAt = time.Now()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO payments (id, user_id, debt_id, amount, note, paid_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		p.ID, p.UserID, p.DebtID, p.Amount, p.Note, p.PaidAt)
	if err != nil {
		return err
	}

	var currentRemaining float64
	var debtType string
	err = tx.QueryRowContext(ctx, `SELECT remaining_amount, type FROM debts WHERE id = $1 AND user_id = $2`, p.DebtID, p.UserID).Scan(&currentRemaining, &debtType)
	if err != nil {
		return err
	}

	newRemaining := currentRemaining - p.Amount
	if newRemaining < 0 {
		newRemaining = 0
	}

	newStatus := "pending"
	var settledAt *time.Time
	if newRemaining == 0 {
		newStatus = "paid"
		now := time.Now()
		settledAt = &now
	} else {
		newStatus = "partial"
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE debts SET remaining_amount = $1, status = $2, settled_at = $3, updated_at = NOW() WHERE id = $4`,
		newRemaining, newStatus, settledAt, p.DebtID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
