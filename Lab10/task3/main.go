package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	hits  map[string]int
	total int
}

func NewCounter() *Counter {
	return &Counter{hits: map[string]int{}}
}

func (c *Counter) Visit(page string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hits[page]++
	c.total++
}

func (c *Counter) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.total
}

func main() {
	c := NewCounter()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Visit("/")
		}()
	}
	wg.Wait()
	fmt.Println(c.Total())
}
