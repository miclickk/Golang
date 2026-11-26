package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrBadFormat = errors.New("bad format")

type Record struct {
	Name  string
	Grade int
}

func parseLine(line string) (Record, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		return Record{}, ErrBadFormat
	}
	name := strings.TrimSpace(parts[0])
	if name == "" {
		return Record{}, ErrBadFormat
	}
	grade, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Record{}, err
	}
	return Record{Name: name, Grade: grade}, nil
}

func ParseAll(lines []string) ([]Record, error) {
	records := make([]Record, 0, len(lines))
	for _, line := range lines {
		rec, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func main() {
	records, err := ParseAll([]string{"Alice, 95", "Bob, 82"})
	fmt.Println(records, err)
}
