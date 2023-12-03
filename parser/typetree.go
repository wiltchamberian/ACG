//to represent the hierarchy of type
//we need a new structure to do it
//this code can be made to metaprogramming code but it
//seems that golang doesn't support this.

package parser

type TreeType[T any] struct {
	Value  T
	parent *TreeType[T]
}

func (s *TreeType[T]) isType(t *TreeType[T]) bool {
	if s != nil && t == nil {
		return false
	}
	for s != t {
		s = s.parent
	}
	if s == t {
		return true
	}
	return false
}
