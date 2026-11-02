package debt

import (
	"api/internal/store"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Debt struct {
	ID              string          `db:"id" json:"id"`
	UserID          string          `db:"user_id" json:"user_id"`
	PersonName      string          `db:"person_name" json:"person_name"`
	Description     string          `db:"description" json:"description"`
	Amount          decimal.Decimal `db:"amount" json:"amount"`
	RemainingAmount decimal.Decimal `db:"remaining_amount" json:"remaining_amount"`
	Type            string          `db:"type" json:"type"`
	Status          string          `db:"status" json:"status"`
	DueDate         *time.Time      `db:"due_date" json:"due_date,omitempty"`
	SettledAt       *time.Time      `db:"settled_at" json:"settled_at,omitempty"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision        int64           `db:"revision" json:"revision"`
}

const (
	maxPersonNameLen  = 100
	maxDescriptionLen = 500
)

func validDebtType(t string) bool {
	return t == "owed_to_me" || t == "i_owe"
}

func validDebtStatus(s string) bool {
	return s == "pending" || s == "partial" || s == "paid"
}

type CreateDebtRequest struct {
	PersonName  string          `json:"person_name"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Type        string          `json:"type"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
}

func (r *CreateDebtRequest) Validate() error {
	r.PersonName = strings.TrimSpace(r.PersonName)
	r.Description = strings.TrimSpace(r.Description)
	switch {
	case r.PersonName == "":
		return fmt.Errorf("%w: person_name is required", store.ErrInvalidInput)
	case len(r.PersonName) > maxPersonNameLen:
		return fmt.Errorf("%w: person_name too long", store.ErrInvalidInput)
	case len(r.Description) > maxDescriptionLen:
		return fmt.Errorf("%w: description too long", store.ErrInvalidInput)
	case !r.Amount.IsPositive():
		return fmt.Errorf("%w: amount must be greater than zero", store.ErrInvalidInput)
	case !validDebtType(r.Type):
		return fmt.Errorf("%w: unknown debt type %q", store.ErrInvalidInput, r.Type)
	}
	return nil
}

type UpdateDebtRequest struct {
	PersonName      string          `json:"person_name"`
	Description     string          `json:"description"`
	Amount          decimal.Decimal `json:"amount"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	Type            string          `json:"type"`
	Status          string          `json:"status"`
	DueDate         *time.Time      `json:"due_date,omitempty"`
	Revision        int64           `json:"revision"`
}

func (r *UpdateDebtRequest) Validate() error {
	r.PersonName = strings.TrimSpace(r.PersonName)
	r.Description = strings.TrimSpace(r.Description)
	switch {
	case r.PersonName == "":
		return fmt.Errorf("%w: person_name is required", store.ErrInvalidInput)
	case len(r.PersonName) > maxPersonNameLen:
		return fmt.Errorf("%w: person_name too long", store.ErrInvalidInput)
	case len(r.Description) > maxDescriptionLen:
		return fmt.Errorf("%w: description too long", store.ErrInvalidInput)
	case !r.Amount.IsPositive():
		return fmt.Errorf("%w: amount must be greater than zero", store.ErrInvalidInput)
	case r.RemainingAmount.IsNegative():
		return fmt.Errorf("%w: remaining_amount cannot be negative", store.ErrInvalidInput)
	case r.RemainingAmount.GreaterThan(r.Amount):
		return fmt.Errorf("%w: remaining_amount cannot exceed amount", store.ErrInvalidInput)
	case !validDebtType(r.Type):
		return fmt.Errorf("%w: unknown debt type %q", store.ErrInvalidInput, r.Type)
	case !validDebtStatus(r.Status):
		return fmt.Errorf("%w: unknown debt status %q", store.ErrInvalidInput, r.Status)
	}
	return nil
}
