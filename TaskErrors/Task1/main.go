package main

// Task 1.1: Accumulator
func makeAccumulator() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

// Task 1.2: Memoize
func memoize(f func(int) int) func(int) int {
	cache := make(map[int]int)
	return func(n int) int {
		if val, ok := cache[n]; ok {
			return val
		}
		res := f(n)
		cache[n] = res
		return res
	}
}

// Task 2: Fibonacci closure
func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		res := a
		a, b = b, a+b
		return res
	}
}
