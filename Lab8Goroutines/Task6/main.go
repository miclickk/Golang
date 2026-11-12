package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func buildReport(ctx context.Context) []byte {
	select {
	case <-time.After(20 * time.Millisecond):
		return make([]byte, 64*1024)
	case <-ctx.Done():
		return nil
	}
}

func handle(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan []byte, 1)

	go func() {
		ch <- buildReport(ctx)
	}()

	select {
	case <-ch:

	case <-ctx.Done():

	}
}

func printStats(round string) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s | Goroutines: %d | Heap: %d MB\n", round, runtime.NumGoroutine(), m.HeapAlloc/1024/1024)
}

func main() {
	printStats("Start          ")

	for r := 1; r <= 3; r++ {
		for i := 0; i < 1000; i++ {
			handle(5 * time.Millisecond)
		}
		time.Sleep(50 * time.Millisecond)
		printStats(fmt.Sprintf("После round %d", r))
	}
}
