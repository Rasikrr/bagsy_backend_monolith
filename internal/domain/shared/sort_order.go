package shared

import "errors"

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

var ErrInvalidSortOrder = errors.New("invalid sort_order value")

func ParseSortOrder(s string) (SortOrder, error) {
	switch SortOrder(s) {
	case SortAsc:
		return SortAsc, nil
	case SortDesc:
		return SortDesc, nil
	default:
		return "", ErrInvalidSortOrder
	}
}

func (s SortOrder) String() string { return string(s) }
