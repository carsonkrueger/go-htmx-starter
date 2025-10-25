package util

// This function should be used sparingly. Returning a pointer to a stack allocated value will cause heap allocations, slowing down application.
func Ptr[T any](v T) *T {
	return &v
}
