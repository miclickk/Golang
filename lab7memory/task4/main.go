package main

import "testing"

type Point struct{ X, Y int }

func NewPointValue(x, y int) Point { return Point{x, y} }

func NewPointPointer(x, y int) *Point { return &Point{x, y} }

func CollectNoCap(n int) []int {
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

func CollectWithCap(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

func BenchmarkNewPointValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPointValue(1, 2)
	}
}

func BenchmarkNewPointPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPointPointer(1, 2)
	}
}

func BenchmarkCollectNoCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CollectNoCap(1000)
	}
}

func BenchmarkCollectWithCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CollectWithCap(1000)
	}
}
