package payment

import (
	"api/internal/store"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID        string          `db:"id" json:"id"`
	UserID    string          `db:"user_id" json:"user_id"`
	DebtID    string          `db:"debt_id" json:"debt_id"`
	Amount    decimal.Decimal `db:"amount" json:"amount"`
	Note      string          `db:"note" json:"note"`
	PaidAt    time.Time       `db:"paid_at" json:"paid_at"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision  int64           `db:"revision" json:"revision"`
}

const maxNoteLen = 500

type CreatePaymentRequest struct {
	DebtID string          `json:"debt_id"`
	Amount decimal.Decimal `json:"amount"`
	Note   string          `json:"note"`
}

func (r *CreatePaymentRequest) Validate() error {
	r.Note = strings.TrimSpace(r.Note)
	switch {
	case strings.TrimSpace(r.DebtID) == "":
		return fmt.Errorf("%w: debt_id is required", store.ErrInvalidInput)
	case !r.Amount.IsPositive():
		return fmt.Errorf("%w: amount must be greater than zero", store.ErrInvalidInput)
	case len(r.Note) > maxNoteLen:
		return fmt.Errorf("%w: note too long", store.ErrInvalidInput)
	}
	return nil
}
