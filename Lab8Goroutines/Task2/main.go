package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("--- 1. Правильный WaitGroup ---")
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}
	wg.Wait()

	fmt.Println("\n--- 2. Исправление Race Condition (через atomic) ---")
	var counter int64
	var wgRace sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wgRace.Add(1)
		go func() {
			defer wgRace.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wgRace.Wait()
	fmt.Printf("Counter (должен быть 1000): %d\n", counter)

	fmt.Println("\n--- 3. Перехват паники в горутине ---")
	var wgPanic sync.WaitGroup
	wgPanic.Add(1)
	go func() {
		defer wgPanic.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Паника перехвачена внутри горутины:", r)
			}
		}()
		var m map[string]int
		m["x"] = 1
	}()
	wgPanic.Wait()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Программа жива и успешно завершилась!")
}
