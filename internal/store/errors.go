// Package store defines persistence-layer errors shared by all repositories.
// Handlers map these to HTTP status codes; repos wrap them with context.
package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrNotFound: the row doesn't exist, is soft-deleted, or belongs to
	// another user. Repos deliberately do not distinguish these cases.
	ErrNotFound = errors.New("store: not found")

	// ErrConflict: optimistic-lock failure; the client's revision is stale.
	// Re-read and retry. HTTP 409.
	ErrConflict = errors.New("store: revision conflict")

	// ErrDuplicate: a uniqueness constraint was violated. HTTP 409.
	ErrDuplicate = errors.New("store: duplicate value")

	// ErrInvalidInput: input failed validation. HTTP 400.
	ErrInvalidInput = errors.New("store: invalid input")
)

// Translate maps driver errors onto the shared sentinels so callers never
// see driver-specific types. Assumes pgx; for lib/pq match *pq.Error on
// the same SQLSTATE codes.
func Translate(err error) error {
	if err == nil {
		return nil
	}
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		switch pg.Code {
		case "23505": // unique_violation
			return ErrDuplicate
		case "23503": // foreign_key_violation
			return ErrInvalidInput
		}
	}
	return err
}
