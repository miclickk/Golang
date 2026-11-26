package main

import "fmt"

func celsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func makeData() []int {
	data := make([]int, 100000)
	for i := range data {
		data[i] = i
	}
	return data
}

func sum(data []int) int {
	total := 0
	for _, v := range data {
		total += v
	}
	return total
}

func main() {
	fmt.Println(celsius(98.6))
	fmt.Println(sum(makeData()))
}
