package main

import "testing"

var (
	sinkF float64
	sinkB bool
)

func BenchmarkCelsiusBN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		celsius(98.6)
	}
}

func BenchmarkCelsiusLoop(b *testing.B) {
	for b.Loop() {
		celsius(98.6)
	}
}

func BenchmarkCelsiusSink(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkF = celsius(98.6)
	}
}

func BenchmarkCelsiusSinkVar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkF = celsius(float64(i))
	}
}

func BenchmarkSumSetupInside(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := makeData()
		sinkB = sum(data) > 0
	}
}

func BenchmarkSumLoop(b *testing.B) {
	data := makeData()
	for b.Loop() {
		sum(data)
	}
}
