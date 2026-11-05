package main

import "fmt"

func depth(n int) int {
	var pad [1024]byte // Меняем на 1024
	pad[n%len(pad)] = byte(n)
	if n == 0 {
		return 0
	}
	return depth(n-1) + 1 + int(pad[(n+1)%len(pad)])
}

func forever(n int) int {
	return forever(n+1) + 1
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Поймали панику:", r)
		}
	}()
	forever(0)
}
