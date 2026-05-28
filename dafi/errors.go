package dafi

import "errors"

var (
	ErrInvalidPaginationValue = errors.New("dafi: invalid pagination value")
	ErrInvalidRelation        = errors.New("dafi: invalid relation")
	ErrInvalidSelectField     = errors.New("dafi: invalid select field")
)
