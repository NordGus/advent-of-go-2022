package structs

import "errors"

type Stack struct {
	top   *node
	count uint
}

type node struct {
	value string
	next  *node
}

func (s *Stack) Peek() string {
	if s.top == nil {
		return ""
	}

	return s.top.value
}

func (s *Stack) Push(value string) {
	s.top = &node{value: value, next: s.top}
	s.count++
}

func (s *Stack) Pop() (string, error) {
	if s.top == nil {
		return "", errors.New("empty stack")
	}

	popped := s.top
	s.top = popped.next
	popped.next = nil
	s.count--

	return popped.value, nil
}

func (s *Stack) Size() uint {
	return s.count
}
