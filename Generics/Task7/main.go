package main

import (
	"fmt"
	"strings"
)

type List[T any] struct {
	val  T
	next *List[T]
}

func (l *List[T]) Push(v T) *List[T] {
	if l == nil {
		return &List[T]{val: v}
	}
	curr := l
	for curr.next != nil {
		curr = curr.next
	}
	curr.next = &List[T]{val: v}
	return l
}

func (l *List[T]) Len() int {
	count := 0
	for curr := l; curr != nil; curr = curr.next {
		count++
	}
	return count
}

func (l *List[T]) String() string {
	if l == nil {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for curr := l; curr != nil; curr = curr.next {
		sb.WriteString(fmt.Sprint(curr.val))
		if curr.next != nil {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func main() {
	var l *List[int]

	fmt.Printf("Initial: %s, len: %d\n", l.String(), l.Len())

	l = l.Push(10)
	l = l.Push(20)
	l = l.Push(30)

	fmt.Printf("After pushes: %s, len: %d\n", l.String(), l.Len())
}
