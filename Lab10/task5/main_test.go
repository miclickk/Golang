package main

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
		{"Привет", "тевирП"},
	}

	for _, tt := range tests {
		got, err := Reverse(tt.in)
		if err != nil {
			t.Fatalf("Reverse(%q) error: %v", tt.in, err)
		}
		if got != tt.want {
			t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func FuzzReverse(f *testing.F) {
	for _, tc := range []string{"Hello, world", " ", "!12345"} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := Reverse(orig)
		if err1 != nil {
			return
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
