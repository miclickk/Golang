package main

import (
	"fmt"
	"time"
)

func run(name string, ch chan int) {
	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }
	done := make(chan struct{})

	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
			fmt.Printf("%s: sent %d at %v\n", name, i, since())
		}
		close(done)
	}()

	time.Sleep(300 * time.Millisecond)
	for i := 0; i < 3; i++ {
		<-ch
		time.Sleep(100 * time.Millisecond)
	}
	<-done
}

func part1() {
	caps := []int{0, 1, 2, 3, 10}
	for _, c := range caps {
		name := fmt.Sprintf("cap-%d", c)
		run(name, make(chan int, c))
		fmt.Println("--------------------")
	}
}

func part2() {
	ch2 := make(chan string, 2)
	ch2 <- "a"
	ch2 <- "b"
	fmt.Println("Buffer 2:", <-ch2, <-ch2)

	done := make(chan struct{})
	go func() {
		ch1 := make(chan string, 1)
		ch1 <- "a"
		ch1 <- "b"
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		fmt.Println("Buffer 1: deadlocks on second send because buffer is full")
	}
}

func main() {
	fmt.Println("--- Task 2.1: Run with different capacities ---")
	part1()

	fmt.Println("\n--- Task 2.2: Buffer 2 vs Buffer 1 ---")
	part2()
}
