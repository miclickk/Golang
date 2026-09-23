package main

import (
	"fmt"
	"strings"
)

func concat(values []string) string {
	totalLength := 0
	for _, v := range values {
		totalLength += len(v)
	}

	var sb strings.Builder
	sb.Grow(totalLength)

	for _, v := range values {
		sb.WriteString(v)
	}
	return sb.String()
}

func main() {
	words := []string{"Go", " ", "is", " ", "fast"}
	result := concat(words)
	fmt.Println("Результат:", result)
}
