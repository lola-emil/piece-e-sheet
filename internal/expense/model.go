package expense

import (
	"api/internal/store"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	ID          string          `db:"id" json:"id"`
	UserID      string          `db:"user_id" json:"user_id"`
	CategoryID  *string         `db:"category_id" json:"category_id"`
	AccountID   *string         `db:"account_id" json:"account_id"`
	Description string          `db:"description" json:"description"`
	Type        string          `db:"type" json:"type"`
	Amount      decimal.Decimal `db:"amount" json:"amount"`
	OccurredAt  time.Time       `db:"occurred_at" json:"occurred_at"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision    int64           `db:"revision" json:"revision"`
}

const maxDescriptionLen = 500

func validExpenseType(t string) bool {
	return t == "expense" || t == "income"
}

type CreateExpenseRequest struct {
	AccountID   *string         `json:"account_id"`
	CategoryID  *string         `json:"category_id"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	OccurredAt  time.Time       `json:"occurred_at"`
}

func (r *CreateExpenseRequest) Validate() error {
	r.Description = strings.TrimSpace(r.Description)
	switch {
	case r.Description == "":
		return fmt.Errorf("%w: description is required", store.ErrInvalidInput)
	case len(r.Description) > maxDescriptionLen:
		return fmt.Errorf("%w: description too long", store.ErrInvalidInput)
	case !r.Amount.IsPositive():
		return fmt.Errorf("%w: amount must be greater than zero", store.ErrInvalidInput)
	case !validExpenseType(r.Type):
		return fmt.Errorf("%w: unknown expense type %q", store.ErrInvalidInput, r.Type)
	}
	return nil
}

type UpdateExpenseRequest struct {
	AccountID   *string         `json:"account_id"`
	CategoryID  *string         `json:"category_id"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Revision    int64           `json:"revision"`
}

func (r *UpdateExpenseRequest) Validate() error {
	r.Description = strings.TrimSpace(r.Description)
	switch {
	case r.Description == "":
		return fmt.Errorf("%w: description is required", store.ErrInvalidInput)
	case len(r.Description) > maxDescriptionLen:
		return fmt.Errorf("%w: description too long", store.ErrInvalidInput)
	case !r.Amount.IsPositive():
		return fmt.Errorf("%w: amount must be greater than zero", store.ErrInvalidInput)
	case !validExpenseType(r.Type):
		return fmt.Errorf("%w: unknown expense type %q", store.ErrInvalidInput, r.Type)
	case r.OccurredAt.IsZero():
		return fmt.Errorf("%w: occurred_at is required", store.ErrInvalidInput)
	}
	return nil
}

type ExpenseFilter struct {
	CategoryID   *string
	AccountID    *string
	StartDate    *time.Time
	EndDate      *time.Time
	MinAmount    *float64
	MaxAmount    *float64
	Search       string
	SortBy       string
	Limit        int
	Offset       int
	NoPagination bool
}

type ExpenseSummary struct {
	Count              int64            `db:"count" json:"count"`
	Total              decimal.Decimal  `db:"total" json:"total"`
	Average            decimal.Decimal  `db:"average" json:"average"`
	LargestAmount      *decimal.Decimal `db:"largest_amount" json:"largest_amount"`
	LargestDescription *string          `db:"largest_description" json:"largest_description"`
}
