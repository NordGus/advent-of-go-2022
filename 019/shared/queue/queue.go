package queue

import (
	"errors"
	"math"
)

const (
	initialItemCapacity = 50
)

var (
	IsFullErr  = errors.New("queue: is full")
	IsEmptyErr = errors.New("queue: is empty")
)

// Queue is an implementation of a priority queue.
type Queue[T any] struct {
	items []node[T]
	count uint
}

// node is the container that stores each item (value) in the Queue with its respective priority.
type node[T any] struct {
	value    T
	priority int
}

// New initializes a Queue for the requested type
func New[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]node[T], 0, initialItemCapacity),
	}
}

// Enqueue adds a new item to the Queue with the given priority. It returns an error if the Queue is full.
func (q *Queue[T]) Enqueue(item T, priority int) error {
	if q.count == math.MaxUint {
		return IsFullErr
	}

	nd := node[T]{value: item, priority: priority}

	q.items = append(q.items, nd)
	q.count++

	q.bubbleUp(q.count - 1)

	return nil
}

// Pop removes the top priority item from the Queue and returns it. Additionally, it returns an error if the Queue is
// empty.
func (q *Queue[T]) Pop() (T, error) {
	var out T

	if q.count == 0 {
		return out, IsEmptyErr
	}

	out = q.items[0].value

	q.items = q.items[1:]
	q.count--

	q.bubbleDown(0)

	return out, nil
}

// Peek returns the top priority item in the Queue. Additionally, it returns an error if the Queue is empty.
func (q *Queue[T]) Peek() (T, error) {
	if q.count == 0 {
		var out T

		return out, IsEmptyErr
	}

	return q.items[0].value, nil
}

// bubbleUp moves the item in the given index up the heap
func (q *Queue[T]) bubbleUp(index uint) {
	var (
		current = index
		parent  = q.parent(index)
	)

	for current > 0 && q.items[current].priority > q.items[parent].priority {
		var tmp = q.items[current]

		q.items[current] = q.items[parent]
		q.items[parent] = tmp

		current = parent
		parent = q.parent(current)
	}
}

// bubbleDown moves the item in the given index down the heap
func (q *Queue[T]) bubbleDown(index uint) {
	var current = index

	for current < q.count && !q.isValidParent(current) {
		var (
			child = q.largestChild(current)
			tmp   = q.items[current]
		)

		q.items[current] = q.items[child]
		q.items[child] = tmp

		current = child
	}
}

// parent is a helper function to calculate the given index parent index in the heap array
func (q *Queue[T]) parent(index uint) uint {
	if index == 0 {
		return 0
	}

	return (index - 1) / 2
}

// isValidParent is a helper function that validates that the item in the given index is a valid parent element in the heap.
func (q *Queue[T]) isValidParent(parent uint) bool {
	var (
		left  = parent*2 + 1
		right = parent*2 + 2

		valid bool
	)

	if left >= q.count {
		// The parent node doesn't have a left child. Which means it doesn't have any children,
		// making it a valid parent by default.

		return true
	}

	valid = q.items[parent].priority >= q.items[left].priority

	if right < q.count {
		// The parent node has a right child. So we need to take it into consideration.

		valid = valid && q.items[parent].priority >= q.items[right].priority
	}

	return valid
}

// largestChild is a helper function that returns the index of the top priority child of the item in the given index.
func (q *Queue[T]) largestChild(parent uint) uint {
	var (
		left  = parent*2 + 1
		right = parent*2 + 2
	)

	if left >= q.count {
		// The parent node doesn't have a left child. Which means it doesn't have any children,
		// so we return itself to prevent errors.

		return parent
	}

	if right >= q.count {
		// The parent node doesn't have a right child. Which means the left child is the largest child by default.

		return left
	}

	if q.items[left].priority > q.items[right].priority {
		return left
	}

	return right
}
