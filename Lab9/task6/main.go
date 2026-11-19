package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Word  string
	Upper string
}

func producer(words []string, done <-chan struct{}) <-chan string {
	jobs := make(chan string)
	go func() {
		defer close(jobs)
		for _, w := range words {
			select {
			case jobs <- w:
			case <-done:
				return
			}
		}
	}()
	return jobs
}

func startWorkers(n int, jobs <-chan string, done <-chan struct{}) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case w, ok := <-jobs:
					if !ok {
						return
					}
					res := Result{Word: w, Upper: strings.ToUpper(w)}
					select {
					case results <- res:
					case <-done:
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func startWorkersLeaky(n int, jobs <-chan string) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range jobs {
				results <- Result{Word: w, Upper: strings.ToUpper(w)}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

var ErrQueueFull = errors.New("queue is full")

func Submit(queue chan<- string, job string, timeout time.Duration) error {
	t := time.NewTimer(timeout)
	defer t.Stop()

	select {
	case queue <- job:
		return nil
	case <-t.C:
		return ErrQueueFull
	}
}

func main() {
	words := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}

	fmt.Println("--- 1. Full Run ---")
	done1 := make(chan struct{})
	jobs1 := producer(words, done1)
	results1 := startWorkers(3, jobs1, done1)
	for r := range results1 {
		fmt.Printf("%s -> %s\n", r.Word, r.Upper)
	}

	fmt.Println("\n--- 2. Cancellation after 3 items ---")
	done2 := make(chan struct{})
	jobs2 := producer(words, done2)
	results2 := startWorkers(3, jobs2, done2)

	count := 0
	for r := range results2 {
		fmt.Printf("got: %s -> %s\n", r.Word, r.Upper)
		count++
		if count == 3 {
			close(done2)
			break
		}
	}

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines after cancellation: %d\n", runtime.NumGoroutine())

	fmt.Println("\n--- 3. Leaked Goroutines without done ---")
	jobsLeaky := make(chan string)
	resultsLeaky := startWorkersLeaky(3, jobsLeaky)

	go func() {
		for _, w := range words {
			jobsLeaky <- w
		}
	}()

	count = 0
	for range resultsLeaky {
		count++
		if count == 3 {
			break
		}
	}

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines with leak: %d\n", runtime.NumGoroutine())

	fmt.Println("\n--- 4. Submit (Backpressure) ---")
	q := make(chan string, 1)

	err1 := Submit(q, "job-1", 50*time.Millisecond)
	fmt.Println("Submit 1:", err1)

	err2 := Submit(q, "job-2", 50*time.Millisecond)
	fmt.Println("Submit 2 (full):", err2)
}
