package main

import (
	"errors"
	"strings"
	"unicode"
)

type Status string

const (
	StatusPending    Status = "Pending"
	StatusProcessing Status = "Processing"
	StatusCompleted  Status = "Completed"
)

var (
	ErrEmptyPassword     = errors.New("password cannot be empty")
	ErrPasswordHasSpaces = errors.New("password cannot contain spaces")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidState      = errors.New("task cannot be processed in its current state")
)

type AnalysisResult struct {
	Length   int
	Digits   int
	Upper    int
	Special  int
	Strength string
}

type Task struct {
	ID       int
	Password string
	State    Status
	Result   AnalysisResult
}

func (t Task) IsCompleted() bool {
	return t.State == StatusCompleted
}

func (t *Task) SetState(newState Status) error {
	if t.State == StatusPending && newState == StatusProcessing {
		t.State = newState
		return nil
	}
	if t.State == StatusProcessing && newState == StatusCompleted {
		t.State = newState
		return nil
	}
	return ErrInvalidState
}

func ValidatePassword(password string) error {
	if len(password) == 0 {
		return ErrEmptyPassword
	}
	if strings.Contains(password, " ") {
		return ErrPasswordHasSpaces
	}
	return nil
}

type Processor interface {
	Process(input string) AnalysisResult
}

type PasswordProcessor struct{}

func (p PasswordProcessor) Process(password string) AnalysisResult {
	res := AnalysisResult{
		Length: len([]rune(password)),
	}

	for _, ch := range password {
		if unicode.IsDigit(ch) {
			res.Digits++
		} else if unicode.IsUpper(ch) {
			res.Upper++
		} else if !unicode.IsLower(ch) {
			res.Special++
		}
	}

	score := 0
	if res.Length >= 8 {
		score++
	}
	if res.Digits > 0 {
		score++
	}
	if res.Upper > 0 {
		score++
	}
	if res.Special > 0 {
		score++
	}

	if score <= 1 {
		res.Strength = "Weak"
	} else if score <= 3 {
		res.Strength = "Medium"
	} else {
		res.Strength = "Strong"
	}

	return res
}
