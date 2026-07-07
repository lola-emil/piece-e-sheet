package expense

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	CategoryID  string    `db:"category_id" json:"category_id"`
	Description string    `db:"description" json:"description"`
	Amount      float32   `db:"amount" json:"amount"`

	OccuredAt time.Time `db:"occurred_at" json:"occurred_at"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	DeletedAt *time.Time `db:"deleted_at" json:"-"`

	Revision int64 `db:"revision"`
}
