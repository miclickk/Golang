package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var nonLetters = regexp.MustCompile(`[^\p{L}]+`)

func topK(counts map[string]int, k int) []string {
	words := make([]string, 0, len(counts))
	for w := range counts {
		words = append(words, w)
	}
	sort.Slice(words, func(i, j int) bool {
		if counts[words[i]] != counts[words[j]] {
			return counts[words[i]] > counts[words[j]]
		}
		return words[i] < words[j]
	})
	if len(words) > k {
		words = words[:k]
	}
	return words
}

func TopWords(text string, k int) []string {
	counts := map[string]int{}
	for _, w := range strings.Fields(text) {
		re := regexp.MustCompile(`[^\p{L}]+`)
		w = strings.ToLower(re.ReplaceAllString(w, ""))
		if w != "" {
			counts[w]++
		}
	}
	return topK(counts, k)
}

func TopWordsV2(text string, k int) []string {
	counts := map[string]int{}
	for _, w := range strings.Fields(text) {
		w = strings.ToLower(nonLetters.ReplaceAllString(w, ""))
		if w != "" {
			counts[w]++
		}
	}
	return topK(counts, k)
}

func TopWordsV3(text string, k int) []string {
	counts := map[string]int{}
	for _, w := range strings.Fields(text) {
		w = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) {
				return unicode.ToLower(r)
			}
			return -1
		}, w)
		if w != "" {
			counts[w]++
		}
	}
	return topK(counts, k)
}

func main() {
	chunk := "Go is fast, simple, and reliable! Привет, мир! привет Go performance profiling benchmark test.\n"
	text := strings.Repeat(chunk, 300)
	_ = os.MkdirAll("testdata", 0o755)
	_ = os.WriteFile(filepath.Join("testdata", "sample.txt"), []byte(text), 0o644)
	fmt.Println(TopWordsV3(text, 5))
}
