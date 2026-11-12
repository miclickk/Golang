package main

import (
	"fmt"
	"time"
)

type Resource struct{ Name string }

func Open(name string) (*Resource, error) {
	fmt.Println("open", name)
	return &Resource{Name: name}, nil
}
func (r *Resource) Close() { fmt.Println("close", r.Name) }

func processAllBad(names []string) {
	for _, name := range names {
		r, _ := Open(name)
		defer r.Close()
	}
}

func processAllGood(names []string) {
	for _, name := range names {
		func() {
			r, _ := Open(name)
			defer r.Close()
		}()
	}
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
