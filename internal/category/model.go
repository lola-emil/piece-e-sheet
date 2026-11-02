package category

import (
	"api/internal/store"
	"fmt"
	"strings"
	"time"
)

type Category struct {
	ID        string     `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Revision  int64      `db:"revision" json:"revision"`
}

const maxNameLen = 100

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

func (r *CreateCategoryRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	switch {
	case r.Name == "":
		return fmt.Errorf("%w: name is required", store.ErrInvalidInput)
	case len(r.Name) > maxNameLen:
		return fmt.Errorf("%w: name exceeds %d characters", store.ErrInvalidInput, maxNameLen)
	}
	return nil
}

func (r *UpdateCategoryRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	switch {
	case r.Name == "":
		return fmt.Errorf("%w: name is required", store.ErrInvalidInput)
	case len(r.Name) > maxNameLen:
		return fmt.Errorf("%w: name exceeds %d characters", store.ErrInvalidInput, maxNameLen)
	}
	return nil
}
