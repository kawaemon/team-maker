package g

import (
	"encoding/json"
	"math/rand"
)

type Slice[T any] struct {
	slice []T
}

func NewSlice[T any]() Slice[T] {
	return Slice[T]{slice: []T{}}
}
func NewSliceWithCapacity[T any](cap int) Slice[T] {
	return Slice[T]{slice: make([]T, 0, cap)}
}
func NewSliceFrom[T any](slice []T) Slice[T] {
	return Slice[T]{slice}
}

func (s *Slice[T]) Push(value T) {
	s.slice = append(s.slice, value)
}
func (s *Slice[T]) Get(index int) T {
	return s.slice[index]
}
func (s *Slice[T]) GetRef(index int) *T {
	return &s.slice[index]
}

func (s *Slice[T]) Len() int {
	return len(s.slice)
}
func (s *Slice[T]) Cap() int {
	return cap(s.slice)
}
func (s *Slice[T]) Slice() []T {
	return s.slice
}
func (r *Slice[T]) IsNil() bool {
	return r.slice == nil
}

func (s *Slice[T]) Shuffle() {
	rand.Shuffle(s.Len(), func(i, j int) {
		s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
	})
}

func (s *Slice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.slice)
}

func (s *Slice[T]) UnmarshalJSON(b []byte) error {
	var sr []T
	if err := json.Unmarshal(b, &sr); err != nil {
		return err
	}

	s.slice = sr
	return nil
}

func Contains[T comparable](s *Slice[T], target T) bool {
	for _, v := range s.slice {
		if v == target {
			return true
		}
	}
	return false
}
