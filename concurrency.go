package queue

import (
	"context"
	"sync"
)

type ConcurrencyInterface[T any] interface {
	Interface[T]
	DequeueBlocking(ctx context.Context) (T, error)
}

type Concurrency[T any] struct {
	Interface[T]
	mu   sync.RWMutex
	cond *sync.Cond
	isClosable bool
}

// Creates a new ConcurrencyQueue
func NewConcurrency[T any](wrappee Interface[T]) *Concurrency[T] {
	q := &Concurrency[T]{
		Interface: wrappee,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Enqueue appends a new element to the tail of the queue.
// Returns error if queue is closed. Signals waiting readers when new data is available.
func (q *Concurrency[T]) Enqueue(value T) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if err := q.Interface.Enqueue(value); err != nil {
		return err
	}
	q.cond.Signal()
	// err = nil
	return nil
}

// Dequeue returns the current value and removes it from the queue.
// Returns hasNext to indicate if more elements remain after dequeuing.
// Non-blocking operation for compatibility.
func (q *Concurrency[T]) Dequeue() (value T, hasNext bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.Dequeue()
}

// DequeueBlocking blocks until a message is available and returns it.
// Context can be used for cancellation and timeout control.
// It avoids spawning helper goroutines by checking ctx.Done() inline.
func (q *Concurrency[T]) DequeueBlocking(ctx context.Context) (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for {
		// If queue has elements, dequeue immediately.
		if q.Interface.Len() > 0 {
			val, _ := q.Interface.Dequeue()
			return val, nil
		}

		// If closed and empty → error.
		// TODO check this
		if closableQueue, ok := q.Interface.(*Closable[T]); ok && closableQueue.IsClosed() {
			var zero T
			return zero, ErrQueueClosed
		}

		// Before waiting, check context cancellation.
		select {
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		default:
			// Nothing cancelled yet; wait for signal.
			q.cond.Wait()
		}
	}
}

// TODO: 
// check with closable
// on close ->  q.cond.Broadcast()!
