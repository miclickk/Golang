package main

import (
	"errors"
	"testing"
)

func TestCreateTaskSuccess(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})

	task, err := manager.CreateTask("MyPass123!")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	if task.State != StatusPending {
		t.Errorf("expected state %s, got %s", StatusPending, task.State)
	}
}

func TestCreateTaskInvalidInput(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"empty password", "", ErrEmptyPassword},
		{"password with spaces", "pass word", ErrPasswordHasSpaces},
		{"only spaces", "   ", ErrPasswordHasSpaces},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.CreateTask(tc.password)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetTaskSuccessAndNotFound(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})

	created, _ := manager.CreateTask("ValidPass1")

	fetched, err := manager.GetTask(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetched.Password != "ValidPass1" {
		t.Errorf("expected password ValidPass1, got %s", fetched.Password)
	}

	_, err = manager.GetTask(999)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
	}
}

func TestProcessTaskAndResults(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})

	tests := []struct {
		name         string
		password     string
		wantLength   int
		wantDigits   int
		wantUpper    int
		wantSpecial  int
		wantStrength string
	}{
		{"weak password", "abcdef", 6, 0, 0, 0, "Weak"},
		{"medium password", "Abcdef12", 8, 2, 1, 0, "Medium"},
		{"strong password", "Abcdef12!", 9, 2, 1, 1, "Strong"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task, err := manager.CreateTask(tc.password)
			if err != nil {
				t.Fatalf("unexpected error on create: %v", err)
			}

			processed, err := manager.ProcessTask(task.ID)
			if err != nil {
				t.Fatalf("unexpected error on process: %v", err)
			}

			if processed.State != StatusCompleted {
				t.Errorf("expected state %s, got %s", StatusCompleted, processed.State)
			}
			if !processed.IsCompleted() {
				t.Errorf("expected IsCompleted() to be true")
			}
			if processed.Result.Length != tc.wantLength {
				t.Errorf("expected length %d, got %d", tc.wantLength, processed.Result.Length)
			}
			if processed.Result.Digits != tc.wantDigits {
				t.Errorf("expected digits %d, got %d", tc.wantDigits, processed.Result.Digits)
			}
			if processed.Result.Upper != tc.wantUpper {
				t.Errorf("expected upper %d, got %d", tc.wantUpper, processed.Result.Upper)
			}
			if processed.Result.Special != tc.wantSpecial {
				t.Errorf("expected special %d, got %d", tc.wantSpecial, processed.Result.Special)
			}
			if processed.Result.Strength != tc.wantStrength {
				t.Errorf("expected strength %s, got %s", tc.wantStrength, processed.Result.Strength)
			}
		})
	}
}

func TestInvalidStateTransition(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})

	task, _ := manager.CreateTask("Secure123!")
	_, err := manager.ProcessTask(task.ID)
	if err != nil {
		t.Fatalf("expected first process to succeed, got %v", err)
	}

	_, err = manager.ProcessTask(task.ID)
	if !errors.Is(err, ErrInvalidState) {
		t.Errorf("expected error %v on second process, got %v", ErrInvalidState, err)
	}
}

func TestMultipleIndependentTasks(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})

	t1, _ := manager.CreateTask("short")
	t2, _ := manager.CreateTask("VeryStrong123!")

	_, _ = manager.ProcessTask(t1.ID)

	fetched1, _ := manager.GetTask(t1.ID)
	fetched2, _ := manager.GetTask(t2.ID)

	if fetched1.State != StatusCompleted || fetched1.Result.Strength != "Weak" {
		t.Errorf("task 1 corrupted: %+v", fetched1)
	}
	if fetched2.State != StatusPending || fetched2.Result.Strength != "" {
		t.Errorf("task 2 should still be pending and empty result: %+v", fetched2)
	}
}
