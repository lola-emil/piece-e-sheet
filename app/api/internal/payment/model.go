package payment

import "time"

type Payment struct {
	ID        string     `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"user_id"`
	DebtID    string     `db:"debt_id" json:"debt_id"`
	Amount    float64    `db:"amount" json:"amount"`
	Note      string     `db:"note" json:"note"`
	PaidAt    time.Time  `db:"paid_at" json:"paid_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision  int64      `db:"revision" json:"revision"`
}

type CreatePaymentRequest struct {
	DebtID string  `json:"debt_id"`
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}
