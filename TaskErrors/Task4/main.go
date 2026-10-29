package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrBadFormat = errors.New("missing comma")

type Record struct {
	Name  string
	Grade int
}

func ParseAll(input string) ([]Record, error) {
	var records []Record
	var errs []error
	lines := strings.Split(strings.TrimSpace(input), "\n")

	for i, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			errs = append(errs, fmt.Errorf("line %d: %w", i+1, ErrBadFormat))
			continue
		}

		grade, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			errs = append(errs, fmt.Errorf("line %d: %w", i+1, err))
			continue
		}

		records = append(records, Record{
			Name:  strings.TrimSpace(parts[0]),
			Grade: grade,
		})
	}
	return records, errors.Join(errs...)
}
