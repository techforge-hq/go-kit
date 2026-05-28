package sqlcraft

import (
	"strconv"
	"strings"
)

// InsertQuery represents an INSERT query.
type InsertQuery struct {
	table            string
	columns          []string
	returningColumns []string
	values           []any
}

// InsertInto creates a new InsertQuery targeting the specified table.
func InsertInto(tableName string) InsertQuery {
	return InsertQuery{
		table:  tableName,
		values: make([]any, 0),
	}
}

// WithColumns sets the columns for the INSERT query.
func (i InsertQuery) WithColumns(columns ...string) InsertQuery {
	i.columns = columns
	return i
}

// WithValues adds values to the INSERT query.
func (i InsertQuery) WithValues(values ...any) InsertQuery {
	i.values = append(i.values, values...)
	return i
}

// Returning adds a RETURNING clause to the query.
func (i InsertQuery) Returning(columns ...string) InsertQuery {
	i.returningColumns = columns
	return i
}

// ToSQL builds the SQL query and returns the Result.
func (i InsertQuery) ToSQL() (Result, error) {
	if len(i.values) == 0 {
		return Result{}, ErrEmptyValues
	}

	if len(i.values)%len(i.columns) != 0 {
		return Result{}, ErrMissMatchValues
	}

	builder := strings.Builder{}
	builder.WriteString("INSERT INTO ")
	builder.WriteString(i.table)
	builder.WriteString(" (")
	builder.WriteString(strings.Join(i.columns, ", "))
	builder.WriteString(") VALUES ")

	valueRowCount := 0
	for index := range i.values {
		valueRowCount++

		if valueRowCount == 1 && index > 0 {
			builder.WriteString(", ")
		}

		if valueRowCount == 1 {
			builder.WriteString("(")
		}

		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(index + 1))

		if valueRowCount == len(i.columns) {
			builder.WriteString(")")
			valueRowCount = 0
		} else {
			builder.WriteString(", ")
		}
	}

	if len(i.returningColumns) > 0 {
		builder.WriteString(" RETURNING ")
		builder.WriteString(strings.Join(i.returningColumns, ", "))
	}

	return Result{
		SQL:  builder.String(),
		Args: i.values,
	}, nil
}
