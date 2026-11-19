package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func part1() {
	runFair := func() {
		a := make(chan int, 1000)
		b := make(chan int, 1000)
		for i := 0; i < 1000; i++ {
			a <- 1
			b <- 2
		}
		countA, countB := 0, 0
		for i := 0; i < 1000; i++ {
			select {
			case <-a:
				countA++
			case <-b:
				countB++
			}
		}
		fmt.Printf("Fair split: a=%d, b=%d\n", countA, countB)
	}

	for i := 0; i < 3; i++ {
		runFair()
	}

	a := make(chan int, 1000)
	b := make(chan int, 1000)
	for i := 0; i < 1000; i++ {
		a <- 1
		b <- 2
	}
	countA, countB := 0, 0
	for i := 0; i < 1000; i++ {
		select {
		case <-a:
			countA++
		default:
			select {
			case <-a:
				countA++
			case <-b:
				countB++
			}
		}
	}
	fmt.Printf("Biased split (prefer a): a=%d, b=%d\n", countA, countB)
}

func fetchWithTimeout(timeout time.Duration) (string, error) {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "response data"
	}()

	select {
	case res := <-ch:
		return res, nil
	case <-time.After(timeout):
		return "", errors.New("timeout")
	}
}

func part2() {
	res, err := fetchWithTimeout(150 * time.Millisecond)
	fmt.Println("Fetch 150ms:", res, err)

	before := runtime.NumGoroutine()
	res, err = fetchWithTimeout(20 * time.Millisecond)
	fmt.Println("Fetch 20ms:", res, err)

	time.Sleep(150 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Printf("Goroutines before timeout: %d, after: %d\n", before, after)
}

func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

func part3() {
	gen := func(start, count int, delay time.Duration) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 0; i < count; i++ {
				time.Sleep(delay)
				out <- start + i
			}
		}()
		return out
	}

	merged := merge(gen(1, 3, 10*time.Millisecond), gen(100, 3, 20*time.Millisecond))
	for v := range merged {
		fmt.Print(v, " ")
	}
	fmt.Println()

	chA := make(chan int)
	chB := make(chan int)
	close(chA)

	go func() {
		time.Sleep(30 * time.Millisecond)
		chB <- 42
		close(chB)
	}()

	zeroCount := 0
	for bActive := true; bActive; {
		select {
		case v, ok := <-chA:
			if !ok && v == 0 {
				zeroCount++
			}
		case _, ok := <-chB:
			if !ok {
				bActive = false
			}
		}
	}
	fmt.Printf("Received %d useless zero values without 'a = nil'\n", zeroCount)
}

func participation() {
	ball := make(chan int)
	done := make(chan struct{})
	var wg sync.WaitGroup

	player := func(name string) {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case v, ok := <-ball:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case ball <- v + 1:
				}
			}
		}
	}

	wg.Add(2)
	go player("P1")
	go player("P2")

	ball <- 0
	time.Sleep(100 * time.Millisecond)
	hits := <-ball
	close(done)
	wg.Wait()

	time.Sleep(20 * time.Millisecond)
	fmt.Printf("Ping-pong hits in 100ms: %d, goroutines active: %d\n", hits, runtime.NumGoroutine())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("--- Task 5.1: Random vs Biased Select ---")
	part1()

	fmt.Println("\n--- Task 5.2: Fetch With Timeout ---")
	part2()

	fmt.Println("\n--- Task 5.3: Merge and Nil Channels ---")
	part3()

	fmt.Println("\n--- Task 5.4: Ping Pong Game ---")
	participation()
}
