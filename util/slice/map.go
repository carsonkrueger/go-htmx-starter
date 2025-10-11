package slice

func MapIdx[T, U any](slice []T, fn func(T, int) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v, i)
	}
	return result
}

func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}
