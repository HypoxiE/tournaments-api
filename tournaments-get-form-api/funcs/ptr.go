package funcs

func GetPtr[T any](data T) *T {
	return &data
}

func OrElsePtr[T any](ptr *T, fallback T) *T {
	if ptr != nil {
		return ptr
	}
	return &fallback
}
func OrElseVal[T any](ptr *T, fallback T) T {
	if ptr != nil {
		return *ptr
	}
	return fallback
}
