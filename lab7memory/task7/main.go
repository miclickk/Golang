package main

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

var bufPool = sync.Pool{
	New: func() any { b := new(bytes.Buffer); b.Grow(4096); return b },
}

func RenderNew(ids []int) int {
	b := new(bytes.Buffer)
	b.Grow(4096)
	for _, id := range ids {
		fmt.Fprintf(b, "id=%d\n", id)
	}
	return b.Len()
}

func RenderPool(ids []int) int {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset() // ОЧИСТКА ВАЖНА!
	for _, id := range ids {
		fmt.Fprintf(buf, "id=%d\n", id)
	}
	n := buf.Len()
	bufPool.Put(buf)
	return n
}

func BenchmarkRenderNew(b *testing.B) {
	ids := make([]int, 200)
	for i := 0; i < b.N; i++ {
		RenderNew(ids)
	}
}

func BenchmarkRenderPool(b *testing.B) {
	ids := make([]int, 200)
	for i := 0; i < b.N; i++ {
		RenderPool(ids)
	}
}
