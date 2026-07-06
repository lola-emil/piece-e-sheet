package category

import "time"

type Category struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`

	Name string `db:"name"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	DeletedAt *time.Time `db:"deleted_at"`

	Revision int64 `db:"revision"`
}
