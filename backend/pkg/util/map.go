package util

// Map is a generic function to map a slice of type T to a slice of type R.
func Map[T any, R any](input []T, mapFunc func(item T, index int) R) []R {
	result := make([]R, len(input))

	for i, item := range input {
		result[i] = mapFunc(item, i)
	}

	return result
}
