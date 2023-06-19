package structs

import "errors"

type list struct {
	head  *node
	tail  *node
	count uint
}

type node struct {
	value int
	next  *node
}

func (l *list) push(item int) {
	n := &node{value: item, next: nil}

	if l.tail == nil {
		l.head = n
		l.tail = n
	} else {
		l.tail.next = n
		l.tail = n
	}

	l.count++
}

func (l *list) pop() (int, error) {
	if l.head == nil {
		return 0, errors.New("empty queue")
	}

	if l.head.next == nil {
		l.tail = nil
	}

	popped := l.head
	l.head = popped.next
	popped.next = nil
	l.count--

	return popped.value, nil
}

func (l *list) toSlice() []int {
	out := make([]int, l.count)
	current := l.head

	for i := 0; current != nil; i++ {
		out[i] = current.value
		current = current.next
	}

	return out
}
