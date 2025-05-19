package internal_type

import (
	"container/list"
)

type Queue[T any] struct {
	q *list.List
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{list.New()}
}

func (q *Queue[T]) Enqueue(value T) {
	q.q.PushBack(value)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.q.Len() == 0 {
		return zero, false
	}

	element := q.q.Front()
	q.q.Remove(element)
	return element.Value.(T), true
}

func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.q.Len() == 0 {
		return zero, false
	}

	return q.q.Front().Value.(T), true
}

func (q *Queue[T]) Len() int {
	return q.q.Len()
}

func (q *Queue[T]) IsEmpty() bool {
	return q.q.Len() == 0
}
