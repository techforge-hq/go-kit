package dafi

type Pagination struct {
	PageNumber uint
	PageSize   uint
}

func (p Pagination) IsZero() bool {
	return p.PageNumber == 0 && p.PageSize == 0
}

func (p Pagination) HasPageNumber() bool {
	return p.PageNumber > 0
}

func (p Pagination) HasPageSize() bool {
	return p.PageSize > 0
}
