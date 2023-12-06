package parser

import "unicode"

func ReverseSlice[T any](s []T) []T {
	var r []T
	for i := len(s) - 1; i >= 0; i-- {
		r = append(r, s[i])
	}
	return r
}

func IsUpperCase(s string) bool {
	for _, char := range s {
		if !unicode.IsUpper(char) {
			return false
		}
	}
	return true
}
