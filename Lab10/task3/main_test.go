package main

import (
	"sync"
	"testing"
)

func TestCounterConcurrent(t *testing.T) {
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
	if got := c.Total(); got != 1000 {
		t.Errorf("Total() = %d, want 1000", got)
	}
}
