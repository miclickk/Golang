package main

import (
	"errors"
	"unicode"
)

type TaskState string

const (
	StatusPending    TaskState = "Pending"
	StatusQueued     TaskState = "Queued"
	StatusProcessing TaskState = "Processing"
	StatusCompleted  TaskState = "Completed"
)

var (
	ErrEmptyPassword     = errors.New("password cannot be empty")
	ErrPasswordHasSpaces = errors.New("password cannot contain spaces")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidState      = errors.New("task cannot be processed in its current state")
	ErrQueueFull         = errors.New("task queue is full")
	ErrManagerStopped    = errors.New("task manager is stopped")
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
	State    TaskState
	Result   AnalysisResult
}

func (t Task) IsCompleted() bool {
	return t.State == StatusCompleted
}

func (t *Task) SetState(newState TaskState) error {
	switch newState {
	case StatusQueued:
		if t.State != StatusPending {
			return ErrInvalidState
		}
	case StatusProcessing:
		if t.State != StatusPending && t.State != StatusQueued {
			return ErrInvalidState
		}
	case StatusCompleted:
		if t.State != StatusProcessing {
			return ErrInvalidState
		}
	default:
		return ErrInvalidState
	}
	t.State = newState
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return ErrEmptyPassword
	}
	for _, ch := range password {
		if unicode.IsSpace(ch) {
			return ErrPasswordHasSpaces
		}
	}
	return nil
}

type Processor interface {
	Process(input string) AnalysisResult
}

type PasswordProcessor struct{}

func (PasswordProcessor) Process(password string) AnalysisResult {
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

	switch {
	case score >= 4:
		res.Strength = "Strong"
	case score >= 2:
		res.Strength = "Medium"
	default:
		res.Strength = "Weak"
	}

	return res
}
