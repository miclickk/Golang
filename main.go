package main

import "fmt"

type MyError struct{}

func (e *MyError) Error() string { return "boom" }

func check() error {
	var err *MyError = nil
	return err
}

func main() {
	err := check()
	fmt.Println(err == nil)
}
