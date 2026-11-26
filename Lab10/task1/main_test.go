package main

import (
	"fmt"
	"testing"
)

func TestLetterA(t *testing.T) {
	got, err := Letter(95)
	if err != nil {
		t.Fatalf("Letter(95) returned error: %v", err)
	}
	if got != "A" {
		t.Errorf("Letter(95) = %q, want %q", got, "A")
	}
}

func TestLetter(t *testing.T) {
	tests := []struct {
		name    string
		score   int
		want    string
		wantErr bool
	}{
		{"top", 100, "A", false},
		{"lowest A", 90, "A", false},
		{"highest B", 89, "B", false},
		{"lowest B", 80, "B", false},
		{"C", 75, "C", false},
		{"D", 60, "D", false},
		{"highest F", 59, "F", false},
		{"zero", 0, "F", false},
		{"negative", -1, "", true},
		{"too big", 101, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Letter(tt.score)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Letter(%d) error = %v, wantErr %v", tt.score, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Letter(%d) = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

func checkLetter(t *testing.T, score int, want string) {
	t.Helper()
	got, err := Letter(score)
	if err != nil || got != want {
		t.Errorf("Letter(%d) = %q, %v; want %q", score, got, err, want)
	}
}

func TestCheckLetter(t *testing.T) {
	checkLetter(t, 90, "A")
	checkLetter(t, 80, "B")
}

func TestAverage(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		in   []int
		want float64
	}{
		"empty": {nil, 0},
		"one":   {[]int{90}, 90},
		"three": {[]int{90, 80, 70}, 80},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := Average(tt.in); got != tt.want {
				t.Errorf("Average(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func ExampleLetter() {
	g, _ := Letter(85)
	fmt.Println(g)
	// Output: B
}

func ExampleAverage() {
	fmt.Println(Average([]int{90, 80, 70}))
	// Output: 80
}
