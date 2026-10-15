package main

import (
	"cmp"
	"fmt"
)

type Number interface {
	~int | ~int64 | ~float64
}

type Score int

func Sum[T Number](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}

func Average[T Number](nums []T) (float64, bool) {
	if len(nums) == 0 {
		return 0, false
	}
	return float64(Sum(nums)) / float64(len(nums)), true
}

func MaxOf[T cmp.Ordered](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	max := items[0]
	for _, v := range items[1:] {
		if v > max {
			max = v
		}
	}
	return max, true
}

func IndexOf[T comparable](items []T, target T) int {
	for i, v := range items {
		if v == target {
			return i
		}
	}
	return -1
}

func Filter[T any](items []T, keep func(T) bool) []T {
	var res []T
	for _, v := range items {
		if keep(v) {
			res = append(res, v)
		}
	}
	return res
}

func main() {
	scores := []Score{80, 95, 70, 100}

	fmt.Println("Sum:", Sum(scores))

	if avg, ok := Average(scores); ok {
		fmt.Println("Average:", avg)
	}

	highScores := Filter(scores, func(s Score) bool { return s >= 75 })
	fmt.Println("High scores (>= 75):", highScores)
}
