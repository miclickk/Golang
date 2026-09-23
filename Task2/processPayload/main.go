package main

import (
	"bytes"
	"fmt"
	"strings"
)

func processPayloadOriginal(data []byte) string {
	trimmed := string(bytes.TrimSpace(data))
	clean := strings.ReplaceAll(trimmed, "\r", "")
	return clean
}

func processPayloadOptimized(data []byte) []byte {
	trimmed := bytes.TrimSpace(data)
	clean := bytes.ReplaceAll(trimmed, []byte("\r"), []byte(""))
	return clean
}

func main() {
	raw := []byte("  hello \r world  ")
	fmt.Printf("Original:  %q\n", processPayloadOriginal(raw))
	fmt.Printf("Optimized: %q\n", processPayloadOptimized(raw))
}
