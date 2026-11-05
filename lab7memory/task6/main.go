package main

import (
	"fmt"
)

type Record struct {
	ID   int
	Data [256]byte
}

var live []*Record

func consume(r *Record) int {
	return r.ID
}

func main() {
	for i := 0; i < 80_000; i++ {
		live = append(live, &Record{ID: 1})
	}

	sum := 0
	for i := 0; i < 10_000_000; i++ {
		r := &Record{ID: 1}
		sum += consume(r)
	}
	fmt.Println("Готово. Sum =", sum)
}
