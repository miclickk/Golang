package main

import "fmt"

func main() {

	numbers := [5]float64{1.2, 2.5, 3.7, 4.1, 5.9}
	fmt.Println("массив float64:", numbers)
	values := [5]int{10, 20, 30}
	fmt.Println("массив с 3 значениями, а вывел 5", values)

	var matrix = [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println("2д-массив:")

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Print(matrix[i][j], " ")
		}
		fmt.Println()
	}
}
