package sqlcraft

import (
	"testing"

	"github.com/techforge-hq/go-kit/dafi"

	"github.com/stretchr/testify/assert"
)

func TestDeleteQuery_ToSQL(t *testing.T) {
	tests := []struct {
		name    string
		query   DeleteQuery
		want    Result
		wantErr bool
	}{
		{
			name:  "simple delete",
			query: DeleteFrom("users"),
			want:  Result{SQL: "DELETE FROM users", Args: []any{}},
		},
		{
			name:  "delete with returning",
			query: DeleteFrom("users").Returning("id"),
			want:  Result{SQL: "DELETE FROM users RETURNING id", Args: []any{}},
		},
		{
			name:  "delete with returning and filters",
			query: DeleteFrom("users").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}).Returning("id"),
			want:  Result{SQL: "DELETE FROM users WHERE email = $1 RETURNING id", Args: []any{"hernan_rm@outlook.es"}},
		},
		{
			name:  "delete with returning and filters IN",
			query: DeleteFrom("users").Where(dafi.Filter{Field: "email", Operator: dafi.In, Value: []string{"hernan_rm@outlook.es", "brownie@gmail.com"}}).Returning("id"),
			want:  Result{SQL: "DELETE FROM users WHERE email IN ($1, $2) RETURNING id", Args: []any{"hernan_rm@outlook.es", "brownie@gmail.com"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.ToSQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
