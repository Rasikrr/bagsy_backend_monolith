package ptr

func Deref[T any](val *T) T {
	if val != nil {
		return *val
	}
	var zero T
	return zero
}

func DerefOr[T any](val *T, fallback T) T {
	if val != nil {
		return *val
	}
	return fallback
}
