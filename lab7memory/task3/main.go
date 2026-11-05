package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func JoinIDsSlow(ids []int) string {
	s := ""
	for i, id := range ids {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprint(id)
	}
	return s
}

func JoinIDsFast(ids []int) string {
	var sb strings.Builder
	sb.Grow(len(ids) * 8)
	var tmp [20]byte
	for i, id := range ids {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.Write(strconv.AppendInt(tmp[:0], int64(id), 10))
	}
	return sb.String()
}

func TestJoinIDsAllocs(t *testing.T) {
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	allocsSlow := testing.AllocsPerRun(1000, func() {
		JoinIDsSlow(ids)
	})
	fmt.Printf("JoinIDsSlow: %f allocs/op\n", allocsSlow)

	allocsFast := testing.AllocsPerRun(1000, func() {
		JoinIDsFast(ids)
	})
	fmt.Printf("JoinIDsFast: %f allocs/op\n", allocsFast)

	if allocsFast > 1 {
		t.Errorf("JoinIDsFast: %v allocs, want at most 1", allocsFast)
	}
}

func BenchmarkJoinIDsSlow(b *testing.B) {
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		JoinIDsSlow(ids)
	}
}

func BenchmarkJoinIDsFast(b *testing.B) {
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		JoinIDsFast(ids)
	}
}
