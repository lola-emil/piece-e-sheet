package expense

import "time"

type Expense struct {
	ID          string  `db:"id"`
	UserID      string  `db:"user_id"`
	CategoryID  string  `db:"category_id"`
	Description string  `db:"description"`
	Amount      float32 `db:"amount"`

	OccuredAt time.Time `db:"occured_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Revision int64 `db:"revision"`
}
