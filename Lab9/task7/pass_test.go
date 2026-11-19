package main

import (
	"fmt"
	"testing"
)

func TestPassSum(t *testing.T) {
	expected := N * (N - 1) / 2
	sizes := []int{0, 1, 100, 10000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("buf-%d", size), func(t *testing.T) {
			if got := passN(size); got != expected {
				t.Fatalf("expected %d, got %d", expected, got)
			}
		})
	}

	t.Run("batches", func(t *testing.T) {
		if got := passBatches(); got != expected {
			t.Fatalf("expected %d, got %d", expected, got)
		}
	})
}

func BenchmarkPass(b *testing.B) {
	sizes := []int{0, 1, 100, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("buf_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				passN(size)
			}
		})
	}

	b.Run("batches_100", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			passBatches()
		}
	})
}
