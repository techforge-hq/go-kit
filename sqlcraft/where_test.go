package sqlcraft

import (
	"testing"

	"github.com/techforge-hq/go-kit/dafi"

	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	tests := []struct {
		name    string
		filters dafi.Filters
		want    Result
		wantErr bool
	}{
		{
			name: "one filter",
			filters: dafi.Filters{
				{Field: "email", Operator: dafi.Equal, Value: "hernan_rm@outlook.es"},
			},
			want:    Result{SQL: " WHERE email = $1", Args: []any{"hernan_rm@outlook.es"}},
			wantErr: false,
		},
		{
			name: "and chaining key",
			filters: dafi.Filters{
				{Field: "email", Operator: dafi.Equal, Value: "hernan_rm@outlook.es"},
				{Field: "nickname", Operator: dafi.Equal, Value: "hernanreyes"},
			},
			want:    Result{SQL: " WHERE email = $1 AND nickname = $2", Args: []any{"hernan_rm@outlook.es", "hernanreyes"}},
			wantErr: false,
		},
		{
			name: "or chaining key",
			filters: dafi.Filters{
				{Field: "email", Operator: dafi.Equal, Value: "hernan_rm@outlook.es", ChainingKey: dafi.Or},
				{Field: "nickname", Operator: dafi.Equal, Value: "hernanreyes"},
			},
			want:    Result{SQL: " WHERE email = $1 OR nickname = $2", Args: []any{"hernan_rm@outlook.es", "hernanreyes"}},
			wantErr: false,
		},
		{
			name: "one condition group",
			filters: dafi.Filters{
				{IsGroupOpen: true, Field: "email", Operator: dafi.Equal, Value: "hernan_rm@outlook.es", ChainingKey: dafi.Or},
				{Field: "nickname", Operator: dafi.Equal, Value: "hernanreyes", IsGroupClose: true},
			},
			want:    Result{SQL: " WHERE (email = $1 OR nickname = $2)", Args: []any{"hernan_rm@outlook.es", "hernanreyes"}},
			wantErr: false,
		},
		{
			name: "two conditions group",
			filters: dafi.Filters{
				{IsGroupOpen: true, Field: "email", Operator: dafi.Equal, Value: "hernan_rm@outlook.es", ChainingKey: dafi.Or},
				{Field: "nickname", Operator: dafi.Equal, Value: "hernanreyes", IsGroupClose: true},
				{IsGroupOpen: true, Field: "phone_number", Operator: dafi.Equal, Value: "12345679", ChainingKey: dafi.Or},
				{Field: "full_name", Operator: dafi.Contains, Value: "Hernan Reyes", IsGroupClose: true},
			},
			want:    Result{SQL: " WHERE (email = $1 OR nickname = $2) AND (phone_number = $3 OR full_name ILIKE $4)", Args: []any{"hernan_rm@outlook.es", "hernanreyes", "12345679", "%Hernan Reyes%"}},
			wantErr: false,
		},
		{
			name: "two conditions group with multiple opening and closing parenthesis",
			filters: dafi.Filters{
				{IsGroupOpen: true, GroupOpenQty: 2, Field: "email", Operator: dafi.Equal, Value: "hernan_rm@outlook.es", ChainingKey: dafi.Or},
				{Field: "nickname", Operator: dafi.Equal, Value: "hernanreyes", IsGroupClose: true, GroupOpenQty: 1},
				{IsGroupOpen: true, GroupOpenQty: 1, Field: "phone_number", Operator: dafi.Equal, Value: "12345679", ChainingKey: dafi.Or},
				{Field: "full_name", Operator: dafi.Contains, Value: "Hernan Reyes", IsGroupClose: true, GroupCloseQty: 2},
			},
			want:    Result{SQL: " WHERE ((email = $1 OR nickname = $2) AND (phone_number = $3 OR full_name ILIKE $4))", Args: []any{"hernan_rm@outlook.es", "hernanreyes", "12345679", "%Hernan Reyes%"}},
			wantErr: false,
		},
		{
			name: "in operator",
			filters: dafi.Filters{
				{Field: "id", Operator: dafi.In, Value: []uint{1, 2, 3}},
			},
			want:    Result{SQL: " WHERE id IN ($1, $2, $3)", Args: []any{uint(1), uint(2), uint(3)}},
			wantErr: false,
		},
		{
			name: "not in operator",
			filters: dafi.Filters{
				{Field: "id", Operator: dafi.NotIn, Value: []uint{1, 2, 3}},
			},
			want:    Result{SQL: " WHERE id NOT IN ($1, $2, $3)", Args: []any{uint(1), uint(2), uint(3)}},
			wantErr: false,
		},
		{
			name: "in operator with float",
			filters: dafi.Filters{
				{Field: "price", Operator: dafi.In, Value: []float64{1.1, 2.2, 3.3}},
			},
			want:    Result{SQL: " WHERE price IN ($1, $2, $3)", Args: []any{1.1, 2.2, 3.3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Where(0, tt.filters...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
