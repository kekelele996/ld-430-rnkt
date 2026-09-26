package repository

import "errors"

// Sentinel errors shared across repositories.
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

// IsNotFound reports whether err is the not-found sentinel or wraps it.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
