package sqlcraft

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/techforge-hq/go-kit/dafi"
)

var psqlOperatorByDafiOperator = map[dafi.FilterOperator]string{
	dafi.Equal:          "=",
	dafi.NotEqual:       "<>",
	dafi.Greater:        ">",
	dafi.GreaterOrEqual: ">=",
	dafi.Less:           "<",
	dafi.LessOrEqual:    "<=",
	dafi.Contains:       "ILIKE",
	dafi.NotContains:    "NOT ILIKE",
	dafi.Is:             "IS",
	dafi.IsNull:         "IS NULL",
	dafi.IsNot:          "IS NOT",
	dafi.IsNotNull:      "IS NOT NULL",
	dafi.In:             "IN",
	dafi.NotIn:          "NOT IN",
	dafi.Default:        "",
}

// WhereSafe maps domain field names to SQL column names.
// If a filter with an unknown domain field name is found it will return ErrInvalidFieldName.
func WhereSafe(initialArgCount int, sqlColumnByDomainField map[string]string, filters ...dafi.Filter) (Result, error) {
	if len(sqlColumnByDomainField) > 0 {
		for i, filter := range filters {
			sqlColumnName, ok := sqlColumnByDomainField[string(filter.Field)]
			if !ok {
				return Result{}, fmt.Errorf("%w: %s", ErrInvalidFieldName, filter.Field)
			}
			filters[i].Field = dafi.FilterField(sqlColumnName)
		}
	}

	return Where(initialArgCount, filters...)
}

// Where builds the WHERE clause.
func Where(initialArgCount int, filters ...dafi.Filter) (Result, error) {
	if len(filters) == 0 {
		return Result{}, nil
	}

	builder := strings.Builder{}
	builder.WriteString(" WHERE ")

	args := []any{}
	argCount := initialArgCount

	for i, filter := range filters {
		if filter.IsGroupOpen {
			for j := 0; j < max(1, filter.GroupOpenQty); j++ {
				builder.WriteString("(")
			}
		}

		operator := filter.Operator
		if operator == "" {
			operator = dafi.Equal
		}

		switch operator {
		case dafi.IsNull, dafi.IsNotNull:
			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(psqlOperatorByDafiOperator[operator])
		case dafi.In, dafi.NotIn:
			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(psqlOperatorByDafiOperator[operator])
			builder.WriteString(" ")

			inResult := In(filter.Value, argCount+1)
			builder.WriteString(inResult.SQL)
			args = append(args, inResult.Args...)
			argCount += len(inResult.Args)
		case dafi.Contains, dafi.NotContains:
			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(psqlOperatorByDafiOperator[operator])
			builder.WriteString(" ")
			builder.WriteString("$")
			builder.WriteString(strconv.Itoa(argCount + 1))

			args = append(args, fmt.Sprintf("%%%v%%", filter.Value))
			argCount++
		default:
			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(psqlOperatorByDafiOperator[operator])
			builder.WriteString(" ")
			builder.WriteString("$")
			builder.WriteString(strconv.Itoa(argCount + 1))

			args = append(args, filter.Value)
			argCount++
		}

		if filter.IsGroupClose {
			for j := 0; j < max(1, filter.GroupCloseQty); j++ {
				builder.WriteString(")")
			}
		}

		if i < len(filters)-1 {
			chainingKey := filter.ChainingKey
			if chainingKey == "" {
				chainingKey = dafi.And
			}
			builder.WriteString(" ")
			builder.WriteString(string(chainingKey))
			builder.WriteString(" ")
		}
	}

	return Result{
		SQL:  builder.String(),
		Args: args,
	}, nil
}
