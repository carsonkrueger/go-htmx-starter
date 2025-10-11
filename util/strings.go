package util

import "strings"

func ToUpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func ToLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func ToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}
	var result string
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_" + string(r)
		} else {
			result += string(r)
		}
	}
	return strings.ToLower(result)
}
