package dafi

type SortType string

const (
	Asc  SortType = "ASC"
	Desc SortType = "DESC"
	None SortType = ""
)

type SortField string

type Sort struct {
	Field SortField
	Type  SortType
}

type Sorts []Sort

func (s Sorts) IsZero() bool {
	return len(s) == 0
}
