package main

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Record
		wantErr error
	}{
		{"valid", "Alice, 95", Record{Name: "Alice", Grade: 95}, nil},
		{"no comma", "Alice 95", Record{}, ErrBadFormat},
		{"extra comma", "Alice, 95, 1", Record{}, ErrBadFormat},
		{"empty name", " , 95", Record{}, ErrBadFormat},
		{"syntax error", "Alice, abc", Record{}, strconv.ErrSyntax},
		{"range error", "Alice, 999999999999999999999", Record{}, strconv.ErrRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLine(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("parseLine(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("parseLine(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseAll(t *testing.T) {
	got, err := ParseAll([]string{"Alice, 95", "Bob, 80"})
	if err != nil {
		t.Fatalf("ParseAll returned error: %v", err)
	}
	if len(got) != 2 || got[0].Name != "Alice" || got[1].Grade != 80 {
		t.Errorf("ParseAll() = %+v", got)
	}

	if _, err := ParseAll([]string{"Alice, 95", "bad"}); !errors.Is(err, ErrBadFormat) {
		t.Errorf("ParseAll(bad) error = %v, want %v", err, ErrBadFormat)
	}
}
