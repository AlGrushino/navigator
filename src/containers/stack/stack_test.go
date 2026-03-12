package stack

import (
	"testing"

	errorspkg "navigator/internal/errors"
)

func TestStack_PushPopTop_Int(t *testing.T) {
	s := Stack[int]()
	if s.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", s.Len())
	}
	s.Push(1)
	s.Push(2)
	if s.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", s.Len())
	}
	if top, _ := s.Top(); top != 2 {
		t.Fatalf("Top() = %d, want 2", top)
	}
	if v, _ := s.Pop(); v != 2 {
		t.Fatalf("Pop() = %d, want 2", v)
	}
	if v, _ := s.Pop(); v != 1 {
		t.Fatalf("Pop() = %d, want 1", v)
	}
	if _, err := s.Pop(); err != errorspkg.ErrEmptyContainer {
		t.Fatalf("expected ErrEmptyContainer, got %v", err)
	}
}

func TestStack_Generic_String(t *testing.T) {
	s := Stack[string]()
	s.Push("a")
	s.Push("b")
	if v, _ := s.Top(); v != "b" {
		t.Fatalf("Top() = %q, want %q", v, "b")
	}
	if v, _ := s.Pop(); v != "b" {
		t.Fatalf("Pop() = %q, want %q", v, "b")
	}
	if v, _ := s.Pop(); v != "a" {
		t.Fatalf("Pop() = %q, want %q", v, "a")
	}
}
