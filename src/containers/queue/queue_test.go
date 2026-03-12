package queue

import (
	"testing"

	errorspkg "navigator/internal/errors"
)

func TestQueue_PushPopFrontBack_Int(t *testing.T) {
	q := Queue[int]()
	q.Push(1)
	q.Push(2)
	q.Push(3)
	if v, _ := q.Front(); v != 1 {
		t.Fatalf("Front() = %d, want 1", v)
	}
	if v, _ := q.Back(); v != 3 {
		t.Fatalf("Back() = %d, want 3", v)
	}

	if v, _ := q.Pop(); v != 1 {
		t.Fatalf("Pop() = %d, want 1", v)
	}
	if v, _ := q.Pop(); v != 2 {
		t.Fatalf("Pop() = %d, want 2", v)
	}
	if v, _ := q.Pop(); v != 3 {
		t.Fatalf("Pop() = %d, want 3", v)
	}
	if _, err := q.Pop(); err != errorspkg.ErrEmptyContainer {
		t.Fatalf("expected ErrEmptyContainer, got %v", err)
	}
}

func TestQueue_GrowAndOrder(t *testing.T) {
	q := Queue[int]()

	// Push more than initialCapacity to force grow.
	for i := 0; i < initialCapacity+5; i++ {
		q.Push(i)
	}
	for i := 0; i < initialCapacity+5; i++ {
		v, err := q.Pop()
		if err != nil {
			t.Fatalf("Pop() error = %v", err)
		}
		if v != i {
			t.Fatalf("Pop() = %d, want %d", v, i)
		}
	}
}

func TestQueue_Generic_String(t *testing.T) {
	q := Queue[string]()
	q.Push("x")
	q.Push("y")
	if v, _ := q.Front(); v != "x" {
		t.Fatalf("Front() = %q, want %q", v, "x")
	}
	if v, _ := q.Back(); v != "y" {
		t.Fatalf("Back() = %q, want %q", v, "y")
	}
}
