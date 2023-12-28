package parser

import "unicode"

type TreeType[T any] struct {
	Value  T
	parent *TreeType[T]
}

func (s *TreeType[T]) isType(t *TreeType[T]) bool {
	if s != nil && t == nil {
		return false
	}
	for s != t && s != nil {
		s = s.parent
	}
	if s == t {
		return true
	}
	return false
}

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

func ReverseSliceInPlace[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return
}

func Do(a ...any) {
	return
}

type Iterator[T any] interface {
	Start() int
	End() int
	Legal(int) bool
	Advance(int) int
	Get(i int) T
}

type ArrayRIter[T any] struct {
	arr []T
}

func (s ArrayRIter[T]) Get(i int) T {
	return s.arr[i]
}

func (s ArrayRIter[T]) Start() int {
	return len(s.arr) - 1
}

func (s ArrayRIter[T]) End() int {
	return -1
}

func (s ArrayRIter[T]) Legal(i int) bool {
	return i >= 0
}

func (s ArrayRIter[T]) Advance(i int) int {
	return (i - 1)
}

type ArrayIter[T any] struct {
	arr []T
}

func (s ArrayIter[T]) Start() int {
	return 0
}

func (s ArrayIter[T]) End() int {
	return len(s.arr)
}

func (s ArrayIter[T]) Legal(i int) bool {
	return i < len(s.arr)
}

func (s ArrayIter[T]) Advance(i int) int {
	return (i + 1)
}

func (s ArrayIter[T]) Get(i int) T {
	return s.arr[i]
}

//mapping function
func Mapping[T any](i int) int {
	return 0
}
