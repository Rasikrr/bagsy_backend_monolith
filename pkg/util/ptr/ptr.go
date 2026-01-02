package ptr

func Pointer[T any](value T) *T {
	return &value
}
