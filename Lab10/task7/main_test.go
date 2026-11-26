package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func loadSample() string {
	b, err := os.ReadFile(filepath.Join("testdata", "sample.txt"))
	if err == nil && len(b) > 0 {
		return string(b)
	}
	chunk := "Go is fast, simple, and reliable! Привет, мир! привет Go performance profiling benchmark test.\n"
	return strings.Repeat(chunk, 300)
}

func TestTopWords(t *testing.T) {
	tests := []struct {
		name string
		text string
		k    int
		want []string
	}{
		{"cyrillic", "Привет, мир! привет", 1, []string{"привет"}},
		{"ties alphabetical", "b a c a b", 2, []string{"a", "b"}},
		{"empty", "", 5, []string{}},
		{"punctuation only", "123 !!! ---", 3, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TopWords(tt.text, tt.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTopWordsVersionsEqual(t *testing.T) {
	sample := loadSample()
	v1 := TopWords(sample, 10)
	v2 := TopWordsV2(sample, 10)
	v3 := TopWordsV3(sample, 10)

	if !reflect.DeepEqual(v1, v2) {
		t.Fatalf("v2 = %v, want %v", v2, v1)
	}
	if !reflect.DeepEqual(v1, v3) {
		t.Fatalf("v3 = %v, want %v", v3, v1)
	}
}

func BenchmarkTopWords(b *testing.B) {
	sample := loadSample()
	for b.Loop() {
		TopWords(sample, 10)
	}
}

func BenchmarkTopWordsV2(b *testing.B) {
	sample := loadSample()
	for b.Loop() {
		TopWordsV2(sample, 10)
	}
}

func BenchmarkTopWordsV3(b *testing.B) {
	sample := loadSample()
	for b.Loop() {
		TopWordsV3(sample, 10)
	}
}
