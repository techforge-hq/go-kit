package database

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestClassifyPgErr_Nil(t *testing.T) {
	assert.NoError(t, classifyPgErr(nil))
}

func TestClassifyPgErr_NoRows(t *testing.T) {
	err := classifyPgErr(pgx.ErrNoRows)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.True(t, errors.Is(err, pgx.ErrNoRows))
}

func TestClassifyPgErr_UniqueViolation(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}
	err := classifyPgErr(pgErr)
	assert.True(t, errors.Is(err, ErrConflict))

	var target *pgconn.PgError
	assert.True(t, errors.As(err, &target))
}

func TestClassifyPgErr_UnknownError(t *testing.T) {
	orig := errors.New("something else")
	err := classifyPgErr(orig)
	assert.Equal(t, orig, err)
}

func TestIsUniqueViolation(t *testing.T) {
	assert.True(t, isUniqueViolation(&pgconn.PgError{Code: "23505"}))
	assert.False(t, isUniqueViolation(&pgconn.PgError{Code: "23503"}))
	assert.False(t, isUniqueViolation(errors.New("not a pg error")))
}
