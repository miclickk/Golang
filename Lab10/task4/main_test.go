package main

import (
	"testing"
	"unicode/utf8"
)

func TestShorten(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want string
	}{
		{"fits", "hello", 10, "hello"},
		{"cut", "hello world", 5, "hello..."},
		{"empty", "", 3, ""},
		{"unicode", "Привет", 2, "Пр..."},
		{"zero", "hello", 0, "..."},
		{"negative", "hello", -5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Shorten(tt.s, tt.n); got != tt.want {
				t.Errorf("Shorten(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

func FuzzShorten(f *testing.F) {
	f.Add("hello world", 5)
	f.Add("Go", 10)
	f.Add("Привет", 2)
	f.Add("abc", -1)

	f.Fuzz(func(t *testing.T, s string, n int) {
		if n > 1000 {
			t.Skip()
		}
		got := Shorten(s, n)
		if n < 0 {
			if got != "" {
				t.Errorf("Shorten(%q, %d) = %q, want empty", s, n, got)
			}
			return
		}
		if utf8.ValidString(s) && !utf8.ValidString(got) {
			t.Errorf("Shorten(%q, %d) = %q: invalid UTF-8", s, n, got)
		}
		if utf8.RuneCountInString(got) > n+3 {
			t.Errorf("Shorten(%q, %d) = %q: too long", s, n, got)
		}
	})
}
