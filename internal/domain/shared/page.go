package shared

// Page is a generic paginated result.
type Page[T any] struct {
	Items []T
	Total int
}
