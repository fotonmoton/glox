package main

type Stack[Item any] interface {
	Push(Item)
	Pop() Item
	Peek() Item
	At(int) Item
	Size() int
	Empty() bool
}

type node[Item any] struct {
	item Item
	next *node[Item]
}

type stack[OfType any] []OfType

func NewStack[OfType any]() Stack[OfType] {
	return &stack[OfType]{}
}

func (s *stack[Item]) Push(item Item) {
	*s = append(*s, item)
}

func (s *stack[Item]) Pop() Item {
	last := s.Peek()
	*s = (*s)[:len(*s)-1]
	return last
}

func (s *stack[Item]) At(idx int) Item {
	return (*s)[idx]
}

func (s *stack[Item]) Peek() Item {
	return (*s)[len(*s)-1]
}

func (s *stack[_]) Size() int {
	return len(*s)
}

func (s *stack[_]) Empty() bool {
	return s.Size() == 0
}
