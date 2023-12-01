package util

// FromPointer returns the value of a pointer or the default value of the type if the pointer is nil.
func FromPointer[T any](v *T) T {
	if v == nil {
		var defaultValue T
		return defaultValue
	}
	return *v
}
