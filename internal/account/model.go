package account

import (
	"api/internal/store"
	"fmt"
	"strings"
	"time"
)

type Account struct {
	ID        string     `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Name      string     `db:"name" json:"name"`
	Type      string     `db:"type" json:"type"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision  int64      `db:"revision" json:"revision"`
}

type CreateAccountRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UpdateAccountRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

const maxNameLen = 100

func validAccountType(t string) bool {
	switch t {
	case "checking", "savings", "credit":
		return true
	}
	return false
}

func (r *CreateAccountRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	switch {
	case r.Name == "":
		return fmt.Errorf("%w: name is required", store.ErrInvalidInput)
	case len(r.Name) > maxNameLen:
		return fmt.Errorf("%w: name exceeds %d characters", store.ErrInvalidInput, maxNameLen)
	case !validAccountType(r.Type):
		return fmt.Errorf("%w: unknown account type %q", store.ErrInvalidInput, r.Type)
	}
	return nil
}
