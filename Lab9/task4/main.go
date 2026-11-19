package main

import (
	"fmt"
	"runtime"
	"time"
)

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{Value: v}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

func NewTree(k int) *Tree {
	var t *Tree
	for i := 1; i <= 10; i++ {
		t = insert(t, i*k)
	}
	return t
}

func walk(t *Tree, ch chan<- int, quit <-chan struct{}) bool {
	if t == nil {
		return true
	}
	if !walk(t.Left, ch, quit) {
		return false
	}
	select {
	case ch <- t.Value:
	case <-quit:
		return false
	}
	return walk(t.Right, ch, quit)
}

func Walk(t *Tree, quit <-chan struct{}) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		walk(t, ch, quit)
	}()
	return ch
}

func Same(t1, t2 *Tree) bool {
	quit := make(chan struct{})
	defer close(quit)

	c1 := Walk(t1, quit)
	c2 := Walk(t2, quit)

	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2

		if ok1 != ok2 || v1 != v2 {
			return false
		}
		if !ok1 {
			break
		}
	}
	return true
}

func main() {
	t1 := NewTree(1)
	t2 := NewTree(1)
	t3 := NewTree(2)

	fmt.Println("Same(t1, t2):", Same(t1, t2))
	fmt.Println("Same(t1, t3):", Same(t1, t3))

	time.Sleep(50 * time.Millisecond)
	base := runtime.NumGoroutine()

	for i := 0; i < 1000; i++ {
		Same(NewTree(1), NewTree(2))
	}

	time.Sleep(50 * time.Millisecond)
	current := runtime.NumGoroutine()
	fmt.Printf("Goroutines before: %d, after 1000 mismatches: %d\n", base, current)
}
