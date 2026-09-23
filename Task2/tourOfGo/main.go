package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	counts := make(map[string]int)
	words := strings.Fields(s)

	for _, word := range words {
		counts[word]++
	}
	return counts
}
func main() {
	input := "I am learning Go and Go is awesome"
	fmt.Println("Результат подсчета слов:")
	fmt.Println(WordCount(input))
}
