package dafi

type (
	FilterField string
	FilterValue any
)

type FilterOperator string

const (
	Equal          FilterOperator = "eq"
	NotEqual       FilterOperator = "ne"
	Greater        FilterOperator = "gt"
	GreaterOrEqual FilterOperator = "gte"
	Less           FilterOperator = "lt"
	LessOrEqual    FilterOperator = "lte"
	Like           FilterOperator = "like"
	In             FilterOperator = "in"
	NotIn          FilterOperator = "nin"
	Contains       FilterOperator = "contains"
	NotContains    FilterOperator = "ncontains"
	Is             FilterOperator = "is"
	IsNull         FilterOperator = "isnull"
	IsNot          FilterOperator = "isn"
	IsNotNull      FilterOperator = "isnnull"
	Default        FilterOperator = "default"
)

type FilterChainingKey string

const (
	And FilterChainingKey = "AND"
	Or  FilterChainingKey = "OR"
)

type Filter struct {
	Module                            string
	IsGroupOpen                       bool
	GroupOpenQty                      int
	Field                             FilterField
	Operator                          FilterOperator
	Value                             FilterValue
	IsGroupClose                      bool
	GroupCloseQty                     int
	ChainingKey                       FilterChainingKey
	OverridePreviousFilterChainingKey FilterChainingKey
}

type Filters []Filter

func (f Filters) IsZero() bool {
	return len(f) == 0
}

func FilterBy(name string, operator FilterOperator, value any) Filters {
	return Filters{
		{
			Field:    FilterField(name),
			Operator: operator,
			Value:    value,
		},
	}
}

func (f Filters) Or(field string, operator FilterOperator, value any) Filters {
	if f.IsZero() {
		return Filters{{Field: FilterField(field), Operator: operator, Value: value}}
	}

	f[len(f)-1].ChainingKey = Or

	return append(f, Filter{
		Field:    FilterField(field),
		Operator: operator,
		Value:    value,
	})
}

func (f Filters) And(field string, operator FilterOperator, value any) Filters {
	if f.IsZero() {
		return Filters{{Field: FilterField(field), Operator: operator, Value: value}}
	}

	f[len(f)-1].ChainingKey = And

	return append(f, Filter{
		Field:    FilterField(field),
		Operator: operator,
		Value:    value,
	})
}

func (f Filters) AndGroup(filters ...Filter) Filters {
	if len(filters) == 0 {
		return f
	}

	if len(f) > 0 {
		f[len(f)-1].ChainingKey = And
	}

	filters[0].IsGroupOpen = true
	filters[len(filters)-1].IsGroupClose = true

	f = append(f, filters...)

	return f
}

func (f Filters) OrGroup(filters ...Filter) Filters {
	if len(filters) == 0 {
		return f
	}

	if len(f) > 0 {
		f[len(f)-1].ChainingKey = Or
	}

	filters[0].IsGroupOpen = true
	filters[len(filters)-1].IsGroupClose = true

	f = append(f, filters...)

	return f
}
