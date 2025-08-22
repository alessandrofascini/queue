package queue

import "errors"

var ErrQueueClosed = errors.New("queue is closed")

type Closable[T any] struct {
	Interface[T]
	closed bool
}

func NewClosableQueue[T any](wrappee Interface[T]) *Closable[T] {
	return &Closable[T]{wrappee, false}
}

func (q *Closable[T]) Enqueue(value T) error {
	if q.closed {
		return ErrQueueClosed
	}
	return q.Interface.Enqueue(value)
}

func (q *Closable[T]) Open() {
	q.closed = false
}

func (q *Closable[T]) Close() {
	q.closed = true
}

func (q *Closable[T]) IsClosed() bool {
	return q.closed
}
