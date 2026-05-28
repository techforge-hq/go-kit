package sqlcraft

import (
	"fmt"
	"reflect"
	"strings"
)

// In builds an IN clause.
func In(value any, initialArgCount int) Result {
	if value == nil {
		return Result{}
	}

	if str, ok := value.(string); ok {
		stringValues := strings.Split(str, ",")
		if len(stringValues) == 0 {
			return Result{}
		}

		args := make([]any, 0, len(stringValues))
		var sql strings.Builder
		sql.WriteString("(")
		for i, v := range stringValues {
			args = append(args, v)
			if i > 0 {
				sql.WriteString(", ")
			}
			fmt.Fprintf(&sql, "$%d", initialArgCount+i)
		}
		sql.WriteString(")")

		return Result{
			SQL:  sql.String(),
			Args: args,
		}
	}

	valSlice := reflect.ValueOf(value)
	if valSlice.Kind() == reflect.Slice {
		if valSlice.Len() == 0 {
			return Result{}
		}

		args := make([]any, 0, valSlice.Len())
		var sql strings.Builder
		sql.WriteString("(")
		for i := 0; i < valSlice.Len(); i++ {
			args = append(args, valSlice.Index(i).Interface())
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(fmt.Sprintf("$%d", initialArgCount+i))
		}
		sql.WriteString(")")

		return Result{
			SQL:  sql.String(),
			Args: args,
		}
	}

	return Result{}
}
