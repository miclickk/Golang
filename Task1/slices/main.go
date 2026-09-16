package main

import "fmt"

func changeFirst(numbers []int) {
	numbers[0] = 100
}

func main() {
	a := []int{1, 2, 3}
	var b []int
	c := make([]int, 3)
	d := a[1:3]
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)
	fmt.Println("d:", d)
	cars := []string{"Lexus", "Toyota", "Ford", "Lixiang"}
	fmt.Println("cars:", cars, "len:", len(cars), "cap:", cap(cars))
	cars = append(cars, "Toyota")
	fmt.Println("after append:", cars, "len:", len(cars), "cap:", cap(cars))
	array := [4]int{10, 20, 30, 40}
	slice := array[1:3]

	slice[0] = 99
	fmt.Println("array:", array)
	fmt.Println("slice:", slice)

	numbers := []int{1, 2, 3}
	changeFirst(numbers)
	fmt.Println("after function:", numbers)
}
