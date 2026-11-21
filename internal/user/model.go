package user

import "time"

type User struct {
	ID          string `db:"id" json:"id"`
	Email       string `db:"email" json:"email"`
	DisplayName string `db:"display_name" json:"display_name"`

	PasswordHash string `db:"password_hash" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
