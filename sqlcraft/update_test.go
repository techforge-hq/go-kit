package sqlcraft

import (
	"testing"

	"github.com/techforge-hq/go-kit/dafi"

	"github.com/stretchr/testify/assert"
)

func TestUpdateQuery_ToSQL(t *testing.T) {
	tests := []struct {
		name    string
		query   UpdateQuery
		want    Result
		wantErr bool
	}{
		{
			name:  "update one field",
			query: Update("employees").WithColumns("salary").WithValues(4000),
			want:  Result{SQL: "UPDATE employees SET salary = $1", Args: []any{4000}},
		},
		{
			name:  "update two fields",
			query: Update("employees").WithColumns("salary", "name").WithValues(4000, "Hernan"),
			want:  Result{SQL: "UPDATE employees SET salary = $1, name = $2", Args: []any{4000, "Hernan"}},
		},
		{
			name:  "update two fields with returning",
			query: Update("employees").WithColumns("salary", "name").WithValues(4000, "Hernan").Returning("id"),
			want:  Result{SQL: "UPDATE employees SET salary = $1, name = $2 RETURNING id", Args: []any{4000, "Hernan"}},
		},
		{
			name:  "update two fields with partial update",
			query: Update("employees").WithColumns("salary", "name").WithValues(4000, "Hernan").WithPartialUpdate(),
			want:  Result{SQL: "UPDATE employees SET salary = COALESCE($1, salary), name = COALESCE($2, name)", Args: []any{4000, "Hernan"}},
		},
		{
			name:  "update without provided values",
			query: Update("employees").WithColumns("salary", "name").WithPartialUpdate(),
			want:  Result{SQL: "UPDATE employees SET salary = COALESCE($1, salary), name = COALESCE($2, name)", Args: []any{}},
		},
		{
			name:    "error mismatch values",
			query:   Update("employees").WithColumns("salary", "name").WithValues("salary").WithPartialUpdate(),
			want:    Result{},
			wantErr: true,
		},
		{
			name:  "partial update with filters",
			query: Update("employees").WithColumns("salary", "name").WithValues(4000, "Hernan").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}).WithPartialUpdate(),
			want:  Result{SQL: "UPDATE employees SET salary = COALESCE($1, salary), name = COALESCE($2, name) WHERE email = $3", Args: []any{4000, "Hernan", "hernan_rm@outlook.es"}},
		},
		{
			name:  "partial update with multiple filters including IN",
			query: Update("employees").WithColumns("salary", "name").WithValues(4000, "Hernan").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}, dafi.Filter{Field: "nickname", Operator: dafi.In, Value: []string{"hernan", "brownie"}}).WithPartialUpdate(),
			want:  Result{SQL: "UPDATE employees SET salary = COALESCE($1, salary), name = COALESCE($2, name) WHERE email = $3 AND nickname IN ($4, $5)", Args: []any{4000, "Hernan", "hernan_rm@outlook.es", "hernan", "brownie"}},
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
