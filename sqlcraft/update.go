package sqlcraft

import (
	"strconv"
	"strings"

	"github.com/techforge-hq/go-kit/dafi"
)

// UpdateQuery represents an UPDATE query.
type UpdateQuery struct {
	table           string
	columns         []string
	returningValues []string
	values          []any

	isPartialUpdate bool

	sqlColumnByDomainField map[string]string
	filters                dafi.Filters
}

// Update creates a new UpdateQuery targeting the specified table.
func Update(table string) UpdateQuery {
	return UpdateQuery{
		table:           table,
		columns:         []string{},
		returningValues: []string{},
		values:          []any{},
	}
}

// WithColumns sets the columns to be updated.
func (u UpdateQuery) WithColumns(columns ...string) UpdateQuery {
	u.columns = columns
	return u
}

// WithValues sets the values for the update.
func (u UpdateQuery) WithValues(values ...any) UpdateQuery {
	u.values = values
	return u
}

// Where adds filters to the UPDATE query.
func (u UpdateQuery) Where(filters ...dafi.Filter) UpdateQuery {
	u.filters = filters
	return u
}

// SQLColumnByDomainField sets the mapping from domain fields to SQL columns.
func (u UpdateQuery) SQLColumnByDomainField(sqlColumnByDomainField map[string]string) UpdateQuery {
	u.sqlColumnByDomainField = sqlColumnByDomainField
	return u
}

// Returning adds a RETURNING clause to the query.
func (u UpdateQuery) Returning(columns ...string) UpdateQuery {
	u.returningValues = columns
	return u
}

// WithPartialUpdate enables partial updates using COALESCE.
func (u UpdateQuery) WithPartialUpdate() UpdateQuery {
	u.isPartialUpdate = true
	return u
}

// ToSQL builds the SQL query and returns the Result.
func (u UpdateQuery) ToSQL() (Result, error) {
	if len(u.values) > 0 && len(u.values) != len(u.columns) {
		return Result{}, ErrMissMatchValues
	}

	builder := strings.Builder{}
	builder.WriteString("UPDATE ")
	builder.WriteString(u.table)
	builder.WriteString(" SET ")

	for i, column := range u.columns {
		if u.isPartialUpdate {
			builder.WriteString(column)
			builder.WriteString(" = COALESCE($")
			builder.WriteString(strconv.Itoa(i + 1))
			builder.WriteString(", ")
			builder.WriteString(column)
			builder.WriteString(")")
		} else {
			builder.WriteString(column)
			builder.WriteString(" = $")
			builder.WriteString(strconv.Itoa(i + 1))
		}

		if i < len(u.columns)-1 {
			builder.WriteString(", ")
		}
	}

	args := u.values
	if len(u.filters) > 0 {
		whereResult, err := WhereSafe(len(u.values), u.sqlColumnByDomainField, u.filters...)
		if err != nil {
			return Result{}, err
		}
		args = append(args, whereResult.Args...)
		builder.WriteString(whereResult.SQL)
	}

	if len(u.returningValues) > 0 {
		builder.WriteString(" RETURNING ")
		builder.WriteString(strings.Join(u.returningValues, ", "))
	}

	return Result{
		SQL:  builder.String(),
		Args: args,
	}, nil
}
