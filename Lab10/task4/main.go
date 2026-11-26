package main

import "fmt"

func Shorten(s string, n int) string {
	if n < 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

func main() {
	fmt.Println(Shorten("hello world", 5))
	fmt.Println(Shorten("Привет мир", 6))
}
