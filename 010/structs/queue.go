package structs

import "errors"

type queue struct {
	head  *node
	tail  *node
	count uint
}

type node struct {
	value instruction
	next  *node
}

func (q *queue) enqueue(item instruction) {
	n := &node{value: item, next: nil}

	if q.tail == nil {
		q.head = n
		q.tail = n
	} else {
		q.tail.next = n
		q.tail = n
	}

	q.count++
}

func (q *queue) dequeue() (instruction, error) {
	if q.head == nil {
		return nil, errors.New("empty queue")
	}

	if q.head.next == nil {
		q.tail = nil
	}

	dequeued := q.head
	q.head = dequeued.next
	dequeued.next = nil
	q.count--

	return dequeued.value, nil
}

func (q *queue) clear() {
	for {
		_, err := q.dequeue()
		if err != nil {
			break
		}
	}
	q.count = 0
}

func (q *queue) size() uint {
	return q.count
}
