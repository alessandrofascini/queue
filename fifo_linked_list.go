package queue

// var ErrQueueEmpty  = errors.New("queue is empty")

type element[T any] struct {
	Value T
	next *element[T]
}

func newElement[T any](value T) *element[T] {
	return &element[T]{value, nil}
}

type FIFOLinkedList[T any] struct {
	head, tail *element[T]
	length int
}

func NewFIFOLinkedList[T any]() *FIFOLinkedList[T] {
	return &FIFOLinkedList[T]{nil, nil, 0}
}

func (q *FIFOLinkedList[T]) Enqueue(value T) error {
	e := newElement(value)
	if q.head == nil {
		// if q.head = nil → q.tail = nil
		q.head, q.tail = e, e
	} else {
		q.tail.next, q.tail = e, e
	}
	q.length++
	return nil
}

func (q *FIFOLinkedList[T]) Dequeue() (T, bool) {
	var v T
	if q.head == nil {
		return v, false 
	}
	v = q.head.Value
	hasNext := q.head.next != nil
	if !hasNext {
		// q.head = q.tail
		q.tail = nil
	}
	q.head = q.head.next
	q.length--
	return v, hasNext
}

func (q *FIFOLinkedList[T]) Len() int {
	return q.length
}

func (q *FIFOLinkedList[T]) IsEmpty() bool {
	return q.length == 0
}

func (q *FIFOLinkedList[T]) Clear() {
	q.head, q.tail, q.length = nil, nil, 0
}
