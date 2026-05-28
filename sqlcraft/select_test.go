package sqlcraft

import (
	"testing"

	"github.com/techforge-hq/go-kit/dafi"

	"github.com/stretchr/testify/assert"
)

func TestSelectQuery_ToSQL(t *testing.T) {
	tests := []struct {
		name    string
		query   SelectQuery
		want    Result
		wantErr bool
	}{
		{
			name:    "error empty columns",
			query:   Select().From("users"),
			want:    Result{SQL: ""},
			wantErr: true,
		},
		{
			name:  "simple select",
			query: Select("first_name", "last_name").From("users"),
			want:  Result{SQL: "SELECT first_name, last_name FROM users", Args: []any{}},
		},
		{
			name:  "select with required columns",
			query: Select("first_name", "last_name").From("users").RequiredColumns("first_name"),
			want:  Result{SQL: "SELECT first_name FROM users", Args: []any{}},
		},
		{
			name:  "select with filters",
			query: Select("first_name", "last_name").From("users").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}),
			want:  Result{SQL: "SELECT first_name, last_name FROM users WHERE email = $1", Args: []any{"hernan_rm@outlook.es"}},
		},
		{
			name:  "select with filters and order by",
			query: Select("first_name", "last_name").From("users").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}).OrderBy(dafi.Sort{Field: "created_at"}),
			want:  Result{SQL: "SELECT first_name, last_name FROM users WHERE email = $1 ORDER BY created_at", Args: []any{"hernan_rm@outlook.es"}},
		},
		{
			name:  "select with filters and order by desc",
			query: Select("first_name", "last_name").From("users").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}).OrderBy(dafi.Sort{Field: "created_at", Type: dafi.Desc}),
			want:  Result{SQL: "SELECT first_name, last_name FROM users WHERE email = $1 ORDER BY created_at DESC", Args: []any{"hernan_rm@outlook.es"}},
		},
		{
			name:  "select with filters, order by desc and limit",
			query: Select("first_name", "last_name").From("users").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}).OrderBy(dafi.Sort{Field: "created_at", Type: dafi.Desc}).Limit(10),
			want:  Result{SQL: "SELECT first_name, last_name FROM users WHERE email = $1 ORDER BY created_at DESC LIMIT 10 OFFSET 0", Args: []any{"hernan_rm@outlook.es"}},
		},
		{
			name:  "select with filters, order by desc, limit and page",
			query: Select("first_name", "last_name").From("users").Where(dafi.Filter{Field: "email", Value: "hernan_rm@outlook.es"}).OrderBy(dafi.Sort{Field: "created_at", Type: dafi.Desc}).Limit(10).Page(2),
			want:  Result{SQL: "SELECT first_name, last_name FROM users WHERE email = $1 ORDER BY created_at DESC LIMIT 10 OFFSET 10", Args: []any{"hernan_rm@outlook.es"}},
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
