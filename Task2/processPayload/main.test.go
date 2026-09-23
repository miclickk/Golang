package main

import "testing"

var sampleData = []byte("   data payload with some \r extra bytes \r\n  ")

func BenchmarkProcessPayloadOriginal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = processPayloadOriginal(sampleData)
	}
}

func BenchmarkProcessPayloadOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = processPayloadOptimized(sampleData)
	}
}
