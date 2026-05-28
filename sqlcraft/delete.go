package sqlcraft

import (
	"strings"

	"github.com/techforge-hq/go-kit/dafi"
)

// DeleteQuery represents a DELETE query.
type DeleteQuery struct {
	table            string
	returningColumns []string

	rawValues []any

	sqlColumnByDomainField map[string]string
	filters                dafi.Filters
}

// DeleteFrom creates a new DeleteQuery targeting the specified table.
func DeleteFrom(table string) DeleteQuery {
	return DeleteQuery{
		table:            table,
		returningColumns: []string{},
	}
}

// Where adds filters to the DELETE query.
func (d DeleteQuery) Where(filters ...dafi.Filter) DeleteQuery {
	d.filters = filters
	return d
}

// SQLColumnByDomainField sets the mapping from domain fields to SQL columns.
func (d DeleteQuery) SQLColumnByDomainField(sqlColumnByDomainField map[string]string) DeleteQuery {
	d.sqlColumnByDomainField = sqlColumnByDomainField
	return d
}

// Returning adds a RETURNING clause to the query.
func (d DeleteQuery) Returning(columns ...string) DeleteQuery {
	d.returningColumns = columns
	return d
}

// ToSQL builds the SQL query and returns the Result.
func (d DeleteQuery) ToSQL() (Result, error) {
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM ")
	builder.WriteString(d.table)

	args := []any{}
	if len(d.filters) > 0 {
		whereResult, err := WhereSafe(len(d.rawValues), d.sqlColumnByDomainField, d.filters...)
		if err != nil {
			return Result{}, err
		}
		args = whereResult.Args
		builder.WriteString(whereResult.SQL)
	}

	if len(d.returningColumns) > 0 {
		builder.WriteString(" RETURNING ")
		builder.WriteString(strings.Join(d.returningColumns, ", "))
	}

	return Result{
		SQL:  builder.String(),
		Args: args,
	}, nil
}
