package sqlcraft

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/techforge-hq/go-kit/dafi"
)

// JoinType represents the type of SQL join.
type JoinType string

const (
	InnerJoinType JoinType = "INNER JOIN"
	LeftJoinType  JoinType = "LEFT JOIN"
	RightJoinType JoinType = "RIGHT JOIN"
)

// Join represents a SQL JOIN clause.
type Join struct {
	Type      JoinType
	Table     string
	Condition string
}

// SelectQuery represents a SELECT query.
type SelectQuery struct {
	table                  string
	columns                []string
	requiredColumns        map[string]struct{}
	sqlColumnByDomainField map[string]string

	filters    dafi.Filters
	sorts      dafi.Sorts
	pagination dafi.Pagination

	groups []string
	joins  []Join
}

// Select creates a new SelectQuery with the specified columns.
func Select(columns ...string) SelectQuery {
	return SelectQuery{
		table:           "",
		columns:         columns,
		requiredColumns: make(map[string]struct{}),
	}
}

// From sets the table for the SELECT query.
func (s SelectQuery) From(table string) SelectQuery {
	s.table = table
	return s
}

// Where adds filters to the SELECT query.
func (s SelectQuery) Where(filters ...dafi.Filter) SelectQuery {
	s.filters = filters
	return s
}

// OrderBy adds sort instructions to the SELECT query.
func (s SelectQuery) OrderBy(sorts ...dafi.Sort) SelectQuery {
	s.sorts = sorts
	return s
}

// Limit sets the limit for the SELECT query.
func (s SelectQuery) Limit(limit uint) SelectQuery {
	s.pagination.PageSize = limit
	return s
}

// Page sets the page number for the SELECT query.
func (s SelectQuery) Page(page uint) SelectQuery {
	s.pagination.PageNumber = page
	return s
}

// RequiredColumns allows selecting just a subset of the columns provided in Select.
func (s SelectQuery) RequiredColumns(columns ...string) SelectQuery {
	for _, col := range columns {
		s.requiredColumns[col] = struct{}{}
	}
	return s
}

// SQLColumnByDomainField sets the mapping from domain fields to SQL columns.
func (s SelectQuery) SQLColumnByDomainField(sqlColumnByDomainField map[string]string) SelectQuery {
	s.sqlColumnByDomainField = sqlColumnByDomainField
	return s
}

// InnerJoin adds an INNER JOIN to the query.
func (s SelectQuery) InnerJoin(table, condition string) SelectQuery {
	return s.addJoin(InnerJoinType, table, condition)
}

// LeftJoin adds a LEFT JOIN to the query.
func (s SelectQuery) LeftJoin(table, condition string) SelectQuery {
	return s.addJoin(LeftJoinType, table, condition)
}

// RightJoin adds a RIGHT JOIN to the query.
func (s SelectQuery) RightJoin(table, condition string) SelectQuery {
	return s.addJoin(RightJoinType, table, condition)
}

func (s SelectQuery) addJoin(joinType JoinType, table, condition string) SelectQuery {
	s.joins = append(s.joins, Join{
		Type:      joinType,
		Table:     table,
		Condition: condition,
	})
	return s
}

// ToSQL builds the SQL query and returns the Result.
func (s SelectQuery) ToSQL() (Result, error) {
	if len(s.columns) == 0 {
		return Result{}, ErrEmptyColumns
	}

	if len(s.sqlColumnByDomainField) > 0 {
		requiredCols := make(map[string]struct{})
		for k := range s.requiredColumns {
			requiredSQLColumn, ok := s.sqlColumnByDomainField[k]
			if !ok {
				return Result{}, fmt.Errorf("%w: %s", ErrInvalidFieldName, k)
			}
			requiredCols[requiredSQLColumn] = struct{}{}
		}
		s.requiredColumns = requiredCols
	}

	builder := strings.Builder{}
	builder.WriteString("SELECT ")

	if len(s.requiredColumns) == 0 {
		builder.WriteString(strings.Join(s.columns, ", "))
	} else {
		selectedCols := make([]string, 0, len(s.requiredColumns))
		for _, col := range s.columns {
			if _, ok := s.requiredColumns[col]; ok {
				selectedCols = append(selectedCols, col)
			}
		}

		if len(selectedCols) == 0 {
			builder.WriteString(strings.Join(s.columns, ", "))
		} else {
			builder.WriteString(strings.Join(selectedCols, ", "))
		}
	}

	builder.WriteString(" FROM ")
	builder.WriteString(s.table)

	for _, join := range s.joins {
		builder.WriteString(" ")
		builder.WriteString(string(join.Type))
		builder.WriteString(" ")
		builder.WriteString(join.Table)
		builder.WriteString(" ON ")
		builder.WriteString(join.Condition)
	}

	args := []any{}
	if len(s.filters) > 0 {
		whereResult, err := WhereSafe(0, s.sqlColumnByDomainField, s.filters...)
		if err != nil {
			return Result{}, err
		}
		args = append(args, whereResult.Args...)
		builder.WriteString(whereResult.SQL)
	}

	if len(s.groups) > 0 {
		groupSQL, err := BuildGroupBy(s.groups, s.sqlColumnByDomainField)
		if err != nil {
			return Result{}, err
		}
		builder.WriteString(groupSQL)
	}

	if len(s.sorts) > 0 {
		builder.WriteString(BuildOrderBy(s.sorts, s.sqlColumnByDomainField))
	}

	builder.WriteString(BuildPagination(s.pagination))

	return Result{
		SQL:  builder.String(),
		Args: args,
	}, nil
}

// BuildOrderBy builds the ORDER BY clause.
func BuildOrderBy(sorts dafi.Sorts, sqlColumnByDomainField map[string]string) string {
	if sorts.IsZero() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(" ORDER BY ")
	for i, sort := range sorts {
		fieldName := string(sort.Field)

		if len(sqlColumnByDomainField) > 0 {
			if sqlColumn, ok := sqlColumnByDomainField[fieldName]; ok {
				fieldName = sqlColumn
			}
		}

		builder.WriteString(fieldName)

		if sort.Type != dafi.None {
			builder.WriteString(" ")
			builder.WriteString(strings.ToUpper(string(sort.Type)))
		}

		if i < len(sorts)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

// BuildPagination builds the LIMIT and OFFSET clauses.
func BuildPagination(pagination dafi.Pagination) string {
	if pagination.HasPageSize() && !pagination.HasPageNumber() {
		pagination.PageNumber = 1
	}

	if pagination.IsZero() {
		return ""
	}

	const maxInt = int(^uint(0) >> 1)

	builder := strings.Builder{}
	builder.WriteString(" LIMIT ")
	if pagination.PageSize > uint(maxInt) {
		builder.WriteString(strconv.Itoa(maxInt))
	} else {
		builder.WriteString(strconv.Itoa(int(pagination.PageSize)))
	}

	if pagination.HasPageNumber() {
		builder.WriteString(" OFFSET ")
		var offset uint
		if pagination.PageNumber > 0 {
			offset = pagination.PageSize * (pagination.PageNumber - 1)
		}
		if offset > uint(maxInt) {
			builder.WriteString(strconv.Itoa(maxInt))
		} else {
			builder.WriteString(strconv.Itoa(int(offset)))
		}
	}

	return builder.String()
}

// BuildGroupBy builds the GROUP BY clause, mapping domain fields to SQL columns when provided.
func BuildGroupBy(groups []string, sqlColumnByDomainField map[string]string) (string, error) {
	if len(sqlColumnByDomainField) > 0 {
		for i, group := range groups {
			sqlColumnName, ok := sqlColumnByDomainField[group]
			if !ok {
				return "", fmt.Errorf("%w: %s", ErrInvalidFieldName, group)
			}
			groups[i] = sqlColumnName
		}
	}

	return " GROUP BY " + strings.Join(groups, ", "), nil
}
