package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func countPrimes(from, to int) int {
	count := 0
	for i := from; i < to; i++ {
		isPrime := true
		for j := 2; j*j <= i; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime && i > 1 {
			count++
		}
	}
	return count
}

func cpuWork(workers int) (int, time.Duration) {
	const N = 600_000
	start := time.Now()
	results := make([]int, workers)
	var wg sync.WaitGroup
	chunk := N / workers

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			results[w] = countPrimes(w*chunk, (w+1)*chunk)
		}(w)
	}
	wg.Wait()

	total := 0
	for _, v := range results {
		total += v
	}
	return total, time.Since(start)
}

func ioWork(workers int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond) // Эмуляция I/O
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	fmt.Printf("Количество ядер (NumCPU): %d\n", runtime.NumCPU())
	fmt.Printf("Текущий GOMAXPROCS: %d\n\n", runtime.GOMAXPROCS(0))

	_, d := cpuWork(1)
	fmt.Printf("CPU work, 1 goroutine: %v\n", d)
	_, d = cpuWork(4)
	fmt.Printf("CPU work, 4 goroutines: %v\n", d)

	fmt.Printf("I/O work, 1000 goroutines: %v\n\n", ioWork(1000))

	runtime.GOMAXPROCS(1)
	fmt.Println("Установлен GOMAXPROCS=1")

	_, d = cpuWork(4)
	fmt.Printf("CPU work, 4 goroutines: %v\n", d)
	fmt.Printf("I/O work, 1000 goroutines: %v\n", ioWork(1000))
}
