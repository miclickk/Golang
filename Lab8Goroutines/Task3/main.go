package main

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func emptyFunc() {}

func main() {
	fmt.Println("--- 1. Потребление памяти ---")
	const N = 100_000
	stop := make(chan struct{})

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-stop
		}()
	}

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	fmt.Printf("Запущено горутин: %d\n", runtime.NumGoroutine())
	stackUsed := memAfter.StackSys - memBefore.StackSys
	fmt.Printf("Память под стек: %d MB (%.2f KB на горутину)\n", stackUsed/1024/1024, float64(stackUsed)/1024/float64(N))

	close(stop)
	wg.Wait()

	fmt.Println("\n--- 2. Бенчмарки ---")
	benchGo := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var w sync.WaitGroup
			w.Add(1)
			go func() { w.Done() }()
			w.Wait()
		}
	})
	fmt.Printf("Создание горутины: %v\n", benchGo)

	benchFunc := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			emptyFunc()
		}
	})
	fmt.Printf("Обычный вызов функции: %v\n", benchFunc)
}
