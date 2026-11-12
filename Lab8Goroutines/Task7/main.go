package main

import (
	"fmt"
	"sync"
	"testing"
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

func CountSeq(n int) int {
	return countPrimes(0, n)
}

func CountPar(n, parts int) int {
	var wg sync.WaitGroup
	results := make([]int, parts)
	chunk := n / parts

	for i := 0; i < parts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = countPrimes(i*chunk, (i+1)*chunk)
		}(i)
	}
	wg.Wait()

	total := 0
	for _, v := range results {
		total += v
	}
	return total
}

func main() {
	const N = 200_000

	fmt.Println("Проверка: CountSeq == CountPar?")
	fmt.Printf("Seq: %d | Par(8): %d\n\n", CountSeq(N), CountPar(N, 8))

	fmt.Println("Запуск бенчмарков...")

	partsToTest := []int{2, 8, 1000, N}
	for _, parts := range partsToTest {
		res := testing.Benchmark(func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CountPar(N, parts)
			}
		})
		fmt.Printf("Горутин (parts): %-6d | Время: %v\n", parts, res.String())
	}
}
