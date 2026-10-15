package main

import (
	"testing"
)

type Number interface {
	~int | ~int64 | ~float64
}

func SumInts(nums []int) int {
	var s int
	for _, n := range nums {
		s += n
	}
	return s
}

func SumGeneric[T Number](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}

func SumAny(nums []any) int {
	total := 0
	for _, n := range nums {
		total += n.(int)
	}
	return total
}

func ToAny(nums []int) []any {
	res := make([]any, len(nums))
	for i, n := range nums {
		res[i] = n
	}
	return res
}

func generateData() []int {
	nums := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		nums[i] = i
	}
	return nums
}

func BenchmarkSumInts(b *testing.B) {
	nums := generateData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumInts(nums)
	}
}

func BenchmarkSumGeneric(b *testing.B) {
	nums := generateData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumGeneric(nums)
	}
}

func BenchmarkSumAny(b *testing.B) {
	nums := generateData()
	anyNums := ToAny(nums)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumAny(anyNums)
	}
}

func BenchmarkSumAnyWithBoxing(b *testing.B) {
	nums := generateData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumAny(ToAny(nums))
	}
}
