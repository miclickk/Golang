package main

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

func runScenario(name string, count int, task func()) {
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task()
		}()
	}
	wg.Wait()
	fmt.Printf("%-30s | Создано потоков ОС: %d\n", name, getThreadCount())
}

func getThreadCount() int {
	p := pprof.Lookup("threadcreate")
	return p.Count()
}

func main() {
	runtime.GOMAXPROCS(4)
	fmt.Printf("Стартовое количество потоков: %d\n", getThreadCount())

	runScenario("100 x time.Sleep", 100, func() {
		time.Sleep(200 * time.Millisecond)
	})

	runScenario("100 x blocking OS thread", 100, func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		time.Sleep(200 * time.Millisecond)
	})
}
