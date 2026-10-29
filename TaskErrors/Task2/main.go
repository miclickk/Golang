package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("student not found")

type GradeError struct {
	Student string
	Grade   int
}

func (e *GradeError) Error() string {
	return fmt.Sprintf("invalid grade %d for %s", e.Grade, e.Student)
}

type GradeBook map[string]int

func (gb GradeBook) Get(name string) (int, error) {
	if g, ok := gb[name]; ok {
		return g, nil
	}
	return 0, ErrNotFound
}

func (gb GradeBook) Set(name string, grade int) error {
	if grade < 0 || grade > 100 {
		return &GradeError{Student: name, Grade: grade}
	}
	gb[name] = grade
	return nil
}

func AddBonus(gb GradeBook, name string, bonus int) error {
	g, err := gb.Get(name)
	if err != nil {
		return fmt.Errorf("add bonus: %w", err)
	}
	if err := gb.Set(name, g+bonus); err != nil {
		return fmt.Errorf("add bonus: %w", err)
	}
	return nil
}

func report(err error) {
	if err == nil {
		fmt.Println("OK")
		return
	}
	var ge *GradeError
	if errors.Is(err, ErrNotFound) {
		fmt.Println("404: not found")
	} else if errors.As(err, &ge) {
		fmt.Printf("400: invalid grade %d for %s\n", ge.Grade, ge.Student)
	} else {
		fmt.Println("500:", err)
	}
}
