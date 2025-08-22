package queue

type Interface[T any] interface {
	// Enqueue appends a new element to the tail of the queue.
	Enqueue(value T) error
	Dequeue() (value T, hasNext bool)
	Len() int
	Clear()
	IsEmpty() bool
}
