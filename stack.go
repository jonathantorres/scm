package main

import (
	llist "container/list"
)

type Stack struct {
	items *llist.List
}

func newStack() *Stack {
	items := llist.New()

	return &Stack{
		items: items,
	}
}

func (s *Stack) push(value interface{}) {
	s.items.PushFront(value)
}

func (s *Stack) pop() interface{} {
	if s.items.Len() == 0 {
		panic("calling pop() on an empty stack")
	}

	f := s.items.Front()
	return s.items.Remove(f)
}
