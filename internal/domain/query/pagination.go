package query

type Page[T any] struct {
	Items []T
	Total int
}

func NewPage[T any](items []T, total int) *Page[T] {
	return &Page[T]{
		Items: items,
		Total: total,
	}
}
