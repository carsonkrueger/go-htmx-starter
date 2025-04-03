package tools

func Ptr[T any](v T) *T {
	return &v
}

func PtrSlice[T any](v []T) []*T {
	result := make([]*T, len(v))
	for i, item := range v {
		result[i] = &item
	}
	return result
}
