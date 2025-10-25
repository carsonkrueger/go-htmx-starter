package slice

func MapIdx[T, R any](slice []T, fn func(T, int) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v, i)
	}
	return result
}

func Map[T, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}
