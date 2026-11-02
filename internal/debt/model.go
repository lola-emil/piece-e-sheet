package debt

import "time"

type Debt struct {
	ID              string     `db:"id" json:"id"`
	UserID          string     `db:"user_id" json:"user_id"`
	PersonName      string     `db:"person_name" json:"person_name"`
	Description     string     `db:"description" json:"description"`
	Amount          float64    `db:"amount" json:"amount"`
	RemainingAmount float64    `db:"remaining_amount" json:"remaining_amount"`
	Type            string     `db:"type" json:"type"`     // "owed_to_me" or "i_owe"
	Status          string     `db:"status" json:"status"` // "pending", "partial", "paid"
	DueDate         *time.Time `db:"due_date" json:"due_date,omitempty"`
	SettledAt       *time.Time `db:"settled_at" json:"settled_at,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision        int64      `db:"revision" json:"revision"`
}

type CreateDebtRequest struct {
	PersonName  string     `json:"person_name"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	Type        string     `json:"type"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type UpdateDebtRequest struct {
	PersonName      string     `json:"person_name"`
	Description     string     `json:"description"`
	Amount          float64    `json:"amount"`
	RemainingAmount float64    `json:"remaining_amount"`
	Type            string     `json:"type"`
	Status          string     `json:"status"`
	DueDate         *time.Time `json:"due_date,omitempty"`
}
