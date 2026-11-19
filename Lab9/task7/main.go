package main

import "fmt"

const N = 10000

func passN(size int) int {
	ch := make(chan int, size)
	go func() {
		for i := 0; i < N; i++ {
			ch <- i
		}
		close(ch)
	}()

	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func passBatches() int {
	ch := make(chan []int, 10)
	go func() {
		batchSize := 100
		for i := 0; i < N; i += batchSize {
			batch := make([]int, batchSize)
			for j := 0; j < batchSize; j++ {
				batch[j] = i + j
			}
			ch <- batch
		}
		close(ch)
	}()

	sum := 0
	for batch := range ch {
		for _, v := range batch {
			sum += v
		}
	}
	return sum
}

func main() {
	expected := N * (N - 1) / 2
	fmt.Printf("Expected sum: %d\n", expected)

	sizes := []int{0, 1, 100, 10000}
	for _, s := range sizes {
		res := passN(s)
		fmt.Printf("passN(buf=%-5d): %d (match: %v)\n", s, res, res == expected)
	}

	resBatch := passBatches()
	fmt.Printf("passBatches()    : %d (match: %v)\n", resBatch, resBatch == expected)
}
