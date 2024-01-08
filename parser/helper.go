package parser

import (
	"encoding/binary"
	"strings"
	"unicode"
)

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

// mapping function
func Mapping[T any](i int) int {
	return 0
}

// FirstLower 字符串首字母小写
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func ReadUint16(ins []byte) uint16 {
	return binary.LittleEndian.Uint16(ins)
}

func ReadInt16(ins []byte) int16 {
	return (int16)(binary.LittleEndian.Uint16(ins))
}

func ReadUint32(ins []byte) uint32 {
	return binary.LittleEndian.Uint32(ins)
}

func ReadInt32(ins []byte) int32 {
	return int32(binary.LittleEndian.Uint32(ins))
}

func ReadUint64(ins []byte) uint64 {
	return binary.LittleEndian.Uint64(ins)
}

func ReadInt64(ins []byte) int64 {
	return int64(binary.LittleEndian.Uint64(ins))
}

func ReadFloat32(ins []byte) float32 {
	return float32(binary.LittleEndian.Uint32(ins))
}

func ReadFloat64(ins []byte) float64 {
	return float64(binary.LittleEndian.Uint64(ins))
}

func WriteUint16(ins []byte, x uint16) {
	binary.LittleEndian.PutUint16(ins, x)
}

func WriteInt16(ins []byte, x int16) {
	binary.LittleEndian.PutUint16(ins, uint16(x))
}

func WriteUint32(ins []byte, x uint32) {
	binary.LittleEndian.PutUint32(ins, x)
}

func WriteInt32(ins []byte, x int32) {
	binary.LittleEndian.PutUint32(ins, uint32(x))
}

func WriteUint64(ins []byte, x uint64) {
	binary.LittleEndian.PutUint64(ins, x)
}

func WriteInt64(ins []byte, x int64) {
	binary.LittleEndian.PutUint64(ins, uint64(x))
}

func WriteFloat32(ins []byte, x float32) {
	binary.LittleEndian.PutUint32(ins, (uint32)(x))
}

func WriteFloat64(ins []byte, x float64) {
	binary.LittleEndian.PutUint64(ins, (uint64)(x))
}
