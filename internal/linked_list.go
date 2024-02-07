package internal

import (
	"errors"
)

type Node[T any] struct {
	Data     T
	Next     *Node[T]
	Previous *Node[T]
}

type LinkedList[T any] struct {
	FirstNode *Node[T]
	LastNode  *Node[T]
	Size      int
	Actual    *Node[T]
}

func NewLinkedList[T any](elems ...T) LinkedList[T] {
	list := LinkedList[T]{}
	for _, elem := range elems {
		list.Append(elem)
	}
	return list
}

func (l *LinkedList[T]) Append(elem T) {
	if l.FirstNode == nil {
		l.FirstNode = &Node[T]{Data: elem, Next: nil}
		l.LastNode = l.FirstNode
	} else {
		NewNode := Node[T]{Data: elem, Next: nil, Previous: l.LastNode}
		l.LastNode.Next = &NewNode
		l.LastNode = &NewNode
	}
	l.Size++
}

func (l *LinkedList[T]) Pop() (*T, error) {
	if l.LastNode == nil {
		return nil, errors.New("no elements in list")
	}
	elem := &l.LastNode.Data
	if l.LastNode == l.FirstNode {
		l.LastNode = nil
		l.FirstNode = nil
	} else {
		l.LastNode = l.LastNode.Previous
		l.LastNode.Next = nil
	}
	l.Size--
	return elem, nil
}

func (l *LinkedList[T]) ResetActual() error {
	if l.FirstNode == nil {
		return errors.New("no elements in list")
	}
	l.Actual = l.FirstNode
	return nil
}

func (l *LinkedList[T]) GetNext() (*T, error) {
	if l.Actual != nil {
		value := l.Actual.Data
		l.Actual = l.Actual.Next
		return &value, nil
	}
	return nil, errors.New("actual wasn't reseted, or the list has reach the end")
}

func (l *LinkedList[T]) DeleteAtIndex(i int) error {
	if i < 0 || i > l.Size {
		return errors.New("index out of bounds")
	}
	if i == 0 {
		l.FirstNode = l.FirstNode.Next
		l.FirstNode.Previous = nil
		return nil
	}
	node := l.FirstNode.Next
	for j := 1; j < i; j++ {
		node = node.Next
	}
	node.Previous = node.Next
	node = nil
	l.Size--
	return nil
}

// Deletes the first instance of the pointer to the element found in the list.
func (l *LinkedList[T]) DeleteElement(elem *T) error {
	node := l.FirstNode
	if node == nil {
		return errors.New("no elements in list")
	}
	if &node.Data == elem {
		l.FirstNode = l.FirstNode.Next
		l.FirstNode.Previous = nil
		l.Size--
		return nil
	}
	node = node.Next
	for node != nil {
		if &node.Data == elem {
			node.Previous.Next = node.Next
			node = nil
			l.Size--
			return nil
		} else {
			node = node.Next
		}
	}
	return errors.New("element not found in list")
}
