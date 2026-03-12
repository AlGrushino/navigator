package stack

import errorspkg "navigator/internal/errors"

// StackInterface is a minimal LIFO container interface.
type StackInterface[T any] interface {
	// Push adds a value on top of the stack.
	Push(value T)
	// Pop returns the top value and removes it from the stack.
	Pop() (T, error)
	// Top returns the top value without removing it.
	Top() (T, error)
	// Len returns the number of elements in the stack.
	Len() int
}

type stackImpl[T any] struct {
	items []T
}

// Stack creates an empty stack.
func Stack[T any]() StackInterface[T] {
	return &stackImpl[T]{}
}

func (s *stackImpl[T]) Push(value T) {
	s.items = append(s.items, value)
}

func (s *stackImpl[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, errorspkg.ErrEmptyContainer
	}
	idx := len(s.items) - 1
	v := s.items[idx]
	s.items = s.items[:idx]
	return v, nil
}

func (s *stackImpl[T]) Top() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, errorspkg.ErrEmptyContainer
	}
	return s.items[len(s.items)-1], nil
}

func (s *stackImpl[T]) Len() int {
	return len(s.items)
}
