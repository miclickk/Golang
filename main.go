package main

import "fmt"

var a = 3
var b = 4

func add(a int, b int) int {
	return a + b
}

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(add(a, b))
}
