package sqlcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert_ToSQL(t *testing.T) {
	tests := []struct {
		name    string
		query   InsertQuery
		want    Result
		wantErr bool
	}{
		{
			name:  "standard insert",
			query: InsertInto("users").WithColumns("first_name", "last_name", "email", "password").WithValues("Hernan", nil, "hernan_rm@outlook.es", "secrethash"),
			want:  Result{SQL: "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)", Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"}},
		},
		{
			name:  "standard insert with returning",
			query: InsertInto("users").WithColumns("first_name", "last_name", "email", "password").WithValues("Hernan", nil, "hernan_rm@outlook.es", "secrethash").Returning("id", "created_at"),
			want:  Result{SQL: "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at", Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"}},
		},
		{
			name:  "insert with multiple row values and returning",
			query: InsertInto("users").WithColumns("first_name", "last_name", "email", "password").WithValues("Hernan", nil, "hernan_rm@outlook.es", "secrethash").WithValues("Brownie", nil, "brownie@gmail.com", "secrethash").Returning("id", "created_at"),
			want:  Result{SQL: "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8) RETURNING id, created_at", Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash", "Brownie", nil, "brownie@gmail.com", "secrethash"}},
		},
		{
			name:    "error empty values",
			query:   InsertInto("users").WithColumns("first_name", "last_name", "email", "password").Returning("id", "created_at"),
			want:    Result{},
			wantErr: true,
		},
		{
			name:    "error mismatch values",
			query:   InsertInto("users").WithColumns("first_name", "last_name", "email", "password").WithValues("hernan").Returning("id", "created_at"),
			want:    Result{},
			wantErr: true,
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
