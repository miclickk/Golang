package main

import (
	"fmt"
	"sync"
)

func wrong() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("WRONG finished with panic:", r)
		}
	}()

	results := make(chan int)
	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				results <- w*10 + i
			}
			close(results)
		}(w)
	}

	sum := 0
	for v := range results {
		sum += v
	}
	wg.Wait()
	fmt.Println("sum:", sum)
}

func collect(workers int) <-chan int {
	results := make(chan int)
	var wg sync.WaitGroup

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				results <- w*10 + i
			}
		}(w)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

func sum(in <-chan int) int {
	total := 0
	for n := range in {
		total += n
	}
	return total
}

func main() {
	fmt.Println("--- Task 3.1: WRONG vs COLLECT ---")
	wrong()

	results := collect(3)
	total := 0
	for v := range results {
		total += v
	}
	fmt.Println("collect sum:", total)

	fmt.Println("\n--- Task 3.2: 3-Stage Pipeline ---")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stage3Sum := sum(square(generator(nums...)))
	fmt.Println("3-stage sum:", stage3Sum)

	fmt.Println("\n--- Task 3.3: 4-Stage Pipeline with filterEven ---")
	stage4Sum := sum(filterEven(square(generator(nums...))))
	fmt.Println("4-stage sum:", stage4Sum)
}
