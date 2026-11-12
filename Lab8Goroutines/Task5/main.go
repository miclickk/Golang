package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func slowQuery() int {
	time.Sleep(50 * time.Millisecond)
	return 42
}

func fetchFixed(timeout time.Duration) int {
	ch := make(chan int, 1)
	go func() {
		ch <- slowQuery()
	}()
	select {
	case v := <-ch:
		return v
	case <-time.After(timeout):
		return 0
	}
}

func startTickerFixed(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):

			}
		}
	}()
}

func main() {
	fmt.Printf("Горутин на старте: %d\n", runtime.NumGoroutine())

	for i := 0; i < 100; i++ {
		fetchFixed(10 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("После 100 fetchFixed: %d\n", runtime.NumGoroutine())

	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 10; i++ {
		startTickerFixed(ctx)
	}
	fmt.Printf("После 10 stoppable tickers: %d\n", runtime.NumGoroutine())

	cancel()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("После cancel(): %d\n", runtime.NumGoroutine())
}
