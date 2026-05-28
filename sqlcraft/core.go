// Package sqlcraft provides a fluent SQL query builder for PostgreSQL.
package sqlcraft

import "errors"

var (
	ErrEmptyValues      = errors.New("sqlcraft: empty values in query")
	ErrEmptyColumns     = errors.New("sqlcraft: empty columns in query")
	ErrMissMatchValues  = errors.New("sqlcraft: mismatch values for given columns")
	ErrInvalidOperator  = errors.New("sqlcraft: invalid dafi operator")
	ErrInvalidFieldName = errors.New("sqlcraft: invalid field name")
)
