package structs

import (
	"errors"
)

type Decoder struct {
	head    *node
	tail    *node
	index   uint
	size    uint
	maxSize uint
}

type node struct {
	value rune
	next  *node
}

func NewDecoder(size uint) Decoder {
	return Decoder{
		maxSize: size,
	}
}

func (d *Decoder) Push(value rune) {
	n := &node{value: value, next: nil}

	if d.tail == nil {
		d.head = n
		d.tail = n
	} else {
		d.tail.next = n
		d.tail = n
	}

	d.size++
	d.index++

	if d.size > d.maxSize {
		d.Pop()
	}
}

func (d *Decoder) Pop() (rune, error) {
	if d.head == nil {
		return rune(0), errors.New("empty queue")
	}

	if d.head.next == nil {
		d.tail = nil
	}

	popped := d.head
	d.head = popped.next
	popped.next = nil
	d.size--

	return popped.value, nil
}

func (d *Decoder) IsStartOfPackage() (bool, uint) {
	if d.size < d.maxSize {
		return false, 0
	}

	for n := d.head; n != nil; n = n.next {
		for c := n.next; c != nil; c = c.next {
			if n.value == c.value {
				return false, 0
			}
		}
	}

	return true, d.index
}

func (d *Decoder) Clear() {
	for {
		_, err := d.Pop()
		if err != nil {
			break
		}
	}
	d.size = 0
	d.index = 0
}
