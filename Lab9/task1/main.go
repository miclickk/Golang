package main

import (
	"fmt"
	"time"
)

func part1() {
	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }

	ch := make(chan string)

	go func() {
		fmt.Println("sender: about to send", since())
		ch <- "ping"
		fmt.Println("sender: send finished", since())
	}()

	fmt.Println("receiver: waiting...", since())
	time.Sleep(500 * time.Millisecond)
	msg := <-ch
	fmt.Println("receiver: got", msg, since())
}

func testCell(name string, fn func()) {
	done := make(chan struct{})
	var pan interface{}

	go func() {
		defer func() {
			pan = recover()
			close(done)
		}()
		fn()
	}()

	select {
	case <-done:
		if pan != nil {
			fmt.Printf("%-24s: PANIC (%v)\n", name, pan)
		} else {
			fmt.Printf("%-24s: OK\n", name)
		}
	case <-time.After(150 * time.Millisecond):
		fmt.Printf("%-24s: DEADLOCK / BLOCKS FOREVER\n", name)
	}
}

func part2() {
	testCell("1. nil <- send", func() {
		var ch chan int
		ch <- 1
	})
	testCell("2. <- nil receive", func() {
		var ch chan int
		<-ch
	})
	testCell("3. close(nil)", func() {
		var ch chan int
		close(ch)
	})
	testCell("4. len/cap nil", func() {
		var ch chan int
		_ = len(ch)
		_ = cap(ch)
	})

	testCell("5. open <- send (unbuf)", func() {
		ch := make(chan int)
		ch <- 1
	})
	testCell("6. <- open receive (unbuf)", func() {
		ch := make(chan int)
		<-ch
	})
	testCell("7. close(open)", func() {
		ch := make(chan int)
		close(ch)
	})
	testCell("8. len/cap open", func() {
		ch := make(chan int, 3)
		ch <- 1
		_ = len(ch)
		_ = cap(ch)
	})

	testCell("9. closed <- send", func() {
		ch := make(chan int)
		close(ch)
		ch <- 1
	})
	testCell("10. <- closed receive", func() {
		ch := make(chan int)
		close(ch)
		<-ch
	})
	testCell("11. close(closed)", func() {
		ch := make(chan int)
		close(ch)
		close(ch)
	})
	testCell("12. len/cap closed", func() {
		ch := make(chan int, 3)
		ch <- 10
		close(ch)
		_ = len(ch)
		_ = cap(ch)
	})
}

func say(s string, done chan<- struct{}) {
	for i := 0; i < 3; i++ {
		time.Sleep(50 * time.Millisecond)
		fmt.Println(s)
	}
	done <- struct{}{}
}

func part3() {
	done := make(chan struct{})
	go say("goroutine", done)
	<-done
}

func main() {
	fmt.Println("--- Task 1.1: Ping ---")
	part1()

	fmt.Println("\n--- Task 1.2: 12 Table Cells ---")
	part2()

	fmt.Println("\n--- Task 1.3: Say with Channel ---")
	part3()
}
