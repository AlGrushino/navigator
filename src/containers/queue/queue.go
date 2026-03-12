package queue

import errorspkg "navigator/internal/errors"

// QueueInterface is a minimal FIFO container interface.
type QueueInterface[T any] interface {
	// Push appends a value to the back of the queue.
	Push(value T)
	// Pop returns the front value and removes it from the queue.
	Pop() (T, error)
	// Front returns the front value without removing it.
	Front() (T, error)
	// Back returns the back value without removing it.
	Back() (T, error)
	// Len returns the number of elements in the queue.
	Len() int
}

type queueImpl[T any] struct {
	buf        []T
	head, tail int
	size       int
}

const initialCapacity = 16

// Queue creates an empty queue.
func Queue[T any]() QueueInterface[T] {
	return &queueImpl[T]{buf: make([]T, initialCapacity)}
}

func (q *queueImpl[T]) Len() int { return q.size }

func (q *queueImpl[T]) Push(value T) {
	if q.size == len(q.buf) {
		q.grow()
	}
	q.buf[q.tail] = value
	q.tail = (q.tail + 1) % len(q.buf)
	q.size++
}

func (q *queueImpl[T]) Pop() (T, error) {
	if q.size == 0 {
		var zero T
		return zero, errorspkg.ErrEmptyContainer
	}
	v := q.buf[q.head]
	q.head = (q.head + 1) % len(q.buf)
	q.size--
	return v, nil
}

func (q *queueImpl[T]) Front() (T, error) {
	if q.size == 0 {
		var zero T
		return zero, errorspkg.ErrEmptyContainer
	}
	return q.buf[q.head], nil
}

func (q *queueImpl[T]) Back() (T, error) {
	if q.size == 0 {
		var zero T
		return zero, errorspkg.ErrEmptyContainer
	}
	idx := q.tail - 1
	if idx < 0 {
		idx = len(q.buf) - 1
	}
	return q.buf[idx], nil
}

func (q *queueImpl[T]) grow() {
	newBuf := make([]T, len(q.buf)*2)
	for i := 0; i < q.size; i++ {
		newBuf[i] = q.buf[(q.head+i)%len(q.buf)]
	}
	q.buf = newBuf
	q.head = 0
	q.tail = q.size
}
