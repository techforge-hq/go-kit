package dafi

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParser_Parse(t *testing.T) {
	defaultOperators := map[FilterOperator]struct{}{
		"eq":        {},
		"ne":        {},
		"gt":        {},
		"gte":       {},
		"lt":        {},
		"lte":       {},
		"like":      {},
		"in":        {},
		"nin":       {},
		"contains":  {},
		"ncontains": {},
		"is":        {},
		"isn":       {},
	}

	type fields struct {
		operators map[FilterOperator]struct{}
	}
	type args struct {
		values url.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Criteria
		wantErr bool
	}{
		{
			name:    "empty values",
			fields:  fields{operators: defaultOperators},
			args:    args{values: url.Values{}},
			want:    Criteria{},
			wantErr: false,
		},
		{
			name:   "basic pagination",
			fields: fields{operators: defaultOperators},
			args: args{values: url.Values{
				"x": []string{"page:1", "limit:10"},
			}},
			want: Criteria{
				Pagination: Pagination{
					PageNumber: 1,
					PageSize:   10,
				},
			},
			wantErr: false,
		},
		{
			name:   "basic filtering with different operators",
			fields: fields{operators: defaultOperators},
			args: args{values: url.Values{
				"name":   []string{"eq:john"},
				"age":    []string{"gt:18"},
				"status": []string{"in:active,pending"},
				"email":  []string{"contains:example.com"},
			}},
			want: Criteria{
				Filters: Filters{
					{Field: "name", Operator: "eq", Value: "john", ChainingKey: And},
					{Field: "age", Operator: "gt", Value: "18", ChainingKey: And},
					{Field: "status", Operator: "in", Value: []string{"active", "pending"}, ChainingKey: And},
					{Field: "email", Operator: "contains", Value: "example.com", ChainingKey: And},
				},
			},
			wantErr: false,
		},
		{
			name:   "single relation",
			fields: fields{operators: defaultOperators},
			args: args{values: url.Values{
				"relations": []string{"user"},
			}},
			want: Criteria{
				Relations: []string{"user"},
			},
			wantErr: false,
		},
		{
			name:   "multiple relations",
			fields: fields{operators: defaultOperators},
			args: args{values: url.Values{
				"relations": []string{"user,role"},
			}},
			want: Criteria{
				Relations: []string{"user", "role"},
			},
			wantErr: false,
		},
		{
			name:   "relations combined with pagination and filters",
			fields: fields{operators: defaultOperators},
			args: args{values: url.Values{
				"relations": []string{"user,role"},
				"name":      []string{"eq:john"},
				"x":         []string{"page:0", "limit:10"},
			}},
			want: Criteria{
				Relations: []string{"user", "role"},
				Filters: Filters{
					{Field: "name", Operator: "eq", Value: "john", ChainingKey: And},
				},
				Pagination: Pagination{
					PageSize: 10,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &QueryParser{
				operators: tt.fields.operators,
			}

			got, err := p.Parse(tt.args.values)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.ElementsMatch(t, tt.want.Filters, got.Filters)
			assert.ElementsMatch(t, tt.want.Sorts, got.Sorts)
			assert.Equal(t, tt.want.Pagination, got.Pagination)
			assert.Equal(t, tt.want.SelectColumns, got.SelectColumns)
			assert.Equal(t, tt.want.FiltersByModule, got.FiltersByModule)
			assert.Equal(t, tt.want.Relations, got.Relations)
		})
	}
}
