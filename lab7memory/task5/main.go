package main

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

func readFile(i int) []byte {
	data := make([]byte, 1<<20) // 1 MB
	copy(data, fmt.Sprintf("HEADER-%04d", i))
	return data
}

func keepHeadersSlice(n int) [][]byte {
	var headers [][]byte
	for i := 0; i < n; i++ {
		data := readFile(i)
		headers = append(headers, data[:11]) // УТЕЧКА!
	}
	return headers
}

func keepHeadersCopy(n int) [][]byte {
	var headers [][]byte
	for i := 0; i < n; i++ {
		data := readFile(i)
		headers = append(headers, bytes.Clone(data[:11])) // ПРАВИЛЬНО
	}
	return headers
}

func TestLeaks(t *testing.T) {
	var m runtime.MemStats

	// Тест утечки
	headers := keepHeadersSlice(100)
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Printf("Утечка (Slice) heap: %v MB\n", m.Alloc/1024/1024)
	runtime.KeepAlive(headers)

	// Тест без утечки
	headersClean := keepHeadersCopy(100)
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Printf("Без утечки (Clone) heap: %v MB\n", m.Alloc/1024/1024)
	runtime.KeepAlive(headersClean)
}
