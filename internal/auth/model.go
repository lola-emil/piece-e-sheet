package auth

import "time"

type User struct {
	ID          string    `db:"id" json:"id"`
	Email       string    `db:"email" json:"email"`
	DisplayName string    `db:"display_name" json:"display_name"`
	Password    string    `db:"password_hash" json:"-"` // Never expose password in JSON
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
