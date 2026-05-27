package database

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Sentinel errors returned by this package.
var (
	ErrNotFound         = errors.New("database: not found")
	ErrConflict         = errors.New("database: conflict")
	ErrConnectionFailed = errors.New("database: connection failed")
	ErrHealthCheck      = errors.New("database: health check failed")
	ErrNilPool          = errors.New("database: nil pool")
)

// isUniqueViolation reports whether err is a PostgreSQL unique_violation (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// classifyPgErr maps known pgx/PostgreSQL errors to sentinel-wrapped errors.
// Unknown errors are returned unchanged so the caller can handle them.
func classifyPgErr(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case isUniqueViolation(err):
		return fmt.Errorf("%w: %w", ErrConflict, err)
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}

	return err
}
