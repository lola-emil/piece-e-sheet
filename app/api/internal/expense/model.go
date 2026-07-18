package expense

import "time"

type Expense struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	CategoryID  *string    `db:"category_id" json:"category_id"` // Nullable
	Description string     `db:"description" json:"description"`
	Amount      float64    `db:"amount" json:"amount"`
	OccurredAt  time.Time  `db:"occurred_at" json:"occurred_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision    int64      `db:"revision" json:"revision"`
}

type CreateExpenseRequest struct {
	CategoryID  *string   `json:"category_id"` // Optional
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type UpdateExpenseRequest struct {
	CategoryID  *string   `json:"category_id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// ExpenseFilter is used for querying expenses
type ExpenseFilter struct {
	CategoryID *string
	StartDate  *time.Time
	EndDate    *time.Time
}
