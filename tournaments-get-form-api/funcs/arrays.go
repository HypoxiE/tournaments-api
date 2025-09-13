package funcs

func GetUnique[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	for _, val := range input {
		if _, ok := seen[val]; !ok {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}
