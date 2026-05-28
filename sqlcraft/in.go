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
		sql := "("
		for i, v := range stringValues {
			args = append(args, v)
			if i > 0 {
				sql += ", "
			}
			sql += fmt.Sprintf("$%d", initialArgCount+i)
		}
		sql += ")"

		return Result{
			SQL:  sql,
			Args: args,
		}
	}

	valSlice := reflect.ValueOf(value)
	if valSlice.Kind() == reflect.Slice {
		if valSlice.Len() == 0 {
			return Result{}
		}

		args := make([]any, 0, valSlice.Len())
		sql := "("
		for i := 0; i < valSlice.Len(); i++ {
			args = append(args, valSlice.Index(i).Interface())
			if i > 0 {
				sql += ", "
			}
			sql += fmt.Sprintf("$%d", initialArgCount+i)
		}
		sql += ")"

		return Result{
			SQL:  sql,
			Args: args,
		}
	}

	return Result{}
}
